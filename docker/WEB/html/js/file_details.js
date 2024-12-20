// Função para formatar o tamanho do ficheiro
function formatFileSize(bytes) {
    if (bytes < 1024) return `${bytes} B`;
    else if (bytes < 1048576) return `${(bytes / 1024).toFixed(2)} KB`;
    else if (bytes < 1073741824) return `${(bytes / 1048576).toFixed(2)} MB`;
    return `${(bytes / 1073741824).toFixed(2)} GB`;
}

// Função para exibir os detalhes do ficheiro
function displayFileDetails(fileData) {
    // Atualizar os detalhes básicos do ficheiro
    document.getElementById('fileName').textContent = fileData.file.name;
    document.getElementById('fileSize').textContent = formatFileSize(fileData.file.size);
    document.getElementById('totalFragments').textContent = fileData.fragments.length;

    // Atualizar a lista de fragmentos e as localizações nos nodes
    const fragmentsList = document.getElementById('fragmentsList');
    fragmentsList.innerHTML = ''; // Limpar a lista antes de adicionar novos itens

    fileData.fragments.forEach(fragment => {
        const listItem = document.createElement('li');
        const nodeAddresses = fragment.nodes.map(node => node.node_address).join(', ');
        listItem.textContent = `Fragmento ${fragment.order} (Hash: ${fragment.hash}) - Nodes: ${nodeAddresses}`;
        fragmentsList.appendChild(listItem);
    });
}

// Função para carregar os detalhes do ficheiro
function fetchFileDetails() {
    // Obter o ID do ficheiro a partir dos parâmetros da URL
    const queryParams = new URLSearchParams(window.location.search);
    const fileId = queryParams.get('file_id');

    if (!fileId) {
        alert('Ficheiro inválido.');
        window.location.href = 'profile.html';
        return;
    }

    // Fazer a requisição para obter os detalhes do ficheiro
    fetch(`http://localhost:8081/user/file/details?file_id=${fileId}`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('authToken')}`
        }
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Erro ao obter detalhes do ficheiro: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            // Exibir os detalhes do ficheiro
            displayFileDetails(data);
        })
        .catch(error => {
            console.error('Erro:', error);
            alert('Erro ao carregar os detalhes do ficheiro. Por favor, tente novamente.');
            window.location.href = 'profile.html';
        });
}

// Função para realizar logout
function backToProfile() {
    
    window.location.href = 'profile.html';
}

// Adicionar evento para carregar os detalhes do ficheiro quando a página estiver pronta
document.addEventListener('DOMContentLoaded', fetchFileDetails);
