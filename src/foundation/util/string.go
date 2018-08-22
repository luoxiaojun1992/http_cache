package util

import (
	"strings"
)

func Serialize(mapData map[string][]string) string {
	var lines []string
	for key, values := range mapData {
		for _, value := range values {
			lines = append(lines, key+"\r\n"+value)
		}
	}
	return strings.Join(lines, "\r\n\r\n")
}

func DeSerialize(data string) map[string]string {
	mapData := make(map[string]string)
	lines := strings.Split(data, "\r\n\r\n")
	for _, line := range lines {
		pair := strings.Split(line, "\r\n")
		mapData[pair[0]] = pair[1]
	}
	return mapData
}
