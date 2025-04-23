//go:build !integration

package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHttpSuccess(t *testing.T) {
	tests := []struct {
		name       string
		data       interface{}
		status     int
		logMsg     string
		wantStatus int
		wantBody   SuccessResponse
	}{
		{
			name:       "valid success",
			data:       map[string]string{"key": "value"},
			status:     http.StatusOK,
			logMsg:     "Success log",
			wantStatus: http.StatusOK,
			wantBody: SuccessResponse{
				BaseResponse: BaseResponse{Status: Success},
				Data:         map[string]interface{}{"key": "value"},
			},
		},
		{
			name:       "nil data",
			data:       nil,
			status:     http.StatusNoContent,
			logMsg:     "No data",
			wantStatus: http.StatusNoContent,
			wantBody: SuccessResponse{
				BaseResponse: BaseResponse{Status: Success},
				Data:         nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			HttpSuccess(w, tt.data, tt.status, tt.logMsg)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			var responseBody SuccessResponse
			err := json.NewDecoder(resp.Body).Decode(&responseBody)
			if err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if !reflect.DeepEqual(responseBody, tt.wantBody) {
				t.Errorf("expected body %v, got %v", tt.wantBody, responseBody)
			}
		})
	}
}

func TestHttpFail(t *testing.T) {
	tests := []struct {
		name       string
		data       interface{}
		status     int
		logMsg     string
		wantStatus int
		wantBody   FailResponse
	}{
		{
			name:       "valid fail",
			data:       "failure reason",
			status:     http.StatusInternalServerError,
			logMsg:     "Failure occurred",
			wantStatus: http.StatusInternalServerError,
			wantBody: FailResponse{
				BaseResponse: BaseResponse{Status: Fail},
				Data:         "failure reason",
			},
		},
		{
			name:       "empty data",
			data:       "",
			status:     http.StatusBadRequest,
			logMsg:     "Bad request",
			wantStatus: http.StatusBadRequest,
			wantBody: FailResponse{
				BaseResponse: BaseResponse{Status: Fail},
				Data:         "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			HttpFail(w, tt.data, tt.status, tt.logMsg)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			var responseBody FailResponse
			err := json.NewDecoder(resp.Body).Decode(&responseBody)
			if err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if responseBody != tt.wantBody {
				t.Errorf("expected body %v, got %v", tt.wantBody, responseBody)
			}
		})
	}
}

func TestHttpError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		status     int
		logMsg     string
		wantStatus int
		wantBody   ErrorResponse
	}{
		{
			name:       "valid error",
			err:        errors.New("something went wrong"),
			status:     http.StatusInternalServerError,
			logMsg:     "Error occurred",
			wantStatus: http.StatusInternalServerError,
			wantBody: ErrorResponse{
				BaseResponse: BaseResponse{Status: Error},
				Message:      "something went wrong",
			},
		},
		{
			name:       "nil error",
			err:        nil,
			status:     http.StatusInternalServerError,
			logMsg:     "Error occurred",
			wantStatus: http.StatusInternalServerError,
			wantBody: ErrorResponse{
				BaseResponse: BaseResponse{Status: Error},
				Message:      "no error message provided",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			HttpError(w, tt.err, tt.status, tt.logMsg)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			var responseBody ErrorResponse
			err := json.NewDecoder(resp.Body).Decode(&responseBody)
			if err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if responseBody != tt.wantBody {
				t.Errorf("expected body %v, got %v", tt.wantBody, responseBody)
			}
		})
	}
}
