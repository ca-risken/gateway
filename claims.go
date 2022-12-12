package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type claimsClient struct {
	region          string
	userIdpKey      string
	idpProviderName []string
	verify          bool
}

type claimsInterface interface {
	getClaims(ctx context.Context, tokenString string) (*jwt.MapClaims, error)
	getUserName(claims *jwt.MapClaims) string
	getUserIdpKey(claims *jwt.MapClaims) string
}

func newClaimsClient(region, userIdpKey string, idpProviderName []string, verify bool) *claimsClient {
	return &claimsClient{
		region:          region,
		userIdpKey:      userIdpKey,
		idpProviderName: idpProviderName,
		verify:          verify,
	}
}

func (c *claimsClient) getClaims(ctx context.Context, tokenString string) (*jwt.MapClaims, error) {
	var claims *jwt.MapClaims
	var err error
	if c.verify {
		claims, err = c.getClaimsForALB(ctx, tokenString)
	} else {
		claims, err = c.getClaimsWithoutVerify(tokenString)
	}
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (c *claimsClient) getClaimsForALB(ctx context.Context, tokenString string) (*jwt.MapClaims, error) {
	jwt.DecodePaddingAllowed = true
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid is not found in jwt header")
		}
		keyURL := fmt.Sprintf(PublicKeyURLTemplate, c.region, kid)
		key, err := fetchALBPublicKey(ctx, keyURL)
		if err != nil {
			return nil, fmt.Errorf("failed to get public key, err: %w", err)
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &claims, nil
}

const (
	PublicKeyURLTemplate = "https://public-keys.auth.elb.%s.amazonaws.com/%s"
)

func fetchALBPublicKey(ctx context.Context, keyURL string) (*ecdsa.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, keyURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to new GET request for %s, err: %w", keyURL, err)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key from %s, err: %w", keyURL, err)
	}
	defer resp.Body.Close()
	pem, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s, err: %w", keyURL, err)
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM(pem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key from %s, err: %w", keyURL, err)
	}

	return publicKey, nil
}

func (c *claimsClient) getClaimsWithoutVerify(tokenString string) (*jwt.MapClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT string pattern")
	}
	// Decode JWT
	claimBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var claims jwt.MapClaims
	if err := json.NewDecoder(bytes.NewBuffer(claimBytes)).Decode(&claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

func (c *claimsClient) getUserName(claims *jwt.MapClaims) string {
	value, ok := (*claims)["username"]
	if !ok {
		return ""
	}
	userName, ok := value.(string)
	if !ok {
		return ""
	}
	for _, idp := range c.idpProviderName {
		if strings.HasPrefix(strings.ToLower(userName), strings.ToLower(idp)+"_") {
			userName = strings.Replace(userName, idp+"_", "", 1)
			break
		}
	}
	return userName
}

func (c *claimsClient) getUserIdpKey(claims *jwt.MapClaims) string {
	value, ok := (*claims)[c.userIdpKey]
	if !ok {
		return ""
	}
	strValue, ok := value.(string)
	if !ok {
		return ""
	}
	return strValue
}
