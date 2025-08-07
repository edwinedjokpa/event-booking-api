package apiresponse

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(message string, data interface{}) ApiResponse {
	return ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func Error(message string, err interface{}) ApiResponse {
	return ApiResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
}
