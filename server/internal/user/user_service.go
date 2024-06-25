package user

import (
	"context"
	"server/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// not a good implementation, store in .env file instead
const (
	secretKey = "ksecretey"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return &CreateUserRes{}, err
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	createdUser, err := s.Repository.CreateUser(ctx, user)
	if err != nil {
		return &CreateUserRes{}, err
	}

	response := &CreateUserRes{
		ID:       strconv.Itoa(int(createdUser.ID)),
		Username: createdUser.Username,
		Email:    createdUser.Email,
	}

	return response, nil
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// get user by email
	user, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &LoginUserRes{}, err
	}

	// compare password
	if err := util.CheckPassword(req.Password, user.Password); err != nil {
		return &LoginUserRes{}, err
	}

	// jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(user.ID)),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &LoginUserRes{}, err
	}

	// respose
	response := &LoginUserRes{
		accessToken: ss,
		ID:          strconv.Itoa(int(user.ID)),
		Username:    user.Username,
	}

	return response, nil
}
