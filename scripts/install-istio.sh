#! /usr/bin/env sh

clusters=("kind-cluster-a" "kind-cluster-b")

ISTIO_VERSION="${ISTIO_VERSION:-1.9.0}"
rm -rf certs

echo "Using Istio version ${ISTIO_VERSION}"

if [ ! -d "istio-${ISTIO_VERSION}" ]; then
    echo -n "Istio version ${ISTIO_VERSION} not available, downloading..."
    curl -sL  https://istio.io/downloadIstio | ISTIO_VERSION=${ISTIO_VERSION} sh - > /dev/null 2>&1

    if [ -d "istio-${ISTIO_VERSION}" ]; then
      echo "OK"
    else
      echo "FAILED"
    fi
fi

echo "Installing Istio in a fully inter-connected multi-primary setup to contexts:"
for cluster in ${clusters[@]}; do
    echo "  ${cluster}"
done

echo "Generate root-ca and cluster certs"
mkdir certs
pushd certs
echo "Generate root-ca"
make -f ../istio-${ISTIO_VERSION}/tools/certs/Makefile.selfsigned.mk  root-ca
for cluster in ${clusters[@]}; do
  echo "Removing cert for cluster ${cluster}"
  rm -rf certs/${cluster}
  echo "Regenerate cert for cluster ${cluster}"
  make -f ../istio-${ISTIO_VERSION}/tools/certs/Makefile.selfsigned.mk  "${cluster}-cacerts"
done
popd


for cluster in ${clusters[@]}; do
    k8sContext="${cluster}"

    echo "Initializing the Istio Operator Controller on ${cluster}"

    kubectl --context ${k8sContext} delete namespace istio-system istio-operator --ignore-not-found=true
    kubectl --context ${k8sContext} create namespace istio-system
    kubectl --context ${k8sContext} create secret generic cacerts -n istio-system \
                    --from-file=certs/${cluster}/ca-cert.pem \
                    --from-file=certs/${cluster}/ca-key.pem \
                    --from-file=certs/${cluster}/root-cert.pem \
                    --from-file=certs/${cluster}/cert-chain.pem

    istio-${ISTIO_VERSION}/bin/istioctl operator init --context ${k8sContext}

    operatorConfig=$(cat <<EOF
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  meshConfig:
    defaultConfig:
      proxyMetadata:
        ISTIO_META_DNS_AUTO_ALLOCATE: "true"
        ISTIO_META_DNS_CAPTURE: "true"
  values:
    global:
      meshID: mesh1
      multiCluster:
        clusterName: ${cluster}
      network: network1
EOF
)     
    echo "Configuring istio with config: ${operatorConfig}"
    echo "${operatorConfig}" | istio-${ISTIO_VERSION}/bin/istioctl install --context="${k8sContext}" -y -f -
done

echo "Installing multi-primary"

for sourceCluster in ${clusters[@]}; do

    for targetCluster in ${clusters[@]}; do

        if [ "${sourceCluster}" = "${targetCluster}" ]; then
            continue
        fi

        
        if [[ $sourceCluster == kind* ]] ; then
            dockerContainerWithKind=${sourceCluster#"kind-"}

            echo "Adding the istio discovery from ${sourceCluster} to ${targetCluster}"
            sourceClusterControlPlaneAPIIP=$(docker inspect ${dockerContainerWithKind}-control-plane | jq .[].NetworkSettings.Networks.kind.IPAddress -r)
            echo "Using ${sourceClusterControlPlaneAPIIP} for the address of the origin cluster"
            "istio-${ISTIO_VERSION}/bin/istioctl" x create-remote-secret \
                --context="${sourceCluster}" \
                --name="${targetCluster}-to-${sourceCluster}" | \
                sed -E 's!server:.*!server: https://'"${sourceClusterControlPlaneAPIIP}"':6443!' | \
                kubectl apply -f - --context="${targetCluster}"
        else 

            echo "Adding the istio discovery from ${sourceCluster} to ${targetCluster}"
            echo "Using the original kubecontext addres for the address of the origin cluster"
            "istio-${ISTIO_VERSION}/bin/istioctl" x create-remote-secret \
                --context="${sourceCluster}" \
                --name="${sourceCluster}" | \
                kubectl apply -f - --context="${targetCluster}"
        fi

    done
    
done
