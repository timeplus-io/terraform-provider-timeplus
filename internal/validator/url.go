// SPDX-License-Identifier: MPL-2.0

package validator

import (
	"context"
	stdurl "net/url"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// validator for validating URLs
type url struct{}

func URL() validator.String {
	return url{}
}

// Description implements validator.String
func (url) Description(_ context.Context) string {
	return "validates input should be a valid URL"
}

// MarkdownDescription implements validator.String
func (u url) MarkdownDescription(ctx context.Context) string {
	return u.Description(ctx)
}

// ValidateString implements validator.String
func (url) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	_, err := stdurl.Parse(req.ConfigValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "invalid URL", err.Error())
	}
}
