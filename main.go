package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/tjvr/go-monzo"
)

var client *storage.Client
var cl monzo.Client

const TOKENFILE = "monzo.token"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ctx := context.Background()

	client, _ = storage.NewClient(ctx)

	CheckMonzoAuthComplete()
	EnsureAccessTokenIsValid()

	tok, err := GetMonzoAccessToken()
	if err != nil {
		panic(err)
	}

	cl = monzo.Client{
		BaseURL:     "https://api.monzo.com",
		AccessToken: tok,
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Write([]byte("unsupported method"))
			return
		}

		EnsureAccessTokenIsValid()

		decoder := json.NewDecoder(r.Body)
		var wh WebhookContent
		err := decoder.Decode(&wh)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid json"))
			return
		}

		if wh.Data.Amount == 0 || wh.Data.Amount == -0 {
			w.Write([]byte("active card check, ignoring"))
			return
		}

		InsertTransactionIntoLunchmoney(wh.Data, false)

		w.Write([]byte("ok"))

	})

	http.HandleFunc("/transactionlistinsert", func(w http.ResponseWriter, r *http.Request) {
		ts, err := cl.Transactions(os.Getenv("ACCOUNT_ID"), true)
		if err != nil {
			panic(err)
		}

		keys, _ := r.URL.Query()["lminsert"]

		for _, t := range ts {
			x, err := json.Marshal(t)

			if err != nil {
				fmt.Println("marshal err")
				fmt.Println(err)
				return
			}

			var trn MonzoTransaction

			err = json.Unmarshal(x, &trn)

			if err == nil {
				if trn.Created.Year() == 2020 {
					if keys[0] == "yess" {
						fmt.Println("insert into LM!")
						//InsertTransactionIntoLunchmoney(trn, false)
					} else {
						fmt.Println("skipping insert")
					}
				} else {
					fmt.Println("skipping, not 2020")
				}
			} else {
				fmt.Println("err occurred")
				fmt.Println(err)
				fmt.Println(string(x))
			}

		}

		w.Write([]byte("done"))
	})

	fmt.Println("Listening")

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func InsertTransactionIntoLunchmoney(transaction MonzoTransaction, dupeCheck bool) {
	if dupeCheck {
		panic("dupeCheck not impl")
		return
	}

	url := "https://dev.lunchmoney.app/v1/transactions"

	tname := ""

	if transaction.Merchant.Name != "" {
		tname = transaction.Merchant.Name
	} else if transaction.Counterparty.Name != "" {
		tname = transaction.Counterparty.Name
	} else if strings.Contains(transaction.Metadata.ExternalID, "CoinJarRule") {
		tname = "Coin Jar"
	} else if transaction.Metadata.Trigger == "ifttt" {
		tname = "Pot (IFTTT)"
	} else if transaction.Metadata.PotID != "" {
		EnsureAccessTokenIsValid()

		pot, err := cl.Pot(transaction.Metadata.PotID)
		if err != nil {
			tname = "Pot - unknown"
		} else {
			tname = "Pot - " + pot.Name
		}
	} else if transaction.Description == "Monzo Plus subscription fee" {
		tname = "Monzo Plus Fee"
	} else if strings.Contains(transaction.Description, "overdraft fees") {
		tname = "Overdraft fees"
	} else if strings.Contains(transaction.Description, "Interest for") {
		tname = "Interest"
	} else {
		tname = transaction.Description
	}

	lmt := LunchMoneyTransaction{
		Currency:   "gbp",
		Payee:      tname,
		Amount:     float64(transaction.Amount) / float64(100),
		Date:       transaction.Created.Format("2006-01-02"),
		Status:     "cleared",
		ExternalID: transaction.ID,
		AssetID:    1446,
	}

	lmj := LunchMoneyTransactionInsert{
		Transactions:      []LunchMoneyTransaction{lmt},
		ApplyRules:        true,
		CheckForRecurring: true,
		DebitAsNegative:   true,
	}

	jsonStr, _ := json.Marshal(lmj)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("LUNCHMONEY_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}