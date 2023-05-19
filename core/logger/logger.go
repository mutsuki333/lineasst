/*
	logger.go
	Purpose: Logging utiliy functions.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/02  v1.0.0 Evan Chen   Initial release
*/

package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"app/core/config"
	"app/core/property"

	"golang.org/x/exp/slog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger constructs a new *slog.Logger with @file and @opt respecting the configs `property.LOG_*`
//
//   - @file is the file name to output to, if `LOG_FILE` is set to true.
//   - Filename can also be a path.
//   - If @file does not end with ".log", it will append ".log" to the end of @file
func NewLogger(file string, opt ...*slog.HandlerOptions) *slog.Logger {

	var op *slog.HandlerOptions

	if len(opt) == 0 || opt[0] == nil {
		op = &slog.HandlerOptions{}
	} else {
		op = opt[0]
	}

	if op.Level == nil {
		op.Level = property.LogLevel
	}
	if op.ReplaceAttr == nil {
		op.ReplaceAttr = replaceAttr
	}

	if !strings.HasSuffix(file, ".log") {
		file = file + ".log"
	}

	var Writer io.Writer
	if config.GetBool(property.LOG_STD) && config.GetBool(property.LOG_FILE) {
		Writer = io.MultiWriter(os.Stdout, rotator(file))
	} else if config.GetBool(property.LOG_STD) {
		Writer = os.Stdout
	} else {
		Writer = rotator(file)
	}

	if config.GetString(property.LOG_FORMAT) == "json" {
		return slog.New(slog.NewJSONHandler(Writer, op))
	} else {
		return slog.New(slog.NewTextHandler(Writer, op))
	}
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if groups != nil {
		return a
	}
	switch a.Key {
	case slog.SourceKey:
		src, ok := a.Value.Any().(*slog.Source)
		if ok {
			return slog.String(a.Key, fmt.Sprintf("%s:%d", getShortFile(src.File), src.Line))
		}
	}
	return a
}

func getShortFile(file string) string {
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			return file[i+1:]
		}
	}
	return file
}

func rotator(name string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   name,
		MaxSize:    config.GetInt(property.LOG_FILE_SIZE),
		MaxAge:     config.GetInt(property.LOG_FILE_AGE),
		MaxBackups: 0,
		LocalTime:  false, Compress: true,
	}
}
