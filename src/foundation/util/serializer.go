package util

import (
	"strings"
)

const CRLF = "\r\n"
const ESCAPED_CRLF  = "\\r\\n"

func Serialize(mapData map[string][]string) string {
	var lines []string
	for key, values := range mapData {
		for _, value := range values {
			lines = append(lines, EscapeCRLF(key)+CRLF+EscapeCRLF(value))
		}
	}
	return strings.Join(lines, strings.Repeat(CRLF, 2))
}

func DeSerialize(data string) map[string]string {
	mapData := make(map[string]string)
	lines := strings.Split(data, strings.Repeat(CRLF, 2))
	for _, line := range lines {
		pair := strings.Split(line, CRLF)
		mapData[UnEscapeCRLF(pair[0])] = UnEscapeCRLF(pair[1])
	}
	return mapData
}

func EscapeCRLF(str string) string {
	return strings.Replace(str, CRLF, ESCAPED_CRLF, -1)
}

func UnEscapeCRLF(escapedStr string) string {
	return strings.Replace(escapedStr, ESCAPED_CRLF, CRLF, -1)
}
