// SPDX-License-Identifier: MPL-2.0

package validator

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// validator for validating a string is a valid JSON object
type jsonObject struct{}

func JsonObject() validator.String {
	return jsonObject{}
}

// Description implements validator.String
func (jsonObject) Description(_ context.Context) string {
	return "validates input should be a valid JSON object"
}

// MarkdownDescription implements validator.String
func (j jsonObject) MarkdownDescription(ctx context.Context) string {
	return j.Description(ctx)
}

// ValidateString implements validator.String
func (jsonObject) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// we can introduce an option `ignoreEmpty` when it's needed
	if req.ConfigValue.ValueString() == "" {
		return
	}
	var v map[string]any
	if err := json.Unmarshal([]byte(req.ConfigValue.ValueString()), &v); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "invalid JSON object", err.Error())
	}
}

// validator for validating a string is a valid JSON array of objects
type jsonArrayOfObject struct{}

func JsonArrayOfObject() validator.String {
	return jsonArrayOfObject{}
}

// Description implements validator.String
func (jsonArrayOfObject) Description(_ context.Context) string {
	return "validates input should be a valid JSON array of objects"
}

// MarkdownDescription implements validator.String
func (j jsonArrayOfObject) MarkdownDescription(ctx context.Context) string {
	return j.Description(ctx)
}

// ValidateString implements validator.String
func (jsonArrayOfObject) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// we can introduce an option `ignoreEmpty` when it's needed
	if req.ConfigValue.ValueString() == "" {
		return
	}
	var v []map[string]any
	if err := json.Unmarshal([]byte(req.ConfigValue.ValueString()), &v); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "invalid JSON array", err.Error())
	}
}
