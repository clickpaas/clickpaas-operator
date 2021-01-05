#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})


bash ../vendor/k8s.io/code-generator/generate-groups.sh "all"   l0calh0st.cn/clickpaas-operator/pkg/client   l0calh0st.cn/clickpaas-operator/pkg/apis   middleware:v1alpha1   --output-base ${SCRIPT_ROOT}/../../../   --go-header-file "${SCRIPT_ROOT}/boilerplate.go.txt"
