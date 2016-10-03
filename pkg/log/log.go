package log

import (
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
)

const (
	TypeInfo = iota
	TypeDebug
	TypeLocal
	TypeFatal
)

var defaultLogger = newLogger()

type Event struct {
	Type   int
	Time   time.Time
	Fields Fields
	Index  []string
}

func (e Event) Append(key string, value string) Event {
	_, ok := e.Fields[key]
	if ok {
		e.Fields[key] = e.Fields[key] + " " + value
	} else {
		e.Fields[key] = value
		e.Index = append(e.Index, key)
	}
	return e
}

func (e Event) Remove(key string) Event {
	delete(e.Fields, key)
	var index []string
	for _, field := range e.Index {
		if field != key {
			index = append(index, field)
		}
	}
	e.Index = index
	return e
}

type Fields map[string]string

type FieldProcessor func(e Event, field interface{}) (Event, bool)

type Observer interface {
	Log(e Event)
}

type logger struct {
	sync.Mutex
	debug     bool
	local     bool
	processor FieldProcessor
	observers map[Observer]struct{}
}

func newLogger() *logger {
	return &logger{
		observers: make(map[Observer]struct{}),
	}
}

func (l *logger) log(typ int, fields []interface{}) {
	if len(fields) == 0 {
		return
	}
	if typ == TypeDebug && !l.debug {
		return
	}
	if typ == TypeLocal && !l.local {
		return
	}
	e := Event{
		Type:   typ,
		Time:   time.Now(),
		Fields: Fields{"pkg": callerPkg()},
	}
	for _, field := range fields {
		e = l.processField(e, field)
	}
	for o, _ := range l.observers {
		o.Log(e)
	}
}

func (l *logger) processField(e Event, field interface{}) Event {
	switch obj := field.(type) {
	case string:
		return e.Append("msg", obj)
	case Fields:
		for k, v := range obj {
			e = e.Append(k, cast.ToString(v))
		}
		return e
	default:
		if l.processor != nil {
			if ee, ok := l.processor(e, field); ok {
				return ee
			}
		}
		if err, ok := field.(error); ok {
			return e.Append("err", err.Error())
		}
		return e.Append("data", cast.ToString(field))
	}
}

func RegisterObserver(o Observer) { defaultLogger.RegisterObserver(o) }
func (l *logger) RegisterObserver(o Observer) {
	l.Lock()
	defer l.Unlock()
	l.observers[o] = struct{}{}
}

func UnregisterObserver(o Observer) { defaultLogger.UnregisterObserver(o) }
func (l *logger) UnregisterObserver(o Observer) {
	l.Lock()
	defer l.Unlock()
	delete(l.observers, o)
}

func SetDebug(debug bool) { defaultLogger.SetDebug(debug) }
func (l *logger) SetDebug(debug bool) {
	l.Lock()
	defer l.Unlock()
	l.debug = debug
}

func SetLocal(local bool) { defaultLogger.SetLocal(local) }
func (l *logger) SetLocal(local bool) {
	l.Lock()
	defer l.Unlock()
	l.local = local
}

func SetFieldProcessor(fn FieldProcessor) { defaultLogger.SetFieldProcessor(fn) }
func (l *logger) SetFieldProcessor(fn FieldProcessor) {
	l.Lock()
	defer l.Unlock()
	l.processor = fn
}

func Info(o ...interface{}) { defaultLogger.Info(o...) }
func (l *logger) Info(o ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.log(TypeInfo, o)
}

func Debug(o ...interface{}) { defaultLogger.Debug(o...) }
func (l *logger) Debug(o ...interface{}) {
	l.Lock()
	defer l.Unlock()
	if !l.debug {
		return
	}
	l.log(TypeDebug, o)
}

func Local(o ...interface{}) { defaultLogger.Local(o...) }
func (l *logger) Local(o ...interface{}) {
	l.Lock()
	defer l.Unlock()
	if !l.local {
		return
	}
	l.log(TypeLocal, o)
}

func Fatal(o ...interface{}) { defaultLogger.Fatal(o...) }
func (l *logger) Fatal(o ...interface{}) {
	l.Lock()
	defer l.Unlock()
	l.log(TypeFatal, o)
	os.Exit(1)
}

func callerPkg() string {
	pc := make([]uintptr, 10)
	runtime.Callers(5, pc)
	f := runtime.FuncForPC(pc[0]).Name()
	base := path.Base(f)
	dir := path.Dir(f)
	dotparts := strings.Split(base, ".")
	pathparts := strings.Split(path.Join(dir, dotparts[0]), "/")
	return pathparts[len(pathparts)-1]
}

type ResponseWriter interface {
	loggingResponseWriter
}

func WrapResponseWriter(w http.ResponseWriter) ResponseWriter {
	return makeLogger(w)
}
