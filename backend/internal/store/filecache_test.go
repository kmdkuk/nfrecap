package store

import (
	"testing"
	"time"

	"github.com/kmdkuk/nfrecap/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileCache_PutAndGet(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	cache := NewFileCache(tmpDir, 0) // no TTL

	workTitle := "Inception"
	typ := "movie"
	expectedMD := model.Metadata{
		Title:   "Inception",
		Runtime: 148,
		Genres:  []string{"Action", "Sci-Fi"},
	}

	// 1. Put
	err := cache.Put(workTitle, typ, expectedMD)
	require.NoError(t, err)

	// 2. Get (Hit)
	gotMD, found, err := cache.Get(workTitle, typ)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, expectedMD, gotMD)

	// 3. Get (Miss)
	_, found, err = cache.Get("Unknown", "movie")
	require.NoError(t, err)
	assert.False(t, found)

	// 4. Persistence check (new instance pointing to same dir)
	cache2 := NewFileCache(tmpDir, 0)
	gotMD2, found2, err := cache2.Get(workTitle, typ)
	require.NoError(t, err)
	assert.True(t, found2)
	assert.Equal(t, expectedMD, gotMD2)
}

func TestFileCache_TTL(t *testing.T) {
	tmpDir := t.TempDir()
	ttl := 100 * time.Millisecond
	cache := NewFileCache(tmpDir, ttl)

	workTitle := "ShortLived"
	typ := "movie"
	md := model.Metadata{Title: "ShortLived"}

	// Put
	require.NoError(t, cache.Put(workTitle, typ, md))

	// Get Immedidate
	_, found, err := cache.Get(workTitle, typ)
	require.NoError(t, err)
	assert.True(t, found)

	// Wait for expiration
	time.Sleep(200 * time.Millisecond)

	// Get Expired
	_, found, err = cache.Get(workTitle, typ)
	require.NoError(t, err)
	assert.False(t, found, "should be expired")
}

func TestFileCache_Sanitization(t *testing.T) {
	tmpDir := t.TempDir()
	cache := NewFileCache(tmpDir, 0)

	// Title with slash should be sanitized to avoid directory issues
	workTitle := "Face/Off" // Slash in title
	typ := "movie"
	md := model.Metadata{Title: "Face/Off"}

	err := cache.Put(workTitle, typ, md)
	require.NoError(t, err)

	gotMD, found, err := cache.Get(workTitle, typ)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, md, gotMD)
}
