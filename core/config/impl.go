// Copyright 2021 magnifier Author.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package config

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type impl struct {
	data map[string]interface{}
	sync.RWMutex
}

func (c *impl) Set(key, value string) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = value
}

func (c *impl) String(key string) (string, error) {
	if v, err := c.getData(key); err == nil {
		if vv, ok := v.(string); ok {
			return vv, nil
		}
	}
	return "", nil
}

func (c *impl) DefaultString(key string, defaultVal string) string {
	if v, err := c.String(key); err == nil {
		return v
	}
	return defaultVal
}

func (c *impl) Strings(key string) ([]string, error) {
	v, err := c.String(key)
	if v == "" || err != nil {
		return nil, err
	}
	return strings.Split(v, ";"), nil
}

func (c *impl) DefaultStrings(key string, defaultVal []string) []string {
	if v, err := c.Strings(key); err == nil {
		return v
	}
	return defaultVal
}

func (c *impl) Int(key string) (int, error) {
	if v, err := c.getData(key); err != nil {
		return 0, err
	} else if vv, ok := v.(int); ok {
		return vv, nil
	} else if vv, ok := v.(int64); ok {
		return int(vv), nil
	}
	return 0, errors.New("not int value")
}

func (c *impl) DefaultInt(key string, defaultVal int) int {
	if v, err := c.Int(key); err == nil {
		return v
	}
	return defaultVal
}

func (c *impl) Int64(key string) (int64, error) {
	v, err := c.getData(key)
	if err != nil {
		return 0, err
	}
	switch val := v.(type) {
	case int:
		return int64(val), nil
	case int64:
		return val, nil
	default:
		return 0, errors.New("not int or int64 value")
	}
}

func (c *impl) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := c.Int64(key); err == nil {
		return v
	}
	return defaultVal
}

func (c *impl) Bool(key string) (bool, error) {
	val, err := c.getData(key)
	if err != nil {
		return false, err
	}
	if val != nil {
		switch v := val.(type) {
		case bool:
			return v, nil
		case string:
			switch v {
			case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
				return true, nil
			case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
				return false, nil
			}
		case int8, int32, int64:
			strV := fmt.Sprintf("%d", v)
			if strV == "1" {
				return true, nil
			} else if strV == "0" {
				return false, nil
			}
		case float64:
			if v == 1.0 {
				return true, nil
			} else if v == 0.0 {
				return false, nil
			}
		}
		return false, fmt.Errorf("parsing %q: invalid syntax", val)
	}
	return false, fmt.Errorf("parsing <nil>: invalid syntax")
}

func (c *impl) DefaultBool(key string, defaultVal bool) bool {
	if v, err := c.Bool(key); err == nil {
		return v
	}
	return defaultVal
}

func (c *impl) Float(key string) (float64, error) {
	if v, err := c.getData(key); err != nil {
		return 0.0, err
	} else if vv, ok := v.(float64); ok {
		return vv, nil
	} else if vv, ok := v.(int); ok {
		return float64(vv), nil
	} else if vv, ok := v.(int64); ok {
		return float64(vv), nil
	}
	return 0.0, errors.New("not float64 value")
}

func (c *impl) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := c.Float(key); err == nil {
		return v
	}
	return defaultVal
}

func (c *impl) DIY(key string) (interface{}, error) {
	return c.getData(key)
}

func (c *impl) Sub(key string) (Config, error) {
	sub, err := c.subMap(key)
	if err != nil {
		return nil, err
	}
	return &impl{
		data: sub,
	}, nil
}

func (c *impl) subMap(key string) (map[string]interface{}, error) {
	tmpData := c.data
	keys := strings.Split(key, ".")
	for idx, k := range keys {
		if v, ok := tmpData[k]; ok {
			switch val := v.(type) {
			case map[string]interface{}:
				tmpData = val
				if idx == len(keys)-1 {
					return tmpData, nil
				}
			default:
				return nil, fmt.Errorf("the key is invalid: %s", key)
			}
		}
	}

	return tmpData, nil
}

func (c *impl) getData(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}
	c.RLock()
	defer c.RUnlock()

	keys := strings.Split(key, ".")
	tmpData := c.data
	for idx, k := range keys {
		if v, ok := tmpData[k]; ok {
			switch val := v.(type) {
			case map[string]interface{}:
				tmpData = val
				if idx == len(keys)-1 {
					return tmpData, nil
				}
			default:
				return v, nil
			}
		}
	}
	return nil, fmt.Errorf("not exist key %q", key)
}
