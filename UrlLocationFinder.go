package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// Struct complète
type IPInfoResponse struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

func UrlLocationFinder(rawURL string) ([]byte, error) {
	rawURL = strings.TrimSpace(rawURL)

	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	hostname := parsedURL.Hostname()
	if hostname == "" {
		return nil, fmt.Errorf("no hostname found in URL")
	}

	ips, err := net.LookupHost(hostname)
	if err != nil || len(ips) == 0 {
		return nil, fmt.Errorf("error resolving hostname: %w", err)
	}
	ip := ips[0]

	// Appel API JSON complet
	apiURL := fmt.Sprintf("https://ipinfo.io/%s?token=139fdd0c754fae", ip)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// On peut unmarshaler pour vérifier, puis renvoyer JSON propre
	var ipInfo IPInfoResponse
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w\nresponse: %s", err, string(body))
	}

	// Retourner le JSON "pretty" de la struct
	resultJSON, err := json.MarshalIndent(ipInfo, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", err)
	}

	return resultJSON, nil
}
