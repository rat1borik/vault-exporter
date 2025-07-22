package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"vault-exporter/internal/utils"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Info(args ...any)
	Error(args ...any)
	Debug(args ...any)
	Writer() io.Writer
	WithField(key string, value any) Logger
}

type LogrusLogger struct {
	entry *logrus.Entry
}

func NewLogrusLogger(path string) *LogrusLogger {

	if path == "" {
		path, _ = utils.ExecPath("logs/vault-agent.log")
	}

	// Настройка ротации логов
	logFile := &lumberjack.Logger{
		Filename:   path, // лог будет автоматически ротироваться
		MaxSize:    10,   // мегабайты
		MaxBackups: 7,    // количество резервных файлов
		MaxAge:     14,   // хранить 14 дней
		Compress:   true, // сжимать старые логи .gz
	}

	log := logrus.New()
	log.SetFormatter(&CustomFormatter{})
	file, _ := os.OpenFile(filepath.Join(path, "1.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	multiWriter := io.MultiWriter(logFile, file, os.Stdout) // stdout в конце, чтобы не валил всю цепочку, когда служба
	log.SetOutput(multiWriter)

	return &LogrusLogger{entry: logrus.NewEntry(log)}
}

func (l *LogrusLogger) Info(args ...any)  { l.entry.Info(args...) }
func (l *LogrusLogger) Error(args ...any) { l.entry.Error(args...) }
func (l *LogrusLogger) Debug(args ...any) { l.entry.Debug(args...) }

func (l *LogrusLogger) WithField(key string, value any) Logger {
	return &LogrusLogger{entry: l.entry.WithField(key, value)}
}

func (l *LogrusLogger) Writer() io.Writer {
	return l.entry.Logger.Writer()
}

// CustomFormatter определяет собственный формат логов
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer
	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	// Пример: [INFO] 2025-07-21 14:33:00 - msg
	fmt.Fprintf(&b, "[%s %s] %s\n", timestamp, entry.Level.String(), entry.Message)

	return b.Bytes(), nil
}
