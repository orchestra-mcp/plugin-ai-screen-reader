package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"github.com/orchestra-mcp/plugin-ai-screen-reader/internal/a11y"
	"google.golang.org/protobuf/types/known/structpb"
)

func GetElementHierarchySchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"app_name": map[string]any{
				"type":        "string",
				"description": "Name of the application to inspect. Uses frontmost app if not specified.",
			},
		},
	})
	return s
}

func GetElementHierarchy() func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if !a11y.IsSupported() {
			return helpers.TextResult("Element hierarchy inspection requires macOS (uses System Events)."), nil
		}

		appName := helpers.GetString(req.Arguments, "app_name")

		var script string
		if appName != "" {
			script = fmt.Sprintf(
				`tell application "System Events" to tell (first process whose name is "%s") to get {name, role, subrole, description} of every UI element`,
				appName,
			)
		} else {
			script = `tell application "System Events" to tell (first process whose frontmost is true) to get {name, role, subrole, description} of every UI element`
		}

		result, err := a11y.RunAppleScript(ctx, script)
		if err != nil {
			return helpers.ErrorResult("applescript_error", fmt.Sprintf("failed to get element hierarchy: %v", err)), nil
		}

		return helpers.TextResult(result), nil
	}
}
