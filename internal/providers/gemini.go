package providers

import (
	"context"
	"log"
	"os"

	"github.com/yyovil/tandem/internal/agent"
	"github.com/yyovil/tandem/internal/settings"
	"github.com/yyovil/tandem/internal/tools"
	"github.com/yyovil/tandem/internal/utils"
	"google.golang.org/genai"
)

type GeminiProvider struct {
	Client *genai.Client
}

func NewGeminiProvider() (GeminiProvider, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})

	if err != nil {
		log.Println("Error creating Gemini client:", err.Error())
		return GeminiProvider{}, err
	}

	return GeminiProvider{
		Client: client,
	}, nil
}

// returns messages in gemini specific api from tandem specified api format.
func (g GeminiProvider) FromMessages(history []agent.Message) any {
	var providerHistory []*genai.Content

	for _, message := range history {
		var content *genai.Content
		switch message.Role {
		case agent.RoleUser:
			content.Role = genai.RoleUser
			// prompt in the user message
			content.Parts = append(content.Parts, genai.NewPartFromText(message.Part.Text))
			// files in the user message
			for _, file := range message.Files {
				content.Parts = append(content.Parts, genai.NewPartFromBytes(file.Data, file.MimeType))
			}
			providerHistory = append(providerHistory, content)

		case agent.RoleAssistant:
			content.Role = genai.RoleModel
			// text content in the assistant message
			if message.Part.Text != "" {
				content.Parts = append(content.Parts, genai.NewPartFromText(message.Part.Text))
			}

			// tool calls in the assistant message
			if len(message.Part.ToolCalls) > 0 {
				for _, toolCall := range message.Part.ToolCalls {
					content.Parts = append(content.Parts, genai.NewPartFromFunctionCall(string(toolCall.Name), toolCall.Args))
				}
			}

			// NOTE: i guess this avoids the message with finish reason with zero content.
			if len(content.Parts) > 0 {
				providerHistory = append(providerHistory, content)
			}

		case agent.RoleTool:
			for _, toolResult := range message.Part.ToolResult {
				content := genai.NewContentFromFunctionResponse(string(toolResult.Name), toolResult.Result, genai.RoleUser)
				providerHistory = append(providerHistory, content)
			}
		}
	}

	return providerHistory
}

func (g GeminiProvider) GetStream(ctx context.Context, messages []agent.Message, settings settings.Settings) <-chan agent.Message {

	geminiMessages := g.FromMessages(messages).([]*genai.Content)
	history := geminiMessages[:len(geminiMessages)-1]
	// NOTE: last message is always going to be user message.
	lastMessage := geminiMessages[len(geminiMessages)-1]

	providerTools := g.GetToolsForProvider(settings.Tools).([]*genai.Tool)

	contentConfig := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				genai.NewPartFromText(settings.GetSystemPrompt()),
			},
		},
	}

	if len(providerTools) > 0 {
		contentConfig.Tools = providerTools
	}

	chat, err := g.Client.Chats.Create(ctx, string(settings.Model), contentConfig, history)
	if err != nil {
		log.Println("Error creating Gemini chat:", err.Error())
		// ADHD: what we should really do here in case a service critical error occurs in here? and i think there're like many different ways this could go wrong.
	}

	messageChan := make(chan agent.Message)

	go func() {
		defer close(messageChan)
		var parts []genai.Part
		for _, part := range lastMessage.Parts {
			parts = append(parts, *part)
		}
		for content, err := range chat.SendMessageStream(ctx, parts...) {
			if err != nil {
				log.Println("error occured while SendMessageStream: ", err)
			}

			message := g.ToMessage(*content.Candidates[0])
			messageChan <- message
		}
	}()

	return messageChan
}

func (g GeminiProvider) GetToolsForProvider(toolNames []tools.ToolName) any {
	geminiTool := &genai.Tool{}
	geminiTool.FunctionDeclarations = make([]*genai.FunctionDeclaration, 0, len(toolNames))
	for _, toolName := range toolNames {
		tool, found := tools.GetTool(toolName)
		if !found {
			log.Println("Tool not found:", toolName)
			continue
		}

		params := g.FromParameters(tool.Parameters).(map[string]*genai.Schema)
		decl := &genai.FunctionDeclaration{
			Name:        string(tool.Name),
			Description: tool.Description,
			Parameters: &genai.Schema{
				Type:       "object",
				Properties: params,
				Required:   tool.Required,
			},
		}

		geminiTool.FunctionDeclarations = append(geminiTool.FunctionDeclarations, decl)
	}

	return []*genai.Tool{geminiTool}
}

// returns message from gemini specific api to tandem specified api format.
func (g GeminiProvider) ToMessage(candidate any) agent.Message {
	_candidate := candidate.(genai.Candidate)
	var geminiContent = *_candidate.Content

	message := agent.Message{
		Role: agent.RoleUser,
	}

	message.TokenCount += _candidate.TokenCount

	if _candidate.FinishReason != "" {
		message.Type = agent.ResponseCompletedMsg
		switch _candidate.FinishReason {
		case genai.FinishReasonStop:
			message.FinishReason = agent.BecauseStop
		case genai.FinishReasonMaxTokens:
			message.FinishReason = agent.BecauseMaxTokens
		case genai.FinishReasonProhibitedContent:
			message.FinishReason = agent.BecausePermissionDenied
		default:
			message.FinishReason = agent.BecauseUnknown
		}
	}

	for _, part := range geminiContent.Parts {
		switch {
		case part.Text != "":
			message.Type = agent.ResponseMsg
			message.Part.Text = part.Text
		case part.FunctionCall != nil:
			toolCall := tools.ToolCall{
				Name: tools.ToolName(part.FunctionCall.Name),
				Id:   part.FunctionCall.ID,
				Args: part.FunctionCall.Args,
			}
			message.Type = agent.ToolCallMsg
			message.Part.ToolCalls = append(message.Part.ToolCalls, toolCall)
			// NOTE: we can add thoughts if ever need. ig its pretty much standard so we got to support them very soon. soon as we get the api key for models supporting thoughts.
		}
	}
	return message
}

// FromParameters converts provider-specific parameters to the tandem agent's parameter format.
func (g GeminiProvider) FromParameters(params tools.ToolParameters) any {
	genaiParams := make(map[string]*genai.Schema, len(params))
	for key, value := range params {
		genaiParams[key] = g.ToSchema(value).(*genai.Schema)
	}
	return genaiParams
}

func (g GeminiProvider) ToSchema(param tools.Param) any {
	schema := &genai.Schema{}
	schema.Description = param.Description
	// NOTE: this is some conning right in here cuz i copied it straight from the genai docs.
	schema.Type = genai.Type(param.Type)

	switch param.Type {
	case utils.TypeArray:
		schema.Type = genai.TypeArray
		schema.Items = g.ToSchema(*param.Items).(*genai.Schema)
	case utils.TypeObject:
		schema.Properties = g.FromParameters(param.Properties).(map[string]*genai.Schema)
	}

	return schema
}
