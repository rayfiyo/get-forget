package main

import (
	"math"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shogo82148/go-mecab"
)

type Memory struct {
	Content           string    `json:"content"`
	Timestamp         time.Time `json:"timestamp"`
	InitialImportance float64   `json:"initialImportance"`
	UseCount          int       `json:"useCount"`
	Keywords          []string  `json:"keywords"`
}

type MemoryStore struct {
	Memories map[string]Memory
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Memories: make(map[string]Memory),
	}
}

// 形態素解析を行い、重要な単語を抽出する
func extractKeywords(text string) ([]string, error) {
	tagger, err := mecab.New(map[string]string{
		"output-format-type": "wakati",
	})
	if err != nil {
		return nil, err
	}
	defer tagger.Destroy()

	node, err := tagger.ParseToNode(text)
	if err != nil {
		return nil, err
	}

	var keywords []string
	for ; node != (mecab.Node{}); node = node.Next() {
		features := strings.Split(node.Feature(), ",")
		if len(features) > 0 {
			pos := features[0] // 品詞情報
			if pos == "名詞" || pos == "動詞" || pos == "形容詞" {
				if len(node.Surface()) > 1 { // 1文字以上の単語のみ
					keywords = append(keywords, node.Surface())
				}
			}
		}
	}
	return keywords, nil
}

// 重要度の計算
func calculateImportance(memory Memory) float64 {
	age := time.Since(memory.Timestamp).Hours()
	ageWeight := math.Exp(-age / 24) // 24時間で重要度が半減
	useCount := float64(memory.UseCount)
	return memory.InitialImportance * ageWeight * math.Log(useCount+1)
}

func main() {
	store := NewMemoryStore()
	router := gin.Default()

	// CORSの設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// メモリの保存
	router.POST("/memory", func(c *gin.Context) {
		var input struct {
			Content string `json:"content"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		keywords, err := extractKeywords(input.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		memory := Memory{
			Content:           input.Content,
			Timestamp:         time.Now(),
			InitialImportance: 50 + rand.Float64()*50, // 50-100の重要度
			UseCount:          1,
			Keywords:          keywords,
		}

		for _, keyword := range keywords {
			if existing, exists := store.Memories[keyword]; exists {
				existing.UseCount++
				store.Memories[keyword] = existing
			} else {
				store.Memories[keyword] = memory
			}
		}

		c.JSON(http.StatusOK, memory)
	})

	// 関連メモリの取得
	router.GET("/memory", func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
			return
		}

		keywords, err := extractKeywords(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var relevantMemories []Memory
		for _, keyword := range keywords {
			if memory, exists := store.Memories[keyword]; exists {
				if calculateImportance(memory) > rand.Float64()*100 { // ランダムな忘却
					relevantMemories = append(relevantMemories, memory)
				}
			}
		}

		c.JSON(http.StatusOK, relevantMemories)
	})

	router.Run(":8080")
}
