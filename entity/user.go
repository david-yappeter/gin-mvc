package entity

import "myapp/tools"

// MODEl: User
type User struct {
	ID               string  `json:"id" gorm:"type:varchar(100);primaryKey"`
	Name             string  `json:"name" gorm:"type:varchar(255);not null"`
	Email            string  `json:"email" gorm:"type:varchar(255);not null"`
	Password         string  `json:"-" gorm:"type:varchar(255);not null"`
	Address          *string `json:"address,omitempty" gorm:"type:text;null;default:NULL"`
	PhoneCountryCode *string `json:"phone_country_code,omitempty" gorm:"type:text;null;default:NULL"`
	Phone            *string `json:"phone,omitempty" gorm:"type:text;null;default:NULL"`
	Balance          float64 `json:"balance" gorm:"type:decimal(15,2);not null;default:0"`
	RememberToken    *string `json:"-" gorm:"type:varchar(255);null;default:NULL"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) ComparePass(s string) error {
	return tools.ComparePass(u.Password, s)
}

// MODEL: UserRegister
type UserRegister struct {
	Name             string  `json:"name" binding:"required,max=100"`
	Email            string  `json:"email" binding:"required,email,max=100"`
	Password         string  `json:"password" binding:"required,max=100"`
	Address          *string `json:"address,omitempty" binding:"omitempty,min=1"`
	PhoneCountryCode *string `json:"phone_country_code,omitempty" binding:"omitempty,min=1"`
	Phone            *string `json:"phone,omitempty" binding:"required_with=PhoneCountryCode,omitempty,min=5"`
}

func (u *UserRegister) HashedPass() string {
	return tools.HashPassword(u.Password)
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,max=100"`
	Password string `json:"password" binding:"required,max=100"`
}
