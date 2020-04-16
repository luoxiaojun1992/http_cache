package util

import (
	"strings"
)

const CR = "\r"
const LF = "\n"

func Serialize(mapData map[string][]string) string {
	var lines []string
	for key, values := range mapData {
		for _, value := range values {
			lines = append(lines, EscapeCRLF(key)+CR+LF+EscapeCRLF(value))
		}
	}
	return strings.Join(lines, strings.Repeat(CR+LF, 2))
}

func DeSerialize(data string) map[string]string {
	mapData := make(map[string]string)
	lines := strings.Split(data, strings.Repeat(CR+LF, 2))
	for _, line := range lines {
		pair := strings.Split(line, CR+LF)
		mapData[UnEscapeCRLF(pair[0])] = UnEscapeCRLF(pair[1])
	}
	return mapData
}

func EscapeCRLF(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "\r", "\\r"), "\n", "\\n")
}

func UnEscapeCRLF(escapedStr string) string {
	return strings.ReplaceAll(strings.ReplaceAll(escapedStr, "\\r", "\r"), "\\n", "\n")
}
