package configfile

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	opengeminiv1 "github.com/openGemini/openGemini-operator/api/v1"
	"github.com/openGemini/openGemini-operator/pkg/naming"
)

type Config struct {
	Common  CommonConfig  `toml:"common"`
	Http    HttpConfig    `toml:"http"`
	Meta    MetaConfig    `toml:"meta"`
	Data    DataConfig    `toml:"data"`
	Logging LoggingConfig `toml:"logging"`
	Gossip  GossipConfig  `toml:"gossip"`
}

type CommonConfig struct {
	MetaJoin []string `toml:"meta-join"`
}

type HttpConfig struct {
	BindAddress string `toml:"bind-address"`
	AuthEnabled bool   `toml:"auth-enabled"`
}

type MetaConfig struct {
	BindAddress     string `toml:"bind-address"`
	HttpBindAddress string `toml:"http-bind-address"`
	RpcBindAddress  string `toml:"rpc-bind-address"`
	Dir             string `toml:"dir"`
	Domain          string `toml:"domain"`
}

type DataConfig struct {
	StoreIngestAddr string `toml:"store-ingest-addr"`
	StoreSelectAddr string `toml:"store-select-addr"`
	StoreDataDir    string `toml:"store-data-dir"`
	StoreWalDir     string `toml:"store-wal-dir"`
	StoreMetaDir    string `toml:"store-meta-dir"`
}

type LoggingConfig struct {
	Path string `toml:"path"`
}

type GossipConfig struct {
	BindAddress   string   `toml:"bind-address"`
	StoreBindPort int      `toml:"store-bind-port"`
	MetaBindPort  int      `toml:"meta-bind-port"`
	Members       []string `toml:"members"`
}

func NewConfiguration(cluster *opengeminiv1.GeminiCluster) (string, error) {
	metaJoinAddrs := []string{}
	metaGossipAddrs := []string{}
	for i := 0; i < int(*cluster.Spec.Meta.Replicas); i++ {
		host := naming.GenerateMetaHeadlessSvc(cluster, i)
		metaJoinAddrs = append(metaJoinAddrs, fmt.Sprintf("%s:8092", host))
		metaGossipAddrs = append(metaGossipAddrs, fmt.Sprintf("%s:8010", host))
	}

	config := Config{
		Common: CommonConfig{
			MetaJoin: metaJoinAddrs,
		},
		Http: HttpConfig{
			BindAddress: "<HOST_IP>:8086",
			AuthEnabled: cluster.GetEnableHttpAuth(),
		},
		Meta: MetaConfig{
			BindAddress:     "<HOST_IP>:8088",
			HttpBindAddress: "<HOST_IP>:8091",
			RpcBindAddress:  "<HOST_IP>:8092",
			Dir:             "/ogdata/meta",
			Domain:          "<META_DOMAIN>",
		},
		Data: DataConfig{
			StoreIngestAddr: "<HOST_IP>:8400",
			StoreSelectAddr: "<HOST_IP>:8401",
			StoreDataDir:    "/ogdata/data",
			StoreWalDir:     "/ogdata/wal",
			StoreMetaDir:    "/ogdata/meta",
		},
		Logging: LoggingConfig{
			Path: "/ogdata/logs",
		},
		Gossip: GossipConfig{
			BindAddress:   "<HOST_IP>",
			StoreBindPort: 8011,
			MetaBindPort:  8010,
			Members:       metaGossipAddrs,
		},
	}

	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(config)
	if err != nil {
		return "", fmt.Errorf("encode configuration to string failed. err: %w", err)
	}
	return buf.String(), nil
}
