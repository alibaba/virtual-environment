#!/bin/bash

# This is a simple automated smoke testing with assumptions:
# 1. kubectl and istioctl has configured to kubernetes cluster properly
# 2. user has push authority to the target image repository (you could change image name via parameter)
# 3. VirtualEnvironment CRD has been installed to cluster (https://alibaba.github.io/virtual-environment/#/en-us/doc/deployment)

# Constants
ci_tag="ci"
operator_name="virtual-env-operator"
webhook_name="virtual-env-admission-webhook"
default_operator_image="virtualenvironment/${operator_name}"
default_webhook_image="virtualenvironment/${webhook_name}"
default_ci_namespace="virtual-env-ci"

usage() {
    cat <<EOF
Usage: ci.sh [flags] [tag] [<name-of-ci-namespace>] [<name-of-operator-image>] [<name-of-webhook-image>]
  supported flags:
    --no-cleanup        keep test pod running after all case finish
    --keep-namespace    keep namespace after cleanup resource
    --include-webhook   also build and deploy webhook component
    --help              show this message
  supported tags:
    @DEPLOY             run from deploy step
    @TEST               run from test step
    @CLEAN or @DELETE   run from clean up step
  default config:
    ci namespace: ${default_ci_namespace}
    operator image: ${default_operator_image}:${ci_tag}
    webhook image: ${default_webhook_image}:${ci_tag}
EOF
}

# Configure
skip_cleanup="N"
with_webhook="N"
keep_namespace="N"

# Parameters
for p in ${@}; do
    if [[ "${p}" =~ ^--.*$ ]]; then
        if [[ "${p}" = "--no-cleanup" ]]; then
            skip_cleanup="Y"
        elif [[ "${p}" = "--include-webhook" ]]; then
            with_webhook="Y"
        elif [[ "${p}" = "--keep-namespace" ]]; then
            keep_namespace="Y"
        elif [[ "${p}" = "--help" ]]; then
            usage
            exit 0
        fi
        shift
    fi
    if [[ "${1}" =~ ^@[A-Z]{1,}$ ]]; then
        action="${1#*@}"
        shift
    fi
done
ns="${1}"
if [[ "${ns}" = "" || "${ns}" = "_" ]]; then
    ns="${default_ci_namespace}"
fi
ci_operator_image="${2}"
if [[ "${ci_operator_image}" = "" || "${ci_operator_image}" = "_" ]]; then
    ci_operator_image="${default_operator_image}:${ci_tag}"
fi
ci_webhook_image="${3}"
if [[ "${ci_webhook_image}" = "" ]]; then
    ci_webhook_image="${default_webhook_image}:${ci_tag}"
fi

# Print context
echo "> Using namespace ${ns}"
echo "> Using operator image ${ci_operator_image}"
if [[ "${with_webhook}" = "Y" ]]; then
    echo "> Using webhook image ${ci_webhook_image}"
fi

echo "---- Begin CI Task ----"

# Jump to specified code location
goto() {
    sed_cmd="sed"
    if [[ "$(uname -s)" = "Darwin" ]]; then
        sed_cmd="gsed"
    fi
    label=$1
    cmd=$(${sed_cmd} -n "/^# >\+ $label:/{:a;n;p;ba};" $0 | grep -v ':$')
    eval "$cmd"
    exit
}

# Shortcuts
if [[ "${action}" = "DEPLOY" ]]; then
    goto DEPLOY_ANCHOR
elif [[ "${action}" = "TEST" ]]; then
    goto TEST_ANCHOR
elif [[ "${action}" = "CLEAN" || "${action}" = "DELETE" ]]; then
    goto CLEAN_UP_ANCHOR
fi

# >>>>>>> BUILD_ANCHOR:

# Generate temporary operator image
make build-operator OPERATOR_IMAGE_AND_VERSION=${ci_operator_image}
if [[ ${?} -ne 0 ]]; then
    echo "Build operator failed !!!"
    exit -1
fi
docker push ${ci_operator_image}
if [[ ${?} -ne 0 ]]; then
    echo "Push operator image failed !!!"
    exit -1
fi

if [[ "${with_webhook}" = "Y" ]]; then
    make build-webhook WEBHOOK_IMAGE_AND_VERSION=${ci_webhook_image}
    if [[ ${?} -ne 0 ]]; then
        echo "Build webhook failed !!!"
        exit -1
    fi
    docker push ${ci_webhook_image}
    if [[ ${?} -ne 0 ]]; then
        echo "Push webhook image failed !!!"
        exit -1
    fi
fi

echo "---- Build OK ----"

# >>>>>>> DEPLOY_ANCHOR:

# Create temporary namespace and put operator into it
kubectl create namespace ${ns}
kubectl apply -f deploy/global/ktenv_crd.yaml
for f in deploy/*.yaml; do
    cat $f | sed "s#${default_operator_image}:.*#${ci_operator_image}#g" | kubectl apply -n ${ns} -f -
done
if [[ "${with_webhook}" = "Y" ]]; then
    cat deploy/global/ktenv_webhook.yaml | sed "s#${default_webhook_image}:.*#${ci_webhook_image}#g" | kubectl apply -f -
fi
kubectl label namespaces ${ns} environment-tag-injection=enabled
echo "---- Operator deployment ready ----"

# Deploy demo apps
kubectl create -n ${ns} deployment sleep --image=virtualenvironment/sleep --dry-run -o yaml \
        | istioctl kube-inject -f - | kubectl apply -n ${ns} -f -
examples/deploy/app.sh apply ${ns}

# Wait for apps ready
declare -i expect_count=9
for i in `seq 50`; do
    count=`kubectl get -n ${ns} pods | awk '{print $3}' | grep 'Running' | wc -l`
    if [[ ${count} -eq ${expect_count} ]]; then
        break
    fi
    echo "waiting ... ${i} (count: ${count}/${expect_count})"
    sleep 10s
done
if [[ ${count} -ne ${expect_count} ]]; then
    echo "Apps deployment failed !!!"
    exit -1
fi
echo "---- Apps deployment ready ----"

# >>>>>>> TEST_ANCHOR:

# Call service and format response
invoke_api() {
    header="${1}"
    kubectl exec -n ${ns} $(kubectl get -n ${ns} pod -l app=sleep -o jsonpath='{.items[0].metadata.name}') -c sleep \
                 -- curl -s -H "ali-env-mark: ${header}" app-js:8080/demo | sed 'N;N;s/\n/, /g'
}

# Check response with expectation
check_result() {
    real="${1}"
    expect="${2}"
    if [[ "${real}" != "${expect}" ]]; then
        echo "Test failed !!!"
        echo "Namespace: ${ns}"
        echo "Real response: ${real}"
        echo "Expectation  : ${expect}"
        exit -1
    fi
}

# Do functional check
res=$(invoke_api dev.proj1)
check_result "$res" "[springboot @ dev.proj1] <-dev.proj1, [go @ dev] <-dev.proj1, [node @ dev.proj1] <-dev.proj1"
echo "passed: case 1"

res=$(invoke_api dev.proj1.feature1)
check_result "$res" "[springboot @ dev.proj1.feature1] <-dev.proj1.feature1, [go @ dev] <-dev.proj1.feature1, [node @ dev.proj1] <-dev.proj1.feature1"
echo "passed: case 2"

res=$(invoke_api dev.proj2)
check_result "$res" "[springboot @ dev] <-dev.proj2, [go @ dev.proj2] <-dev.proj2, [node @ dev] <-dev.proj2"
echo "passed: case 3"

res=$(invoke_api dev)
check_result "$res" "[springboot @ dev] <-dev, [go @ dev] <-dev, [node @ dev] <-dev"
echo "passed: case 4"

res=$(invoke_api)
check_result "$res" "[springboot @ dev] <-dev, [go @ dev] <-dev, [node @ dev] <-empty"
echo "passed: case 5"
echo "---- Functional check OK ----"

# >>>>>>> CLEAN_UP_ANCHOR:

# Clean up everything
if [[ "${skip_cleanup}" != "Y" ]]; then
    examples/deploy/app.sh delete ${ns}
    kubectl delete -n ${ns} deployment sleep
    for f in deploy/*.yaml; do kubectl delete -n ${ns} -f ${f}; done
    kubectl label namespaces ${ns} environment-tag-injection-
    if [[ "${keep_namespace}" != "Y" ]]; then
        kubectl delete namespace ${ns}
    fi
    echo "---- Clean up OK ----"
fi
