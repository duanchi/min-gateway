package response

type IntegrationKeys struct {
}

type IntegrationResponse struct {
	Id              string `json:"id"`
	Alias           string `json:"alias"`
	Url             string `json:"url"`
	Protocol        string `json:"protocol"`
	DataType        string `json:"data_type"`
	RequestMethod   string `json:"request_method"`
	RequestTemplate string `json:"request_template"`
}
