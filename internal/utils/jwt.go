package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func CreateToken(subject string, privateKeyBase64 string) (string, error) {
	// 1. Base64 解码私钥
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode private key: %v", err)
	}

	// 2. 解析 PKCS#8 格式私钥
	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// 3. 转换为 *rsa.PrivateKey
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("not an RSA private key")
	}

	// 4. 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub": subject,
	})

	// 5. 使用 RSA 私钥签名
	return token.SignedString(rsaPrivateKey)
}

// ParseToken 使用 RSA 公钥验证 JWT 并返回 subject
func ParseToken(tokenString string, publicKeyBase64 string) (string, error) {
	// 1. 移除 "Bearer " 前缀（如果有）
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 2. Base64 解码公钥
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %v", err)
	}

	// 3. 解析 PKIX 格式公钥
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}

	// 4. 转换为 *rsa.PublicKey
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("not an RSA public key")
	}

	// 5. 解析并验证 JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 检查签名算法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return rsaPublicKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	// 6. 提取 subject
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if subject, ok := claims["sub"].(string); ok {
			return subject, nil
		}
		return "", errors.New("subject not found in token")
	}

	return "", errors.New("invalid token")
}
