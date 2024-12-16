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
        "password": password 
    };

    try {
        
        const response = await fetch('api/users/add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        });

    
        if (response.ok) {
            const result = await response.json();
            alert('Registado com sucesso! Bem-vindo, ' + result.username);
            window.location.href = 'index.html'; // Redirecionar para a página de login após o registo
        } else {
            const errorText = await response.text();
            alert('Erro ao registar: ' + errorText);
        }
    } catch (error) {
        console.error('Erro:', error);
        alert('Erro ao registar. Tente novamente mais tarde.');
    }
});
