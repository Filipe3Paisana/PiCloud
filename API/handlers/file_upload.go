package handlers

import (
    "fmt"
    "io"
    "net/http"
    "encoding/hex"
    "encoding/base64"
    "crypto/md5"
    
    "api/helpers"
    "api/utils"
)

const availability = 0.999
const failureRate = 0.2

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
        fileID, err := helpers.SaveFileInfo(fileName, fileSize, userID) 
        if err != nil {
            http.Error(w, fmt.Sprintf("Erro ao salvar informações do arquivo: %v", err), http.StatusInternalServerError)
            return 
        }

        fmt.Printf("ID do Arquivo salvo: %d\n", fileID)

        numberOfFragments := helpers.CalculateNumberOfFragments(fileSize)
        fmt.Printf("Número de fragmentos: %d\n", numberOfFragments)

        
        err = helpers.TestFragmentAndReassemble(fileContent, fileSize, numberOfFragments)
        if err != nil {
            http.Error(w, fmt.Sprintf("Erro ao testar a integridade do arquivo: %v", err), http.StatusInternalServerError)
            return
        }

        // Fragmentar o arquivo
        fragments, err := helpers.FragmentFile(fileContent, fileSize, numberOfFragments)
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
            helpers.SaveFragmentInfo(fileID, i+1, hashString) 

        }
        numberOfNodes := helpers.GetNumberOfOnlineNodes()
        replicationFactor := helpers.CalcReplicationFactor(availability, failureRate, numberOfNodes)
        fmt.Printf("Fator de Replicação: %d\n", replicationFactor)

        err = helpers.DistributeFragments(fileID, numberOfFragments, fragments, replicationFactor)
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

