![image](https://user-images.githubusercontent.com/62708917/163804624-513018a3-ce45-4220-8c01-5a2f7bdce94c.png)

# Spend-Juice Go Library

## Introduction
This is a Go wrapper around the [API](https://docs.spendjuice.org/reference/) for [Cards Integrators by Spend-Juice](https://spendjuice.org/).

## Installation
To install, run

``` go get github.com/bushaHQ/spend-juice-go```

## Import Package
The base class for this package is 'spend-juice-go'. To use this class, add:

```
 import (
 	"github.com/bushaHQ/spend-juice-go"
 )
 ```

## Initialization

To use Spend Juice, instantiate Spend-Juice with your public key. We recommend that you store your secret key in an environment variable named, ```JUICE_PRIVATE_KEY```. See example below.
 ```
	err := godotenv.Load("./.env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	client := juice.NewClient()
 ```
You can override the default settings by passing in the following parameters:

```
    client.SetBaseURL("your-base-url")
    client.SetDebug(false)
    client.SetHTTPClient(&http.Client{
        Timeout: your-timeout,
    })
    client.SetAuth(os.Getenv({JUICE_PRIVATE_KEY}))
```

# Card Integration Methods
This is the documentation for all of the components of card Integrator

**Methods Included:**

* ```.RegisterAccount```

* ```.UpdateAccount```

* ```.RegisterUser```

* ```.ListUsers```

* ```.CreateCard```

* ```.ListCards```

* ```.GetCard```

* ```.CreditCard```

* ```.DebitCard```

* ```.FreezeCard```

* ```.UnfreezeCard```

* ```.ListTransactions```

* ```.GetTransaction```

### ```.RegisterAccount(data RegisterAccountData) (AccountResp, error)```
This is called to create a card integrator account. The payload should be of type ```juice.RegisterAccountData```. See below for  ```juice.RegisterAccountData``` definition

```
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
```
A sample register call is:

```
    payload := juice.RegisterAccountData{
        BusinessAddress:    "Ajah",
        BusinessName:       "Algoro",
        Chain:              "ETH",
        ContactNumber:      "+2349034384669",
        Country:            "NG",
        Domain:             "https://ajalekoko.com",
        Email:              "ajalenkoko@gmail.com",
        FirstName:          "Olusola",
        FloatCurrencies:    []string{"USD"},
        LastName:           "Alao",
        Password:           "@Password",
        RegistrationNumber: "RC-5467898",
        WebhookUrl:         "https://webhook.site/7bdf91c4-6e84-4ff1-a3e1-185f138247c1",
    }

    response, err := client.RegisterAccount(payload)

    if err != nil {
        panic(err)
    }
    
    fmt.Println(response)
```
#### Sample Response

```
    {{Lekki Ikate Algo Math ETH +2349034384660 NG https://boro.com alaoolusolae@gmail.com Olusola [USD] 27de9f46-726a-4499-aa62-27c3ed274026 Alao RC-546787 {0xbaf2e14f27c106f0d078b397af9c7eebf1611d46 ETH USD}}}
```

### ```.UpdateAccount(webhook, businessAddress, domain string) (AccountResp, error)```
This is called to update the integrator account.

A sample validate call is:

```
    response, err := client.UpdateAccount("https://webhook.site/043c6db2-5c17-4885-b769-4491ce3b0b0e", "New Road Lincoln street", "https://olusolaa.tech")
    if err != nil {
        panic(err)
    }
    fmt.Println(response)
```
#### Sample Response

```
   {{New Road Lincoln street algo math  +2349034384662 NG https://olusolaa.tech email@gmail.com Olusola [USD] 8de0c7a2-0004-4420-899b-f8d89c81f82b Alao 12345 {0x7af2fa93d4069655098a24dd055ff7aa51bab531 ETH USD}}}
```

### ```.RegisterUser(data RegisterUserData, accountId string) (UserResp, error)```
This is called to create an account for user requesting a card. The payload should be of type ```spend-juice-go.RegisterUserData```. See below for  ```spend-juice-go.RegisterUserData``` definition
```
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
```
IdNumber is the user identity card number (in sandbox environment use 00000000000 as your BVN number)

A sample register-user call is:

```
    payload := juice.RegisterUserData{
        Address: juice.UserAddress{
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
    }
        
    accountId := "27de9f46-726a-4499-aa62-27c3ed274026"

    response, err := client.RegisterUser(payload, id)
    if err != nil {
        panic(err)
    }
    fmt.Println(response)
```
#### Sample Response

```
{{{Lagos NG Lekki Phase 1 Saltana Park Lagos 101233} false 27de9f46-726a-4499-aa62-27c3ed274026 user2@gmail.com Olusola be2c7d1c-c02a-4925-a7c4-4c5b4fc579f1 00000000000 BVN Alao +2348023547675 false}}
```
### ```.ListUsers(limit, page int) (usersResp, error)```
This is called to get a list of card users attached to an integrator account. Limit is the number of users dispaly per page.

A sample list users call is:

```
    response, err := client.ListUsers(10, 1)
    if err != nil {
        panic(err)
    }
    fmt.Println(response)
```
#### Sample Response

```
{1 1 1 [{{Lagos NG Lekki Phase 1 <nil> <nil> 101233} false 27de9f46-726a-4499-aa62-27c3ed274026 user1@gmail.com Olusola 1c607ba6-4a59-405a-bf63-55cb76078ade 00000000000 BVN Alao +2348023547672 true}]}
```
