package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomAccessClaims struct {
	ID   string `json:"id"`
	Role uint32 `json:"role"`
	jwt.RegisteredClaims
}

type CustomRefreshClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

type JWT struct {
	AccessSignKey  string
	RefreshSignKey string
	AccessExpire   time.Duration
	RefreshExpire  time.Duration
}

func NewJWT(AccessSignKey, RefreshSignKey string, AccessExpire, RefreshExpire time.Duration) *JWT {
	return &JWT{
		AccessSignKey,
		RefreshSignKey,
		AccessExpire,
		RefreshExpire,
	}
}

func (j *JWT) SignAccessToken(userID string, role uint32) (string, error) {
	now := time.Now()
	claims := CustomAccessClaims{
		ID:   userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.AccessExpire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	at, err := token.SignedString([]byte(j.AccessSignKey))
	if err != nil {
		return "", err
	}

	return at, nil
}

func (j *JWT) SignRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := CustomRefreshClaims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.RefreshExpire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rt, err := token.SignedString([]byte(j.RefreshSignKey))
	if err != nil {
		return "", err
	}

	return rt, nil
}

func (j *JWT) ValidateAccessToken(token string) (string, uint32, error) {
	ac := CustomAccessClaims{}
	t, err := jwt.ParseWithClaims(token, &ac, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.AccessSignKey), nil
	})
	if err != nil {
		return "", 0, errors.New("unauthorized")
	}

	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return "", 0, errors.New("unauthorized")
	}

	if t.Valid {
		if ac.ExpiresAt.Time.Unix() < time.Now().UTC().Unix() {
			return "", 0, errors.New("token expired")
		}

		if ac.ID == "" || ac.Role == 0 {
			return "", 0, errors.New("unauthorized")
		}

		return ac.ID, ac.Role, nil
	}

	return "", 0, errors.New("unauthorized")
}

func (j *JWT) ValidateRefreshToken(token string) (string, error) {
	rc := CustomRefreshClaims{}
	t, err := jwt.ParseWithClaims(token, &rc, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.RefreshSignKey), nil
	})
	if err != nil {
		return "", errors.New("unauthorized")
	}

	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return "", errors.New("unauthorized")
	}

	if t.Valid {
		if rc.ExpiresAt.Time.Unix() < time.Now().UTC().Unix() {
			return "", errors.New("unauthorized")
		}

		if rc.ID == "" {
			return "", errors.New("unauthorized")
		}

		return rc.ID, nil
	}

	return "", errors.New("unauthorized")
}
