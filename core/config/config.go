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

type Config interface {
	Set(key, val string)

	// support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	String(key string) (string, error)
	// get string slice
	Strings(key string) ([]string, error)
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	// support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	DefaultString(key string, defaultVal string) string
	// get string slice
	DefaultStrings(key string, defaultVal []string) []string
	DefaultInt(key string, defaultVal int) int
	DefaultInt64(key string, defaultVal int64) int64
	DefaultBool(key string, defaultVal bool) bool
	DefaultFloat(key string, defaultVal float64) float64

	// DIY return the original value
	DIY(key string) (interface{}, error)

	Sub(key string) (Config, error)
}
