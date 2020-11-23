package mdw

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func LogRequest(log *log.Logger) Middleware {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			log.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
			next(w, r, p)
		}
	}
}
