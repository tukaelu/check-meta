package checkmeta

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mackerelio/golib/pluginutil"
	"github.com/natefinch/atomic"
)

type cache struct {
	Options   interface{} `json:"options"`
	Expected  interface{} `json:"expected"`
	UpdatedAt int64       `json:"updated_at"`
}

func getCacheFile(args []string) string {
	return filepath.Join(
		filepath.Join(pluginutil.PluginWorkDir(), "check-meta"),
		fmt.Sprintf(
			"check-meta-%x.json",
			md5.Sum([]byte(strings.Join(args, " "))),
		),
	)
}

func saveCache(file string, cache *cache) (*cache, error) {
	cache.UpdatedAt = time.Now().Unix()
	b, _ := json.Marshal(cache)
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return nil, err
	}
	return cache, atomic.WriteFile(file, bytes.NewReader(b))
}

func loadCache(file string) (*cache, error) {
	cache := &cache{}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return cache, err
	}
	err = json.Unmarshal(b, cache)
	if err != nil {
		return nil, fmt.Errorf("cache file is corrupted: %s", err.Error())
	}
	return cache, nil
}
