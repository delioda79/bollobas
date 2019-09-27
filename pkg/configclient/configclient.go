package configclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	patronHTTP "github.com/beatlabs/patron/trace/http"
)

// ConfigClient is a rest configuration client.
//
// It requires the endpoint and path of the rest endpoint
//
// On Success it will a bool field: mixpanel_setting, indicating if it's turned on or off
type ConfigClient struct {
	baseURL    string
	requestKey string
	path       string
	httpClient patronHTTP.Client
}

// New constructs a new Config Client
// the endpoint and path construct the final url
func New(baseURL string, key string, path string, opts ...patronHTTP.OptionFunc) (*ConfigClient, error) {
	httpClient, err := patronHTTP.New(opts...)
	if err != nil {
		return nil, err
	}
	return &ConfigClient{
		baseURL:    baseURL,
		httpClient: httpClient,
		requestKey: key,
		path:       path,
	}, nil
}

// GetSettings gets config settings from rest
func (dc *ConfigClient) GetSettings(ctx context.Context) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", dc.baseURL+dc.path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SECURE-ROUTE-KEY", dc.requestKey) //env me

	resp, err := dc.httpClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	newConf := map[string]interface{}{}
	err = json.Unmarshal(respBody, &newConf)
	if err != nil {
		return nil, err
	}

	return newConf, nil
}
