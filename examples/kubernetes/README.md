# Deploy AiAdmin to Kubernetes

## Prerequisites

- AiAdmin: https://AiAdmin.com/download
- Kubernetes cluster. This example will use Google Kubernetes Engine.

## Steps

1. Create the AiAdmin namespace, deployment, and service

   ```bash
   kubectl apply -f cpu.yaml
   ```

## (Optional) Hardware Acceleration

Hardware acceleration in Kubernetes requires NVIDIA's [`k8s-device-plugin`](https://github.com/NVIDIA/k8s-device-plugin) which is deployed in Kubernetes in form of daemonset. Follow the link for more details.

Once configured, create a GPU enabled AiAdmin deployment.

```bash
kubectl apply -f gpu.yaml
```

## Test

1. Port forward the AiAdmin service to connect and use it locally

   ```bash
   kubectl -n AiAdmin port-forward service/AiAdmin 11434:80
   ```

1. Pull and run a model, for example `orca-mini:3b`

   ```bash
   AiAdmin run orca-mini:3b
   ```