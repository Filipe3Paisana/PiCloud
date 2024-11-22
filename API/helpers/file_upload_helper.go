package helpers

import (
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "bytes"
    "math"
    "math/rand"
    "time"
    "strconv"

    "api/models"
    "api/db"
)


func SaveFileInfo(name string, size int64, userID int) (int, error) {
    var fileID int
    query := "INSERT INTO Files (name, size, user_id) VALUES ($1, $2, $3) RETURNING id"
    err := db.DB.QueryRow(query, name, size, userID).Scan(&fileID) // Usando db.DB para acessar a instância do banco de dados
    if err != nil {
        return 0, err
    }
    return fileID, nil
}

func CalculateNumberOfFragments(fileSize int64) int {
    if fileSize <= 0 {
        return 0 // Caso o ficheiro esteja vazio ou com tamanho inválido
    }
    
    const targetFragmentSize = 1024 * 1024 // 1MB
    numberOfFragments := int(fileSize / targetFragmentSize)
    if fileSize%targetFragmentSize != 0 {
        numberOfFragments++ // Adiciona mais um fragmento para o resto do ficheiro
    }
    return numberOfFragments
}

func FragmentFile(fileContent []byte, fileSize int64, numFragments int) ([][]byte, error) {
    if numFragments <= 0 {
        return nil, fmt.Errorf("Número de fragmentos deve ser maior que zero")
    }

    var fragments [][]byte
    fragmentSize := int(fileSize) / numFragments
    remainder := int(fileSize) % numFragments

    start := 0

    for i := 0; i < numFragments; i++ {
        end := start + fragmentSize
        if i == numFragments-1 {
            // O último fragmento contém o resto
            end += remainder
        }
        fragments = append(fragments, fileContent[start:end])
        start = end
    }

    return fragments, nil
}

func TestFragmentAndReassemble(fileContent []byte, fileSize int64, numFragments int) error {
    // Fragmentar o ficheiro
    fragments, err := FragmentFile(fileContent, fileSize, numFragments)
    if err != nil {
        return fmt.Errorf("Erro ao fragmentar o ficheiro: %v", err)
    }

    // Reconstituir o ficheiro a partir dos fragmentos
    var reassembledContent []byte
    for _, fragment := range fragments {
        reassembledContent = append(reassembledContent, fragment...)
    }

    // Verificar se o conteúdo reconstituído é igual ao conteúdo original
    if !bytes.Equal(fileContent, reassembledContent) {
        return fmt.Errorf("O ficheiro reconstituído não é idêntico ao original")
    }

    fmt.Println("Teste bem-sucedido: o ficheiro foi fragmentado e reconstituído corretamente")
    return nil
}

func SaveFragmentInfo(fileID int, fragmentOrder int, hash string) error {
    query := "INSERT INTO FileFragments (file_id, fragment_hashes, fragment_order) VALUES ($1, $2, $3)"
    _, err := db.DB.Exec(query, fileID, hash, fragmentOrder)
    if err != nil {
        return err
    }
    return nil
}

func CalcReplicationFactor(availability float64, failureRate float64, numberOfNodes int) int {
    if availability <= 0 || availability >= 1 || failureRate <= 0 || failureRate >= 1 || numberOfNodes <= 0 {
        fmt.Println("Parâmetros inválidos fornecidos para o cálculo do fator de replicação.")
        return 0
    }

    // Calcular o fator de replicação usando a fórmula
    numerator := math.Log(1 - availability)
    denominator := math.Log(failureRate / float64(numberOfNodes))

    replicationFactor := numerator / denominator

    // Arredondar para cima para garantir o mínimo necessário de réplicas
    return int(math.Ceil(replicationFactor)) //TODO definir se este é o número de réplicas (sem contar com o original) ou o número de fragmentos total
}

func DistributeFragments(fileID int, numberOfFragments int, fragments [][]byte, replicationFactor int) error {
    availableNodes := GetOnlineNodesList()
    if len(availableNodes) == 0 {
        return fmt.Errorf("Nenhum nó disponível para distribuir os fragmentos")
    }

    for i := 1; i <= numberOfFragments; i++ {
        selectedNodes := SelectNodesForFragment(availableNodes, replicationFactor)

        for _, node := range selectedNodes {
            err := SendFragmentToNode(fileID, i, fragments[i-1], node.NodeID) // Passando o conteúdo real do fragmento
            if err != nil {
                fmt.Printf("Erro ao enviar fragmento %d para o nó %d: %v\n", i, node.NodeID, err)
                continue
            }

            err = SaveDistributionInfo(fileID, i, node.NodeAddress)
            if err != nil {
                fmt.Printf("Erro ao guardar informações de distribuição para o fragmento %d no nó %s: %v\n", i, node.NodeAddress, err)
                continue
            }
        }
    }
    return nil
}

func SelectNodesForFragment(availableNodes []models.Node, replicationFactor int) []models.Node {
    if replicationFactor >= len(availableNodes) {
        // Se o fator de replicação é maior ou igual ao número de nós disponíveis, retorna todos os nós
        return availableNodes
    }

    // Inicializar o gerador de números aleatórios
    rand.Seed(time.Now().UnixNano())

    // Embaralhar a lista de nós para uma seleção aleatória
    rand.Shuffle(len(availableNodes), func(i, j int) {
        availableNodes[i], availableNodes[j] = availableNodes[j], availableNodes[i]
    })

    // Selecionar um número de nós igual ao fator de replicação
    return availableNodes[:replicationFactor]
}

func SendFragmentToNode(fileID int, fragmentOrder int, fragmentContent []byte, nodeID int) error {
    // Construir a URL do nó
    nodeURL := fmt.Sprintf("http://node%d:8082/fragments/upload", nodeID)

    var body bytes.Buffer
    writer := multipart.NewWriter(&body)

    // Adicionar campos ao formulário para passar informações adicionais
    writer.WriteField("file_id", strconv.Itoa(fileID))
    writer.WriteField("fragment_order", strconv.Itoa(fragmentOrder))

    // Adicionar o fragmento ao formulário
    part, err := writer.CreateFormFile("fragment", fmt.Sprintf("file_%dfragment_%d", fileID, fragmentOrder))
    if err != nil {
        return err
    }

    // Escreve o conteúdo real do fragmento
    _, err = io.Copy(part, bytes.NewReader(fragmentContent))
    if err != nil {
        return err
    }

    err = writer.Close()
    if err != nil {
        return err
    }

    // Criar e enviar requisição HTTP POST
    req, err := http.NewRequest("POST", nodeURL, &body)
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("falha ao enviar fragmento para o nó: %d, código de status: %d", nodeID, resp.StatusCode)
    }

    return nil
}

func SaveDistributionInfo(fileID int, fragmentOrder int, nodeAddress string) error {
    // Primeiro, verificar se o fragmento já existe e obter seu fragment_id
    var fragmentID int
    err := db.DB.QueryRow("SELECT fragment_id FROM FileFragments WHERE file_id = $1 AND fragment_order = $2", fileID, fragmentOrder).Scan(&fragmentID)
    if err != nil {
        return fmt.Errorf("erro ao obter fragmento: %v", err)
    }

    // Obter o ID do nó com base no endereço
    var nodeID int
    err = db.DB.QueryRow("SELECT id FROM Nodes WHERE node_address = $1", nodeAddress).Scan(&nodeID)
    if err != nil {
        return fmt.Errorf("erro ao obter nó: %v", err)
    }

    // Inserir na tabela FragmentLocation para registar onde o fragmento foi armazenado
    query := "INSERT INTO FragmentLocation (fragment_id, node_id) VALUES ($1, $2)"
    _, err = db.DB.Exec(query, fragmentID, nodeID)
    if err != nil {
        return fmt.Errorf("erro ao guardar informações de localização do fragmento: %v", err)
    }
    return nil
}