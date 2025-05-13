package domain

type Client struct {
	Id           int32  `json:"id" xml:"id"`
	UserId       int32  `json:"user_id" xml:"user_id"`
	Name         string `json:"name" xml:"name"`
	Document     string `json:"document" xml:"document"`
	Password     string `json:"password" xml:"password"`
	RegisterDate string `json:"register_date" xml:"register_date"`
	Status       bool   `json:"status" xml:"status"`
}

func (c *Client) GetId() int32                        { return c.Id }
func (c *Client) GetName() string                     { return c.Name }
func (c *Client) GetDocument() string                 { return c.Document }
func (c *Client) GetPassword() string                 { return c.Password }
func (c *Client) GetRegisterDate() string             { return c.RegisterDate }
func (c *Client) GetStatus() bool                     { return c.Status }
func (c *Client) SetId(id int32)                      { c.Id = id }
func (c *Client) SetName(name string)                 { c.Name = name }
func (c *Client) SetDocument(document string)         { c.Document = document }
func (c *Client) SetPassword(password string)         { c.Password = password }
func (c *Client) SetRegisterDate(registerDate string) { c.RegisterDate = registerDate }
func (c *Client) SetStatus(status bool)               { c.Status = status }
