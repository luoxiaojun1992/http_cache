package util

import (
	"strings"
)

func Serialize(map_data map[string][]string) string {
	lines := []string{}
	for key, values := range map_data {
		for _, value := range values {
			lines = append(lines, key+"\r\n"+value)
		}
	}
	return strings.Join(lines, "\r\n\r\n")
}

func DeSerialize(data string) map[string]string {
	map_data := make(map[string]string)
	lines := strings.Split(data, "\r\n\r\n")
	for _, line := range lines {
		pair := strings.Split(line, "\r\n")
		map_data[pair[0]] = pair[1]
	}
	return map_data
}
