package book

import (
	"context"
	"time"

	"github.com/DoNewsCode/skeleton/internal/entities"
	"github.com/DoNewsCode/std/pkg/srverr"
	"github.com/DoNewsCode/std/pkg/srvhttp"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Handler struct {
	Redis redis.UniversalClient
}

func (b Handler) Find(c *gin.Context) {
	book := c.Query("book")
	val, err := b.Redis.Get(c, book).Result()
	if err != nil {
		srvhttp.EncodeError(c.Writer, srverr.NotFoundErr(err, "book %s not found", book))
		return
	}
	srvhttp.EncodeResponse(c.Writer, entities.Book{
		BookName: val,
	})
}

func Seed(client redis.UniversalClient) error {
	_, err := client.Set(context.Background(), "book", "Thinking In Java", time.Hour).Result()
	return err
}
