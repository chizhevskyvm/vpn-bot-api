package config

type Config struct {
	Server   Server   `yaml:"server"`
	Vpn      Vpn      `yaml:"vpn"`
	Telegram Telegram `yaml:"telegram"`
}

type Server struct {
	Address  string `yaml:"address"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Network  string `yaml:"network"`
}

type Vpn struct {
	OpenVpn OpenVpn `yaml:"openVPN"`
}

type Telegram struct {
	Token string `yaml:"token"`
}

type OpenVpn struct {
	Enabled    bool   `yaml:"enabled"`
	RemoteIP   string `yaml:"remoteIP"`
	RemotePort int    `yaml:"remotePort"`
	Auth       string `yaml:"auth"`
	Cipher     string `yaml:"cipher"`
	Template   string `yaml:"template"`
}
