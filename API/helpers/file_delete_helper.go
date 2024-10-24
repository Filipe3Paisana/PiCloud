package helpers

import (
    "fmt"
    "net/http"
	"io"

	"api/db"
)


func DeleteFileFragments(fileID int) error {
	// Lista os nós disponíveis
	availableNodes := GetOnlineNodesList()
	if len(availableNodes) == 0 {
		return fmt.Errorf("Nenhum nó disponível para deletar os fragmentos")
	}

	// Obter os fragmentos e suas localizações da base de dados
	rows, err := db.DB.Query(`
		SELECT f.fragment_order, n.id, n.node_address
		FROM FileFragments AS f 
		JOIN FragmentLocation AS fl ON f.fragment_id = fl.fragment_id
		JOIN Nodes AS n ON fl.node_id = n.id 
		WHERE f.file_id = $1`, fileID)
	if err != nil {
		return fmt.Errorf("Erro ao buscar fragmentos na base de dados: %v", err)
	}
	defer rows.Close()

	// Iterar pelos fragmentos e remover de todos os nós
	for rows.Next() {
		var fragmentOrder int
		var nodeID int
		var nodeAddress string

		if err := rows.Scan(&fragmentOrder, &nodeID, &nodeAddress); err != nil {
			return fmt.Errorf("Erro ao escanear fragmento: %v", err)
		}

		// Apagar o fragmento do nó
		err := DeleteFragmentOnNode(fileID, fragmentOrder, nodeID)
		if err != nil {
			fmt.Printf("Erro ao remover fragmento %d do nó %d: %v\n", fragmentOrder, nodeID, err)
			continue
		}

		// Remover a informação de distribuição da base de dados
		err = DeleteDistributionInfo(fileID, fragmentOrder, nodeID)
		if err != nil {
			fmt.Printf("Erro ao remover informação de distribuição para o fragmento %d no nó %s: %v\n", fragmentOrder, nodeAddress, err)
			continue
		}
	}

	// Remover as informações do arquivo da base de dados
	err = DeleteFileInfo(fileID)
	if err != nil {
		return fmt.Errorf("Erro ao remover informações do arquivo na base de dados: %v", err)
	}

	return nil
}

func DeleteFragmentOnNode(fileID, fragmentOrder, nodeID int) error {
    // Construir a URL para o node que lida com a exclusão do fragmento
    deleteURL := fmt.Sprintf("http://node%d:8082/fragments/delete?file_id=%d&fragment_order=%d", nodeID, fileID, fragmentOrder)

    req, err := http.NewRequest("DELETE", deleteURL, nil)
    if err != nil {
        return fmt.Errorf("erro ao criar requisição para deletar fragmento no nó %d: %v", nodeID, err)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("erro ao conectar ao nó %d para deletar fragmento: %v", nodeID, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("falha ao deletar fragmento no nó %d: %s", nodeID, string(body))
    }

    fmt.Printf("Fragmento %d do arquivo %d deletado do nó %d com sucesso\n", fragmentOrder, fileID, nodeID)
    return nil
}

// DeleteDistributionInfo remove informações de fragmentos da tabela de distribuição na base de dados.
func DeleteDistributionInfo(fileID, fragmentOrder, nodeID int) error {
	query := "DELETE FROM FragmentLocation WHERE fragment_id = (SELECT fragment_id FROM FileFragments WHERE file_id = $1 AND fragment_order = $2) AND node_id = $3"
	_, err := db.DB.Exec(query, fileID, fragmentOrder, nodeID)
	if err != nil {
		return fmt.Errorf("Erro ao remover informações de distribuição do fragmento: %v", err)
	}
	return nil
}

// DeleteFileInfo remove informações do arquivo e todos os fragmentos relacionados.
func DeleteFileInfo(fileID int) error {
	// Remover todas as informações de fragmentos do arquivo
	_, err := db.DB.Exec("DELETE FROM FileFragments WHERE file_id = $1", fileID)
	if err != nil {
		return fmt.Errorf("Erro ao remover fragmentos do arquivo: %v", err)
	}

	// Remover o próprio arquivo
	_, err = db.DB.Exec("DELETE FROM Files WHERE id = $1", fileID)
	if err != nil {
		return fmt.Errorf("Erro ao remover arquivo: %v", err)
	}
	return nil
}