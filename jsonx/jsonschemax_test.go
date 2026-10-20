package jsonx

import (
	"encoding/json"
	"strings"
	"testing"
)

type closedSchemaFixture struct {
	Name string `json:"name"`
}

func TestAllowAdditionalProperties(t *testing.T) {
	t.Parallel()

	schema := SchemaFrom[closedSchemaFixture]()
	if !strings.Contains(string(schema), `"additionalProperties":false`) {
		t.Fatalf("reflected schema should close additional properties, got %s", schema)
	}

	opened := AllowAdditionalProperties(schema)
	if !strings.Contains(string(opened), `"additionalProperties":true`) {
		t.Fatalf("jsonx.AllowAdditionalProperties() should open additional properties, got %s", opened)
	}
}

func TestAllowAdditionalPropertiesKeepsUnreadableSchema(t *testing.T) {
	t.Parallel()

	for _, schema := range []json.RawMessage{json.RawMessage(`not json`), json.RawMessage(`{"$ref":"#/$defs/Missing"}`)} {
		if got := AllowAdditionalProperties(schema); string(got) != string(schema) {
			t.Fatalf("jsonx.AllowAdditionalProperties(%s) = %s, want it unchanged", schema, got)
		}
	}
}

func TestSchemaID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		schema json.RawMessage
		want   string
	}{
		{
			name:   "extracts definition key from ref path",
			schema: json.RawMessage(`{"$ref":"#/$defs/MyType"}`),
			want:   "MyType",
		},
		{
			name:   "nested ref path returns base",
			schema: json.RawMessage(`{"$ref":"#/$defs/deep/Nested"}`),
			want:   "Nested",
		},
		{
			name:   "invalid JSON returns empty string",
			schema: json.RawMessage(`not-json`),
			want:   "",
		},
		{
			name:   "missing ref returns dot",
			schema: json.RawMessage(`{}`),
			want:   ".",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := SchemaID(tc.schema)
			if got != tc.want {
				t.Fatalf("jsonx.SchemaID() = %q, want %q", got, tc.want)
			}
		})
	}
}
