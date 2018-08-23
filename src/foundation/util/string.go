package util

import (
	"strings"
)

const LF = "\r\n"

func Serialize(mapData map[string][]string) string {
	var lines []string
	for key, values := range mapData {
		for _, value := range values {
			lines = append(lines, key+LF+value)
		}
	}
	return strings.Join(lines, strings.Repeat(LF, 2))
}

func DeSerialize(data string) map[string]string {
	mapData := make(map[string]string)
	lines := strings.Split(data, strings.Repeat(LF, 2))
	for _, line := range lines {
		pair := strings.Split(line, LF)
		mapData[pair[0]] = pair[1]
	}
	return mapData
}
