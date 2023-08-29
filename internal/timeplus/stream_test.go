// SPDX-License-Identifier: MPL-2.0

package timeplus_test

import (
	"os"
	"testing"

	"github.com/timeplus-io/terraform-provider-timeplus/internal/timeplus"
)

func newClient(t *testing.T) *timeplus.Client {
	apiKey := os.Getenv("TIMEPLUS_API_KEY")
	if len(apiKey) == 0 {
		t.Log("no API key used")
	} else {
		t.Logf("found API key: %s[=== scrubbed ===]", apiKey[:8])
	}
	c, err := timeplus.NewClient("latest", apiKey, timeplus.ClientOptions{
		BaseURL: "https://dev.timeplus.cloud",
	})
	if err != nil {
		t.Fatalf("unable to create `timeplus.Client`: %v", err)
	}
	return c
}

func TestCreateStream(t *testing.T) {
	c := newClient(t)
	s := timeplus.Stream{
		Name:        "test_stream",
		Description: "before update",
		Columns: []timeplus.Column{
			{
				Name:    "col_1",
				Type:    "string",
				Default: "empty",
			},
			{
				Name: "col_2",
				Type: "integer",
			},
		},
	}

	defer t.Run("DeleteStream", func(t *testing.T) {
		if err := c.DeleteStream(&s); err != nil {
			t.Fatalf("DeleteStream failed: %v", err)
		}
	})

	if err := c.CreateStream(&s); err != nil {
		t.Fatalf("CreateStream failed: %v", err)
	}

	if s.Description != "before update" {
		t.Fatalf("expected description %q, but got %q", "before update", s.Description)
	}

	t.Run("UpdateStream", func(t *testing.T) {
		s.Description = "after update"
		if err := c.UpdateStream(&s); err != nil {
			t.Fatalf("UpdateStream failed: %v", err)
		}

		if s.Description != "after update" {
			t.Fatalf("expected description %q, but got %q", "after update", s.Description)
		}
	})
}
