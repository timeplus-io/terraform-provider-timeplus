// SPDX-License-Identifier: MPL-2.0

package timeplus

type Column struct {
	Name                    string `json:"name" binding:"required" example:"name"`
	Type                    string `json:"type" binding:"required" example:"string"`
	Default                 string `json:"default,omitempty"`
	CompressionCodec        string `json:"compression_codec,omitempty"`
	Codec                   string `json:"codec"`
	TTLExpression           string `json:"ttl_expression,omitempty"`
	SkippingIndexExpression string `json:"skipping_index_expression,omitempty"`
}

type StreamMode string

const (
	StreamModeAppend      StreamMode = "append"
	StreamModeChangeLog   StreamMode = "changelog"
	StreamModeChangeLogKV StreamMode = "changelog_kv"
	StreamModeVersionedKV StreamMode = "versioned_kv"
)

type Stream struct {
	// Stream name should only contain a maximum of 64 letters, numbers, or _, and start with a letter
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Columns     []Column `json:"columns"`
	// This column will be used as the event time if specified
	EventTimeColumn string `json:"event_time_column,omitempty"`
	// The timezone of the `EventTimeColumn`
	TimestampTimezone       string `json:"event_time_timezone"`
	Shards                  int    `json:"shards,omitempty"`
	ReplicationFactor       int    `json:"replication_factor,omitempty"`
	OrderByExpression       string `json:"order_by_expression,omitempty"`
	OrderByGranularity      string `json:"order_by_granularity,omitempty"`
	PartitionByGranularity  string `json:"partition_by_granularity,omitempty"`
	HistoricalTTLExpression string `json:"ttl_expression,omitempty"` // API request
	TTL                     string `json:"ttl"`                      // API response

	// Default to `append`.
	Mode string `json:"mode,omitempty"`

	// Expression of primary key, required in `changelog_kv` and `versioned_kv` mode
	PrimaryKey string `json:"primary_key,omitempty"`

	// The max size a stream can grow. Any non-positive value means unlimited size. Default to 10 GiB.
	RetentionBytes int `json:"logstore_retention_bytes,omitempty" example:"10737418240"`

	// The max time the data can be retained in the stream. Any non-positive value means unlimited time. Default to 7 days.
	RetentionMS int `json:"logstore_retention_ms,omitempty" example:"604800000"`
}

// resourceID implements resource
func (s Stream) resourceID() string {
	return s.Name
}

// resourcePath implements resource
func (Stream) resourcePath() string {
	return "streams"
}

func (c *Client) CreateStream(s *Stream) error {
	if err := c.post(s); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteStream(s *Stream) error {
	return c.delete(s)
}

func (c *Client) UpdateStream(s *Stream) error {
	return c.patch(s)
}

func (c *Client) GetStream(name string) (Stream, error) {
	s := Stream{Name: name}
	err := c.get(&s)
	s.HistoricalTTLExpression = s.TTL
	return s, err
}
