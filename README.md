# GoTesseract

A Go port of [pytesseract](https://github.com/madmaze/pytesseract). GoTesseract provides Go bindings for Google's Tesseract-OCR Engine.

## Installation

```bash
go get github.com/thedesertm/gotesseract
```

### Prerequisites
- Tesseract >= 3.05
- Go >= 1.16

## Usage

```go
import "github.com/thedesertm/gotesseract"

// Initialize client
client, err := tesseract.NewClient(tesseract.Config{
    TesseractPath: "C:\\Program Files\\Tesseract-OCR\\tesseract.exe", // Windows path
    Language:      "eng",
    Timeout:       30 * time.Second,
})

// Basic text extraction
text, err := client.ImageToString(img, "eng")

// Get bounding boxes
boxes, err := client.ImageToBoxes(img, "eng")
for _, box := range boxes {
    fmt.Printf("Character '%c' at (%d,%d,%d,%d)\n", 
        box.Char, box.Left, box.Top, box.Right, box.Bottom)
}

// Get other formats
hocr, err := client.ImageToExtension(img, "eng", "hocr")
pdf, err := client.ImageToExtension(img, "eng", "pdf")
tsv, err := client.ImageToExtension(img, "eng", "tsv")
```

## Command Line Example

```bash
go run cmd/example/main.go -image=sample.png -output=boxes "C:\Program Files\Tesseract-OCR\tesseract.exe"
```

## Features
- Text extraction
- Bounding box detection
- Multiple output formats (Text, hOCR, PDF, TSV)
- Configurable timeouts
- Custom Tesseract path
- Language selection

## License
MIT License - see [LICENSE](LICENSE)

## Credits
Port of [pytesseract](https://github.com/madmaze/pytesseract) to Go.
