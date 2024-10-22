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
    
    const maxSize = 100 * 1024 * 1024; // 100 MB
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
            Authorization: `Bearer ${localStorage.getItem('authToken')}` // Apenas o token
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

function logout() {
    localStorage.removeItem('authToken'); 
    window.location.href = 'index.html'; 
}
