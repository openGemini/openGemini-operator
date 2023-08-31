package configfile

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateConfig_ERROR(t *testing.T) {
	tmpl := `
[common] 
meta-join = ["opengemini-meta-01.test.svc.cluster.local:8092
[meta]
  bind-address = "{{HOST_IP}}:8088"`
	newConf := `
[meta]
   bind-address = "{{addr}}:8088"`

	_, err := UpdateConfig(tmpl, newConf)
	assert.Contains(t, err.Error(), "strings cannot contain newlines")

	tmpl = `
[common] 
meta-join = ["opengemini-meta-01.test.svc.cluster.local:8092"]
[meta]
  bind-address = "{{HOST_IP}}:8088"`
	newConf = `
[meta]
   bind-address = "{{addr}}`

	_, err = UpdateConfig(tmpl, newConf)
	assert.Contains(t, err.Error(), "unexpected EOF")
}

func Test_UpdateConfig_OK(t *testing.T) {
	tmpl := `
[common] 
meta-join = ["opengemini-meta-01.test.svc.cluster.local:8092", "opengemini-meta-02.test.svc.cluster.local:8092", "opengemini-meta-03.test.svc.cluster.local:8092"]
cpu-num = 8  
memory-size = "32G"

[meta]
  bind-address = "{{HOST_IP}}:8088"
  http-bind-address = "{{HOST_IP}}:8091"
  rpc-bind-address = "{{addr}}:8092"
  dir = "/opt/openGemini/meta"
  domain = "{domain}"`
	newConf := `
[meta]
   bind-address = "{{addr}}:8088"
   heartbeat-timeout = "1s"`

	confData, err := UpdateConfig(tmpl, newConf)
	assert.NoError(t, err)

	var confToml map[string]map[string]interface{}
	_, err = toml.Decode(confData, &confToml)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(confToml["common"]))
	assert.Equal(t, 6, len(confToml["meta"]))
	assert.Equal(t, "{{HOST_IP}}:8088", confToml["meta"]["bind-address"])
	assert.Equal(t, "1s", confToml["meta"]["heartbeat-timeout"])
}

func Test_UpdateConfig_Nested_OK(t *testing.T) {
	tmpl := `
[common] 
meta-join = ["opengemini-meta-01.test.svc.cluster.local:8092", "opengemini-meta-02.test.svc.cluster.local:8092", "opengemini-meta-03.test.svc.cluster.local:8092"]
cpu-num = 8  
memory-size = "32G"

[data]
  store-ingest-addr = "{{HOST_IP}}:8400"
  store-select-addr = "{{HOST_IP}}:8401"

[data.ops-monitor]
      store-http-addr = "{{HOST_IP}}:9999"`
	newConf := `
[data]
  store-ingest-addr = "{{addr}}:8400"
  store-select-addr = "{{addr}}:8401"

[data.ops-monitor]
      store-http-addr = "{{addr}}:9999"
      auth-enabled = false`

	confData, err := UpdateConfig(tmpl, newConf)
	assert.NoError(t, err)

	var confToml map[string]map[string]interface{}
	_, err = toml.Decode(confData, &confToml)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(confToml["common"]))
	assert.Equal(t, 3, len(confToml["data"]))
	assert.Equal(t, "{{HOST_IP}}:8400", confToml["data"]["store-ingest-addr"])
	assert.Equal(t, false, confToml["data"]["ops-monitor"].(map[string]interface{})["auth-enabled"])
	assert.Equal(t, "{{HOST_IP}}:9999", confToml["data"]["ops-monitor"].(map[string]interface{})["store-http-addr"])
}
