package proxy

import (
	"sync"
)

type ProxyConfig struct {
	Protocol string
	Host     string
	Port     int
	SSL      bool
}

type ConfigCache struct {
	configs map[string]ProxyConfig
	mu      sync.RWMutex
}

func NewConfigCache() *ConfigCache {
	return &ConfigCache{
		configs: make(map[string]ProxyConfig),
	}
}

func (cc *ConfigCache) Get(host string) (ProxyConfig, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	config, exists := cc.configs[host]
	return config, exists
}

func (cc *ConfigCache) Set(host string, config ProxyConfig) {
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
	// Erstelle eine Kopie der Konfigurationen
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
