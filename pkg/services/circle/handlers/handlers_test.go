package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Shopify/sarama"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spatiumsocialis/infra/configs/services/circle/config"
	"github.com/spatiumsocialis/infra/pkg/common"
	"github.com/spatiumsocialis/infra/pkg/common/auth"
	"github.com/spatiumsocialis/infra/pkg/common/kafka"
	"github.com/spatiumsocialis/infra/pkg/services/circle/models"
	"github.com/stretchr/testify/assert"
)

var s *common.Service

var saramaConfig *sarama.Config
var producer sarama.AsyncProducer
var validUsers = []auth.User{
	{ID: "1"},
	{ID: "2"},
	{ID: "3"},
}

var validCircle = models.Circle{ID: "123", Users: validUsers}
var testUID = "TEST_UID"
var testToken = &auth.Token{UID: testUID}

func TestMain(m *testing.M) {
	if err := common.LoadEnv(); err != nil {
		log.Fatalln(err)
	}
	os.Setenv("DB_PROVIDER", "sqlite3")
	os.Setenv("DB_CONNECTION_STRING", ":memory:")
	saramaConfig = sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	producer = kafka.NewNullAsyncProducer()
	os.Exit(m.Run())
}

func TestGetCircle_Valid(t *testing.T) {
	s = common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, producer, config.ProductionTopic)
	err := s.DB.Create(&validCircle).Error
	assert.Nil(t, err)
	token := &auth.Token{UID: "1"}
	r := httptest.NewRequest("GET", "/irrelevant", nil)
	ctx := auth.AddTokenTo(context.Background(), token)
	w := httptest.NewRecorder()
	GetCircle(s).ServeHTTP(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
	var responseCircle models.CircleResponse
	err = json.Unmarshal(w.Body.Bytes(), &responseCircle)
	assert.Nil(t, err)
	assert.Equal(t, validCircle.ID, responseCircle.ID)
	assert.Equal(t, len(validCircle.Users), len(responseCircle.Users))

}

func TestGetCircle_NoCircle(t *testing.T) {
	s = common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, producer, config.ProductionTopic)
	err := s.DB.Create(&validCircle).Error
	assert.Nil(t, err)
	token := &auth.Token{UID: "99"}
	r := httptest.NewRequest("GET", "/irrelevant", nil)
	ctx := auth.AddTokenTo(context.Background(), token)
	w := httptest.NewRecorder()
	GetCircle(s).ServeHTTP(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
	var responseCircle models.CircleResponse
	err = json.Unmarshal(w.Body.Bytes(), &responseCircle)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(responseCircle.Users))
	assert.Equal(t, token.UID, responseCircle.Users[0].UID)
}

func TestAddToCircle_Valid(t *testing.T) {
	s = common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, producer, config.ProductionTopic)
	err := s.DB.Create(&validCircle).Error
	assert.Nil(t, err)
	token := &auth.Token{UID: "4"}
	payload, err := json.Marshal(map[string]string{"id": "123"})
	assert.Nil(t, err)
	r := httptest.NewRequest("PATCH", "/irrelevant", bytes.NewBuffer(payload))
	ctx := auth.AddTokenTo(context.Background(), token)
	w := httptest.NewRecorder()
	AddToCircle(s).ServeHTTP(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)
	var responseCircle models.CircleResponse
	err = json.Unmarshal(w.Body.Bytes(), &responseCircle)
	assert.Equal(t, 4, len(responseCircle.Users))
}

func TestAddToCircle_Invalid(t *testing.T) {
	s = common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, producer, config.ProductionTopic)
	err := s.DB.Create(&validCircle).Error
	assert.Nil(t, err)
	token := &auth.Token{UID: "4"}
	payload, err := json.Marshal(map[string]string{"id": "9999"})
	assert.Nil(t, err)
	r := httptest.NewRequest("PATCH", "/irrelevant", bytes.NewBuffer(payload))
	ctx := auth.AddTokenTo(context.Background(), token)
	w := httptest.NewRecorder()
	AddToCircle(s).ServeHTTP(w, r.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
