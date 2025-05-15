package domain

type Client struct {
	ClientId     uint32 `json:"ClientID" xml:"ClientID"`
	UserId       uint32 `json:"UserID" xml:"UserID" validate:"required"`
	Name         string `json:"Name" xml:"Name" validate:"required"`
	Document     string `json:"Document" xml:"Document" validate:"required"`
	Password     string `json:"Password" xml:"Password" validate:"required"`
	RegisterDate string `json:"RegisterDate" xml:"RegisterDate"`
	Status       bool   `json:"Status" xml:"Status"`
}

func (c *Client) GetClientId() uint32                 { return c.ClientId }
func (c *Client) GetUserID() uint32                   { return c.UserId }
func (c *Client) GetName() string                     { return c.Name }
func (c *Client) GetDocument() string                 { return c.Document }
func (c *Client) GetPassword() string                 { return c.Password }
func (c *Client) GetRegisterDate() string             { return c.RegisterDate }
func (c *Client) GetStatus() bool                     { return c.Status }
func (c *Client) SetClientId(id uint32)               { c.ClientId = id }
func (c *Client) SetUserId(id uint32)                 { c.UserId = id }
func (c *Client) SetName(name string)                 { c.Name = name }
func (c *Client) SetDocument(document string)         { c.Document = document }
func (c *Client) SetPassword(password string)         { c.Password = password }
func (c *Client) SetRegisterDate(registerDate string) { c.RegisterDate = registerDate }
func (c *Client) SetStatus(status bool)               { c.Status = status }
