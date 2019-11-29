#!/usr/bin/env bash

action=${1}
if [[ "${action}" = "" ]]; then
    echo "action parameter required"
    exit -1
fi

# Service
for s in app-js app-go app-java; do
    cat service.yaml | sed -e "s/service-name-placeholder/${s}/g" | kubectl ${action} -f -
done

# Deployment
e=dev
for s in app-js app-go app-java; do
    cat ${s}.yaml | sed -e "s/service-name-env-placeholder/${s}-${e}/g" -e "s/service-name-placeholder/${s}/g" -e "s/app-env-placeholder/${e}/g" | kubectl ${action} -f -
done
e=dev-proj1
for s in app-js app-java; do
    cat ${s}.yaml | sed -e "s/service-name-env-placeholder/${s}-${e}/g" -e "s/service-name-placeholder/${s}/g" -e "s/app-env-placeholder/${e}/g" | kubectl ${action} -f -
done
e=dev-proj2
for s in app-go; do
    cat ${s}.yaml | sed -e "s/service-name-env-placeholder/${s}-${e}/g" -e "s/service-name-placeholder/${s}/g" -e "s/app-env-placeholder/${e}/g" | kubectl ${action} -f -
done
e=dev-proj1-feature1
for s in app-java; do
    cat ${s}.yaml | sed -e "s/service-name-env-placeholder/${s}-${e}/g" -e "s/service-name-placeholder/${s}/g" -e "s/app-env-placeholder/${e}/g" | kubectl ${action} -f -
done

kubectl ${action} -f virtualenv.yaml
