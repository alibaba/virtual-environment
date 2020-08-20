#!/bin/bash

# This is a simple automated smoke testing with assumptions:
# 1. kubectl and istioctl has configured to kubernetes cluster properly
# 2. user has push authority to the target image repository (you could change image name via parameter)
# 3. VirtualEnvironment CRD has been installed to cluster (https://alibaba.github.io/virtual-environment/#/en-us/doc/deployment)
#
# Usage: ci.sh [<name-of-ci-image>] [<name-of-ci-namespace>]

# Constants
operator_name="virtual-env-operator"
default_image="virtualenvironment/${operator_name}"
default_tag="ci"

# Parameters
if [[ "${1}" =~ ^[A-Z]{1,}$ ]]; then
    action="${1}"
    shift
fi
ci_image="${1}"
if [[ "${ci_image}" = "" || "${ci_image}" = "_" ]]; then
    ci_image="${default_image}:${default_tag}"
fi
ns="${2:-virtual-env-ci}"

echo "> Using image $ci_image"
echo "> Using namespace $ns"

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
if [[ ${action} = "DEPLOY" ]]; then
    goto DEPLOY_ANCHOR
elif [[ ${action} = "TEST" ]]; then
    goto TEST_ANCHOR
elif [[ ${action} = "CLEAN" || ${action} = "DELETE" ]]; then
    goto CLEAN_UP_ANCHOR
fi

# >>>>>>> BUILD_ANCHOR:

# Generate temporary operator image
operator-sdk build --go-build-args "-o build/_output/operator/${operator_name}" --image-build-args "--no-cache" ${ci_image}
if [[ ${?} != 0 ]]; then
    echo "Build failed !!!"
    exit -1
fi
docker push ${ci_image}
echo "---- Build OK ----"

# >>>>>>> DEPLOY_ANCHOR:

# Create temporary namespace and put operator into it
kubectl create namespace ${ns}
for f in deploy/*.yaml; do
    cat $f | sed "s#${default_image}:[^ ]*#${ci_image}#g" | kubectl apply -n ${ns} -f -
done
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
examples/deploy/app.sh delete ${ns}
kubectl delete -n ${ns} deployment sleep
for f in deploy/*.yaml; do kubectl delete -n ${ns} -f ${f}; done
kubectl delete namespace ${ns}
echo "---- Clean up OK ----"
