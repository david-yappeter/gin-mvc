package service

import (
	"context"
	"myapp/config"
	"myapp/entity"
	"strings"
	"time"

	"github.com/google/uuid"
)

func UserCreate(ctx context.Context, input entity.User) (*entity.User, error) {
	// Google UUID
	input.ID = uuid.New().String()
	input.Email = strings.ToLower(input.Email)

	if _, err := UserGetByEmail(ctx, input.Email); err != nil {
		if err != ErrRecordNotFound {
			return nil, err
		}
	} else {
		return nil, ErrRecordFound
	}

	db := config.GetDB()

	if err := db.Model(input).Create(&input).Error; err != nil {
		return nil, err
	}

	return &input, nil
}

func UserGetByID(ctx context.Context, id string) (*entity.User, error) {
	db := config.GetDB()

	var user entity.User
	if err := db.Model(user).Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func UserGetByEmail(ctx context.Context, email string) (*entity.User, error) {
	db := config.GetDB()

	var user entity.User
	if err := db.Model(user).Where("email like ?", email).Take(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func UserRegister(ctx context.Context, input entity.User) (string, string, error) {
	createUser, err := UserCreate(ctx, input)
	if err != nil {
		return "", "", err
	}

	accToken, remToken, err := JwtNew(createUser.ID)
	if err != nil {
		return "", "", err
	}

	_, err = UserUpdateRememberToken(ctx, createUser.ID, remToken)
	if err != nil {
		return "", "", err
	}

	return accToken, remToken, err
}

func UserUpdateRememberToken(ctx context.Context, id string, token string) (string, error) {
	db := config.GetDB()

	if err := db.Model(entity.User{}).Where("id = ?", id).Update("remember_token", token).Error; err != nil {
		return "", nil
	}

	return "Success", nil
}

func UserLogin(ctx context.Context, input entity.UserLogin) (string, string, error) {
	getUser, err := UserGetByEmail(ctx, input.Email)
	if err != nil {
		return "", "", ErrRecordNotFound
	}

	if err = getUser.ComparePass(input.Password); err != nil {
		return "", "", err
	}

	var rememberToken string
	if getUser.RememberToken != nil {
		_, err := JwtValidate(*getUser.RememberToken)
		if err != nil {
			rememberToken, err = JwtCreateRememberToken(getUser.ID)
			if err != nil {
				return "", "", err
			}
			_, err = UserUpdateRememberToken(ctx, getUser.ID, rememberToken)
			if err != nil {
				return "", "", err
			}
		} else {
			rememberToken = *getUser.RememberToken
		}
	}

	accToken, err := jwtCreate(getUser.ID, time.Minute*15)

	return accToken, rememberToken, err
}
