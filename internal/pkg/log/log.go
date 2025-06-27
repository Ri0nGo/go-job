package log

import (
	"github.com/natefinch/lumberjack"
	"log/slog"
	"os"
)

func InitSlog(level slog.Level, source bool) {
	_ = os.MkdirAll("./logs", 0755)

	writer := &lumberjack.Logger{
		Filename:   "./logs/app.log", // 日志文件路径
		MaxSize:    100,              // 单个日志文件最大100MB
		MaxBackups: 10,               // 最多保留5个旧文件
		MaxAge:     30,               // 日志保留30天
		Compress:   true,             // 是否压缩归档日志
	}
	logger := slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level:     level,
		AddSource: source,
	}))

	slog.SetDefault(logger)
}
