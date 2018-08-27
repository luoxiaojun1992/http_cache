package extension_concrete

import . "github.com/luoxiaojun1992/http_cache/src/foundation/environment"

type Registry struct {
}

func (r *Registry) StartUp() {
	//
}

func (r *Registry) ShutDown() {
	//
}

func (r *Registry) IsEnabled() int {
	return EnvInt("REGISTRY_SWITCH", 0)
}
