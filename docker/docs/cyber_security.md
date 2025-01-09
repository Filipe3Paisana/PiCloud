# Tópicos para Segurança

### docker scout
 - para ver as vulnerabilidades das imagens e retirar o que não está a ser utilizado, criando uma imagem adaptada para o nosso projeto e mais segura. (limpeza de packages, libraries desnecessárias) 
 - container as non-root, ou seja, o container não tem todas as permissões, não sendo possível comprometer o sistema

### Kubernetes
 - implementar encriptação entre pods
 - restringir as permissões dos pods
 - limitar comunicação entre os pods (Que pods podem receber informação e com quem é que estes podem falar) NetworkPolicy garantindo que o frontend apenas fala com o backend 
 - Utilizar ServiceMesh para configurar a nivel aplicacional, basicamente proxy que intrepreta a intrada e saida de trafego, controlo e monitorização de comunicações, implementações de timeouts ect ...
 - Alem disso utilizar o ServiceMesh para Ativar a mTLS entre pods para garantir encriptação entre pods
 - Utilização e secret dever ser acegurada pelo outem AWS ou apens encriptar Secrets não sei bem como seria isso
 - Utilizar Argo CD para garantir consistencia nos clusters


### fragmentos
 - fazendo um hash dos fragmentos, conseguimos garantir a integridade dos mesmos para reconstruir o ficheiro corretamente. 
 - encriptar informação dos fragmentos, para que não sejam acedidos pelo host do node. TODO


### CloudFlare
 - para prevenir ataques de ddos 
 - O tráfego é criptografado e inspecionado pela Cloudflare TLS.
 - Adicionar rate Limiting(numero de requisições por minuto a um IP)
 



### Ecriptação da Base de Dados
 - Passar as passwords para hash (atualmente com um custo de 10, vaz sentido aumentar se sim para quanto?)

    Cost	Tempo por Hash (ms)	Total para Força Bruta (anos)
    4	1 ms	~6.91 horas
    6	10 ms	~28.78 dias
    8	100 ms	~7.88 anos
    10	300 ms	~23.65 anos
    12	1.2 segundos	~94.6 anos
    14	4.8 segundos	~378.6 anos
    16	19.2 segundos	~1,514.5 anos
    Foi considerado uma password de 8 caracteres alfanuméricos


   #### Impletações 

   Objetivo:
   A função validateUserInput valida os dados fornecidos pelos usuários (nome de usuário, e-mail e senha) para garantir que estejam no formato adequado, protegendo contra vulnerabilidades como XSS e entradas inválidas.

   Validações Implementadas:

   Nome de Usuário:
   Permite apenas letras, números e sublinhados, com até 50 caracteres.
   E-mail:
   Verifica se o e-mail tem o formato correto (ex.: usuario@dominio.com).
   Senha:
   Exige pelo menos 8 caracteres e pelo menos 1 número.
   Benefícios de Segurança:

   Prevenção de XSS e SQL Injection.
   Garante que os dados não contenham entradas maliciosas.
   Impede senhas fracas e entradas inválidas.
   Código:

   func validateUserInput(user *models.User) error {
      // Validação do nome de usuário
      usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{1,50}$`)
      if !usernameRegex.MatchString(user.Username) {
         return errors.New("Username inválido: deve conter apenas letras, números ou sublinhados e ter no máximo 50 caracteres")
      }

      // Validação de e-mail
      emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
      if !emailRegex.MatchString(user.Email) {
         return errors.New("E-mail inválido: formato de e-mail incorreto")
      }

      // Validação de senha
      if len(user.Password) < 8 || !strings.ContainsAny(user.Password, "0123456789") {
         return errors.New("Senha inválida: deve conter pelo menos 8 caracteres e incluir pelo menos um número")
      }

      return nil
   }


   