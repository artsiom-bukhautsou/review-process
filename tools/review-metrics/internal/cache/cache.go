// Package cache is a tiny on-disk cache of raw API responses keyed by request URL.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Cache stores response bodies as files named by the SHA-256 of their key.
// A nil *Cache is valid and caches nothing.
type Cache struct {
	dir string
}

// New returns a cache rooted at dir, creating it if needed.
func New(dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("creating cache dir %s: %w", dir, err)
	}
	return &Cache{dir: dir}, nil
}

// DefaultDir returns $XDG_CACHE_HOME/review-metrics (or the OS user cache dir).
func DefaultDir() (string, error) {
	if x := os.Getenv("XDG_CACHE_HOME"); x != "" {
		return filepath.Join(x, "review-metrics"), nil
	}
	d, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("locating user cache dir: %w", err)
	}
	return filepath.Join(d, "review-metrics"), nil
}

func (c *Cache) path(key string) string {
	sum := sha256.Sum256([]byte(key))
	h := hex.EncodeToString(sum[:])
	return filepath.Join(c.dir, h[:2], h+".json")
}

// Get returns the cached body for key, if present.
func (c *Cache) Get(key string) ([]byte, bool) {
	if c == nil {
		return nil, false
	}
	b, err := os.ReadFile(c.path(key))
	if err != nil {
		return nil, false
	}
	return b, true
}

// Put stores body under key. Writes are atomic (temp file + rename).
func (c *Cache) Put(key string, body []byte) error {
	if c == nil {
		return nil
	}
	p := c.path(key)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return fmt.Errorf("creating cache subdir: %w", err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(p), ".tmp-*")
	if err != nil {
		return fmt.Errorf("creating cache temp file: %w", err)
	}
	if _, err := tmp.Write(body); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return fmt.Errorf("writing cache entry: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return fmt.Errorf("closing cache entry: %w", err)
	}
	if err := os.Rename(tmp.Name(), p); err != nil && !errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("committing cache entry: %w", err)
	}
	return nil
}
