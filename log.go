// KIProtect Hyper
// Copyright (C) 2021-2023 KIProtect GmbH
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package hyper

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

type Logger struct {
}

type Level log.Level

const (
	PanicLogLevel = Level(log.PanicLevel)
	FatalLogLevel = Level(log.FatalLevel)
	ErrorLogLevel = Level(log.ErrorLevel)
	WarnLogLevel  = Level(log.WarnLevel)
	InfoLogLevel  = Level(log.InfoLevel)
	DebugLogLevel = Level(log.DebugLevel)
	TraceLogLevel = Level(log.TraceLevel)
)

func ParseLevel(level string) (Level, error) {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return PanicLogLevel, fmt.Errorf("error parsing level: %w", err)
	}
	return Level(lvl), err
}

type IRISFormatter struct {
	service string
}

func SetLogFormat(format string, service string) error {
	switch format {
	case "iris":
		log.SetFormatter(&IRISFormatter{service: service})
	default:
		return fmt.Errorf("unknown log format: %s", format)
	}
	return nil
}

func (c *IRISFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Format(time.RFC3339Nano)
	loglevel := strings.ToUpper(entry.Level.String())
	pid := os.Getpid()
	thread := "[nan]" // thread ID is not applicable to Golang
	function := "[nan]"
	if entry.Caller != nil {
		function = entry.Caller.Function
	}
	message := entry.Message
	return []byte(fmt.Sprintf("%s %s %d %s %s %s : %s\n", timestamp, loglevel, pid, c.service, thread, function, message)), nil
}

func (l *Logger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func (l *Logger) Info(args ...interface{}) {
	log.Info(args...)
}

func (l *Logger) Warning(args ...interface{}) {
	log.Warning(args...)
}

func (l *Logger) SetLevel(level Level) {
	log.SetLevel(log.Level(level))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	log.Error(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	log.Debug(args...)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

func (l *Logger) Trace(args ...interface{}) {
	log.Trace(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

var Log = Logger{}
