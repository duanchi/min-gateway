package response

type RouteUrl struct {
	Type  string `json:"type"`
	Match string `json:"match"`
}

type RouteResponse struct {
	Id               string            `json:"id"`
	Url              RouteUrl          `json:"url"`
	Method           []string          `json:"method"`
	ServiceId        string            `json:"service_id"`
	Authorize        bool              `json:"authorize"`
	AuthorizePrefix  string            `json:"authorize_prefix"`
	AuthorizeTypeKey string            `json:"authorize_type_key"`
	CustomToken      bool              `json:"custom_token"`
	Rewrite          map[string]string `json:"rewrite"`
	Order            int64             `json:"order"`
}
