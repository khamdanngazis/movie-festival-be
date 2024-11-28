package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"movie-festival-be/package/logging"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique request ID using UUID.
		requestID := uuid.New().String()
		startTime := time.Now()

		// Add the request ID to the request context.
		ctx := context.WithValue(r.Context(), logging.RequestIDKey, requestID)

		// Set the updated context in the request.
		r = r.WithContext(ctx)

		// Log the request details with the generated request ID using logrus.
		field := logrus.Fields{
			"request_id": requestID,
			"method":     r.Method,
			"uri":        r.RequestURI,
			"proto":      r.Proto,
			"body":       extractRequestBody(r),
		}
		logging.LogCustomField(logrus.InfoLevel, field, "Incoming request")

		// Create a custom response writer to capture the response status and body
		responseRecorder := NewResponseRecorder(w)

		// Call the next handler in the chain.
		next.ServeHTTP(responseRecorder, r)

		// Log the response details
		field = logrus.Fields{
			"request_id": requestID,
			"status":     responseRecorder.Status(),
			"body":       responseRecorder.Body(),
			"duration":   fmt.Sprintf("%v", time.Since(startTime)),
		}

		logging.LogCustomField(logrus.InfoLevel, field, "Outgoing response")
	})
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestID := uuid.New().String()
	startTime := time.Now()
	// Add the request ID to the request context.
	ctx = context.WithValue(ctx, logging.RequestIDKey, requestID)
	field := logrus.Fields{
		"request_id": requestID,
		"method":     fmt.Sprintf("%s", info.FullMethod),
		"proto":      "proto-buffer",
		"message":    grpcMessageToJSON(req),
	}
	logging.LogCustomField(logrus.InfoLevel, field, "Incoming request")

	// Call the actual handler to process the request
	resp, err := handler(ctx, req)

	// Log the response and duration
	field = logrus.Fields{
		"request_id": requestID,
		"method":     fmt.Sprintf("%s", info.FullMethod),
		"message":    grpcMessageToJSON(resp),
		"duration":   fmt.Sprintf("%v", time.Since(startTime)),
	}

	logging.LogCustomField(logrus.InfoLevel, field, "Outgoing response")

	return resp, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func grpcMessageToJSON(s interface{}) string {
	result := ""
	if s != nil {
		reqJSON, _ := protojson.Marshal(s.(proto.Message))
		result = string(reqJSON)
	}
	return result
}

func extractRequestBody(r *http.Request) string {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "Error reading request body"
	}
	defer r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return string(body)
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{w, http.StatusOK, bytes.Buffer{}}
}

// WriteHeader captures the status code.
func (r *ResponseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Write captures the response body.
func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Status returns the captured status code.
func (r *ResponseRecorder) Status() int {
	return r.status
}

// Body returns the captured response body.
func (r *ResponseRecorder) Body() string {
	return r.body.String()
}
