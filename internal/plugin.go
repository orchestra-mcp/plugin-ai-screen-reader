package internal

import (
	"github.com/orchestra-mcp/sdk-go/plugin"
	"github.com/orchestra-mcp/plugin-ai-screen-reader/internal/tools"
)

// ToolsPlugin registers all ai.screen-reader tools with the plugin builder.
type ToolsPlugin struct{}

// RegisterTools registers all 6 screen reader tools on the given plugin builder.
func (tp *ToolsPlugin) RegisterTools(builder *plugin.PluginBuilder) {
	builder.RegisterTool("list_windows",
		"List all visible windows and their associated processes",
		tools.ListWindowsSchema(), tools.ListWindows())

	builder.RegisterTool("get_focused_element",
		"Get the currently focused UI element including its role, title, and value",
		tools.GetFocusedElementSchema(), tools.GetFocusedElement())

	builder.RegisterTool("get_accessibility_tree",
		"Get the full accessibility element tree for an application",
		tools.GetAccessibilityTreeSchema(), tools.GetAccessibilityTree())

	builder.RegisterTool("find_element",
		"Find a UI element by its accessibility label and optional role",
		tools.FindElementSchema(), tools.FindElement())

	builder.RegisterTool("get_window_elements",
		"Get all UI elements within a specific window by its title",
		tools.GetWindowElementsSchema(), tools.GetWindowElements())

	builder.RegisterTool("get_element_hierarchy",
		"Get the accessibility element hierarchy (name, role, subrole, description) for an application",
		tools.GetElementHierarchySchema(), tools.GetElementHierarchy())
}
