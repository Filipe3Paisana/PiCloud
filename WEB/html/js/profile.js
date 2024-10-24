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
        alert('Você precisa estar logado para acessar esta página.');
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
        alert('Por favor, selecione um arquivo para upload.');
        return;
    }
    const maxSize = 100 * 1024 * 1024; // 5 MB
    if (file.size > maxSize) {
        alert('O arquivo excede o tamanho máximo permitido de 5MB.');
        return;
    }

    const formData = new FormData();
    formData.append('file', file);

    const uploadMessage = document.createElement('div');
    uploadMessage.textContent = 'Loading...';
    document.body.appendChild(uploadMessage);
    console.log("FormData created, starting fetch...");

    fetch('http://localhost:8081/user/upload', {
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
    })
    .catch(error => {
        console.error('Erro:', error);
        alert('Erro ao enviar o arquivo. Por favor, tente novamente.');
    })
    .finally(() => {
        console.log("Removing upload message");
        uploadMessage.remove();
    });
}
function fetchUserFiles() {
    fetch('http://localhost:8081/user/files', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('authToken')}`
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Erro ao obter a lista de arquivos: ' + response.statusText);
        }
        return response.json();
    })
    .then(data => {
        console.log("Arquivos recebidos:", data);
        if (!data || !Array.isArray(data)) {
            console.error("Dados inválidos recebidos:", data);
            alert('Erro ao obter a lista de arquivos. Por favor, tente novamente.');
            return;
        }
        displayFiles(data);
    })
    .catch(error => {
        console.error('Erro:', error);
        alert('Erro ao obter a lista de arquivos. Por favor, tente novamente.');
    });
}

function displayFiles(files) {
    const filesList = document.getElementById('filesList');
    filesList.innerHTML = ''; // Limpa o conteúdo existente

    if (!files || files.length === 0) {
        filesList.textContent = 'Você não possui arquivos.';
        return;
    }

    files.forEach(file => {
        const listItem = document.createElement('li');
        listItem.textContent = `${file.name} (${formatFileSize(file.size)})`;

        const downloadButton = document.createElement('button');
        downloadButton.textContent = 'Download';
        downloadButton.onclick = () => downloadFile(file.id); 

        listItem.appendChild(downloadButton);

        filesList.appendChild(listItem);
    });
}

function formatFileSize(bytes) {
    if (bytes < 1024) return `${bytes} B`;
    else if (bytes < 1048576) return `${(bytes / 1024).toFixed(2)} KB`;
    else if (bytes < 1073741824) return `${(bytes / 1048576).toFixed(2)} MB`;
    else return `${(bytes / 1073741824).toFixed(2)} GB`;
}


function formatFileSize(bytes) {
    if (bytes < 1024) return `${bytes} B`;
    else if (bytes < 1048576) return `${(bytes / 1024).toFixed(2)} KB`;
    else if (bytes < 1073741824) return `${(bytes / 1048576).toFixed(2)} MB`;
    else return `${(bytes / 1073741824).toFixed(2)} GB`;
}

function downloadFile(fileID) {
    const url = `http://localhost:8081/user/download?file_id=${fileID}`;
    
    fetch(url, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${localStorage.getItem('authToken')}`
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Erro ao baixar o arquivo: ' + response.statusText);
        }
        return response.blob();
    })
    .then(blob => {
        // Criar URL para o blob e baixar o arquivo
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        
        // Definir o nome do arquivo para download (você pode melhorar para usar o nome real do arquivo)
        a.download = `arquivo_${fileID}`; 
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.URL.revokeObjectURL(url);
    })
    .catch(error => {
        console.error('Erro ao baixar o arquivo:', error);
        alert('Erro ao baixar o arquivo. Por favor, tente novamente.');
    });
}

function logout() {
    localStorage.removeItem('authToken'); 
    window.location.href = 'index.html'; 
}

window.onload = function() {
    const token = localStorage.getItem('authToken');

    if (!token) {
        alert('Você precisa estar logado para acessar esta página.');
        window.location.href = 'index.html';
        return;
    }
    greetUser(token);

    fetchUserFiles();
};


function logout() {
    localStorage.removeItem('authToken'); 
    window.location.href = 'index.html'; 
}
