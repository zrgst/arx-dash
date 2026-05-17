package arx

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client er ARX HTTP-klienten.
//
// Den kjenner til:
// - base URL
// - Basic Auth
// - self-signed TLS hvis aktivert
type Client struct {
	baseURL  string
	username string
	password string
	http     *http.Client
}

// NewClient lager en ARX client.
func NewClient(baseURL string, username string, password string, allowSelfSigned bool) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	if allowSelfSigned {
		transport.TLSClientConfig = &tls.Config{
			// Matcher Java-eksempelet sin "trust all certificates"-oppførsel.
			// Dette bør kun brukes i test/lab/internt miljø.
			InsecureSkipVerify: true,
		}
	}

	return &Client{
		baseURL:  strings.TrimRight(baseURL, "/"),
		username: username,
		password: password,
		http: &http.Client{
			Timeout:   15 * time.Second,
			Transport: transport,
		},
	}
}

// ExportPersonsRaw henter rå XML fra /arx/export.
func (c *Client) ExportPersonsRaw(ctx context.Context) ([]byte, error) {
	return c.postForm(ctx, "/arx/export", url.Values{})
}

// ExportPersons henter og parser personer/kort/adgangsgrupper.
func (c *Client) ExportPersons(ctx context.Context) (PersonsExport, error) {
	data, err := c.ExportPersonsRaw(ctx)
	if err != nil {
		return PersonsExport{}, err
	}

	return ParsePersonsExport(data)
}

func (c *Client) postForm(ctx context.Context, path string, form url.Values) ([]byte, error) {
	endpoint := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ARX request failed: %s\n%s", resp.Status, string(body))
	}

	return body, nil
}
