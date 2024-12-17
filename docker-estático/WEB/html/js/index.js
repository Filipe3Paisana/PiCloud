document.getElementById('loginForm').addEventListener('submit', async function(e) {
    e.preventDefault(); 

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    
    if (email === "" || password === "") {
        alert("Por favor, preencha todos os campos!");
        return;
    }

    try {
        const response = await fetch('http://localhost:8081/users/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email: email, password: password })
        });

        console.log('Status da resposta:', response.status);
        console.log('Resposta completa:', response);

        if (response.ok) {
            const result = await response.json();
            console.log('Dados da resposta:', result);

            // Armazenar o token no LocalStorage
            localStorage.setItem('authToken', result.token); // Armazena o token

            alert(`Login efetuado com sucesso! Bem-vindo(a), ${email}`);
            window.location.href = 'profile.html'; // Redireciona para o perfil
        } else {
            const errorData = await response.text();
            console.error('Erro ao fazer login:', errorData);
            alert(`Erro ao fazer login: ${errorData}`);
        }
    } catch (error) {
        console.error('Erro de rede:', error);
        alert('Erro ao tentar fazer login. Tente novamente mais tarde.');
    }
});
