// Package tesseract provides Go bindings for the Tesseract OCR engine
package tesseract

import (
	"context"
	"fmt"
	"image"
	"os"
	"path/filepath"
)

// SupportedExtension represents a supported output format and its configuration
type SupportedExtension struct {
	config  string
	version string // Minimum Tesseract version required
}

// Supported Tesseract output extensions and their configurations
var supportedExtensions = map[string]SupportedExtension{
	"hocr": {"tessedit_create_hocr=1", "3.05"},
	"xml":  {"tessedit_create_alto=1", "4.1.0"},
	"tsv":  {"tessedit_create_tsv=1", "3.05"},
	"pdf":  {"tessedit_create_pdf=1", "3.05"},
}

// ErrUnsupportedExtension indicates requested output format is not supported
var ErrUnsupportedExtension = fmt.Errorf("unsupported output extension")

// ImageToExtension performs OCR and returns output in specified format
func (c *Client) ImageToExtension(img image.Image, lang, extension string) (string, error) {
	if err := validateImageFormat(img); err != nil {
		return "", err
	}

	extConfig, ok := supportedExtensions[extension]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedExtension, extension)
	}

	ctx := context.Background()
	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	tmpDir, err := createTempDir()
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	outBase := filepath.Join(tmpDir, "output")
	args := []string{outBase}

	if lang != "" {
		if err := validateLanguage(lang); err != nil {
			return "", err
		}
		args = append(args, "-l", lang)
	}

	args = append(args, "-c", extConfig.config)

	if _, err := runOCR(ctx, img, lang, args); err != nil {
		return "", err
	}

	outputFile := outBase + "." + extension
	output, err := os.ReadFile(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to read output file: %w", err)
	}

	return string(output), nil
}
