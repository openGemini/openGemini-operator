package configfile

import (
	"bytes"

	"gopkg.in/ini.v1"
)

func NewConfiguration() (string, error) {
	cfg := ini.Empty()

	common, _ := cfg.NewSection("common")
	common.NewKey("meta-join", `["hostname1:8092", "hostname2:8092", "hostname3:8092"]`)

	http, _ := cfg.NewSection("http")
	http.NewKey("bind-address", "192.168.0.1:8086")

	meta, _ := cfg.NewSection("meta")
	meta.NewKey("bind-address", "192.168.0.1:8088")
	meta.NewKey("http-bind-address", "192.168.0.1:8091")
	meta.NewKey("rpc-bind-address", "192.168.0.1:8092")
	meta.NewKey("dir", "/tmp/openGemini/data/meta")

	data, _ := cfg.NewSection("data")
	data.NewKey("store-ingest-addr", "192.168.0.1:8400")
	data.NewKey("store-select-addr", "192.168.0.1:8401")
	data.NewKey("store-data-dir", "/tmp/openGemini/data")
	data.NewKey("store-wal-dir", "/tmp/openGemini/data")
	data.NewKey("store-meta-dir", "/tmp/openini/data/meta")

	gossip, _ := cfg.NewSection("gossip")
	gossip.NewKey("bind-address", "192.168.0.1")
	gossip.NewKey("store-bind-port", "8011")
	gossip.NewKey("meta-bind-port", "8010")
	gossip.NewKey("members", `["192.168.0.1:8010", "192.168.0.2:8010", "192.168.0.3:8010"]`)

	buf := new(bytes.Buffer)
	if _, err := cfg.WriteTo(buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
