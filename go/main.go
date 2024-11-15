// main.go
package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Memory struct {
	Content           string    `json:"content"`
	Timestamp         time.Time `json:"timestamp"`
	InitialImportance float64   `json:"initialImportance"`
	UseCount          int       `json:"useCount"`
}

type Message struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatResponse struct {
	Message      Message           `json:"message"`
	Memories     map[string]Memory `json:"memories"`
	ForgottenKey string            `json:"forgottenKey,omitempty"`
}

type Server struct {
	memories map[string]Memory
	mutex    sync.RWMutex
}

func NewServer() *Server {
	server := &Server{
		memories: make(map[string]Memory),
	}

	// Start forgetting process
	go server.startForgettingProcess()

	return server
}

func (s *Server) startForgettingProcess() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		s.forgetMemories()
	}
}

func (s *Server) forgetMemories() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for key, memory := range s.memories {
		importance := s.calculateImportance(memory)
		if rand.Float64() > importance/100 {
			delete(s.memories, key)
		}
	}
}

func (s *Server) calculateImportance(memory Memory) float64 {
	age := time.Since(memory.Timestamp).Hours()
	ageWeight := math.Exp(-age / 24) // Decay over 24 hours
	return memory.InitialImportance * ageWeight * math.Log(float64(memory.UseCount+1))
}

func (s *Server) handleMessage(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Process message and update memories
	response := s.processMessage(input.Content)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) processMessage(content string) ChatResponse {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Create new message
	message := Message{
		ID:        time.Now().Format(time.RFC3339Nano),
		Type:      "ai",
		Content:   "",
		Timestamp: time.Now(),
	}

	// Process input and update memories
	words := splitIntoWords(content)
	var rememberedContent string

	for _, word := range words {
		if len(word) > 3 {
			if memory, exists := s.memories[word]; exists {
				rememberedContent = memory.Content
				memory.UseCount++
				s.memories[word] = memory
			} else {
				s.memories[word] = Memory{
					Content:           content,
					Timestamp:         time.Now(),
					InitialImportance: 50 + rand.Float64()*50,
					UseCount:          1,
				}
			}
		}
	}

	// Generate response
	if rememberedContent != "" {
		message.Content = "以前の会話を思い出しました: " + rememberedContent
	} else {
		message.Content = "申し訳ありません。関連する記憶が曖昧です..."
	}

	return ChatResponse{
		Message:  message,
		Memories: s.memories,
	}
}

func splitIntoWords(text string) []string {
	// Implement word splitting logic here
	// This is a simple implementation - you might want to use a proper tokenizer
	return strings.Fields(text)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	server := NewServer()
	router := mux.NewRouter()

	router.HandleFunc("/api/chat", server.handleMessage).Methods("POST")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})

	handler := c.Handler(router)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
