package mdw

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Recover - prevents server from crashing by adding recovery functionality
func Recover(log *log.Logger) Middleware {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			defer func() {
				// Use the builtin recover function to check if there has been a
				// panic or not. If there has...
				if err := recover(); err != nil {
					// Set a "Connection: close" header on the response and return error
					log.Printf("recovering from: %s\n", err)
					w.Header().Set("Connection", "close")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("something's wrong, try again later..."))
				}
			}()
			next(w, r, p)
		}
	}
}
