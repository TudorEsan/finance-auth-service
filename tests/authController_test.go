package tests

import (
	"auth-service/routes"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient() *mongo.Client {
	opts := &memongo.Options{
		MongoVersion: "6.0.0",
	}
	if runtime.GOARCH == "arm64" {
		if runtime.GOOS == "darwin" {
			// Only set the custom url as workaround for arm64 macs
			opts.DownloadURL = "https://fastdl.mongodb.org/osx/mongodb-macos-x86_64-5.0.0.tgz"
		}
	}
	server, err := memongo.StartWithOptions(opts)
	if err != nil {
		fmt.Println(err)
		panic("Error while starting the in memory mongodb")
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(server.URI()))
	if err != nil {
		panic("Error while creating client")
	}
	testCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = client.Connect(testCtx)
	if err != nil {
		panic("Error while connecting to the client")
	}
	return client
}

var inMemoryMongo = GetMongoClient()

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	gr := r.Group("")
	l := hclog.New(&hclog.LoggerOptions{
		Name: "TEST-AUTH",
	})
	routes.AuthRoutes(gr, l, inMemoryMongo)
	return r
}

func TestAuthService(t *testing.T) {
	t.Setenv("JWT_SECRET", "NoSecret")
	router := setupRouter()

	t.Run("Signup with email to short", func(t *testing.T) {
		w := httptest.NewRecorder()
		fmt.Println(1)
		body := []byte(`{
			"username": "tu",
			"password": "tudoresan",
			"email": "tudor@tud.com"
		}`)
		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("Signup with password too short", func(t *testing.T) {
		w := httptest.NewRecorder()
		fmt.Println(1)
		body := []byte(`{
			"username": "tudoresan",
			"password": "ad",
			"email": "tudo@rtud.com"
		}`)
		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("Signup with unvalid email", func(t *testing.T) {
		w := httptest.NewRecorder()
		fmt.Println(1)
		body := []byte(`{
			"username": "tudoresan",
			"password": "tudoresan",
			"email": "tudortud.com"
		}`)
		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
	})

	t.Run("Testing Signup when credentials are good", func(t *testing.T) {
		w := httptest.NewRecorder()
		fmt.Println(2)
		body := []byte(`{
			"username": "tudoresan",
			"password": "tudoresan",
			"email": "tudor@tud.com"
		}`)
		req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		if err != nil {
			t.Error("Error on new request", err)
		}
		router.ServeHTTP(w, req)
		t.Log(w.Body)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("Testing Login when credentials are good", func(t *testing.T) {
		w := httptest.NewRecorder()
		fmt.Println(3)
		body := []byte(`{
			"username": "tudoresan",
			"password": "tudoresan"
		}`)
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		if err != nil {
			t.Error("Error on new request", err)
		}
		router.ServeHTTP(w, req)
		t.Log(w)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("Signup with same email", func(t *testing.T) {
		w := httptest.NewRecorder()
		fmt.Println(1)
		body := []byte(`{
			"username": "tudoresan",
			"password": "tudoresan",
			"email": "tudor@tud.com"
		}`)
		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		router.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code)
	})
}

func TestLogin(t *testing.T) {

}
