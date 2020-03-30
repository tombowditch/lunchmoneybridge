package main

import "time"

type WebhookContent struct {
	Type string           `json:"type"`
	Data MonzoTransaction `json:"data"`
}

type MonzoTransaction struct {
	ID          string    `json:"id"`
	Created     time.Time `json:"created"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	Fees        struct {
	} `json:"fees"`
	Currency string `json:"currency"`
	Merchant struct {
		ID       string    `json:"id"`
		GroupID  string    `json:"group_id"`
		Created  time.Time `json:"created"`
		Name     string    `json:"name"`
		Logo     string    `json:"logo"`
		Emoji    string    `json:"emoji"`
		Category string    `json:"category"`
		Online   bool      `json:"online"`
		Atm      bool      `json:"atm"`
		Address  struct {
			ShortFormatted string  `json:"short_formatted"`
			Formatted      string  `json:"formatted"`
			Address        string  `json:"address"`
			City           string  `json:"city"`
			Region         string  `json:"region"`
			Country        string  `json:"country"`
			Postcode       string  `json:"postcode"`
			Latitude       float64 `json:"latitude"`
			Longitude      float64 `json:"longitude"`
			ZoomLevel      int     `json:"zoom_level"`
			Approximate    bool    `json:"approximate"`
		} `json:"address"`
		Updated  time.Time `json:"updated"`
		Metadata struct {
			CreatedForMerchant     string `json:"created_for_merchant"`
			CreatedForTransaction  string `json:"created_for_transaction"`
			EnrichedFromSettlement string `json:"enriched_from_settlement"`
			FoursquareCategory     string `json:"foursquare_category"`
			FoursquareCategoryIcon string `json:"foursquare_category_icon"`
			FoursquareID           string `json:"foursquare_id"`
			FoursquareWebsite      string `json:"foursquare_website"`
			GooglePlacesIcon       string `json:"google_places_icon"`
			GooglePlacesID         string `json:"google_places_id"`
			GooglePlacesName       string `json:"google_places_name"`
			SuggestedName          string `json:"suggested_name"`
			SuggestedTags          string `json:"suggested_tags"`
			TwitterID              string `json:"twitter_id"`
			Website                string `json:"website"`
		} `json:"metadata"`
		DisableFeedback bool `json:"disable_feedback"`
	} `json:"merchant"`
	Notes    string `json:"notes"`
	Metadata struct {
		LedgerInsertionID       string `json:"ledger_insertion_id"`
		MastercardApprovalType  string `json:"mastercard_approval_type"`
		MastercardAuthMessageID string `json:"mastercard_auth_message_id"`
		MastercardCardID        string `json:"mastercard_card_id"`
		MastercardLifecycleID   string `json:"mastercard_lifecycle_id"`
		Mcc                     string `json:"mcc"`
		ExternalID              string `json:"external_id"`
		Trigger                 string `json:"trigger"`
		PotID                   string `json:"pot_id"`
	} `json:"metadata"`
	Labels         interface{} `json:"labels"`
	AccountBalance int         `json:"account_balance"`
	Attachments    interface{} `json:"attachments"`
	International  interface{} `json:"international"`
	Category       string      `json:"category"`
	Categories     interface{} `json:"categories"`
	IsLoad         bool        `json:"is_load"`
	Settled        string      `json:"settled"`
	LocalAmount    int         `json:"local_amount"`
	LocalCurrency  string      `json:"local_currency"`
	Updated        time.Time   `json:"updated"`
	AccountID      string      `json:"account_id"`
	UserID         string      `json:"user_id"`
	Counterparty   struct {
		Name string `json:"name"`
	} `json:"counterparty"`
	Scheme                     string `json:"scheme"`
	DedupeID                   string `json:"dedupe_id"`
	Originator                 bool   `json:"originator"`
	IncludeInSpending          bool   `json:"include_in_spending"`
	CanBeExcludedFromBreakdown bool   `json:"can_be_excluded_from_breakdown"`
	CanBeMadeSubscription      bool   `json:"can_be_made_subscription"`
	CanSplitTheBill            bool   `json:"can_split_the_bill"`
	CanAddToTab                bool   `json:"can_add_to_tab"`
	AmountIsPending            bool   `json:"amount_is_pending"`
}

type MonzoToken struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Type         string `json:"token_type"`
	UserID       string `json:"user_id"`
}

type LunchMoneyTransactionInsert struct {
	Transactions      []LunchMoneyTransaction `json:"transactions"`
	ApplyRules        bool                    `json:"apply_rules"` // set to true
	CheckForRecurring bool                    `json:"check_for_recurring"`
	DebitAsNegative   bool                    `json:"debit_as_negative"` // if true, will assume negative amoutn values denote expenses and positive amoutn values denote credits (set to true if just passing thru monzo amount, payments are negative)
}

type LunchMoneyTransaction struct {
	Date        string  `json:"date"`
	CategoryID  int     `json:"category_id"`
	Payee       string  `json:"payee"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	AssetID     int     `json:"asset_id"`
	RecurringID int     `json:"recurring_id"`
	Notes       string  `json:"notes"`
	Status      string  `json:"status"`
	ExternalID  string  `json:"external_id"`
}
