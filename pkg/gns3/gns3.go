/*
Copyright 2022 Ciena Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gns3

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

type Gns3 struct {
}

func Connect() *Gns3 {
	return &Gns3{}
}

func (g *Gns3) Get(path string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("http://%s/%s", viper.GetString("address"), path), nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("insecure-skip-verify")},
	}
	client := &http.Client{Transport: tr}
	if err != nil {
		return fmt.Errorf("req: %w", err)
	}
	req.SetBasicAuth(viper.GetString("username"), viper.GetString("password"))
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}
	defer resp.Body.Close()
	if int(resp.StatusCode/100) != 2 {
		var httpErr HttpError
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&httpErr)
		if err != nil {
			return fmt.Errorf("error decode: %w", err)
		}
		return &httpErr
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
}

func (g *Gns3) Delete(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		fmt.Sprintf("http://%s/%s", viper.GetString("address"), path), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(viper.GetString("username"), viper.GetString("password"))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("insecure-skip-verify")},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	defer resp.Body.Close()
	if int(resp.StatusCode/100) != 2 {
		var httpErr HttpError
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&httpErr)
		if err != nil {
			return fmt.Errorf("error decode: %w", err)
		}
		return &httpErr
	}
	return nil
}

func (g *Gns3) Post(path string, contentType string, in interface{}, out interface{}) error {
	var req *http.Request
	var resp *http.Response
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
	defer cancel()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("insecure-skip-verify")},
	}
	client := &http.Client{Transport: tr}

	if in == nil {
		req, err = http.NewRequestWithContext(ctx, http.MethodPost,
			fmt.Sprintf("http://%s/%s", viper.GetString("address"), path), nil)
		if err != nil {
			return err
		}
		req.SetBasicAuth(viper.GetString("username"), viper.GetString("password"))
		resp, err = client.Do(req)
	} else {
		buf := new(bytes.Buffer)
		encoder := json.NewEncoder(buf)
		err = encoder.Encode(in)
		if err != nil {
			return fmt.Errorf("encoding: %w", err)
		}
		req, err = http.NewRequestWithContext(ctx, http.MethodPost,
			fmt.Sprintf("http://%s/%s", viper.GetString("address"), path), buf)
		if err != nil {
			return err
		}
		req.Header = map[string][]string{
			"content-type": {"application/json"},
		}
		req.SetBasicAuth(viper.GetString("username"), viper.GetString("password"))
		resp, err = client.Do(req)
	}
	if err != nil {
		return fmt.Errorf("sent: %w", err)
	}
	defer resp.Body.Close()
	if int(resp.StatusCode/100) != 2 {
		var httpErr HttpError
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&httpErr)
		if err != nil {
			return fmt.Errorf("error decode: %w", err)
		}
		return &httpErr
	}
	if out != nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(out)
		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}
	}
	return nil
}

func (g *Gns3) Put(path string, contentType string, in interface{}, out interface{}) error {
	var req *http.Request
	var resp *http.Response
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
	defer cancel()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: viper.GetBool("insecure-skip-verify")},
	}
	client := &http.Client{Transport: tr}

	if in == nil {
		req, err = http.NewRequestWithContext(ctx, http.MethodPut,
			fmt.Sprintf("http://%s/%s", viper.GetString("address"), path), nil)
		if err != nil {
			return err
		}
		req.SetBasicAuth(viper.GetString("username"), viper.GetString("password"))
		resp, err = client.Do(req)
	} else {
		buf := new(bytes.Buffer)
		encoder := json.NewEncoder(buf)
		err = encoder.Encode(in)
		if err != nil {
			return fmt.Errorf("encoding: %w", err)
		}
		req, err = http.NewRequestWithContext(ctx, http.MethodPut,
			fmt.Sprintf("http://%s/%s", viper.GetString("address"), path), buf)
		if err != nil {
			return err
		}
		req.Header = map[string][]string{
			"Content-type": {"application/json"},
		}
		req.SetBasicAuth(viper.GetString("username"), viper.GetString("password"))
		resp, err = client.Do(req)
	}
	if err != nil {
		return fmt.Errorf("sent: %w", err)
	}
	defer resp.Body.Close()
	if int(resp.StatusCode/100) != 2 {
		var httpErr HttpError
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&httpErr)
		if err != nil {
			return fmt.Errorf("error decode: %w", err)
		}
		return &httpErr
	}
	if out != nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(out)
		if err != nil {
			return fmt.Errorf("decode: %w", err)
		}
	}
	return nil
}
