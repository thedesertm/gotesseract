// Package tesseract provides Go bindings for the Tesseract OCR engine
package tesseract

import (
	"fmt"
	"strings"
)

// Common errors returned by the tesseract package
var (
	// ErrTesseractNotFound indicates Tesseract binary is not installed or not in PATH
	ErrTesseractNotFound = fmt.Errorf("tesseract is not installed or not in PATH")

	// ErrInvalidImage indicates the provided image is nil or invalid
	ErrInvalidImage = fmt.Errorf("invalid or unsupported image format")

	// ErrProcessTimeout indicates OCR operation exceeded timeout
	ErrProcessTimeout = fmt.Errorf("OCR process timeout")

	// ErrLanguageNotFound indicates requested language data is not installed
	ErrLanguageNotFound = fmt.Errorf("specified language data not found")

	// ErrUnsupportedFormat indicates image format is not supported by Tesseract
	ErrUnsupportedFormat = fmt.Errorf("unsupported image format")

	// ErrInvalidConfig indicates invalid Tesseract configuration
	ErrInvalidConfig = fmt.Errorf("invalid Tesseract config")

	// ErrEmptyOutput indicates OCR produced no output
	ErrEmptyOutput = fmt.Errorf("OCR produced no output")

	// ErrInvalidOutput indicates OCR output is invalid or corrupted
	ErrInvalidOutput = fmt.Errorf("invalid OCR output")
)

// OCRError represents a Tesseract process error with exit code and stderr output
type OCRError struct {
	// Code is the process exit code
	Code int

	// Stderr contains error output from Tesseract
	Stderr string
}

// Error implements the error interface for OCRError
func (e *OCRError) Error() string {
	msg := strings.TrimSpace(e.Stderr)
	if msg == "" {
		msg = "unknown error"
	}
	return fmt.Sprintf("tesseract error (code %d): %s", e.Code, msg)
}

// IsTimeout returns true if the error indicates a timeout
func IsTimeout(err error) bool {
	return err == ErrProcessTimeout
}

// IsNotFound returns true if the error indicates Tesseract is not installed
func IsNotFound(err error) bool {
	return err == ErrTesseractNotFound
}
