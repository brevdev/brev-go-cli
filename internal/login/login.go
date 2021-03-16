package login

import (
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/brevdev/brev-go-cli/internal/config"
)

const (
	COTTER_ENDPOINT         = "https://js.cotter.app/app"
	COTTER_BACKEND_ENDPOINT = "https://www.cotter.app/api/v0"
	LOCAL_PORT              = "8395"
	LOCAL_ENDPOINT          = "http://localhost:" + LOCAL_PORT
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

func authenticateWithCotter() error {
	cotterCodeVerifier := generateCodeVerifier()

	cotterURL, err := buildCotterAuthURL(cotterCodeVerifier)
	if err != nil {
		return err
	}

	fmt.Println(cotterURL)

	err = openInDefaultBrowser(cotterURL)
	if err != nil {
		return err
	}

	token, err := captureCotterToken(cotterCodeVerifier)
	if err != nil {
		return err
	}

	// TODO: write to file
	tokenString, _ := json.Marshal(token)
	fmt.Printf(string(tokenString))

	return nil
}

func buildCotterAuthURL(code_verifier string) (string, error) {
	state := generateStateValue()
	code_challenge := generateCodeChallenge(code_verifier)

	cotterRequest, err := http.NewRequest("GET", COTTER_ENDPOINT, nil)
	if err != nil {
		return "", err
	}

	q := cotterRequest.URL.Query()
	q.Add("api_key", getCotterAPIKey())
	q.Add("redirect_url", LOCAL_ENDPOINT)
	q.Add("state", state)
	q.Add("code_challenge", code_challenge)
	q.Add("type", "EMAIL")
	q.Add("auth_method", "MAGIC_LINK")

	cotterRequest.URL.RawQuery = q.Encode()
	return cotterRequest.URL.String(), nil
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

func requestCotterToken(code string, challenge_id string, code_verifier string) (*CotterOauthToken, error) {
	challenge_id_int, err := strconv.Atoi(challenge_id)
	if err != nil {
		return nil, err
	}

	requestBody, err := json.Marshal(cotterTokenRequestPayload{
		CodeVerifier:      code_verifier,
		AuthorizationCode: code,
		ChallengeId:       challenge_id_int,
		RedirectURL:       LOCAL_ENDPOINT,
	})

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", COTTER_BACKEND_ENDPOINT+"/verify/get_identity?oauth_token=true", bytes.NewBuffer(requestBody))
	request.Header.Set("API_KEY_ID", getCotterAPIKey())
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	var tokenResponse cotterTokenResponseBody
	err = json.NewDecoder(response.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &tokenResponse.OauthToken, nil
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
