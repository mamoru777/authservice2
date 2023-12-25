package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
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

func TestService_updateAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRep := userrepository.NewMockIUserRepository(ctrl)
	sessionRep := sessionrepository.NewMockISessionRepository(ctrl)
	loggerRep := mylogger.NewMockILogger(ctrl)
	logger := mylogger.New(&logrus.Logger{})
	tokens := jwttokens.NewMockITokens(ctrl)

	testTable := []struct {
		name             string
		request          *gatewayapi.UpdateAccessTokenRequest
		expectedErr      error
		expectedResponse *gatewayapi.UpdateAccessTokenResponse
		x_request_id     string
	}{
		{
			name: "FullInfo",
			request: &gatewayapi.UpdateAccessTokenRequest{
				Userid:       "a7add5d7-5c51-479f-96aa-c34a7df5a16d",
				RefreshToken: "testRefreshToken",
			},
			x_request_id:     "testRequestId",
			expectedErr:      nil,
			expectedResponse: &gatewayapi.UpdateAccessTokenResponse{AccessToken: "testAccessToken", RefreshToken: "testRefreshToken"},
		},
		{
			name: "WithoutUserId",
			request: &gatewayapi.UpdateAccessTokenRequest{
				Userid:       "",
				RefreshToken: "testRefreshToken",
			},
			x_request_id:     "testRequestId",
			expectedErr:      errors.New("userId в функции UpdateAccessToken имеет nil значение"),
			expectedResponse: &gatewayapi.UpdateAccessTokenResponse{},
		},
		{
			name: "WithoutRefreshToken",
			request: &gatewayapi.UpdateAccessTokenRequest{
				Userid:       "a7add5d7-5c51-479f-96aa-c34a7df5a16d",
				RefreshToken: "",
			},
			x_request_id:     "testRequestId",
			expectedErr:      errors.New("refreshToken в функции UpdateAccessToken имеет nil значение"),
			expectedResponse: &gatewayapi.UpdateAccessTokenResponse{},
		},
		{
			name: "WrongUserId",
			request: &gatewayapi.UpdateAccessTokenRequest{
				Userid:       "wrondId",
				RefreshToken: "testRefreshToken",
			},
			x_request_id:     "testRequestId",
			expectedErr:      errors.New("Не удалось конвертировать userId в uuid"),
			expectedResponse: &gatewayapi.UpdateAccessTokenResponse{},
		},
		{
			name: "WrongRefreshToken",
			request: &gatewayapi.UpdateAccessTokenRequest{
				Userid:       "a7add5d7-5c51-479f-96aa-c34a7df5a16d",
				RefreshToken: "wrongRefreshToken",
			},
			x_request_id:     "testRequestId",
			expectedErr:      errors.New("refreshToken из запроса в функции UpdateAccessToken не совпадает с refreshToken из бд"),
			expectedResponse: &gatewayapi.UpdateAccessTokenResponse{},
		},
	}

	for _, testCase := range testTable {
		convey.Convey(testCase.name, t, func() {
			userService := New(userRep, sessionRep, logger, tokens)
			ctx := context.WithValue(context.Background(), "x_request_id", testCase.x_request_id)
			_ = loggerRep.EXPECT().Init().Return(nil).AnyTimes()
			_ = sessionRep.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&sessionrepository.Session{RefreshToken: "testRefreshToken"}, nil).AnyTimes()
			_ = tokens.EXPECT().CreateAccessToken(gomock.Any()).Return("testAccessToken", nil).AnyTimes()
			_ = tokens.EXPECT().VerifyToken(gomock.Any(), gomock.Any()).Return(jwt.MapClaims{}, nil).AnyTimes()
			response, err := userService.UpdateAccessToken(ctx, testCase.request)
			convey.So(err, convey.ShouldEqual, testCase.expectedErr)
			convey.So(response, convey.ShouldResemble, testCase.expectedResponse)
		})
	}
}
