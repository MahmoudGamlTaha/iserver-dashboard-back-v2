package handlers

import (
	"encoding/json"
	"enterprise-architect-api/config"
	"enterprise-architect-api/models"
	"enterprise-architect-api/services"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileObjectsHandler struct {
	fileObjectsService *services.FileObjectsService
}

func NewFileObjectsHandler(fileObjectsService *services.FileObjectsService) *FileObjectsHandler {
	return &FileObjectsHandler{fileObjectsService: fileObjectsService}
}

// ConvertVisioToSVGHandler handles Visio file upload and converts it to SVG
func (fh *FileObjectsHandler) ConvertVisioToSVGHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")

	// Validate method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ConversionResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(config.MaxUploadSize); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ConversionResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to parse form: %v", err),
		})
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ConversionResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to get uploaded file: %v", err),
		})
		return
	}
	defer file.Close()

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".vsdx" && ext != ".vsd" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ConversionResponse{
			Success: false,
			Error:   "Unsupported file format. Please upload .vsdx or .vsd files",
		})
		return
	}

	// Create temporary file
	tempFile, err := os.CreateTemp(os.Getenv("uploadDir"), "visio-*"+ext)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ConversionResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to create temporary file: %v", err),
		})
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name()) // Clean up temp file

	// Copy uploaded content to temporary file
	if _, err := io.Copy(tempFile, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ConversionResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to save uploaded file: %v", err),
		})
		return
	}

	// Convert Visio to SVG
	svgContent, err := fh.fileObjectsService.ConvertVisioToSVG(tempFile.Name())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ConversionResponse{
			Success: false,
			Error:   fmt.Sprintf("Conversion failed: %v", err),
		})
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ConversionResponse{
		Success: true,
		Message: "Successfully converted Visio to SVG",
		SVG:     svgContent,
	})
}
