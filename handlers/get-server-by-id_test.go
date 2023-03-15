package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/minhmannh2001/sms/entity"
	"github.com/minhmannh2001/sms/test"
	"github.com/stretchr/testify/assert"
)

var dbName, mock_dbName string

// test
func TestGetUserId(t *testing.T) {

	// Load env variables
	err := godotenv.Load()
	if err != nil {
		godotenv.Load("./../.env")
	}

	mock_dbName = os.Getenv("MOCK_DBNAME")
	dbName = os.Getenv("DBNAME")

	// Create mock database
	test.MockedDB(test.CREATE)

	// Replace current connection with the mock database connection
	serverDatabase.ReplaceConnection(mock_dbName)

	// Do our test

	var server = entity.Server{
		Name:     "server1",
		Ipv4:     "172.22.0.11",
		User:     "root",
		Password: "helloworld",
		Status:   "Up",
	}

	serverController.CreateServer(&server)

	w := httptest.NewRecorder()
	ctx := test.GetTestGinContext(w)

	// Configure path params
	params := []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}

	// Configure query params
	u := url.Values{}

	test.MockSimpleGet(ctx, params, u)

	GetServerById(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)

	// Drop mock database
	defer test.MockedDB(test.DROP)

	// Replace mock database connection with the real one
	serverDatabase.ReplaceConnection(dbName)
}

func TestGetUserNonexistedId(t *testing.T) {

	// Load env variables
	err := godotenv.Load()
	if err != nil {
		godotenv.Load("./../.env")
	}

	mock_dbName = os.Getenv("MOCK_DBNAME")
	dbName = os.Getenv("DBNAME")

	// Create mock database
	test.MockedDB(test.CREATE)

	// Replace current connection with the mock database connection
	serverDatabase.ReplaceConnection(mock_dbName)

	// Do our test

	var server = entity.Server{
		Name:     "server1",
		Ipv4:     "172.22.0.11",
		User:     "root",
		Password: "helloworld",
		Status:   "Up",
	}

	serverController.CreateServer(&server)

	w := httptest.NewRecorder()
	ctx := test.GetTestGinContext(w)

	//configure path params
	params := []gin.Param{
		{
			Key:   "id",
			Value: "-1",
		},
	}

	// configure query params
	u := url.Values{}

	test.MockSimpleGet(ctx, params, u)

	GetServerById(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)

	// Drop mock database
	defer test.MockedDB(test.DROP)

	// Replace mock database connection with the real one
	serverDatabase.ReplaceConnection(dbName)
}
