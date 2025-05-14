package domain

type Account struct {
	AccountId       int32   `json:"id" xml:"id"`
	ClientId        int32   `json:"client_id" xml:"client_id"`
	UserId          int32   `json:"user_id" xml:"user_id"`
	AccountPassword string  `json:"account_password" xml:"account_password"`
	AccountBalance  float64 `json:"account_balance" xml:"account_balance"`
	Status          bool    `json:"status" xml:"status"`
}

func (a *Account) GetAccountId() int32                { return a.AccountId }
func (a *Account) GetClientId() int32                 { return a.ClientId }
func (a *Account) GetUserId() int32                   { return a.UserId }
func (a *Account) GetAccountPassword() string         { return a.AccountPassword }
func (a *Account) GetAccountBalance() float64         { return a.AccountBalance }
func (a *Account) GetStatus() bool                    { return a.Status }
func (a *Account) SetAccountId(id int32)              { a.AccountId = id }
func (a *Account) SetClientId(id int32)               { a.ClientId = id }
func (a *Account) SetUserId(id int32)                 { a.UserId = id }
func (a *Account) SetAccountPassword(password string) { a.AccountPassword = password }
func (a *Account) SetAccountBalance(balance float64)  { a.AccountBalance = balance }
func (a *Account) SetStatus(status bool)              { a.Status = status }

func (a *Account) GetBalance() float64 {
	return a.AccountBalance
}

func (a *Account) AddBalance(amount float64) {
	a.AccountBalance += amount
}

func (a *Account) SubtractBalance(amount float64) {
	a.AccountBalance -= amount
}

func (a *Account) TransferBalance(amount float64, targetAccount *Account) {
	a.SubtractBalance(amount)
	targetAccount.AddBalance(amount)
}

func (a *Account) ValidatePassword(password string) bool {
	return a.AccountPassword == password
}

func (a *Account) UpdatePassword(newPassword string) {
	a.AccountPassword = newPassword
}

func (a *Account) IsActive() bool {
	return a.Status
}

func (a *Account) Activate() {
	a.Status = true
}

func (a *Account) Deactivate() {
	a.Status = false
}
