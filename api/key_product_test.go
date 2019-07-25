package api

import (
	"bytes"
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo/v4"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"net/http/httptest"
	"testing"
)

type KeyProductTestSuite struct {
	suite.Suite
	router *keyProductRoute
	api    *Api
}

func Test_keyProduct(t *testing.T) {
	suite.Run(t, new(KeyProductTestSuite))
}

func (suite *KeyProductTestSuite) SetupTest() {
	suite.api = &Api{
		Http:           echo.New(),
		validate:       validator.New(),
		billingService: mock.NewBillingServerOkMock(),
		authUser: &AuthUser{
			Id: "ffffffffffffffffffffffff",
		},
	}

	suite.api.authUserRouteGroup = suite.api.Http.Group(apiAuthUserGroupPath)
	suite.router = &keyProductRoute{Api: suite.api}
}

func (suite *KeyProductTestSuite) TearDownTest() {}

func (suite *KeyProductTestSuite) TestProject_GetPlatformList_ValidationError() {
	req := httptest.NewRequest(http.MethodGet, "/platforms", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := suite.api.Http.NewContext(req, rsp)

	ctx.SetPath("/platforms")

	err := suite.router.getPlatformsList(ctx)
	assert.Error(suite.T(), err)
	e, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), 400, e.Code)
	assert.NotEmpty(suite.T(), e.Message)
}

func (suite *KeyProductTestSuite) TestProject_RemovePlatform_Ok() {
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := suite.api.Http.NewContext(req, rsp)

	ctx.SetPath("/key-products/:key_product_id/platforms/:platform_id")
	ctx.SetParamNames("key_product_id", "platform_id")
	ctx.SetParamValues(bson.NewObjectId().Hex(), "steam")

	err := suite.router.removePlatformForKeyProduct(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_PublishKeyProduct_Ok() {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := suite.api.Http.NewContext(req, rsp)

	ctx.SetPath("/key-products/:key_product_id/publish")
	ctx.SetParamNames("key_product_id")
	ctx.SetParamValues(bson.NewObjectId().Hex())

	err := suite.router.publishKeyProduct(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_GetPlatformList_Ok() {
	req := httptest.NewRequest(http.MethodGet, "/platforms?limit=10", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := suite.api.Http.NewContext(req, rsp)

	ctx.SetPath("/platforms?limit=10")

	err := suite.router.getPlatformsList(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_GetListKeyProduct_Ok() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := suite.api.Http.NewContext(req, rsp)

	err := suite.router.getKeyProductList(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_GetKeyProduct_ValidationError() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := suite.api.Http.NewContext(req, rsp)

	ctx.SetPath("/key-products/:key_product_id")
	ctx.SetParamNames("key_product_id")
	ctx.SetParamValues("")

	err := suite.router.getKeyProductById(ctx)
	assert.Error(suite.T(), err)
	e, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), 400, e.Code)
	assert.NotEmpty(suite.T(), e.Message)
}

func (suite *KeyProductTestSuite) TestProject_GetKeyProduct_Ok() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := suite.api.Http.NewContext(req, rsp)

	ctx.SetPath("/key-products/:key_product_id")
	ctx.SetParamNames("key_product_id")
	ctx.SetParamValues(bson.NewObjectId().Hex())

	err := suite.router.getKeyProductById(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_CreateKeyProduct_Ok() {
	body := &grpc.CreateOrUpdateKeyProductRequest{
		MerchantId:      bson.NewObjectId().Hex(),
		Name:            map[string]string{"en": "A", "ru": "А"},
		Description:     map[string]string{"en": "A", "ru": "А"},
		DefaultCurrency: "RUB",
		ProjectId:       bson.NewObjectId().Hex(),
		Sku:             "some_sku",
		Enabled:         false,
	}

	b, err := json.Marshal(&body)
	assert.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := suite.api.Http.NewContext(req, rsp)

	err = suite.router.createKeyProduct(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_ChangeKeyProduct_ValidationError() {
	body := &grpc.CreateOrUpdateKeyProductRequest{
		Id:              bson.NewObjectId().Hex(),
		MerchantId:      bson.NewObjectId().Hex(),
		DefaultCurrency: "RUB",
		ProjectId:       bson.NewObjectId().Hex(),
		Sku:             "some_sku",
		Enabled:         false,
	}

	b, err := json.Marshal(&body)
	assert.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := suite.api.Http.NewContext(req, rsp)

	ctx.SetPath("/key-products/:key_product_id")
	ctx.SetParamNames("key_product_id")
	ctx.SetParamValues(body.Id)

	err = suite.router.changeKeyProduct(ctx)
	assert.Error(suite.T(), err)
	e, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), 400, e.Code)
	assert.NotEmpty(suite.T(), e.Message)
}

func (suite *KeyProductTestSuite) TestProject_ChangeKeyProduct_Ok() {
	body := &grpc.CreateOrUpdateKeyProductRequest{
		Id:              bson.NewObjectId().Hex(),
		MerchantId:      bson.NewObjectId().Hex(),
		Name:            map[string]string{"en": "A", "ru": "А"},
		Description:     map[string]string{"en": "A", "ru": "А"},
		DefaultCurrency: "RUB",
		ProjectId:       bson.NewObjectId().Hex(),
		Sku:             "some_sku",
		Enabled:         false,
	}

	b, err := json.Marshal(&body)
	assert.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := suite.api.Http.NewContext(req, rsp)

	ctx.SetPath("/key-products/:key_product_id")
	ctx.SetParamNames("key_product_id")
	ctx.SetParamValues(body.Id)

	err = suite.router.changeKeyProduct(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *KeyProductTestSuite) TestProject_CreateKeyProduct_ValidationError() {
	body := &grpc.CreateOrUpdateKeyProductRequest{
		MerchantId:      bson.NewObjectId().Hex(),
		Description:     map[string]string{"en": "A", "ru": "А"},
		DefaultCurrency: "RUB",
		ProjectId:       bson.NewObjectId().Hex(),
		Sku:             "some_sku",
		Enabled:         false,
	}

	b, err := json.Marshal(&body)
	assert.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := suite.api.Http.NewContext(req, rsp)

	err = suite.router.createKeyProduct(ctx)

	assert.Error(suite.T(), err)
	e, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), 400, e.Code)
	assert.NotEmpty(suite.T(), e.Message)
}