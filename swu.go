package swu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
)

const (
	Endpoint        = "https://api.sendwithus.com/api/v1"
	APIHeaderClient = "golang-0.0.1"
)

type Client struct {
	Client *http.Client
	apiKey string
	URL    string
}

type Template struct {
	ID       string     `json:"id,omitempty"`
	Tags     []string   `json:"tags,omitempty"`
	Created  int64      `json:"created,omitempty"`
	Versions []*Version `json:"versions,omitempty"`
	Name     string     `json:"name,omitempty"`
}

type Version struct {
	Name      string `json:"name,omitempty"`
	ID        string `json:"id,omitempty"`
	Created   int64  `json:"created,omitempty"`
	HTML      string `json:"html,omitempty"`
	Text      string `json:"text,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Published bool   `json:"published,omitempty"`
}

type Email struct {
	ID          string        `json:"email_id,omitempty"`
	Recipient   *Recipient    `json:"recipient,omitempty"`
	CC          []*Recipient  `json:"cc,omitempty"`
	BCC         []*Recipient  `json:"bcc,omitempty"`
	Sender      *Sender       `json:"sender,omitempty"`
	EmailData   interface{}   `json:"email_data,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
	Inline      *Attachment   `json:"inline,omitempty"`
	Files       []*Attachment `json:"files,omitempty"`
	ESPAccount  string        `json:"esp_account,omitempty"`
	VersionName string        `json:"version_name,omitempty"`
}

type SendResponse struct {
	Success   bool   `json:"success"`
	Status    string `json:"status"`
	ReceiptID string `json:"receipt_id"`
	Email     struct {
		Name        string `json:"name"`
		VersionName string `json:"version_name"`
		Locale      string `json:"locale"`
	} `json:"email"`
}

type Drip struct {
	Recipient  *Recipient   `json:"recipient,omitempty"`
	CC         []*Recipient `json:"cc,omitempty"`
	BCC        []*Recipient `json:"bcc,omitempty"`
	EmailData  interface{}  `json:"email_data,omitempty"`
	Sender     *Sender      `json:"sender,omitempty"`
	Tags       []string     `json:"tags,omitempty"`
	ESPAccount string       `json:"esp_account,omitempty"`
}

type Recipient struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Sender struct {
	Recipient
	ReplyTo string `json:"reply_to,omitempty"`
}

type Attachment struct {
	ID   string `json:"id,omitempty"`
	Data string `json:"data,omitempty"`
}

type LogEvent struct {
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

type LogQuery struct {
	Count      int   `json:"count,omitempty" url:"count,omitempty"`
	Offset     int   `json:"offset,omitempty" url:"offset,omitempty"`
	CreatedGT  int64 `json:"created_gt,omitempty" url:"created_gt,omitempty"`
	CreatedGTE int64 `json:"created_gte,omitempty" url:"created_gte,omitempty"`
	CreatedLT  int64 `json:"created_lt,omitempty" url:"created_lt,omitempty"`
	CreatedLTE int64 `json:"created_lte,omitempty" url:"created_lte,omitempty"`
}

type Log struct {
	LogEvent
	ID               string `json:"id,omitempty"`
	RecipientName    string `json:"recipient_name,omitempty"`
	RecipientAddress string `json:"recipient_address,omitempty"`
	Status           string `json:"status,omitempty"`
	EmailID          string `json:"email_id,omitempty"`
	EmailName        string `json:"email_name,omitempty"`
	EmailVersion     string `json:"email_version,omitempty"`
	EventsURL        string `json:"events_url,omitempty"`
}

type LogResend struct {
	Success bool   `json:"success,omitempty"`
	Status  string `json:"status,omitempty"`
	ID      string `json:"log_id,omitempty"`
	Email   struct {
		Name        string `json:"name"`
		VersionName string `json:"version_name"`
	} `json:"email"`
}

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("swu.go: Status code: %d, Error: %s", e.Code, e.Message)
}

func New(apiKey string) *Client {
	return &Client{
		Client: http.DefaultClient,
		apiKey: apiKey,
		URL:    Endpoint,
	}
}

func (c *Client) Templates() ([]*Template, error) {
	return c.Emails()
}

func (c *Client) Emails() ([]*Template, error) {
	var result []*Template
	return result, c.makeRequest("GET", "/templates", nil, &result)
}

func (c *Client) GetTemplate(id string) (*Template, error) {
	var result Template
	return &result, c.makeRequest("GET", "/templates/"+id, nil, &result)
}

func (c *Client) GetTemplateVersion(id, version string) (*Version, error) {
	var result Version
	return &result, c.makeRequest("GET", "/templates/"+id+"/versions/"+version, nil, &result)
}

func (c *Client) UpdateTemplateVersion(id, version string, template *Version) (*Version, error) {
	var result Version
	payload, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}
	return &result, c.makeRequest("PUT", "/templates/"+id+"/versions/"+version, bytes.NewReader(payload), &result)
}

func (c *Client) CreateTemplate(template *Version) (*Template, error) {
	var result Template
	payload, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}
	return &result, c.makeRequest("POST", "/templates", bytes.NewReader(payload), &result)
}

func (c *Client) CreateTemplateVersion(id string, template *Version) (*Template, error) {
	var result Template
	payload, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}
	return &result, c.makeRequest("POST", "/templates/"+id+"/versions", bytes.NewReader(payload), &result)
}

func (c *Client) Send(email *Email) (SendResponse, error) {
	payload, err := json.Marshal(email)
	if err != nil {
		return SendResponse{}, err
	}
	var sr SendResponse
	err = c.makeRequest("POST", "/send", bytes.NewReader(payload), &sr)
	return sr, err
}

func (c *Client) BeginDrip(dripID string, drip *Drip) error {
	payload, err := json.Marshal(drip)
	if err != nil {
		return err
	}
	return c.makeRequest("POST", "/drip_campaigns/"+dripID+"/activate", bytes.NewReader(payload), nil)
}

func (c *Client) CancelDrip(dripID string, recipientEmail string) error {
	payload, err := json.Marshal(map[string]interface{}{"recipient_address": recipientEmail})
	if err != nil {
		return err
	}
	return c.makeRequest("POST", "/drip_campaigns/"+dripID+"/deactivate", bytes.NewReader(payload), nil)
}

func (c *Client) GetLogs(q *LogQuery) ([]*Log, error) {
	var result []*Log
	payload, _ := query.Values(q)
	return result, c.makeRequest("GET", "/logs?"+payload.Encode(), nil, &result)
}

func (c *Client) GetLog(id string) (*Log, error) {
	var result Log
	return &result, c.makeRequest("GET", "/logs/"+id, nil, &result)
}

func (c *Client) GetLogEvents(id string) (*LogEvent, error) {
	var result LogEvent
	return &result, c.makeRequest("GET", "/logs/"+id+"/events", nil, &result)
}

func (c *Client) ResendLog(id string) (*LogResend, error) {
	result := &LogResend{
		ID: id,
	}
	payload, _ := json.Marshal(result)
	return result, c.makeRequest("POST", "/resend", bytes.NewReader(payload), result)
}

func (c *Client) makeRequest(method, endpoint string, body io.Reader, result interface{}) error {
	r, _ := http.NewRequest(method, c.URL+endpoint, body)
	r.SetBasicAuth(c.apiKey, "")
	r.Header.Set("X--API-CLIENT", APIHeaderClient)
	res, err := c.Client.Do(r)
	if err != nil {
		var code int
		if res != nil {
			code = res.StatusCode
		}
		return &Error{
			Code:    code,
			Message: err.Error(),
		}
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &Error{
			Code:    res.StatusCode,
			Message: err.Error(),
		}
	}
	if res.StatusCode >= 300 {
		return &Error{
			Code:    res.StatusCode,
			Message: string(b),
		}
	}
	if result != nil {
		return json.Unmarshal(b, result)
	}
	return nil
}
