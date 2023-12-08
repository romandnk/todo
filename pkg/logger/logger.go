package logger

//go:generate mockgen -source=logger.go -destination=mock/mock.go logger

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
}
