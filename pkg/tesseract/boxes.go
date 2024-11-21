package tesseract

import (
	"bufio"
	"context"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Box struct {
	Char   rune
	Left   int
	Bottom int
	Right  int
	Top    int
	Page   int
}

func (c *Client) ImageToBoxes(img image.Image, lang string) ([]Box, error) {
	tmpDir, err := createTempDir()
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	imgPath, err := saveImage(tmpDir, img)
	if err != nil {
		return nil, err
	}

	outBase := filepath.Join(tmpDir, "output")
	args := []string{
		imgPath,
		outBase,
		"-c", "tessedit_create_boxfile=1",
		"batch.nochop", "makebox",
	}

	if lang != "" {
		if err := validateLanguage(lang); err != nil {
			return nil, err
		}
		args = append(args, "-l", lang)
	}
	ctx := context.Background()
	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}
	// set context for the command
	cmd := exec.CommandContext(ctx, TesseractCmd, args...)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("tesseract error: %w", err)
	}

	boxFile := outBase + ".box"
	return parseBoxFile(boxFile)
}

func parseBoxFile(filename string) ([]Box, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var boxes []Box
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 6 {
			continue
		}

		char := []rune(fields[0])[0]
		left, _ := strconv.Atoi(fields[1])
		bottom, _ := strconv.Atoi(fields[2])
		right, _ := strconv.Atoi(fields[3])
		top, _ := strconv.Atoi(fields[4])
		page, _ := strconv.Atoi(fields[5])

		boxes = append(boxes, Box{
			Char:   char,
			Left:   left,
			Bottom: bottom,
			Right:  right,
			Top:    top,
			Page:   page,
		})
	}

	return boxes, scanner.Err()
}
