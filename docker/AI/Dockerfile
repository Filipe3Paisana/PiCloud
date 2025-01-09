# Usar uma imagem base do Python
FROM python:3.9-slim

# Definir diretório de trabalho dentro do contêiner
WORKDIR /app

# Instalar dependências do sistema (para LightGBM e outros)
RUN apt-get update && apt-get install -y \
    libgomp1 \
    && rm -rf /var/lib/apt/lists/*

# Copiar todos os arquivos para o diretório de trabalho dentro do contêiner
COPY . /app

# Instalar as dependências listadas no arquivo requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Definir o comando que será executado ao iniciar o contêiner (executar o modelo AI)
CMD ["python", "PiCloud.py"]