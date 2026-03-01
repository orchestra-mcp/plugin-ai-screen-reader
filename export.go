package aiscreenreader

import (
	"github.com/orchestra-mcp/plugin-ai-screen-reader/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// Register adds all screen reader tools to the builder.
func Register(builder *plugin.PluginBuilder) {
	tp := &internal.ToolsPlugin{}
	tp.RegisterTools(builder)
}
