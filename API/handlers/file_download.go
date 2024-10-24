package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	
	"api/helpers"
	"api/utils"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w, r)

	// Extraindo o ID do arquivo a partir da URL
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

	// Receber fragmentos e reconstituir o arquivo
	fileContent, fileName, err := helpers.ReceiveFragments(fileID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao reconstruir o arquivo: %v", err), http.StatusInternalServerError)
		return
	}

	helpers.SendFile(w, fileName, fileContent)
}


