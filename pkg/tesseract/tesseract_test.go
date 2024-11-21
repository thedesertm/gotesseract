package tesseract

import (
	"image"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "Default config",
			config: Config{
				Timeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "Custom config",
			config: Config{
				Language: "eng",
				Timeout:  5 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "Invalid tesseract path",
			config: Config{
				TesseractPath: "/invalid/path",
				Timeout:       5 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestImageToString(t *testing.T) {
	client, err := NewClient(Config{
		Timeout: 30 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create test data directory if it doesn't exist
	testDataDir := "testdata"
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		t.Fatalf("Failed to create test data directory: %v", err)
	}

	// Test image path
	imagePath := filepath.Join(testDataDir, "sample.png")
	f, err := os.Open(imagePath)
	if err != nil {
		t.Skipf("Skipping test: test image not found at %s: %v", imagePath, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("Failed to decode test image: %v", err)
	}

	tests := []struct {
		name    string
		img     image.Image
		lang    string
		wantErr bool
	}{
		{
			name:    "Valid image",
			img:     img,
			lang:    "eng",
			wantErr: false,
		},
		{
			name:    "Nil image",
			img:     nil,
			lang:    "eng",
			wantErr: true,
		},
		{
			name:    "Invalid language",
			img:     img,
			lang:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := client.ImageToString(tt.img, tt.lang)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && text == "" {
				t.Error("ImageToString() returned empty string")
			}
		})
	}
}

func TestGetAvailableLanguages(t *testing.T) {
	langs, err := GetAvailableLanguages()
	if err != nil {
		t.Fatalf("GetAvailableLanguages() error = %v", err)
	}
	if len(langs) == 0 {
		t.Error("GetAvailableLanguages() returned empty list")
	}
}

func TestImageToExtension(t *testing.T) {
	client, err := NewClient(Config{
		Timeout: 30 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	img, err := loadTestImage(t)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		extension string
		wantErr   bool
	}{
		{"PDF", "pdf", false},
		{"HOCR", "hocr", false},
		{"Invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := client.ImageToExtension(img, "eng", tt.extension)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageToExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && output == "" {
				t.Error("ImageToExtension() returned empty output")
			}
		})
	}
}

// Helper function to load test image
func loadTestImage(t *testing.T) (image.Image, error) {
	t.Helper()
	f, err := os.Open(filepath.Join("testdata", "sample.png"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}
