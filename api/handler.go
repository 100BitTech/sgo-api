package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/samber/oops"
)

func JsonHandlerIO[T any](
	c *gin.Context,
	im T,
	handler func(im T) (any, error),
) {
	if err := c.ShouldBindBodyWith(&im, binding.JSON); err != nil {
		c.Error(oops.Wrap(err))
		return
	}

	if data, err := handler(im); err != nil {
		c.Error(oops.Wrap(err))
	} else {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data})
	}
}

func JsonHandlerI[T any](
	c *gin.Context,
	im T,
	handler func(im T) error,
) {
	if err := c.ShouldBindBodyWith(&im, binding.JSON); err != nil {
		c.Error(oops.Wrap(err))
		return
	}

	if err := handler(im); err != nil {
		c.Error(oops.Wrap(err))
	} else {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK})
	}
}

func JsonHandlerO(
	c *gin.Context,
	handler func() (any, error),
) {
	if data, err := handler(); err != nil {
		c.Error(oops.Wrap(err))
	} else {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data})
	}
}

func JsonContextHandlerIO[T any](
	c *gin.Context,
	im T,
	handler func(ctx context.Context, im T) (any, error),
) {
	if err := c.ShouldBindBodyWith(&im, binding.JSON); err != nil {
		c.Error(oops.Wrap(err))
		return
	}

	if data, err := handler(c.Request.Context(), im); err != nil {
		c.Error(oops.Wrap(err))
	} else {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data})
	}
}

func JsonContextHandlerI[T any](
	c *gin.Context,
	im T,
	handler func(ctx context.Context, im T) error,
) {
	if err := c.ShouldBindBodyWith(&im, binding.JSON); err != nil {
		c.Error(oops.Wrap(err))
		return
	}

	if err := handler(c.Request.Context(), im); err != nil {
		c.Error(oops.Wrap(err))
	} else {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK})
	}
}

func JsonContextHandlerO(
	c *gin.Context,
	handler func(ctx context.Context) (any, error),
) {
	if data, err := handler(c.Request.Context()); err != nil {
		c.Error(oops.Wrap(err))
	} else {
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data})
	}
}
