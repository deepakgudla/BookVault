package server

import (
	"net/http"

	"github.com/deepakgudla/bookvault/internal/config"
	"github.com/deepakgudla/bookvault/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server contains the dependencies and route configuration for the HTTP API.
type Server struct {
	config *config.Config
	// db             *gorm.DB
	logger         *zerolog.Logger
	authService    services.AuthServiceInterace
	productService services.ProductServiceInterface
	userService    services.UserServiceInterface
	uploadService  services.UploadServiceInterface
	cartService    services.CartServiceInterface
	orderService   services.OrderServiceInterface
}

// New creates an HTTP server with its service dependencies.
func New(cfg *config.Config,
	// db *gorm.DB,
	logger *zerolog.Logger,
	authService services.AuthServiceInterace,
	productService services.ProductServiceInterface,
	userService services.UserServiceInterface,
	uploadService services.UploadServiceInterface,
	cartService services.CartServiceInterface,
	orderService services.OrderServiceInterface,
) *Server {
	return &Server{
		config: cfg,
		// db:             db,
		logger:         logger,
		authService:    authService,
		productService: productService,
		userService:    userService,
		uploadService:  uploadService,
		cartService:    cartService,
		orderService:   orderService,
	}
}

// SetupRoutes configures and returns the HTTP router.
func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	// middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// routes
	router.GET("/health", s.HealthCheck)

	// doc routes
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.StaticFile("/api-docs", "./docs/rapidoc.html")

	router.Static("/uploads", "./uploads")

	router.GET("/playground", s.playgroundHandler())

	router.GET("/playground/public", s.playgroundPublicHandler())
	router.GET("/playground/protected", s.playgroundProtectedHandler())

	graphqlPublic := router.Group("/graphql/public")
	graphqlPublic.Use(s.graphqlPublicMiddleware())
	graphqlPublic.POST("/", s.graphqlHandler())

	graphqlProtected := router.Group("/graphql")
	graphqlProtected.Use(s.authMiddleware())
	graphqlProtected.Use(s.graphqlMiddleware())
	graphqlProtected.POST("/", s.graphqlHandler())

	api := router.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/register", s.register)
	auth.POST("/login", s.login)
	auth.POST("logout", s.logout)
	auth.POST("/refresh", s.refreshToken)

	protected := api.Group("/")
	protected.Use(s.authMiddleware())
	{
		users := protected.Group("/users")
		{
			userRoutes := users
			userRoutes.GET("/profile", s.getProfile)
			userRoutes.PUT("/profile", s.updateProfile)
		}

		categories := protected.Group("/categories")
		{
			categoryRoutes := categories
			categoryRoutes.POST("/", s.adminMiddleware(), s.createCategory)
			categoryRoutes.PUT("/:id", s.adminMiddleware(), s.updateCategory)
			categoryRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
		}

		products := protected.Group("/products")
		{
			productRoutes := products
			productRoutes.POST("/", s.adminMiddleware(), s.createProduct)
			productRoutes.PUT("/:id", s.adminMiddleware(), s.updateProduct)
			productRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
			productRoutes.POST("/:id/images", s.adminMiddleware(), s.uploadProductImage)
		}

		carts := protected.Group("/carts")
		{
			cartRoutes := carts
			cartRoutes.GET("/", s.getCart)
			cartRoutes.POST("/items", s.addToCart)
			cartRoutes.PUT("/items/:id", s.updateCartItem)
			cartRoutes.DELETE("/items/:id", s.removeFromCart)
		}

		orders := protected.Group("/orders")
		{
			orderRoutes := orders
			orderRoutes.POST("/orders", s.createOrder)
			orderRoutes.GET("/", s.getOrders)
			orderRoutes.GET("/:id", s.getOrder)
		}
	}

	api.GET("/categories", s.getCategories)
	api.GET("/products", s.getProducts)
	api.GET("/products/:id", s.getProduct)

	return router
}

// HealthCheck reports that the API is available.
func (s *Server) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization,")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
