package tesseract

import (
	"strings"
)

func parseTSV(data string) map[string][]string {
	result := make(map[string][]string)
	rows := strings.Split(strings.TrimSpace(data), "\n")
	if len(rows) < 2 {
		return result
	}

	headers := strings.Split(rows[0], "\t")
	for _, header := range headers {
		result[header] = make([]string, 0)
	}

	for _, row := range rows[1:] {
		cols := strings.Split(row, "\t")
		for i, col := range cols {
			if i < len(headers) {
				result[headers[i]] = append(result[headers[i]], col)
			}
		}
	}
	return result
}
