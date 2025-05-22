package models

import (
	"encoding/json"
	"net/url"
	"time"
)

// Custom type to handle URL parsing
type CustomURL struct {
	raw *url.URL
}

func (u CustomURL) String() string {
	if u.raw == nil {
		return ""
	}
	return u.raw.String()
}

func (u CustomURL) IsNil() bool {
	return u.raw == nil
}

// UnmarshalJSON to parse and validate the URL
func (u *CustomURL) UnmarshalJSON(data []byte) error {
	var rawURL string
	if err := json.Unmarshal(data, &rawURL); err != nil {
		return err
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return err
	}

	u.raw = parsedURL
	return nil
}

func (u CustomURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.raw.String())
}

// Main URL struct to represent the object in storage
type URL struct {
	URL       CustomURL `json:"URL"`
	ExpiresAt time.Time `json:"expires_at"`
}
