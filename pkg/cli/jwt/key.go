package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"strconv"
)

type Jwt struct {
	base.BaseCommand
}

func (s *Jwt) Name() string {
	return "gen:jwt-key"
}

func (s *Jwt) Description() string {
	return "生成jwt秘钥"
}

func (s *Jwt) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short:   "l",
				Long:    "length",
				Default: "32",
			},
			"秘钥长度, 如: 32",
			false,
		},
	}
}

func (s *Jwt) Execute(args []string) {
	values := s.ParseFlags(s.Name(), args, s.Help())
	length := values["length"]
	// 字符串转int
	l, err := strconv.Atoi(length)
	if err != nil {
		flag.Errorf("jwt secret key length must be a number")
		return
	}

	err = facade.Config().Set("jwt.key", s.generateJWTKey(l))
	if err != nil {
		flag.Errorf("jwt secret key reset failed: %s", err.Error())
	} else {
		flag.Infof("jwt secret key has been reset")
	}
}

func (s *Jwt) generateJWTKey(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		flag.Errorf("jwt secret key generate failed: %s", err.Error())
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

func init() {
	cli.Register(&Jwt{})
}
