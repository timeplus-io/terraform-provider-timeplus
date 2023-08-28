// SPDX-License-Identifier: MPL-2.0

package timeplus

type Panel struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	// e.g. {"x":0,"y":0,"w":6,"h":2,"nextX":6,"nextY":2}
	Position map[string]any `json:"position"`

	// e.g. `chart`, `markdown`
	VizType string `json:"viz_type"`

	// For chart, the viz_content is the SQL
	// For markdown, the viz_content is the markdown itself
	VizContent string `json:"viz_content"`

	// The JSON configuration of the viz
	// For chart, it is `{ "chart_type": "line", ...  }`
	// For markdown, it is `{ "favour": "github", ...  }`
	VizConfig map[string]any `json:"viz_config"`
}

type Dashboard struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Panels      []Panel `json:"panels"`
}

// redashboardID implements redashboard
func (s Dashboard) resourceID() string {
	return s.ID
}

// redashboardPath implements redashboard
func (Dashboard) resourcePath() string {
	return "dashboards"
}

func (c *Client) CreateDashboard(d *Dashboard) error {
	return c.post(d)
}

func (c *Client) DeleteDashboard(id string) error {
	return c.delete(Dashboard{ID: id})
}

func (c *Client) UpdateDashboard(s *Dashboard) error {
	return c.put(s)
}

func (c *Client) GetDashboard(id string) (Dashboard, error) {
	s := Dashboard{ID: id}
	err := c.get(&s)
	return s, err
}
