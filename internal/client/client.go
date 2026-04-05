package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DefaultAPIURL = "https://api.authfi.app"

// Client is the AuthFI API client used by the Terraform provider.
type Client struct {
	APIURL     string
	APIKey     string
	Tenant     string
	HTTPClient *http.Client
}

// New creates a new AuthFI API client.
func New(apiKey, tenant string) *Client {
	return &Client{
		APIURL: DefaultAPIURL,
		APIKey: apiKey,
		Tenant: tenant,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// --- Platform API (superadmin) ---

func (c *Client) platformURL(path string) string {
	return fmt.Sprintf("%s/platform%s", c.APIURL, path)
}

// --- Management API (tenant-scoped) ---

func (c *Client) mgmtURL(path string) string {
	return fmt.Sprintf("%s/manage/v1/%s%s", c.APIURL, c.Tenant, path)
}

// --- HTTP helpers ---

func (c *Client) do(method, url string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// --- Project CRUD ---

type Project struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Slug    string            `json:"slug"`
	Region  string            `json:"region,omitempty"`
	EnvType string            `json:"env_type,omitempty"`
	Plan    string            `json:"plan,omitempty"`
	Tags    map[string]string `json:"tags,omitempty"`
}

func (c *Client) CreateProject(p *Project) (*Project, error) {
	body, status, err := c.do("POST", c.platformURL("/projects"), p)
	if err != nil {
		return nil, err
	}
	if status != 201 && status != 200 {
		return nil, fmt.Errorf("create project: %s (status %d)", string(body), status)
	}
	var result Project
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode project: %w", err)
	}
	return &result, nil
}

func (c *Client) GetProject(id string) (*Project, error) {
	body, status, err := c.do("GET", c.platformURL("/projects/"+id), nil)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, nil
	}
	if status != 200 {
		return nil, fmt.Errorf("get project: %s (status %d)", string(body), status)
	}
	var result Project
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode project: %w", err)
	}
	return &result, nil
}

func (c *Client) DeleteProject(id string) error {
	_, status, err := c.do("DELETE", c.platformURL("/projects/"+id), nil)
	if err != nil {
		return err
	}
	if status != 200 && status != 204 && status != 404 {
		return fmt.Errorf("delete project: status %d", status)
	}
	return nil
}

// --- Organization CRUD ---

type Organization struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Slug     string            `json:"slug"`
	EnvType  string            `json:"env_type,omitempty"`
	Tags     map[string]string `json:"tags,omitempty"`
	LogoURL  *string           `json:"logo_url,omitempty"`
	Primary  *string           `json:"primary_color,omitempty"`
}

func (c *Client) CreateOrganization(o *Organization) (*Organization, error) {
	body, status, err := c.do("POST", c.mgmtURL("/organizations"), o)
	if err != nil {
		return nil, err
	}
	if status != 201 && status != 200 {
		return nil, fmt.Errorf("create org: %s (status %d)", string(body), status)
	}
	var result Organization
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode org: %w", err)
	}
	return &result, nil
}

func (c *Client) GetOrganization(id string) (*Organization, error) {
	body, status, err := c.do("GET", c.mgmtURL("/organizations/"+id), nil)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, nil
	}
	if status != 200 {
		return nil, fmt.Errorf("get org: %s (status %d)", string(body), status)
	}
	var result Organization
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode org: %w", err)
	}
	return &result, nil
}

func (c *Client) UpdateOrganization(id string, o *Organization) (*Organization, error) {
	body, status, err := c.do("PATCH", c.mgmtURL("/organizations/"+id), o)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("update org: %s (status %d)", string(body), status)
	}
	var result Organization
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode org: %w", err)
	}
	return &result, nil
}

func (c *Client) DeleteOrganization(id string) error {
	_, status, err := c.do("DELETE", c.mgmtURL("/organizations/"+id), nil)
	if err != nil {
		return err
	}
	if status != 200 && status != 204 && status != 404 {
		return fmt.Errorf("delete org: status %d", status)
	}
	return nil
}

// --- Application CRUD ---

type Application struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret,omitempty"`
	AppType      string   `json:"app_type"`
	CallbackURLs []string `json:"allowed_callback_urls,omitempty"`
	LogoutURLs   []string `json:"allowed_logout_urls,omitempty"`
	Origins      []string `json:"allowed_origins,omitempty"`
}

func (c *Client) CreateApplication(a *Application) (*Application, error) {
	body, status, err := c.do("POST", c.mgmtURL("/applications"), a)
	if err != nil {
		return nil, err
	}
	if status != 201 && status != 200 {
		return nil, fmt.Errorf("create app: %s (status %d)", string(body), status)
	}
	var result Application
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode app: %w", err)
	}
	return &result, nil
}

func (c *Client) GetApplication(id string) (*Application, error) {
	body, status, err := c.do("GET", c.mgmtURL("/applications/"+id), nil)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, nil
	}
	if status != 200 {
		return nil, fmt.Errorf("get app: %s (status %d)", string(body), status)
	}
	var result Application
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode app: %w", err)
	}
	return &result, nil
}

func (c *Client) UpdateApplication(id string, a *Application) (*Application, error) {
	body, status, err := c.do("PATCH", c.mgmtURL("/applications/"+id), a)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("update app: %s (status %d)", string(body), status)
	}
	var result Application
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode app: %w", err)
	}
	return &result, nil
}

func (c *Client) DeleteApplication(id string) error {
	_, status, err := c.do("DELETE", c.mgmtURL("/applications/"+id), nil)
	if err != nil {
		return err
	}
	if status != 200 && status != 204 && status != 404 {
		return fmt.Errorf("delete app: status %d", status)
	}
	return nil
}

// --- Connection CRUD ---

type Connection struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Type     string `json:"type"`
	ClientID string `json:"client_id,omitempty"`
}

func (c *Client) CreateConnection(conn *Connection) (*Connection, error) {
	body, status, err := c.do("POST", c.mgmtURL("/connections"), conn)
	if err != nil {
		return nil, err
	}
	if status != 201 && status != 200 {
		return nil, fmt.Errorf("create connection: %s (status %d)", string(body), status)
	}
	var result Connection
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode connection: %w", err)
	}
	return &result, nil
}

func (c *Client) GetConnection(id string) (*Connection, error) {
	body, status, err := c.do("GET", c.mgmtURL("/connections/"+id), nil)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, nil
	}
	if status != 200 {
		return nil, fmt.Errorf("get connection: %s (status %d)", string(body), status)
	}
	var result Connection
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode connection: %w", err)
	}
	return &result, nil
}

func (c *Client) DeleteConnection(id string) error {
	_, status, err := c.do("DELETE", c.mgmtURL("/connections/"+id), nil)
	if err != nil {
		return err
	}
	if status != 200 && status != 204 && status != 404 {
		return fmt.Errorf("delete connection: status %d", status)
	}
	return nil
}

// --- Role CRUD ---

type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

func (c *Client) CreateRole(r *Role) (*Role, error) {
	body, status, err := c.do("POST", c.mgmtURL("/roles"), r)
	if err != nil {
		return nil, err
	}
	if status != 201 && status != 200 {
		return nil, fmt.Errorf("create role: %s (status %d)", string(body), status)
	}
	var result Role
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode role: %w", err)
	}
	return &result, nil
}

func (c *Client) GetRole(id string) (*Role, error) {
	body, status, err := c.do("GET", c.mgmtURL("/roles/"+id), nil)
	if err != nil {
		return nil, err
	}
	if status == 404 {
		return nil, nil
	}
	if status != 200 {
		return nil, fmt.Errorf("get role: %s (status %d)", string(body), status)
	}
	var result Role
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode role: %w", err)
	}
	return &result, nil
}

func (c *Client) DeleteRole(id string) error {
	_, status, err := c.do("DELETE", c.mgmtURL("/roles/"+id), nil)
	if err != nil {
		return err
	}
	if status != 200 && status != 204 && status != 404 {
		return fmt.Errorf("delete role: status %d", status)
	}
	return nil
}
