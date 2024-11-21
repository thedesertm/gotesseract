// Package tesseract provides Go bindings for the Tesseract OCR engine
package tesseract

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	// TesseractCmd is the command used to execute Tesseract OCR
	TesseractCmd = "tesseract"

	// defaultEncoding defines the character encoding for OCR output
	defaultEncoding = "utf-8"

	// supportedFormats lists the image formats that can be processed
	supportedFormats = map[string]bool{
		"jpeg": true, "jpg": true, "png": true,
		"pbm": true, "pgm": true, "ppm": true,
		"tiff": true, "tif": true, "bmp": true,
		"gif": true, "webp": true,
	}
)

// Config holds the configuration options for Tesseract OCR operations
type Config struct {
	// TesseractPath specifies custom path to tesseract executable
	TesseractPath string

	// Language specifies OCR language(s) (e.g., "eng", "eng+fra")
	Language string

	// ConfigFile path to custom Tesseract configuration file
	ConfigFile string

	// Timeout sets maximum duration for OCR operations
	Timeout time.Duration

	// Nice sets process priority (Unix-like systems only)
	Nice int

	// OutputType specifies the format of OCR output
	OutputType OutputType
}

// OutputType defines the available OCR output formats
type OutputType int

const (
	// OutputString returns OCR result as string
	OutputString OutputType = iota

	// OutputBytes returns raw OCR output bytes
	OutputBytes

	// OutputDict returns OCR result as key-value pairs
	OutputDict
)

// SetTesseractCmd sets the path to the Tesseract executable
// It validates the path exists and is executable
func SetTesseractCmd(cmd string) error {
	if _, err := os.Stat(cmd); err != nil {
		if !filepath.IsAbs(cmd) {
			if path, err := exec.LookPath(cmd); err == nil {
				TesseractCmd = path
				return nil
			}
		}
		return errors.New("tesseract binary not found")
	}
	TesseractCmd = cmd
	return nil
}
