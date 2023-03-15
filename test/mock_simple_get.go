package test

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func MockSimpleGet(c *gin.Context, params gin.Params, u url.Values) {
	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", 1)

	// set path params
	c.Params = params

	// set query params
	c.Request.URL.RawQuery = u.Encode()
}
