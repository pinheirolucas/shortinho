package shortener

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type responseError struct {
	Kind    string `json:"kind,omitempty"`
	Message string `json:"message,omitempty"`
}

type response struct {
	Error *responseError `json:"error,omitempty"`
	Data  interface{}    `json:"data,omitempty"`
}

// Handlers register all the handlers for shortener related actions
type Handlers struct {
	service  Service
	slugSize int
}

// NewHandlers creates the handlers
func NewHandlers() (*Handlers, error) {
	h := new(Handlers)
	if err := h.init(); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Handlers) init() error {
	s, err := NewService()
	if err != nil {
		return err
	}

	slugSize := viper.GetInt("app.slug-size")
	if slugSize < 4 {
		return errors.New("minimum slug size is 4")
	}

	h.slugSize = slugSize

	h.service = s
	return nil
}

// Register the handlers
func (h *Handlers) Register(root *gin.RouterGroup) {
	shortener := root.Group("/shortener")
	shortener.POST("/", h.Create)
	shortener.PATCH("/activate/:slug", h.Activate)
	shortener.PATCH("/deactivate/:slug", h.Deactivate)
	shortener.PATCH("/reset/:slug", h.Reset)

	shortener.GET("/:slug", h.Stats)

	root.GET("/r/:slug/", h.Redirect)
}

// Redirect handle the url redirects
func (h *Handlers) Redirect(c *gin.Context) {
	slug := c.Param("slug")

	link, err := h.service.Get(slug)
	switch err {
	case nil:
		// continue
	case ErrNotFound:
		c.Status(http.StatusNotFound)
		return
	default:
		log.Error().
			Err(err).
			Msg("unable to get link info")

		c.Status(http.StatusInternalServerError)
		return
	}

	if !link.Active || link.MaxHits != 0 && link.MaxHits == link.Hits {
		c.Status(http.StatusNotFound)
		return
	}

	err = h.service.Hit(slug)
	if err != nil {
		log.Error().
			Err(err).
			Msg("unable hit the link")

		c.Status(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusPermanentRedirect, link.TargetURL)
}

// Create handle the link creation
func (h *Handlers) Create(c *gin.Context) {
	link := new(Link)
	res := new(response)

	if err := c.BindJSON(link); err != nil {
		res.Error = &responseError{
			Kind:    "body_read",
			Message: "unable to parse JSON body",
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	link.Slug = strings.TrimSpace(link.Slug)
	if link.Slug == "" {
		slug, err := gonanoid.Nanoid(h.slugSize)
		if err != nil {
			log.Error().
				Err(err).
				Msg("unable to generate slug")

			res.Error = &responseError{
				Kind:    "slug_generate",
				Message: err.Error(),
			}
			c.JSON(http.StatusInternalServerError, slug)
			return
		}
		link.Slug = slug
	}

	if err := validation.Validate(link.TargetURL, validation.Required, is.URL); err != nil {
		res.Error = &responseError{
			Kind:    "invalid_target_url",
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	created, err := h.service.New(link)
	switch err {
	case nil:
		// continue
	case ErrAlreadyExists:
		res.Error = &responseError{
			Kind:    "already_exists",
			Message: err.Error(),
		}
		c.JSON(http.StatusConflict, res)
		return
	default:
		res.Error = &responseError{
			Kind:    "unknown",
			Message: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.JSON(http.StatusCreated, created)
}

// Activate handle link activation
func (h *Handlers) Activate(c *gin.Context) {
	slug := c.Param("slug")
	res := new(response)

	err := h.service.Activate(slug)
	switch err {
	case nil:
		// continue
	case ErrNotFound:
		c.Status(http.StatusNotFound)
		return
	default:
		log.Error().
			Err(err).
			Msg("unable to activate")

		res.Error = &responseError{
			Kind:    "unknown",
			Message: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.Status(http.StatusOK)
}

// Deactivate handle link deactivation
func (h *Handlers) Deactivate(c *gin.Context) {
	slug := c.Param("slug")
	res := new(response)

	err := h.service.Deactivate(slug)
	switch err {
	case nil:
		// continue
	case ErrNotFound:
		c.Status(http.StatusNotFound)
		return
	default:
		log.Error().
			Err(err).
			Msg("unable to deactivate")

		res.Error = &responseError{
			Kind:    "unknown",
			Message: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.Status(http.StatusOK)
}

// Reset handle link stats reset
func (h *Handlers) Reset(c *gin.Context) {
	slug := c.Param("slug")
	res := new(response)

	err := h.service.Reset(slug)
	switch err {
	case nil:
		// continue
	case ErrNotFound:
		c.Status(http.StatusNotFound)
		return
	default:
		log.Error().
			Err(err).
			Msg("unable to reset")

		res.Error = &responseError{
			Kind:    "unknown",
			Message: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.Status(http.StatusOK)
}

// Stats of the link
func (h *Handlers) Stats(c *gin.Context) {
	slug := c.Param("slug")
	res := new(response)

	link, err := h.service.Get(slug)
	switch err {
	case nil:
		// continue
	case ErrNotFound:
		c.Status(http.StatusNotFound)
		return
	default:
		log.Error().
			Err(err).
			Msg("unable to get stats")

		res.Error = &responseError{
			Kind:    "unknown",
			Message: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	fmt.Println(link)

	c.JSON(http.StatusOK, link)
}
