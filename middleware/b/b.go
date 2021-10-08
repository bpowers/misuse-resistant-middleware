package b

import (
	"fmt"
	"net/http"

	"goedge/middleware/a"
)

// NewEchoNameHandler echos the users name -- it depends on NameExtractorMiddleware running before it.
func NewEchoNameHandler(key a.NameKey) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		name := a.GetName(req.Context(), key)
		_, _ = fmt.Fprintf(rw, "Hello, %s!\n", name)
	})
}
