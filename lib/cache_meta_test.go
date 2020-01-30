package checkmeta

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadCache(t *testing.T) {
	file := "testdata/valid_cache.json"
	c, err := loadCache(file)

	assert.NoErrorf(t, err, "loadCache(%q) = %v; want no error", file, err)
	assert.NotNilf(t, *c, "loadCache(%q) = %v; cache want not nil", file, *c)
}

func TestLoadInvalidCache(t *testing.T) {
	file := "testdata/invalid_cache.json"
	c, err := loadCache(file)

	assert.NoErrorf(t, err, "loadCache(%q) = %v; want no error", file, err)
	assert.Nilf(t, c.Expected, "loadCache(%q) = %v; cache.Expected want nil", file, *c)
}

func TestLoadBrokenCache(t *testing.T) {
	file := "testdata/broken_cache.json"
	c, err := loadCache(file)

	assert.Errorf(t, err, "loadCache(%q) = %v; want error", file, err)
	assert.Nilf(t, c, "loadCache(%q) = %v; cache want nil", file, c)
}

func TestLoadCacheIfFileNotExist(t *testing.T) {
	file := "testdata/cache_file_does_not_exist"
	_, err := loadCache(file)

	assert.NoErrorf(t, err, "loadCache(%q) = %v; want no error", file, err)
}

func TestSaveCacheString(t *testing.T) {
	file := "testdata/saved_cache"
	defer func() {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}()

	c := &cache{
		Options:  []string{"--namespace", "foobar", "--key", "version"},
		Expected: "string_value",
	}
	saveAndLoad(t, file, c)
}

func TestSaveCacheNumber(t *testing.T) {
	file := "testdata/saved_cache"
	defer func() {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}()

	c := &cache{
		Options:  []string{"--namespace", "foobar", "--key", "version"},
		Expected: float64(1000),
	}
	saveAndLoad(t, file, c)
}

func TestSaveCacheBoolean(t *testing.T) {
	file := "testdata/saved_cache"
	defer func() {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}()

	c := &cache{
		Options:  []string{"--namespace", "foobar", "--key", "version"},
		Expected: true,
	}
	saveAndLoad(t, file, c)
}

func TestOverwriteSaveCache(t *testing.T) {
	file := "testdata/saved_cache"
	defer func() {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}()

	err := ioutil.WriteFile(file, []byte(`{"dummy": "dummy"}`), 0644)
	if err != nil {
		t.Errorf("WriteFile: %v", err)
		return
	}
	c := &cache{
		Options:  []string{"--namespace", "foobar", "--key", "version"},
		Expected: "string_value",
	}
	saveAndLoad(t, file, c)
}

func saveAndLoad(t *testing.T, file string, c *cache) {
	t.Helper()

	var c0 *cache
	var err error
	if c0, err = saveCache(file, c); err != nil {
		t.Errorf("saveCache(%v) = %v; want nil", file, *c)
		return
	}
	c1, err := loadCache(file)
	if err != nil {
		t.Errorf("loadCache: %v", err)
		return
	}
	if !reflect.DeepEqual(c0, c1) {
		t.Errorf("saveCache(%v) -> loadCache() = %v", c, c1)
	}
}
