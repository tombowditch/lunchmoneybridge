package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tjvr/go-monzo"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func CheckMonzoAuthComplete() {
	ctx := context.Background()

	rc, err := client.Bucket(os.Getenv("BUCKET_NAME")).Object(TOKENFILE).NewReader(ctx)
	if err != nil {
		fmt.Println("[monzo] [auth] looks like you need to auth...")
		fmt.Println("[monzo] [auth] Please navigate to http://127.0.0.1:45679/auth")
		StartMonzoAuthWebserver()
	}
	defer rc.Close()
}

func StartMonzoAuthWebserver() {
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("redirecting you to auth page...")

		http.Redirect(w, r, "https://auth.monzo.com/?client_id="+os.Getenv("MONZO_CLIENT_ID")+"&redirect_uri=http://127.0.0.1:45679/auth/return&response_type=code&state=QblZvk", 302)
	})

	http.HandleFunc("/auth/return", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != "QblZvk" {
			w.Write([]byte("invalid state"))
			return
		}

		monzoCode := r.URL.Query().Get("code")

		hclient := &http.Client{}
		data := url.Values{}
		data.Set("grant_type", "authorization_code")
		data.Set("client_id", os.Getenv("MONZO_CLIENT_ID"))
		data.Set("client_secret", os.Getenv("MONZO_CLIENT_SECRET"))
		data.Set("redirect_uri", "http://127.0.0.1:45679/auth/return")
		data.Set("code", monzoCode)

		request, err := http.NewRequest("POST", "https://api.monzo.com/oauth2/token", strings.NewReader(data.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		response, err := hclient.Do(request)

		if err != nil {
			fmt.Println("could not get auth code")
			panic(err)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		if response.Status != "200 OK" {
			panic("auth code return is not '200 OK' (" + response.Status + ")")
		}

		token := MonzoToken{}
		err = json.Unmarshal(body, &token)
		if err != nil {
			panic(err)
		}

		xt, err := json.Marshal(token)
		if err != nil {
			panic(err)
		}

		ctx := context.Background()
		wc := client.Bucket(os.Getenv("BUCKET_NAME")).Object(TOKENFILE).NewWriter(ctx)

		if _, err = io.Copy(wc, strings.NewReader(string(xt))); err != nil {
			w.Write([]byte("failed"))
			return
		}

		if err := wc.Close(); err != nil {
			w.Write([]byte("failed"))
			return
		}

		w.Write([]byte("Verification successful. Before running cardpot again, please go to your Monzo app and allow access"))

		go func() {
			fmt.Println("written monzo authentication data")
			fmt.Println("")
			fmt.Println("##################################################")
			fmt.Println(" ")
			fmt.Println("PLEASE NOTE")
			fmt.Println("Before running cardpot again, check your phone")
			fmt.Println("There should be a notification from Monzo - please ALLOW access")
			fmt.Println(" ")
			fmt.Println("##################################################")

			os.Exit(0)
		}()
	})

	fmt.Println("Listening")

	http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), nil)

}

func GetMonzoTokenFromFile() (*MonzoToken, error) {
	ctx := context.Background()
	rc, err := client.Bucket(os.Getenv("BUCKET_NAME")).Object(TOKENFILE).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	if rc == nil {
		return nil, errors.New("no file found/rc nil")
	}

	tok := &MonzoToken{}
	err = json.NewDecoder(rc).Decode(tok)
	if err != nil {
		return nil, err
	}

	return tok, nil
}

func RenewMonzoToken() {
	tok, err := GetMonzoTokenFromFile()
	if err != nil {
		panic("could not renew: " + err.Error())
	}

	hclient := &http.Client{}
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", os.Getenv("MONZO_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("MONZO_CLIENT_SECRET"))
	data.Set("refresh_token", tok.RefreshToken)

	request, err := http.NewRequest("POST", "https://api.monzo.com/oauth2/token", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := hclient.Do(request)

	if err != nil {
		fmt.Println("could not get auth code")
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	if response.Status != "200 OK" {
		panic("auth code return is not '200 OK' (" + response.Status + ")")
	}

	token := MonzoToken{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		panic(err)
	}

	xt, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	wc := client.Bucket(os.Getenv("BUCKET_NAME")).Object(TOKENFILE).NewWriter(ctx)
	if _, err = io.Copy(wc, strings.NewReader(string(xt))); err != nil {
		panic(err)
		return
	}
	if err := wc.Close(); err != nil {
		panic(err)
		return
	}

	fmt.Println("renewed & written")

	newTok, err := GetMonzoAccessToken()
	if err != nil {
		panic(err)
	}

	fmt.Println("updated client object to use new token")

	cl = monzo.Client{
		BaseURL:     "https://api.monzo.com",
		AccessToken: newTok,
	}

}

func GetMonzoAccessToken() (string, error) {
	tok, err := GetMonzoTokenFromFile()
	if err != nil {
		return "", err
	}

	return tok.AccessToken, nil
}

func EnsureAccessTokenIsValid() {
	tok, err := GetMonzoTokenFromFile()
	if err != nil {
		panic(err)
	}

	hclient := &http.Client{}

	request, err := http.NewRequest("GET", "https://api.monzo.com/ping/whoami", nil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Bearer "+tok.AccessToken)

	response, err := hclient.Do(request)

	if err != nil {
		fmt.Println("could not get auth code")
		panic(err)
	}

	if response.Status != "200 OK" {
		fmt.Println("/ping/whoami failed: renewing token")
		RenewMonzoToken()
	} else {
		fmt.Println("/ping/whoami success: continuing")
	}

}
