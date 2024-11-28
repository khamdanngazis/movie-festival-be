package handlers_test

import (
	"bytes"
	"encoding/json"
	"movie-festival-be/package/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type LoginResponse struct {
	Data struct {
		Token string `json:"token"`
	}
}

func TestLoginHandler(t *testing.T) {

	request := `{"email" : "admin@movie-festival.com","password" : "Symantec2121"}`

	httpRouter.POST("/auth/login", authHandler.Login)

	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte(request)))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	httpRouter.Mux().Use(middleware.LoggingMiddleware)
	httpRouter.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, response.Data.Token)
}
