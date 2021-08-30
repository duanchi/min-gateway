package types

type ServicesMap map[string]Service

type Gray struct {
	Uri	string `json:"uri"`
	Id string `json:"id"`
}

type Service struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Display string `json:"display"`
	LoadBalance string `json:"load_balance"`
	Instances []string `json:"instances"`
	Gray []Gray `json:"gray"`
}