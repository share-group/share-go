package jwt

import (
	"crypto/rsa"
	"embed"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"os"
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
// publicCertPath-公钥证书地址; privateCertPath-私钥证书地址; signingMethod-加密算法
func NewJwt(publicCertPath, privateCertPath string, signingMethod *jwt.SigningMethodRSA) *Jwt {
	j := &Jwt{signingMethod: signingMethod}

	privateKeyBytes, err := os.ReadFile(privateCertPath)
	if err != nil {
		panic(err)
	}

	j.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		panic(err)
	}

	// 从文件加载公钥
	publicKeyBytes, err := os.ReadFile(publicCertPath)
	if err != nil {
		panic(err)
	}

	j.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		panic(err)
	}

	return j
}

// 初始化一个带有证书的jwt实例
//
// publicCertPath-公钥证书地址; privateCertPath-私钥证书地址; signingMethod-加密算法
func NewJwtWithEmbed(publicCertPath, privateCertPath embed.FS, signingMethod *jwt.SigningMethodRSA) *Jwt {
	j := &Jwt{signingMethod: signingMethod}

	privateKeyBytes, err := publicCertPath.ReadFile("rsa_public_key.pem")
	if err != nil {
		panic(err)
	}

	j.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		panic(err)
	}

	// 从文件加载公钥
	publicKeyBytes, err := privateCertPath.ReadFile("rsa_private_key.pem")
	if err != nil {
		panic(err)
	}

	j.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		panic(err)
	}

	return j
}

// 加密数据
//
// data-待加密的数据; duration-有效时长(单位：毫秒)
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
	err = json.Unmarshal([]byte(token.Claims.(*jwtClaims).Data), &result)
	return result, err
}
