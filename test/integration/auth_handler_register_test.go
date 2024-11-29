package handlers_test

import (
	"bytes"
	"movie-festival-be/package/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {

	request := `{
		"email" : "joe@gmail.com",
		"password" : "Symantec2121",
		"name" : "joe"
	}`

	httpRouter.POST("/auth/register", authHandler.Register)

	req, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte(request)))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	httpRouter.Mux().Use(middleware.LoggingMiddleware)
	httpRouter.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

}
