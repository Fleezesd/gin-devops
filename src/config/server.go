package config

import (
	"os"
	"time"

	"github.com/fleezesd/gin-devops/src/common"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
)

type ServerConfig struct {
	HttpAddr    string          `yaml:"http_addr"`
	LogLevel    string          `yaml:"log_level"`
	LogFilePath string          `yaml:"log_file_path"`
	Logger      *otelzap.Logger `yaml:"-"`
	JWT         *JWT            `yaml:"jwt"`
	Mysql       *mysql.Config   `yaml:"mysql"` // 直接用gorm中的mysql driver
}

type JWT struct {
	SigningKey      string        `yaml:"signing_key" json:"signing_key" `   // jwt签名 密码加盐
	ExpiresTime     string        `yaml:"expires_time" json:"expires_time" ` // 过期时间
	ExpiresDuration time.Duration `yaml:"-"`                                 // 过期时间段
	BufferTime      string        `yaml:"buffer_time" json:"buffer_time" `   // 缓冲时间
	BufferDuration  time.Duration `yaml:"-"`                                 // 缓冲时间段
	Issuer          string        `yaml:"issuer" json:"issuer" `             // 签发者
}

func LoadServer(filename string) (*ServerConfig, error) {
	cfg := &ServerConfig{}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}

	// 解析jwt过期&缓冲时间
	if exd, err := time.ParseDuration(cfg.JWT.ExpiresTime); err != nil {
		return nil, err
	} else {
		cfg.JWT.ExpiresDuration = exd
	}

	if bfd, err := time.ParseDuration(cfg.JWT.BufferTime); err != nil {
		return nil, err
	} else {
		cfg.JWT.BufferDuration = bfd
	}

	return cfg, nil
}

func GetServerConfig(c *gin.Context) *ServerConfig {
	return c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*ServerConfig)
}
