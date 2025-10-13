package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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
	module string
	file   *os.File
}

// LoggerConfig 日志配置结构
type LoggerConfig struct {
	Module   string       // 模块名称
	FilePath string       // 文件路径，为空则输出到标准输出
	Level    logrus.Level // 日志级别
}

// NewLogger 创建新的日志实例
func NewLogger(config LoggerConfig) (*Logger, error) {
	logger := logrus.New()
	logger.SetFormatter(&CustomFormatter{Module: config.Module})
	logger.SetLevel(config.Level)

	var logFile *os.File
	var err error

	// 如果指定了文件路径，则输出到文件
	if config.FilePath != "" {
		// 确保目录存在
		dir := getDir(config.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建日志目录失败: %v", err)
		}

		// 打开或创建日志文件
		logFile, err = os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("打开日志文件失败: %v", err)
		}
		logger.SetOutput(logFile)
	} else {
		// 输出到标准输出
		logger.SetOutput(os.Stdout)
	}

	return &Logger{
		Logger: logger,
		module: config.Module,
		file:   logFile,
	}, nil
}

// getDir 从文件路径中提取目录
func getDir(filePath string) string {
	if idx := strings.LastIndex(filePath, "/"); idx != -1 {
		return filePath[:idx]
	}
	if idx := strings.LastIndex(filePath, "\\"); idx != -1 {
		return filePath[:idx]
	}
	return "."
}

// WithModule 创建子模块日志器
func (l *Logger) WithModule(subModule string) *Logger {
	fullModule := l.module + "." + subModule

	// 创建新的日志实例，继承输出配置
	newLogger := &Logger{
		Logger: logrus.New(),
		module: fullModule,
		file:   l.file, // 继承文件句柄
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

// SetOutputFile 动态设置输出文件
func (l *Logger) SetOutputFile(filePath string) error {
	// 关闭旧文件（如果存在）
	if l.file != nil {
		l.file.Close()
	}

	// 创建目录
	dir := getDir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 打开新文件
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	l.file = logFile
	l.Logger.SetOutput(logFile)
	return nil
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
		Module:   module,
		FilePath: "",
		Level:    logrus.InfoLevel,
	})
	return logger
}
