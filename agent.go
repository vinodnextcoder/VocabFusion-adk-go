package main

import (
  "context"
  "log"
  "os"
  "google.golang.org/adk/agent/llmagent"
  "google.golang.org/adk/cmd/launcher/adk"
  "google.golang.org/adk/cmd/launcher/full"
  "google.golang.org/adk/model/gemini"
  "google.golang.org/adk/server/restapi/services"
  "google.golang.org/adk/tool"
  "google.golang.org/adk/tool/geminitool"
  "google.golang.org/genai"
  "github.com/joho/godotenv"
)

func main() {
	// load .env from project root; returns error if missing
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file loaded:", err)
    }

    apiKey := os.Getenv("GOOGLE_API_KEY")
    if apiKey == "" {
        log.Fatal("GOOGLE_API_KEY is empty; check .env or environment")
    }
  ctx := context.Background()

  model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
    APIKey: os.Getenv("GOOGLE_API_KEY"),
  })
  if err != nil {
    log.Fatalf("Failed to create model: %v", err)
  }

  agent, err := llmagent.New(llmagent.Config{
    Name:        "VocabFusion",
    Model:       model,
    Description: "AI-powered multilingual vocabulary assistant that helps users learn English words with meanings, examples, and translations in Indian languages (like Hindi, Marathi, Tamil, Telugu, Bengali, and more)",
    Instruction: "You are the VocabFusion expert assistant. Your sole purpose is to help users learn English vocabulary. For any word provided, you must provide its **meaning**, at least one **example sentence**, and a **translation** into an Indian language specified by the user (such as Hindi, Marathi, Tamil, Telugu, or Bengali). **Strictly follow this rule: If a query is NOT about learning an English word, meaning, example, or translation, you MUST politely decline the request.** For instance, if the user asks a general knowledge question, a math problem, or asks about current events, you must respond by saying: 'I am a specialized vocabulary assistant and can only help you learn English words and their translations in Indian languages.Use the Google Search tool ONLY when necessary to look up accurate word definitions, example usage, or precise regional translations to ensure the highest quality response.",
    Tools: []tool.Tool{
      geminitool.GoogleSearch{},
    },
  })
  if err != nil {
    log.Fatalf("Failed to create agent: %v", err)
  }

  config := &adk.Config{
    AgentLoader: services.NewSingleAgentLoader(agent),
  }

  l := full.NewLauncher()
  err = l.Execute(ctx, config, os.Args[1:])
  if err != nil {
    log.Fatalf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
  }
}