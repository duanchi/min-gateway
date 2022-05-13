package response

type Instance struct {
	Uri         string `json:"uri"`
	Id          string `json:"id"`
	IsEphemeral bool   `json:"is_ephemeral"`
	IsOnline    bool   `json:"is_online"`
}

type ServiceResponse struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Instances []Instance `json:"instances"`
	Gray      []Instance `json:"gray"`
}
