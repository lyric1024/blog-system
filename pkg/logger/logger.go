package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(levelStr, outputFile string) {
	var level zapcore.Level
	switch levelStr {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	// 编码器配置（带颜色、时间、调用位置）
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 控制台彩色日志
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// var cores []zapcore.Core

	// 控制台输出
	// cores = append(cores, zapcore.NewCore(
	// 	consoleEncoder,
	// 	zapcore.Lock(os.Stdout),
	// 	level,
	// ))

	// 开启开发模式（自动加调用栈、行号）
	core := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level)
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// 提供全局访问方法
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Sync() {
	_ = log.Sync()
}