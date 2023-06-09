package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogTypeDaily = "daily"
)

var (
	L    *Logger
	once sync.Once
)

type Conf struct {
	Filename  string `mapstructure:"filename"`
	Level     string `mapstructure:"level"`
	Type      string `mapstructure:"type"`
	MaxSize   int    `mapstructure:"max_size"`
	MaxBackup int    `mapstructure:"max_backup"`
	MaxAge    int    `mapstructure:"max_age"`
	Compress  bool   `mapstructure:"compress"`
	IsLocal   bool   `mapstructure:"is_local"`
}

type Logger struct {
	*zap.Logger
	conf *Conf
}

func (l *Logger) getLogWriter() zapcore.WriteSyncer {
	if l.conf.Type == LogTypeDaily {
		logName := time.Now().Format("2006-01-02.log")
		l.conf.Filename = strings.ReplaceAll(l.conf.Filename, "logs.log", logName)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   l.conf.Filename,
		MaxSize:    l.conf.MaxSize,
		MaxBackups: l.conf.MaxBackup,
		MaxAge:     l.conf.MaxAge,
		Compress:   l.conf.Compress,
	}

	if l.conf.IsLocal {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}

func (l *Logger) getEncoder() zapcore.Encoder {
	encoderConf := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "Logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     l.customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if l.conf.IsLocal {
		encoderConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConf)
	}
	return zapcore.NewJSONEncoder(encoderConf)
}

func (l *Logger) customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func (l *Logger) setConf(conf *Conf) *Logger {
	l.conf = conf
	return l
}

func (l *Logger) build() *Logger {
	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(l.conf.Level)); err != nil {
		panic("日志级别不正确")
	}

	core := zapcore.NewCore(l.getEncoder(), l.getLogWriter(), level)

	l.Logger = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	zap.ReplaceGlobals(l.Logger)
	return l
}

func New(conf *Conf) *Logger {
	return new(Logger).setConf(conf).build()
}
