package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const baseURL = "https://api.rollbar.com/api/1"

type Client struct {
	httpClient *http.Client
	token      string
}

func New() (*Client, error) {
	token := os.Getenv("ROLLBAR_READ_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("ROLLBAR_READ_TOKEN environment variable is not set")
	}
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
	}, nil
}

func (c *Client) doRequest(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Rollbar-Access-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// Item represents a Rollbar item (error group)
type Item struct {
	ID               string `json:"id"`
	Counter          int    `json:"counter"`
	Title            string `json:"title"`
	Level            string `json:"level"`
	Status           string `json:"status"`
	Environment      string `json:"environment"`
	TotalOccurrences int    `json:"total_occurrences"`
	LastOccurrence   string `json:"last_occurrence"`
	FirstOccurrence  string `json:"first_occurrence"`
}

// ItemsResponse represents the API response for items list
type ItemsResponse struct {
	Err    int `json:"err"`
	Result struct {
		Items []struct {
			ID                       json.Number `json:"id"`
			Counter                  int         `json:"counter"`
			Title                    string      `json:"title"`
			Level                    string      `json:"level"`
			Status                   string      `json:"status"`
			Environment              string      `json:"environment"`
			TotalOccurrences         int         `json:"total_occurrences"`
			LastOccurrenceTimestamp  int64       `json:"last_occurrence_timestamp"`
			FirstOccurrenceTimestamp int64       `json:"first_occurrence_timestamp"`
		} `json:"items"`
	} `json:"result"`
}

// ItemResponse represents the API response for a single item
type ItemResponse struct {
	Err    int `json:"err"`
	Result struct {
		ID                       json.Number `json:"id"`
		Counter                  int         `json:"counter"`
		Title                    string      `json:"title"`
		Level                    string      `json:"level"`
		Status                   string      `json:"status"`
		Environment              string      `json:"environment"`
		TotalOccurrences         int         `json:"total_occurrences"`
		LastOccurrenceTimestamp  int64       `json:"last_occurrence_timestamp"`
		FirstOccurrenceTimestamp int64       `json:"first_occurrence_timestamp"`
	} `json:"result"`
}

// Occurrence represents a single occurrence
type Occurrence struct {
	ID        int64                  `json:"id"`
	ItemID    string                 `json:"item_id"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// OccurrencesResponse represents the API response for occurrences list
type OccurrencesResponse struct {
	Err    int `json:"err"`
	Result struct {
		Instances []struct {
			ID        int64                  `json:"id"`
			ItemID    json.Number            `json:"item_id"`
			Timestamp int64                  `json:"timestamp"`
			Data      map[string]interface{} `json:"data"`
		} `json:"instances"`
	} `json:"result"`
}

// OccurrenceResponse represents the API response for a single occurrence
type OccurrenceResponse struct {
	Err    int `json:"err"`
	Result struct {
		ID        int64                  `json:"id"`
		ItemID    json.Number            `json:"item_id"`
		Timestamp int64                  `json:"timestamp"`
		Data      map[string]interface{} `json:"data"`
	} `json:"result"`
}

func formatTimestamp(ts int64) string {
	return time.Unix(ts, 0).UTC().Format("2006-01-02 15:04:05")
}

// GetItems returns a list of items with optional filters
func (c *Client) GetItems(limit int, status, level, env string) ([]Item, error) {
	params := url.Values{}
	if status != "" {
		params.Set("status", status)
	}
	if level != "" {
		params.Set("level", level)
	}
	if env != "" {
		params.Set("environment", env)
	}

	endpoint := "/items"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	body, err := c.doRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var resp ItemsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Err != 0 {
		return nil, fmt.Errorf("API returned error code: %d", resp.Err)
	}

	items := make([]Item, 0, len(resp.Result.Items))
	for i, item := range resp.Result.Items {
		if i >= limit {
			break
		}
		items = append(items, Item{
			ID:               item.ID.String(),
			Counter:          item.Counter,
			Title:            item.Title,
			Level:            item.Level,
			Status:           item.Status,
			Environment:      item.Environment,
			TotalOccurrences: item.TotalOccurrences,
			LastOccurrence:   formatTimestamp(item.LastOccurrenceTimestamp),
			FirstOccurrence:  formatTimestamp(item.FirstOccurrenceTimestamp),
		})
	}

	return items, nil
}

// GetItem returns details of a specific item
func (c *Client) GetItem(itemID string) (*Item, error) {
	body, err := c.doRequest("/item/" + itemID)
	if err != nil {
		return nil, err
	}

	var resp ItemResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Err != 0 {
		return nil, fmt.Errorf("API returned error code: %d", resp.Err)
	}

	return &Item{
		ID:               resp.Result.ID.String(),
		Counter:          resp.Result.Counter,
		Title:            resp.Result.Title,
		Level:            resp.Result.Level,
		Status:           resp.Result.Status,
		Environment:      resp.Result.Environment,
		TotalOccurrences: resp.Result.TotalOccurrences,
		LastOccurrence:   formatTimestamp(resp.Result.LastOccurrenceTimestamp),
		FirstOccurrence:  formatTimestamp(resp.Result.FirstOccurrenceTimestamp),
	}, nil
}

// GetOccurrences returns occurrences for a specific item
func (c *Client) GetOccurrences(itemID string, limit int) ([]Occurrence, error) {
	endpoint := "/item/" + itemID + "/instances"

	body, err := c.doRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var resp OccurrencesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Err != 0 {
		return nil, fmt.Errorf("API returned error code: %d", resp.Err)
	}

	occurrences := make([]Occurrence, 0, len(resp.Result.Instances))
	for i, occ := range resp.Result.Instances {
		if i >= limit {
			break
		}
		occurrences = append(occurrences, Occurrence{
			ID:        occ.ID,
			ItemID:    occ.ItemID.String(),
			Timestamp: formatTimestamp(occ.Timestamp),
			Data:      occ.Data,
		})
	}

	return occurrences, nil
}

// GetOccurrence returns a single occurrence
func (c *Client) GetOccurrence(occurrenceID string) (*Occurrence, error) {
	body, err := c.doRequest("/instance/" + occurrenceID)
	if err != nil {
		return nil, err
	}

	var resp OccurrenceResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Err != 0 {
		return nil, fmt.Errorf("API returned error code: %d", resp.Err)
	}

	return &Occurrence{
		ID:        resp.Result.ID,
		ItemID:    resp.Result.ItemID.String(),
		Timestamp: formatTimestamp(resp.Result.Timestamp),
		Data:      resp.Result.Data,
	}, nil
}

// GetOccurrenceRaw returns raw occurrence data
func (c *Client) GetOccurrenceRaw(occurrenceID string) ([]byte, error) {
	return c.doRequest("/instance/" + occurrenceID)
}

// FormatOccurrenceData extracts essential fields from occurrence data
func FormatOccurrenceData(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	if env, ok := data["environment"]; ok {
		result["environment"] = env
	}
	if level, ok := data["level"]; ok {
		result["level"] = level
	}

	// Extract message
	if body, ok := data["body"].(map[string]interface{}); ok {
		if msg, ok := body["message"].(map[string]interface{}); ok {
			if msgBody, ok := msg["body"]; ok {
				result["message"] = msgBody
			}
		}
		// Extract exception info
		if trace, ok := body["trace"].(map[string]interface{}); ok {
			if exc, ok := trace["exception"].(map[string]interface{}); ok {
				result["exception_class"] = exc["class"]
				result["exception_message"] = exc["message"]
			}
			if frames, ok := trace["frames"].([]interface{}); ok {
				backtrace := make([]string, 0, len(frames))
				for _, f := range frames {
					if frame, ok := f.(map[string]interface{}); ok {
						filename := frame["filename"]
						lineno := frame["lineno"]
						method := frame["method"]
						backtrace = append(backtrace, fmt.Sprintf("%v:%v in %v", filename, lineno, method))
					}
				}
				result["backtrace"] = backtrace
			}
		}
		// Extract trace_chain for chained exceptions
		if traceChain, ok := body["trace_chain"].([]interface{}); ok && len(traceChain) > 0 {
			if trace, ok := traceChain[0].(map[string]interface{}); ok {
				if exc, ok := trace["exception"].(map[string]interface{}); ok {
					result["exception_class"] = exc["class"]
					result["exception_message"] = exc["message"]
				}
				if frames, ok := trace["frames"].([]interface{}); ok {
					backtrace := make([]string, 0, len(frames))
					for _, f := range frames {
						if frame, ok := f.(map[string]interface{}); ok {
							filename := frame["filename"]
							lineno := frame["lineno"]
							method := frame["method"]
							backtrace = append(backtrace, fmt.Sprintf("%v:%v in %v", filename, lineno, method))
						}
					}
					result["backtrace"] = backtrace
				}
			}
		}
	}

	// Extract server info
	if server, ok := data["server"].(map[string]interface{}); ok {
		serverInfo := make(map[string]interface{})
		if host, ok := server["host"]; ok {
			serverInfo["host"] = host
		}
		if root, ok := server["root"]; ok {
			serverInfo["root"] = root
		}
		if pid, ok := server["pid"]; ok {
			serverInfo["pid"] = pid
		}
		if len(serverInfo) > 0 {
			result["server"] = serverInfo
		}
	}

	// Extract request info
	if req, ok := data["request"].(map[string]interface{}); ok {
		reqInfo := make(map[string]interface{})
		if url, ok := req["url"]; ok {
			reqInfo["url"] = url
		}
		if method, ok := req["method"]; ok {
			reqInfo["method"] = method
		}
		if userIP, ok := req["user_ip"]; ok {
			reqInfo["user_ip"] = userIP
		}
		if params, ok := req["params"]; ok {
			reqInfo["params"] = params
		}
		if headers, ok := req["headers"]; ok {
			reqInfo["headers"] = headers
		}
		if len(reqInfo) > 0 {
			result["request"] = reqInfo
		}
	}

	return result
}

// IntToString converts int64 to string
func IntToString(n int64) string {
	return strconv.FormatInt(n, 10)
}
