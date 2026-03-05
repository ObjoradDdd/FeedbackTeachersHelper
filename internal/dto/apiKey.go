package dto

type AddAPIKeyRequest struct {
	APIKey string `json:"api_key"`
}

type AddApiKeyResponse struct {
	Message string `json:"message"`
}
