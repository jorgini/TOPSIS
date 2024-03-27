package usecase

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
	"webApp/app/configs"
	"webApp/app/entity"
	"webApp/app/repository"
)

const (
	charset = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890`!@#$%^&*()_+{}:[];'./,<>?=-~"
)

type SessionService struct {
	repo repository.Session
}

func NewSessionService(repo repository.Session) *SessionService {
	return &SessionService{
		repo: repo,
	}
}

func (s *SessionService) GenerateToken(ctx context.Context, uid int64, cfg *configs.AppConfig) (entity.Tokens, error) {
	claims := jwt.MapClaims{
		"user_id": uid,
		"exp":     time.Now().Add(cfg.AccessTTL).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	access, err := accessToken.SignedString([]byte(cfg.SignKey))
	if err != nil {
		return entity.Tokens{}, err
	}

	refresh := make([]byte, 32)
	gen := rand.New(rand.NewSource(time.Now().Unix()))
	for i := range refresh {
		refresh[i] = charset[gen.Intn(len(charset))]
	}

	session := entity.Session{
		Uid:       uid,
		Token:     string(refresh),
		ExpiredAt: time.Now().Add(cfg.RefreshTTL),
	}

	if err := s.repo.InsertRefreshToken(ctx, session); err != nil {
		return entity.Tokens{}, err
	}

	return entity.Tokens{Access: access, Refresh: string(refresh)}, nil
}

func (s *SessionService) RefreshToken(ctx context.Context, refresh string, cfg *configs.AppConfig) (entity.Tokens, error) {
	uid, err := s.repo.GetUIDByToken(ctx, refresh)
	if err != nil {
		return entity.Tokens{}, err
	}

	return s.GenerateToken(ctx, uid, cfg)
}
