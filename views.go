package main

import (
	"net/http"
	"encoding/json"
	"regexp"
	"errors"
	"fmt"
	"log"
	"time"
)

type ValidationReponse struct {
	responseStatus	int		// Won't be marshaled
	RequestedURL	string	`json:"requested_url"`
	StatusCode		int		`json:"status_code"`
	ErrorMessage	string	`json:"error_message"`
	DebugMessage	string 	`json:"debug_message"`
}

func injectAccessControlHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")

	injectAccessControlHeaders(w)

	w.WriteHeader(code)
	w.Write(response)
}

func urlHasHypertextScheme(url string ) bool {
	if matched, err := regexp.MatchString("\\A(https?://)", url); err != nil || !matched {
		log.Printf("Requested url does not contain a protocol scheme")
		return false
	}
	return true
}

func (a *App)findUsableScheme(requestURL string) (schemedURL string, response *http.Response, err error) {

	httpsPrefixed := fmt.Sprintf("https://%s", requestURL)
	httpPrefixed := fmt.Sprintf("http://%s", requestURL)

	httpsResp, err := a.performGET(httpsPrefixed)
	if err != nil {
		log.Printf("GET to %s returned and error", httpsPrefixed)
		httpResp, err := a.performGET(httpPrefixed)
		if err != nil {
			return httpPrefixed, nil, err
		}
		return httpPrefixed, httpResp, nil
	}
	return httpsPrefixed, httpsResp, nil
}

func (a *App) performGET(url string) (r *http.Response, e error) {
	timeoutDuration, err := time.ParseDuration(fmt.Sprintf("%ds", a.requestTimeoutSec))
	if err != nil {
		timeoutDuration = time.Duration(3 * time.Second)
	}

	client := http.Client{
		Timeout: timeoutDuration,
		Transport: &http.Transport{
		},
	}


	resp, err := client.Get(url)

	if err == nil {
		log.Printf("Status code for GET to %s: %d", url, resp.StatusCode)
	}

	return resp, err
}

func (a *App) requestURL(rawReqURL string, preflightRegex string) (StatusCode int, finalURL string, err error) {
	if preflightRegex != "" {
		if matched, err := regexp.MatchString(preflightRegex, rawReqURL); err != nil || !matched {
			return -1, rawReqURL, errors.New("unable to validate URL")
		}
	}

	if urlHasHypertextScheme(rawReqURL) {
		resp, err := a.performGET(rawReqURL)
		if err != nil {
			return -1, rawReqURL, err
		}

		return resp.StatusCode, rawReqURL, nil
	}

	schemedURL, finalResp, err := a.findUsableScheme(rawReqURL)
	if err != nil {
		return -1, rawReqURL, err
	}
	return finalResp.StatusCode, schemedURL, nil
}


func (a *App) genericURLValidationView(w http.ResponseWriter, r *http.Request) {
	valURL := r.URL.Query().Get("url")

	if valURL == "" {
		respondWithJSON(w, http.StatusBadRequest, &ValidationReponse{
			StatusCode: -1,
			ErrorMessage: "Unable to validate empty URL",
			DebugMessage: "'url' query parameter must be filled",
		})
		return
	}

	urlStatus, requestedURL, err := a.requestURL(valURL, "")

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, &ValidationReponse{
			StatusCode: -1,
			RequestedURL: requestedURL,
			ErrorMessage: "An error ocurred while attempting to validate URL",
			DebugMessage: err.Error(),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, &ValidationReponse{
		StatusCode:	urlStatus,
		RequestedURL: requestedURL,
	})
}

func (a *App) facebookValidationView(w http.ResponseWriter, r *http.Request) {
	valURL := r.URL.Query().Get("url")

	if valURL == "" {
		respondWithJSON(w, http.StatusBadRequest, &ValidationReponse{
			StatusCode: -1,
			ErrorMessage: "Unable to validate empty URL",
			DebugMessage: "'url' query parameter must be filled",
		})
		return
	}

	reqStatus, requestedURL, err := a.requestURL(valURL, fbURLRegex)

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, &ValidationReponse{
			StatusCode: -1,
			RequestedURL: requestedURL,
			ErrorMessage: "An error ocurred while attempting to validate URL",
			DebugMessage: err.Error(),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, &ValidationReponse{
		StatusCode:	reqStatus,
		RequestedURL: requestedURL,
	})
}


func (a *App) linkedInValidationView(w http.ResponseWriter, r *http.Request) {
	valURL := r.URL.Query().Get("url")

	if valURL == "" {
		respondWithJSON(w, http.StatusBadRequest, &ValidationReponse{
			StatusCode: -1,
			ErrorMessage: "Unable to validate empty URL",
			DebugMessage: "'url' query parameter must be filled",
		})
		return
	}

	reqStatus, requestedURL, err := a.requestURL(valURL, linkedInURLRegex)

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, &ValidationReponse{
			StatusCode: -1,
			RequestedURL: requestedURL,
			ErrorMessage: "An error ocurred while attempting to validate URL",
			DebugMessage: err.Error(),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, &ValidationReponse{
		StatusCode:	reqStatus,
		RequestedURL: requestedURL,
	})
}

func (a *App) twitterValidationView(w http.ResponseWriter, r *http.Request) {
	testHandle := r.URL.Query().Get("handle")

	if testHandle == "" {
		respondWithJSON(w, http.StatusBadRequest, &ValidationReponse{
			StatusCode: -1,
			ErrorMessage: "Unable to validate empty handle",
			DebugMessage: "'handle' query parameter must be filled",
		})
		return
	}

	twitterTestURL := fmt.Sprintf(twitterURLFmtStr, testHandle)

	reqStatus, requestedURL, err := a.requestURL(twitterTestURL, "")

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, &ValidationReponse{
			StatusCode: -1,
			RequestedURL: requestedURL,
			ErrorMessage: "An error ocurred while attempting to validate handle",
			DebugMessage: err.Error(),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, &ValidationReponse{
		StatusCode:	reqStatus,
		RequestedURL: requestedURL,
	})
}

func (a *App) instagramValidationView(w http.ResponseWriter, r *http.Request) {
	testUsername := r.URL.Query().Get("username")

	if testUsername == "" {
		respondWithJSON(w, http.StatusBadRequest, &ValidationReponse{
			StatusCode: -1,
			ErrorMessage: "Unable to validate empty username",
			DebugMessage: "'username' query parameter must be filled",
		})
		return
	}

	instagramTestURL := fmt.Sprintf(instagramURLFmtStr, testUsername)

	reqStatus, requestedURL, err := a.requestURL(instagramTestURL, "")

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, &ValidationReponse{
			StatusCode: -1,
			RequestedURL: requestedURL,
			ErrorMessage: "An error ocurred while attempting to validate username",
			DebugMessage: err.Error(),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, &ValidationReponse{
		StatusCode:	reqStatus,
		RequestedURL: requestedURL,
	})
}
