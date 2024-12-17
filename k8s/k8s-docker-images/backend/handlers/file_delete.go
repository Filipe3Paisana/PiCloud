package handlers

import (
    "fmt"
    "net/http"
    "strconv"

    "api/helpers"
    "api/utils"
)

// DeleteFileFragments apaga todos os fragmentos relacionados a um ficheiro específico, incluindo réplicas.
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
    utils.EnableCors(w, r)

    // Validar o método HTTP
    if r.Method != http.MethodDelete {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    // Extraindo o ID do ficheiro a partir dos parâmetros
    fileIDStr := r.URL.Query().Get("file_id")
    if fileIDStr == "" {
        http.Error(w, "ID do ficheiro não fornecido", http.StatusBadRequest)
        return
    }

    fileID, err := strconv.Atoi(fileIDStr)
    if err != nil {
        http.Error(w, "ID do ficheiro inválido", http.StatusBadRequest)
        return
    }

    // Chamar a lógica para deletar o ficheiro e seus fragmentos
    err = helpers.DeleteFileFragments(fileID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Erro ao eliminar o ficheiro: %v", err), http.StatusInternalServerError)
        return
    }

    // Retornar sucesso
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "ficheiro eliminado com sucesso."}`))
}