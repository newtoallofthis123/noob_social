package routes

import "github.com/gin-gonic/gin"

func (api *ApiServer) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			c.Redirect(302, "/login")
			return
		}

		session, err := api.store.GetSessionById(cookie)
		if err != nil {
			c.Redirect(302, "/login")
			return
		}

		// Sets the user in the context so that it can be used anywhere.
		c.Set("user_id", session.UserId)

		c.Next()
	}
}
