package auth

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/brevdev/brev-go-cli/internal/config"
	"github.com/brevdev/brev-go-cli/internal/files"
	"github.com/brevdev/brev-go-cli/internal/requests"
)

const (
	COTTER_ENDPOINT         = "https://js.cotter.app/app"
	COTTER_BACKEND_ENDPOINT = "https://www.cotter.app/api/v0"
	COTTER_TOKEN_ENDPOINT   = "https://www.cotter.app/api/v0/token"
	COTTER_JWKS_ENDPOINT    = "https://www.cotter.app/api/v0/token/jwks"
	LOCAL_PORT              = "8395"
	LOCAL_ENDPOINT          = "http://localhost:" + LOCAL_PORT

	BREV_CREDENTIALS_FILE = "credentials.json"
)

type cotterTokenRequestPayload struct {
	CodeVerifier      string `json:"code_verifier"`
	AuthorizationCode string `json:"authorization_code"`
	ChallengeId       int    `json:"challenge_id"`
	RedirectURL       string `json:"redirect_url"`
}

type cotterTokenResponseBody struct {
	OauthToken CotterOauthToken `json:"oauth_token"`
}

type CotterOauthToken struct {
	AccessToken  string `json:"access_token"`
	AuthMethod   string `json:"auth_method"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

//go:embed success.html
var successHTML embed.FS

// Login performs a full round trip to Cotter:
//   1. Generate a code verifier
//   2. Open the Brev+Cotter auth URL in the default browser
//   3. Start a local web server awaiting a redirect to localhost
//   4. Capture the Cotter token upon redirect
//   5. Write the Cotter token to a file in the hidden brev directory
func Login() error {
	cotterCodeVerifier := generateCodeVerifier()

	cotterURL, err := buildCotterAuthURL(cotterCodeVerifier)
	if err != nil {
		return err
	}

	// TODO: pretty print URL?
	fmt.Println(cotterURL)

	err = openInDefaultBrowser(cotterURL)
	if err != nil {
		return err
	}

	token, err := captureCotterToken(cotterCodeVerifier)
	if err != nil {
		return err
	}

	err = writeTokenToBrevConfigFile(token)
	if err != nil {
		return err
	}

	return nil
}

// GetToken reads the previously-persisted token from the filesystem but may issue a round
// trip request to Cotter if the token is determined to have expired:
//   1. Read the Cotter token from the hidden brev directory
//   2. Determine if the token is valid
//   3. If valid, return
//   4. If invalid, issue a refresh request to Cotter
//   5. Write the refreshed Cotter token to a file in the hidden brev directory
func GetToken() (*CotterOauthToken, error) {
	token, err := getTokenFromBrevConfigFile()
	if err != nil {
		return nil, err
	}

	if token.isValid() {
		return token, nil
	}

	refreshedToken, err := refreshCotterToken(token)
	if err != nil {
		return nil, err
	}

	err = writeTokenToBrevConfigFile(refreshedToken)
	if err != nil {
		return nil, err
	}

	return refreshedToken, nil
}

func (t *CotterOauthToken) isValid() bool {
	var claims map[string]interface{}

	jwtToken, err := jwt.ParseSigned(t.AccessToken)
	if err != nil {
		return false
	}

	jwks, err := fetchCotterPublicKeySet()
	if err != nil {
		return false
	}

	err = jwtToken.Claims(jwks, &claims)
	if err != nil {
		return false
	}

	now := time.Now().Unix()
	exp := int64(claims["exp"].(float64))
	if now > exp {
		return false
	}

	return true
}

func fetchCotterPublicKeySet() (*jose.JSONWebKeySet, error) {
	request := requests.RESTRequest{
		Method:   "GET",
		Endpoint: COTTER_JWKS_ENDPOINT,
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var cotterJWKS jose.JSONWebKeySet
	err = response.DecodePayload(&cotterJWKS)
	if err != nil {
		return nil, err
	}

	return &cotterJWKS, nil
}

func buildCotterAuthURL(code_verifier string) (string, error) {
	state := generateStateValue()
	code_challenge := generateCodeChallenge(code_verifier)

	request := &requests.RESTRequest{
		Method:   "GET",
		Endpoint: COTTER_ENDPOINT,
		QueryParams: []requests.QueryParam{
			{"api_key", getCotterAPIKey()},
			{"redirect_url", LOCAL_ENDPOINT},
			{"state", state},
			{"code_challenge", code_challenge},
			{"type", "EMAIL"},
			{"auth_method", "MAGIC_LINK"},
		},
	}

	httpRequest, err := request.BuildHTTPRequest()
	if err != nil {
		return "", err
	}

	return httpRequest.URL.String(), nil
}

func openInDefaultBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		return errors.New(fmt.Sprintf("Unsupported runtime: %s", runtime.GOOS))
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func captureCotterToken(code_verifier string) (*CotterOauthToken, error) {
	m := http.NewServeMux()
	s := http.Server{Addr: ":" + LOCAL_PORT, Handler: m}

	var token *CotterOauthToken
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		errorParm := q.Get("error")
		if errorParm != "" {
			fmt.Println(errorParm)
			// panic!
		}

		code := q.Get("code")
		challenge_id := q.Get("challenge_id")

		cotterToken, err := requestCotterToken(code, challenge_id, code_verifier)
		if err != nil {
			// panic!
		}

		token = cotterToken

		content, _ := successHTML.ReadFile("success.html")
		w.Write(content)

		err = s.Shutdown(context.Background())
		if err != nil {
			// panic!
		}
	})

	// blocks
	s.ListenAndServe()
	return token, nil
}

func writeTokenToBrevConfigFile(token *CotterOauthToken) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	brevCredentialsFile := home + "/" + config.GetBrevDirectory() + "/" + BREV_CREDENTIALS_FILE

	err = files.OverwriteJSON(brevCredentialsFile, token)
	if err != nil {
		return err
	}

	return nil
}

func getTokenFromBrevConfigFile() (*CotterOauthToken, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	brevCredentialsFile := home + "/" + config.GetBrevDirectory() + "/" + BREV_CREDENTIALS_FILE

	var token CotterOauthToken
	err = files.ReadJSON(brevCredentialsFile, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func requestCotterToken(code string, challenge_id string, code_verifier string) (*CotterOauthToken, error) {
	challenge_id_int, err := strconv.Atoi(challenge_id)
	if err != nil {
		return nil, err
	}

	request := &requests.RESTRequest{
		Method:   "POST",
		Endpoint: COTTER_BACKEND_ENDPOINT + "/verify/get_identity",
		Headers: []requests.Header{
			{"API_KEY_ID", getCotterAPIKey()},
			{"Content-Type", "application/json"},
		},
		QueryParams: []requests.QueryParam{
			{"oauth_token", "true"},
		},
		Payload: cotterTokenRequestPayload{
			CodeVerifier:      code_verifier,
			AuthorizationCode: code,
			ChallengeId:       challenge_id_int,
			RedirectURL:       LOCAL_ENDPOINT,
		},
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var tokenResponse cotterTokenResponseBody
	response.DecodePayload(&tokenResponse)

	return &tokenResponse.OauthToken, nil
}

func refreshCotterToken(oldToken *CotterOauthToken) (*CotterOauthToken, error) {
	request := &requests.RESTRequest{
		Method:   "POST",
		Endpoint: COTTER_TOKEN_ENDPOINT + "/" + getCotterAPIKey(),
		Headers: []requests.Header{
			{"API_KEY_ID", getCotterAPIKey()},
			{"Content-Type", "application/json"},
		},
		Payload: map[string]string{
			"grant_type":    "refresh_token",
			"refresh_token": oldToken.RefreshToken,
		},
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var token CotterOauthToken
	response.DecodePayload(&token)
	return &token, nil
}

func generateStateValue() string {
	return randomAlphabetical(10)
}

func generateCodeVerifier() string {
	code_verifier_bytes := randomBytes(32)
	code_verifier_raw := base64.URLEncoding.EncodeToString(code_verifier_bytes)
	code_verifier := strings.TrimRight(code_verifier_raw, "=")

	return code_verifier
}

func generateCodeChallenge(code_verifier string) string {
	sha256Hasher := sha256.New()
	sha256Hasher.Write([]byte(code_verifier))
	challenge_bytes := sha256Hasher.Sum(nil)
	challenge_raw := base64.URLEncoding.EncodeToString(challenge_bytes)
	challenge := strings.TrimRight(challenge_raw, "=")

	return challenge
}

func getCotterAPIKey() string {
	return config.GetCotterAPIKey()
}
