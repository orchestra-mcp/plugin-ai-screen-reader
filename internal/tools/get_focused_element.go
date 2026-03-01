package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"github.com/orchestra-mcp/plugin-ai-screen-reader/internal/a11y"
	"google.golang.org/protobuf/types/known/structpb"
)

func GetFocusedElementSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	})
	return s
}

func GetFocusedElement() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if !a11y.IsSupported() {
			return helpers.TextResult("Focused element inspection requires macOS (uses System Events)."), nil
		}

		script := `tell application "System Events" to get {role, title, value} of focused UI element of first application process whose frontmost is true`
		result, err := a11y.RunAppleScript(ctx, script)
		if err != nil {
			return helpers.ErrorResult("applescript_error", fmt.Sprintf("failed to get focused element: %v", err)), nil
		}

		return helpers.TextResult(result), nil
	}
}
