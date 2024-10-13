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
    const userEmail = userData.email; 

    
    document.getElementById('username').textContent = userName; 

    if (userName) {
        alert(`Olá, ${userName}! Seu email é ${userEmail}.`);
    } else {
        alert('Olá, usuário!');
    }
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



function logout() {
    localStorage.removeItem('authToken'); 
    window.location.href = 'index.html'; 
}
