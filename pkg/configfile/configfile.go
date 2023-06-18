package configfile

import (
	"bytes"
	"fmt"
	"strings"

	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	"github.com/openGemini/openGemini-operator/pkg/naming"
	"gopkg.in/ini.v1"
)

func NewConfiguration(cluster *opengeminiv1.GeminiCluster) (string, error) {
	metaJoinAddrs := []string{}
	metaGossipAddrs := []string{}
	for i := 0; i < int(*cluster.Spec.Meta.Replicas); i++ {
		host := naming.GenerateMetaHeadlessSvc(cluster, i)
		metaJoinAddrs = append(metaJoinAddrs, fmt.Sprintf(`"%s:8092"`, host))
		metaGossipAddrs = append(metaGossipAddrs, fmt.Sprintf(`"%s:8010"`, host))
	}
	metaJoinAddrsStr := "[" + strings.Join(metaJoinAddrs, ", ") + "]"
	metaGossipAddrsStr := "[" + strings.Join(metaGossipAddrs, ", ") + "]"

	cfg := ini.Empty()

	common, _ := cfg.NewSection("common")
	common.NewKey("meta-join", metaJoinAddrsStr)

	http, _ := cfg.NewSection("http")
	http.NewKey("bind-address", `"<HOST_IP>:8086"`)

	meta, _ := cfg.NewSection("meta")
	meta.NewKey("bind-address", `"<HOST_IP>:8088"`)
	meta.NewKey("http-bind-address", `"<HOST_IP>:8091"`)
	meta.NewKey("rpc-bind-address", `"<HOST_IP>:8092"`)
	meta.NewKey("dir", `"/ogdata/meta"`)
	meta.NewKey("domain", `"<META_DOMAIN>"`)

	data, _ := cfg.NewSection("data")
	data.NewKey("store-ingest-addr", `"<HOST_IP>:8400"`)
	data.NewKey("store-select-addr", `"<HOST_IP>:8401"`)
	data.NewKey("store-data-dir", `"/ogdata/data"`)
	data.NewKey("store-wal-dir", `"/ogdata/wal"`)
	data.NewKey("store-meta-dir", `"/ogdata/meta"`)

	logging, _ := cfg.NewSection("logging")
	logging.NewKey("path", `"/ogdata/logs"`)

	gossip, _ := cfg.NewSection("gossip")
	gossip.NewKey("bind-address", `"<HOST_IP>"`)
	gossip.NewKey("store-bind-port", "8011")
	gossip.NewKey("meta-bind-port", "8010")
	gossip.NewKey("members", metaGossipAddrsStr)

	buf := new(bytes.Buffer)
	if _, err := cfg.WriteTo(buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
