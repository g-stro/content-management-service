package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Status string

const (
	Success Status = "success"
	Fail    Status = "fail"
	Error   Status = "error"
)

// BaseResponse is the core response struct, including the status
type BaseResponse struct {
	Status Status `json:"status"`
}

// SuccessResponse is for successful requests
type SuccessResponse struct {
	BaseResponse
	Data interface{} `json:"data"`
}

// FailResponse is for requests that failed due to client error
type FailResponse struct {
	BaseResponse
	Data interface{} `json:"data"`
}

// ErrorResponse is for requests that failed due to server error
type ErrorResponse struct {
	BaseResponse
	Message string `json:"message"`
}

func HttpSuccess(w http.ResponseWriter, data interface{}, status int, logMsg string) {
	slog.Info(logMsg, "data", data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(SuccessResponse{
		BaseResponse: BaseResponse{Status: Success},
		Data:         data,
	})
	if err != nil {
		slog.Error("failed to encode successful response data into JSON", "error", err)
	}
}

func HttpFail(w http.ResponseWriter, data interface{}, status int, logMsg string) {
	slog.Error(logMsg, "data", data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(FailResponse{
		BaseResponse: BaseResponse{Status: Fail},
		Data:         data,
	})
	if err != nil {
		slog.Error("failed to encode failure response data into JSON", "error", err)
	}
}

func HttpError(w http.ResponseWriter, err error, status int, logMsg string) {
	slog.Error(logMsg, "error", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encodeErr := json.NewEncoder(w).Encode(ErrorResponse{
		BaseResponse: BaseResponse{Status: Error},
		Message:      err.Error(),
	})
	if encodeErr != nil {
		slog.Error("failed to encode error response data into JSON", "error", encodeErr)
	}
}
