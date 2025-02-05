# Resumo da Arquitetura para Cloud Privada com Fragmentação e Recuperação de Ficheiros

## Componentes Principais da Arquitetura

1. **Dois Datacenters (Ativo-Passivo) com Kubernetes**:
   - Cada **datacenter** será um cluster Kubernetes localizado em locais diferentes, sendo um **ativo** e outro **passivo**, configurados como **Ativo-Passivo** para alta disponibilidade.
   - **Cloudflare** gere o balanceamento de carga e o failover entre os dois datacenters. Enquanto o datacenter ativo lida com todas as requisições, o passivo é ativado apenas em caso de falha do ativo.
   - Dentro de cada datacenter, haverá **dois pods** em um esquema **Ativo-Ativo**, onde cada pod conterá um **serviço WEB**, uma **API** e uma **base de dados** local.
   - **Kubernetes** gere a comunicação interna entre os pods dentro de cada datacenter, utilizando a rede interna para equilíbrio de carga e descoberta de serviços.

2. **Fragmentação de Ficheiros (APIs)**:
   - As APIs em cada datacenter são responsáveis por **fragmentar os ficheiros** grandes em partes menores. Essas partes (fragmentos) são identificadas com **metadados**, que incluem informações sobre a localização, ordem e réplicas.
   - O processo de **fragmentação** utiliza técnicas como **Erasure Coding** ou **Sharding**, garantindo que, mesmo com a perda de alguns fragmentos, o ficheiro possa ser reconstruído.

3. **Distribuição de Fragmentos em Nodes**:
   - Os fragmentos são distribuídos para **nodes de armazenamento** que podem estar **dentro ou fora dos datacenters**.
   - A API escolhe os nodes de acordo com a disponibilidade e a proximidade geográfica, garantindo uma distribuição eficiente e equilibrada.
   - Cada fragmento pode ter **várias réplicas** distribuídas entre diferentes nodes para garantir a resiliência. Isto significa que, se um node falhar, os fragmentos ainda estarão acessíveis em outros nodes.

4. **Replicação e Monitorização de Fragmentos**:
   - **Replicação**: A API implementa a replicação dos fragmentos em múltiplos nodes para garantir que, mesmo em caso de falha de um ou mais nodes, os fragmentos possam ser recuperados de suas réplicas.
   - **Monitorização**: A API monitora a saúde dos nodes e a disponibilidade dos fragmentos. Se um nodo falhar, novas réplicas podem ser geradas automaticamente e redistribuídas para outros nodes disponíveis.
   - A replicação de fragmentos não requer comunicação contínua entre os datacenters, tornando a sincronização eficiente e assíncrona.

5. **Reconstrução de Ficheiros**:
   - Quando solicitado, a API consulta base de dados com os **metadados** para obter a localização dos fragmentos e, então, recupera os fragmentos dos nodes necessários.
   - A recuperação dos fragmentos é feita de forma **paralela** para otimizar a velocidade da reconstrução.
   - A API verifica a **integridade** dos fragmentos usando a sua **hash**, e, uma vez recuperados, reconstrói o ficheiro original.
   - Se algum fragmento estiver inacessível, a API tenta aceder a uma réplica ou emite um erro caso muitos fragmentos estejam ausentes.

6. **Nodes Distribuídos (Dentro e Fora da Rede dos Datacenters)**:
   - Os nodes que armazenam os fragmentos podem estar localizados tanto **dentro** quanto **fora** dos datacenters. Cada nodo é identificado por um **ID único** e oferece uma interface para a API solicitar fragmentos.
   - Esses nodes podem estar em diferentes **redes** ou localizações geográficas, permitindo maior flexibilidade na distribuição dos fragmentos.
   - A comunicação entre os nodes e as APIs dos datacenters será feita de forma **segura** (via HTTPS ou outro protocolo de criptografia) para proteger os dados em trânsito.

7. **Segurança e Comunicação**:
   - **Criptografia de Dados**: Os fragmentos devem ser armazenados de forma criptografada nos nodes, garantindo que, mesmo que um nodo seja comprometido, os dados permaneçam protegidos.
   - **Autenticação dos Nodes**: Apenas nodes autenticados e confiáveis podem participar do sistema de armazenamento e recuperação.
   - **Comunicação Segura**: Todo o tráfego entre os datacenters, nodes e APIs será protegido por SSL/TLS, e a API pode usar **Cloudflare Tunnel** ou outro meio seguro para garantir que a comunicação entre as partes distribuídas seja segura e eficiente.

8. **Failover e Redundância**:
   - Em caso de falha do datacenter **ativo**, o Cloudflare automaticamente redireciona o tráfego para o datacenter **passivo**, que assume o papel de ativo, mantendo a continuidade do serviço.
   - Como a base de dados dos metadados é replicada entre os datacenters, o sistema pode rapidamente aceder às informações necessárias para reconstruir ficheiros e gerir fragmentos.

## Fluxo de Funcionamento

1. **Fragmentação**: O ficheiro é carregado na API, que fragmenta o ficheiro em múltiplas partes e gera metadados para cada fragmento.
2. **Distribuição**: Os fragmentos são distribuídos para nodes disponíveis, e réplicas são criadas para garantir redundância.
3. **Armazenamento**: Os fragmentos são armazenados de forma criptografada e distribuídos de maneira eficiente pelos nodes.
4. **Monitorização e Replicação**: A API monitora a saúde dos nodes e a integridade dos fragmentos, gerando novas réplicas quando necessário.
5. **Reconstrução**: Quando um ficheiro é solicitado, a API recupera os fragmentos necessários dos nodes, verifica a sua integridade e reconstrói o ficheiro original.

## Resumo Final
- **Dois datacenters**, geridos via **Cloudflare** em um esquema **Ativo-Passivo**, garantem alta disponibilidade e failover.
- Dentro de cada datacenter, o sistema é **Ativo-Ativo**, com dois pods executando os serviços de **WEB**, **API** e **BD**, e utilizando o **Kubernetes** para gestão interna.
- Os ficheiros são fragmentados pela API, e os fragmentos são distribuídos e replicados em **nodes distribuídos**, tanto dentro quanto fora dos datacenters, garantindo **resiliência** e **segurança**.
- A comunicação entre nodes e APIs é feita de forma **segura**, com a possibilidade de usar **Cloudflare Tunnel** para criar uma rede confiável.
- O sistema é capaz de **reconstruir ficheiros** ao solicitar os fragmentos de diferentes nodes, garantindo a recuperação eficiente e a proteção dos dados.

Esta arquitetura distribui a carga e aumenta a **resiliência** do sistema, com capacidade de lidar com falhas de nodes ou datacenters, mantendo sempre a capacidade de recuperação dos ficheiros fragmentados.
