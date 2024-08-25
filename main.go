package main

import (
	"net/http"
	"time"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	libs "spajam/libs"
)

type ChatMessage struct {
    ID         int64     `json:"id"`
    MessageText string   `json:"message_text"`
    SenderID   int       `json:"sender_id"`
    IsHome     bool      `json:"is_home"`
    CreatedAt  time.Time `json:"created_at"`
}

func main() {
	godotenv.Load()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(30)))
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 1,
	}))
	e.Use(middleware.RemoveTrailingSlash())

	v1 := e.Group("/v1")
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	v1.GET("/hello", Hello)
	v1.GET("/getchats", GetChats)

	serverPort := ":3000"
	e.Logger.Fatal(e.Start(serverPort))
}

func Hello(c echo.Context) error {
	return c.JSON(http.StatusOK, libs.ErrorResponse("Hello, World!"))
}

func GetChats(c echo.Context) error {
	dsn := os.Getenv("DSN")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var msgs []ChatMessage
	db.Find(&msgs)
	return c.JSON(http.StatusOK, msgs)
}

/*package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	"spajam/chat"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading env target")
	}

	// コマンドライン引数から質問テキストを取得する
	if len(os.Args) < 2 {
		panic("too few arguments")
	}
	content := os.Args[1]

	secret := os.Getenv("OPENAI_API_KEY")

	// リソース節約のためにタイムアウトを設定する
	timeout := 15 * time.Second

	// トークン節約のために応答の最大トークンを設定する
	maxTokens := 500

	// チャットに使用するモデルのID
	modelID := "gpt-4o-mini"

	c := chat.NewChatCompletions(modelID, secret, maxTokens, timeout)
	res, err := c.AskOneQuestion(content)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("In %d / Out %d / Total %d tokens\n", res.Usage.PromptTokens, res.Usage.CompletionTokens, res.Usage.TotalTokens)
	for _, v := range res.Choices {
		fmt.Printf("[%s]: %s\n", v.Message.Role, v.Message.Content)
	}

}
*/
