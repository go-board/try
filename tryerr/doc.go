// Package tryerr provides small structured errors that keep Go's native error
// chain intact.
//
// Use New for a structured root error and Wrap for adding structured context to
// an existing error. Wrap returns nil when the cause is nil, so it can be used
// directly in ordinary error propagation paths.
//
// Errors can carry a stable non-zero integer code, string attributes, and a
// captured stack. WithAttrs copies its input map. Code 0 is reserved for unset
// values. Code returns -1 when no non-zero code is found in the chain. Attrs
// merges attributes from outer to inner errors, keeping the outer value when
// keys collide.
//
// New and Wrap capture one caller frame by default. WithStackDepth changes how
// many frames are captured, and WithStackSkip skips additional callers relative
// to the caller of New or Wrap. StackFrames returns the first available stack as
// an iter.Seq[runtime.Frame]. Error includes the code and first captured frame
// in the formatted error text when one is available.
//
// The extractor functions primarily read Error values. They also accept small
// same-shape providers in the chain: interface{ Code() int },
// interface{ Attrs() map[string]string }, and
// interface{ StackFrames() iter.Seq[runtime.Frame] }.
package tryerr
