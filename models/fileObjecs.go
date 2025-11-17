package models

type ConversionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	SVG     string `json:"svg,omitempty"`
	Error   string `json:"error,omitempty"`
}
