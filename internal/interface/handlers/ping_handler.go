package handlers

import (
	"movie-festival-be/package/middleware"
	"net/http"
)

type Pinghandlers struct {
}

func NewPinghandlers() *Pinghandlers {
	return &Pinghandlers{}
}

func (h *Pinghandlers) Ping(w http.ResponseWriter, r *http.Request) {

	middleware.WriteResponse(w, http.StatusOK, "", "Pong")
}
