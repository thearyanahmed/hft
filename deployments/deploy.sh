#!/bin/bash

# Apply the Kubernetes manifests including Ingress
kubectl apply -f deployment.yaml -f service.yaml

# Wait for the deployment to be ready (optional)
kubectl wait --for=condition=available deployments/hf-test-app-deployment --timeout=120s

