package helpers

import (
    "fmt"
    "net"
    "syscall"
)



// Função para obter o endereço IP local
func GetLocalIPAddress() (string, error) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "", err
    }

    for _, addr := range addrs {
        if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
            if ipNet.IP.To4() != nil {
                return ipNet.IP.String(), nil
            }
        }
    }
    return "", fmt.Errorf("não foi possível obter o IP")
}

// Função para obter o uso de disco
func GetDiskUsage(path string) (total uint64, free uint64, err error) {
    var stat syscall.Statfs_t
    err = syscall.Statfs(path, &stat)
    if err != nil {
        return 0, 0, err
    }

    total = stat.Blocks * uint64(stat.Bsize)       // Capacidade total
    free = stat.Bavail * uint64(stat.Bsize)        // Capacidade disponível
    return total, free, nil
}