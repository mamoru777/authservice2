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

func TestService_signIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRep := userrepository.NewMockIUserRepository(ctrl)
	sessionRep := sessionrepository.NewMockISessionRepository(ctrl)
	loggerRep := mylogger.NewMockILogger(ctrl)
	tokens := jwttokens.NewMockITokens(ctrl)
	logger := mylogger.New(&logrus.Logger{})
	testTable := []struct {
		name             string
		request          *gatewayapi.SignInRequest
		expectedErr      error
		expectedResponse *gatewayapi.SignInResponse
		x_request_id     string
	}{
		{
			name: "FullInfo",
			request: &gatewayapi.SignInRequest{
				Login:    "testlogin",
				Password: "111111",
			},
			x_request_id:     "testRequestId",
			expectedErr:      nil,
			expectedResponse: &gatewayapi.SignInResponse{AccessToken: "TestAccessToken", RefreshToken: "TestRefreshToken", IsExist: true, IsSignedUp: false},
		},
		{
			name: "WithoutLogin",
			request: &gatewayapi.SignInRequest{
				Login:    "",
				Password: "111111",
			},
			x_request_id:     "testRequestId",
			expectedErr:      errors.New("Логин должен состоять минимум из 4 символов"),
			expectedResponse: &gatewayapi.SignInResponse{},
		},
		{
			name: "WithoutPassword",
			request: &gatewayapi.SignInRequest{
				Login:    "testlogin",
				Password: "",
			},
			x_request_id:     "testRequestId",
			expectedErr:      errors.New("Пароль должен состоять минимум из 6 символов"),
			expectedResponse: &gatewayapi.SignInResponse{},
		},
	}
	for _, testCase := range testTable {
		convey.Convey(testCase.name, t, func() {
			userService := New(userRep, sessionRep, logger, tokens)
			ctx := context.WithValue(context.Background(), "x_request_id", testCase.x_request_id)
			_ = loggerRep.EXPECT().Init().Return(nil).AnyTimes()
			_ = userRep.EXPECT().GetByUserAndPassword(gomock.Any(), testCase.request.Login).Return(&userrepository.User{Password: []byte("$2a$10$eu9gKHpYI36Q7OVZKGCgzuWiimSOUwsVuoXgL5JpuYZIOGg0n2Qvi")}, nil).AnyTimes()
			_ = sessionRep.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&sessionrepository.Session{}, nil).AnyTimes()
			_ = sessionRep.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			_ = tokens.EXPECT().CreateRefreshToken(gomock.Any()).Return("TestRefreshToken", nil).AnyTimes()
			_ = tokens.EXPECT().CreateAccessToken(gomock.Any()).Return("TestAccessToken", nil).AnyTimes()
			//$2a$10$eu9gKHpYI36Q7OVZKGCgzuWiimSOUwsVuoXgL5JpuYZIOGg0n2Qvi
			response, err := userService.SignIn(ctx, testCase.request)
			convey.So(err, convey.ShouldEqual, testCase.expectedErr)
			convey.So(response, convey.ShouldResemble, testCase.expectedResponse)
		})
	}
}
