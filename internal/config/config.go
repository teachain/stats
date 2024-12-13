package config

type DatabaseConfig struct {
	Database   string `json:"database" yaml:"database"`
	Ip         string `json:"ip" yaml:"ip"`
	Port       int    `json:"port" yaml:"port"`
	User       string `json:"user" yaml:"user"`
	Password   string `json:"password" yaml:"password"`
	Charset    string `json:"charset" yaml:"charset"`     //字符编码
	ParseTime  bool   `json:"parseTime" yaml:"parseTime"` //是否转化时间
	Loc        string `json:"loc" yaml:"loc"`             //时区设置
	Key        string `json:"key" yaml:"key"`
	Iv         string `json:"iv" yaml:"iv"`
	Ciphertext bool   `json:"ciphertext" yaml:"ciphertext"`
}

type Config struct {
	DB      *DatabaseConfig `json:"db" yaml:"db"`
	NodeURL string          `json:"node_url" yaml:"node_url"`
}
