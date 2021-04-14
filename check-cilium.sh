for context in kind-cluster-a kind-cluster-b; do                                                                                      
    for pod in $(kubectl --context ${context} -n kube-system get pod -l k8s-app=cilium -o jsonpath='{.items[*].metadata.name}'); do
      kubectl --context ${context} -n kube-system exec ${pod} -- cilium status;
    done
  done


# ./mongo my-replica-set-0-svc.mdb.svc.cluster.local:27017
# ./mongo my-replica-set-1-svc.mdb.svc.cluster.local:27017
# ./mongo my-replica-set-2-svc.mdb.svc.cluster.local:27017