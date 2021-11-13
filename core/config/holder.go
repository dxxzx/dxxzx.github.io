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
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var holder Config

func FromYaml(filename string) error {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return FromYamlData(raw)
}

func FromYamlData(raw []byte) error {
	data := make(map[string]interface{})
	err := yaml.Unmarshal(raw, &data)
	if err != nil {
		return err
	}
	holder = &impl{data: ExpandValueForMap(data)}
	return nil
}

// ExpandValueForMap convert all string value.
func ExpandValueForMap(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		switch value := v.(type) {
		case map[string]interface{}:
			m[k] = ExpandValueForMap(value)
		case map[interface{}]interface{}:
			tmp := make(map[string]interface{}, len(value))
			for k2, v2 := range value {
				tmp[k2.(string)] = v2
			}
			m[k] = ExpandValueForMap(tmp)
		}
	}
	return m
}

func String(key string) (string, error) {
	return holder.String(key)
}

func DefaultString(key string, defaultVal string) string {
	return holder.DefaultString(key, defaultVal)
}

func Strings(key string) ([]string, error) {
	return holder.Strings(key)
}

func DefaultStrings(key string, defaultVal []string) []string {
	return holder.DefaultStrings(key, defaultVal)
}

func Int(key string) (int, error) {
	return holder.Int(key)
}

func DefaultInt(key string, defaultVal int) int {
	return holder.DefaultInt(key, defaultVal)
}

func Int64(key string) (int64, error) {
	return holder.Int64(key)
}

func DefaultInt64(key string, defaultVal int64) int64 {
	return holder.DefaultInt64(key, defaultVal)
}

func Bool(key string) (bool, error) {
	return holder.Bool(key)
}

func DefaultBool(key string, defaultVal bool) bool {
	return holder.DefaultBool(key, defaultVal)
}

func Float(key string) (float64, error) {
	return holder.Float(key)
}

func DefaultFloat(key string, defaultVal float64) float64 {
	return holder.DefaultFloat(key, defaultVal)
}

func DIY(key string) (interface{}, error) {
	return holder.DIY(key)
}

func Sub(key string) (Config, error) {
	return holder.Sub(key)
}
