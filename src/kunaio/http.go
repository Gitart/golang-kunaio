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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"
)

var (
	gClient *http.Client
)

// Module initialization hook.
func init() {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Minute,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          5,
		IdleConnTimeout:       35 * time.Minute,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	gClient = &http.Client{
		Transport: transport,
		Timeout:   4 * time.Second,
	}
}

// Send a GET request to the server.
func doGet(url string) (interface{}, error) {
	resp, err := gClient.Get(url)
	if err != nil {
		return nil, err
	}
	return readResp(resp)
}

// Send a POST request to the server.
func doPost(url string) (interface{}, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return Order{}, err
	}
	resp, err := gClient.Do(req)
	if err != nil {
		return Order{}, err
	}
	return readResp(resp)
}

// Check HTTP response and decode JSON object.
func readResp(r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	if r.StatusCode/100 != 2 {
		if err, ok := decodeError(r.Body); ok {
			return nil, fmt.Errorf("%s; %s", r.Status, err)
		}
		return nil, fmt.Errorf("server returned HTTP"+
			" response code %s", r.Status)
	}
	return DecodeJSON(r.Body)
}

type Args []struct {
	Key   string
	Value string
}

func (a Args) Len() int {
	return len(a)
}

func (a Args) Less(i, j int) bool {
	return a[i].Key < a[j].Key
}

func (a Args) Swap(i, j int) {
	t := a[i]
	a[i] = a[j]
	a[j] = t
}

// Generate URL for private API request.
func privURL(method, url, access_key, secret_key string, args Args) string {
	if args == nil {
		args = Args{}
	}
	args = append(args, Args{
		{"access_key", access_key},
		{"tonce", fmt.Sprintf("%d000", time.Now().Unix())},
	}...)
	sort.Sort(args)
	query := ""
	for _, e := range args {
		query += fmt.Sprintf("&%s=%s", e.Key, e.Value)
	}
	query = strings.Trim(query, "&")
	h := hmac.New(sha256.New, []byte(secret_key))
	h.Write([]byte(fmt.Sprintf("%s|%s|%s", method, url, query)))
	secret := hex.EncodeToString(h.Sum(nil))
	url = fmt.Sprintf("%s?%s&signature=%s", url, query, secret)
	return fmt.Sprintf("%s%s", gBaseURL, url)
}
