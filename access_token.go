package main

import (
	"crypto/rand"
	"encoding/base64"
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

// encodeAccessToken is encoding AccessToken. Format: urlEncode({project_id}@{access_token_id}@{plain_text})
func encodeAccessToken(projectID, accessTokenID uint32, plainText string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprint(projectID) + "@" + fmt.Sprint(accessTokenID) + "@" + plainText))
}

// decodeAccessToken is decoding AccessToken, and return access_token_id, plain_text. Format: urlEncode({project_id}@{access_token_id}@{plain_text})
func decodeAccessToken(accessToken string) (projectID, accessTokenID uint32, plainText string) {
	tokenBody, err := base64.RawURLEncoding.DecodeString(accessToken)
	if err != nil {
		appLogger.Warnf("Failed to decode access token, err=%+v", err)
		return 0, 0, ""
	}
	parts := strings.SplitN(string(tokenBody), "@", 3)
	if len(parts) != 3 {
		return 0, 0, ""
	}
	pID, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		appLogger.Warnf("Failed to parse project_id, id=%s, err=%+v", parts[0], err)
		return 0, 0, ""
	}
	aID, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		appLogger.Warnf("Failed to parse access_token_id, id=%s, err=%+v", parts[0], err)
		return 0, 0, ""
	}
	return uint32(pID), uint32(aID), parts[2]
}
