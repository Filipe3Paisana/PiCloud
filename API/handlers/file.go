package handlers

import (
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "bytes"

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
        
        const maxFileSize = 10 * 1024 * 1024 // 10MB
        if fileHeader.Size > maxFileSize {
            http.Error(w, "Arquivo excede o tamanho máximo permitido de 10MB", http.StatusBadRequest)
            return
        }
        
        // Registrar informações do arquivo na base de dados
        fileID, err := saveFileInfo(fileHeader.Filename, fileHeader.Size, 2) // Supondo user_id = 1 para este exemplo
        if err != nil {
            http.Error(w, fmt.Sprintf("Erro ao salvar informações do arquivo: %v", err), http.StatusInternalServerError)
            return 
        }

        fmt.Printf("ID do Arquivo salvo: %d\n", fileID)

        numberOfFragments := calculateNumberOfFragments(fileHeader.Size)
        fmt.Printf("Número de fragmentos: %d\n", numberOfFragments)
        
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

func saveFileInfo(name string, size int64, userID int) (int, error) {
    var fileID int
    query := "INSERT INTO Files (name, size, user_id) VALUES ($1, $2, $3) RETURNING id"
    err := db.DB.QueryRow(query, name, size, userID).Scan(&fileID) // Usando db.DB para acessar a instância do banco de dados
    if err != nil {
        return 0, err
    }
    return fileID, nil
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
