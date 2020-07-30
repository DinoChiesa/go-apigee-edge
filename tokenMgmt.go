package apigee

import (
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"os"
	"os/user"
  "net/url"
	"net/http"
	"fmt"
	"time"
	"path"
	"strings"
)

const (
	tokenStashFile = "~/.apigee-edge-tokens"
)

type ISOTimes struct {
	IssuedAt string `json:"issued_at,omitempty"`
	Expires string `json:"expires,omitempty"`
}

type AuthToken struct {
	AccessToken *string `json:"access_token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	TokenType *string `json:"token_type,omitempty"`
	Lifetime int64 `json:"expires_in,omitempty"`
	Expires int64 `json:"expires,omitempty"`
  Scope *string `json:"scope,omitempty"`
  Id *string `json:"jti,omitempty"`
	IssuedAt int64 `json:"issued_at,omitempty"`
	ISO *ISOTimes `json:"ISO,omitempty"`
}

func resolveAnyTildes(v string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	if strings.HasPrefix(v, "~/") {
    v = path.Join(dir, v[2:])
	}	else if v == "~" {
    v = dir
	}
	return v
}

func ReadTokenStash() (map[string]*AuthToken, error) {
	filepath := resolveAnyTildes(tokenStashFile)
   _, e := os.Stat(filepath)
  if os.IsNotExist(e) {
		fmt.Printf("token stash file does not exist\n")
		return make(map[string]*AuthToken), nil
	}

	file, e := ioutil.ReadFile(filepath)
	if e != nil {
		return nil, e
	}
	var tokenStash map[string]*AuthToken
	e = json.Unmarshal(file, &tokenStash)
	if tokenStash == nil {
		tokenStash = make(map[string]*AuthToken)
	}
	return tokenStash, e
}


func tokenStashKey(c *ApigeeClient) string {
	key := c.auth.Username + "##" + c.BaseURL.String() + "##" + getLoginBaseUrl(c)
	//fmt.Printf("token stash key: %s\n", key)
  return key
}

func CurrentToken(c *ApigeeClient) (*AuthToken, error) {
  stash, e := ReadTokenStash()
	if e != nil {
		return nil, e
	}
	key := tokenStashKey(c)
	val, ok := stash[key]
	if ok {
		return val, nil
	}
	return nil, nil
}

func IsInvalidOrExpired(token *AuthToken) bool {
	if token.AccessToken == nil || token.Expires == 0 || token.IssuedAt == 0 {
		return true
	}
	nowSecondsSinceEpoch := time.Now().Unix()
	adjustmentInMilliseconds := int64(30 * 1000)
	adjustedNow := nowSecondsSinceEpoch + adjustmentInMilliseconds
	invalidOrExpired := (token.Expires < adjustedNow)
  return invalidOrExpired
}

func enhanceToken(token *AuthToken) *AuthToken {
  var iso ISOTimes
  if token.AccessToken != nil {
		parts := strings.Split(*token.AccessToken, ",")
    if len(parts) == 3 {
			var payload []byte
			payload, e := base64.RawStdEncoding.DecodeString(parts[1])
			if e == nil {
				var claims map[string]interface{}
				e = json.Unmarshal(payload, &claims)
				if e == nil {
					// The issued_at and expires_in properties on the token
					// WRAPPER are inconsistent with the actual token. So let's
					// overwrite them.
					if val, ok := claims["iat"]; ok {
						v := val.(int64)
						t := time.Unix(v, 0)
						token.IssuedAt = t.UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
						iso.IssuedAt = t.Format(time.RFC3339)
					}
					if val, ok := claims["exp"]; ok {
						v := val.(int64)
						t := time.Unix(v, 0)
						v = t.UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
						v2 := int64(claims["iat"].(int64))
						token.Lifetime = v - v2
						token.Expires = v
						iso.Expires = t.Format(time.RFC3339)
					}
				}
      }
		} else {
			// not a JWT; probably a googleapis opaque oauth token
			if token.IssuedAt != 0 {
				v := int64(token.IssuedAt * int64(time.Millisecond))
				t := time.Unix(0, v)
				iso.IssuedAt = t.Format(time.RFC3339)
				if token.Lifetime != 0 {
					t := time.Unix(0, (token.IssuedAt + token.Lifetime * 1000) * int64(time.Millisecond))
					token.Expires = t.UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
					iso.Expires = t.Format(time.RFC3339)
				}
			}
		}
  }
  token.ISO = &iso;
  return token;
}


func StashToken(c *ApigeeClient, newToken *AuthToken) (map[string]*AuthToken, error) {
  stash, e := ReadTokenStash()
	if e != nil {
		return nil, e
	}
	key := tokenStashKey(c)
	stash[key] = newToken;  // possibly overwrite an existing entry

	keptTokens := make(map[string]*AuthToken)
	for key, element := range stash {
		if !IsInvalidOrExpired(element) {
			keptTokens[key] = enhanceToken(element)
		}
	}

	json, e := json.MarshalIndent(keptTokens, "", "  ")
	if e != nil {
		return keptTokens, e
	}

	filename := resolveAnyTildes(tokenStashFile)
 	file, e := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0600)
	if e != nil {
		return keptTokens, e
	}
	defer file.Close()
	file.Write(json)
  return keptTokens, nil
}

func GetToken(c *ApigeeClient) (*AuthToken, error) {
	// read the token stash and return the stashed token if not expired
	currentToken, e := CurrentToken(c)
	if e != nil {
		return nil, e
	}
	if currentToken != nil {
		return currentToken, nil
	}
	return GetNewToken(c)
}

func getLoginBaseUrl(c *ApigeeClient) string {
	if c.LoginBaseUrl != "" {
		return c.LoginBaseUrl
	}
	return "https://login.apigee.com"
}

func GetNewToken(c *ApigeeClient) (*AuthToken, error) {
	// TODO: support GAAMBO
	// if (arg1.config) {
	//   return postGoogleapisTokenEndpoint(conn, arg1.config, cb);
	// }
	loginBaseUrl := getLoginBaseUrl(c)
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("username", c.auth.Username)
	form.Add("password", c.auth.Password)

	req, e := http.NewRequest("POST", loginBaseUrl + "/oauth/token", strings.NewReader(form.Encode()))

	// This method always tries to directly login to apigee IDP.
	// TODO: add support for 2FA and for Apigee SSO.
	//
	// else if (arg1.passcode) {
	//   // exchange passcode for token
	//   formparams = merge(formparams, { response_type: 'token', passcode: arg1.passcode });
	// }
	// else if (arg1.mfa_token) {
	//     formparams = merge(formparams, { mfa_token: arg1.mfa_token });
	// }

  if e != nil {
    return nil, e
  }

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Add("Accept", "application/json")
  req.Header.Add("Authorization", "Basic ZWRnZWNsaTplZGdlY2xpc2VjcmV0")
	var v AuthToken
	_, e = c.Do(req, &v)
  if e != nil {
    return nil, e // errors.Wrap(e, "while getting new token:")
  }
	v.IssuedAt = time.Now().Unix() * 1000
	v.Expires = v.IssuedAt + v.Lifetime*1000
	c.auth.Token = *v.AccessToken
	_, e = StashToken(c, &v)

	return &v, e
}
