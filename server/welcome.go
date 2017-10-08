package server

import (
	"fmt"
	"net/http"
)

type welcomeHandler struct {
	id string
}

func (h welcomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"Message\": \"Welcome to the chain. This is node %s\"}", h.id)))
}
