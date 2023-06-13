package service

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/ahfrd/grpc/auth-client/src/proto/auth"
	repository "github.com/ahfrd/grpc/auth-client/src/repository/auth"
	helpers "github.com/ahfrd/grpc/auth-client/src/utils/helpers"
)

type ControllerAuth struct {
	repository.AuthenticationTable
	Jwt helpers.JwtWrapper
	auth.UnimplementedAuthServiceServer
}

type Response struct {
	Token string `json:"token"`
}

func (o *ControllerAuth) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	phoneNumber := req.PhoneNumber
	password := req.Password

	validatePhoneNumber := helpers.RegexNumber(phoneNumber)
	if validatePhoneNumber == false {
		return &auth.RegisterResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: "General Errors",
		}, nil
	}
	hashPassword := helpers.PasswordHash256(phoneNumber, password)

	checkData, db, err := o.ValidateDataByPhoneNumber(phoneNumber)
	db.Close()
	if err != nil {
		return &auth.RegisterResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: err.Error(),
		}, nil
	}
	if checkData.PhoneNumber != "" {
		return &auth.RegisterResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: "Data Exist",
		}, nil
	}
	_, db, err = o.RegisterUser(phoneNumber, hashPassword)
	db.Close()
	if err != nil {
		return &auth.RegisterResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: err.Error(),
		}, nil
	}

	return &auth.RegisterResponse{
		RessponseCode:   http.StatusAccepted,
		ResponseMessage: "Succsess",
	}, nil

}

func (o *ControllerAuth) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	password := req.Password
	phoneNumber := req.PhoneNumber
	validatePhoneNumber := helpers.RegexNumber(phoneNumber)
	if validatePhoneNumber == false {
		return &auth.LoginResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: "General Errors",
		}, nil
	}
	hashPassword := helpers.PasswordHash256(phoneNumber, password)

	checkUser, db, err := o.SelectByPhoneNumber(phoneNumber)
	db.Close()
	if err != nil {
		return &auth.LoginResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: err.Error(),
		}, nil
	}
	if checkUser.Id == "" {
		return &auth.LoginResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: "Username / Password salah",
		}, nil
	}
	fmt.Println(checkUser.PhoneNumber)

	currentTime := time.Now()
	if currentTime.String() <= checkUser.NextLoginDate {
		return &auth.LoginResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: "You can login again at " + checkUser.NextLoginDate,
		}, nil
	}
	fmt.Println(hashPassword)
	fmt.Println(checkUser)
	if checkUser.Password != hashPassword {
		count := checkUser.LoginRetry + 1
		if count > 3 {
			math := int(math.Pow((float64(count)-3), 2) * 600)
			login_again := time.Now().Add(time.Second * time.Duration(math))
			_, db, err := o.UpdateRetryNextLogin(phoneNumber, login_again.String(), count)
			db.Close()
			if err != nil {
				return &auth.LoginResponse{
					RessponseCode:   http.StatusBadRequest,
					ResponseMessage: err.Error(),
				}, nil
			}
		} else {
			_, db, err := o.UpdateRetryLogin(phoneNumber, count)
			db.Close()
			if err != nil {
				return &auth.LoginResponse{
					RessponseCode:   http.StatusBadRequest,
					ResponseMessage: err.Error(),
				}, nil
			}
		}
		return &auth.LoginResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: "Username / Password salah",
		}, nil
	}
	idInt, _ := strconv.Atoi(checkUser.Id)
	token, _ := o.Jwt.GenerateToken(idInt, checkUser.PhoneNumber)
	if err != nil {
		return &auth.LoginResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: err.Error(),
		}, nil
	}
	_, db, err = o.UpdateRetryLogin(phoneNumber, 0)
	db.Close()
	if err != nil {
		return &auth.LoginResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: err.Error(),
		}, nil
	}
	return &auth.LoginResponse{
		RessponseCode:   http.StatusContinue,
		ResponseMessage: "Succsses",
		ResponseData: &auth.TokenLogin{
			Token: token,
		},
	}, nil
}

func (o *ControllerAuth) Validate(ctx context.Context, req *auth.ValidateRequest) (*auth.ValidateResponse, error) {
	claims, err := o.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &auth.ValidateResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: err.Error(),
		}, nil
	}
	checkEmail, db, err := o.ValidateDataByPhoneNumber(claims.PhoneNumber)
	db.Close()
	if err != nil {
		return &auth.ValidateResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: err.Error(),
		}, nil
	}
	if checkEmail.PhoneNumber == "" {
		return &auth.ValidateResponse{
			RessponseCode:   http.StatusBadRequest,
			ResponseMessage: "Data Not Found",
		}, nil
	}
	idInt, _ := strconv.Atoi(checkEmail.Id)

	return &auth.ValidateResponse{
		RessponseCode:   http.StatusOK,
		ResponseMessage: err.Error(),
		UserId:          int64(idInt),
	}, nil
}
