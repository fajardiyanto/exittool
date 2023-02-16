package models

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type GCSResponse struct {
	OriginalURL string `json:"original_url"`
	CompressURL string `json:"compress_url"`
}
