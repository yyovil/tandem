package styles

import (
	"fmt"
	"os"
	"testing"
)

func TestActualTableRenderingDemo(t *testing.T) {
	testMarkdown := `
# Vulnerability Assessment Report

## Discovered Vulnerabilities

| Severity | Vulnerability | Impact | Recommendation |
|----------|---------------|---------|----------------|
| High | Samba SMB 3.x Local Privilege Escalation (CVE-2016-2124) | Allows local attacker to escalate privileges to root. | Upgrade Samba to a patched version, disable is_known_pipename feature |
| Medium | CUPS 1.7 Multiple Vulnerabilities (CVEs varying) | Known vulnerabilities related to arbitrary file reading and privilege escalation. | Upgrade CUPS to the latest stable version. Restrict access to web interface. |
| High | WEBrick httpd 1.3.1 Directory Traversal | HTTP (WEBrick httpd 1.3.1) directory traversal vulnerability | Upgrade Ruby and WEBrick. Implement strict input validation |

## Summary

This vulnerability scan discovered **3 vulnerabilities** across the target system, with **2 rated as High severity** and **1 as Medium severity**.
`

	renderer := GetMarkdownRenderer(100)
	if renderer == nil {
		fmt.Println("Failed to create markdown renderer")
		os.Exit(1)
	}

	rendered, err := renderer.Render(testMarkdown)
	if err != nil {
		t.Fatalf("Failed to render markdown: %v", err)
	}

	fmt.Println("Improved Table Rendering Demo:")
	fmt.Println("=" + string(make([]rune, 100))[:99])
	fmt.Println(rendered)
	fmt.Println("=" + string(make([]rune, 100))[:99])

	// Basic validation
	if len(rendered) == 0 {
		t.Error("Rendered markdown is empty")
	}
}