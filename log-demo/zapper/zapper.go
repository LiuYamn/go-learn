package zapper

import (
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	DebugLevel  = zapcore.DebugLevel
	InfoLevel   = zapcore.InfoLevel
	WarnLevel   = zapcore.WarnLevel
	ErrorLevel  = zapcore.ErrorLevel
	DPanicLevel = zapcore.DPanicLevel
	PanicLevel  = zapcore.PanicLevel
	FatalLevel  = zapcore.FatalLevel
)

type Logger struct {
	z zap.Logger
}

func NewLogger(logFolderName, logFileName string, level zapcore.Level) Logger {
	zapLog := newZapLog(logFolderName, logFileName, level)
	logger := Logger{z: zapLog}
	return logger
}

func newZapLog(logFolderName, logFileName string, level zapcore.Level) zap.Logger {
	err := checkDir("./" + logFolderName)
	if err != nil {
		log.Fatal(err)
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	CurrLevel := zap.NewAtomicLevelAt(level)

	customCfg := zap.Config{
		Level:            CurrLevel,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr", logFolderName + "/" + logFileName + ".log"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := customCfg.Build()
	return *logger
}

func (l *Logger) Debug(key string, value interface{}) {
	l.z.Debug("", zap.String("c", getRaw()), zap.Any(key, value))
}

func (l *Logger) Info(key string, value interface{}) {
	l.z.Info("", zap.String("c", getRaw()), zap.Any(key, value))
}

func (l *Logger) Warn(key string, value interface{}) {
	l.z.Warn("", zap.String("c", getRaw()), zap.Any(key, value))
}

func (l *Logger) Error(key string, value interface{}) {
	l.z.Error("", zap.String("c", getRaw()), zap.Any(key, value))
}

func (l *Logger) Fatal(key string, value interface{}) {
	l.z.Fatal("", zap.String("c", getRaw()), zap.Any(key, value))
}

// 日志所需的时间格式化
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// 检查目录,不存在则创建
func checkDir(path string) error {
	isExists := exists(path)
	if isExists {
		return nil
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil

}

//Exists 判断所给路径文件/文件夹是否存在
func exists(onePath string) bool {
	_, err := os.Stat(onePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func getRaw() (tmpString string) {
	_, filePath, line, ok := runtime.Caller(2)
	tmp := strings.Split(filePath, "/")
	filePath = tmp[len(tmp)-2] + "/" + tmp[len(tmp)-1]
	if ok {
		tmpString = filePath + ":" + strconv.Itoa(line)
	}
	return tmpString
}
