package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type FileObjectsService struct {
}

func NewFileObjectsService() *FileObjectsService {
	return &FileObjectsService{}
}

// findLibreOfficePath finds the LibreOffice executable path
func findLibreOfficePath() (string, error) {
	// Check environment variable first
	if customPath := os.Getenv("LIBREOFFICE_PATH"); customPath != "" {
		if _, err := os.Stat(customPath); err == nil {
			return customPath, nil
		}
	}

	// Platform-specific defaults
	var candidates []string

	if runtime.GOOS == "windows" {
		candidates = []string{
			`C:\Program Files\LibreOffice\program\soffice.exe`,
		}
	} else {
		// Linux/Mac
		candidates = []string{
			"/usr/bin/libreoffice",
			"/usr/local/bin/libreoffice",
			"/usr/bin/soffice",
		}
	}

	// Check each candidate path
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Try to find in PATH
	path, err := exec.LookPath("soffice")
	if err == nil {
		return path, nil
	}

	path, err = exec.LookPath("libreoffice")
	if err == nil {
		return path, nil
	}

	return "", fmt.Errorf("LibreOffice not found. Please install LibreOffice or set LIBREOFFICE_PATH environment variable")
}

// ConvertVisioToSVG converts a Visio file to SVG format using LibreOffice
func (s *FileObjectsService) ConvertVisioToSVG(visioPath string) (string, error) {
	// Check if file exists
	if _, err := os.Stat(visioPath); os.IsNotExist(err) {
		return "", fmt.Errorf("visio file not found: %s", visioPath)
	}

	// Verify file extension
	ext := strings.ToLower(filepath.Ext(visioPath))
	if ext != ".vsdx" && ext != ".vsd" {
		return "", fmt.Errorf("unsupported file format: %s. Please provide .vsdx or .vsd file", ext)
	}

	// Find LibreOffice executable
	libreOfficePath, err := findLibreOfficePath()
	if err != nil {
		return "", err
	}

	// Create temporary output directory
	outputDir, err := os.MkdirTemp(os.Getenv("uploadDir"), "visio-svg-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Construct LibreOffice conversion command
	cmd := exec.Command(libreOfficePath,
		"--headless",
		"--convert-to", "svg",
		"--outdir", outputDir,
		visioPath,
	)

	// Run conversion
	output, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(outputDir) // Clean up on error
		return "", fmt.Errorf("conversion failed: %v\nOutput: %s", err, string(output))
	}

	// Find the converted SVG file
	files, err := os.ReadDir(outputDir)
	if err != nil {
		os.RemoveAll(outputDir) // Clean up on error
		return "", fmt.Errorf("failed to read output directory: %v", err)
	}

	var svgFile string
	for _, file := range files {
		if strings.ToLower(filepath.Ext(file.Name())) == ".svg" && !file.IsDir() {
			svgFile = filepath.Join(outputDir, file.Name())
			break
		}
	}

	if svgFile == "" {
		os.RemoveAll(outputDir) // Clean up on error
		return "", fmt.Errorf("no SVG file generated in output directory")
	}

	// Read the SVG content
	svgContent, err := os.ReadFile(svgFile)
	if err != nil {
		os.RemoveAll(outputDir) // Clean up on error
		return "", fmt.Errorf("failed to read SVG file: %v", err)
	}

	// Clean up temp directory after successful read
	os.RemoveAll(outputDir)

	return string(svgContent), nil
}
