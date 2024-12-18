# Usar uma imagem base do Python
FROM python:3.9-slim

# Definir diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiar todos os arquivos para o diretório de trabalho dentro do contêiner
COPY . /app

# Instalar as dependências listadas no arquivo requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Definir o comando que será executado ao iniciar o contêiner (executar o modelo AI)
CMD ["python", "PiCloud.py"]
