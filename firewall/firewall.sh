#!/bin/sh

# Limpa regras existentes
iptables -F

# Permite tráfego de entrada na porta 80 e 8080
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -j ACCEPT

# Bloqueia todo o resto
iptables -A INPUT -j DROP

# Permite tráfego de saída
iptables -A OUTPUT -j ACCEPT

# Mantém o container ativo
tail -f /dev/null
