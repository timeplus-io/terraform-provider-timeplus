// SPDX-License-Identifier: MPL-2.0

package timeplus

type UDFType string

const (
	UDFTypeJavascript UDFType = "javascript"
	UDFTypeRemote     UDFType = "remote"
)

type UDFAuthMethod string

const (
	UDFAuthHeader UDFAuthMethod = "auth_header"
	UDFAuthNone   UDFAuthMethod = "none"
)

type UDFArgument struct {
	Name string `json:"name" binding:"required" example:"val"`
	Type string `json:"type" binding:"required" example:"float64"`
}

type UDFAuthContext struct {
	Name  string `json:"key_name"`
	Value string `json:"key_value"`
}

type UDF struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        UDFType `json:"type"`

	// The input argument of the UDF, the type should be Timeplus data types,  not javascript types.
	// Only int8/16/32/64, uint8/16/32/64 are supported.
	Arguments []UDFArgument `json:"arguments"`

	// The erturn type of the UDF
	//   * For UDA: if it returns a single value, the return type is the corresponding data type of Timeplus.
	//     It supports the same types of input arguments, except for datetime, it only supports DateTime64(3).
	ReturnType string `json:"return_type" example:"float64"`

	// Only valid when `type` is `remote`.
	URL string `json:"url,omitempty"`

	// Only valid when `type` is `remote`.
	// This field is used to set the authentication method for remote UDF. It can be either `auth_header` or `none`.
	// When `auth_header` is set, you can configure `auth_context` to specify the HTTP header that be sent the remote URL
	AuthMethod UDFAuthMethod `json:"auth_method,omitempty"`

	// Only valid when `type` is `remote` and `auth_method` is `auth_header`
	AuthContext *UDFAuthContext `json:"auth_context,omitempty"`

	// Only valid when type is 'javascript'. Whether it is an aggregation function.
	IsAggrFunction bool `json:"is_aggregation,omitempty"`

	Source string `json:"source,omitempty"`
}

// resourceID implements resource
func (u UDF) resourceID() string {
	return u.Name
}

// resourcePath implements resource
func (UDF) resourcePath() string {
	return "udfs"
}

func (c *Client) CreateUDF(u *UDF) error {
	return c.post(u)
}

func (c *Client) DeleteUDF(name string) error {
	return c.delete(&UDF{Name: name})
}

func (c *Client) UpdateUDF(u *UDF) error {
	return c.put(u)
}

func (c *Client) GetUDF(name string) (UDF, error) {
	u := UDF{Name: name}
	err := c.get(&u)
	return u, err
}
