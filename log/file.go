package log

import (
	"errors"
	"fmt"
	gstr "goutil/strings"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// 最小保存的天数
	minKeepDay = 24 * time.Hour
)

var (
	errFileClosed = errors.New("file has been closed")
)

// FileConfig 是 NewFile 的参数
type FileConfig struct {
	// 日志保存的根目录
	RootDir string `json:"rootDir" yaml:"rootDir" validate:"required,dirpath"`
	// 日期目录名称格式，默认是 20060102
	DirNameFormat string `json:"dirNameFormat" yaml:"dirNameFormat"`
	// 文件名称格式，默认是 150405.000000
	FileNameFormat string `json:"fileNameFormat" yaml:"fileNameFormat"`
	// 每一份日志文件的最大字节，单位，k/m/g/t
	FileMaxSize string `json:"fileMaxSize" yaml:"fileMaxSize" validate:"required"`
	// 保存的最大天数，最小是 1 天
	MaxKeepDay int `json:"maxKeepDay" yaml:"maxKeepDay" validate:"required,min=1"`
	// 内存输出到磁盘的间隔，最小是 1 毫秒，设置太小没有意义
	// 但是如果文件大小达到 FileMaxSize ，那么立即输出
	FlushInterval time.Duration `json:"flushInterval" yaml:"flushInterval" validate:"required,min=1000000"`
	// 检查过期文件的间隔，最小是 1 秒
	CheckExpireInterval time.Duration `json:"checkExpireInterval" yaml:"checkExpireInterval" validate:"required,min=1000000000"`
	// 输出到文件的同时，是否输出到控制台，out/err
	Std string `json:"std" yaml:"std" validate:"omitempty,oneof=out err"`
}

// NewFile 返回一个 File 实例
func NewFile(conf *FileConfig) (*File, error) {
	fileMaxSize, err := gstr.StringToByte(conf.FileMaxSize)
	if err != nil {
		return nil, errors.New("error file max size format")
	}
	// 实例
	f := new(File)
	f.rootDir = conf.RootDir
	f.maxFileSize = int64(fileMaxSize)
	f.exit = make(chan struct{})
	f.maxKeepDuraion = time.Duration(conf.MaxKeepDay) * minKeepDay
	switch conf.Std {
	case "err":
		f.std = os.Stderr
	case "out":
		f.std = os.Stdout
	}
	f.dirNameFormat = conf.DirNameFormat
	f.fileNameFormat = conf.FileNameFormat
	if f.dirNameFormat == "" {
		f.dirNameFormat = "20060102"
	}
	if f.fileNameFormat == "" {
		f.fileNameFormat = "20060102150405.999999"
	}
	y, m, d := time.Now().Date()
	f.dateY, f.dateM, f.dateD = y, int(m), d
	// 输出协程
	f.wait.Add(1)
	go f.syncRoutine(conf.FlushInterval)
	// 检查过期协程
	f.wait.Add(1)
	go f.checkExpireRoutine(conf.CheckExpireInterval)
	//
	return f, nil
}

// File 实现了 io.Writer 接口，可以作为 Logger 的输出
// 首先将 log 缓存在内存中，每隔一段时间，
// 或者内存的数据到了指定的字节，才将数据输出到磁盘，提高性能
// 日志目录结构是，root/date/time.ms
// 还会自动删除磁盘上时间超过指定天数的文件
type File struct {
	lock sync.Mutex
	wait sync.WaitGroup
	// 退出协程通知
	exit chan struct{}
	// 是否已关闭标志
	closed bool
	// 日志文件的根目录
	rootDir string
	// 内存数据
	data []byte
	// 当前打开的文件
	file *os.File
	// 最大的保存天数
	maxKeepDuraion time.Duration
	// 当前磁盘文件的字节
	curFileSize int64
	// 磁盘文件的最大字节
	maxFileSize int64
	// 控制台输出
	std io.Writer
	// 目录格式
	dirNameFormat string
	// 文件格式
	fileNameFormat string
	// 日期，用于检查日期更换
	dateY, dateM, dateD int
}

// Write 是 io.Writer 接口
func (f *File) Write(b []byte) (int, error) {
	f.lock.Lock()
	// 关闭了
	if f.closed {
		f.lock.Unlock()
		return 0, errFileClosed
	}
	// 内存
	f.data = append(f.data, b...)
	f.curFileSize += int64(len(b))
	// 内存数据达到最大了
	if f.curFileSize >= f.maxFileSize {
		f.curFileSize = 0
		// 保存
		f.flushFile()
		f.closeFile()
		f.openFile()
	}
	f.lock.Unlock()
	// 输出控制台
	if f.std != nil {
		_, _ = f.std.Write(b)
	}
	//
	return len(b), nil
}

// checkExpireRoutine 检查过期文件
func (f *File) checkExpireRoutine(dur time.Duration) {
	// 计时器
	timer := time.NewTicker(dur)
	defer func() {
		// 异常
		_ = recover()
		// 计时器
		timer.Stop()
		// 结束
		f.wait.Done()
	}()
	for !f.closed {
		select {
		case now := <-timer.C:
			// 检查过期
			f.checkExpire(&now)
			// 检查日期更换
			f.checkDate(&now)
			// 计时器
			timer.Reset(dur)
		case <-f.exit:
			// 退出信号
			return
		}
	}
}

// syncRoutine 输出磁盘
func (f *File) syncRoutine(dur time.Duration) {
	// 计时器
	timer := time.NewTicker(dur)
	defer func() {
		// 异常
		_ = recover()
		// 计时器
		timer.Stop()
		// 剩余的数据
		f.lock.Lock()
		f.flushFile()
		f.closeFile()
		f.lock.Unlock()
		// 结束
		f.wait.Done()
	}()
	// 先打开文件准备
	f.openLastFile()
	for !f.closed {
		select {
		case <-timer.C:
			// 保存
			f.lock.Lock()
			f.flushFile()
			f.lock.Unlock()
			// 计时器
			timer.Reset(dur)
		case <-f.exit:
			// 退出信号
			return
		}
	}
}

// close 设置标记，然后通知
func (f *File) close() error {
	f.lock.Lock()
	// 已经关闭
	if f.closed {
		f.lock.Unlock()
		return errFileClosed
	}
	f.closed = true
	f.lock.Unlock()
	// 结束协程通知
	close(f.exit)
	//
	return nil
}

// Close 实现 io.Closer 接口，同步内存到磁盘，等待协程退出
func (f *File) Close() error {
	err := f.close()
	if err != nil {
		return err
	}
	// 等待退出
	f.wait.Wait()
	//
	return nil
}

// checkExpire 检查过期文件
func (f *File) checkExpire(now *time.Time) {
	// 读取根目录下的所有文件
	dirEntries, err := os.ReadDir(f.rootDir)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// 应该删除的时间
	delTime := now.Add(-f.maxKeepDuraion)
	// 循环检查
	for i := 0; i < len(dirEntries); i++ {
		entry := dirEntries[i]
		fi, err := entry.Info()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		// 文件时间小于删除时间
		if fi.ModTime().Sub(delTime) < 0 {
			err = os.RemoveAll(filepath.Join(f.rootDir, fi.Name()))
			if nil != err {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

// checkDate 检查日期更换
func (f *File) checkDate(now *time.Time) {
	y, m, d := now.Date()
	if y != f.dateY || int(m) != f.dateM || d != f.dateD {
		// 保存
		f.lock.Lock()
		f.flushFile()
		f.closeFile()
		f.openFile()
		f.lock.Unlock()
	}
	f.dateY, f.dateM, f.dateD = y, int(m), d
}

// flush 将内存的数据保存到硬盘，如果写入失败，数据会丢弃
func (f *File) flushFile() {
	if len(f.data) > 0 {
		_, err := f.file.Write(f.data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		f.data = f.data[:0]
	}
}

// open 打开一个新的文件
func (f *File) openFile() {
	now := time.Now()
	// 创建目录，root/date
	dateDir := filepath.Join(f.rootDir, now.Format(f.dirNameFormat))
	err := os.MkdirAll(dateDir, os.ModePerm)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// 创建日志文件，root/date/time.ms
	timeFile := filepath.Join(dateDir, now.Format(f.fileNameFormat)+".log")
	f.file, err = os.OpenFile(timeFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
	}
}

// closeFile 关闭当前文件
func (f *File) closeFile() {
	if nil != f.file {
		f.file.Close()
		f.file = nil
	}
}

// openLastFile 打开上一个最新的文件
func (f *File) openLastFile() {
	now := time.Now()
	// 创建目录，root/date
	dateDir := filepath.Join(f.rootDir, now.Format(f.dirNameFormat))
	err := os.MkdirAll(dateDir, os.ModePerm)
	if nil != err {
		panic(err)
	}
	// 读取根目录下的所有文件
	dirEntries, err := os.ReadDir(dateDir)
	if nil != err {
		panic(err)
	}
	// 没有文件
	fileName := now.Format(f.fileNameFormat)
	if len(dirEntries) > 0 {
		// 循环检查
		dirEntry := dirEntries[0]
		lastFI, err := dirEntry.Info()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			lastTime := lastFI.ModTime()
			// 找出最新的文件时间
			for i := 1; i < len(dirEntries); i++ {
				dirEntry := dirEntries[i]
				fi, err := dirEntry.Info()
				if err != nil {
					panic(err)
				}
				m := fi.ModTime()
				if m.After(lastTime) {
					lastTime = m
					lastFI = fi
				}
			}
			// 最新的大小
			if lastFI.Size() < f.maxFileSize {
				fileName = lastFI.Name()
			}
		}
	}
	// 创建日志文件，root/date/time.ms
	timeFile := filepath.Join(dateDir, fileName+".log")
	f.file, err = os.OpenFile(timeFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if nil != err {
		panic(err)
	}
	fi, err := f.file.Stat()
	if nil != err {
		panic(err)
	}
	f.curFileSize = fi.Size()
}
