package types

type Provider string

type User struct {
	Id        string                 `json:"id"`
	Provider  Provider               `json:"provider"`
	Domain    string                 `json:"domain"`
	Username  string                 `json:"username"`
	CreatedAt int64                  `json:"createdAt"`
	Contents  map[string]interface{} `json:"contents"`
}
