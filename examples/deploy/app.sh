#!/usr/bin/env bash

basepath=$(cd `dirname $0`; pwd)

action=${1}
namespace=${2}
if [[ "${namespace}" = "" ]]; then
    namespace=`kubectl config get-contexts | grep '^\*' | awk '{print $NF}'`
fi

if [[ "${action}" = "" ]]; then
    echo "action parameter required"
    exit -1
fi

function apply_deployment()
{
    cat ${basepath}/${s}.yaml | sed -e "s/service-name-env-placeholder/${s}-${e}/g" \
                        -e "s/service-name-placeholder/${s}/g" \
                        -e "s/app-env-placeholder/${e}/g" \
                  | istioctl kube-inject -f - \
                  | kubectl ${action} -n ${namespace} -f -
}

# Service
for s in app-js app-go app-java; do
    cat ${basepath}/service.yaml | sed -e "s/service-name-placeholder/${s}/g" | kubectl ${action} -n ${namespace} -f -
done

# Deployment
e=dev
for s in app-js app-go app-java; do
    apply_deployment
done
e='dev.proj1'
for s in app-js app-java; do
    apply_deployment
done
e='dev.proj2'
for s in app-go; do
    apply_deployment
done
e='dev.proj1.feature1'
for s in app-java; do
    apply_deployment
done

kubectl ${action} -n ${namespace} -f ${basepath}/virtualenv.yaml
