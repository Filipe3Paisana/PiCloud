package handlers

import (
	"fmt"
    "io"
    "mime/multipart"
    "net/http"
    "bytes"

	"api/utils"
)

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

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    utils.EnableCors(w, r)

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method == http.MethodPost {
        fmt.Println("Recebendo arquivo...")

        // Tamanho máximo do arquivo (10MB)
        err := r.ParseMultipartForm(10 << 20)
        if err != nil {
            http.Error(w, "Erro ao processar o arquivo", http.StatusBadRequest)
            return
        }

        // Obter o arquivo do formulário
        file, fileHeader, err := r.FormFile("file")
        if err != nil {
            http.Error(w, "Erro ao obter o arquivo", http.StatusBadRequest)
            return
        }
        defer file.Close()

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
