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

const orgTokenPrefix = "risken_org_"

// generateAccessToken return random accessToken text
func generateAccessToken() string {
	buf := make([]byte, 64)
	_, _ = rand.Read(buf)
	return base64.RawURLEncoding.EncodeToString(buf)
}

// encodeAccessToken is encoding AccessToken. Format: urlEncode({unit_id}@{access_token_id}@{plain_text})
func encodeAccessToken(unitID, accessTokenID uint32, plainText string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprint(unitID) + "@" + fmt.Sprint(accessTokenID) + "@" + plainText))
}

func encodeOrgAccessToken(unitID, accessTokenID uint32, plainText string) string {
	return orgTokenPrefix + encodeAccessToken(unitID, accessTokenID, plainText)
}

// decodeAccessToken is decoding AccessToken, and return access_token_id, plain_text. Format: urlEncode({unit_id}@{access_token_id}@{plain_text})
func decodeAccessToken(ctx context.Context, accessToken string) (unitID, accessTokenID uint32, plainText string, err error) {
	tokenBody, err := base64.RawURLEncoding.DecodeString(accessToken)
	if err != nil {
		appLogger.Warnf(ctx, "Failed to decode access token, err=%+v", err)
		return 0, 0, "", err
	}
	parts := strings.SplitN(string(tokenBody), "@", 3)
	if len(parts) != 3 {
		return 0, 0, "", errors.New("invalid token, token must contain three values")
	}
	uID, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		appLogger.Warnf(ctx, "Failed to parse unit_id, id=%s, err=%+v", parts[0], err)
		return 0, 0, "", err
	}
	aID, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		appLogger.Warnf(ctx, "Failed to parse access_token_id, id=%s, err=%+v", parts[1], err)
		return 0, 0, "", err
	}
	return uint32(uID), uint32(aID), parts[2], nil
}

func decodeOrgAccessToken(ctx context.Context, accessToken string) (unitID, accessTokenID uint32, plainText string, err error) {
	if !isOrgAccessToken(accessToken) {
		appLogger.Warnf(ctx, "Invalid organization token prefix")
		return 0, 0, "", errors.New("invalid organization token")
	}
	return decodeAccessToken(ctx, strings.TrimPrefix(accessToken, orgTokenPrefix))
}

func isOrgAccessToken(token string) bool {
	return strings.HasPrefix(token, orgTokenPrefix)
}
