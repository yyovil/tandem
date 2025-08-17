package styles

import (
	"fmt"
	"testing"
)

func TestTableMarkdownRendering(t *testing.T) {
	// Test markdown with a table similar to the one shown in the issue
	testMarkdown := `
# Vulnerability Report

| Severity | Vulnerability | Impact | Recommendation |
|----------|---------------|---------|----------------|
| High | Samba SMB 3.x Upgrade | Allows local attacker to escalate privileges to root. | Upgrade Samba to a patched version is_known_pipename |
| Medium | CUPS 1.7 Multiple Vulnerabilities (CVEs varying) | Known vulnerabilities related to arbitrary file reading and privilege escalation. | Upgrade CUPS to the latest stable version. Restrict access to web interface. |
| High | WEBrick httpd 1.3.1 Directory Traversal | HTTP (WEBrick httpd | Upgrade Ruby and WEBrick. Implement strict input |

This table demonstrates the formatting issue described in the GitHub issue.
`

	// Test with current renderer
	renderer := GetMarkdownRenderer(80)
	if renderer == nil {
		t.Fatal("Failed to create markdown renderer")
	}

	rendered, err := renderer.Render(testMarkdown)
	if err != nil {
		t.Fatalf("Failed to render markdown: %v", err)
	}

	fmt.Println("Current table rendering:")
	fmt.Println("=" + string(make([]rune, 80))[:79])
	fmt.Println(rendered)
	fmt.Println("=" + string(make([]rune, 80))[:79])

	// Basic checks that a table was rendered
	if len(rendered) == 0 {
		t.Error("Rendered markdown is empty")
	}

	// The rendered content should contain table elements
	// Note: We're not checking for perfect formatting yet, just that it renders something
	if len(rendered) < len(testMarkdown)/2 {
		t.Error("Rendered markdown seems too short, likely rendering failed")
	}
}