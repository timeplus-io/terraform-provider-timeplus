// SPDX-License-Identifier: MPL-2.0

package timeplus

type Source struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Stream      string         `json:"stream"`
	Type        string         `json:"type"`
	Properties  map[string]any `json:"properties"`
}

// resourceID implements resource
func (s Source) resourceID() string {
	return s.ID
}

// resourcePath implements resource
func (Source) resourcePath() string {
	return "sources"
}

func (c *Client) CreateSource(s *Source) error {
	return c.post(s)
}

func (c *Client) DeleteSource(s *Source) error {
	return c.delete(s)
}

func (c *Client) UpdateSource(s *Source) error {
	return c.put(s)
}

func (c *Client) GetSource(id string) (Source, error) {
	s := Source{ID: id}
	err := c.get(&s)
	return s, err
}
