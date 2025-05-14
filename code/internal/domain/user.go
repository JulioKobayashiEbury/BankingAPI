package domain

type User struct {
	Id           int32  `json:"UserID" xml:"UserID"`
	Name         string `json:"Name" xml:"Name" validate:"required"`
	Document     string `json:"Document" xml:"Document" validate:"required"`
	Password     string `json:"Password" xml:"Password" validate:"required"`
	RegisterDate string `json:"RegisterDate" xml:"RegisterDate"`
	Status       bool   `json:"Status" xml:"Status"`
}

func (u *User) GetId() int32                        { return u.Id }
func (u *User) GetName() string                     { return u.Name }
func (u *User) GetDocument() string                 { return u.Document }
func (u *User) GetPassword() string                 { return u.Password }
func (u *User) GetRegisterDate() string             { return u.RegisterDate }
func (u *User) GetStatus() bool                     { return u.Status }
func (u *User) SetId(id int32)                      { (*u).Id = id }
func (u *User) SetName(name string)                 { u.Name = name }
func (u *User) SetDocument(document string)         { u.Document = document }
func (u *User) SetPassword(password string)         { u.Password = password }
func (u *User) SetRegisterDate(registerDate string) { u.RegisterDate = registerDate }
func (u *User) SetStatus(status bool)               { u.Status = status }
