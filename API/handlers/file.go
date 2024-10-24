package handlers

import (
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "bytes"
    "encoding/hex"
    "encoding/base64"
    "crypto/md5"
    "math"
    "math/rand"
    "time"
    "strconv"

    "api/utils"
    "api/models"
    "api/db"
)

const availability = 0.999
const failureRate = 0.1

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    utils.EnableCors(w, r)
    
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    if r.Method == http.MethodPost {
        fmt.Println("Recebendo arquivo...")
        
        // Obter o arquivo do formulário
        file, fileHeader, err := r.FormFile("file")
        if err != nil {
            http.Error(w, "Erro ao obter o arquivo", http.StatusBadRequest)
            return
        }
        defer file.Close()

        // Ler o conteudo do ficheiro e calcular o tamanho
        fileContent, err := io.ReadAll(file)
        if err != nil {
            http.Error(w, "Erro ao ler o arquivo", http.StatusInternalServerError)
            return
        }
        fileSize := int64(len(fileContent))
        fileName := fileHeader.Filename

        const maxFileSize = 1000 * 1024 * 1024 // 100MB
        if fileSize > maxFileSize {
            http.Error(w, "Arquivo excede o tamanho máximo permitido de 10MB", http.StatusBadRequest)
            return
        }

        userID, err := utils.ExtractUserIDFromJWT(r) // Extrair ID JWT
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }
        
        // Informações do arquivo na base de dados
        fileID, err := saveFileInfo(fileName, fileSize, userID) 
        if err != nil {
            http.Error(w, fmt.Sprintf("Erro ao salvar informações do arquivo: %v", err), http.StatusInternalServerError)
            return 
        }

        fmt.Printf("ID do Arquivo salvo: %d\n", fileID)

        numberOfFragments := calculateNumberOfFragments(fileSize)
        fmt.Printf("Número de fragmentos: %d\n", numberOfFragments)

        
        err = testFragmentAndReassemble(fileContent, fileSize, numberOfFragments)
        if err != nil {
            http.Error(w, fmt.Sprintf("Erro ao testar a integridade do arquivo: %v", err), http.StatusInternalServerError)
            return
        }

        // Fragmentar o arquivo
        fragments, err := fragmentFile(fileContent, fileSize, numberOfFragments)
        if err != nil {
            http.Error(w, fmt.Sprintf("Erro ao fragmentar o arquivo: %v", err), http.StatusInternalServerError)
            return
        }

        
        for i, fragment := range fragments {
            // Calcular hash MD5 para verificar a integridade
            hash := md5.Sum(fragment)
            hashString := hex.EncodeToString(hash[:])

            // Converter para base64 para visualizar o conteúdo de forma legível
            encoded := base64.StdEncoding.EncodeToString(fragment)
            if len(encoded) > 20 {
                encoded = encoded[:20] + "..." // Mostrar apenas os primeiros 20 caracteres
            }

            fmt.Printf("Fragmento %d: Tamanho = %d bytes, Hash MD5 = %s, Conteúdo (base64) = %s\n", i+1, len(fragment), hashString, encoded)
            saveFragmentInfo(fileID, i+1, hashString) 

        }
        numberOfNodes := GetNumberOfOnlineNodes()
        replicationFactor := calcReplicationFactor(availability, failureRate, numberOfNodes)
        fmt.Printf("Fator de Replicação: %d\n", replicationFactor)

        err = distributeFragments(fileID, numberOfFragments, fragments)
        if err != nil {
            http.Error(w, fmt.Sprint("Erro ao distribuir fragmentos pelos nodes: %v", err), http.StatusInternalServerError)
            return
        }

        // // Enviar o arquivo para o node
        // err = sendFileToNode(file, fileHeader.Filename)
        // if err != nil {
        //     http.Error(w, fmt.Sprintf("Erro ao enviar arquivo para o node: %v", err), http.StatusInternalServerError)
        //     return
        // }
        
        // Responder com sucesso
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "Arquivo enviado com sucesso para o node."}`))
        return
    }
    
    // Responder com erro para métodos não permitidos
    http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
}

func saveFileInfo(name string, size int64, userID int) (int, error) {
    var fileID int
    query := "INSERT INTO Files (name, size, user_id) VALUES ($1, $2, $3) RETURNING id"
    err := db.DB.QueryRow(query, name, size, userID).Scan(&fileID) // Usando db.DB para acessar a instância do banco de dados
    if err != nil {
        return 0, err
    }
    return fileID, nil
}

func calculateNumberOfFragments(fileSize int64) int {
    if fileSize <= 0 {
        return 0 // Caso o arquivo esteja vazio ou com tamanho inválido
    }
    
    const targetFragmentSize = 1024 * 1024 // 1MB
    numberOfFragments := int(fileSize / targetFragmentSize)
    if fileSize%targetFragmentSize != 0 {
        numberOfFragments++ // Adiciona mais um fragmento para o resto do arquivo
    }
    return numberOfFragments
}

func fragmentFile(fileContent []byte, fileSize int64, numFragments int) ([][]byte, error) {
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

func testFragmentAndReassemble(fileContent []byte, fileSize int64, numFragments int) error {
    // Fragmentar o arquivo
    fragments, err := fragmentFile(fileContent, fileSize, numFragments)
    if err != nil {
        return fmt.Errorf("Erro ao fragmentar o arquivo: %v", err)
    }

    // Reconstituir o arquivo a partir dos fragmentos
    var reassembledContent []byte
    for _, fragment := range fragments {
        reassembledContent = append(reassembledContent, fragment...)
    }

    // Verificar se o conteúdo reconstituído é igual ao conteúdo original
    if !bytes.Equal(fileContent, reassembledContent) {
        return fmt.Errorf("O arquivo reconstituído não é idêntico ao original")
    }

    fmt.Println("Teste bem-sucedido: o arquivo foi fragmentado e reconstituído corretamente")
    return nil
}

func saveFragmentInfo(fileID int, fragmentOrder int, hash string) error {
    query := "INSERT INTO FileFragments (file_id, fragment_hashes, fragment_order) VALUES ($1, $2, $3)"
    _, err := db.DB.Exec(query, fileID, hash, fragmentOrder)
    if err != nil {
        return err
    }
    return nil
}

func calcReplicationFactor(availability float64, failureRate float64, numberOfNodes int) int {
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

func distributeFragments(fileID int, numberOfFragments int, fragments [][]byte) error {
    availableNodes := GetOnlineNodesList()
    if len(availableNodes) == 0 {
        return fmt.Errorf("Nenhum nó disponível para distribuir os fragmentos")
    }

    replicationFactor := calcReplicationFactor(availability, failureRate, len(availableNodes))

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
                fmt.Printf("Erro ao salvar informações de distribuição para o fragmento %d no nó %s: %v\n", i, node.NodeAddress, err)
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

    // Inserir na tabela FragmentLocation para registrar onde o fragmento foi armazenado
    query := "INSERT INTO FragmentLocation (fragment_id, node_id) VALUES ($1, $2)"
    _, err = db.DB.Exec(query, fragmentID, nodeID)
    if err != nil {
        return fmt.Errorf("erro ao salvar informações de localização do fragmento: %v", err)
    }
    return nil
}





















func sendFileToNode(file multipart.File, filename string) error {
    var body bytes.Buffer
    writer := multipart.NewWriter(&body)

    part, err := writer.CreateFormFile("fragment", filename)
    if err != nil {
        return err
    }

    // Copiar o conteúdo do arquivo para a requisição
    _, err = io.Copy(part, file)
    if err != nil {
        return err
    }

    err = writer.Close()
    if err != nil {
        return err
    }

    // Enviar a requisição POST para o node
    nodeURL := "http://node1:8082/fragments/upload" // URL do Node
    req, err := http.NewRequest("POST", nodeURL, &body)
    if err != nil {
        return err
    }

    // Adicionar o cabeçalho de tipo de conteúdo
    req.Header.Set("Content-Type", writer.FormDataContentType())

    // Fazer a requisição
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("falha ao enviar arquivo para o node: %s", resp.Status)
    }

    return nil
}
