// Wrapper over the slog package

package slogger

import (
	"io"
	"reflect"
	"strconv"

	"golang.org/x/exp/slog"
)

type Logger struct {
	Logger *slog.Logger
}

type Arguments struct {
	args any
}

// Generic logger wrapper to log events as JSON to a file, stdout, or any io.Writer implementation
func NewLogger(f io.Writer) *Logger {
	return &Logger{Logger: slog.New(slog.NewJSONHandler(f))}
}

// LogEvent will log a given event with, or without arguments
// Will accept a map or variadic array of strings
func (p *Logger) LogEvent(logType string, msg string, args ...any) {
	arguments := NewArgumentMapper(args).mapArguments().attributeBuilder()

	switch logType {
	case "warn":
		p.Logger.Warn(msg, arguments...)
	case "debug":
		p.Logger.Debug(msg, arguments...)
	case "info":
		p.Logger.Info(msg, arguments...)
	}
}

// Logs error messages
func (p *Logger) LogError(msg string, err error, args ...any) error {
	p.Logger.Error(msg, err, NewArgumentMapper(args).mapArguments().attributeBuilder()...)

	return err
}

// NewArgumentMapper creates a pointer to Arguments struct which is processed and passed into the logger
func NewArgumentMapper(args []any) *Arguments {
	return &Arguments{args: args}
}

// Generates an argument map to be passed into the log even
func (a *Arguments) mapArguments() *Arguments {
	var args []any
	if a, ok := a.args.([]any); ok {
		args = a
	}

	argMap := make(map[string]any)

	// Only 1 argument will create a key value pair of empty string
	// or if the singular argument is a map, the map will be passed on
	if len(args) == 1 {
		if isMap(args[0]) {
			return a.SetMap(args[0].(map[string]any))
		}

		addArgsToMap(args[0], "", argMap)

		return a.SetMap(argMap)
	}

	// Will handle slice of any values including maps.
	// If maps are passed within the []any slice they are passed over
	for i := 0; i < len(args)-1; i += 2 {
		if isMap(args[i]) && !isMap(args[i+1]) {
			addArgsToMap(args[i+1], "", argMap)
			addMapArgsToMap(args[i], argMap)

			continue
		} else if !isMap(args[i]) && isMap(args[i+1]) {
			addArgsToMap(args[i], "", argMap)
			addMapArgsToMap(args[i+1], argMap)

			continue
		} else if isMap(args[i]) && isMap(args[i+1]) {
			addMapArgsToMap(joinMaps(args[i+1], args[i]), argMap)

			continue
		}

		addArgsToMap(args[i], args[i+1], argMap)
	}

	// If odd length, insert last elements onto tail under misc to avoid data loss
	if len(args)%2 != 0 {
		// There will only ever be 1 if odd.
		arg := args[len(args)-1]
		if isMap(arg) {
			addMapArgsToMap(arg, argMap)
		} else {
			addArgsToMap("miscFields", arg, argMap)
		}
	}

	return a.SetMap(argMap)
}

func (a *Arguments) SetMap(value map[string]any) *Arguments {
	a.args = value

	return a
}

func joinMaps(maps ...any) map[string]any {
	joinMap := make(map[string]any)

	for _, iMap := range maps {
		if isMap(iMap) {
			if v, ok := iMap.(map[string]any); ok {
				for key, value := range v {
					joinMap[key] = value
				}
			}
		}
	}

	return joinMap
}

// Converts incoming attributes arguments args into slog Atttributes to be further mapped within record.setAttrsFromArgs()
func (a *Arguments) attributeBuilder() []any {
	var args map[string]any
	var ok bool

	if args, ok = a.args.(map[string]any); !ok {
		return nil
	}

	if len(args) == 0 {
		return nil
	}

	var attrs []any

	for k, v := range args {
		attrs = append(attrs, slog.Attr{Key: k, Value: slog.AnyValue(v)})
	}

	return attrs
}

func addMapArgsToMap(data any, argMap map[string]any) {
	if isMap(data) {
		if data, ok := data.(map[string]any); ok {
			for k, v := range data {
				argMap[k] = v
			}
		}
	}
}

// Handle alternate types passed in as keys
func addArgsToMap(key, value any, argMap map[string]any) {
	switch val := key.(type) {
	case string:
		argMap[val] = value
	case int:
		argMap[strconv.Itoa(val)] = value
	case bool:
		argMap[strconv.FormatBool(val)] = value
	}
}

func isMap(value any) bool {
	return reflect.ValueOf(value).Kind() == reflect.Map
}
