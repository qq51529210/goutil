# 日志轮子
模仿标准库的 log 造了一个轮子，增加了 Level 和 TraceID 。
# 日志头
提供自定义日志头接口，默认有三种日志头格式可以使用。  
appID 和 traceID 如果为空串不输出，也可以实现 HeaderFormater 自定义日志头格式。
- defalut 输出格式：level appID traceID log
- FileNameStackHeaderFormater 输出格式：level appID traceID fileName:fileLine log
- FilePathStackHeaderFormater 输出格式：level appID traceID filepath:fileLine log

# 输出
默认 Logger 是输出到 os.Stdout ，可以自己指定 io.Writer 。[file.go](./file.go) 实现了输出到文件。

# usage
看 [logger_test.go](./logger_test.go) 文件。

# benchmark

```go
goos: darwin
goarch: amd64
pkg: github.com/qq51529210/log
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
Benchmark_My_Logger-4            1039316              1083 ns/op             216 B/op          2 allocs/op
Benchmark_Std_Logger-4           1100016              1086 ns/op             240 B/op          4 allocs/op
Benchmark_My_Logger_f-4          1065927              1158 ns/op             216 B/op          2 allocs/op
Benchmark_Std_Logger_f-4          878757              1160 ns/op             248 B/op          4 allocs/op
PASS
ok      github.com/qq51529210/log       7.449s
```

