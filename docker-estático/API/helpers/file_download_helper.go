package helpers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"crypto/md5"
	"encoding/hex"


	"api/db"
)


func ReceiveFragments(fileID int) ([]byte, string, error) {
	// Buscar os fragmentos armazenados nos diferentes nós para o ficheiro
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

		fragment, err := GetFragmentsFromNode(fileID, fragmentOrder, nodeID, fragmentHash)
		if err != nil {
			return nil, "", fmt.Errorf("erro ao obter fragmento do nó %d: %v", nodeID, err)
		}
		fragments = append(fragments, fragment)
	}

	// Reconstituir o ficheiro a partir dos fragmentos TODO verificar integridade do ficheiro
	fileContent, err := ReassembleFile(fragments)
	if err != nil {
		return nil, "", fmt.Errorf("erro ao montar o ficheiro: %v", err)
	}

	// Recuperar o nome do ficheiro original
	fileName, err = GetFileName(fileID)
	if err != nil {
		return nil, "", fmt.Errorf("erro ao obter nome do ficheiro: %v", err)
	}

	return fileContent, fileName, nil
}

func GetFragmentsFromNode(fileID, fragmentOrder int, nodeID int, fragmentHash string) ([]byte, error) {
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

    // Calcular o hash MD5 do fragmento recebido
    hash := md5.Sum(fragment)
    calculatedHash := hex.EncodeToString(hash[:])

    // Verificar se o hash calculado coincide com o armazenado
    if calculatedHash != fragmentHash {
        RemoveFragment(fileID, fragmentOrder)
        return nil, fmt.Errorf("hash do fragmento adulterado no nó %d: esperado %s, recebido %s", nodeID, fragmentHash, calculatedHash)
    }

    return fragment, nil
}
func ReassembleFile(fragments [][]byte) ([]byte, error) {
	var fileContent bytes.Buffer

	for _, fragment := range fragments {
		_, err := fileContent.Write(fragment)
		if err != nil {
			return nil, fmt.Errorf("erro ao reconstituir ficheiro: %v", err)
		}
	}

	return fileContent.Bytes(), nil
}

func SendFile(w http.ResponseWriter, fileName string, fileContent []byte) {
	// Definir headers para download
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(fileContent)))

	// Escrever o conteúdo do ficheiro no ResponseWriter
	_, err := w.Write(fileContent)
	if err != nil {
		http.Error(w, "Erro ao enviar o ficheiro", http.StatusInternalServerError)
	}
}

func GetFileName(fileID int) (string, error) {
	var fileName string
	err := db.DB.QueryRow("SELECT name FROM Files WHERE id=$1", fileID).Scan(&fileName)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar nome do ficheiro: %v", err)
	}
	return fileName, nil
}

func RemoveFragment(fileID, fragmentOrder int) {
	_, err := db.DB.Exec("DELETE FROM FileFragments WHERE file_id=$1 AND fragment_order=$2", fileID, fragmentOrder)
	if err != nil {
		fmt.Printf("Erro ao remover fragmento %d do ficheiro %d: %v", fragmentOrder, fileID, err)
	}
}