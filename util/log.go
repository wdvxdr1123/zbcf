package util

import (
	"bytes"
	"strings"

	"github.com/sirupsen/logrus"
)

// LogFormat specialize for go-cqhttp
type LogFormat struct{}

// Format implements logrus.Formatter
func (f LogFormat) Format(entry *logrus.Entry) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteByte('[')
	buf.WriteString(entry.Time.Format("2006-01-02 15:04:05"))
	buf.WriteString("] [")
	buf.WriteString(strings.ToUpper(entry.Level.String()))
	buf.WriteString("]: ")
	buf.WriteString(entry.Message)
	buf.WriteString(" \n")
	ret := append([]byte(nil), buf.Bytes()...) // copy buffer
	return ret, nil
}
