#!/bin/bash
# Prepare version number for next release

VERSION=$(cat Makefile | grep '^VERSION \?' | sed -e 's/.* \(v[0-9.]*\)/\1/g')
echo "Current version is: ${VERSION}"
read -p "Next version should be: " NEXT

for f in Makefile; do
    sed -i '' "s/= ${VERSION}/= ${NEXT}/" $f
done
for f in deploy/global/ktenv_webhook.yaml deploy/ktenv_operator.yaml docs/zh-cn/doc/private-repo.md; do
    sed -i '' "s/:${VERSION}/:${NEXT}/" $f
done
for f in docs/zh-cn/doc/deployment.md docs/en-us/doc/deployment.md; do
    sed -i '' "s/${VERSION}/${NEXT}/g" $f
done
