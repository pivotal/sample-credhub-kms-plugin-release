#!/bin/bash

function set_bash_error_handling() {
    set -euo pipefail
}

function go_to_project_root_directory() {
    local -r script_dir=$( dirname "${BASH_SOURCE[0]}")

    cd "$script_dir/.."
}

function run_tests() {
    export GOPATH=${PWD}
    pushd src/github.com/pivotal/sample-credhub-kms-plugin >/dev/null
        scripts/setup_dev_grpc_certs.sh
        echo -e "CA is available at $(pwd)/grpc-kms-certs/grpc_kms_ca_cert.pem"
        ginkgo -r -p
        rm -rf grpc-kms-certs
    popd >/dev/null
}

main() {
    set_bash_error_handling
    go_to_project_root_directory
    run_tests
}

main
