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


    "api/utils"
    "api/db"
)

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

        const maxFileSize = 100 * 1024 * 1024 // 10MB
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

        // Enviar o arquivo para o node
        err = sendFileToNode(file, fileHeader.Filename)
        if err != nil {
            http.Error(w, fmt.Sprintf("Erro ao enviar arquivo para o node: %v", err), http.StatusInternalServerError)
            return
        }
        
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

func calcReplicationFactor(numberOfNodes int) int {
    // Implementar a função para calcular o fator de replicação
    return 0
}

func replicateFragment(fragment multipart.File, filename string) error {
    // Implementar a função para replicar o fragmento em outros nodes
    return nil
}

func saveReplicaInfo(fileID int, fragmentNumber int, nodeID int) error {
    // Implementar a função para salvar as informações da réplica na base de dados
    return nil
}

func distributeFragments(fileID int, numberOfFragments int) error {
    // Implementar a função para distribuir os fragmentos entre os nodes
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
