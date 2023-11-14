package configfile

var EntrypointScript = `#!/bin/env bash
set -e

dataDir=${DATA_DIR:-/ogdata}
mkdir -p "$dataDir"/{conf,logs,data,meta,wal}

configMount=${CONFIG_PATH:-/etc/opengemini/opengemini.conf}
if [[ ! -f $configMount ]]; then
    echo "Configuration file does not exist: $configMount"
    exit 1
fi

config=$dataDir/conf/opengemini.conf
cp -f "$configMount" "$config"

hostIp=$HOST_IP
if [[ -z $hostIp ]]; then
    hostIp=$(/sbin/ifconfig -a | grep inet | grep -v 127.0.0.1 | grep -v inet6 | awk '{print $2;exit}' | tr -d "addr:")
fi
sed -i "s/<HOST_IP>/$hostIp/g" "$config"

metaDomain=$META_DOMAIN
sed -i "s/<META_DOMAIN>/$metaDomain/g" "$config"
sqlDomain=$SQL_DOMAIN
sed -i "s/<SQL_DOMAIN>/$sqlDomain/g" "$config"
storeDomain=$STORE_DOMAIN
sed -i "s/<STORE_DOMAIN>/$storeDomain/g" "$config"

app=ts-"${APP:-sql}"
exec $app --config "$config"
`
