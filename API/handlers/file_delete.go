package handlers

// import (
// 	"api/models"
// )


// func DeleteFragments(fileID int, fragments [][]byte, replicationFactor int) error {
//     availableNodes := GetOnlineNodesList()
//     if len(availableNodes) == 0 {
//         return fmt.Errorf("Nenhum nó disponível para distribuir os fragmentos")
//     }

//     for i := 1; i <= numberOfFragments; i++ {
//         selectedNodes := SelectNodesWithFragment(availableNodes, fileID)

//         for _, node := range selectedNodes {
//             err := DeleteFragments(fileID, i, fragments[i-1], node.NodeID) // Passando o conteúdo real do fragmento
//             if err != nil {
//                 fmt.Printf("Erro ao eliminar fragmento %d do nó %d: %v\n", i, node.NodeID, err)
//                 continue
//             }

//             err = DeleteDistributionInfo(fileID, i, node.NodeAddress)
//             if err != nil {
//                 fmt.Printf("Erro ao salvar informações de distribuição para o fragmento %d no nó %s: %v\n", i, node.NodeAddress, err)
//                 continue
//             }
//         }
//     }
//     return nil
// }


// SelectNodesWithFragment(availableNodes models.Node, fileID int) {
// 	return nil
// }

// DeleteFragments(fileID int, i int, fragments int, nodeID int) {
// 	return nil
// }

// DeleteDistributionInfo(fileID int, i int, nodeAddress string) {
// 	return nil
// }