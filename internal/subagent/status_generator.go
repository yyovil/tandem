package subagent

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/yaydraco/tandem/internal/config"
)

// StatusGenerator generates meaningful status messages for different agent types
type StatusGenerator struct{}

// GenerateStatusText creates dynamic status messages based on agent type and progress
func (sg *StatusGenerator) GenerateStatusText(agentName config.AgentName, task string, progress int) string {
	switch agentName {
	case config.Reconnoiter:
		return sg.generateReconStatus(task, progress)
	case config.VulnerabilityScanner:
		return sg.generateVulnScanStatus(task, progress)
	case config.Exploiter:
		return sg.generateExploitStatus(task, progress)
	case config.Reporter:
		return sg.generateReportStatus(task, progress)
	default:
		return fmt.Sprintf("Processing task: %s", strings.Split(task, "\n")[0])
	}
}

func (sg *StatusGenerator) generateReconStatus(task string, progress int) string {
	phases := []string{
		"Initializing reconnaissance phase...",
		"Identifying target systems and services...",
		"Enumerating open ports and services...",
		"Gathering system information and fingerprinting...",
		"Mapping network topology and discovering hosts...",
		"Analyzing service versions and configurations...",
		"Documenting findings and preparing reconnaissance report...",
		"Reconnaissance phase completed successfully",
	}
	
	if progress < 10 {
		return phases[0]
	} else if progress < 25 {
		return phases[1]
	} else if progress < 40 {
		return phases[2]
	} else if progress < 60 {
		return phases[3]
	} else if progress < 75 {
		return phases[4]
	} else if progress < 90 {
		return phases[5]
	} else if progress < 100 {
		return phases[6]
	}
	return phases[7]
}

func (sg *StatusGenerator) generateVulnScanStatus(task string, progress int) string {
	phases := []string{
		"Initializing vulnerability scanning engine...",
		"Loading vulnerability signatures and patterns...",
		"Scanning for common vulnerabilities (CVEs)...",
		"Analyzing service configurations for weaknesses...",
		"Testing for authentication bypasses and injection flaws...",
		"Checking for privilege escalation opportunities...",
		"Correlating findings and assessing risk levels...",
		"Generating vulnerability assessment report...",
		"Vulnerability scan completed with findings documented",
	}
	
	phaseIndex := min(progress*len(phases)/100, len(phases)-1)
	return phases[phaseIndex]
}

func (sg *StatusGenerator) generateExploitStatus(task string, progress int) string {
	phases := []string{
		"Analyzing identified vulnerabilities for exploitation...",
		"Selecting appropriate exploit techniques and payloads...",
		"Preparing exploitation framework and tools...",
		"Attempting controlled exploitation within RoE boundaries...",
		"Escalating privileges and maintaining access...",
		"Documenting proof-of-concept and impact assessment...",
		"Cleaning up exploitation artifacts and traces...",
		"Exploitation phase completed with documented evidence",
	}
	
	phaseIndex := min(progress*len(phases)/100, len(phases)-1)
	return phases[phaseIndex]
}

func (sg *StatusGenerator) generateReportStatus(task string, progress int) string {
	phases := []string{
		"Gathering findings from all assessment phases...",
		"Analyzing and correlating vulnerability data...",
		"Categorizing findings by severity and impact...",
		"Creating executive summary and technical details...",
		"Generating remediation recommendations...",
		"Formatting report with evidence and screenshots...",
		"Performing quality assurance and technical review...",
		"Finalizing penetration testing report",
	}
	
	phaseIndex := min(progress*len(phases)/100, len(phases)-1)
	return phases[phaseIndex]
}

// GenerateRandomProgress simulates dynamic progress for demo purposes
func (sg *StatusGenerator) GenerateRandomProgress(current int) int {
	if current >= 100 {
		return 100
	}
	
	// Add some randomness to make it feel more realistic
	increment := rand.Intn(15) + 5 // 5-20% increments
	newProgress := current + increment
	
	if newProgress > 100 {
		return 100
	}
	return newProgress
}

// FormatProgress returns a formatted progress string
func (sg *StatusGenerator) FormatProgress(progress int) string {
	if progress <= 0 {
		return "0%"
	}
	if progress >= 100 {
		return "100%"
	}
	return fmt.Sprintf("%d%%", progress)
}

// GetEstimatedTimeRemaining estimates time remaining based on progress and duration
func (sg *StatusGenerator) GetEstimatedTimeRemaining(progress int, startTime time.Time) string {
	if progress <= 0 {
		return "Calculating..."
	}
	if progress >= 100 {
		return "Completed"
	}
	
	elapsed := time.Since(startTime)
	rate := float64(progress) / elapsed.Seconds()
	remaining := time.Duration((100.0-float64(progress))/rate) * time.Second
	
	if remaining > time.Hour {
		return fmt.Sprintf("~%dh %dm", int(remaining.Hours()), int(remaining.Minutes())%60)
	} else if remaining > time.Minute {
		return fmt.Sprintf("~%dm", int(remaining.Minutes()))
	}
	return fmt.Sprintf("~%ds", int(remaining.Seconds()))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}