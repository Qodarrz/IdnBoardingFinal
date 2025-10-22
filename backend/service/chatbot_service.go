package service

import (
	"context"
	"errors"
	"os"

	"github.com/joho/godotenv"
	genai "google.golang.org/genai"
)

var ErrMissingAPIKey = errors.New("missing GOOGLE_API_KEY environment variable")

type GeminiService struct {
	client *genai.Client
	ctx    context.Context
}

func NewGeminiService() (*GeminiService, error) {
	_ = godotenv.Load()

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &GeminiService{client: client, ctx: ctx}, nil
}

func (s *GeminiService) Generate(prompt string) (string, error) {
	systemInstruction := `
Kamu adalah asisten untuk aplikasi Greenflow.
Greenflow adalah aplikasi web untuk tracking karbon, misi lingkungan, penukaran poin, dan sistem badge.
Jawaban harus selalu relevan dengan konteks Greenflow: users, activity logs, missions, carbon vehicles, electronics, store, dan points system.
Jangan jawab di luar domain ini.
`

	contentText := systemInstruction + "\n\nPertanyaan: " + prompt

	resp, err := s.client.Models.GenerateContent(
		s.ctx,
		"gemini-2.5-pro",
		[]*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{Text: contentText},
				},
			},
		},
		nil,
	)

	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "Tidak ada jawaban yang dihasilkan.", nil
	}

	return resp.Candidates[0].Content.Parts[0].Text, nil
}
