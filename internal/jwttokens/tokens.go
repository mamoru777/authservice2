package jwttokens

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func New() *Tokens {
	return &Tokens{}
}
