// Package tesseract provides Go bindings for the Tesseract OCR engine
package tesseract

import (
	"context"
	"fmt"
	"image"
	"time"
)

// Client manages Tesseract OCR operations
type Client struct {
	config Config
}

// NewClient creates a new Tesseract client with the given configuration
func NewClient(cfg Config) (*Client, error) {
	if cfg.TesseractPath != "" {
		if err := SetTesseractCmd(cfg.TesseractPath); err != nil {
			return nil, err
		}
	}
	if err := checkVersion(); err != nil {
		return nil, err
	}
	return &Client{config: cfg}, nil
}

// ImageToString performs OCR on an image using the default client
func ImageToString(img image.Image, lang string) (string, error) {
	return DefaultClient.ImageToString(img, lang)
}

// ImageToString performs OCR on an image and returns the extracted text
func (c *Client) ImageToString(img image.Image, lang string) (string, error) {
	if err := validateImageFormat(img); err != nil {
		return "", err
	}

	ctx := context.Background()
	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	out, err := runOCR(ctx, img, lang, []string{"stdout"})
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ImageToOutput performs OCR and returns the result in the specified format
func (c *Client) ImageToOutput(img image.Image, lang string, outputType OutputType) (interface{}, error) {
	if err := validateImageFormat(img); err != nil {
		return "", err
	}

	ctx := context.Background()
	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	out, err := runOCR(ctx, img, lang, []string{"stdout"})
	if err != nil {
		return nil, err
	}

	switch outputType {
	case OutputString:
		return string(out), nil
	case OutputBytes:
		return out, nil
	case OutputDict:
		return parseTSV(string(out)), nil
	default:
		return nil, fmt.Errorf("unsupported output type")
	}
}

// ImageToFile performs OCR and saves the result to the specified file
func (c *Client) ImageToFile(img image.Image, lang, outputFile string) error {
	if err := validateImageFormat(img); err != nil {
		return err
	}

	ctx := context.Background()
	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	_, err := runOCR(ctx, img, lang, []string{outputFile})
	return err
}

// DefaultClient provides a pre-configured client instance
var DefaultClient = &Client{
	config: Config{
		Timeout: 30 * time.Second,
	},
}
