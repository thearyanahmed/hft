#!/bin/bash

# Apply the Kubernetes manifests including Ingress
kubectl apply -f deployment.yaml -f service.yaml -f ingress.yaml

# Wait for the deployment to be ready (optional)
kubectl wait --for=condition=available deployments/hf-test-deployment --timeout=120s

# Get the external IP using Minikube's tunneling feature
# minikube tunnel

# Now, you can access your application using the specified host (hf-app.example.com in this example)

