document.getElementById('signupForm').addEventListener('submit', async function(event) {
    event.preventDefault(); // Prevenir o comportamento padrão do formulário

    // Capturar os dados do formulário
    const name = document.getElementById('name').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    // Validar se as senhas correspondem
    if (password !== confirmPassword) {
        alert('As senhas não correspondem.');
        return;
    }

    // Montar o objeto de dados para enviar
    const userData = {
        "username": name,
        "email": email,
        "password": password // Aqui você envia a senha, mas seria melhor aplicar hashing no lado do servidor
    };

    try {
        // Fazer a requisição POST para a API de backend
        const response = await fetch('http://localhost:8081/users/add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        });

        // Verificar a resposta
        if (response.ok) {
            const result = await response.json();
            alert('Usuário registrado com sucesso! Bem-vindo, ' + result.username);
            window.location.href = 'index.html'; // Redirecionar para a página de login após o registro
        } else {
            const errorText = await response.text();
            alert('Erro ao registrar: ' + errorText);
        }
    } catch (error) {
        console.error('Erro:', error);
        alert('Erro ao registrar. Tente novamente mais tarde.');
    }
});
