package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v7/internal/usage"
)

// MyUsageHandler returns aggregated usage statistics for the authenticated API key.
func MyUsageHandler(c *gin.Context) {
	apiKey, _ := c.Get("userApiKey")
	key, ok := apiKey.(string)
	if !ok || key == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing API key"})
		return
	}

	stats := usage.GetRequestStatistics().GetAPIStats(key)
	if stats == nil {
		c.JSON(http.StatusOK, gin.H{
			"total_requests": 0,
			"success_count":  0,
			"failure_count":  0,
			"total_tokens":   0,
			"models":          map[string]interface{}{},
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
