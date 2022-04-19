package juice

import (
	"bytes"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	cl *Client
)

type MockHttpClient struct {
	DoFunc  func(req *http.Request) (*http.Response, error)
	Timeout time.Duration
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func init() {
	err := godotenv.Load("./.env")
	if err != nil && os.Getenv("ENV") == "" {
		panic(err)
	}
	cl = NewClient()

}

func TestClient_RegisterAccount(t *testing.T) {
	type args struct {
		payload RegisterAccountData
	}
	tests := []struct {
		name           string
		mockHttpClient MockHttpClient
		args           args
		want           AccountResp
		wantErr        bool
		errMsg         string
	}{
		{
			name: "Create a card integrator account (Success)",
			mockHttpClient: MockHttpClient{

				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/card-integrators/register-integrator" {
						t.Errorf("Expected to request '/card-integrators/register-integrator', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(
						`{
							"data": {
								"business_address": "Lekki Ikate",
								"business_name": "Algo Math",
								"chain": "ETH",
								"contact_number": "+2349099435568",
								"country": "NG",
								"domain": "https://boro.com",
								"email": "boro@gmail.com",
								"first_name": "Olusola",
								"float_currencies": [
									"USD"
								],
								"id": "27de9f46-726a-4499-aa62-27c3ed274026",
								"last_name": "Alao",
								"registration_number": "RC-546787",
								"usdc_address": {
									"address": "0xbaf2e14f27c106f0d078b397af9c7eebf1611d46",
									"chain": "ETH",
									"currency": "USD"
								}
							}
						}`)))
					return &http.Response{
						StatusCode: 201,
						Body:       responseBody,
					}, nil
				},
			},
			args: args{
				payload: RegisterAccountData{
					BusinessAddress:    "Lekki Ikate",
					BusinessName:       "Algo Math",
					Chain:              "ETH",
					ContactNumber:      "+2349099435568",
					Country:            "NG",
					Domain:             "https://boro.com",
					Email:              "boro@gmail.com",
					FirstName:          "Olusola",
					FloatCurrencies:    []string{"USD"},
					LastName:           "Alao",
					Password:           "@Password",
					RegistrationNumber: "RC-546787",
					WebhookUrl:         "https://webhook.site/7bdf91c4-6e84-4ff1-a3e1-185f138247c1",
				},
			},
			want: AccountResp{
				Data: Account{
					BusinessAddress:    "Lekki Ikate",
					BusinessName:       "Algo Math",
					Chain:              "ETH",
					ContactNumber:      "+2349099435568",
					Country:            "NG",
					Domain:             "https://boro.com",
					Email:              "boro@gmail.com",
					FirstName:          "Olusola",
					FloatCurrencies:    []string{"USD"},
					Id:                 "27de9f46-726a-4499-aa62-27c3ed274026",
					LastName:           "Alao",
					RegistrationNumber: "RC-546787",
					UsdcAddress: UsdcAddress{
						Address:  "0xbaf2e14f27c106f0d078b397af9c7eebf1611d46",
						Chain:    "ETH",
						Currency: "USD",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Create a card integrator account (Error: user already exists)",
			mockHttpClient: MockHttpClient{

				DoFunc: func(r *http.Request) (*http.Response, error) {
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(
						`{"errors": {"message": "Email or phone number already exists"}}`)))
					return &http.Response{
						StatusCode: 400,
						Body:       responseBody,
					}, nil
				},
			},
			args: args{
				payload: RegisterAccountData{
					BusinessAddress:    "1 Ademola Tokunbo Str, Lagos",
					BusinessName:       "algo math",
					Chain:              "ETH",
					ContactNumber:      "+2349034384662",
					Country:            "NG",
					Domain:             "https://olusolaa.tech",
					Email:              "email@gmail.com",
					FirstName:          "Olusola",
					FloatCurrencies:    []string{"USD"},
					LastName:           "Alao",
					Password:           "password",
					RegistrationNumber: "12345",
					WebhookUrl:         "https://webhook.site/043c6db2-5c17-4885-b769-4491ce3b0b0e",
				},
			},
			errMsg:  "email or phone number already exists",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(&tt.mockHttpClient)
			got, err := cl.RegisterAccount(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("RegisterAccount() actual error = %v, expected error %v", err, tt.errMsg)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterAccount() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestClient_UpdateAccount(t *testing.T) {

	type args struct {
		webhook     string
		businessAdd string
		domain      string
	}
	tests := []struct {
		name           string
		mockHttpClient HTTPClient
		args           args
		want           AccountResp
		wantErr        bool
		errMsg         string
	}{
		{
			name: "update integration (success)",
			mockHttpClient: &MockHttpClient{

				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/card-integrators/update" {
						t.Errorf("Expected to request '/card-integrators/update', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(`{
		"data":{"business_address":"New Road Lincoln street","business_name":"Algo Math",
		"chain":"ETH","contact_number":"+2349034384660","country":"NG","domain":"https://olusolaa.tech",
		"email":"alaoolusolae@gmail.com","first_name":"Olusola","float_currencies":["USD"],
		"id":"8de0c7a2-0004-4420-899b-f8d89c81f82b","last_name":"Alao","registration_number":"RC-546789",
		"usdc_address":{"address":"0x7af2fa93d4069655098a24dd055ff7aa51bab531","chain":"ETH","currency":"USD"
		}}} <nil> 200`)))
					return &http.Response{
						StatusCode: 200,
						Body:       responseBody,
					}, nil
				},
				Timeout: 0,
			},
			args: args{
				businessAdd: "New Road Lincoln street",
				domain:      "https://olusolaa.tech",
				webhook:     "https://webhook.site/043c6db2-5c17-4885-b769-4491ce3b0b0e",
			},
			want: AccountResp{
				Account{
					"New Road Lincoln street",
					"Algo Math",
					"ETH",
					"+2349034384660",
					"NG",
					"https://olusolaa.tech",
					"alaoolusolae@gmail.com",
					"Olusola",
					[]string{"USD"},
					"8de0c7a2-0004-4420-899b-f8d89c81f82b",
					"Alao",
					"RC-546789",
					UsdcAddress{
						"0x7af2fa93d4069655098a24dd055ff7aa51bab531",
						"ETH",
						"USD",
					},
				},
			},
		},
		{
			name: "update integration (error; invalid domain name)",
			mockHttpClient: &MockHttpClient{

				DoFunc: func(r *http.Request) (*http.Response, error) {
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(`{
					"errors":{"domain":["This field must be a valid URL."]},"message":"Unprocessable entity"}`)))
					return &http.Response{
						StatusCode: 422,
						Body:       responseBody,
					}, nil
				},
				Timeout: 0,
			},
			args: args{
				businessAdd: "New address 1",
				domain:      "wrong_domain.com",
				webhook:     "https://webhook.site/043c6db2-5c17-4885-b769-4491ce3b0b0e",
			},
			want:    AccountResp{},
			wantErr: true,
			errMsg:  "unprocessable entity this field must be a valid url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(tt.mockHttpClient)
			got, err := cl.UpdateAccount(tt.args.webhook, tt.args.businessAdd, tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("UpdateAccount() actual error = %v, expected error %v", err, tt.errMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_RegisterUser(t *testing.T) {
	type args struct {
		payload RegisterUserData
		id      string
	}
	tests := []struct {
		name           string
		mockHttpClient MockHttpClient
		args           args
		want           UserResp
		wantErr        bool
		errMsg         string
	}{
		{
			name: "Create an account for your users requesting a card (Success)",
			mockHttpClient: MockHttpClient{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/card-integrators/27de9f46-726a-4499-aa62-27c3ed274026/register-user" {
						t.Errorf("Expected to request '/card-integrators/27de9f46-726a-4499-aa62-27c3ed274026/register-user', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(
						`{
						   "data": {
							   "address": {
								   "city": "Lagos",
								   "country": "NG",
								   "line1": "Lekki Phase 1",
								   "line2": "Saltana Park",
								   "state": "Lagos",
								   "zip_code": "101233"
							   },
							   "archived": false,
							   "card_integrator_id": "27de9f46-726a-4499-aa62-27c3ed274026",
							   "email": "user2@gmail.com",
							   "first_name": "Olusola",
							   "id": "be2c7d1c-c02a-4925-a7c4-4c5b4fc579f1",
							   "id_number": "00000000000",
							   "id_type": "BVN",
							   "last_name": "Alao",
							   "phone_number": "+2348023547675",
							   "verified": false
						   }
						}`)))
					return &http.Response{
						StatusCode: 201,
						Body:       responseBody,
					}, nil
				},
			},
			args: args{
				payload: RegisterUserData{
					Address: UserAddress{
						City:    "Lagos",
						Country: "NG",
						Line1:   "Lekki Phase 1",
						Line2:   "Saltana Park",
						State:   "Lagos",
						ZipCode: "101233",
					},
					Email:       "user2@gmail.com",
					FirstName:   "Olusola",
					IdNumber:    "00000000000",
					IdType:      "BVN",
					LastName:    "Alao",
					PhoneNumber: "+2348023547675",
				},
				id: "27de9f46-726a-4499-aa62-27c3ed274026",
			},
			want: UserResp{
				Data: User{
					Address: UserAddress{
						City:    "Lagos",
						Country: "NG",
						Line1:   "Lekki Phase 1",
						Line2:   "Saltana Park",
						State:   "Lagos",
						ZipCode: "101233",
					},
					Archived:         false,
					CardIntegratorId: "27de9f46-726a-4499-aa62-27c3ed274026",
					Email:            "user2@gmail.com",
					FirstName:        "Olusola",
					Id:               "be2c7d1c-c02a-4925-a7c4-4c5b4fc579f1",
					IdNumber:         "00000000000",
					IdType:           "BVN",
					LastName:         "Alao",
					PhoneNumber:      "+2348023547675",
					Verified:         false,
				},
			},
			wantErr: false,
		},

		{
			name: "Create an account for your users requesting a card (Error: user already_exists)",
			mockHttpClient: MockHttpClient{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"message": "already_exists"}`)))
					return &http.Response{
						StatusCode: 400,
						Body:       responseBody,
					}, nil
				},
			},
			args: args{
				payload: RegisterUserData{
					Address: UserAddress{
						City:    "Abuja",
						Country: "NG",
						Line1:   "1 Ademola Tokunbo Str",
						Line2:   "suite 206",
						State:   "Lagos",
						ZipCode: "101233",
					},
					Email:       "user1@gmail.com",
					FirstName:   "Joe",
					IdNumber:    "00000000000",
					IdType:      "BVN",
					LastName:    "Doe",
					PhoneNumber: "+2349034384664",
					UserPhoto:   " ",
				},
				id: "8de0c7a2-0004-4420-899b-f8d89c81f82b",
			},
			wantErr: true,
			errMsg:  "already_exists",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(&tt.mockHttpClient)
			got, err := cl.RegisterUser(tt.args.payload, tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("RegisterUser() actual error = %v, expected error %v", err, tt.errMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListUsers(t *testing.T) {

	type args struct {
		limit int
		page  int
	}
	tests := []struct {
		name           string
		mockHttpClient MockHttpClient
		args           args
		want           UsersResp
		wantErr        bool
	}{
		{
			name: "Get a list of card users attached to the account (Success)",
			mockHttpClient: MockHttpClient{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/card-integrators/card-users" {
						t.Errorf("Expected to request '/card-integrators/card-users', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(
						`{
							"data": [
								{
									"address": {
										"city": "Lagos",
										"country": "NG",
										"line1": "Lekki Phase 1",
										"line2": null,
										"state": null,
										"zip_code": "101233"
									},
									"archived": false,
									"card_integrator_id": "27de9f46-726a-4499-aa62-27c3ed274026",
									"email": "user1@gmail.com",
									"first_name": "Olusola",
									"id": "1c607ba6-4a59-405a-bf63-55cb76078ade",
									"id_number": "00000000000",
									"id_type": "BVN",
									"last_name": "Alao",
									"phone_number": "+2348023547672",
									"verified": true
								}
							],
							"page": 1,
							"total": 1,
							"total_pages": 1
						}`)))
					return &http.Response{
						StatusCode: 200,
						Body:       responseBody,
					}, nil
				},
			},
			args: args{
				limit: 10,
				page:  1,
			},
			want: UsersResp{
				Page:       1,
				Total:      1,
				TotalPages: 1,
				Data: []User{
					{
						Address: UserAddress{
							City:    "Lagos",
							Country: "NG",
							Line1:   "Lekki Phase 1",
							Line2:   nil,
							State:   nil,
							ZipCode: "101233",
						},
						Archived:         false,
						CardIntegratorId: "27de9f46-726a-4499-aa62-27c3ed274026",
						Email:            "user1@gmail.com",
						FirstName:        "Olusola",
						Id:               "1c607ba6-4a59-405a-bf63-55cb76078ade",
						IdNumber:         "00000000000",
						IdType:           "BVN",
						LastName:         "Alao",
						PhoneNumber:      "+2348023547672",
						Verified:         true,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cl.httpClient = &tt.mockHttpClient
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(&tt.mockHttpClient)
			got, err := cl.ListUsers(tt.args.limit, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_CreateCard(t *testing.T) {
	type args struct {
		req CreateCardData
	}
	tests := []struct {
		name           string
		mockHttpClient MockHttpClient
		args           args
		want           CreateCardResp
		wantErr        bool
	}{
		{
			name: "create/order a virtual card for a user (Success)",
			mockHttpClient: MockHttpClient{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/cards/create-virtual-card" {
						t.Errorf("Expected to request '/cards/create-virtual-card', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(`{
						"data": {
							"balance": 0,
							"business_id": "27de9f46-726a-4499-aa62-27c3ed274026",
							"card_name": "Olusola Alao",
							"card_number": "5368988843030561",
							"card_type": "virtual",
							"currency": "USD",
							"cvv2": "149",
							"design_type": "Aurora",
							"expiry": "2022-05-17T00:00:00.000Z",
							"id": "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
							"provider": "union54",
							"single_use": false,
							"status": "active",
							"user_id": "be2c7d1c-c02a-4925-a7c4-4c5b4fc579f1",
							"valid": "05/22"
						}
					}`)))
					return &http.Response{
						StatusCode: 201,
						Body:       responseBody,
					}, nil
				},
			},
			args: args{req: CreateCardData{
				DesignType:       "Aurora",
				SingleUse:        false,
				Source:           "integrator",
				CardIntegratorId: "27de9f46-726a-4499-aa62-27c3ed274026",
				Currency:         "USD",
				JuiceUserId:      "be2c7d1c-c02a-4925-a7c4-4c5b4fc579f1",
				Validity:         30,
			}},
			want: CreateCardResp{
				Data: Card{
					Balance:    0,
					BusinessId: "27de9f46-726a-4499-aa62-27c3ed274026",
					CardName:   "Olusola Alao",
					CardNumber: "5368988843030561",
					CardType:   "virtual",
					Currency:   "USD",
					Cvv2:       "149",
					DesignType: "Aurora",
					Expiry:     time.Date(2022, 05, 17, 00, 00, 00, +0000, time.UTC),
					Id:         "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
					Provider:   "union54",
					SingleUse:  false,
					Status:     "active",
					UserId:     "be2c7d1c-c02a-4925-a7c4-4c5b4fc579f1",
					Valid:      "05/22",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(&tt.mockHttpClient)
			got, err := cl.CreateCard(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateCard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListCards(t *testing.T) {
	type args struct {
		limit int
		page  int
		id    string
	}
	tests := []struct {
		name           string
		mockHttpClient MockHttpClient
		args           args
		want           []CardResp
		wantErr        bool
	}{
		{
			name: "get a particular user's cards (success)",
			mockHttpClient: MockHttpClient{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/cards" {
						t.Errorf("Expected to request '/cards', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(
						`[
							{
								"balance": 19900,
								"card_number": "5368988843030561",
								"card_type": "virtual",
								"cvv2": "149",
								"expiry": "2022-05-17T00:00:00.000Z",
								"id": "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
								"single_use": false,
								"status": "active",
								"valid": "05/22"
							}
						]`,
					)))
					return &http.Response{
						StatusCode: 200,
						Body:       responseBody,
					}, nil
				},
				Timeout: 0,
			},
			args: args{
				limit: 15,
				page:  1,
				id:    "be2c7d1c-c02a-4925-a7c4-4c5b4fc579f1",
			},
			want: []CardResp{
				{
					Balance:    19900,
					CardNumber: "5368988843030561",
					CardType:   "virtual",
					Cvv2:       "149",
					Expiry:     time.Date(2022, 05, 17, 00, 00, 00, +0000, time.UTC),
					Id:         "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
					SingleUse:  false,
					Status:     "active",
					Valid:      "05/22",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(&tt.mockHttpClient)
			got, err := cl.ListCards(tt.args.limit, tt.args.page, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListCards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListCards() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetCard(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name           string
		mockHttpClient MockHttpClient
		args           args
		want           CardResp
		wantErr        bool
	}{
		{
			name: "Get a particular card (Success)",
			mockHttpClient: MockHttpClient{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/cards/0c7ca765-764c-4f62-9c35-ac3e2abcee01" {
						t.Errorf("Expected to request '/cards/0c7ca765-764c-4f62-9c35-ac3e2abcee01', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(
						`{
							"balance": 19900,
							"card_number": "5368988843030561",
							"card_type": "virtual",
							"cvv2": "149",
							"expiry": "2022-05-17T00:00:00.000Z",
							"id": "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
							"single_use": false,
							"status": "active",
							"valid": "05/22"
						}`,
					)))
					return &http.Response{
						StatusCode: 200,
						Body:       responseBody,
					}, nil
				},
				Timeout: 0,
			},
			args: args{
				id: "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
			},
			want: CardResp{
				Balance:    19900,
				CardNumber: "5368988843030561",
				CardType:   "virtual",
				Cvv2:       "149",
				Expiry:     time.Date(2022, 05, 17, 00, 00, 00, +0000, time.UTC),
				Id:         "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
				SingleUse:  false,
				Status:     "active",
				Valid:      "05/22",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(&tt.mockHttpClient)
			got, err := cl.GetCard(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_CreditCard(t *testing.T) {
	type args struct {
		req PaymentData
	}
	tests := []struct {
		name           string
		mockHttpClient MockHttpClient
		args           args
		want           CardResp
		wantErr        bool
	}{
		{
			name: "Top up a card for a user (Success)",
			mockHttpClient: MockHttpClient{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/cards/credit/balance" {
						t.Errorf("Expected to request '/cards/credit/balance', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(
						`{
						"balance": 20000,
						"card_number": "5368988843030561",
						"card_type": "virtual",
						"cvv2": "149",
						"expiry": "2022-05-17T00:00:00.000Z",
						"id": "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
						"single_use": false,
						"status": "active",
						"valid": "05/22"
					}`)))
					return &http.Response{
						StatusCode: 200,
						Body:       responseBody,
					}, nil
				},
				Timeout: 0,
			},
			args: args{PaymentData{
				Source: "integrator",
				Amount: 20000,
				CardId: "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
			}},
			want: CardResp{
				Balance:    20000,
				CardNumber: "5368988843030561",
				CardType:   "virtual",
				Cvv2:       "149",
				Expiry:     time.Date(2022, 05, 17, 00, 00, 00, +0000, time.UTC),
				Id:         "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
				SingleUse:  false,
				Status:     "active",
				Valid:      "05/22",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(&tt.mockHttpClient)
			got, err := cl.CreditCard(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreditCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreditCard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DebitCard(t *testing.T) {
	type args struct {
		req PaymentData
	}
	tests := []struct {
		name           string
		mockHttpClient MockHttpClient
		args           args
		want           CardResp
		wantErr        bool
	}{
		{
			name: "Debit a card for a user (Success)",
			mockHttpClient: MockHttpClient{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/cards/debit/balance" {
						t.Errorf("Expected to request '/cards/debit/balance', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(
						`{
							"balance": 19900,
							"card_number": "5368988843030561",
							"card_type": "virtual",
							"cvv2": "149",
							"expiry": "2022-05-17T00:00:00.000Z",
							"id": "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
							"single_use": false,
							"status": "active",
							"valid": "05/22"
						}`)))
					return &http.Response{
						StatusCode: 200,
						Body:       responseBody,
					}, nil
				},
				Timeout: 0,
			},
			args: args{PaymentData{
				Source: "integrator",
				Amount: 100,
				CardId: "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
			}},
			want: CardResp{
				Balance:    19900,
				CardNumber: "5368988843030561",
				CardType:   "virtual",
				Cvv2:       "149",
				Expiry:     time.Date(2022, 05, 17, 00, 00, 00, +0000, time.UTC),
				Id:         "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
				SingleUse:  false,
				Status:     "active",
				Valid:      "05/22",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(&tt.mockHttpClient)
			got, err := cl.DebitCard(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DebitCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DebitCard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_FreezeCard(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name           string
		mockHttpClient MockHttpClient
		args           args
		want           CardResp
		wantErr        bool
	}{
		{
			name: "freeze the given card (Success)",
			mockHttpClient: MockHttpClient{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/cards/0c7ca765-764c-4f62-9c35-ac3e2abcee01/freeze" {
						t.Errorf("Expected to request '/cards/0c7ca765-764c-4f62-9c35-ac3e2abcee01/freeze', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(
						`{
							"balance": 19900,
							"card_number": "5368988843030561",
							"card_type": "virtual",
							"cvv2": "149",
							"expiry": "2022-05-17T00:00:00.000Z",
							"id": "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
							"single_use": false,
							"status": "inactive",
							"valid": "05/22"
						}`,
					)))
					return &http.Response{
						StatusCode: 200,
						Body:       responseBody,
					}, nil
				},
				Timeout: 0,
			},
			args: args{
				id: "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
			},
			want: CardResp{
				Balance:    19900,
				CardNumber: "5368988843030561",
				CardType:   "virtual",
				Cvv2:       "149",
				Expiry:     time.Date(2022, 05, 17, 00, 00, 00, +0000, time.UTC),
				Id:         "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
				SingleUse:  false,
				Status:     "inactive",
				Valid:      "05/22",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(&tt.mockHttpClient)
			got, err := cl.FreezeCard(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FreezeCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FreezeCard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UnfreezeCard(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name           string
		mockHttpClient MockHttpClient
		args           args
		want           CardResp
		wantErr        bool
	}{
		{
			name: "unfreeze the given card (Success)",
			mockHttpClient: MockHttpClient{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/cards/0c7ca765-764c-4f62-9c35-ac3e2abcee01/unfreeze" {
						t.Errorf("Expected to request '/cards/0c7ca765-764c-4f62-9c35-ac3e2abcee01/unfreeze', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(
						`{
							"balance": 19900,
							"card_number": "5368988843030561",
							"card_type": "virtual",
							"cvv2": "149",
							"expiry": "2022-05-17T00:00:00.000Z",
							"id": "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
							"single_use": false,
							"status": "active",
							"valid": "05/22"
						}`,
					)))
					return &http.Response{
						StatusCode: 200,
						Body:       responseBody,
					}, nil
				},
				Timeout: 0,
			},
			args: args{
				id: "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
			},
			want: CardResp{
				Balance:    19900,
				CardNumber: "5368988843030561",
				CardType:   "virtual",
				Cvv2:       "149",
				Expiry:     time.Date(2022, 05, 17, 00, 00, 00, +0000, time.UTC),
				Id:         "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
				SingleUse:  false,
				Status:     "active",
				Valid:      "05/22",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(&tt.mockHttpClient)
			got, err := cl.UnfreezeCard(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnfreezeCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnfreezeCard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListTransactions(t *testing.T) {
	type args struct {
		id    string
		param Param
	}
	tests := []struct {
		name           string
		mockHttpClient MockHttpClient
		args           args
		want           TransactionsResp
		wantErr        bool
	}{
		{
			name: "Get paginated transactions for the given card (Success)",
			mockHttpClient: MockHttpClient{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/cards/0c7ca765-764c-4f62-9c35-ac3e2abcee01/transactions" {
						t.Errorf("Expected to request '/cards/0c7ca765-764c-4f62-9c35-ac3e2abcee01/transactions', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(
						`{
							"data": [
								{
									"amount": 100,
									"card_balance_after": 19900,
									"card_balance_before": 20000,
									"conversion_rate": 1,
									"created_at": "2022-04-17T20:55:36.798Z",
									"credit_currency": null,
									"credit_id": null,
									"debit_currency": "USD",
									"debit_id": "34d1f60e-0ff2-4e02-9b82-b1faa10f9153",
									"id": "9b14c12e-2e20-4a4e-8508-ab816b8e575c",
									"narrative": null,
									"type": "debit"
								},
								{
									"amount": 20000,
									"card_balance_after": 20000,
									"card_balance_before": 0,
									"conversion_rate": 1,
									"created_at": "2022-04-17T20:41:32.483Z",
									"credit_currency": "USD",
									"credit_id": "34d1f60e-0ff2-4e02-9b82-b1faa10f9153",
									"debit_currency": null,
									"debit_id": null,
									"id": "e406f24c-6651-44fa-8566-e03e2a1c3623",
									"narrative": null,
									"type": "credit"
								}
							],
							"message": "Ok",
							"next_page": null
						}`,
					)))
					return &http.Response{
						StatusCode: 200,
						Body:       responseBody,
					}, nil
				},
				Timeout: 0,
			},

			args: args{
				id: "0c7ca765-764c-4f62-9c35-ac3e2abcee01",
				param: Param{
					Limit: 50,
					Page:  1,
				},
			},
			want: TransactionsResp{
				Data: []Transaction{
					{
						Amount:            100,
						CardBalanceAfter:  19900,
						CardBalanceBefore: 20000,
						ConversionRate:    1,
						CreatedAt:         time.Date(2022, 04, 17, 20, 55, 36, 798000000, time.UTC),
						CreditCurrency:    nil,
						CreditId:          nil,
						DebitCurrency:     "USD",
						DebitId:           "34d1f60e-0ff2-4e02-9b82-b1faa10f9153",
						Id:                "9b14c12e-2e20-4a4e-8508-ab816b8e575c",
						Narrative:         nil,
						Type:              "debit",
					},
					{
						Amount:            20000,
						CardBalanceAfter:  20000,
						CardBalanceBefore: 0,
						ConversionRate:    1,
						CreatedAt:         time.Date(2022, 04, 17, 20, 41, 32, 483000000, time.UTC),
						CreditCurrency:    "USD",
						CreditId:          "34d1f60e-0ff2-4e02-9b82-b1faa10f9153",
						DebitCurrency:     nil,
						DebitId:           nil,
						Id:                "e406f24c-6651-44fa-8566-e03e2a1c3623",
						Narrative:         nil,
						Type:              "credit",
					},
				},
				Message:  "Ok",
				NextPage: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(&tt.mockHttpClient)
			got, err := cl.ListTransactions(tt.args.id, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTransactions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetTransaction(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name           string
		mockHttpClient MockHttpClient
		args           args
		want           TransactionResp
		wantErr        bool
	}{
		{
			name: "Get a transaction information using its reference (Success)",
			mockHttpClient: MockHttpClient{
				DoFunc: func(r *http.Request) (*http.Response, error) {
					if r.URL.Path != "/cards/transaction/9b14c12e-2e20-4a4e-8508-ab816b8e575c" {
						t.Errorf("Expected to request '/cards/transaction/9b14c12e-2e20-4a4e-8508-ab816b8e575c', got: %s", r.URL.Path)
					}
					if r.Header.Get("Content-Type") != "application/json" {
						t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
					}
					responseBody := ioutil.NopCloser(bytes.NewReader([]byte(
						`{
							"data": {
								"amount": 100,
								"card_balance_after": 19900,
								"card_balance_before": 20000,
								"conversion_rate": 1,
								"created_at": "2022-04-17T20:55:36.798Z",
								"credit_currency": null,
								"credit_id": null,
								"debit_currency": "USD",
								"debit_id": "34d1f60e-0ff2-4e02-9b82-b1faa10f9153",
								"id": "9b14c12e-2e20-4a4e-8508-ab816b8e575c",
								"narrative": null,
								"type": "debit"
							},
							"message": "Ok"
						}`,
					)))
					return &http.Response{
						StatusCode: 200,
						Body:       responseBody,
					}, nil
				},
				Timeout: 0,
			},
			args: args{
				id: "9b14c12e-2e20-4a4e-8508-ab816b8e575c",
			},
			want: TransactionResp{
				Data: Transaction{
					Amount:            100,
					CardBalanceAfter:  19900,
					CardBalanceBefore: 20000,
					ConversionRate:    1,
					CreatedAt:         time.Date(2022, 04, 17, 20, 55, 36, 798000000, time.UTC),
					CreditCurrency:    nil,
					CreditId:          nil,
					DebitCurrency:     "USD",
					DebitId:           "34d1f60e-0ff2-4e02-9b82-b1faa10f9153",
					Id:                "9b14c12e-2e20-4a4e-8508-ab816b8e575c",
					Narrative:         nil,
					Type:              "debit",
				},
				Message: "Ok",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl.SetHTTPClient(&tt.mockHttpClient)
			got, err := cl.GetTransaction(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_health(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "Health check",
			want: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cl.Health()
			if (err != nil) != tt.wantErr {
				t.Errorf("health() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("health() got = %v, want %v", got, tt.want)
			}
		})
	}
}
