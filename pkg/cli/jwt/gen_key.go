package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"gin/app/facade"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"strconv"
)

type Jwt struct {
	base.Command
}

// NewJwt 创建命令
func NewJwt() *Jwt {
	return &Jwt{
		Command: base.NewCommand(base.Config{
			Name:        "jwt-gen:key",
			Description: "生成jwt秘钥",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "l", Long: "length", Default: "32"},
					Desc:     "秘钥长度, 如: 32",
					Required: false,
				},
			},
		}),
	}
}

func (s *Jwt) Handle(values map[string]string) error {
	length := values["length"]
	// 字符串转int
	l, err := strconv.Atoi(length)
	if err != nil {
		terminal.Errorf("jwt secret key length must be a number")
		return nil

	}

	err = facade.Config().Set("jwt.key", s.generateJWTKey(l))
	if err != nil {
		terminal.Errorf("jwt secret key reset failed: %s", err.Error())
	} else {
		terminal.Infof("jwt secret key has been reset")
	}
	return nil
}

func (s *Jwt) generateJWTKey(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		terminal.Errorf("jwt secret key generate failed: %s", err.Error())
	}
	return base64.StdEncoding.EncodeToString(bytes)
}
