package tryerr

import (
	"errors"
	"fmt"
	"iter"
	"runtime"
	"strings"
	"testing"
)

var sentinelErr = errors.New("sentinel failure")

type joinedError []error

func (e joinedError) Error() string { return "joined" }

func (e joinedError) Unwrap() []error { return []error(e) }

type codedError struct {
	code int
}

func (e codedError) Error() string { return "coded" }

func (e codedError) Code() int { return e.code }

type attributedError struct {
	attrs map[string]string
}

func (e attributedError) Error() string { return "attributed" }

func (e attributedError) Attrs() map[string]string { return e.attrs }

type stackFrameSeqError struct {
	frames []runtime.Frame
}

func (e stackFrameSeqError) Error() string { return "stack frame seq" }

func (e stackFrameSeqError) StackFrames() iter.Seq[runtime.Frame] {
	return func(yield func(runtime.Frame) bool) {
		for _, frame := range e.frames {
			if !yield(frame) {
				return
			}
		}
	}
}

func TestNewCreatesStructuredError(t *testing.T) {
	err := New(
		"load config",
		WithCode(1001),
		WithAttr("component", "config"),
	)

	var structured *Error
	if !errors.As(err, &structured) {
		t.Fatalf("errors.As(err, *Error) = false")
	}
	if got := err.Error(); !strings.HasPrefix(got, "load config [code=1001] [") {
		t.Fatalf("err.Error() = %q, want load config with code and frame", got)
	}
	if got := err.Error(); !strings.Contains(got, "TestNewCreatesStructuredError") {
		t.Fatalf("err.Error() = %q, want frame function TestNewCreatesStructuredError", got)
	}
	if got, want := structured.Message(), "load config"; got != want {
		t.Fatalf("Message() = %q, want %q", got, want)
	}
	if got, want := structured.Code(), 1001; got != want {
		t.Fatalf("Code() = %d, want %d", got, want)
	}
	if got := structured.Attrs()["component"]; got != "config" {
		t.Fatalf("Attrs()[component] = %q, want config", got)
	}
	if unwrapped := errors.Unwrap(err); unwrapped != nil {
		t.Fatalf("errors.Unwrap(err) = %v, want nil", unwrapped)
	}
}

func TestWrapPreservesErrorIdentity(t *testing.T) {
	err := Wrap(
		sentinelErr,
		"read config",
		WithCode(2001),
	)

	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false; got %v", err)
	}
	if got := err.Error(); !strings.HasPrefix(got, "read config [code=2001] [") {
		t.Fatalf("err.Error() = %q, want read config with code and frame", got)
	}
	if got := err.Error(); !strings.Contains(got, "TestWrapPreservesErrorIdentity") {
		t.Fatalf("err.Error() = %q, want frame function TestWrapPreservesErrorIdentity", got)
	}
	if got := err.Error(); !strings.HasSuffix(got, ": sentinel failure") {
		t.Fatalf("err.Error() = %q, want wrapped sentinel failure suffix", got)
	}
	var structured *Error
	if !errors.As(err, &structured) {
		t.Fatalf("errors.As(err, *Error) = false")
	}
	if unwrapped := structured.Unwrap(); unwrapped != sentinelErr {
		t.Fatalf("Unwrap() = %v, want sentinelErr", unwrapped)
	}
}

func TestWrapNilReturnsNil(t *testing.T) {
	if err := Wrap(nil, "ignored", WithCode(9999)); err != nil {
		t.Fatalf("Wrap(nil, ...) = %v, want nil", err)
	}
}

func TestDefaultCode(t *testing.T) {
	err := New("plain")

	var structured *Error
	if !errors.As(err, &structured) {
		t.Fatalf("errors.As(err, *Error) = false")
	}
	if got := structured.Code(); got != defaultCode {
		t.Fatalf("Code() = %d, want %d", got, defaultCode)
	}
	if got := Code(err); got != defaultCode {
		t.Fatalf("Code(err) = %d, want %d", got, defaultCode)
	}
	if got := err.Error(); !strings.HasPrefix(got, "plain [code=-1] [") {
		t.Fatalf("err.Error() = %q, want plain with default code and frame", got)
	}
}

func TestEmptyWrappedMessageUsesCauseMessage(t *testing.T) {
	err := Wrap(sentinelErr, "")
	if got, want := err.Error(), "sentinel failure"; got != want {
		t.Fatalf("err.Error() = %q, want %q", got, want)
	}
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("errors.Is(err, sentinelErr) = false")
	}
}

func TestOptionsCopyAttrs(t *testing.T) {
	attrs := map[string]string{"region": "cn"}
	err := New("publish", WithAttrs(attrs))
	attrs["region"] = "us"

	gotAttrs := Attrs(err)
	if got, want := gotAttrs["region"], "cn"; got != want {
		t.Fatalf("Attrs(err)[region] = %q, want %q", got, want)
	}

	gotAttrs["region"] = "sg"
	if got, want := Attrs(err)["region"], "cn"; got != want {
		t.Fatalf("Attrs(err)[region] after caller mutation = %q, want %q", got, want)
	}
}

func TestWithAttrsCopiesInputWhenOptionIsCreated(t *testing.T) {
	attrs := map[string]string{"region": "cn"}
	option := WithAttrs(attrs)
	attrs["region"] = "us"

	err := New("publish", option)
	if got, want := Attrs(err)["region"], "cn"; got != want {
		t.Fatalf("Attrs(err)[region] = %q, want %q", got, want)
	}
}

func TestChainHelpersPreferOuterMetadata(t *testing.T) {
	inner := New(
		"open file",
		WithCode(3001),
		WithAttr("component", "storage"),
		WithAttr("path", "inner"),
	)
	outer := Wrap(
		inner,
		"load config",
		WithCode(4001),
		WithAttr("path", "outer"),
		WithAttr("user", "alice"),
	)

	if got, want := Code(outer), 4001; got != want {
		t.Fatalf("Code(outer) = %d, want %d", got, want)
	}

	attrs := Attrs(outer)
	want := map[string]string{
		"component": "storage",
		"path":      "outer",
		"user":      "alice",
	}
	if len(attrs) != len(want) {
		t.Fatalf("len(Attrs(outer)) = %d, want %d: %v", len(attrs), len(want), attrs)
	}
	for key, value := range want {
		if got := attrs[key]; got != value {
			t.Fatalf("Attrs(outer)[%q] = %q, want %q", key, got, value)
		}
	}
}

func TestCodeSkipsDefaultOuterCode(t *testing.T) {
	inner := New("open file", WithCode(3001))
	outer := Wrap(inner, "load config")

	if got, want := Code(outer), 3001; got != want {
		t.Fatalf("Code(outer) = %d, want %d", got, want)
	}
}

func TestCodeSkipsZeroCode(t *testing.T) {
	inner := New("open file", WithCode(3001))
	outer := Wrap(inner, "load config", WithCode(0))

	if got, want := Code(outer), 3001; got != want {
		t.Fatalf("Code(outer) = %d, want %d", got, want)
	}

	var structured *Error
	if !errors.As(outer, &structured) {
		t.Fatalf("errors.As(outer, *Error) = false")
	}
	if got, want := structured.Code(), defaultCode; got != want {
		t.Fatalf("Code() = %d, want %d", got, want)
	}
}

func TestChainHelpersTraverseJoinedErrors(t *testing.T) {
	err := joinedError{
		sentinelErr,
		New("structured", WithCode(5001), WithAttr("source", "joined")),
	}

	if got, want := Code(err), 5001; got != want {
		t.Fatalf("Code(err) = %d, want %d", got, want)
	}
	if got := Attrs(err)["source"]; got != "joined" {
		t.Fatalf("Attrs(err)[source] = %q, want joined", got)
	}
}

func TestCodeReadsProviderFromChain(t *testing.T) {
	err := fmt.Errorf("outer: %w", codedError{code: 7001})

	if got, want := Code(err), 7001; got != want {
		t.Fatalf("Code(err) = %d, want %d", got, want)
	}
}

func TestCodeSkipsZeroProviderCode(t *testing.T) {
	err := joinedError{
		codedError{code: 0},
		codedError{code: 7001},
	}

	if got, want := Code(err), 7001; got != want {
		t.Fatalf("Code(err) = %d, want %d", got, want)
	}
}

func TestAttrsReadsProviderFromChain(t *testing.T) {
	inner := attributedError{attrs: map[string]string{
		"component": "provider",
		"path":      "inner",
	}}
	outer := Wrap(inner, "outer", WithAttr("path", "outer"))

	attrs := Attrs(outer)
	if got, want := attrs["component"], "provider"; got != want {
		t.Fatalf("Attrs(err)[component] = %q, want %q", got, want)
	}
	if got, want := attrs["path"], "outer"; got != want {
		t.Fatalf("Attrs(err)[path] = %q, want %q", got, want)
	}
}

func TestStackReadsStackFrameSeqProviderFromChain(t *testing.T) {
	err := fmt.Errorf("outer: %w", stackFrameSeqError{frames: []runtime.Frame{
		{Function: "providerStackFrameSeq", File: "provider.go", Line: 23},
	}})

	frames := collectFrames(StackFrames(err))
	if got, want := len(frames), 1; got != want {
		t.Fatalf("len(StackFrames(err)) = %d, want %d", got, want)
	}
	if got, want := frames[0].Function, "providerStackFrameSeq"; got != want {
		t.Fatalf("StackFrames(err)[0].Function = %q, want %q", got, want)
	}
}

func TestNewCapturesDefaultStack(t *testing.T) {
	err := New("with stack")

	frames := collectFrames(StackFrames(err))
	if got, want := len(frames), 1; got != want {
		t.Fatalf("len(StackFrames(err)) = %d, want %d", got, want)
	}
	if !strings.Contains(frames[0].Function, "TestNewCapturesDefaultStack") {
		t.Fatalf("first frame function = %q, want TestNewCapturesDefaultStack", frames[0].Function)
	}
}

func TestErrorStackFramesReturnsIterator(t *testing.T) {
	err := New("with stack")

	var structured *Error
	if !errors.As(err, &structured) {
		t.Fatalf("errors.As(err, *Error) = false")
	}

	frames := collectFrames(structured.StackFrames())
	if got, want := len(frames), 1; got != want {
		t.Fatalf("len((*Error).StackFrames()) = %d, want %d", got, want)
	}
	if !strings.Contains(frames[0].Function, "TestErrorStackFramesReturnsIterator") {
		t.Fatalf("first frame function = %q, want TestErrorStackFramesReturnsIterator", frames[0].Function)
	}
}

func TestWrapCapturesDefaultStack(t *testing.T) {
	err := Wrap(sentinelErr, "with stack")

	frames := collectFrames(StackFrames(err))
	if got, want := len(frames), 1; got != want {
		t.Fatalf("len(StackFrames(err)) = %d, want %d", got, want)
	}
	if !strings.Contains(frames[0].Function, "TestWrapCapturesDefaultStack") {
		t.Fatalf("first frame function = %q, want TestWrapCapturesDefaultStack", frames[0].Function)
	}
}

func TestDefaultStackFrameStartsAtConstructorCaller(t *testing.T) {
	tests := []struct {
		name string
		call func() (error, runtime.Frame)
	}{
		{
			name: "new",
			call: newAtDefaultStackFrameCallSite,
		},
		{
			name: "wrap",
			call: wrapAtDefaultStackFrameCallSite,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, want := tt.call()

			frames := collectFrames(StackFrames(err))
			if got, want := len(frames), 1; got != want {
				t.Fatalf("len(StackFrames(err)) = %d, want %d", got, want)
			}
			assertFrame(t, frames[0], want)
		})
	}
}

func TestStackDepthOption(t *testing.T) {
	err := New("with stack", WithStackDepth(2))

	frames := collectFrames(StackFrames(err))
	if got, want := len(frames), 2; got != want {
		t.Fatalf("len(StackFrames(err)) = %d, want %d", got, want)
	}
	if !strings.Contains(frames[0].Function, "TestStackDepthOption") {
		t.Fatalf("first frame function = %q, want TestStackDepthOption", frames[0].Function)
	}
}

func TestStackDepthOptionMinimum(t *testing.T) {
	err := New("with stack", WithStackDepth(0))

	frames := collectFrames(StackFrames(err))
	if got, want := len(frames), 1; got != want {
		t.Fatalf("len(StackFrames(err)) = %d, want %d", got, want)
	}
}

func TestStackSkipOption(t *testing.T) {
	err := newWithStackSkip(1)

	frames := collectFrames(StackFrames(err))
	if got, want := len(frames), 1; got != want {
		t.Fatalf("len(StackFrames(err)) = %d, want %d", got, want)
	}
	if !strings.Contains(frames[0].Function, "TestStackSkipOption") {
		t.Fatalf("first frame function = %q, want TestStackSkipOption", frames[0].Function)
	}
	if got := err.Error(); !strings.Contains(got, "TestStackSkipOption") {
		t.Fatalf("err.Error() = %q, want frame function TestStackSkipOption", got)
	}
}

func TestStackSkipOptionMinimum(t *testing.T) {
	err := newWithStackSkip(-1)

	frames := collectFrames(StackFrames(err))
	if got, want := len(frames), 1; got != want {
		t.Fatalf("len(StackFrames(err)) = %d, want %d", got, want)
	}
	if !strings.Contains(frames[0].Function, "newWithStackSkip") {
		t.Fatalf("first frame function = %q, want newWithStackSkip", frames[0].Function)
	}
}

func TestMissingMetadata(t *testing.T) {
	err := errors.New("plain")

	if got := Code(err); got != defaultCode {
		t.Fatalf("Code(plain) = %d, want %d", got, defaultCode)
	}
	if attrs := Attrs(err); attrs != nil {
		t.Fatalf("Attrs(plain) = %v, want nil", attrs)
	}
	for frame := range StackFrames(err) {
		t.Fatalf("StackFrames(plain) yielded frame %v, want none", frame)
	}
}

func newWithStackSkip(skip int) error {
	return New("with stack", WithStackSkip(skip))
}

//go:noinline
func newAtDefaultStackFrameCallSite() (error, runtime.Frame) {
	want := nextLineFrame()
	err := New("with stack")
	return err, want
}

//go:noinline
func wrapAtDefaultStackFrameCallSite() (error, runtime.Frame) {
	want := nextLineFrame()
	err := Wrap(sentinelErr, "with stack")
	return err, want
}

//go:noinline
func nextLineFrame() runtime.Frame {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("runtime.Caller failed")
	}
	return runtime.Frame{
		Function: runtime.FuncForPC(pc).Name(),
		File:     file,
		Line:     line + 1,
	}
}

func assertFrame(t *testing.T, got, want runtime.Frame) {
	t.Helper()
	if got.Function != want.Function || got.File != want.File || got.Line != want.Line {
		t.Fatalf(
			"frame = {Function:%q File:%q Line:%d}, want {Function:%q File:%q Line:%d}",
			got.Function,
			got.File,
			got.Line,
			want.Function,
			want.File,
			want.Line,
		)
	}
}

func collectFrames(seq iter.Seq[runtime.Frame]) []runtime.Frame {
	var frames []runtime.Frame
	for frame := range seq {
		frames = append(frames, frame)
	}
	return frames
}
