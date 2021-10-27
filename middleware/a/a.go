package a

import (
	"context"
	goji "goji.io"
	"net/http"

	"goji.io/pat"
)

// NameKey is the type of value needed to access the name string stashed in the context for later middlewares.
// The only way a user of this package can get a non-nil value of NameKey is by calling NewAMiddleware().Register()
type NameKey interface {
	// because this is unexported, nobody outside this package can create a type that
	// satisfies the NameKey interface
	nameKeyMarker()
}

// nameKey is (the only) type that be used as a NameKey interface
type nameKey int

func (nameKey) nameKeyMarker() {}

// the single instance of nameKey.  should be `const`, but go doesn't let interfaces be const
var nameKeyName NameKey = nameKey(0)

type NameExtractorMiddleware struct {
	patternName string
}

// NewNameExtractorMiddleware returns a middleware to store the named path component (from goji)
// in the context.  (this is a simple example middleware to illustrate the misuse resistant pattern
// in this package)
func NewNameExtractorMiddleware(name string) *NameExtractorMiddleware {
	return &NameExtractorMiddleware{
		patternName: name,
	}
}

// Register calls `mux.Use` on this middleware, returning a NameKey that can be used to access
// the data we store in the context.
func (nem *NameExtractorMiddleware) Register(mux *goji.Mux) NameKey {
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			name := pat.Param(req, nem.patternName)
			ctx := context.WithValue(req.Context(), nameKeyName, name)

			next.ServeHTTP(rw, req.WithContext(ctx))
		})
	})

	return nameKeyName
}

// GetName returns the value the NameExtractorMiddleware stuck in the context as a string.
// It panics if you give it a nil NameKey -- but the only way to do that would be to add a
// line like `var empty NameKey` and pass that in to GetName, which would be sort of hard to
// do by accident?  The goal here is to be misuse-resistant, not misuse-proof (unless that is
// possible and not too baroque)
func GetName(ctx context.Context, nameKey NameKey) string {
	if nameKey == nil {
		panic("don't construct a nil NameKey")
	}
	return ctx.Value(nameKey).(string)
}
