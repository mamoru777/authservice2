package service

import (
	"context"
	"errors"
	"github.com/mamoru777/authservice2/internal/jwttokens"
	"github.com/mamoru777/authservice2/internal/mylogger"
	"github.com/mamoru777/authservice2/internal/repositories/sessionrepository"
	"github.com/mamoru777/authservice2/internal/repositories/userrepository"
	"github.com/sirupsen/logrus"
	"github.com/smartystreets/goconvey/convey"

	gatewayapi "github.com/mamoru777/authservice2/pkg/gateway-api"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestService_signUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRep := userrepository.NewMockIUserRepository(ctrl)
	sessionRep := sessionrepository.NewMockISessionRepository(ctrl)
	loggerRep := mylogger.NewMockILogger(ctrl)
	logger := mylogger.New(&logrus.Logger{})
	tokens := jwttokens.NewMockITokens(ctrl)

	testTable := []struct {
		name             string
		request          *gatewayapi.SignUpRequest
		expectedErr      error
		expectedResponse *gatewayapi.SignUpResponse
		x_request_id     string
	}{
		{
			name: "FullInfo",
			request: &gatewayapi.SignUpRequest{
				Login:    "testlogin",
				Password: "testpassword",
				Email:    "test@example.com",
			},
			x_request_id:     "testRequestId",
			expectedErr:      nil,
			expectedResponse: &gatewayapi.SignUpResponse{},
		},
		{
			name: "WithoutLogin",
			request: &gatewayapi.SignUpRequest{
				Login:    "",
				Password: "testpassword",
				Email:    "test@example.com",
			},
			x_request_id:     "testRequestId",
			expectedErr:      errors.New("Логин должен состоять минимум из 4 символов"),
			expectedResponse: &gatewayapi.SignUpResponse{},
		},
		{
			name: "WithoutPassword",
			request: &gatewayapi.SignUpRequest{
				Login:    "testlogin",
				Password: "",
				Email:    "test@example.com",
			},
			x_request_id:     "testRequestId",
			expectedErr:      errors.New("Пароль должен состоять минимум из 6 символов"),
			expectedResponse: &gatewayapi.SignUpResponse{},
		},
		{
			name: "WithoutEmail",
			request: &gatewayapi.SignUpRequest{
				Login:    "testlogin",
				Password: "testpassword",
				Email:    "",
			},
			x_request_id:     "testRequestId",
			expectedErr:      errors.New("Поле почта обязательна для заполнения"),
			expectedResponse: &gatewayapi.SignUpResponse{},
		},
	}

	for _, testCase := range testTable {
		convey.Convey(testCase.name, t, func() {
			userService := New(userRep, sessionRep, logger, tokens)
			ctx := context.WithValue(context.Background(), "x_request_id", testCase.x_request_id)
			_ = userRep.EXPECT().GetByEmail(gomock.Any(), testCase.request.Email).Return(&userrepository.User{}, nil).AnyTimes()
			_ = loggerRep.EXPECT().Init().Return(nil).AnyTimes()
			_ = userRep.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			_ = sessionRep.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			_ = tokens.EXPECT().CreateRefreshToken(gomock.Any()).Return("TestRefreshToken", nil).AnyTimes()
			response, err := userService.SignUp(ctx, testCase.request)
			convey.So(err, convey.ShouldEqual, testCase.expectedErr)
			convey.So(response, convey.ShouldResemble, testCase.expectedResponse)
		})
	}

}
