package tryerr

import (
	"fmt"
	"iter"
	"maps"
	"runtime"
)

const (
	defaultStackDepth = 1
	defaultCode       = -1
	stackHelperSkip   = 3
)

// Error is a structured error that preserves Go's native error chain.
type Error struct {
	code    int
	message string
	attrs   map[string]string
	cause   error
	stack   []uintptr
}

type options struct {
	code       int
	attrs      map[string]string
	stackDepth int
	stackSkip  int
}

// Option configures an Error created by New or Wrap.
type Option func(*options)

// New creates a structured error with message.
func New(message string, opts ...Option) error {
	cfg := applyOptions(opts)
	err := &Error{
		message: message,
		code:    cfg.code,
		attrs:   cfg.attrs,
		stack:   captureStack(cfg.stackSkip, cfg.stackDepth),
	}
	return err
}

// Wrap creates a structured error that wraps cause.
//
// A nil cause returns nil so callers can use Wrap directly in ordinary error
// propagation paths without manufacturing a new failure.
func Wrap(cause error, message string, opts ...Option) error {
	if cause == nil {
		return nil
	}
	cfg := applyOptions(opts)
	err := &Error{
		message: message,
		code:    cfg.code,
		attrs:   cfg.attrs,
		cause:   cause,
		stack:   captureStack(cfg.stackSkip, cfg.stackDepth),
	}
	return err
}

func applyOptions(opts []Option) options {
	cfg := options{stackDepth: defaultStackDepth}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

// WithCode attaches a stable non-zero integer error code.
func WithCode(code int) Option {
	return func(o *options) {
		o.code = code
	}
}

// WithAttr attaches one structured attribute.
func WithAttr(key, value string) Option {
	return func(o *options) {
		if o.attrs == nil {
			o.attrs = make(map[string]string, 1)
		}
		o.attrs[key] = value
	}
}

// WithAttrs attaches structured attributes.
func WithAttrs(attrs map[string]string) Option {
	attrs = maps.Clone(attrs)
	return func(o *options) {
		if len(attrs) == 0 {
			return
		}
		if o.attrs == nil {
			o.attrs = make(map[string]string, len(attrs))
		}
		maps.Copy(o.attrs, attrs)
	}
}

// WithStackDepth sets how many caller frames are captured for diagnostics.
//
// Depth values below 1 are treated as 1. The first captured frame is the caller
// of New or Wrap.
func WithStackDepth(depth int) Option {
	return func(o *options) {
		if depth < 1 {
			depth = defaultStackDepth
		}
		o.stackDepth = depth
	}
}

// WithStackSkip sets how many caller frames to skip before capturing the stack.
//
// Skip is relative to the caller of New or Wrap. Values below 0 are treated as
// 0, so the captured stack never starts inside package internals.
func WithStackSkip(skip int) Option {
	return func(o *options) {
		if skip < 0 {
			skip = 0
		}
		o.stackSkip = skip
	}
}

// Error returns the message with code, stack location, and wrapped cause.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	message := e.message
	if message == "" {
		if e.cause == nil {
			return ""
		}
		return e.cause.Error()
	}
	message = fmt.Sprintf("%s [code=%d]", message, e.Code())
	frames := runtime.CallersFrames(e.stack)
	frame, _ := frames.Next()
	if frame.PC != 0 {
		message = fmt.Sprintf("%s [%s %s:%d]", message, frame.Function, frame.File, frame.Line)
	}
	if e.cause == nil {
		return message
	}
	return message + ": " + e.cause.Error()
}

// Unwrap returns the wrapped cause.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// Message returns the local error message.
func (e *Error) Message() string {
	if e == nil {
		return ""
	}
	return e.message
}

// Code returns the attached error code, or -1 when none is configured.
func (e *Error) Code() int {
	if e == nil || e.code == 0 {
		return defaultCode
	}
	return e.code
}

// Attrs returns a copy of the attached attributes.
func (e *Error) Attrs() map[string]string {
	if e == nil {
		return nil
	}
	return maps.Clone(e.attrs)
}

// StackFrames returns the captured stack as runtime frames.
func (e *Error) StackFrames() iter.Seq[runtime.Frame] {
	if e == nil {
		return emptyFrames
	}
	return framesOf(e.stack)
}

// Code returns the first explicitly configured code found in err's chain.
//
// If no code was configured, Code returns -1.
func Code(err error) int {
	for next := range walkErrors(err) {
		if structured, ok := next.(*Error); ok {
			if structured.code != 0 {
				return structured.code
			}
			continue
		}
		if provider, ok := next.(interface{ Code() int }); ok {
			if code := provider.Code(); code != 0 {
				return code
			}
		}
	}
	return defaultCode
}

// Attrs returns merged structured attributes from err's chain.
//
// Outer errors take precedence when multiple errors provide the same key.
func Attrs(err error) map[string]string {
	var attrs map[string]string
	for next := range walkErrors(err) {
		var nextAttrs map[string]string
		if structured, ok := next.(*Error); ok {
			nextAttrs = structured.attrs
		} else if provider, ok := next.(interface{ Attrs() map[string]string }); ok {
			nextAttrs = provider.Attrs()
		}
		if len(nextAttrs) == 0 {
			continue
		}
		if attrs == nil {
			attrs = make(map[string]string, len(nextAttrs))
		}
		for key, value := range nextAttrs {
			if _, exists := attrs[key]; !exists {
				attrs[key] = value
			}
		}
	}
	return attrs
}

// StackFrames returns the first captured stack found in err's chain.
func StackFrames(err error) iter.Seq[runtime.Frame] {
	for next := range walkErrors(err) {
		if structured, ok := next.(*Error); ok {
			if len(structured.stack) > 0 {
				return framesOf(structured.stack)
			}
			continue
		}
		if provider, ok := next.(interface {
			StackFrames() iter.Seq[runtime.Frame]
		}); ok {
			if nextFrames := provider.StackFrames(); nextFrames != nil {
				return nextFrames
			}
		}
	}
	return emptyFrames
}

func walkErrors(err error) iter.Seq[error] {
	return func(yield func(error) bool) {
		var walk func(error) bool
		walk = func(err error) bool {
			if err == nil {
				return true
			}
			if !yield(err) {
				return false
			}

			switch unwrapped := err.(type) {
			case interface{ Unwrap() []error }:
				for _, child := range unwrapped.Unwrap() {
					if !walk(child) {
						return false
					}
				}
			case interface{ Unwrap() error }:
				return walk(unwrapped.Unwrap())
			}
			return true
		}
		walk(err)
	}
}

var emptyFrames iter.Seq[runtime.Frame] = func(func(runtime.Frame) bool) {}

func framesOf(stack []uintptr) iter.Seq[runtime.Frame] {
	if len(stack) == 0 {
		return emptyFrames
	}
	stack = append([]uintptr(nil), stack...)
	return func(yield func(runtime.Frame) bool) {
		frames := runtime.CallersFrames(stack)
		for {
			frame, more := frames.Next()
			if !yield(frame) || !more {
				return
			}
		}
	}
}

func captureStack(skip, depth int) []uintptr {
	stack := make([]uintptr, depth)
	n := runtime.Callers(stackHelperSkip+skip, stack)
	if n == 0 {
		return nil
	}
	return stack[:n:n]
}
