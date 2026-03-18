package dto

type AddAPIKeyRequest struct {
	APIKey string `json:"api_key" validate:"required,max=512"`
}

type AddApiKeyResponse struct {
	Message string `json:"message"`
}
