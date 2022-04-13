package request

type ServiceRequest map[string]Service

type Instance struct {
	Uri string `json:"uri"`
	Id  string `json:"id"`
}

type Service struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Instances []Instance `json:"instances"`
	Gray      []Instance `json:"gray"`
}
