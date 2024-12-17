package service_test

import (
	"context"
	"testing"

	"github.com/marvinmarpol/golang-boilerplate/internal/pkg/utils/cryptho"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/command"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/domain/mask"
	mock_mask "github.com/marvinmarpol/golang-boilerplate/internal/protection/domain/mask/mock"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/query"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/service"
	"github.com/sirupsen/logrus"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx            context.Context
	mockCtrl       *gomock.Controller
	mockRepository *mock_mask.MockRepository
	commands       command.Commands
	queries        query.Queries
	svc            service.Services
}

// SetupSuite runs once before any tests are executed.
func (suite *ServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	suite.mockCtrl = gomock.NewController(suite.T()) // Initialize gomock
	suite.mockRepository = mock_mask.NewMockRepository(suite.mockCtrl)

	// load public and private key for KEK
	publicKey, err := cryptho.LoadRSAPublicKeyFromFile("../../../infrastructures/kek/public_key.pem")
	if err != nil {
		logrus.Error(err)
		return
	}
	privateKey, err := cryptho.LoadRSAPrivateKeyFromFile("../../../infrastructures/kek/private_key.pem")
	if err != nil {
		logrus.Error(err)
		return
	}

	createMaskHandler := command.NewCreateMaskHandler(suite.mockRepository)
	updateTokenHandler := command.NewUpdateTokenHandler(suite.mockRepository)
	suite.commands = command.Commands{
		CreateMaskCommand:  createMaskHandler,
		UpdateTokenCommand: updateTokenHandler,
	}

	getCypherHandler := query.NewGetCypherHandler(suite.mockRepository)
	getMaskHandler := query.NewGetMaskHandler(suite.mockRepository)
	getTokenHandler := query.NewGetTokenHandler(suite.mockRepository)
	suite.queries = query.Queries{
		GetCypherQuery: getCypherHandler,
		GetMaskQuery:   getMaskHandler,
		GetTokenQuery:  getTokenHandler,
	}

	// Initialize the service or any required dependencies here
	suite.svc = service.NewServiceImpl(suite.commands, suite.queries, publicKey, privateKey)
}

// TearDownSuite runs once after all tests are executed.
func (suite *ServiceTestSuite) TearDownSuite() {
	suite.mockCtrl.Finish() // Ensure all expectations are met
}

// SetupTest runs before each individual test.
func (suite *ServiceTestSuite) SetupTest() {
	// Reset the state if needed before each test
}

func (suite *ServiceTestSuite) TestDeidentify() {
	// Example input data
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    12345,
			"name":  "encrypt-JohnDoe",
			"email": "encrypt-johndoe@example.com",
		},
	}

	decryptPrefix := "decrypt-"

	// Setup mock expectations for the Commands.CreateMaskCommand.Handle method
	suite.mockRepository.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).Times(2) // Expect 2 calls for `name` and `email` fields

	// Call the method being tested
	suite.svc.Deidentify(suite.ctx, data)

	// Assert that the values in the map are correctly encrypted (mocking the encryption behavior)
	userData := data["user"].(map[string]interface{})
	assert.Contains(suite.T(), userData["name"], decryptPrefix)
	assert.Contains(suite.T(), userData["email"], decryptPrefix)
}

func (suite *ServiceTestSuite) TestReidentify() {
	// Example input data
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"email": "decrypt-ebi?iG)vhOAO{@LBm>i",
			"id":    12345,
			"name":  "decrypt-tVp_%uQ",
		},
	}

	decryptPrefix := "decrypt-"

	// Setup mock expectations for the Commands.CreateMaskCommand.Handle method
	suite.mockRepository.
		EXPECT().
		FindByToken(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, token string) (mask.Mask, error) {
			return mask.Mask{
				Cypher: "4cLkwq_bLxSOIax2vFib-R2R4ofhLiPHiABtk6VMIe32ksX1Cn1u7YhDHAAme0M=",
				Key:    "D1yDFcSHen41anZaHLfRWNGoD33Ko1N6lhl/kGzqx8iDJRqzKQd3+u3K8a8LmzvSeG2itlmgmcQP4XH44CYDHWL06raMwIJ+mb90z3JznPWV/X1heOwj5EimvNLo9ar30aAVUc0d9+uY6h5vu3buJEHbxBq6ThWDIIdRMYyvKUg=",
			}, nil // Return empty mask on the first call
		}).
		Times(2)

	// Call the method being tested
	suite.svc.Reidentify(suite.ctx, data)

	// Assert that the values in the map are correctly encrypted (mocking the encryption behavior)
	userData := data["user"].(map[string]interface{})
	assert.NotContains(suite.T(), userData["name"], decryptPrefix)
	assert.NotContains(suite.T(), userData["email"], decryptPrefix)
	assert.Equal(suite.T(), "johndoe@example.com", userData["name"])
}

// TestServiceTestSuite runs the test suite
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
