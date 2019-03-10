#!/bin/bash
DIR_PATH=$(dirname $(readlink -f $0))
LOG_DIR="/var/log/golang-proxy-server"

set_ecbackend_config() {
    local conf="/etc/golang-proxy-server/config.yaml"
    sed -i "s#env_host:.*#env_host: $ENV_HOST#g" $conf
    sed -i "s#env_port:.*#env_port: $ENV_PORT#g" $conf
    sed -i "s#env_timeout:.*#env_timeout: $ENV_TIMEOUT#g" $conf

    sed -i "s#external_url:.*#external_url: $EXTERNAL_URL#g" $conf
    sed -i "s#external_method:.*#external_method: $EXTERNAL_METHOD#g" $conf
    sed -i "s#external_limit_per:.*#external_limit_per: $EXTERNAL_LIMIT_PER#g" $conf
    sed -i "s#external_request_timeout:.*#external_request_timeout: $EXTERNAL_REQUEST_TIMEOUT#g" $conf
    sed -i "s#external_request_queue:.*#external_request_queue: $EXTERNAL_REQUEST_QUEUE#g" $conf
}

main() {
    CONFIGUREFILE="/etc/golang-proxy-server/configure.conf"
    if [ ! -f $CONFIGUREFILE ] ; then
        echo "$(date +'[%d/%b/%Y %T]') File not exist : \"$CONFIGUREFILE\""
        exit 1
    fi
    source $CONFIGUREFILE

    echo "$(date +'[%d/%b/%Y %T]') Start Init Backend API Environment" >> $LOG_DIR/golang-proxy-server_system.log
    set_ecbackend_config
    echo "$(date +'[%d/%b/%Y %T]') Start Service" >> $LOG_DIR/golang-proxy-server_system.log
    echo "" >> $LOG_DIR/golang-proxy-server.log    
    echo "" >> $LOG_DIR/golang-proxy-server.log    
    echo "" >> $LOG_DIR/golang-proxy-server.log    
    exec golang-proxy-server 1>>$LOG_DIR/golang-proxy-server.log 2>>$LOG_DIR/golang-proxy-server.err.log
}

main "$@"
