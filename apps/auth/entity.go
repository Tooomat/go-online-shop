package auth

import (
	"strings"
	"time"

	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/Tooomat/go-online-shop/utility"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	ROLE_Super_Admin Role = "super_admin"
	ROLE_Admin       Role = "admin"
	ROLE_User        Role = "user"
)

type AuthEntity struct {
	Id        int       `db:"id"`
	Email     string    `db:"email"`
	PublicId  uuid.UUID `db:"public_id"` //security untuk keperluan token
	Password  string    `db:"password"`
	Role      Role      `db:"role"`
	CreatedAt time.Time `db:"created_time"`
	UpdateAt  time.Time `db:"update_time"`
}

func NewFromSeed(req RequestPayLoadSuperAdmin) AuthEntity {
	return AuthEntity{
		Email:     req.Email,
		PublicId:  uuid.New(),
		Password:  req.Password,
		Role:      ROLE_Super_Admin,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
}

// step2(register): menangkap req dari user dalam bentuk json kemudian di masukkan ke struct
func NewFromRegisterRequest(req RegisterRequestPayload) AuthEntity {
	return AuthEntity{
		Email:     req.Email,
		PublicId:  uuid.New(), //bernilai random
		Password:  req.Password,
		Role:      ROLE_User,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
}

func NewFromLoginRequest(req LoginRequestPayLoad) AuthEntity {
	return AuthEntity{
		Email:    req.Email,
		Password: req.Password,
	}
}

// validasi
func (a AuthEntity) AuthIsValid() (err error) {
	if err = a.ValidateEmail(); err != nil {
		return
	}
	if err = a.ValidatePassword(); err != nil {
		return
	}
	return
}

// validasi email
func (a AuthEntity) ValidateEmail() (err error) {
	if a.Email == "" { //email kosong
		return response.ErrEmailRequired
	}

	if emails := strings.Split(a.Email, "@"); len(emails) != 2 { //reyy@roy@gmail.com
		return response.ErrEmailInvalid
	}

	return nil
}

// validasi password
func (a AuthEntity) ValidatePassword() (err error) {
	if a.Password <= "" {
		return response.ErrPassRequired
	}

	if len(a.Password) <= 6 {
		return response.ErrPassInvalid
	}
	return nil
}

func (a AuthEntity) IsExsist() bool {
	return a.Id != 0
}

// hash password
func (a *AuthEntity) EncriyptPassword(salt int) (err error) {

	encriyptedPass, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(encriyptedPass)

	return nil
}

// verify password
func (a AuthEntity) VerifyPasswordFromEncrypted(plain string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(plain))
}
func (a AuthEntity) VerifyPasswordFromPlain(encryp string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(encryp), []byte(a.Password))
}
func (a AuthEntity) GenerateAccessToken(secret string) (tokenString string, err error) {
	return utility.CreateAccessToken(a.PublicId.String(), string(a.Role), secret)
}
func (a AuthEntity) GenerateRefreshToken(secret string) (tokenString string, err error) {
	return utility.CreateRefreshToken(a.PublicId.String(), string(a.Role), secret)
}
