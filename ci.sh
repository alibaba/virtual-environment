#!/bin/bash

# This is a simple automated smoke testing with assumptions:
# 1. kubectl and istioctl has configured to kubernetes cluster properly
# 2. user has push authority to ${image} repository (or you could change it to other name)
# 3. VirtualEnvironment CRD has been installed to cluster (with `kubectl apply -f deploy/crds/*_crd.yaml`)
#
# Usage: ci.sh [<name-of-temporary-image-tag>] [<name-of-temporary-namespace>]

# Parameters
if [[ "${1}" =~ ^[A-Z]{1,}$ ]]; then
    action="${1}"
    shift
fi
tag="${1:-ci}"
ns="${2:-virtual-env-ci}"

# Constants
operator_name="virtual-env-operator"
image="virtualenvironment/${operator_name}"
full_image_name="${image}:${tag}"
echo "---- Begin CI Test ----"

# Jump to specified code location
function goto
{
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
elif [[ ${action} = "CLEAN" ]]; then
    goto CLEAN_UP_ANCHOR
fi

# >>>>>>> BUILD_ANCHOR:

# Generate temporary operator image
operator-sdk build --go-build-args "-o build/_output/bin/${operator_name}" ${full_image_name}
if [[ ${?} != 0 ]]; then
    echo "Build failed !!!"
    exit -1
fi
docker push ${full_image_name}
echo "---- Build OK ----"

# >>>>>>> DEPLOY_ANCHOR:

# Create temporary namespace and put operator into it
kubectl create namespace ${ns}
for f in deploy/*.yaml; do
    cat $f | sed "s#${image}:[^ ]*#${full_image_name}#g" | kubectl apply -n ${ns} -f -
done
echo "---- Operator deployment ready ----"

# Deploy demo apps
kubectl create -n ${ns} deployment sleep --image=virtualenvironment/sleep --dry-run -o yaml \
        | istioctl kube-inject -f - | kubectl apply -n ${ns} -f -
examples/deploy/app.sh apply ${ns}

# Wait for apps ready
for i in `seq 50`; do
    count=`kubectl get -n $ns pods | awk '{print $3}' | grep 'Running' | wc -l`
    if [[ ${count} -eq 9 ]]; then
        break
    fi
    echo "waiting ... ${i} (count: ${count})"
    sleep 10s
done
if [[ ${count} -ne 9 ]]; then
    echo "Apps deployment failed !!!"
    exit -1
fi
echo "---- Apps deployment ready ----"

# >>>>>>> TEST_ANCHOR:

# Call service and format response
function invoke_api()
{
    header="${1}"
    kubectl exec -n ${ns} $(kubectl get -n ${ns} pod -l app=sleep -o jsonpath='{.items[0].metadata.name}') -c sleep \
                 -- curl -s -H "ali-env-mark: ${header}" app-js:8080/demo | sed 'N;N;s/\n/, /g'
}

# Check response with expectation
function check_result()
{
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
res=$(invoke_api dev-proj1)
check_result "$res" "[springboot @ dev-proj1] <-dev-proj1, [go @ dev] <-dev-proj1, [node @ dev-proj1] <-dev-proj1"
echo "passed: case 1"

res=$(invoke_api dev-proj1-feature1)
check_result "$res" "[springboot @ dev-proj1-feature1] <-dev-proj1-feature1, [go @ dev] <-dev-proj1-feature1, [node @ dev-proj1] <-dev-proj1-feature1"
echo "passed: case 2"

res=$(invoke_api dev-proj2)
check_result "$res" "[springboot @ dev] <-dev-proj2, [go @ dev-proj2] <-dev-proj2, [node @ dev] <-dev-proj2"
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
