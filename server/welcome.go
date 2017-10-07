package server

import (
	"fmt"
	"net/http"
)

type welcomeHandler struct {
	id string
}

func (h welcomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprint("Welcome to the chain. This is node ", h.id)))
}
