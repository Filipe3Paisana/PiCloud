package handlers

import (
	"fmt"
    "encoding/json"
    "net/http"
    "strconv"

	"api/utils"
	"api/models"
	"api/db"
)

// GetFileDetailsHandler retorna os detalhes de um ficheiro, incluindo fragmentos e os nodes onde estão armazenados
func GetFileDetailsHandler(w http.ResponseWriter, r *http.Request) {
    utils.EnableCors(w, r) // Habilitar CORS para a requisição

    // Validar o método HTTP
    if r.Method != http.MethodGet {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    // Extraindo o ID do ficheiro dos parâmetros
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

    // Buscar informações básicas do ficheiro
    var file models.File
    err = db.DB.QueryRow(`SELECT id, name, size FROM Files WHERE id = $1`, fileID).Scan(&file.ID, &file.Name, &file.Size)
    if err != nil {
        http.Error(w, fmt.Sprintf("Ficheiro não encontrado: %v", err), http.StatusNotFound)
        return
    }

    // Buscar fragmentos e nodes associados
    rows, err := db.DB.Query(`
        SELECT ff.fragment_id, ff.fragment_hashes, ff.fragment_order, n.location
        FROM FileFragments ff
        JOIN FragmentLocation fl ON ff.fragment_id = fl.fragment_id
        JOIN Nodes n ON fl.node_id = n.id
        WHERE ff.file_id = $1`, fileID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Erro ao buscar fragmentos: %v", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    type FileFragment struct {
        FragmentID int          `json:"fragment_id"`
        Hash       string       `json:"hash"`
        Order      int          `json:"order"`
        Nodes      []models.Node `json:"nodes"`
    }

    fragmentsMap := make(map[int]*FileFragment)
    for rows.Next() {
        var fragmentID int
        var fragmentHash string
        var fragmentOrder int
        var Location string

        if err := rows.Scan(&fragmentID, &fragmentHash, &fragmentOrder, &Location); err != nil {
            http.Error(w, fmt.Sprintf("Erro ao processar fragmentos: %v", err), http.StatusInternalServerError)
            return
        }

        if _, exists := fragmentsMap[fragmentID]; !exists {
            fragmentsMap[fragmentID] = &FileFragment{
                FragmentID: fragmentID,
                Hash:       fragmentHash,
                Order:      fragmentOrder,
                Nodes:      []models.Node{},
            }
        }

        fragmentsMap[fragmentID].Nodes = append(fragmentsMap[fragmentID].Nodes, models.Node{
            Location: Location,
        })
    }

    var fragments []FileFragment
    for _, fragment := range fragmentsMap {
        fragments = append(fragments, *fragment)
    }

    // Construir a resposta JSON
    response := struct {
        File      models.File     `json:"file"`
        Fragments []FileFragment `json:"fragments"`
    }{
        File:      file,
        Fragments: fragments,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}