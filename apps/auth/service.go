package auth

//main
import (
	"context"
	"log"
	"time"

	"github.com/Tooomat/go-online-shop/external/cache"
	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/Tooomat/go-online-shop/internal/configs"
	"github.com/Tooomat/go-online-shop/utility"
)

type authService struct {
	repo Repository //implement repository
	rdb  cache.RedisRepository
}

func newAuthService(repo Repository, rdb cache.RedisRepository) authService {
	return authService{
		repo: repo,
		rdb:  rdb,
	}
}

// mengapa menggunakan struct method? lebih rapi, dan pengelompokan (mengambil method, variabel dari file lain)
func (s authService) RegisterSeedService(ctx context.Context, req RequestPayLoadSuperAdmin) (err error) {
	authEntity := NewFromSeed(req)

	//validasi inputan super_admin
	if err = authEntity.AuthIsValid(); err != nil {
		return
	}

	//enkripsi password
	if err = authEntity.EncriyptPassword(int(configs.Cfg.App.Encryption.Salt)); err != nil {
		return
	}

	//cek apakah terdapat super_admin
	count, err := s.repo.CekSuperAdmin(ctx)
	if err != nil {
		return
	}

	if count == 0 { //tidak ada
		if err = s.repo.CreatedAuth(ctx, authEntity); err != nil {
			return
		}
		log.Println("✅ super Admin berhasil di-seed")
	} else {
		log.Println("⚠️ super Admin sudah ada, skip seeding")
	}
	return
}

// (register) -> 3 steps (req user(json) -> input to struct(db) -> mapping to databases)
func (s authService) registerService(ctx context.Context, req RegisterRequestPayload) (err error) {
	authEntity := NewFromRegisterRequest(req)

	//validasi inputan user
	if err = authEntity.AuthIsValid(); err != nil {
		return
	}

	//encryption pass
	if err = authEntity.EncriyptPassword(int(configs.Cfg.App.Encryption.Salt)); err != nil {
		return
	}

	//cari apakah ada email yang sama pada db
	model, err := s.repo.GetAuthByEmail(ctx, authEntity.Email)
	if err != nil { 
		if err != response.ErrNotFound {
			return
		}
	}

	if model.IsExsist() {
		return response.ErrEmailAlreadyUsed
	}

	return s.repo.CreatedAuth(ctx, authEntity)
}

// login user
func (s authService) loginService(ctx context.Context, req LoginRequestPayLoad) (accessToken, refreshToken string, err error) {
	authEntity := NewFromLoginRequest(req)

	//validate input user
	if err = authEntity.ValidateEmail(); err != nil {
		return
	}
	if err = authEntity.ValidatePassword(); err != nil {
		return
	}

	//mencari email yang sama
	model, err := s.repo.GetAuthByEmail(ctx, authEntity.Email)
	if err != nil {
		return
	}

	//apakah password sama atau tidak
	if err = authEntity.VerifyPasswordFromPlain(model.Password); err != nil {
		err = response.ErrPassNotMatch
		return
	}

	//generate token
	accessToken, err = model.GenerateAccessToken(configs.Cfg.App.Encryption.JWTAccessSecret)
	if err != nil {
		return
	}

	refreshToken, err = model.GenerateRefreshToken(configs.Cfg.App.Encryption.JWTRefreshSecret)
	if err != nil {
		return
	}

	// save refres token to redis
	err = s.rdb.SaveRefreshToken(ctx, model.PublicId.String(), refreshToken, utility.REFRESH_TOKEN_TIME_TO_LIFE)
	if err != nil {
		return
	}

	return
}

func (s authService) LogoutService(ctx context.Context, accessToken string, exp time.Time, id string) (err error) {
	// sisa TTL
	ttl := time.Until(exp)
	if ttl < 0 {
		ttl = 0
	}

	// blacklist token di redis
	err = s.rdb.BlacklistAccessToken(ctx, accessToken, ttl)
	if err != nil {
		return 
	}

	// remove refresh token dari redis
	userKey := "refresh:" + id
	err = s.rdb.DeleteKey(ctx, userKey)
	if err != nil {
		return
	}

	return
}

func (s authService) RefreshAccessService(ctx context.Context, rtCookie string) (newAccessToken string, err error) {
	//validasi refresh token di cookie
	id, role, err := utility.ParseRefreshToken(rtCookie, configs.Cfg.App.Encryption.JWTRefreshSecret)
	if err != nil {
		err = response.ErrRefreshTokenInvalid
		return
	}

	// validasi refresh token di redis
	res, err := s.rdb.GetRefreshToken(ctx, id)
	if err != nil || res != rtCookie {
		err = response.ErrRefreshTokenEXP
		return
	}

	// create access token
	newAccessToken, err = utility.CreateAccessToken(id, role, configs.Cfg.App.Encryption.JWTAccessSecret)
	if err != nil {
		return
	}

	return
}
