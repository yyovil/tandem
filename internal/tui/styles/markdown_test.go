package styles

import (
	"fmt"
	"strings"
	"testing"
)

func TestTableMarkdownRenderingComparison(t *testing.T) {
	// Test markdown with a table similar to the one shown in the issue
	testMarkdown := `
# Vulnerability Report

| Severity | Vulnerability | Impact | Recommendation |
|----------|---------------|---------|----------------|
| High | Samba SMB 3.x Upgrade | Allows local attacker to escalate privileges to root. | Upgrade Samba to a patched version is_known_pipename |
| Medium | CUPS 1.7 Multiple Vulnerabilities (CVEs varying) | Known vulnerabilities related to arbitrary file reading and privilege escalation. | Upgrade CUPS to the latest stable version. Restrict access to web interface. |
| High | WEBrick httpd 1.3.1 Directory Traversal | HTTP (WEBrick httpd | Upgrade Ruby and WEBrick. Implement strict input |

This table demonstrates the improved formatting after the fix.
`

	// Test with improved renderer
	renderer := GetMarkdownRenderer(80)
	if renderer == nil {
		t.Fatal("Failed to create markdown renderer")
	}

	rendered, err := renderer.Render(testMarkdown)
	if err != nil {
		t.Fatalf("Failed to render markdown: %v", err)
	}

	fmt.Println("IMPROVED table rendering:")
	fmt.Println("=" + string(make([]rune, 80))[:79])
	fmt.Println(rendered)
	fmt.Println("=" + string(make([]rune, 80))[:79])

	// Basic checks that the table was rendered properly
	if len(rendered) == 0 {
		t.Error("Rendered markdown is empty")
	}

	// Check that we have table characters (improved formatting should include these)
	if !containsAny(rendered, []string{"│", "─", "┼"}) {
		t.Error("Rendered table does not contain expected table formatting characters")
	}

	// Check that we have proper structure - headers should be present
	if !containsAny(rendered, []string{"Severity", "Vulnerability", "Impact", "Recommendation"}) {
		t.Error("Rendered table does not contain expected headers")
	}
}

func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}