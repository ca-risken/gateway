package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// generateAccessToken return random accessToken text
func generateAccessToken() string {
	buf := make([]byte, 64)
	_, _ = rand.Read(buf)
	return base64.RawURLEncoding.EncodeToString(buf)
}

// encodeAccessToken is encoding AccessToken. Format: urlEncode({owner_id}@{access_token_id}@{plain_text})
func encodeAccessToken(ownerID, accessTokenID uint32, plainText string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprint(ownerID) + "@" + fmt.Sprint(accessTokenID) + "@" + plainText))
}

// decodeAccessToken is decoding AccessToken, and return access_token_id, plain_text. Format: urlEncode({owner_id}@{access_token_id}@{plain_text})
func decodeAccessToken(ctx context.Context, accessToken string) (ownerID, accessTokenID uint32, plainText string, err error) {
	tokenBody, err := base64.RawURLEncoding.DecodeString(accessToken)
	if err != nil {
		appLogger.Warnf(ctx, "Failed to decode access token, err=%+v", err)
		return 0, 0, "", err
	}
	parts := strings.SplitN(string(tokenBody), "@", 3)
	if len(parts) != 3 {
		return 0, 0, "", errors.New("invalid token, token must contain three values")
	}
	pID, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		appLogger.Warnf(ctx, "Failed to parse owner_id, raw=%s, err=%+v", parts[0], err)
		return 0, 0, "", err
	}
	aID, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		appLogger.Warnf(ctx, "Failed to parse access_token_id, raw=%s, err=%+v", parts[1], err)
		return 0, 0, "", err
	}
	return uint32(pID), uint32(aID), parts[2], nil
}
