package subagent

import (
	"context"
	"testing"

	"github.com/yaydraco/tandem/internal/config"
)

func TestActivityService(t *testing.T) {
	service := NewService()
	ctx := context.Background()

	// Test starting an activity
	activity, err := service.StartActivity(ctx, "session1", "parent1", config.Reconnoiter, "Test reconnaissance task")
	if err != nil {
		t.Fatalf("Failed to start activity: %v", err)
	}

	if activity.ID != "session1" {
		t.Errorf("Expected activity ID to be 'session1', got %s", activity.ID)
	}

	if activity.AgentName != config.Reconnoiter {
		t.Errorf("Expected agent name to be Reconnoiter, got %s", activity.AgentName)
	}

	if activity.Status != StatusStarting {
		t.Errorf("Expected status to be starting, got %s", activity.Status)
	}

	// Test updating activity
	err = service.UpdateActivity(ctx, activity.ID, StatusRunning, "Scanning targets...", "50%")
	if err != nil {
		t.Fatalf("Failed to update activity: %v", err)
	}

	// Test getting activity
	retrieved, err := service.GetActivity(ctx, activity.ID)
	if err != nil {
		t.Fatalf("Failed to get activity: %v", err)
	}

	if retrieved.Status != StatusRunning {
		t.Errorf("Expected status to be running, got %s", retrieved.Status)
	}

	if retrieved.Progress != "50%" {
		t.Errorf("Expected progress to be 50%%, got %s", retrieved.Progress)
	}

	// Test completing activity
	err = service.CompleteActivity(ctx, activity.ID, true, "Scan completed successfully")
	if err != nil {
		t.Fatalf("Failed to complete activity: %v", err)
	}

	retrieved, err = service.GetActivity(ctx, activity.ID)
	if err != nil {
		t.Fatalf("Failed to get activity: %v", err)
	}

	if retrieved.Status != StatusCompleted {
		t.Errorf("Expected status to be completed, got %s", retrieved.Status)
	}

	if retrieved.CanAbort {
		// Should not be able to abort completed tasks
		t.Errorf("Expected CanAbort to be false for completed task")
	}
}

func TestActivityAbort(t *testing.T) {
	service := NewService()
	ctx := context.Background()

	// Start an activity
	activity, err := service.StartActivity(ctx, "session2", "parent2", config.VulnerabilityScanner, "Test vulnerability scan")
	if err != nil {
		t.Fatalf("Failed to start activity: %v", err)
	}

	// Set it to running
	err = service.UpdateActivity(ctx, activity.ID, StatusRunning, "Scanning for vulnerabilities...", "25%")
	if err != nil {
		t.Fatalf("Failed to update activity: %v", err)
	}

	// Test aborting
	err = service.AbortActivity(ctx, activity.ID)
	if err != nil {
		t.Fatalf("Failed to abort activity: %v", err)
	}

	retrieved, err := service.GetActivity(ctx, activity.ID)
	if err != nil {
		t.Fatalf("Failed to get activity: %v", err)
	}

	if retrieved.Status != StatusAborted {
		t.Errorf("Expected status to be aborted, got %s", retrieved.Status)
	}

	if retrieved.CanAbort {
		t.Errorf("Expected CanAbort to be false for aborted task")
	}
}

func TestGetActiveActivities(t *testing.T) {
	service := NewService()
	ctx := context.Background()

	// Start multiple activities
	activity1, _ := service.StartActivity(ctx, "session1", "parent1", config.Reconnoiter, "Task 1")
	activity2, _ := service.StartActivity(ctx, "session2", "parent1", config.VulnerabilityScanner, "Task 2")
	_, _ = service.StartActivity(ctx, "session3", "parent1", config.Exploiter, "Task 3")

	// Set one to running
	service.UpdateActivity(ctx, activity1.ID, StatusRunning, "Running...", "50%")

	// Complete one
	service.CompleteActivity(ctx, activity2.ID, true, "Done")

	// Get active activities
	activeActivities := service.GetActiveActivities(ctx)

	// Should have 2 active activities (starting and running)
	expectedCount := 2
	if len(activeActivities) != expectedCount {
		t.Errorf("Expected %d active activities, got %d", expectedCount, len(activeActivities))
	}

	// Check that the completed one is not in the list
	for _, activity := range activeActivities {
		if activity.ID == activity2.ID {
			t.Error("Completed activity should not be in active activities list")
		}
	}
}