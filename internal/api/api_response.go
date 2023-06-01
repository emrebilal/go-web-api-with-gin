package api

type ApiResponse struct {
	Data    *interface{} `json:"Data"`
	Message string       `json:"Message"`
}

func RespondOk(data interface{}) *ApiResponse {
	apiResponse := ApiResponse{
		Data:    &data,
		Message: "Success",
	}

	return &apiResponse
}

func RespondError(message string) *ApiResponse {
	apiResponse := ApiResponse{
		Data:    nil,
		Message: message,
	}

	return &apiResponse
}
