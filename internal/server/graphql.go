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
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

func (s *Server) GraphQLHandler() *handler.Server {

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

	serve.Use(extension.Introspection{})
	serve.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	return serve

}

func (s *Server) graphqlHandler() gin.HandlerFunc {
	h := s.GraphQLHandler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
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

func (s *Server) graphqlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, _ := c.Get("user_id")
		userEmail, _ := c.Get("user_email")
		userRole, _ := c.Get("user_role")

		ctx := context.WithValue(c.Request.Context(), "user_id", userID)
		ctx = context.WithValue(c.Request.Context(), "user_email", userEmail)
		ctx = context.WithValue(c.Request.Context(), "user_role", userRole)

		c.Request = c.Request.WithContext(ctx)

		c.Next()

	}
}
