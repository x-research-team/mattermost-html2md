package api

import (
	"net/http"

	"github.com/x-research-team/mattermost-html2md/pkg/models/request"
	"github.com/x-research-team/mattermost-html2md/pkg/models/response"

	"github.com/gin-gonic/gin"
)

func (a *api) Send(c *gin.Context, req request.Send) {
	err := a.service.Send(c.Request.Context(), req.Body.Text)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err))
		return
	}

	c.Status(http.StatusNoContent)
}
