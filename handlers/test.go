package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/ivybridge-3c33/golang-elastic-apm/bootstrap"
	"github.com/ivybridge-3c33/golang-elastic-apm/entities"
)

type TestHandler struct {
	mysql bootstrap.MySQL
	redis bootstrap.RedisDB
}

func (h *TestHandler) Home(c *gin.Context) {
	ct := c.Request.Context()
	models := []*entities.Test{}
	index := "TEST_INDEX"
	if cache, err := h.redis.DB().Get(ct, index).Result(); err != nil {
		if err == redis.Nil {
			db := h.mysql.DB().WithContext(ct)
			if err := db.Find(&models).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			t, _ := time.ParseDuration("1m")
			j, _ := json.Marshal(models)

			if err := h.redis.DB().Set(ct, index, j, t).Err(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		if err := json.Unmarshal([]byte(cache), &models); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": models,
	})
}
