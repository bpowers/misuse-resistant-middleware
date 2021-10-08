package a

import (
	"context"
	goji "goji.io"
	"net/http"

	"goji.io/pat"
)

// NameKey is the type of value needed to access the data A stores in the context for later middlewares.
// There is no way to get a value that satisfies NameKey outside of this package except by calling
// NewAMiddleware().Handler()
type NameKey interface {
	// because this is unexported, nobody outside this package can create a type that
	// satisfies the NameKey interface
	nameKeyMarker()
}

// nameKey is a type that be used as a NameKey interface
type nameKey int
func (nameKey) nameKeyMarker() {}

// the single instance of nameKey.  should be `const`, but go doesn't let interfaces be const
var nameKeyName NameKey = nameKey(0)

type NameExtractorMiddleware struct {
	patternName string
}

func NewNameExtractorMiddleware(name string) *NameExtractorMiddleware {
	return &NameExtractorMiddleware{
		patternName: name,
	}
}

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

// GetName returns the value the NameExtractorMiddleware stuck in the context as a string
func GetName(ctx context.Context, nameKey NameKey) string {
	return ctx.Value(nameKey).(string)
}
