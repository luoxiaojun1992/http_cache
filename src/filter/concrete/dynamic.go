package filter_concrete

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type DynamicContent struct {
}

func (dc *DynamicContent) Handle(body string) string {
	if !strings.Contains(body, "<dynamic>") {
		return body
	}

	re := regexp.MustCompile(`\<dynamic\>.+\</dynamic\>`)
	dynamicTags := re.FindAllString(body, -1)
	if len(dynamicTags) <= 0 {
		return body
	}

	dynamicContents := make(map[string]string)
	var wg sync.WaitGroup
	var mutexLock sync.Mutex
	for i, dynamicTag := range dynamicTags {
		if i >= 10 {
			break
		}
		go func() {
			defer wg.Done()

			dynamicUrl := strings.Replace(dynamicTag, "<dynamic>", "", 1)
			dynamicUrl = strings.Replace(dynamicUrl, "</dynamic>", "", 1)
			if val, ok := dynamicContents[dynamicUrl]; ok {
				mutexLock.Lock()
				body = strings.Replace(body, dynamicTag, val, 1)
				mutexLock.Unlock()
				return
			}

			resp, err := http.Get(dynamicUrl)
			if err == nil {
				defer resp.Body.Close()

				dynamicContent, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					dynamicContentStr := string(dynamicContent)
					dynamicContents[dynamicUrl] = dynamicContentStr
					mutexLock.Lock()
					body = strings.Replace(body, dynamicTag, dynamicContentStr, 1)
					mutexLock.Unlock()
				}
			}
		}()
		wg.Add(1)
	}
	wg.Wait()

	return body
}

func (dc *DynamicContent) IsRequest() bool {
	return false
}
