package shortener

import "github.com/gin-gonic/gin"

// Handlers register all the handlers for shortener related actions
type Handlers struct{}

// NewHandlers creates the handlers
func NewHandlers() (*Handlers, error) {
	h := new(Handlers)
	if err := h.init(); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Handlers) init() error {
	return nil
}

// Register the handlers
func (h *Handlers) Register(group *gin.RouterGroup) {}

// Redirect handle the url redirects
func (h *Handlers) Redirect(c *gin.Context) {

}
