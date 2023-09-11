package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// 默认的日志文件最大字节
	defaultMaxFileSize = 1024 * 1024 * 10
	// 最小保存的天数
	minKeepDay = 24 * time.Hour
	// 最小同步间隔
	minSyncDur = 100 * time.Millisecond
	// 目录格式
	dirNameFormat = "20060102"
	// 文件格式
	fileNameFormat = "20060102150405.000000"
)

var (
	errFileClosed = errors.New("file has been closed")
)

// FileConfig 是 NewFile 的参数。
type FileConfig struct {
	// 日志保存的根目录
	RootDir string `json:"rootDir" yaml:"rootDir" validate:"required"`
	// 每一份日志文件的最大字节，使用 1.5/K/M/G/T 这样的字符表示
	MaxFileSize string `json:"maxFileSize" yaml:"maxFileSize"`
	// 保存的最大天数，最小是1天
	MaxKeepDay int `json:"maxKeepDay" yaml:"maxKeepDay" validate:"required,min=1"`
	// 同步到磁盘的时间间隔，单位秒，最小是 1
	// 但是如果文件大小达到 MaxFileSize ，那么立即同步
	SyncInterval int `json:"syncInterval" yaml:"syncInterval" validate:"required,min=1"`
	// 是否输出到控制台，out/err
	Std string `json:"std" yaml:"std" validate:"omitempty,oneof=out err"`
}

// NewFile 返回一个 File 实例。
func NewFile(conf *FileConfig) (*File, error) {
	// 解析文件大小
	size, err := ParseSize(conf.MaxFileSize)
	if err != nil {
		return nil, err
	}
	if size < 1 {
		size = defaultMaxFileSize
	}
	// 文件保存的时长
	keepDuraion := time.Duration(conf.MaxKeepDay) * minKeepDay
	if keepDuraion < minKeepDay {
		keepDuraion = minKeepDay
	}
	// 同步时长
	syncDur := time.Duration(conf.SyncInterval) * time.Second
	if syncDur < minSyncDur {
		syncDur = minSyncDur
	}
	// 实例
	f := new(File)
	f.rootDir = conf.RootDir
	f.maxFileSize = int(size)
	f.exit = make(chan struct{})
	f.maxKeepDuraion = keepDuraion
	switch conf.Std {
	case "err":
		f.std = os.Stderr
	case "out":
		f.std = os.Stdout
	}
	// 先打开文件准备
	f.openLast()
	// 启动同步协程
	f.wait.Add(1)
	go f.syncLoop(syncDur)
	return f, nil
}

// File 实现了 io.Writer 接口，可以作为 Logger 的输出。
// File 首先会将 log 保存在内存中，后台启动一个同步协程，每隔一段时间将数据同步到磁盘。
// 如果内存的数据到了最大，会立即同步。
// 在同步的同时，File 还会自动删除磁盘上时间超过指定天数的文件。
// 目录结构是，root/date/time.ms
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
	curFileSize int
	// 磁盘文件的最大字节
	maxFileSize int
	// 控制台输出
	std io.Writer
}

// Write 是 io.Writer 接口。
func (f *File) Write(b []byte) (int, error) {
	f.lock.Lock()
	// 关闭了
	if f.closed {
		f.lock.Unlock()
		return 0, errFileClosed
	}
	// 添加到内存
	f.data = append(f.data, b...)
	f.curFileSize += len(b)
	// 如果内存数据达到最大了，换新文件输出
	if f.curFileSize >= f.maxFileSize {
		f.curFileSize = 0
		f.flush()
		f.close()
		f.open()
	}
	f.lock.Unlock()
	if f.std != nil {
		f.std.Write(b)
	}
	return len(b), nil
}

// syncLoop 运行在一个协程中。
func (f *File) syncLoop(syncDur time.Duration) {
	syncTimer := time.NewTicker(syncDur)
	defer func() {
		syncTimer.Stop()
		f.lock.Lock()
		f.flush()
		f.close()
		f.lock.Unlock()
		f.wait.Done()
	}()
	checkTime := time.Now()
	// 程序退出
	quit := make(chan os.Signal, 1)
	// 先检查一次过期
	f.check(&checkTime)
	for !f.closed {
		select {
		case now := <-syncTimer.C:
			// 检查过期
			if now.Sub(checkTime) > time.Hour {
				f.check(&checkTime)
				checkTime = now
			}
			// 同步时间
			f.lock.Lock()
			f.flush()
			f.lock.Unlock()
			// 计时器
			syncTimer.Reset(syncDur)
		case <-f.exit:
			// 退出信号
			return
		case <-quit:
			return
		}
	}
}

// Close 实现 io.Closer 接口，同步内存到磁盘，等待协程退出。
func (f *File) Close() error {
	f.lock.Lock()
	if f.closed {
		f.lock.Unlock()
		return errFileClosed
	}
	f.closed = true
	f.lock.Unlock()
	// 结束协程通知。
	close(f.exit)
	// 等待退出。
	f.wait.Wait()
	// 同步数据，并关闭文件。
	f.flush()
	f.close()
	// 返回
	return nil
}

// check 检查过期文件。
func (f *File) check(now *time.Time) {
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

// flush 将内存的数据保存到硬盘，如果写入失败，数据会丢弃。
func (f *File) flush() {
	_, err := f.file.Write(f.data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	f.data = f.data[:0]
}

// open 打开一个新的文件
func (f *File) open() {
	now := time.Now()
	// 创建目录，root/date
	dateDir := filepath.Join(f.rootDir, now.Format(dirNameFormat))
	err := os.MkdirAll(dateDir, os.ModePerm)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// 创建日志文件，root/date/time.ms
	timeFile := filepath.Join(dateDir, now.Format(fileNameFormat))
	f.file, err = os.OpenFile(timeFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
	}
}

// close 关闭当前文件
func (f *File) close() {
	if nil != f.file {
		f.file.Close()
		f.file = nil
	}
}

// openLast 打开上一个最新的文件
func (f *File) openLast() {
	now := time.Now()
	// 创建目录，root/date
	dateDir := filepath.Join(f.rootDir, now.Format(dirNameFormat))
	err := os.MkdirAll(dateDir, os.ModePerm)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// 读取根目录下的所有文件
	dirEntries, err := os.ReadDir(dateDir)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	// 没有文件
	fileName := now.Format(fileNameFormat)
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
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				m := fi.ModTime()
				if m.After(lastTime) {
					lastTime = m
					lastFI = fi
				}
			}
			// 最新的大小
			if lastFI.Size() < int64(f.maxFileSize) {
				fileName = lastFI.Name()
			}
		}
	}
	// 创建日志文件，root/date/time.ms
	timeFile := filepath.Join(dateDir, fileName)
	f.file, err = os.OpenFile(timeFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
	}
	fi, err := f.file.Stat()
	if nil != err {
		fmt.Fprintln(os.Stderr, err)
	}
	f.curFileSize = int(fi.Size())
}
