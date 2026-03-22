#!/bin/sh

CONFIG_PATH="./config/config.yaml"

function read_key() {
    local key="$2"
    local flag=0
    while read -r LINE; do
        if [[ $flag == 0 ]]; then
            if [[ "$LINE" == *"$key:"* ]]; then
                if [[ "$LINE" == *" "* ]]; then
                    echo "$LINE" | awk -F " " '{print $2}'
                    return
                else
                    continue
                fi
            fi
        fi
    done < "$1"
}

# 若已在环境中设置 ETCD_ADDR（例如 compose 中指向 etcd:2379），则不再覆盖
if [ -z "$ETCD_ADDR" ]; then
  export ETCD_ADDR=$(read_key $CONFIG_PATH "etcd-addr")
fi

sh ./output/${service}/bootstrap.sh