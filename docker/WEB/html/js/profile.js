function parseJwt(token) {
    const base64Url = token.split('.')[1]; // Pega a parte do payload
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/'); // Corrige a formatação
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload); // Converte de string para objeto
}

function greetUser(token) {
    const userData = parseJwt(token);
    const userName = userData.username; 
    //const userEmail = userData.email; 

    document.getElementById('username').textContent = userName; 
}

window.onload = function() {
    const token = localStorage.getItem('authToken');

    if (!token) {
        alert('Precisa estar autenticado para aceder a esta página.');
        window.location.href = 'index.html'; 
        return;
    }
    greetUser(token);

    
};

function uploadFile() {
    const fileInput = document.getElementById('fileInput');
    const file = fileInput.files[0];
    console.log("Selected file:", file);

    if (!file) {
        alert('Por favor, selecione um ficheiro para upload.');
        return;
    }

    const maxSize = 1000 * 1024 * 1024; // 100 MB
    if (file.size > maxSize) {
        alert('O ficheiro excede o tamanho máximo permitido de 100MB.');
        return;
    }

    const formData = new FormData();
    formData.append('file', file);

    const uploadMessage = document.createElement('div');
    uploadMessage.textContent = 'Loading...';
    document.body.appendChild(uploadMessage);
    console.log("FormData created, starting fetch...");

    fetch('http://localhost:8080/user/upload', {
        method: 'POST',
        headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}` 
        },
        body: formData
    })
    .then(response => {
        console.log("Fetch response:", response);
        if (!response.ok) {
            throw new Error('Erro ao fazer upload: ' + response.statusText);
        }
        return response.json();
    })
    .then(data => {
        alert(data.message);
        console.log("Upload success message:", data.message);
        fetchUserFiles()
    })
    .catch(error => {
        console.error('Erro:', error);
        alert('Erro ao enviar o ficheiro. Por favor, tente novamente.');
    })
    .finally(() => {
        console.log("Removing upload message");
        uploadMessage.remove();
    });
}
function fetchUserFiles() {
    fetch('http://localhost:8080/user/files', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('authToken')}`
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Erro ao obter a lista de ficheiros: ' + response.statusText);
        }
        return response.json();
    })
    .then(data => {
        console.log("ficheiros recebidos:", data);
        if (!data || !Array.isArray(data)) {
            console.error("Dados inválidos recebidos:", data);
            alert('Erro ao obter a lista de ficheiros. Por favor, tente novamente.');
            return;
        }
        displayFiles(data);
    })
    .catch(error => {
        console.error('Erro:', error);
        alert('Erro ao obter a lista de ficheiros. Por favor, tente novamente.');
    });
}

function displayFiles(files) {
    const filesList = document.getElementById('filesList');
    filesList.innerHTML = ''; // Limpa o conteúdo existente

    if (!files || files.length === 0) {
        filesList.innerHTML = '<p>Você não possui ficheiros.</p>';
        return;
    }

    files.forEach(file => {
        const fileCard = document.createElement('div');
        fileCard.classList.add('file-card');

        const fileName = document.createElement('h4');
        fileName.textContent = `${file.name} (${formatFileSize(file.size)})`;

        const iconButtons = document.createElement('div');
        iconButtons.classList.add('icon-buttons');

        const downloadButton = document.createElement('button');
        downloadButton.innerHTML = '<i class="fas fa-download"></i>';
        downloadButton.onclick = () => downloadFile(file.id);

        const deleteButton = document.createElement('button');
        deleteButton.innerHTML = '<i class="fas fa-trash-alt"></i>';
        deleteButton.onclick = () => deleteFile(file.id);

        iconButtons.appendChild(downloadButton);
        iconButtons.appendChild(deleteButton);

        fileCard.appendChild(fileName);
        fileCard.appendChild(iconButtons);

        filesList.appendChild(fileCard);
    });
}

function formatFileSize(bytes) {
    if (bytes < 1024) return `${bytes} B`;
    else if (bytes < 1048576) return `${(bytes / 1024).toFixed(2)} KB`;
    else if (bytes < 1073741824) return `${(bytes / 1048576).toFixed(2)} MB`;
    else return `${(bytes / 1073741824).toFixed(2)} GB`;
}

function downloadFile(fileID) {
    const url = `http://localhost:8080/user/download?file_id=${fileID}`;
    
    fetch(url, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('authToken')}`
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Erro ao descarregar o ficheiro: ' + response.statusText);
        }
        return response.blob();
    })
    .then(blob => {
        // Criar URL para o blob e baixar o ficheiro
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        
        // Definir o nome do ficheiro para download (você pode melhorar para usar o nome real do ficheiro)
        a.download = `ficheiro_${fileID}`; 
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.URL.revokeObjectURL(url);
    })
    .catch(error => {
        console.error('Erro ao descarregar o ficheiro:', error);
        alert('Erro ao descarregar o ficheiro. Por favor, tente novamente.');
    });
}

function deleteFile(fileID) {
    const confirmDelete = confirm("Tem certeza que deseja eliminar este ficheiro?");
    if (!confirmDelete) return;

    const url = `http://localhost:8080/user/delete?file_id=${fileID}`;

    fetch(url, {
        method: 'DELETE',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('authToken')}`
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Erro ao eliminar o ficheiro: ' + response.statusText);
        }
        return response.json();
    })
    .then(data => {
        alert(data.message);
        console.log("ficheiro eliminado com sucesso:", data.message);
        // Atualizar a lista de ficheiros após eliminar
        fetchUserFiles();
    })
    .catch(error => {
        console.error('Erro ao eliminar o ficheiro:', error);
        alert('Erro ao eliminar o ficheiro. Por favor, tente novamente.');
    });
}

function logout() {
    localStorage.removeItem('authToken'); 
    window.location.href = 'index.html'; 
}

window.onload = function() {
    const token = localStorage.getItem('authToken');

    if (!token) {
        alert('Precisa de estar autenticado para aceder a esta página.');
        window.location.href = 'index.html';
        return;
    }
    greetUser(token);

    fetchUserFiles();
};