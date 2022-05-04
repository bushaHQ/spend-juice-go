package juice

import "time"

type RegisterAccountData struct {
	FloatCurrencies    []string `json:"float_currencies"`
	BusinessAddress    string   `json:"business_address"`
	BusinessName       string   `json:"business_name"`
	Chain              string   `json:"chain"`
	ContactNumber      string   `json:"contact_number"`
	Country            string   `json:"country"`
	Domain             string   `json:"domain"`
	Email              string   `json:"email"`
	FirstName          string   `json:"first_name"`
	LastName           string   `json:"last_name"`
	Password           string   `json:"password"`
	RegistrationNumber string   `json:"registration_number"`
	WebhookUrl         string   `json:"webhook_url"`
}

type UpdateAccountData struct {
	WebhookUrl      string `json:"webhook_url"`
	BusinessAddress string `json:"business_address"`
	Domain          string `json:"domain"`
}

type RegisterUserData struct {
	Address     UserAddress `json:"address"`
	Email       string      `json:"email"`
	FirstName   string      `json:"first_name"`
	IdNumber    string      `json:"id_number"`
	IdType      string      `json:"id_type"`
	LastName    string      `json:"last_name"`
	PhoneNumber string      `json:"phone_number"`
	UserPhoto   string      `json:"user_photo,omitempty"`
}

type Param struct {
	Limit int
	Page  int
}

type CreateCardData struct {
	DesignType       string `json:"design_type"`
	SingleUse        bool   `json:"single_use"`
	Source           string `json:"source"`
	CardIntegratorId string `json:"card_integrator_id"`
	Currency         string `json:"currency"`
	JuiceUserId      string `json:"juice_user_id"`
	Validity         int    `json:"validity"`
}

type PaymentData struct {
	Source string `json:"source"`
	Amount int    `json:"amount"`
	CardId string `json:"card_id"`
}

type MockTransactionData struct {
	Amount int    `json:"amount"`
	Type   string `json:"type"`
}

type AccountResp struct {
	Data Account `json:"data"`
}

type TopUpFloatData struct {
	Amount int `json:"amount"`
}

type UsersResp struct {
	Page       int    `json:"page"`
	Total      int    `json:"total"`
	TotalPages int    `json:"total_pages"`
	Data       []User `json:"data"`
}

type UserResp struct {
	Data User `json:"data"`
}

type CardResp struct {
	Balance    int       `json:"balance"`
	CardNumber string    `json:"card_number"`
	CardType   string    `json:"card_type"`
	Cvv2       string    `json:"cvv2"`
	Expiry     time.Time `json:"expiry"`
	Id         string    `json:"id"`
	SingleUse  bool      `json:"single_use"`
	Status     string    `json:"status"`
	Valid      string    `json:"valid"`
}

type CreateCardResp struct {
	Data Card `json:"data"`
}

type TransactionResp struct {
	Data    Transaction `json:"data"`
	Message string      `json:"message"`
}

type TransactionsResp struct {
	Data     []Transaction `json:"data"`
	Message  string        `json:"message"`
	NextPage interface{}   `json:"next_page"`
}

type BalanceResp struct {
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
	Id       string `json:"id"`
}

type Resp struct {
	Message string `json:"message"`
}
