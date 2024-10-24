package handlers

import (
    "fmt"
    "net/http"
    "strconv"

    "api/helpers"
    "api/utils"
)

// DeleteFileFragments apaga todos os fragmentos relacionados a um arquivo específico, incluindo réplicas.
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
    utils.EnableCors(w, r)

    // Validar o método HTTP
    if r.Method != http.MethodDelete {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    // Extraindo o ID do arquivo a partir dos parâmetros
    fileIDStr := r.URL.Query().Get("file_id")
    if fileIDStr == "" {
        http.Error(w, "ID do arquivo não fornecido", http.StatusBadRequest)
        return
    }

    fileID, err := strconv.Atoi(fileIDStr)
    if err != nil {
        http.Error(w, "ID do arquivo inválido", http.StatusBadRequest)
        return
    }

    // Chamar a lógica para deletar o arquivo e seus fragmentos
    err = helpers.DeleteFileFragments(fileID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Erro ao deletar o arquivo: %v", err), http.StatusInternalServerError)
        return
    }

    // Retornar sucesso
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "Arquivo deletado com sucesso."}`))
}