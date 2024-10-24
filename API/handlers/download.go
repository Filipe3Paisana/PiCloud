package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"api/db"
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
	fileContent, fileName, err := receiveFragments(fileID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao reconstruir o arquivo: %v", err), http.StatusInternalServerError)
		return
	}

	sendFile(w, fileName, fileContent)
}


func receiveFragments(fileID int) ([]byte, string, error) {
	// Buscar os fragmentos armazenados nos diferentes nós para o arquivo
	rows, err := db.DB.Query(`
		SELECT f.fragment_order, f.fragment_hashes, n.id
		FROM FileFragments AS f 
		JOIN FragmentLocation AS fl ON f.fragment_id = fl.fragment_id
		JOIN Nodes AS n ON fl.node_id = n.id 
		WHERE f.file_id = $1 AND status = 'online'
		ORDER BY f.fragment_order`, fileID)
	if err != nil {
		return nil, "", fmt.Errorf("erro ao buscar fragmentos na base de dados: %v", err)
	}
	defer rows.Close()

	var fragments [][]byte
	var fileName string

	for rows.Next() {
		var fragmentOrder int
		var fragmentHash string
		var nodeID int

		if err := rows.Scan(&fragmentOrder, &fragmentHash, &nodeID); err != nil {
			return nil, "", fmt.Errorf("erro ao escanear fragmento: %v", err)
		}

		fragment, err := getFragmentsFromNode(fileID, fragmentOrder, nodeID)
		if err != nil {
			return nil, "", fmt.Errorf("erro ao obter fragmento do nó %d: %v", nodeID, err)
		}
		fragments = append(fragments, fragment)
	}

	// Reconstituir o arquivo a partir dos fragmentos
	fileContent, err := reassembleFile(fragments)
	if err != nil {
		return nil, "", fmt.Errorf("erro ao montar o arquivo: %v", err)
	}

	// Recuperar o nome do arquivo original
	fileName, err = getFileName(fileID)
	if err != nil {
		return nil, "", fmt.Errorf("erro ao obter nome do arquivo: %v", err)
	}

	return fileContent, fileName, nil
}

func getFragmentsFromNode(fileID, fragmentOrder int, nodeID int) ([]byte, error) {
	// Construir URL para obter fragmento do nó
	fragmentURL := fmt.Sprintf("http://node%d:8082/fragments/download?file_id=%d&fragment_order=%d", nodeID, fileID, fragmentOrder)

	resp, err := http.Get(fragmentURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao nó %d: %v", nodeID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("falha ao baixar fragmento do nó %d, status: %d", nodeID, resp.StatusCode)
	}

	fragment, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do nó: %v", err)
	}

	return fragment, nil
}

func reassembleFile(fragments [][]byte) ([]byte, error) {
	var fileContent bytes.Buffer

	for _, fragment := range fragments {
		_, err := fileContent.Write(fragment)
		if err != nil {
			return nil, fmt.Errorf("erro ao reconstituir arquivo: %v", err)
		}
	}

	return fileContent.Bytes(), nil
}

func sendFile(w http.ResponseWriter, fileName string, fileContent []byte) {
	// Definir headers para download
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(fileContent)))

	// Escrever o conteúdo do arquivo no ResponseWriter
	_, err := w.Write(fileContent)
	if err != nil {
		http.Error(w, "Erro ao enviar o arquivo", http.StatusInternalServerError)
	}
}

func getFileName(fileID int) (string, error) {
	var fileName string
	err := db.DB.QueryRow("SELECT name FROM Files WHERE id=$1", fileID).Scan(&fileName)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar nome do arquivo: %v", err)
	}
	return fileName, nil
}
