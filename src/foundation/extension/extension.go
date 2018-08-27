package extension

import . "github.com/luoxiaojun1992/http_cache/src/foundation/extension/concrete"

var extensions []extensionProto

func StartUp() {
	extensions = []extensionProto{&Registry{}}

	for _, extensionConcrete := range extensions {
		if extensionConcrete.IsEnabled() == 0 {
			continue
		}

		extensionConcrete.StartUp()
	}
}

func ShutDown() {
	for _, extensionConcrete := range extensions {
		if extensionConcrete.IsEnabled() == 0 {
			continue
		}

		extensionConcrete.ShutDown()
	}
}
