package api

import (
	"net/http"

	"github.com/x-research-team/mattermost-html2md/pkg/models/request"
	"github.com/x-research-team/mattermost-html2md/pkg/models/response"

	"github.com/gin-gonic/gin"
)

func (a *api) Send(c *gin.Context, req request.Webhook) {
	err := a.service.SendAPI(c.Request.Context(), req.Body.Text, req.Body.Channel)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err))
		return
	}

	c.Status(http.StatusNoContent)
}
