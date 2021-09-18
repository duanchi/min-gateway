package types

type RoutesMap map[string]Route

type Route struct {
	Id string `json:"id"`
	Url struct{
		Type string `json:"type"`
		Match string `json:"match"`
	} `json:"url"`
	Method []string `json:"method"`
	Service string `json:"service"`
	Authorize bool `json:"authorize"`
	AuthorizePrefix string `json:"authorize_prefix"`
	AuthorizeTypeKey string `json:"authorize_type_key"`
	CustomToken bool `json:"custom_token"`
	Rewrite map[string]string `json:"rewrite"`
	Description string `json:"description"`
	Order int `json:"order"`
}

type RoutesArray []Route

func (this RoutesArray) Len () int {
	return len(this)
}

func (this RoutesArray) Swap (i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this RoutesArray) Less (i, j int) bool {
	return this[i].Order < this[j].Order
}