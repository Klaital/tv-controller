package vlcclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Addr         string // in the form of host:port, e.g. "localhost:8090"
	HttpUser     string // it seems VLC ensures this is always empty string
	HttpPassword string // VLC requires this to be set at runtime
}

func (c Client) constructUrl(endpoint string, queryParams map[string]string) string {
	urlBuilder := strings.Builder{}
	urlBuilder.WriteString(c.Addr)
	urlBuilder.WriteString(endpoint)
	if len(queryParams) > 0 {
		urlBuilder.WriteString("?")
		queryTokens := make([]string, 0, len(queryParams))
		// TODO: sort the params so that the order is deterministic. Go randomizes each range's ordering
		for k, v := range queryParams {
			queryTokens = append(
				queryTokens,
				fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)),
			)
		}
		urlBuilder.WriteString(strings.Join(queryTokens, "&"))
	}
	return urlBuilder.String()
}

func (c Client) Do(file string, params map[string]string, outResponse any) (err error) {
	req, err := http.NewRequest(http.MethodGet, c.constructUrl(file, params), nil)
	if err != nil {
		return fmt.Errorf("Failed to construct request: %w", err)
	}
	req.SetBasicAuth(c.HttpUser, c.HttpPassword)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to send http request to VLC: %w", err)
	}
	defer func() {
		resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP code indicates error: %s", resp.Status)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read in response body: %w", err)
	}

	// Deserialize the response body into the given struct
	if outResponse != nil {
		err = json.Unmarshal(respBody, outResponse)
		if err != nil {
			return fmt.Errorf("failed to parse response into given type: %w", err)
		}
	}
	// Success!
	return nil
}
