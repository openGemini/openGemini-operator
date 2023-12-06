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
	Domain      string `toml:"domain"`
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
	Domain          string `toml:"domain"`
}

type LoggingConfig struct {
	Path string `toml:"path"`
}

type GossipConfig struct {
	Enabled       bool     `toml:"enabled"`
	BindAddress   string   `toml:"bind-address"`
	StoreBindPort int      `toml:"store-bind-port"`
	MetaBindPort  int      `toml:"meta-bind-port"`
	Members       []string `toml:"members"`
}

func NewBaseConfiguration(cluster *opengeminiv1.OpenGeminiCluster) (string, error) {
	metaJoinAddrs := []string{}
	metaGossipAddrs := []string{}
	for i := 0; i < int(*cluster.Spec.Meta.Replicas); i++ {
		metaInstance := naming.GenerateMetaInstance(cluster, i)
		host := naming.GenerateInstanceFQDN(metaInstance.Name, metaInstance.Namespace)
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
			Domain:      "<SQL_DOMAIN>",
		},
		Meta: MetaConfig{
			BindAddress:     "<HOST_IP>:8088",
			HttpBindAddress: "<HOST_IP>:8091",
			RpcBindAddress:  "<HOST_IP>:8092",
			Dir:             naming.DataMountPath + "/meta",
			Domain:          "<META_DOMAIN>",
		},
		Data: DataConfig{
			StoreIngestAddr: "<HOST_IP>:8400",
			StoreSelectAddr: "<HOST_IP>:8401",
			StoreDataDir:    naming.DataMountPath + "/data",
			StoreWalDir:     naming.DataMountPath + "/wal",
			StoreMetaDir:    naming.DataMountPath + "/meta",
			Domain:          "<STORE_DOMAIN>",
		},
		Logging: LoggingConfig{
			Path: naming.DataMountPath + "/logs",
		},
		Gossip: GossipConfig{
			Enabled:       true,
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

func UpdateConfig(tmpl, newConf string) (string, error) {
	var tmplToml map[string]interface{}
	var newToml map[string]interface{}

	_, err := toml.Decode(tmpl, &tmplToml)
	if err != nil {
		return "", err
	}

	_, err = toml.Decode(newConf, &newToml)
	if err != nil {
		return "", err
	}

	mapDeepCopy(newToml, tmplToml)

	var buf bytes.Buffer
	err = toml.NewEncoder(&buf).Encode(newToml)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func mapDeepCopy(dst, src map[string]interface{}) {
	for k, v := range src {
		if inv, ok := v.(map[string]interface{}); !ok {
			dst[k] = v
		} else {
			if _, ok := dst[k]; !ok {
				dst[k] = make(map[string]interface{})
			}
			mapDeepCopy(dst[k].(map[string]interface{}), inv)
		}
	}
}

// TODO: delete these codes below
func Merge(data ...string) (string, error) {
	output := make(map[string]interface{})
	for _, dt := range data {
		var tmp map[string]interface{}
		_, err := toml.Decode(dt, &tmp)
		if err != nil {
			return "", fmt.Errorf("error in '%s': %w", dt, err)
		}

		output = mergeMaps(output, tmp)
	}

	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(output)
	if err != nil {
		return "", fmt.Errorf("encode configuration to string failed. err: %w", err)
	}
	return buf.String(), nil
}

func mergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Merge keys from both maps
	for key, value1 := range map1 {
		if value2, ok := map2[key]; ok {
			// Check if both values are maps
			if subMap1, isSubMap1 := value1.(map[string]interface{}); isSubMap1 {
				if subMap2, isSubMap2 := value2.(map[string]interface{}); isSubMap2 {
					// Recursive merge for sub-maps
					result[key] = mergeMaps(subMap1, subMap2)
					continue
				}
			}

			// Check if both values are arrays
			if slice1, isSlice1 := value1.([]interface{}); isSlice1 {
				if slice2, isSlice2 := value2.([]interface{}); isSlice2 {
					// Combine arrays, remove duplicates
					result[key] = mergeArrays(slice1, slice2)
					continue
				}
			}

			// Default: value2 overwrites value1
			result[key] = value2
		} else {
			result[key] = value1
		}
	}

	// Add keys that only exist in map2
	for key, value2 := range map2 {
		if _, exists := map1[key]; !exists {
			result[key] = value2
		}
	}

	return result
}

func mergeArrays(slice1, slice2 []interface{}) []interface{} {
	merged := make([]interface{}, len(slice1)+len(slice2))

	copy(merged, slice1)
	for _, value := range slice2 {
		// Add only unique values from slice2
		if !containsValue(merged, value) {
			merged = append(merged, value)
		}
	}

	return merged
}

func containsValue(slice []interface{}, value interface{}) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
