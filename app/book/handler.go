package book

import (
	"context"
	"time"

	"github.com/DoNewsCode/core/srvhttp"
	"github.com/DoNewsCode/core/unierr"
	"github.com/DoNewsCode/skeleton/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Handler struct {
	Redis redis.UniversalClient
}

func (b Handler) Find(c *gin.Context) {
	encoder := srvhttp.NewResponseEncoder(c.Writer)
	book := c.Query("book")
	val, err := b.Redis.Get(c, book).Result()
	if err != nil {
		encoder.EncodeError(unierr.NotFoundErr(err, "book %s not found", book))
		return
	}
	encoder.EncodeResponse(entities.Book{
		BookName: val,
	})
}

func Seed(client redis.UniversalClient) error {
	_, err := client.Set(context.Background(), "book", "Thinking In Java", time.Hour).Result()
	return err
}
