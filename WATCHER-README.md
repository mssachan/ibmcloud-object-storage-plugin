## Steps to deploy Persistent Volume watcher for updating firewall rules

1. Build Watcher image

```
docker build -t provisioner-builder --pull -f images/watcher/Dockerfile.builder .
docker run provisioner-builder /bin/true
docker cp `docker ps -q -n=1`:/root/ca-certs.tar.gz ./
docker cp `docker ps -q -n=1`:/root/watcher.tar.gz ./
docker build \
     --build-arg git_commit_id=${GIT_COMMIT_SHA} \
     --build-arg git_remote_url=${GIT_REMOTE_URL} \
     --build-arg build_date=${BUILD_DATE} \
     -t <image>:<tag> -f ./images/provisioner/Dockerfile .
rm -f watcher.tar.gz
rm -f ca-certs.tar.gz
```

2. Push watcher image to registry
```
docker push <image>:<tag>
```

3. Deploy watcher pod
```
kubectl apply -f deploy/watcher-sa.yaml
kubectl apply -f deploy/watcher.yaml
```

4. Verification
```
kubectl get pods -n kube-system | grep ibmcloud-object-storage-plugin-pv-watcher
kubectl logs -n kube-system <watcher_pod_name>
```
