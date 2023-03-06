package logger

import "fmt"

type FakeLogger struct{}

func (l FakeLogger) Debugf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (l FakeLogger) Debug(args ...interface{}) {
	fmt.Println(args...)
}

func (l FakeLogger) Warnf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (l FakeLogger) Warn(args ...interface{}) {
	fmt.Println(args...)
}

func (l FakeLogger) Errorf(s string, i ...interface{}) {
	fmt.Println(fmt.Sprintf(s, i...))
}

func (l FakeLogger) Error(i ...interface{}) {
	fmt.Println(i...)
}

func (l FakeLogger) Infof(s string, i ...interface{}) {
	fmt.Println(fmt.Sprintf(s, i...))
}

func (l FakeLogger) Info(i ...interface{}) {
	fmt.Println(i...)
}
