package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/HC74/kratosx/config"
	rds "github.com/HC74/kratosx/core/redis"
	"sync"
	"time"

	kratosJwt "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

type Jwt interface {
	NewToken(m map[string]any) (string, error)
	Parse(ctx context.Context, dst any) error
	ParseMapClaims(ctx context.Context) (map[string]any, error)
	IsWhitelist(path string) bool
	IsBlacklist(token string) bool
	AddBlacklist(token string)
	GetToken(ctx context.Context) string
	SetToken(ctx context.Context, token string) context.Context
	Renewal(ctx context.Context) (string, error)
}

type jwt struct {
	conf *config.JWT
	rw   sync.RWMutex
}

var (
	instance *jwt
	tokenKey struct{}
)

const (
	blackPrefix = "token_black"
)

// Instance 获取email对象实例
func Instance() Jwt {
	return instance
}

func Init(conf *config.JWT) {
	if conf == nil {
		return
	}
	instance = &jwt{conf: conf}
}

// NewToken is create jwt []byte
func (j *jwt) NewToken(m map[string]any) (string, error) {
	if j == nil {
		return "", errors.New("jwt config not enable or configure")
	}

	m["exp"] = jwtv4.NewNumericDate(time.Now().Add(j.conf.Expire + time.Second)) // 过期时间
	m["nbf"] = jwtv4.NewNumericDate(time.Now())                                  // 生效时间
	m["iat"] = jwtv4.NewNumericDate(time.Now())                                  // 签发时间

	keyFunc := func(token *jwtv4.Token) (any, error) {
		return []byte(j.conf.Secret), nil
	}

	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, jwtv4.MapClaims(m))

	key, err := keyFunc(token)
	if err != nil {
		return "", err
	}

	return token.SignedString(key)
}

func (j *jwt) Parse(ctx context.Context, dst any) error {
	token, is := kratosJwt.FromContext(ctx)
	if !is {
		return errors.New("token miss")
	}
	claims, is := token.(jwtv4.MapClaims)
	if !is {
		return errors.New("token format error")
	}

	body, err := json.Marshal(claims)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, dst)
}

func (j *jwt) ParseMapClaims(ctx context.Context) (map[string]any, error) {
	token, is := kratosJwt.FromContext(ctx)
	if !is {
		return nil, errors.New("token miss")
	}
	claims, is := token.(jwtv4.MapClaims)
	if !is {
		return nil, errors.New("token format error")
	}
	return claims, nil
}

func (j *jwt) GetToken(ctx context.Context) string {
	token, _ := ctx.Value(tokenKey).(string)
	return token
}

func (j *jwt) SetToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

// Renewal 续订
func (j *jwt) Renewal(ctx context.Context) (string, error) {
	res, err := j.ParseMapClaims(ctx)
	if err != nil {
		return "", err
	}

	// 判断token失效是否超过10s
	exp := int64(res["exp"].(float64))
	now := time.Now().Unix()
	if exp > now {
		return "", errors.New("token is alive")
	}

	if now-exp > int64(j.conf.Renewal.Seconds()) {
		return "", errors.New("token is over max renewal time")
	}

	return j.NewToken(res)
}

func (j *jwt) IsWhitelist(path string) bool {
	j.rw.RLock()
	defer j.rw.RUnlock()
	return j.conf.Whitelist[path]
}

func (j *jwt) IsBlacklist(token string) bool {
	rd := rds.Instance().Get(j.conf.Redis)
	if rd == nil {
		return false
	}
	is, _ := rd.HExists(context.Background(), blackPrefix, token).Result()
	return is
}

func (j *jwt) AddBlacklist(token string) {
	rd := rds.Instance().Get(j.conf.Redis)
	rd.HSet(context.Background(), blackPrefix, token, 1, j.conf.Expire)
}
