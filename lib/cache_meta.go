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

	"github.com/mackerelio/golib/pluginutil"
	"github.com/natefinch/atomic"
)

type cache struct {
	Options  interface{} `json:"options"`
	Expected interface{} `json:"expected"`
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

func saveCache(file string, cache *cache) error {
	b, _ := json.Marshal(cache)
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}
	return atomic.WriteFile(file, bytes.NewReader(b))
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
