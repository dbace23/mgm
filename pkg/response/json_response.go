package response

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	Details   any    `json:"details,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func Success(message string, data any) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func Error(err string, message string, details any) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   err,
		Message: message,
		Details: details,
	}
}

func ErrorWithRequestID(err string, message string, details any, requestID string) ErrorResponse {
	return ErrorResponse{
		Success:   false,
		Error:     err,
		Message:   message,
		Details:   details,
		RequestID: requestID,
	}
}
