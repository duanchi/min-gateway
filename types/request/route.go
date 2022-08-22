package request

type RouteUrl struct {
	Type  string `json:"type"`
	Match string `json:"match"`
}

type RouteBlueTag struct {
	Tag       string `json:"tag"`
	ServiceId string `json:"service_id"`
}

type RouteRequest struct {
	Id               string            `json:"id"`
	Url              RouteUrl          `json:"url"`
	Method           []string          `json:"method"`
	ServiceId        string            `json:"service_id"`
	Authorize        bool              `json:"authorize"`
	AuthorizePrefix  string            `json:"authorize_prefix"`
	AuthorizeTypeKey string            `json:"authorize_type_key"`
	CustomToken      bool              `json:"custom_token"`
	Rewrite          map[string]string `json:"rewrite"`
	BlueTagKey       string            `json:"blue_tag_key"`
	Blue             []RouteBlueTag    `json:"blue"`
	Description      string            `json:"description"`
	Timeout          int64             `json:"timeout"`
	Order            int64             `json:"order"`
}
