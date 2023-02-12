package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	constants "github.com/ostd2000/Auth/constants"
)


type Claims struct {
	Exp    int64
	Type   string
	UserID string
}

type Header struct {
	alg string
	typ string
}


// Encodes given byte array to base64url string.
func Base64UrlEncode(data []byte) string {
  result :=  base64.StdEncoding.EncodeToString(data)

	// 62nd char of encoding.
	result = strings.Replace(result, "+", "-", -1)

	// 63rd char of encoding.
	result = strings.Replace(result, "/", "_", -1)

	// Removes any trailing '='s
	result = strings.Replace(result, "=", "", -1)

	return result
}


// Decodes base64url string to byte array.
func Base64UrlDecode(data string) ([]byte, error) {
	// 62 char of encoding.
	data = strings.Replace(data, "-", "+", -1)

	// 63 char of encoding.
  data = strings.Replace(data, "_", "/", -1)

	// Pad with trailing '='s.
	switch(len(data) % 4) {
		// No padding.
	  case 0:

		// 2 pad chars.
  	case 2: data += "=="

		// 1 pad char.
	  case 3: data += "="
	}

	return base64.StdEncoding.DecodeString(data)
}


func TokenSignature(msg []byte, key []byte) string {
  mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	
	expectedMAC := mac.Sum(nil)	
  
	return string(expectedMAC)
}


var secretKey string = os.Getenv("SECRET_KEY")


func GenerateAccessToken(userID string) (signedAccessToken string, err error) {
    header := &Header{
		alg: constants.JWT_SIGNING_ALGORITHM,
	    typ: "JWT",	
	}

	jsonHeader, err := json.Marshal(header)

	accessClaims := &Claims{
		// Registered claims
		Exp: time.Now().UTC().Add(constants.ACCESS_TOKEN_EXP_TIME).Unix(),

		// Public claims
		Type: "access_token",
		UserID: userID,
	}

	// "jsonAccessClaims" will be a byte array.
	jsonAccessClaims, err := json.Marshal(accessClaims)

	if err != nil {
        log.Panic(err)			

		return 
	}

	msg := Base64UrlEncode(jsonHeader) + "." + Base64UrlEncode(jsonAccessClaims)
  
	signedAccessToken = TokenSignature([]byte(msg), []byte(secretKey))

	return
}


func GenerateRefreshToken(userID string) (signedRefreshToken string, err error) {
  header := &Header{
		alg: constants.JWT_SIGNING_ALGORITHM,
		typ: "JWT",
	}

	jsonHeader, err := json.Marshal(header) 

  refreshClaims := &Claims{
		// Registerd claims
		Exp: time.Now().UTC().Add(constants.REFRESH_TOKEN_EXP_TIME).Unix(),

		// Public cliams
		Type: "refresh_token",
		UserID: userID,
	}

	// "jsonRefreshCaims" will be a byte array.
	jsonRefreshClaims, err := json.Marshal(refreshClaims)

	if err != nil {
		log.Panic(err)

		return
	}

	msg := Base64UrlEncode(jsonHeader) + "." + Base64UrlEncode(jsonRefreshClaims)

	signedRefreshToken = TokenSignature([]byte(msg), []byte(secretKey))

	return
} 