package server

import (
	"fmt"
	"net/http"
)

type WelcomeHandler struct {
	id string
}

func (h WelcomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint("Welcome to the chain. This is node ", h.id)))
}
