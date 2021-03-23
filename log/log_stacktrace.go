package log

import (
	"fmt"

	"github.com/pkg/errors"
)

// modified version of zerolog.ErrorStackMarshaler

type state struct {
	b []byte
}

// Write implement fmt.Formatter interface.
func (s *state) Write(b []byte) (n int, err error) {
	s.b = b

	return len(b), nil
}

// Width implement fmt.Formatter interface.
func (s *state) Width() (wid int, ok bool) {
	return 0, false
}

// Precision implement fmt.Formatter interface.
func (s *state) Precision() (prec int, ok bool) {
	return 0, false
}

// Flag implement fmt.Formatter interface.
func (s *state) Flag(_ int) bool {
	return false
}

func frameField(f errors.Frame, s *state, c rune) string {
	f.Format(s, c)

	return string(s.b)
}

func marshalStack(err error) interface{} {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	var sterr stackTracer

	ok := errors.As(err, &sterr)
	if !ok {
		return nil
	}

	st := sterr.StackTrace()
	s := &state{}
	out := make([]string, 0, len(st))

	for _, frame := range st {
		out = append(out, fmt.Sprintf("%s:%s %s",
			frameField(frame, s, 's'),
			frameField(frame, s, 'd'),
			frameField(frame, s, 'n'),
		))
	}

	return out
}
