package domain

type Account struct {
	AccountId       uint32  `json:"AccountID" xml:"AccountID"`
	ClientId        uint32  `json:"ClientID" xml:"ClientID" validate:"required"`
	UserId          uint32  `json:"UserID" xml:"UserID" validate:"required"`
	AgencyID        uint32  `json:"AgencyID" xml:"AgencyID" validate:"required"`
	AccountPassword string  `json:"Password" xml:"Password" validate:"required"`
	AccountBalance  float64 `json:"Balance" xml:"Balance"`
	Status          bool    `json:"Status" xml:"Status"`
}

func (a *Account) GetAgencyID() uint32                { return a.AgencyID }
func (a *Account) GetAccountId() uint32               { return a.AccountId }
func (a *Account) GetAccountClientId() uint32         { return a.ClientId }
func (a *Account) GetAccountUserId() uint32           { return a.UserId }
func (a *Account) GetAccountPassword() string         { return a.AccountPassword }
func (a *Account) GetAccountBalance() float64         { return a.AccountBalance }
func (a *Account) GetStatus() bool                    { return a.Status }
func (a *Account) SetAgencyID(id uint32)              { a.AgencyID = id }
func (a *Account) SetAccountId(id uint32)             { a.AccountId = id }
func (a *Account) SetAccountClientId(id uint32)       { a.ClientId = id }
func (a *Account) SetAccountUserId(id uint32)         { a.UserId = id }
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
