package message

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/yyovil/tandem/internal/db"
	"github.com/yyovil/tandem/internal/models"
	"github.com/yyovil/tandem/internal/pubsub"
)

type Message struct {
	ID        string
	Role      MessageRole
	SessionID string
	Parts     []ContentPart
	Model     models.ModelID
	CreatedAt int64
	UpdatedAt int64
}

type Service interface {
	pubsub.Subscriber[Message]
	Create(ctx context.Context, sessionID string, params CreateMessageParams) (Message, error)
	Update(ctx context.Context, message Message) error
	Get(ctx context.Context, id string) (Message, error)
	List(ctx context.Context, sessionID string) ([]Message, error)
	Delete(ctx context.Context, id string) error
	DeleteSessionMessages(ctx context.Context, sessionID string) error
}

type service struct {
	*pubsub.Broker[Message]
	q db.Querier
}

type CreateMessageParams struct {
	Role  MessageRole
	Parts []ContentPart
	Model models.ModelID
}

func (s *service) Create(ctx context.Context, sessionID string, params CreateMessageParams) (Message, error) {
	if params.Role != Assistant {
		params.Parts = append(params.Parts, Finish{
			Reason: BecauseStop,
		})
	}
	partsJSON, err := marshallParts(params.Parts)
	if err != nil {
		return Message{}, err
	}
	dbMessage, err := s.q.CreateMessage(ctx, db.CreateMessageParams{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Role:      string(params.Role),
		Parts:     string(partsJSON),
		Model:     sql.NullString{String: string(params.Model), Valid: true},
	})
	if err != nil {
		return Message{}, err
	}
	message, err := s.fromDBItem(dbMessage)
	if err != nil {
		return Message{}, err
	}
	s.Publish(pubsub.CreatedEvent, message)
	return message, nil
}

type partWrapper struct {
	Type partType    `json:"type"`
	Data ContentPart `json:"data"`
}

func marshallParts(parts []ContentPart) ([]byte, error) {
	wrappedParts := make([]partWrapper, len(parts))

	for i, part := range parts {
		var typ partType

		switch part.(type) {
		case BinaryContent:
			typ = binaryType
		case ReasoningContent:
			typ = reasoningType
		case TextContent:
			typ = textType
		case ToolCall:
			typ = toolCallType
		case ToolResult:
			typ = toolResultType
		case Finish:
			typ = finishType
		default:
			return nil, fmt.Errorf("unknown part type: %T", part)
		}

		wrappedParts[i] = partWrapper{
			Type: typ,
			Data: part,
		}
	}
	return json.Marshal(wrappedParts)
}

// NOTE: messages are stringified first and then stored in the db.
func unmarshallParts(data []byte) ([]ContentPart, error) {
	temp := []json.RawMessage{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	parts := make([]ContentPart, 0)

	for _, rawPart := range temp {
		var wrapper struct {
			Type partType        `json:"type"`
			Data json.RawMessage `json:"data"`
		}

		if err := json.Unmarshal(rawPart, &wrapper); err != nil {
			return nil, err
		}

		switch wrapper.Type {
		case reasoningType:
			part := ReasoningContent{}
			if err := json.Unmarshal(wrapper.Data, &part); err != nil {
				return nil, err
			}
			parts = append(parts, part)
		case textType:
			part := TextContent{}
			if err := json.Unmarshal(wrapper.Data, &part); err != nil {
				return nil, err
			}
			parts = append(parts, part)
		case toolCallType:
			part := ToolCall{}
			if err := json.Unmarshal(wrapper.Data, &part); err != nil {
				return nil, err
			}
			parts = append(parts, part)
		case toolResultType:
			part := ToolResult{}
			if err := json.Unmarshal(wrapper.Data, &part); err != nil {
				return nil, err
			}
			parts = append(parts, part)
		case finishType:
			part := Finish{}
			if err := json.Unmarshal(wrapper.Data, &part); err != nil {
				return nil, err
			}
			parts = append(parts, part)
		default:
			return nil, fmt.Errorf("unknown part type: %s", wrapper.Type)
		}

	}

	return parts, nil
}

func (s *service) Update(ctx context.Context, message Message) error {
	parts, err := marshallParts(message.Parts)
	if err != nil {
		return err
	}
	finishedAt := sql.NullInt64{}
	if f := message.FinishPart(); f != nil {
		finishedAt.Int64 = f.Time
		finishedAt.Valid = true
	}
	err = s.q.UpdateMessage(ctx, db.UpdateMessageParams{
		ID:         message.ID,
		Parts:      string(parts),
		FinishedAt: finishedAt,
	})
	if err != nil {
		return err
	}
	s.Publish(pubsub.UpdatedEvent, message)
	return nil
}

func (m *Message) FinishPart() *Finish {
	for _, part := range m.Parts {
		if c, ok := part.(Finish); ok {
			return &c
		}
	}
	return nil
}

func (m *Message) FinishReason() FinishReason {
	for _, part := range m.Parts {
		if c, ok := part.(Finish); ok {
			return c.Reason
		}
	}
	return ""
}

func (s *service) Get(ctx context.Context, id string) (Message, error) {
	dbMessage, err := s.q.GetMessage(ctx, id)
	if err != nil {
		return Message{}, err
	}
	return s.fromDBItem(dbMessage)
}

func (s *service) List(ctx context.Context, sessionID string) ([]Message, error) {
	dbMessages, err := s.q.ListMessagesBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	messages := make([]Message, len(dbMessages))
	for i, dbMessage := range dbMessages {
		messages[i], err = s.fromDBItem(dbMessage)
		if err != nil {
			return nil, err
		}
	}
	return messages, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	message, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	err = s.q.DeleteMessage(ctx, message.ID)
	if err != nil {
		return err
	}
	s.Publish(pubsub.DeletedEvent, message)
	return nil
}

func (s *service) DeleteSessionMessages(ctx context.Context, sessionID string) error {
	messages, err := s.List(ctx, sessionID)
	if err != nil {
		return err
	}
	for _, message := range messages {
		if message.SessionID == sessionID {
			err = s.Delete(ctx, message.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *service) fromDBItem(item db.Message) (Message, error) {
	parts, err := unmarshallParts([]byte(item.Parts))
	if err != nil {
		return Message{}, err
	}
	return Message{
		ID:        item.ID,
		SessionID: item.SessionID,
		Role:      MessageRole(item.Role),
		Parts:     parts,
	}, nil
}

func NewService(q db.Querier) Service {
	return &service{
		Broker: pubsub.NewBroker[Message](),
		q:      q,
	}
}

func (m *Message) ToolCalls() []ToolCall {
	toolCalls := make([]ToolCall, 0)
	for _, part := range m.Parts {
		if c, ok := part.(ToolCall); ok {
			toolCalls = append(toolCalls, c)
		}
	}
	return toolCalls
}

func (m *Message) Content() TextContent {
	for _, part := range m.Parts {
		if c, ok := part.(TextContent); ok {
			return c
		}
	}
	return TextContent{}
}

func (m *Message) BinaryContent() []BinaryContent {
	binaryContents := make([]BinaryContent, 0)
	for _, part := range m.Parts {
		if c, ok := part.(BinaryContent); ok {
			binaryContents = append(binaryContents, c)
		}
	}
	return binaryContents
}

func (m *Message) ToolResults() []ToolResult {
	toolResults := make([]ToolResult, 0)
	for _, part := range m.Parts {
		if c, ok := part.(ToolResult); ok {
			toolResults = append(toolResults, c)
		}
	}
	return toolResults
}

func (m *Message) IsFinished() bool {
	for _, part := range m.Parts {
		if _, ok := part.(Finish); ok {
			return true
		}
	}
	return false
}

func (m *Message) IsThinking() bool {
	if m.ReasoningContent().Thinking != "" && m.Content().Text == "" && !m.IsFinished() {
		return true
	}
	return false
}

func (m *Message) ReasoningContent() ReasoningContent {
	for _, part := range m.Parts {
		if c, ok := part.(ReasoningContent); ok {
			return c
		}
	}
	return ReasoningContent{}
}
