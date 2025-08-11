package subagent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yaydraco/tandem/internal/config"
	"github.com/yaydraco/tandem/internal/pubsub"
)

// ActivityStatus represents the current status of a subagent
type ActivityStatus string

const (
	StatusStarting  ActivityStatus = "starting"
	StatusRunning   ActivityStatus = "running"
	StatusCompleted ActivityStatus = "completed"
	StatusError     ActivityStatus = "error"
	StatusAborted   ActivityStatus = "aborted"
)

// Activity represents a single subagent task
type Activity struct {
	ID              string            `json:"id"`
	SessionID       string            `json:"session_id"`
	ParentID        string            `json:"parent_id"`
	AgentName       config.AgentName  `json:"agent_name"`
	Task            string            `json:"task"`
	Status          ActivityStatus    `json:"status"`
	StatusText      string            `json:"status_text"`
	Progress        string            `json:"progress"`
	ProgressPercent int               `json:"progress_percent"`
	StartedAt       time.Time         `json:"started_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	CompletedAt     *time.Time        `json:"completed_at,omitempty"`
	Error           string            `json:"error,omitempty"`
	CanAbort        bool              `json:"can_abort"`
	EstimatedTime   string            `json:"estimated_time,omitempty"`
}

// ActivityEvent represents events published by the activity service
type ActivityEvent struct {
	Type     string    `json:"type"`
	Activity Activity  `json:"activity"`
}

// Service manages subagent activities
type Service interface {
	pubsub.Subscriber[ActivityEvent]
	
	// StartActivity creates and tracks a new subagent activity
	StartActivity(ctx context.Context, sessionID, parentID string, agentName config.AgentName, task string) (*Activity, error)
	
	// UpdateActivity updates the status and progress of an activity
	UpdateActivity(ctx context.Context, activityID string, status ActivityStatus, statusText, progress string) error
	
	// CompleteActivity marks an activity as completed
	CompleteActivity(ctx context.Context, activityID string, success bool, result string) error
	
	// AbortActivity cancels a running activity
	AbortActivity(ctx context.Context, activityID string) error
	
	// GetActiveActivities returns all currently active activities
	GetActiveActivities(ctx context.Context) []Activity
	
	// GetActivity returns a specific activity by ID
	GetActivity(ctx context.Context, activityID string) (*Activity, error)
	
	// IsActivityActive checks if an activity is currently running
	IsActivityActive(ctx context.Context, activityID string) bool
	
	// SetCancelFunc stores a cancel function for an activity
	SetCancelFunc(activityID string, cancelFunc context.CancelFunc)
}

type service struct {
	*pubsub.Broker[ActivityEvent]
	activities      map[string]*Activity
	cancelFuncs     map[string]context.CancelFunc
	statusGenerator *StatusGenerator
	mu              sync.RWMutex
}

func (s *service) StartActivity(ctx context.Context, sessionID, parentID string, agentName config.AgentName, task string) (*Activity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	activity := &Activity{
		ID:              sessionID, // Use session ID as activity ID for simplicity
		SessionID:       sessionID,
		ParentID:        parentID,
		AgentName:       agentName,
		Task:            task,
		Status:          StatusStarting,
		StatusText:      s.statusGenerator.GenerateStatusText(agentName, task, 0),
		Progress:        "0%",
		ProgressPercent: 0,
		StartedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CanAbort:        true,
		EstimatedTime:   "Calculating...",
	}
	
	s.activities[activity.ID] = activity
	
	s.Publish(pubsub.CreatedEvent, ActivityEvent{
		Type:     "activity_started",
		Activity: *activity,
	})
	
	return activity, nil
}

func (s *service) UpdateActivity(ctx context.Context, activityID string, status ActivityStatus, statusText, progress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	activity, exists := s.activities[activityID]
	if !exists {
		return nil // Activity not found, silently ignore
	}
	
	// Parse progress percentage
	progressPercent := activity.ProgressPercent
	if progress != "" {
		// Try to extract percentage from progress string
		var parsed int
		if n, err := fmt.Sscanf(progress, "%d%%", &parsed); n == 1 && err == nil {
			progressPercent = parsed
		}
	}
	
	// Generate smart status text if not provided
	if statusText == "" {
		statusText = s.statusGenerator.GenerateStatusText(activity.AgentName, activity.Task, progressPercent)
	}
	
	activity.Status = status
	activity.StatusText = statusText
	activity.Progress = progress
	activity.ProgressPercent = progressPercent
	activity.UpdatedAt = time.Now()
	activity.EstimatedTime = s.statusGenerator.GetEstimatedTimeRemaining(progressPercent, activity.StartedAt)
	
	s.Publish(pubsub.UpdatedEvent, ActivityEvent{
		Type:     "activity_updated",
		Activity: *activity,
	})
	
	return nil
}

func (s *service) CompleteActivity(ctx context.Context, activityID string, success bool, result string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	activity, exists := s.activities[activityID]
	if !exists {
		return nil // Activity not found, silently ignore
	}
	
	now := time.Now()
	activity.UpdatedAt = now
	activity.CompletedAt = &now
	activity.CanAbort = false
	
	if success {
		activity.Status = StatusCompleted
		activity.StatusText = "Task completed successfully"
		activity.Progress = "100%"
	} else {
		activity.Status = StatusError
		activity.StatusText = "Task failed"
		activity.Error = result
	}
	
	s.Publish(pubsub.UpdatedEvent, ActivityEvent{
		Type:     "activity_completed",
		Activity: *activity,
	})
	
	// Remove from active activities after a delay
	go func() {
		time.Sleep(30 * time.Second)
		s.mu.Lock()
		delete(s.activities, activityID)
		s.mu.Unlock()
	}()
	
	return nil
}

func (s *service) AbortActivity(ctx context.Context, activityID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	activity, exists := s.activities[activityID]
	if !exists {
		return nil // Activity not found, silently ignore
	}
	
	if !activity.CanAbort {
		return nil // Cannot abort this activity
	}
	
	// Call the cancel function if available
	if cancelFunc, exists := s.cancelFuncs[activityID]; exists {
		cancelFunc()
		delete(s.cancelFuncs, activityID)
	}
	
	now := time.Now()
	activity.Status = StatusAborted
	activity.StatusText = "Task aborted by user"
	activity.UpdatedAt = now
	activity.CompletedAt = &now
	activity.CanAbort = false
	
	s.Publish(pubsub.UpdatedEvent, ActivityEvent{
		Type:     "activity_aborted",
		Activity: *activity,
	})
	
	return nil
}

func (s *service) SetCancelFunc(activityID string, cancelFunc context.CancelFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cancelFuncs[activityID] = cancelFunc
}

func (s *service) GetActiveActivities(ctx context.Context) []Activity {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var activities []Activity
	for _, activity := range s.activities {
		if activity.Status == StatusStarting || activity.Status == StatusRunning {
			activities = append(activities, *activity)
		}
	}
	
	return activities
}

func (s *service) GetActivity(ctx context.Context, activityID string) (*Activity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if activity, exists := s.activities[activityID]; exists {
		return activity, nil
	}
	
	return nil, nil
}

func (s *service) IsActivityActive(ctx context.Context, activityID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if activity, exists := s.activities[activityID]; exists {
		return activity.Status == StatusStarting || activity.Status == StatusRunning
	}
	
	return false
}

func NewService() Service {
	return &service{
		Broker:          pubsub.NewBroker[ActivityEvent](),
		activities:      make(map[string]*Activity),
		cancelFuncs:     make(map[string]context.CancelFunc),
		statusGenerator: &StatusGenerator{},
	}
}