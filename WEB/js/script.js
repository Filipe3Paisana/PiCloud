document.getElementById('loginForm').addEventListener('submit', function(e) {
    e.preventDefault(); 

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    
    if (email === "" || password === "") {
        alert("Por favor, preencha todos os campos!");
        return;
    }
    
    alert(`Login efetuado com sucesso! Bem-vindo(a), ${email}`);
});
