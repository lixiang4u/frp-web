package model

type Vhost struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	CustomDomain string `json:"custom_domain"`
	LocalAddr    string `json:"local_addr"`
	RemotePort   int    `json:"remote_port"`
	CrtPath      string `json:"crt_path"`
	KeyPath      string `json:"key_path"`
	Status       bool   `json:"status"` //true.开启，false.关闭
	CreatedAt    int64  `json:"created_at"`
}
