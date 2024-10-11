#!/bin/bash

# Aplica o ConfigMap atualizado para NGINX
echo "Aplicando ConfigMap para NGINX..."
kubectl apply -f k8s-deployment/nginx/nginx-configmap.yaml
#kubectl apply -f k8s-deployment/web/web-nginx-configmap.yaml

# Aplica todos os ficheiros de YAML para os deployments
echo "Aplicando os Deployments do Kubernetes..."
kubectl apply -f k8s-deployment/api/api-deployment.yaml
kubectl apply -f k8s-deployment/node/node-deployment.yaml
kubectl apply -f k8s-deployment/web/web-deployment.yaml
kubectl apply -f k8s-deployment/nginx/nginx-deployment.yaml
kubectl apply -f k8s-deployment/db/db-statefulset.yaml
kubectl apply -f k8s-deployment/db/postgres-secret.yaml
kubectl apply -f k8s-deployment/grafana/grafana-deployment.yaml

# Aguarda alguns segundos para o Kubernetes iniciar os pods
echo "Aguarde enquanto os pods são iniciados..."
sleep 5

# Verifica o estado dos pods
echo "Verificando o estado dos pods..."
kubectl get pods

# Verifica o estado dos serviços
echo "Verificando o estado dos serviços..."
kubectl get services

echo "Se houver algum problema com os pods, use os seguintes comandos para verificar os logs:"
echo "kubectl logs <nome-do-pod>"
echo "kubectl describe pod <nome-do-pod>"
