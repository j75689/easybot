package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// NoMethodHandler handle method not found
func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method Not Allowed"})
	}
}

// NoRouteHandler handle route path not found
func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page NotFound"})
	}
}

// RouteSkipperFunc check route path func
type RouteSkipperFunc func(c *gin.Context) bool

// Prefix struct
type Prefix struct{ Method, Path string }

// AllowPathPrefixSkipper if path prefix skip
func AllowPathPrefixSkipper(prefixes ...Prefix) RouteSkipperFunc {
	return func(c *gin.Context) bool {
		path := strings.ToUpper(c.Request.URL.Path)
		method := strings.ToUpper(c.Request.Method)
		for _, p := range prefixes {
			prefixPath := strings.ToUpper(p.Path)
			prefixMethod := strings.ToUpper(p.Method)
			if prefixMethod == "ANY" {
				if strings.HasPrefix(path, prefixPath) {
					return true
				}
			} else {
				if prefixMethod == method && strings.HasPrefix(path, prefixPath) {
					return true
				}
			}
		}
		return false
	}
}

// NoAllowPathPrefixSkipper if path prefix not skip
func NoAllowPathPrefixSkipper(prefixes ...struct{ Method, Path string }) RouteSkipperFunc {
	return func(c *gin.Context) bool {
		path := strings.ToUpper(c.Request.URL.Path)
		method := strings.ToUpper(c.Request.Method)
		for _, p := range prefixes {
			prefixPath := strings.ToUpper(p.Path)
			prefixMethod := strings.ToUpper(p.Method)
			if prefixMethod == "ANY" {
				if strings.HasPrefix(path, prefixPath) {
					return false
				}
			} else {
				if prefixMethod == method && strings.HasPrefix(path, prefixPath) {
					return false
				}
			}
		}
		return true
	}
}
