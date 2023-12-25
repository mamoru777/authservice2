package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mamoru777/authservice2/internal/jwttokens"
	"github.com/mamoru777/authservice2/internal/mylogger"
	"github.com/mamoru777/authservice2/internal/repositories/sessionrepository"
	"github.com/mamoru777/authservice2/internal/repositories/userrepository"

	gatewayapi "github.com/mamoru777/authservice2/pkg/gateway-api"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type UserService struct {
	gatewayapi.UnimplementedUsrServiceServer
	userrep    userrepository.IUserRepository
	sessionrep sessionrepository.ISessionRepository
	logger     *mylogger.Logger
	tokens     jwttokens.ITokens
}

func New(userrep userrepository.IUserRepository, sessionrep sessionrepository.ISessionRepository, logger *mylogger.Logger, tokens jwttokens.ITokens) *UserService {
	return &UserService{
		userrep:    userrep,
		sessionrep: sessionrep,
		logger:     logger,
		tokens:     tokens,
	}
}

func (us *UserService) SignUp(ctx context.Context, request *gatewayapi.SignUpRequest) (*gatewayapi.SignUpResponse, error) {
	xRequestId := ctx.Value("x_request_id")
	xRequestIdString, ok := xRequestId.(string)
	if !ok {
		err := errors.New("Не получилось извлечь xRequestId из контекста")
		us.logger.Logger.Error(err)
		return &gatewayapi.SignUpResponse{}, err
	}

	if len(request.Login) < 4 {
		err := errors.New("Логин должен состоять минимум из 4 символов")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.SignUpResponse{}, err
	}
	if len(request.Password) < 6 {
		err := errors.New("Пароль должен состоять минимум из 6 символов")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.SignUpResponse{}, err
	}
	if request.Email == "" {
		err := errors.New("Поле почта обязательна для заполнения")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.SignUpResponse{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось хэшировать пароль ", err)
	}
	user := userrepository.User{
		Login:    request.Login,
		Email:    request.Email,
		Password: hashedPassword,
	}
	err = us.userrep.Create(ctx, &user)
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось зарегистрировать пользователя ", err)
		return &gatewayapi.SignUpResponse{}, err
	}
	userForSession, err := us.userrep.GetByEmail(ctx, request.Email)
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось получить пользователя по email ", err)
		return &gatewayapi.SignUpResponse{}, err
	}
	refreshToken, err := us.tokens.CreateRefreshToken(userForSession.Id)
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось создать refresh token ", err)
		return &gatewayapi.SignUpResponse{}, err
	}
	dateNow := time.Now()
	expireAt := dateNow.AddDate(0, 1, 0)

	err = us.sessionrep.Create(ctx, &sessionrepository.Session{
		UsrId:        userForSession.Id,
		RefreshToken: refreshToken,
		ExpireAt:     expireAt,
	})
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось создать сессию для пользователя ", err)
		return &gatewayapi.SignUpResponse{}, err
	}
	return &gatewayapi.SignUpResponse{}, nil

}

func (us *UserService) SignIn(ctx context.Context, request *gatewayapi.SignInRequest) (*gatewayapi.SignInResponse, error) {
	xRequestId := ctx.Value("x_request_id")
	xRequestIdString, ok := xRequestId.(string)
	if !ok {
		err := errors.New("Не получилось извлечь xRequestId из контекста")
		us.logger.Logger.Error(err)
		return &gatewayapi.SignInResponse{}, err
	}

	if len(request.Login) < 4 {
		err := errors.New("Логин должен состоять минимум из 4 символов")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.SignInResponse{}, err
	}
	if len(request.Password) < 6 {
		err := errors.New("Пароль должен состоять минимум из 6 символов")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.SignInResponse{}, err
	}
	login := request.Login
	password := request.Password
	user, err := us.userrep.GetByUserAndPassword(ctx, login)
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось получить пользователя по логину и паролю, либо пользователя не существует", err)
		return &gatewayapi.SignInResponse{IsExist: false}, err
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось сделать проверку паролей, либо пароль неверный", err)
		return &gatewayapi.SignInResponse{}, err
	}
	session, err := us.sessionrep.Get(ctx, user.Id)
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось получить запись сессии для пользователя", err)
		return &gatewayapi.SignInResponse{}, err
	}
	refreshToken, err := us.tokens.CreateRefreshToken(user.Id)
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось создать refresh token", err)
		return &gatewayapi.SignInResponse{}, err
	}
	dateNow := time.Now()
	expireAt := dateNow.AddDate(0, 1, 0)
	session.RefreshToken = refreshToken
	session.ExpireAt = expireAt
	err = us.sessionrep.Update(ctx, session)
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось обновить информацию о сессии", err)
		return &gatewayapi.SignInResponse{}, err
	}
	accessToken, err := us.tokens.CreateAccessToken(user.Id)
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", "Не удалось создать access token", err)
		return &gatewayapi.SignInResponse{}, err
	}
	return &gatewayapi.SignInResponse{RefreshToken: refreshToken, AccessToken: accessToken, IsSignedUp: user.IsSignedUp, IsExist: true}, nil
}

func (us *UserService) UpdateAccessToken(ctx context.Context, request *gatewayapi.UpdateAccessTokenRequest) (*gatewayapi.UpdateAccessTokenResponse, error) {
	xRequestId := ctx.Value("x_request_id")
	xRequestIdString, ok := xRequestId.(string)
	if !ok {
		err := errors.New("Не получилось извлечь xRequestId из контекста")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.UpdateAccessTokenResponse{}, err
	}
	us.logger.Logger.Println("Запрос № ", xRequestIdString, " ", "Использована функиця обновления access token")
	if request.Userid == "" {
		err := errors.New("userId в функции UpdateAccessToken имеет nil значение")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.UpdateAccessTokenResponse{}, err
	}
	if request.RefreshToken == "" {
		err := errors.New("refreshToken в функции UpdateAccessToken имеет nil значение")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.UpdateAccessTokenResponse{}, err
	}
	userUuid, err := StringToUuuid(request.Userid)
	if err != nil {
		err := errors.New("Не удалось конвертировать userId в uuid")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.UpdateAccessTokenResponse{}, err
	}
	session, err := us.sessionrep.Get(ctx, userUuid)
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " Не удалось получить запись сессии из бд", err)
		err := errors.New("Не удалось получить запись сессии из бд")
		return &gatewayapi.UpdateAccessTokenResponse{}, err
	}
	if session.RefreshToken != request.RefreshToken {
		err := errors.New("refreshToken из запроса в функции UpdateAccessToken не совпадает с refreshToken из бд")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.UpdateAccessTokenResponse{}, err
	}
	_, refreshErr := us.tokens.VerifyToken(session.RefreshToken, us.logger)
	if refreshErr != nil {
		return nil, status.Errorf(codes.Unauthenticated, "неверный refresh token")
	}
	newAccessToken, newErr := us.tokens.CreateAccessToken(userUuid)
	if newErr != nil {
		return nil, status.Errorf(codes.Unauthenticated, "не удалось обновить access token")
	}
	us.logger.Logger.Println("Запрос № ", xRequestIdString, " ", "Новый access token - ", newAccessToken)
	return &gatewayapi.UpdateAccessTokenResponse{AccessToken: newAccessToken, RefreshToken: request.RefreshToken}, nil
}

func (us *UserService) IsLoginExist(ctx context.Context, req *gatewayapi.IsLoginExistRequest) (*gatewayapi.IsLoginExistResponse, error) {
	xRequestId := ctx.Value("x_request_id")
	xRequestIdString, ok := xRequestId.(string)
	if !ok {
		err := errors.New("Не получилось извлечь xRequestId из контекста")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.IsLoginExistResponse{}, err
	}
	if req.Login == "" {
		err := errors.New("login в функции IsLoginExist имеет nil значение")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.IsLoginExistResponse{}, err
	}
	isExist, _ := us.userrep.GetByLoginCheck(ctx, req.Login)
	return &gatewayapi.IsLoginExistResponse{IsExist: isExist}, nil
}

func (us *UserService) IsEmailExist(ctx context.Context, req *gatewayapi.IsEmailExistRequest) (*gatewayapi.IsEmailExistResponse, error) {
	xRequestId := ctx.Value("x_request_id")
	xRequestIdString, ok := xRequestId.(string)
	if !ok {
		err := errors.New("Не получилось извлечь xRequestId из контекста")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.IsEmailExistResponse{}, err
	}
	if req.Email == "" {
		err := errors.New("login в функции IsMailExist имеет nil значение")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.IsEmailExistResponse{}, err
	}
	isExist, _ := us.userrep.GetByEmailCheck(ctx, req.Email)
	return &gatewayapi.IsEmailExistResponse{IsExist: isExist}, nil
}

func (us *UserService) ChangeStatus(ctx context.Context, req *gatewayapi.ChangeStatusRequest) (*gatewayapi.ChangeStatusResponse, error) {
	xRequestId := ctx.Value("x_request_id")
	xRequestIdString, ok := xRequestId.(string)
	if !ok {
		err := errors.New("Не получилось извлечь xRequestId из контекста")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.ChangeStatusResponse{}, err
	}
	if req.Login == "" {
		err := errors.New("login в функции ChangeStatus имеет nil значение")
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		return &gatewayapi.ChangeStatusResponse{}, err
	}
	usr, err := us.userrep.GetByLogin(ctx, req.Login)
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		err = errors.New("Не удалось извлечь пользователя по логину из бд")
		return &gatewayapi.ChangeStatusResponse{}, err
	}
	usr.IsSignedUp = true
	err = us.userrep.Update(ctx, usr)
	if err != nil {
		us.logger.Logger.Error("Запрос № ", xRequestIdString, " ", err)
		err = errors.New("Не удалось поменять статус пользователя")
		return &gatewayapi.ChangeStatusResponse{}, err
	}
	return &gatewayapi.ChangeStatusResponse{}, nil
}

func StringToUuuid(value string) (uuid.UUID, error) {
	emptyUUID := uuid.UUID{}
	uuid, err := uuid.Parse(value)
	if err != nil {
		return emptyUUID, err
	}
	return uuid, err
}
