package logger

type Interface interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

type Logger struct{}

func New() Interface {
	return &Logger{}
}

func (l *Logger) Info(msg string, args ...interface{}) {
}

func (l *Logger) Error(msg string, args ...interface{}) {
}

func (l *Logger) Debug(msg string, args ...interface{}) {
}