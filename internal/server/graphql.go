package server

import (
	"context"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/deepakgudla/bookvault/graph"
	"github.com/deepakgudla/bookvault/graph/resolver"
	"github.com/deepakgudla/bookvault/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

// GraphQLRequest describes a GraphQL HTTP request.
type GraphQLRequest struct {
	Query         string                 `json:"query" example:"query { products { edges { node { id name } } } }"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse is the standard GraphQL HTTP response envelope.
type GraphQLResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors interface{} `json:"errors,omitempty"`
}

// GraphQLHandler builds the executable GraphQL HTTP handler.
func (s *Server) GraphQLHandler() *handler.Server {
	if s.graphqlServer != nil {
		return s.graphqlServer
	}

	r := resolver.NewResolver(
		s.authService,
		s.userService,
		s.productService,
		s.cartService,
		s.orderService,
	)

	schema := graph.NewExecutableSchema(graph.Config{Resolvers: r})

	serve := handler.New(schema)

	serve.AddTransport(transport.Options{})
	serve.AddTransport(transport.GET{})
	serve.AddTransport(transport.POST{})
	serve.AddTransport(transport.MultipartForm{})

	serve.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	if s.config.Environment != "production" {
		serve.Use(extension.Introspection{})
	}
	serve.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	s.graphqlServer = serve
	return s.graphqlServer

}

// @Summary Public GraphQL endpoint
// @Description Execute public GraphQL queries. The GraphQL schema defines available queries and mutations.
// @Tags GraphQL
// @Accept json
// @Produce json
// @Param request body GraphQLRequest true "GraphQL request"
// @Success 200 {object} GraphQLResponse
// @Failure 400 {object} GraphQLResponse
// @Router /graphql/public [post]
func (s *Server) graphqlPublicEndpoint(c *gin.Context) {
	s.GraphQLHandler().ServeHTTP(c.Writer, c.Request)
}

// @Summary Protected GraphQL endpoint
// @Description Execute authenticated GraphQL queries and mutations. The GraphQL schema defines available operations.
// @Tags GraphQL
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GraphQLRequest true "GraphQL request"
// @Success 200 {object} GraphQLResponse
// @Failure 400 {object} GraphQLResponse
// @Failure 401 {object} GraphQLResponse
// @Router /graphql [post]
func (s *Server) graphqlProtectedEndpoint(c *gin.Context) {
	s.GraphQLHandler().ServeHTTP(c.Writer, c.Request)
}

func (s *Server) playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("Graphql playground", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) playgroundPublicHandler() gin.HandlerFunc {
	h := playground.Handler("Graphql public playground", "/graphql/public")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) playgroundProtectedHandler() gin.HandlerFunc {
	h := playground.Handler("Graphql protected playground", "/graphql/")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// graphqlMiddleware is for the PROTECTED endpoint — auth is guaranteed by authMiddleware() beforehand.
func (s *Server) graphqlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "user_id not found in Gin context"})
			return
		}
		userEmail, _ := c.Get("user_email")
		userRole, _ := c.Get("user_role")

		ctx := context.WithValue(c.Request.Context(), utils.UserIDKey, userID)
		ctx = context.WithValue(ctx, utils.UserEmailKey, userEmail)
		ctx = context.WithValue(ctx, utils.UserRoleKey, userRole)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// graphqlPublicMiddleware is for the PUBLIC endpoint — auth is optional.
// Resolvers that require a user (Me, Cart, Orders, mutations, etc.) will
// still fail via GetUserIDFromContext, but public fields like products/categories work fine.
func (s *Server) graphqlPublicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if userID, exists := c.Get("user_id"); exists {
			ctx = context.WithValue(ctx, utils.UserIDKey, userID)
		}
		if userEmail, exists := c.Get("user_email"); exists {
			ctx = context.WithValue(ctx, utils.UserEmailKey, userEmail)
		}
		if userRole, exists := c.Get("user_role"); exists {
			ctx = context.WithValue(ctx, utils.UserRoleKey, userRole)
		}
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
