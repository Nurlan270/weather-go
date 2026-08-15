package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Nurlan270/weather-go/internal/config"
	"github.com/Nurlan270/weather-go/internal/logger"

	bc "github.com/allegro/bigcache"
)

type Cache struct {
	bc *bc.BigCache
}

func New(cfg config.Cache, logger *logger.Logger) *Cache {
	cacheCfg := bc.Config{
		Shards:           512,
		LifeWindow:       cfg.ExpiresIn,
		CleanWindow:      2 * time.Second,
		Verbose:          true,
		Logger:           logger,
		HardMaxCacheSize: 150 << 20, // 150 MiB
	}

	bigCache, err := bc.NewBigCache(cacheCfg)
	if err != nil {
		log.Fatal(err)
	}

	cache := &Cache{
		bc: bigCache,
	}

	return cache
}

// Set is a wrapper around bigcache.BigCache's Set method.
// It provides serialization using json package.
func (c *Cache) Set(key string, value any) error {
	entry, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize entry: %w", err)
	}

	if err = c.bc.Set(key, entry); err != nil {
		return fmt.Errorf("failed to set entry: %w", err)
	}

	return nil
}

// GetInto is a wrapper around bigcache.BigCache's Get method.
// It deserializes value under key into dst using json package; dst should be a reference.
func (c *Cache) GetInto(key string, dst any) error {
	value, err := c.bc.Get(key)
	if err != nil {
		return fmt.Errorf("value with %s key does not exists in cache: %w", key, err)
	}

	if err = json.Unmarshal(value, dst); err != nil {
		return fmt.Errorf("failed to unserialize into destination: %w", err)
	}

	return nil
}

func (c *Cache) Del(key string) error {
	if err := c.bc.Delete(key); err != nil {
		return err
	}

	return nil
}

func BuildKey(key string, a ...any) string {
	return fmt.Sprintf(key, a...)
}
