# Deployment
An example deployment for kubernetes using the docker production image. To deploy use the following order
```
kubectl apply -f srodb.yaml
kubectl apply -f srokeys.yaml
kubectl apply -f kubernetes.yaml
```

An example Istio `VirtualService` has also been created for easy deployment to Istio.