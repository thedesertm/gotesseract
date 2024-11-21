package tesseract

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var versionRegex = regexp.MustCompile(`tesseract\s+v?(\d+)\.(\d+)\.(\d+)(?:\.\d+)?`)

func checkVersion() error {
	cmd := exec.Command(TesseractCmd, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get tesseract version: %w", err)
	}

	// Extract first line which contains version
	firstLine := strings.SplitN(string(out), "\n", 2)[0]
	matches := versionRegex.FindStringSubmatch(firstLine)
	if matches == nil {
		return fmt.Errorf("unrecognized tesseract version format: %s", firstLine)
	}

	// Parse version numbers
	major := matches[1]
	minor := matches[2]

	// Check minimum version requirement
	if major < "3" || (major == "3" && minor < "05") {
		return fmt.Errorf("tesseract version %s.%s is not supported (minimum 3.05 required)", major, minor)
	}

	return nil
}

func createTempDir() (string, error) {
	dir, err := os.MkdirTemp("", "tesseract_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	return dir, nil
}

func saveImage(dir string, img image.Image) (string, error) {
	outPath := filepath.Join(dir, "input.png")
	f, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	if _, err := f.Write(buf.Bytes()); err != nil {
		return "", err
	}

	return outPath, nil
}

func GetAvailableLanguages() ([]string, error) {
	cmd := exec.Command(TesseractCmd, "--list-langs")
	out, err := cmd.Output()
	if err != nil {
		return nil, ErrTesseractNotFound
	}
	// this based on the os for windows it is \r\n and for linux it is \n
	if strings.Contains(string(out), "\r\n") {
		out = bytes.ReplaceAll(out, []byte("\r\n"), []byte("\n"))
	}
	langs := strings.Split(string(out), "\n")
	var langsOutput []string

	for i, l := range langs {
		if i == 0 {
			continue
		}
		if len(l) > 0 {
			langsOutput = append(langsOutput, l)
		}
	}
	return langsOutput, nil
}

func runOCR(ctx context.Context, img image.Image, lang string, args []string) ([]byte, error) {
	tmpDir, err := createTempDir()
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	imgPath, err := saveImage(tmpDir, img)
	if err != nil {
		return nil, err
	}

	cmdArgs := append([]string{imgPath}, args...)
	if lang != "" {
		cmdArgs = append(cmdArgs, "-l", lang)
	}

	cmd := exec.CommandContext(ctx, TesseractCmd, cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, &OCRError{Code: exitErr.ExitCode(), Stderr: string(exitErr.Stderr)}
		}
		return nil, err
	}

	return out, nil
}

func validateImageFormat(img image.Image) error {
	if img == nil {
		return errors.New("nil image")
	}
	return nil
}

func validateLanguage(lang string) error {
	langs, err := GetAvailableLanguages()
	if err != nil {
		return err
	}
	for _, l := range langs {
		if l == lang {
			return nil
		}
	}
	return ErrLanguageNotFound
}
