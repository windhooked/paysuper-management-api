package api

import (
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/paysuper/paysuper-management-api/internal/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type OrderTestSuite struct {
	suite.Suite
	router *orderRoute
	api    *Api
}

func Test_Order(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (suite *OrderTestSuite) SetupTest() {
	suite.api = &Api{
		Http:           echo.New(),
		validate:       validator.New(),
		billingService: mock.NewBillingServerOkMock(),
		authUser: &AuthUser{
			Id: "ffffffffffffffffffffffff",
		},
	}

	suite.api.authUserRouteGroup = suite.api.Http.Group(apiAuthUserGroupPath)
	suite.router = &orderRoute{Api: suite.api}

	err := suite.api.validate.RegisterValidation("uuid", suite.api.UuidValidator)
	assert.NoError(suite.T(), err, "Uuid validator registration failed")
}

func (suite *OrderTestSuite) TearDownTest() {}

func (suite *OrderTestSuite) TestOrder_GetRefund_Ok() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds/:refund_id")
	ctx.SetParamNames(requestParameterOrderId, requestParameterRefundId)
	ctx.SetParamValues(bson.NewObjectId().Hex(), bson.NewObjectId().Hex())

	err := suite.router.getRefund(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *OrderTestSuite) TestOrder_GetRefund_RefundIdEmpty_Error() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds/:refund_id")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(bson.NewObjectId().Hex())

	err := suite.router.getRefund(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), errorIncorrectRefundId, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_GetRefund_OrderIdEmpty_Error() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds/:refund_id")
	ctx.SetParamNames(requestParameterRefundId)
	ctx.SetParamValues(bson.NewObjectId().Hex())

	err := suite.router.getRefund(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), errorIncorrectOrderId, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_GetRefund_BillingServerError() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds/:refund_id")
	ctx.SetParamNames(requestParameterOrderId, requestParameterRefundId)
	ctx.SetParamValues(bson.NewObjectId().Hex(), bson.NewObjectId().Hex())

	suite.router.billingService = mock.NewBillingServerSystemErrorMock()

	err := suite.router.getRefund(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), errorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_GetRefund_BillingServer_RefundNotFound_Error() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds/:refund_id")
	ctx.SetParamNames(requestParameterOrderId, requestParameterRefundId)
	ctx.SetParamValues(bson.NewObjectId().Hex(), bson.NewObjectId().Hex())

	suite.router.billingService = mock.NewBillingServerErrorMock()

	err := suite.router.getRefund(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusNotFound, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ListRefunds_Ok() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(bson.NewObjectId().Hex())

	err := suite.router.listRefunds(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *OrderTestSuite) TestOrder_ListRefunds_BindError() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds")

	err := suite.router.listRefunds(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), errorIncorrectOrderId, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ListRefunds_BillingServerError() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(bson.NewObjectId().Hex())

	suite.router.billingService = mock.NewBillingServerSystemErrorMock()
	err := suite.router.listRefunds(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), errorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateRefund_Ok() {
	data := `{"amount": 10, "reason": "test"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(bson.NewObjectId().Hex())

	err := suite.router.createRefund(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *OrderTestSuite) TestOrder_CreateRefund_BindError() {
	data := `{"amount": 10, "reason": "test"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds")

	err := suite.router.createRefund(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), errorIncorrectOrderId, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateRefund_ValidationError() {
	data := `{"amount": -10, "reason": "test"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(bson.NewObjectId().Hex())

	err := suite.router.createRefund(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), "Amount", httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateRefund_BillingServerError() {
	data := `{"amount": 10, "reason": "test"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(bson.NewObjectId().Hex())

	suite.router.billingService = mock.NewBillingServerSystemErrorMock()
	err := suite.router.createRefund(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), errorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CreateRefund_BillingServer_CreateError() {
	data := `{"amount": 10, "reason": "test"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/order/:order_id/refunds")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(bson.NewObjectId().Hex())

	suite.router.billingService = mock.NewBillingServerErrorMock()
	err := suite.router.createRefund(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangeLanguage_Ok() {
	body := `{"lang": "en"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/language")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	err := suite.router.changeLanguage(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *OrderTestSuite) TestOrder_ChangeLanguage_OrderIdEmpty_Error() {
	body := `{"lang": "en"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	err := suite.router.changeLanguage(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), errorIncorrectOrderId, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangeLanguage_BindError() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/language")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	err := suite.router.changeLanguage(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), errorQueryParamsIncorrect, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangeLanguage_ValidationError() {
	body := `{"lang": "en"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/language")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues("some_value")

	err := suite.router.changeLanguage(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), "OrderId", httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangeLanguage_BillingServerSystemError() {
	body := `{"lang": "en"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/language")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	suite.router.billingService = mock.NewBillingServerSystemErrorMock()
	err := suite.router.changeLanguage(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), errorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangeLanguage_BillingServerErrorResult_Error() {
	body := `{"lang": "en"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/language")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	suite.router.billingService = mock.NewBillingServerErrorMock()
	err := suite.router.changeLanguage(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangePaymentAccount_Ok() {
	body := `{"method_id": "000000000000000000000000", "account": "4000000000000002"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/customer")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	err := suite.router.changeCustomer(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *OrderTestSuite) TestOrder_ChangePaymentAccount_OrderIdEmpty_Error() {
	body := `{"method_id": "000000000000000000000000", "account": "4000000000000002"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	err := suite.router.changeCustomer(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), errorIncorrectOrderId, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangePaymentAccount_BindError() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/customer")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	err := suite.router.changeCustomer(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), errorQueryParamsIncorrect, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangePaymentAccount_ValidationError() {
	body := `{"method_id": "some_value", "account": "4000000000000002"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/customer")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	err := suite.router.changeCustomer(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), "MethodId", httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangePaymentAccount_BillingServerSystemError() {
	body := `{"method_id": "000000000000000000000000", "account": "4000000000000002"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/customer")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	suite.router.billingService = mock.NewBillingServerSystemErrorMock()
	err := suite.router.changeCustomer(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), errorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_ChangePaymentAccount_BillingServerErrorResult_Error() {
	body := `{"method_id": "000000000000000000000000", "account": "4000000000000002"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/customer")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	suite.router.billingService = mock.NewBillingServerErrorMock()
	err := suite.router.changeCustomer(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_Ok() {
	body := `{"country": "US", "city": "Washington", "zip": "98001"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/billing_address")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	err := suite.router.processBillingAddress(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rsp.Code)
	assert.NotEmpty(suite.T(), rsp.Body.String())
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_OrderIdEmpty_Error() {
	body := `{"country": "US", "city": "Washington", "zip": "98001"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	err := suite.router.processBillingAddress(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), errorIncorrectOrderId, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_BindError() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/billing_address")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	err := suite.router.processBillingAddress(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), errorQueryParamsIncorrect, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_ValidationError() {
	body := `{"country": "some_value", "city": "Washington", "zip": "98001"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/billing_address")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	err := suite.router.processBillingAddress(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Regexp(suite.T(), "Country", httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_BillingServerSystemError() {
	body := `{"country": "US", "city": "Washington", "zip": "98001"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/billing_address")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	suite.router.billingService = mock.NewBillingServerSystemErrorMock()
	err := suite.router.processBillingAddress(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusInternalServerError, httpErr.Code)
	assert.Equal(suite.T(), errorUnknown, httpErr.Message)
}

func (suite *OrderTestSuite) TestOrder_CalculateAmounts_BillingServerErrorResult_Error() {
	body := `{"country": "US", "city": "Washington", "zip": "98001"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rsp := httptest.NewRecorder()
	ctx := e.NewContext(req, rsp)

	ctx.SetPath("/api/v1/orders/:order_id/billing_address")
	ctx.SetParamNames(requestParameterOrderId)
	ctx.SetParamValues(uuid.New().String())

	suite.router.billingService = mock.NewBillingServerErrorMock()
	err := suite.router.processBillingAddress(ctx)
	assert.Error(suite.T(), err)

	httpErr, ok := err.(*echo.HTTPError)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), http.StatusBadRequest, httpErr.Code)
	assert.Equal(suite.T(), mock.SomeError, httpErr.Message)
}