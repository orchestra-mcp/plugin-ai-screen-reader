package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"github.com/orchestra-mcp/plugin-ai-screen-reader/internal/a11y"
	"google.golang.org/protobuf/types/known/structpb"
)

func FindElementSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"label": map[string]any{
				"type":        "string",
				"description": "Accessibility label or title of the element to find",
			},
			"role": map[string]any{
				"type":        "string",
				"description": "Accessibility role to filter by (e.g. button, textField). Optional.",
			},
		},
		"required": []any{"label"},
	})
	return s
}

func FindElement() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if !a11y.IsSupported() {
			return helpers.TextResult("Element search requires macOS (uses System Events)."), nil
		}

		if err := helpers.ValidateRequired(req.Arguments, "label"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		label := helpers.GetString(req.Arguments, "label")
		role := helpers.GetString(req.Arguments, "role")

		var script string
		if role != "" {
			script = fmt.Sprintf(
				`tell application "System Events" to tell (first process whose frontmost is true) to get UI elements whose title is "%s" and role is "%s"`,
				label, role,
			)
		} else {
			script = fmt.Sprintf(
				`tell application "System Events" to tell (first process whose frontmost is true) to get UI elements whose title is "%s"`,
				label,
			)
		}

		result, err := a11y.RunAppleScript(ctx, script)
		if err != nil {
			return helpers.ErrorResult("applescript_error", fmt.Sprintf("failed to find element: %v", err)), nil
		}

		return helpers.TextResult(result), nil
	}
}
