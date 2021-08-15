package logger

import (
	"fmt"
	"time"

	"galaxyzeta.io/engine/infra/concurrency/lock"
)

type ILogSink interface {
	Write(string)
	Flush()
}

type AsyncLoggerSink struct {
	buffer        []byte
	pos           int
	flushDuration time.Duration
	timer         time.Ticker
	mu            *lock.SpinLock
}

type StdoutSink struct{}

// -----------------------

func NewStdoutSink() StdoutSink {
	return StdoutSink{}
}

func (sink StdoutSink) Write(str string) {
	fmt.Printf(str)
}

func (sink StdoutSink) Flush() {
	// do nothing
}

// -----------------------

func NewAsyncLoggerSink(cacheSize int, duration time.Duration) *AsyncLoggerSink {
	sink := &AsyncLoggerSink{
		buffer:        make([]byte, cacheSize),
		pos:           0,
		flushDuration: duration,
		timer:         *time.NewTicker(duration),
		mu:            &lock.SpinLock{},
	}
	go sink.logFlushRoutine()
	return sink
}

func (sink *AsyncLoggerSink) logFlushRoutine() {
	for {
		<-sink.timer.C
		sink.Flush()
	}
}

func (sink *AsyncLoggerSink) Write(str string) {
	sink.mu.Lock()
	afterLen := sink.pos + len(str)
	if afterLen > len(sink.buffer) {
		// will trigger flush
		sink.timer.Reset(sink.flushDuration)
		sink.doFlushNoLock()
	} else {
		copy(sink.buffer[sink.pos:], []byte(str))
		sink.pos = afterLen
	}
	sink.mu.Unlock()
}

func (sink *AsyncLoggerSink) Flush() {
	sink.mu.Lock()
	sink.doFlushNoLock()
	sink.mu.Unlock()
}

func (sink *AsyncLoggerSink) doFlushNoLock() {
	fmt.Print(string(sink.buffer[:sink.pos]))
	sink.pos = 0
}
