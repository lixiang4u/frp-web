package model

type Vhost struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	CustomDomain string `json:"custom_domain"`
	LocalAddr    string `json:"local_addr"`
	CrtPath      string `json:"crt_path"`
	KeyPath      string `json:"key_path"`
}
