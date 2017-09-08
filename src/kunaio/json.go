// Copyright 2017 Aleksey Morarash <tuxofil@gmail.com>
//
// Licensed under the BSD 2 Clause License (the "License");
// you may not use the file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://opensource.org/licenses/BSD-2-Clause
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kunaio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

func DecodeJSON(r io.Reader) (interface{}, error) {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()
	var v interface{}
	if err := decoder.Decode(&v); err != nil {
		return nil, err
	}
	debugLog("decoded JSON: %#v", v)
	return v, nil
}

func jsonGetTime(v interface{}) (time.Time, error) {
	if v == nil {
		return time.Time{}, errors.New(
			"expected int but NIL found")
	}
	switch v.(type) {
	case json.Number:
		t, err := v.(json.Number).Int64()
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(t, 0), nil
	}
	return time.Time{}, fmt.Errorf(
		"expected int but %#v (%T) found", v, v)
}

var supportedTimeFormats = []string{
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05-0700",
	"2006-01-02T15:04:05-07:00",
}

func jsonGetTimeFromText(v interface{}) (time.Time, error) {
	if v == nil {
		return time.Time{}, errors.New(
			"expected int but NIL found")
	}
	switch v.(type) {
	case string:
		for _, f := range supportedTimeFormats {
			t, err := time.Parse(f, v.(string))
			if err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("can't parse time: %s",
			v.(string))
	}
	return time.Time{}, fmt.Errorf(
		"expected time but %#v (%T) found", v, v)
}

func jsonGetString(v interface{}) (string, error) {
	if v == nil {
		return "", errors.New(
			"expected string but NIL found")
	}
	switch v.(type) {
	case string:
		return v.(string), nil
	}
	return "", fmt.Errorf(
		"expected string but %#v (%T) found", v, v)
}

func jsonGetBool(v interface{}) (bool, error) {
	if v == nil {
		return false, errors.New(
			"expected bool but NIL found")
	}
	switch v.(type) {
	case bool:
		return v.(bool), nil
	}
	return false, fmt.Errorf(
		"expected bool but %#v (%T) found", v, v)
}

func jsonGetFloat(v interface{}) (float64, error) {
	if v == nil {
		return 0, errors.New("expected float but NIL found")
	}
	switch v.(type) {
	case json.Number:
		return v.(json.Number).Float64()
	case string:
		return strconv.ParseFloat(v.(string), 64)
	}
	return 0, fmt.Errorf(
		"expected float but %#v (%T) found", v, v)
}

func jsonGetFloatDef(v interface{}, def float64) (float64, error) {
	if v == nil {
		return def, nil
	}
	switch v.(type) {
	case json.Number:
		return v.(json.Number).Float64()
	case string:
		return strconv.ParseFloat(v.(string), 64)
	}
	return 0, fmt.Errorf(
		"expected float but %#v (%T) found", v, v)
}

func jsonGetInt(v interface{}) (int, error) {
	if v == nil {
		return 0, errors.New("expected int but NIL found")
	}
	switch v.(type) {
	case json.Number:
		i, err := v.(json.Number).Int64()
		if err != nil {
			return 0, err
		}
		return int(i), nil
	}
	return 0, fmt.Errorf(
		"expected int but %#v (%T) found", v, v)
}

func jsonGetMap(v interface{}) (map[string]interface{}, error) {
	if v == nil {
		return nil, errors.New(
			"expected dict but NIL found")
	}
	switch v.(type) {
	case map[string]interface{}:
		return v.(map[string]interface{}), nil
	}
	return nil, fmt.Errorf(
		"expected dict but %#v (%T) found", v, v)
}

func jsonGetList(v interface{}) ([]interface{}, error) {
	if v == nil {
		return nil, errors.New(
			"expected list but NIL found")
	}
	switch v.(type) {
	case []interface{}:
		return v.([]interface{}), nil
	}
	return nil, fmt.Errorf(
		"expected list but %#v (%T) found", v, v)
}
