package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"github.com/orchestra-mcp/plugin-ai-screen-reader/internal/a11y"
	"google.golang.org/protobuf/types/known/structpb"
)

func GetWindowElementsSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"window_title": map[string]any{
				"type":        "string",
				"description": "Title of the window to inspect",
			},
		},
		"required": []any{"window_title"},
	})
	return s
}

func GetWindowElements() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if !a11y.IsSupported() {
			return helpers.TextResult("Window element inspection requires macOS (uses System Events)."), nil
		}

		if err := helpers.ValidateRequired(req.Arguments, "window_title"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		windowTitle := helpers.GetString(req.Arguments, "window_title")
		script := fmt.Sprintf(
			`tell application "System Events" to tell first window whose name is "%s" to get entire contents`,
			windowTitle,
		)

		result, err := a11y.RunAppleScript(ctx, script)
		if err != nil {
			return helpers.ErrorResult("applescript_error", fmt.Sprintf("failed to get window elements: %v", err)), nil
		}

		return helpers.TextResult(result), nil
	}
}
