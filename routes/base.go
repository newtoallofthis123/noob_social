package routes

import "github.com/gin-gonic/gin"

func (api *ApiServer) handleAuthCheck(c *gin.Context) {
	sessionId, err := c.Cookie("session_id")
	if err != nil {
		c.String(500, err.Error())
		return
	}

	_, err = api.store.GetSessionById(sessionId)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.String(200, "OK")
}
