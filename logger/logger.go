package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	maxFileSize = 10 * 1024 * 1024 // 10MB
)

// CustomFormatter 自定义日志格式器
type CustomFormatter struct {
	Module string // 模块名称
}

// Format 实现 logrus.Formatter 接口
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 获取进程ID
	pid := os.Getpid()

	// 格式化时间
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 构建日志行
	logLine := fmt.Sprintf("[%s] [%d] [%s]: \"%s\"\n",
		timestamp,
		pid,
		f.Module,
		entry.Message,
	)

	return []byte(logLine), nil
}

// Logger 包装logrus.Logger，提供模块化支持和文件输出
type Logger struct {
	*logrus.Logger
	module   string
	file     *os.File
	logDir   string
	filePath string
}

// LoggerConfig 日志配置结构
type LoggerConfig struct {
	Module string       // 模块名称
	LogDir string       // 日志目录，为空则输出到标准输出
	Level  logrus.Level // 日志级别
}

// NewLogger 创建新的日志实例
func NewLogger(config LoggerConfig) (*Logger, error) {
	logger := logrus.New()
	logger.SetFormatter(&CustomFormatter{Module: config.Module})
	logger.SetLevel(config.Level)

	var logFile *os.File
	var filePath string
	var err error

	// 如果指定了文件路径，则输出到文件
	if config.LogDir != "" {
		// 确保目录存在
		if err := os.MkdirAll(config.LogDir, 0755); err != nil {
			return nil, fmt.Errorf("创建日志目录失败: %v", err)
		}

		// 获取或创建日志文件
		logFile, filePath, err = getOrCreateLogFile(config.LogDir)
		if err != nil {
			return nil, fmt.Errorf("获取日志文件失败: %v", err)
		}
		logger.SetOutput(logFile)
	} else {
		// 输出到标准输出
		logger.SetOutput(os.Stdout)
	}

	return &Logger{
		Logger:   logger,
		module:   config.Module,
		file:     logFile,
		logDir:   config.LogDir,
		filePath: filePath,
	}, nil
}

// getOrCreateLogFile 获取或创建日志文件
func getOrCreateLogFile(logDir string) (*os.File, string, error) {
	// 确保目录路径以分隔符结尾
	if !strings.HasSuffix(logDir, string(filepath.Separator)) {
		logDir += string(filepath.Separator)
	}

	// 查找目录中现有的日志文件
	files, err := findLogFiles(logDir)
	if err != nil {
		return nil, "", err
	}

	// 如果有日志文件，检查最新的文件大小
	if len(files) > 0 {
		latestFile := files[len(files)-1]
		fileInfo, err := os.Stat(latestFile)
		if err == nil && fileInfo.Size() < maxFileSize {
			// 文件未满，继续使用
			file, err := os.OpenFile(latestFile, os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return nil, "", err
			}
			return file, latestFile, nil
		}
	}

	// 创建新日志文件
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s.log", timestamp)
	filePath := filepath.Join(logDir, filename)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, "", err
	}

	return file, filePath, nil
}

// findLogFiles 查找目录中的日志文件并按时间排序
func findLogFiles(logDir string) ([]string, error) {
	var logFiles []string

	entries, err := os.ReadDir(logDir)
	if err != nil {
		// 如果目录不存在，返回空列表
		if os.IsNotExist(err) {
			return logFiles, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".log") {
			fullPath := filepath.Join(logDir, entry.Name())
			logFiles = append(logFiles, fullPath)
		}
	}

	// 按文件名（时间戳）排序
	sort.Strings(logFiles)
	return logFiles, nil
}

// checkAndRotateFile 检查当前文件大小，如果需要则轮转文件
func (l *Logger) checkAndRotateFile() error {
	if l.file == nil || l.logDir == "" {
		return nil
	}

	// 获取当前文件信息
	fileInfo, err := os.Stat(l.filePath)
	if err != nil {
		return err
	}

	// 如果文件超过最大大小，创建新文件
	if fileInfo.Size() >= maxFileSize {
		l.file.Close()

		newFile, newFilePath, err := getOrCreateLogFile(l.logDir)
		if err != nil {
			return err
		}

		l.file = newFile
		l.filePath = newFilePath
		l.Logger.SetOutput(newFile)
	}

	return nil
}

// 重写日志输出方法，在写入前检查文件大小
func (l *Logger) Log(level logrus.Level, args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Log(level, args...)
}

func (l *Logger) Logf(level logrus.Level, format string, args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Logf(level, format, args...)
}

func (l *Logger) Logln(level logrus.Level, args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Logln(level, args...)
}

// 实现各个级别的日志方法
func (l *Logger) Debug(args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Error(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Fatal(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Panic(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Errorf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Fatalf(format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Panicf(format, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Debugln(args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Infoln(args...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Warnln(args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Errorln(args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Fatalln(args...)
}

func (l *Logger) Panicln(args ...interface{}) {
	l.checkAndRotateFile()
	l.Logger.Panicln(args...)
}

// WithModule 创建子模块日志器
func (l *Logger) WithModule(subModule string) *Logger {
	fullModule := l.module + "." + subModule

	// 创建新的日志实例，继承输出配置
	newLogger := &Logger{
		Logger:   logrus.New(),
		module:   fullModule,
		file:     l.file,
		logDir:   l.logDir,
		filePath: l.filePath,
	}

	newLogger.SetFormatter(&CustomFormatter{Module: fullModule})
	newLogger.SetLevel(l.GetLevel())

	// 继承输出目标
	if l.file != nil {
		newLogger.SetOutput(l.file)
	} else {
		newLogger.SetOutput(os.Stdout)
	}

	return newLogger
}

// AddOutput 添加额外的输出目标
func (l *Logger) AddOutput(writer io.Writer) {
	// 如果当前输出是文件，创建多写器
	if l.file != nil {
		l.Logger.SetOutput(io.MultiWriter(l.file, writer))
	} else {
		l.Logger.SetOutput(io.MultiWriter(os.Stdout, writer))
	}
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// SimpleLogger 快速创建简单日志器（输出到控制台）
func SimpleLogger(module string) *Logger {
	logger, _ := NewLogger(LoggerConfig{
		Module: module,
		LogDir: "",
		Level:  logrus.InfoLevel,
	})
	return logger
}
