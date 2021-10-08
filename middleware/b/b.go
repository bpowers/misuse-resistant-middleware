package b

import (
	"fmt"
	"net/http"

	"goedge/middleware/a"
)

// NewEchoNameHandler echos the users name, extracted by the NameExtractorMiddleware.
func NewEchoNameHandler(key a.NameKey) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		name := a.GetName(req.Context(), key)
		fmt.Fprintf(rw, "Hello, %s!\n", name)
	})
}
