package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/saint0x/file-storage-app/backend/internal/models"
	"github.com/sashabaranov/go-openai"
)

type Processor struct {
	client *openai.Client
}

func NewProcessor(apiKey string) *Processor {
	return &Processor{
		client: openai.NewClient(apiKey),
	}
}

type FileOrganizationRequest struct {
	Files []models.File `json:"files"`
}

type FileOrganizationResponse struct {
	Folders []struct {
		Name  string   `json:"name"`
		Files []string `json:"files"`
	} `json:"folders"`
}

func (p *Processor) OrganizeFiles(ctx context.Context, req FileOrganizationRequest) (*FileOrganizationResponse, error) {
	prompt := fmt.Sprintf(`Organize the following files into folders that emphasize calm, comfort, and energy. The organization should feel "breathable":

Files:
%s

Respond with a JSON object containing an array of folders, each with a name and an array of file names. For example:
{
  "folders": [
    {
      "name": "Calm Workspace",
      "files": ["document1.pdf", "image2.jpg"]
    },
    {
      "name": "Energizing Projects",
      "files": ["project3.docx", "presentation4.pptx"]
    }
  ]
}`, formatFilesForPrompt(req.Files))

	resp, err := p.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error creating chat completion: %w", err)
	}

	var result FileOrganizationResponse
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling AI response: %w", err)
	}

	return &result, nil
}

func formatFilesForPrompt(files []models.File) string {
	var result string
	for _, file := range files {
		result += fmt.Sprintf("- %s (%s)\n", file.Name, file.ContentType)
	}
	return result
}

// Add this method to the Processor struct
func (p *Processor) Close() error {
	// Add any cleanup logic here
	return nil
}
