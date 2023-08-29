// SPDX-License-Identifier: MPL-2.0

package timeplus

type View struct {
	Name        string
	Description string
	Query       string
}

func (v *View) toAPIModel() viewAPIModel {
	if v == nil {
		return viewAPIModel{}
	}

	return viewAPIModel{
		Name:         v.Name,
		Description:  v.Description,
		Query:        v.Query,
		Materialized: false,
	}
}

func (v *View) fromAPIModel(m viewAPIModel) {
	if v == nil {
		*v = View{}
	}
	v.Name = m.Name
	v.Description = m.Description
	v.Query = m.Query
}

func (c *Client) CreateView(v *View) error {
	m := v.toAPIModel()
	return c.createView(&m)
}

func (c *Client) DeleteView(v *View) error {
	m := v.toAPIModel()
	return c.deleteView(&m)
}

func (c *Client) UpdateView(v *View) error {
	m := v.toAPIModel()
	return c.updateView(&m)
}

func (c *Client) GetView(name string) (v View, err error) {
	m := viewAPIModel{Name: name}
	if err = c.get(&m); err != nil {
		return
	}
	v.fromAPIModel(m)
	return
}

type MaterializedView struct {
	View

	TargetStream   string
	RetentionBytes int
	RetentionMS    int
	TTLExpression  string
}

func (v *MaterializedView) toAPIModel() viewAPIModel {
	if v == nil {
		return viewAPIModel{}
	}

	return viewAPIModel{
		Name:           v.Name,
		Description:    v.Description,
		Query:          v.Query,
		Materialized:   true,
		TargetStream:   v.TargetStream,
		TTLExpression:  v.TTLExpression,
		RetentionBytes: v.RetentionBytes,
		RetentionMS:    v.RetentionMS,
	}
}

func (v *MaterializedView) fromAPIModel(m viewAPIModel) {
	if v == nil {
		*v = MaterializedView{}
	}
	v.Name = m.Name
	v.Description = m.Description
	v.Query = m.Query
	v.TargetStream = m.TargetStream
	v.TTLExpression = m.TTL
	v.RetentionBytes = m.RetentionBytes
	v.RetentionMS = m.RetentionMS
}

func (c *Client) CreateMaterializedView(v *MaterializedView) error {
	m := v.toAPIModel()
	if err := c.post(&m); err != nil {
		return err
	}
	v.fromAPIModel(m)
	return nil
}

func (c *Client) DeleteMaterializedView(v *MaterializedView) error {
	m := v.toAPIModel()
	return c.delete(&m)
}

func (c *Client) UpdateMaterializedView(v *MaterializedView) error {
	m := v.toAPIModel()
	if err := c.patch(&m); err != nil {
		return err
	}
	v.fromAPIModel(m)
	return nil
}

func (c *Client) GetMaterializedView(name string) (v MaterializedView, err error) {
	m := viewAPIModel{Name: name}
	if err = c.get(&m); err != nil {
		return
	}
	v.fromAPIModel(m)
	return
}

type viewAPIModel struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Query          string `json:"query,omitempty"` // can't set query when updating
	Materialized   bool   `json:"materialized"`
	TargetStream   string `json:"target_stream,omitempty"`
	RetentionBytes int    `json:"logstore_retention_bytes,omitempty"`
	RetentionMS    int    `json:"logstore_retention_ms,omitempty"`
	TTLExpression  string `json:"ttl_expression,omitempty"`
	TTL            string `json:"ttl,omitempty"` // ttl is the field name from API responses
}

// resourceID implements resource
func (v viewAPIModel) resourceID() string {
	return v.Name
}

// resourcePath implements resource
func (viewAPIModel) resourcePath() string {
	return "views"
}

func (c *Client) createView(v *viewAPIModel) error {
	return c.post(v)
}

func (c *Client) deleteView(v *viewAPIModel) error {
	return c.delete(v)
}

func (c *Client) updateView(v *viewAPIModel) error {
	return c.patch(v)
}

func (c *Client) getView(name string) (viewAPIModel, error) {
	v := viewAPIModel{Name: name}
	err := c.get(&v)
	return v, err
}
