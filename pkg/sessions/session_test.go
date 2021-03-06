package sessions

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/cozy/cozy-stack/pkg/config"
	"github.com/cozy/cozy-stack/pkg/instance"
)

var JWTSecret = []byte("foobar")

func TestMain(m *testing.M) {
	config.UseTestFile()
	conf := config.GetConfig()
	conf.Authentication = make(map[string]interface{})
	confAuth := make(map[string]interface{})
	confAuth["jwt_secret"] = base64.StdEncoding.EncodeToString(JWTSecret)
	conf.Authentication[config.DefaultInstanceContext] = confAuth

	delegatedInst = &instance.Instance{Domain: "external.notmycozy.net"}
	os.Exit(m.Run())
}
