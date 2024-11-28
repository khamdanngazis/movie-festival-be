package router

import (
	"movie-festival-be/package/logging"
	"movie-festival-be/package/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

var muxDispatcher = mux.NewRouter()

type muxRouter struct{}

// NewMuxRouter creates a new muxRouter instance
func NewMuxRouter() Router {
	return &muxRouter{}
}

func (*muxRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("GET")
}
func (*muxRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("POST")
}

func (*muxRouter) PUT(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("PUT")
}

func (*muxRouter) DELETE(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("DELETE")
}

func (*muxRouter) GETWithMiddleware(uri string, f func(w http.ResponseWriter, r *http.Request), middlewares ...mux.MiddlewareFunc) {
	subRouter := muxDispatcher.PathPrefix(uri).Subrouter()
	subRouter.Use(middlewares...)
	subRouter.HandleFunc("", f).Methods("GET")
}

func (*muxRouter) POSTWithMiddleware(uri string, f func(w http.ResponseWriter, r *http.Request), middlewares ...mux.MiddlewareFunc) {
	subRouter := muxDispatcher.PathPrefix(uri).Subrouter()
	subRouter.Use(middlewares...)
	subRouter.HandleFunc("", f).Methods("POST")
}

func (*muxRouter) PUTWithMiddleware(uri string, f func(w http.ResponseWriter, r *http.Request), middlewares ...mux.MiddlewareFunc) {
	subRouter := muxDispatcher.PathPrefix(uri).Subrouter()
	subRouter.Use(middlewares...)
	subRouter.HandleFunc("", f).Methods("PUT")
}

func (*muxRouter) DELETEWithMiddleware(uri string, f func(w http.ResponseWriter, r *http.Request), middlewares ...mux.MiddlewareFunc) {
	subRouter := muxDispatcher.PathPrefix(uri).Subrouter()
	subRouter.Use(middlewares...)
	subRouter.HandleFunc("", f).Methods("DELETE")
}

func (m *muxRouter) Mux() *mux.Router {
	return muxDispatcher
}

// SERVE method to start the server
func (*muxRouter) SERVE(port string) {
	logging.Log.Infof("Http server listening on port %s", port)
	// Apply Logging Middleware globally
	muxDispatcher.Use(middleware.LoggingMiddleware)
	//muxDispatcher.Use(middleware.AuthMiddleware)
	http.ListenAndServe(port, muxDispatcher)
}
