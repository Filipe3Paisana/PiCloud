import json
import os
import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_squared_error, r2_score
from matplotlib import pyplot as plt
import matplotlib.pyplot as plt

from lazypredict.Supervised import LazyRegressor
from sklearn.model_selection import train_test_split
# from lazypredict.Supervised import LazyClassifier
# from sklearn.model_selection import train_test_split


# Caminho para o arquivo JSON (ajuste de acordo com o local do seu arquivo dentro do contêiner)
file_path = '/app/prometheus_data/all_metrics_sim.json'

# Verificar se o arquivo existe
if not os.path.exists(file_path):
    raise FileNotFoundError(f"O arquivo {file_path} não foi encontrado!")

# Carregar dados do arquivo JSON
with open(file_path, 'r') as f:
    data = json.load(f)

print("Current Working Directory:", os.getcwd())

# Estruturar os dados em formato adequado para o modelo
metrics_data = []
for entry in data:
    node_id = entry.get('node_id', 'unknown')
    timestamp = entry.get('timestamp', None)
    
    # Usar .get() para evitar KeyError se a chave não existir
    cpu_usage = entry.get('metrics', {}).get('CPU', [{}])[0].get('value', np.nan)  # Valor padrão se não existir
    disk_usage = np.mean([disk.get('value', np.nan) for disk in entry.get('metrics', {}).get('Disk', [])] or [0])  # Evitar lista vazia
    memory_available = entry.get('metrics', {}).get('Memory', [{}])[0].get('value', np.nan)
    response_time = entry.get('metrics', {}).get('ResponseTime', [{}])[0].get('value', np.nan)

    hour_of_day = entry.get('metrics', {}).get('HourOfDay', [{}])[0].get('value', np.nan)
    minute = entry.get('metrics', {}).get('Minute', [{}])[0].get('value', np.nan)
    
    
    # Verificar se temos valores válidos para as métricas (evitar NaN)
    if np.isnan(cpu_usage) or np.isnan(disk_usage) or np.isnan(memory_available):
        continue  # Ignorar entradas incompletas
    
    # Adicionar os dados de cada node
    metrics_data.append({
        'node_id': node_id,
        'timestamp': timestamp,
        'cpu_usage': cpu_usage,
        'disk_usage': disk_usage,
        'memory_available': memory_available,
        'response_time': response_time,
        'hour_of_day': hour_of_day,
        'minute': minute
    })

# Criar DataFrame
df = pd.DataFrame(metrics_data)

# Certifique-se de que todos os dados de 'df' são apenas os reais do 'all_metrics'

# Adicionar uma coluna para o tempo de resposta (aqui você usaria os dados reais ou uma simulação baseada nas métricas)
#df['response_time'] = np.random.uniform(0.1, 5, len(df))  # Aqui ainda usamos uma variável simulada, se não houver dados reais

# ------------------- PREPARAÇÃO DOS DADOS -------------------

# Criar rótulos binários (1 = Alta Latência, 0 = Baixa Latência)
df['high_latency'] = (df['response_time'] > 3).astype(int)

# Separar variáveis independentes (X) e variável dependente (y)
X = df[['cpu_usage', 'disk_usage', 'memory_available', 'hour_of_day', 'minute']]  # Não usamos mais 'tasks_interval' se não estiver disponível
# y = df['high_latency']  # Variável de saída binária
y = df['response_time']  # O que queremos prever

# Dividir dados em treino e teste (aqui vamos manter apenas as amostras reais)
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# ------------------- TREINO DO MODELO -------------------

# Criar modelo de Regressão Linear
#model = LinearRegression() #TODO utilizar lazy predict
reg = LazyRegressor(verbose=0, ignore_warnings=True, custom_metric=None)
models, predictions = reg.fit(X_train, X_test, y_train, y_test)

print(models)

# Escolher o melhor modelo (primeira linha dos resultados)
best_model_name = models.index[0]
best_model = reg.models[best_model_name]  # Obter o modelo treinado

print(f"\nO melhor modelo encontrado: {best_model_name}")

# Treinar o modelo
#model.fit(X_train, y_train) # isto é código da professora

# ------------------- AVALIAÇÃO -------------------

# Prever no conjunto de teste
y_pred = best_model.predict(X_test)

# Avaliar o desempenho do modelo
mse = mean_squared_error(y_test, y_pred)
r2 = r2_score(y_test, y_pred)

print(f"Erro Quadrático Médio (MSE): {mse}")
print(f"R2 Score: {r2}")

# ------------------- PREVISÕES -------------------

# Fazer previsões para os nodes reais
future_data = df[['cpu_usage', 'disk_usage', 'memory_available', 'hour_of_day', 'minute']]  # Dados reais para previsões

# Fazer previsões para os nodes reais
future_predictions = best_model.predict(future_data)

print("\nPrevisões de Tempo de Resposta para Dados Reais (Nodes presentes no arquivo):")
for i, prediction in enumerate(future_predictions):
    # Garantir que as previsões não sejam negativas
    prediction = max(0, prediction)
    print(f"Node {df['node_id'].iloc[i]}: Tempo de Resposta Estimado = {prediction:.2f} segundos")

# ------------------- VISUALIZAÇÃO -------------------

# Visualizar previsão vs real
# Valores reais
plt.scatter(range(len(y_test)), y_test, color='blue', label='Valores Reais', alpha=0.7)

# Valores previstos
plt.scatter(range(len(y_pred)), y_pred, color='orange', label='Valores Previstos', alpha=0.7)

# Adicionar legendas e rótulos
plt.xlabel("Amostras")
plt.ylabel("Tempo de Resposta (s)")
plt.title("Previsão vs Real (Tempo de Resposta)")
plt.legend()
plt.savefig('/foo.png')
plt.show()

# ------------------- REDISTRIBUIÇÃO -------------------

# Redistribuir tarefas com base nas previsões
# Regras: Se previsão > 3s, redistribuir tarefas para o node com menor CPU
for i, row in enumerate(future_predictions):
    if row > 3:
        print(f"Redistribuir tarefas do Node {df['node_id'].iloc[i]}: Previsão de tempo de resposta alta ({row:.2f}s).")
    else:
        print(f"Node {df['node_id'].iloc[i]} opera dentro do limite ({row:.2f}s).")