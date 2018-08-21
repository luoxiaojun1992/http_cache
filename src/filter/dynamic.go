package filter

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
	dynamic_tags := re.FindAllString(body, -1)
	if len(dynamic_tags) <= 0 {
		return body
	}

	dynamic_contents := make(map[string]string)
	var wg sync.WaitGroup
	var mutex_lock sync.Mutex
	for i, dynamic_tag := range dynamic_tags {
		if i >= 10 {
			break
		}
		go func() {
			defer wg.Done()

			dynamic_url := strings.Replace(dynamic_tag, "<dynamic>", "", 1)
			dynamic_url = strings.Replace(dynamic_url, "</dynamic>", "", 1)
			if val, ok := dynamic_contents[dynamic_url]; ok {
				mutex_lock.Lock()
				body = strings.Replace(body, dynamic_tag, val, 1)
				mutex_lock.Unlock()
				return
			}

			resp, err := http.Get(dynamic_url)
			if err == nil {
				defer resp.Body.Close()

				dynamic_content, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					dynamic_content_str := string(dynamic_content)
					dynamic_contents[dynamic_url] = dynamic_content_str
					mutex_lock.Lock()
					body = strings.Replace(body, dynamic_tag, dynamic_content_str, 1)
					mutex_lock.Unlock()
				}
			}
		}()
		wg.Add(1)
	}
	wg.Wait()

	return body
}
