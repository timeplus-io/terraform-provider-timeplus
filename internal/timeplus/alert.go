// SPDX-License-Identifier: MPL-2.0

package timeplus

type Alert struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Severity    int            `json:"severity"`
	Action      string         `json:"action"`
	Properties  map[string]any `json:"properties"`
	TriggerSQL  string         `json:"trigger_sql"`
	ResolveSQL  string         `json:"resolve_sql"`
}

// resourceID implements resource
func (a Alert) resourceID() string {
	return a.ID
}

// resourcePath implements resource
func (Alert) resourcePath() string {
	return "alerts"
}

func (c *Client) CreateAlert(s *Alert) error {
	return c.post(s)
}

func (c *Client) DeleteAlert(s *Alert) error {
	return c.delete(s)
}

func (c *Client) UpdateAlert(s *Alert) error {
	return c.put(s)
}

func (c *Client) GetAlert(id string) (Alert, error) {
	s := Alert{ID: id}
	err := c.get(&s)
	return s, err
}
