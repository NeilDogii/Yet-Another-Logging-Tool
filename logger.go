package main

type Logger struct {
	Level string
	logs  []string
}

func NewLogger(level string) *Logger {
	return &Logger{
		Level: level,
		logs:  []string{},
	}
}

func (l *Logger) Log(message string) {
	l.logs = append(l.logs, message)
}

func (l *Logger) GetLogs() []string {
	return l.logs
}
