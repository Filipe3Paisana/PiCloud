#!/bin/bash

kubectl apply -f WEB/.

kubectl apply -f DB/secret/.
kubectl apply -f DB/master/.
#kubectl apply -f DB/replica/.

kubectl apply -f API/.
# Aguarda alguns segundos para o Kubernetes iniciar os pods
echo "Aguardar que os pods seijam iniciados..."
sleep 5

# Verifica o estado dos pods
echo "Verificar o estado dos pods..."
kubectl get pods

# Verifica o estado dos serviços
echo "Verificar o estado dos serviços..."
kubectl get services

echo "Comando em caso de problemas nos pods:"
echo "kubectl logs <nome-do-pod>"
echo "kubectl describe pod <nome-do-pod>"
