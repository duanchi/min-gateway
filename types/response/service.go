package response

type Instance struct {
	Uri string `json:"uri"`
	Id  string `json:"id"`
}

type ServiceResponse struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Instances []Instance `json:"instances"`
	Gray      []Instance `json:"gray"`
}
