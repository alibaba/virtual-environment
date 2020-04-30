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

hinit() {
    rm -f /tmp/hashmap.$1
}

hput() {
    echo "$2 $3" >> /tmp/hashmap.$1
}

hget() {
    grep "^$2 " /tmp/hashmap.$1 | awk '{ print $2 };'
}

apply_deployment() {
    cat ${basepath}/deployment.yaml | sed -e "s/service-name-env-placeholder/${s}-${e}/g" \
                        -e "s/service-name-placeholder/${s}/g" \
                        -e "s/app-env-placeholder/${e}/g" \
                        -e "s/app-image-placeholder/`hget images ${s}`/g" \
                        -e "s#app-url-placeholder#`hget urls ${s}`#g" \
                  | istioctl kube-inject -f - \
                  | kubectl ${action} -n ${namespace} -f -
}

# Init parameters
hinit images
hput images app-js js-demo-debug
hput images app-go go-demo-debug
hput images app-java java-demo-debug
hinit urls
hput urls app-js http://app-go:8080/demo
hput urls app-go http://app-java:8080/demo
hput urls app-java

# Create 3 kinds of Service
for s in app-js app-go app-java; do
    cat ${basepath}/service.yaml | sed -e "s/service-name-placeholder/${s}/g" | kubectl ${action} -n ${namespace} -f -
done

# Create Deployment for each env mark
e='dev'
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

# Create a VirtualEnvironment
kubectl ${action} -n ${namespace} -f ${basepath}/virtualenv.yaml
