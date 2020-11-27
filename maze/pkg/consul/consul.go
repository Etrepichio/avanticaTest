package consul



import (
	"github.com/pkg/errors"
	stdconsul "github.com/hashicorp/consul/api"
)


func SafeConsulGet(consul *stdconsul.KV, key string) ([]byte, error) {
	kvp, _, err := consul.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kvp == nil {
		return nil, errors.Errorf("consul missing key: %v", key)
	}
	return kvp.Value, nil
}


func OpenConsul(addr string) (*stdconsul.KV, error) {
	consulConfig := stdconsul.DefaultConfig()
	if len(addr) > 0 {
		consulConfig.Address = addr
	}
	consulClient, err := stdconsul.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}
	return consulClient.KV(), nil
}