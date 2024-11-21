package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/thedesertm/gotesseract/pkg/tesseract"
)

var (
	imagePath  = flag.String("image", "", "Path to input image")
	language   = flag.String("lang", "eng", "OCR language")
	timeout    = flag.Duration("timeout", 30*time.Second, "OCR timeout duration")
	outputType = flag.String("output", "text", "Output type: text, tsv, hocr, pdf")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] tesseract-path\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()

	if *imagePath == "" || len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	tesseractPath := args[0]
	client, err := tesseract.NewClient(tesseract.Config{
		TesseractPath: tesseractPath,
		Language:      *language,
		Timeout:       *timeout,
	})
	if err != nil {
		log.Fatalf("Failed to initialize OCR client: %v", err)
	}

	img, err := loadImage(*imagePath)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}

	switch *outputType {
	case "text":
		text, err := client.ImageToString(img, *language)
		if err != nil {
			log.Fatalf("OCR failed: %v", err)
		}
		fmt.Printf("Text Result:\n%s\n", text)

	case "tsv":
		output, err := client.ImageToExtension(img, *language, "tsv")
		if err != nil {
			log.Fatalf("TSV extraction failed: %v", err)
		}
		fmt.Printf("TSV Result:\n%s\n", output)

	case "hocr":
		output, err := client.ImageToExtension(img, *language, "hocr")
		if err != nil {
			log.Fatalf("HOCR extraction failed: %v", err)
		}
		fmt.Printf("HOCR Result:\n%s\n", output)

	case "pdf":
		output, err := client.ImageToExtension(img, *language, "pdf")
		if err != nil {
			log.Fatalf("PDF creation failed: %v", err)
		}
		fmt.Printf("PDF Result:\n%s\n", output)

	case "boxes":
		boxes, err := client.ImageToBoxes(img, *language)
		if err != nil {
			log.Fatalf("Box extraction failed: %v", err)
		}
		for _, box := range boxes {
			fmt.Printf("Character '%c' at (%d,%d,%d,%d) on page %d\n",
				box.Char, box.Left, box.Top, box.Right, box.Bottom, box.Page)
		}

	default:
		log.Fatalf("Unsupported output type: %s", *outputType)
	}
}

func loadImage(path string) (image.Image, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(absPath)
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
