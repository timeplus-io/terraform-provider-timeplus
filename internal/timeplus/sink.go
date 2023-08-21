// SPDX-License-Identifier: MPL-2.0

package timeplus

type Sink struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	SQL         string `json:"sql"`   // request
	Query       string `json:"query"` // response

	// Additional configurations such as broker url and etc. should be passed through `properties`
	Type string `json:"type"`

	// Additional properties that required to write the data to the sink (e.g. broker url). Please refer to the sinks documentation
	Properties map[string]any `json:"properties"`
}

// resourceID implements resource
func (s Sink) resourceID() string {
	return s.ID
}

// resourcePath implements resource
func (Sink) resourcePath() string {
	return "sinks"
}

func (c *Client) CreateSink(s *Sink) error {
	return c.post(s)
}

func (c *Client) DeleteSink(s *Sink) error {
	return c.delete(s)
}

func (c *Client) UpdateSink(s *Sink) error {
	return c.put(s)
}

func (c *Client) GetSink(id string) (Sink, error) {
	s := Sink{ID: id}
	err := c.get(&s)
	s.SQL = s.Query // caused by inconsistent API design
	return s, err
}
