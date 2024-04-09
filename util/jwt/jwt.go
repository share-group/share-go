package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"time"
)

type Jwt struct {
	privateKey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	signingMethod *jwt.SigningMethodRSA
}

type jwtClaims struct {
	Data string `json:"data"`
	jwt.StandardClaims
}

// 初始化一个带有证书的jwt实例
//
// publicKeyBytes-公钥证书; privateKeyBytes-私钥证书; signingMethod-加密算法
func NewJwt(publicKeyBytes, privateKeyBytes []byte, signingMethod *jwt.SigningMethodRSA) *Jwt {
	var err error
	j := &Jwt{signingMethod: signingMethod}

	// 加载公钥
	j.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		panic(err)
	}

	// 加载私钥
	j.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		panic(err)
	}

	return j
}

// 加密数据
//
// data-待加密的数据; duration-有效时长(单位：秒)
func (j *Jwt) Encrypt(data any, duration int64) (string, error) {
	b, _ := json.Marshal(data)
	claims := jwtClaims{Data: string(b), StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Unix() + duration}}
	token := jwt.NewWithClaims(j.signingMethod, claims)
	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解密数据
//
// data-待解密的数据
func (j *Jwt) Decrypt(data string) (map[string]any, error) {
	token, err := jwt.ParseWithClaims(data, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.publicKey, nil
	})

	if err != nil {
		return make(map[string]any), err
	}

	result := make(map[string]any)
	claimsBytes, _ := json.Marshal(token.Claims)
	err = json.Unmarshal(claimsBytes, &result)
	delete(result, "data")
	err = json.Unmarshal([]byte(token.Claims.(*jwtClaims).Data), &result)
	return result, err
}
