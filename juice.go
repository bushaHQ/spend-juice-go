package juice

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Account struct {
	BusinessAddress    string      `json:"business_address"`
	BusinessName       string      `json:"business_name"`
	Chain              string      `json:"chain"`
	ContactNumber      string      `json:"contact_number"`
	Country            string      `json:"country"`
	Domain             string      `json:"domain"`
	Email              string      `json:"email"`
	FirstName          string      `json:"first_name"`
	FloatCurrencies    []string    `json:"float_currencies"`
	Id                 string      `json:"id"`
	LastName           string      `json:"last_name"`
	RegistrationNumber string      `json:"registration_number"`
	UsdcAddress        UsdcAddress `json:"usdc_address"`
}

type UsdcAddress struct {
	Address  string `json:"address"`
	Chain    string `json:"chain"`
	Currency string `json:"currency"`
}

type User struct {
	Address          UserAddress `json:"address"`
	Archived         bool        `json:"archived"`
	CardIntegratorId string      `json:"card_integrator_id"`
	Email            string      `json:"email"`
	FirstName        string      `json:"first_name"`
	Id               string      `json:"id"`
	IdNumber         string      `json:"id_number"`
	IdType           string      `json:"id_type"`
	LastName         string      `json:"last_name"`
	PhoneNumber      string      `json:"phone_number"`
	Verified         bool        `json:"verified"`
}

type UserAddress struct {
	City    string      `json:"city"`
	Country string      `json:"country"`
	Line1   string      `json:"line1"`
	Line2   interface{} `json:"line2"`
	State   interface{} `json:"state"`
	ZipCode string      `json:"zip_code"`
}

type Card struct {
	Balance    int       `json:"balance"`
	BusinessId string    `json:"business_id"`
	CardName   string    `json:"card_name"`
	CardNumber string    `json:"card_number"`
	CardType   string    `json:"card_type"`
	Currency   string    `json:"currency"`
	Cvv2       string    `json:"cvv2"`
	DesignType string    `json:"design_type"`
	Expiry     time.Time `json:"expiry"`
	Id         string    `json:"id"`
	Provider   string    `json:"provider"`
	SingleUse  bool      `json:"single_use"`
	Status     string    `json:"status"`
	UserId     string    `json:"user_id"`
	Valid      string    `json:"valid"`
}

type Transaction struct {
	Amount            int         `json:"amount"`
	CardBalanceAfter  int         `json:"card_balance_after"`
	CardBalanceBefore int         `json:"card_balance_before"`
	ConversionRate    int         `json:"conversion_rate"`
	CreatedAt         time.Time   `json:"created_at"`
	CreditCurrency    interface{} `json:"credit_currency"`
	CreditId          interface{} `json:"credit_id"`
	DebitCurrency     interface{} `json:"debit_currency"`
	DebitId           interface{} `json:"debit_id"`
	Id                string      `json:"id"`
	Narrative         interface{} `json:"narrative"`
	Type              string      `json:"type"`
}

// RegisterAccount creates a card integrator account
func (cl *Client) RegisterAccount(data RegisterAccountData) (AccountResp, error) {
	var res AccountResp
	err := cl.post("/card-integrators/register-integrator", data, &res)
	return res, err
}

// UpdateAccount updates the card integrator account
func (cl *Client) UpdateAccount(webhook, businessAddress, domain string) (AccountResp, error) {
	var res AccountResp
	err := cl.patch("/card-integrators/update", &UpdateAccountData{WebhookUrl: webhook, BusinessAddress: businessAddress, Domain: domain}, &res)
	return res, err
}

// RegisterUser creates an account for user requesting a card
func (cl *Client) RegisterUser(data RegisterUserData, accountId string) (UserResp, error) {
	var res UserResp
	err := cl.post(fmt.Sprintf("/card-integrators/%s/register-user", accountId), data, &res)
	return res, err
}

// ListUsers gets list of card users attached to an account
func (cl *Client) ListUsers(limit, page int) (UsersResp, error) {
	var res UsersResp
	err := cl.get("/card-integrators/card-users", Param{Limit: limit, Page: page}, &res)
	return res, err
}

// CreateCard creates a card for a user
func (cl *Client) CreateCard(data CreateCardData) (CreateCardResp, error) {
	var res CreateCardResp
	err := cl.post("/cards/create-virtual-card", data, &res)
	return res, err
}

// ListCards gets a list of cards of a user
func (cl *Client) ListCards(limit, page int, userId string) ([]CardResp, error) {
	var res []CardResp
	err := cl.get(fmt.Sprintf("/cards?juice_user_id=%s&limit=%d&page=%d", userId, limit, page), nil, &res)
	return res, err
}

// GetCard gets a particular card
func (cl *Client) GetCard(cardId string) (CardResp, error) {
	var res CardResp
	err := cl.get(fmt.Sprintf("/cards/%s", cardId), nil, &res)
	return res, err
}

// CreditCard top-up a card for a user
func (cl *Client) CreditCard(data PaymentData) (CardResp, error) {
	var res CardResp
	err := cl.patch("/cards/credit/balance", data, &res)
	return res, err
}

// DebitCard debits a card for a user
func (cl *Client) DebitCard(data PaymentData) (CardResp, error) {
	var res CardResp
	err := cl.patch("/cards/debit/balance", data, &res)
	return res, err
}

// FreezeCard freezes a card for a user
func (cl *Client) FreezeCard(cardId string) (CardResp, error) {
	var res CardResp
	err := cl.patch(fmt.Sprintf("/cards/%s/freeze", cardId), nil, &res)
	return res, err
}

// UnfreezeCard unfreezes a card for a user
func (cl *Client) UnfreezeCard(cardId string) (CardResp, error) {
	var res CardResp
	err := cl.patch(fmt.Sprintf("/cards/%s/unfreeze", cardId), nil, &res)
	return res, err
}

// ListTransactions gets paginated transactions for the given card
func (cl *Client) ListTransactions(cardId string, param Param) (TransactionsResp, error) {
	var res TransactionsResp
	err := cl.get(fmt.Sprintf("/cards/%s/transactions", cardId), param, &res)
	return res, err
}

// GetTransaction gets a particular transaction
func (cl *Client) GetTransaction(trxId string) (TransactionResp, error) {
	var res TransactionResp
	err := cl.get(fmt.Sprintf("/cards/transaction/%s", trxId), nil, &res)
	return res, err
}

func (cl Client) Health() (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/health/live", cl.baseURL))
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), err
}
