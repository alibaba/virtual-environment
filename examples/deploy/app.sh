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
    apply_pods "deployment"
}

apply_statefulset() {
    apply_pods "statefulset"
}

apply_pods() {
    type=${1}
    ee=`echo ${e} | sed -e "s/\./-/g"`
    cat ${basepath}/${type}.yaml | sed -e "s/service-name-env-placeholder/${s}-${ee}/g" \
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

# Create required Deployments for each env mark
s='app-js'
for e in 'dev' 'dev.proj1'; do
    apply_deployment
done
s='app-go'
for e in 'dev' 'dev.proj2'; do
    apply_deployment
done

# Create required Deployments for each env mark
s='app-java'
for e in 'dev' 'dev.proj1' 'dev.proj1.feature1'; do
    apply_statefulset
done

# Create a VirtualEnvironment
kubectl ${action} -n ${namespace} -f ${basepath}/virtualenv.yaml
