package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
)

func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

// https://blog.canopas.com/golang-unit-tests-with-test-gin-context-80e1ac04adcd
// https://mayursinhsarvaiya.medium.com/how-to-mock-postgresql-database-for-unit-testing-in-golang-gorm-b690a4e4bc85
// https://medium.com/geekculture/easily-run-your-unit-test-with-golang-gin-postgres-8a402a29e3f6
// https://www.educative.io/answers/how-to-measure-test-coverage-in-go
