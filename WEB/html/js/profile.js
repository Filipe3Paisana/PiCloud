// profile.js

// Função para decodificar um token JWT
function parseJwt(token) {
    const base64Url = token.split('.')[1]; // Pega a parte do payload
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/'); // Corrige a formatação
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload); // Converte de string para objeto
}

// Verifica se o token está presente no LocalStorage
window.onload = function() {
    const token = localStorage.getItem('authToken');

    if (!token) {
        alert('Você precisa estar logado para acessar esta página.');
        window.location.href = 'index.html'; 
        return;
    }

    
    const userData = parseJwt(token);
    const username = userData.sub; // Ou qualquer que seja a propriedade que você armazena o nome do usuário

    // Atualiza o título da página com o nome do usuário
    document.title = `PiCloud - ${username}`;

    
    
};



function logout() {
    localStorage.removeItem('authToken'); // Remove o token no logout
    window.location.href = 'login.html'; // Redireciona para a página de login
}
