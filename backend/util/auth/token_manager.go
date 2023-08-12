package auth

import (
	"errors"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// JWTの操作
type ITokenManager interface {
	// UserIDを含む期限付きのトークンを生成する
	CreateToken(userID string, duration time.Duration) (string, error)

	// トークンからUserIDを取得する
	GetUserID(token string) (string, error)
}

type TokenManager struct {
	// 署名アルゴリズム
	signAlg jwa.SignatureAlgorithm

	// 暗号化アルゴリズム
	encryptAlg jwa.KeyEncryptionAlgorithm

	// RSA秘密鍵
	privateKey jwk.RSAPrivateKey

	// 署名元の名前
	issuer string
}

func NewTokenManager(issuer string, keyPath string) (*TokenManager, error) {
	if issuer == "" || keyPath == "" {
		return nil, errors.New("error: invalid parameter")
	}
	// RSA秘密鍵ファイル(PEM形式)を読み込む
	src, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, errors.New("error: failed to read private-key-file")
	}
	// 秘密鍵をPEM形式からJWKに変換する
	key, err := jwk.ParseKey(src, jwk.WithPEM(true))
	if err != nil {
		return nil, errors.New("error: failed to parse pem file")
	}
	var jwkPrivateKey jwk.RSAPrivateKey
	var ok bool
	if jwkPrivateKey, ok = key.(jwk.RSAPrivateKey); !ok {
		return nil, errors.New("error: failed to parse jwk-private-key from rsa-private-key")
	}
	return &TokenManager{
		signAlg:    jwa.RS256,
		encryptAlg: jwa.RSA_OAEP,
		issuer:     issuer,
		privateKey: jwkPrivateKey,
	}, nil
}

func (m *TokenManager) CreateToken(userID string, duration time.Duration) (string, error) {
	if userID == "" {
		return "", errors.New("error: invalid token parameter")
	}
	// 秘密鍵から公開鍵を取得する
	publickKey, err := m.privateKey.PublicKey()
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()

	// トークンに情報を含める
	token, err := jwt.NewBuilder().Issuer(m.issuer).IssuedAt(now).Subject(userID).Expiration(now.Add(duration)).Build()
	if err != nil {
		return "", err
	}
	// 秘密鍵を使用してトークンの署名を行う
	signed, err := jwt.Sign(token, jwt.WithKey(m.signAlg, m.privateKey))
	if err != nil {
		return "", err
	}
	// 公開鍵を使用して暗号化を行う
	encrypted, err := jwe.Encrypt(signed, jwe.WithKey(m.encryptAlg, publickKey))
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}

func (m *TokenManager) GetUserID(token string) (string, error) {
	// 秘密鍵から公開鍵を取得する
	publickKey, err := m.privateKey.PublicKey()
	if err != nil {
		return "", err
	}
	// 秘密鍵を使用して復号化を行う
	decrypted, err := jwe.Decrypt([]byte(token), jwe.WithKey(m.encryptAlg, m.privateKey))
	if err != nil {
		return "", err
	}
	// 公開鍵を使用してトークンの署名を検証する
	verifyed, err := jwt.Parse(decrypted, jwt.WithKey(m.signAlg, publickKey))
	if err != nil {
		return "", errors.New("error: failed to verify token")
	}
	// トークンからuserIDを取得する
	return verifyed.Subject(), nil
}
