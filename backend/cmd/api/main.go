package main

import (
	"backend/internal/handlers"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/internal/websocket"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db := repository.NewPostgresDB()

	hub := websocket.NewHub()
	go hub.Run()

	// Создаём репозиторий один раз и переиспользуем
	userRepo := repository.NewUserRepository(db)

	// User слой
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Auth слой
	authService := service.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Post слой
	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService)

	// Relations слой
	relRepo := repository.NewRelationRepository(db)
	relService := service.NewRelationService(relRepo, postRepo)
	relHandler := handlers.NewRelationHandler(relService)

	//Chat слой
	chatRepo := repository.NewChatRepository(db)
	chatService := service.NewChatService(chatRepo, hub)
	chatHandler := handlers.NewChatHandler(chatService)

	// Interaction слой
	interRepo := repository.NewInteractionRepository(db)
	interService := service.NewInteractionService(interRepo)
	interHandler := handlers.NewInteractionHandler(interService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Публичные роуты
	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// Защищенные роуты (нужен JWT)
	api := r.Group("/api/v1")
	api.Use(handlers.AuthMiddleware())
	{
		api.POST("/posts", postHandler.CreatePost)
		api.GET("/feed", postHandler.GetFeed)
		api.POST("/posts/:id/like", interHandler.Like)
		api.DELETE("/posts/:id/like", interHandler.Unlike)
		api.POST("/posts/:id/comment", interHandler.Comment)
		api.GET("/posts/:id/comments", interHandler.GetComments)
		api.POST("/follow/:id", relHandler.Follow)
		api.GET("/my-feed", relHandler.GetPersonalFeed)
		api.GET("/users/search", userHandler.Search)
		api.GET("/chats/:id/messages", chatHandler.GetHistory)
		api.POST("/chats/private", chatHandler.CreatePrivateChat)
		api.Use(handlers.AuthMiddleware())
		{
			api.GET("/ws", func(c *gin.Context) {
				userIDStr, _ := c.Get("user_id")
				userID, _ := uuid.Parse(userIDStr.(string))

				// Передаем chatService в обработчик
				websocket.HandleWS(hub, chatService, c.Writer, c.Request, userID)
			})
		}
		r.Run(":" + os.Getenv("PORT"))
	}
}
