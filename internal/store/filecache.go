package store

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/kmdkuk/nfrecap/internal/model"
)

type FileCache struct {
	dir string
}

func NewFileCache(dir string) *FileCache {
	_ = os.MkdirAll(dir, 0o755)
	return &FileCache{dir: dir}
}

func DefaultCacheDir() string {
	// cross-platform enough for MVP
	base, err := os.UserCacheDir()
	if err != nil || base == "" {
		return ".cache/nfrecap"
	}
	return filepath.Join(base, "nfrecap")
}

func (c *FileCache) Get(workTitle string, typ string) (model.Metadata, bool, error) {
	p := c.path(workTitle, typ)
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return model.Metadata{}, false, nil
		}
		return model.Metadata{}, false, err
	}
	var md model.Metadata
	if err := json.Unmarshal(b, &md); err != nil {
		return model.Metadata{}, false, err
	}
	return md, true, nil
}

func (c *FileCache) Put(workTitle string, typ string, md model.Metadata) error {
	p := c.path(workTitle, typ)
	b, err := json.MarshalIndent(md, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o644)
}

func (c *FileCache) path(workTitle string, typ string) string {
	key := strings.ToLower(strings.TrimSpace(typ)) + "|" + strings.TrimSpace(workTitle)
	h := sha1.Sum([]byte(key))
	name := hex.EncodeToString(h[:]) + ".json"
	return filepath.Join(c.dir, name)
}
