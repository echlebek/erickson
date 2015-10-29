package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var Redirect = errors.New("redirect")

type Client struct {
	*http.Client
	URL          *url.URL
	redirectedTo string // stash location from redirect here
}

// New creates a new erickson client.
func New(URL string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	url, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Client: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// TODO FIXME !!!
					// Use CURL_CA_BUNDLE or user supplied certs.
					InsecureSkipVerify: true,
				},
			},
		},
		URL: url,
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		client.redirectedTo = req.URL.String()
		return nil
	}
	return client, nil
}

// Token gets a one-time X-CSRF-Token
func (c *Client) Token() (string, error) {
	req, err := http.NewRequest("HEAD", c.URL.String(), nil)
	if err != nil {
		return "", err
	}
	response, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	token := response.Header.Get("X-CSRF-Token")
	return token, nil
}

func (c *Client) Authenticate(username, password string) (*http.Response, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	form := url.Values{"username": {username}, "password": {password}}
	req, err := http.NewRequest("POST", c.URL.String()+"/login", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-CSRF-Token", token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Post username and password
	response, err := c.Do(req)
	if err != nil {
		return response, err
	}
	return response, response.Body.Close()
}

func (c *Client) Session(sessionKey string) error {
	cookie := http.Cookie{
		Name:     "erickson",
		Value:    sessionKey,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		return err
	}
	req.AddCookie(&cookie)
	if _, err := c.Do(req); err != nil {
		return err
	}
	return nil
}

// PostReview posts a review and returns its location.
func (c *Client) PostReview(diff, username, commitmsg, repo string) (string, error) {
	token, err := c.Token()
	form := url.Values{
		"diff":       {diff},
		"submitter":  {username},
		"commitmsg":  {commitmsg},
		"repository": {repo},
	}
	req, err := http.NewRequest("POST", c.URL.String()+"/reviews", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-CSRF-Token", token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return "", err
	}
	response, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		return "", errors.New(fmt.Sprintf("HTTP status %d", response.StatusCode))
	}
	return c.redirectedTo, nil
}
