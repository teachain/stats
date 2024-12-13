package config

import (
	"fmt"
	"github.com/teachain/stats/pkg/utils"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
)

func MustLoadConfig(filename string) (*Config, error) {
	// 读取YAML文件
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// 解析YAML数据
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func DataSource(c *DatabaseConfig) string {
	format := "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s"
	password := c.Password
	if c.Ciphertext {
		var err error
		password, err = utils.DecryptWithBase64([]byte(c.Key), []byte(c.Iv), c.Password)
		if err != nil {
			panic(err)
		}
	}
	dsn := fmt.Sprintf(format, c.User, password, c.Ip, c.Port, c.Database, c.Charset, c.ParseTime, url.QueryEscape(c.Loc))
	return dsn
}
