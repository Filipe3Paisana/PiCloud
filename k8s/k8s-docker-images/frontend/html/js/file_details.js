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
    // Calcular o número de nodes únicos
    const uniqueNodes = new Set();
    fileData.fragments.forEach(fragment => {
        fragment.nodes.forEach(node => {
            uniqueNodes.add(node.location); // Use "location" como identificador único
        });
    });
    document.getElementById('totalNodes').textContent = uniqueNodes.size;


    // Ordenar os fragmentos por ordem crescente de `fragment.order`
    const sortedFragments = fileData.fragments.sort((a, b) => a.order - b.order);

    // Atualizar a lista de fragmentos como cards
    const fragmentsList = document.getElementById('fragmentsList');
    fragmentsList.innerHTML = ''; // Limpar a lista antes de adicionar novos itens

    sortedFragments.forEach(fragment => {
        const fragmentCard = document.createElement('div');
        fragmentCard.classList.add('file-card');

        const fragmentInfo = document.createElement('div');
        fragmentInfo.classList.add('clickable-area');

        // Informações do fragmento
        const fragmentTitle = document.createElement('h4');
        fragmentTitle.textContent = `Fragmento ${fragment.order}`;
        
        const fragmentHash = document.createElement('p');
        fragmentHash.innerHTML = `<strong>Hash:</strong> ${fragment.hash}`;

        const fragmentLocation = document.createElement('p');
        const location = fragment.nodes.map(node => node.location).join(', ');
        fragmentLocation.innerHTML = `<strong>Nodes:</strong> ${location}`;

        // Adicionar informações ao card
        fragmentInfo.appendChild(fragmentTitle);
        fragmentInfo.appendChild(fragmentHash);
        fragmentInfo.appendChild(fragmentLocation);
        fragmentCard.appendChild(fragmentInfo);

        fragmentsList.appendChild(fragmentCard);
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
    fetch(`/api/user/file/details?file_id=${fileId}`, {
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
