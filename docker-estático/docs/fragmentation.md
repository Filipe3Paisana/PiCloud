# Implementação de Redundância em Fragmentação de ficheiros

## Cálculo do Número de Cópias

Para um sistema com `N` nodes e uma tolerância de falha de 60%, podemos usar a seguinte fórmula para calcular o número de cópias:

**Número de cópias** = ⌈N/(1 - tolerância de falhas)⌉

Por exemplo, se `N` é 10 e a tolerância de falhas é 60% (0,6), o número de cópias necessárias para cada fragmento seria:

**Número de cópias** = ⌈10/0,4⌉ = 25 cópias

## Estrutura de Redundância

1. **Número Total de Nodes**: Denote como `N` o número total de nodes.
2. **Tolerância a Falhas**: Para permitir a falha de até 60% dos nodes, você deve garantir que existam fragmentos suficientes para reconstruir o ficheiro mesmo que 60% dos nodes falhem.
3. **Número de Cópias**: A fórmula para calcular o número de cópias necessárias pode ser ajustada para ser um pouco mais complexa. Para 60% de falhas, uma abordagem é garantir que haja cópias suficientes para que os 40% restantes possam garantir a recuperação.

## Atualização da Função de Cálculo de Cópias

A função que calcula o número de cópias pode ser implementada da seguinte forma:

```go
import "math"

// Calcula o número de cópias necessárias para garantir a redundância
func calculateCopies(totalNodes int, failureTolerance float64) int {
    if failureTolerance < 0 || failureTolerance >= 1 {
        panic("A tolerância a falhas deve estar entre 0 e 1 (exclusivo).")
    }

    // Número de cópias necessário
    copies := int(math.Ceil(float64(totalNodes) / (1 - failureTolerance)))
    return copies
}
