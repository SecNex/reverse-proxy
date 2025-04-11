package proxy

import (
	"log"
	"sync"

	"github.com/secnex/reverse-proxy/cert"
)

type ProxyConfig struct {
	Protocol string
	Host     string
	Port     int
	SSL      bool
}

type ConfigCache struct {
	configs     map[string]ProxyConfig
	mu          sync.RWMutex
	db          *DBManager
	certManager *cert.CertManager
}

func NewConfigCache(db *DBManager, certManager *cert.CertManager) *ConfigCache {
	return &ConfigCache{
		configs:     make(map[string]ProxyConfig),
		db:          db,
		certManager: certManager,
	}
}

func (cc *ConfigCache) Get(host string) (ProxyConfig, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	config, exists := cc.configs[host]
	return config, exists
}

func (cc *ConfigCache) Set(host string, config ProxyConfig) {
	log.Println("Setting config for host:", host)

	if config.SSL {
		cc.certManager.GenerateSelfSignedCert(host)
	}
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.configs[host] = config
}

func (cc *ConfigCache) Delete(host string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	delete(cc.configs, host)
}

func (cc *ConfigCache) GetAll() map[string]ProxyConfig {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	configs := make(map[string]ProxyConfig)
	for host, config := range cc.configs {
		configs[host] = config
	}
	return configs
}

func (cc *ConfigCache) Update(newConfigs map[string]ProxyConfig) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.configs = newConfigs
}

func (cc *ConfigCache) LoadFromDB() error {
	websites, err := cc.db.GetAllWebsites()
	if err != nil {
		return err
	}

	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.configs = make(map[string]ProxyConfig)
	for _, website := range websites {
		if website.Active {
			cc.configs[website.Domain] = ProxyConfig{
				Protocol: website.Protocol,
				Host:     website.Host,
				Port:     website.Port,
				SSL:      website.SSL,
			}
		}
	}
	return nil
}
