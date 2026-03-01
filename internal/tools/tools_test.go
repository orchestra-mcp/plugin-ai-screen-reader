package tools

import (
	"context"
	"runtime"
	"strings"
	"testing"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

// callTool invokes a tool handler with the provided args map and returns the
// response. It fails the test immediately if NewStruct or the handler itself
// returns a Go-level error.
func callTool(t *testing.T, handler func(context.Context, *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error), args map[string]any) *pluginv1.ToolResponse {
	t.Helper()
	var s *structpb.Struct
	if args != nil {
		var err error
		s, err = structpb.NewStruct(args)
		if err != nil {
			t.Fatalf("NewStruct: %v", err)
		}
	}
	resp, err := handler(context.Background(), &pluginv1.ToolRequest{Arguments: s})
	if err != nil {
		t.Fatalf("handler returned Go error: %v", err)
	}
	return resp
}

func isError(resp *pluginv1.ToolResponse) bool { return resp != nil && !resp.Success }
func errorCode(resp *pluginv1.ToolResponse) string {
	if resp == nil {
		return ""
	}
	return resp.GetErrorCode()
}
func isMacOS() bool { return runtime.GOOS == "darwin" }

// requiresMacOSMessage returns true if the response text contains the canonical
// "requires macOS" phrase returned by all tools on non-darwin platforms.
func requiresMacOSMessage(resp *pluginv1.ToolResponse) bool {
	if resp == nil || resp.Result == nil {
		return false
	}
	fields := resp.Result.GetFields()
	if fields == nil {
		return false
	}
	textVal, ok := fields["text"]
	if !ok {
		return false
	}
	return strings.Contains(textVal.GetStringValue(), "macOS")
}

// --- get_accessibility_tree ---

func TestGetAccessibilityTree_NoApp(t *testing.T) {
	resp := callTool(t, GetAccessibilityTree(), nil)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if isMacOS() {
		// On macOS the tool runs AppleScript: accept either success or
		// applescript_error (osascript may be denied in CI).
		if isError(resp) && errorCode(resp) != "applescript_error" {
			t.Errorf("unexpected error code %q, want applescript_error or success", errorCode(resp))
		}
	} else {
		// On non-macOS the tool must succeed and explain macOS is required.
		if isError(resp) {
			t.Errorf("expected success on non-macOS, got error: %s", errorCode(resp))
		}
		if !requiresMacOSMessage(resp) {
			t.Errorf("expected response to mention macOS requirement, got: %v", resp)
		}
	}
}

func TestGetAccessibilityTree_WithAppName(t *testing.T) {
	resp := callTool(t, GetAccessibilityTree(), map[string]any{"app_name": "Finder"})
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if isMacOS() {
		if isError(resp) && errorCode(resp) != "applescript_error" {
			t.Errorf("unexpected error code %q, want applescript_error or success", errorCode(resp))
		}
	} else {
		if isError(resp) {
			t.Errorf("expected success on non-macOS, got error: %s", errorCode(resp))
		}
		if !requiresMacOSMessage(resp) {
			t.Errorf("expected response to mention macOS requirement, got: %v", resp)
		}
	}
}

// --- get_focused_element ---

func TestGetFocusedElement_NoArgs(t *testing.T) {
	resp := callTool(t, GetFocusedElement(), nil)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if isMacOS() {
		if isError(resp) && errorCode(resp) != "applescript_error" {
			t.Errorf("unexpected error code %q, want applescript_error or success", errorCode(resp))
		}
	} else {
		if isError(resp) {
			t.Errorf("expected success on non-macOS, got error: %s", errorCode(resp))
		}
		if !requiresMacOSMessage(resp) {
			t.Errorf("expected response to mention macOS requirement, got: %v", resp)
		}
	}
}

// --- find_element ---

// TestFindElement_MissingLabel verifies the validation_error path for the
// missing required "label" field. This path is only reachable on macOS because
// the tool returns early with a success message on other platforms.
func TestFindElement_MissingLabel(t *testing.T) {
	resp := callTool(t, FindElement(), nil)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if isMacOS() {
		// IsSupported() passes, so ValidateRequired runs.
		if !isError(resp) {
			t.Fatal("expected error response when label is missing on macOS")
		}
		if errorCode(resp) != "validation_error" {
			t.Errorf("expected validation_error, got %q", errorCode(resp))
		}
	} else {
		// IsSupported() returns false before validation runs; tool succeeds with
		// a "requires macOS" message.
		if isError(resp) {
			t.Errorf("expected success on non-macOS, got error: %s", errorCode(resp))
		}
		if !requiresMacOSMessage(resp) {
			t.Errorf("expected macOS message on non-macOS, got: %v", resp)
		}
	}
}

func TestFindElement_WithLabel(t *testing.T) {
	resp := callTool(t, FindElement(), map[string]any{"label": "OK"})
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if isMacOS() {
		if isError(resp) && errorCode(resp) != "applescript_error" {
			t.Errorf("unexpected error code %q, want applescript_error or success", errorCode(resp))
		}
	} else {
		if isError(resp) {
			t.Errorf("expected success on non-macOS, got error: %s", errorCode(resp))
		}
		if !requiresMacOSMessage(resp) {
			t.Errorf("expected macOS message on non-macOS, got: %v", resp)
		}
	}
}

func TestFindElement_WithRole(t *testing.T) {
	resp := callTool(t, FindElement(), map[string]any{"label": "OK", "role": "button"})
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if isMacOS() {
		if isError(resp) && errorCode(resp) != "applescript_error" {
			t.Errorf("unexpected error code %q, want applescript_error or success", errorCode(resp))
		}
	} else {
		if isError(resp) {
			t.Errorf("expected success on non-macOS, got error: %s", errorCode(resp))
		}
		if !requiresMacOSMessage(resp) {
			t.Errorf("expected macOS message on non-macOS, got: %v", resp)
		}
	}
}

// --- get_window_elements ---

// TestGetWindowElements_MissingRequired verifies the validation_error path for
// the missing required "window_title" field. Only reachable on macOS.
func TestGetWindowElements_MissingRequired(t *testing.T) {
	resp := callTool(t, GetWindowElements(), nil)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if isMacOS() {
		if !isError(resp) {
			t.Fatal("expected error response when window_title is missing on macOS")
		}
		if errorCode(resp) != "validation_error" {
			t.Errorf("expected validation_error, got %q", errorCode(resp))
		}
	} else {
		if isError(resp) {
			t.Errorf("expected success on non-macOS, got error: %s", errorCode(resp))
		}
		if !requiresMacOSMessage(resp) {
			t.Errorf("expected macOS message on non-macOS, got: %v", resp)
		}
	}
}

func TestGetWindowElements_Valid(t *testing.T) {
	resp := callTool(t, GetWindowElements(), map[string]any{"window_title": "Test Window"})
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if isMacOS() {
		if isError(resp) && errorCode(resp) != "applescript_error" {
			t.Errorf("unexpected error code %q, want applescript_error or success", errorCode(resp))
		}
	} else {
		if isError(resp) {
			t.Errorf("expected success on non-macOS, got error: %s", errorCode(resp))
		}
		if !requiresMacOSMessage(resp) {
			t.Errorf("expected macOS message on non-macOS, got: %v", resp)
		}
	}
}

// --- get_element_hierarchy ---

// GetElementHierarchy has no required args (app_name is optional); the tool
// falls back to the frontmost process when app_name is absent.

func TestGetElementHierarchy_NoArgs(t *testing.T) {
	resp := callTool(t, GetElementHierarchy(), nil)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if isMacOS() {
		if isError(resp) && errorCode(resp) != "applescript_error" {
			t.Errorf("unexpected error code %q, want applescript_error or success", errorCode(resp))
		}
	} else {
		if isError(resp) {
			t.Errorf("expected success on non-macOS, got error: %s", errorCode(resp))
		}
		if !requiresMacOSMessage(resp) {
			t.Errorf("expected macOS message on non-macOS, got: %v", resp)
		}
	}
}

func TestGetElementHierarchy_Valid(t *testing.T) {
	resp := callTool(t, GetElementHierarchy(), map[string]any{"app_name": "Finder"})
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if isMacOS() {
		if isError(resp) && errorCode(resp) != "applescript_error" {
			t.Errorf("unexpected error code %q, want applescript_error or success", errorCode(resp))
		}
	} else {
		if isError(resp) {
			t.Errorf("expected success on non-macOS, got error: %s", errorCode(resp))
		}
		if !requiresMacOSMessage(resp) {
			t.Errorf("expected macOS message on non-macOS, got: %v", resp)
		}
	}
}

// --- list_windows ---
//
// Note: the source file list_windows.go is named with the "_windows" OS suffix,
// which causes Go's implicit file-name build constraint to exclude it on all
// non-Windows platforms. Tests for ListWindows therefore live in a separate
// build-constrained file (tools_list_windows_test.go) so that this file
// compiles on all platforms. See tools_list_windows_test.go.
