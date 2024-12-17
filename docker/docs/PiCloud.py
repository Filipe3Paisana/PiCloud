
#####################################################################################
# Otimização da Atividade Computacional dos Nodes
#####################################################################################

#1.  Machine Learning Supervisionado (Regressão):
# Objectivo: prever quais nodes estarão sobrecarregados em horários específicos e redistribuir tarefas para otimizar o desempenho.
# Algoritmo: Regressão Linear
    #Dados Necessários:
    #Logs históricos de atividade computacional dos nodes (CPU, memória, tempo de resposta).
    #Horários específicos e tarefas processadas por node.
    #Estados de congestionamento passados.
#Processo:
    #Usar os dados para treinar um modelo que prevê a carga futura de cada node em função do tempo e do tipo de tarefa.
    #Integrar o modelo para redistribuir dinamicamente tarefas baseando-se nas previsões.


import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_squared_error, r2_score
import matplotlib.pyplot as plt

# Dados simulados para exemplo
data = {
    'cpu_usage': np.random.uniform(20, 100, 1000),  # % de CPU
    'disk_usage': np.random.uniform(10, 500, 1000),  # Ocupação em MG
    'memory_available': np.random.uniform(10, 90, 1000),  # % RAM disponível
    'uptime': np.random.uniform(1, 100, 1000),  # Em horas
    'response_time': np.random.uniform(0.1, 5, 1000),  # Tempo de resposta em segundos
    'tasks_interval': np.random.randint(1, 50, 1000),  # Número de tarefas por intervalo
    'congestion': np.random.randint(0, 2, 1000)  # Congestionamento (0 = não, 1 = sim)
}

# Criar DataFrame
df = pd.DataFrame(data)

# Adicionar intervalo de horários fictícios
df['interval'] = pd.date_range(start='2023-01-01', periods=len(df), freq='h')

# ------------------- PREPARAÇÃO DOS DADOS -------------------

# Separar variáveis independentes (X) e variável dependente (y)
X = df[['cpu_usage', 'disk_usage', 'memory_available', 'uptime', 'tasks_interval']]
y = df['response_time']  # O que queremos prever

# Dividir dados em treino e teste
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# ------------------- TREINO DO MODELO -------------------

# Criar modelo de Regressão Linear
model = LinearRegression()

# Treinar o modelo
model.fit(X_train, y_train)

# ------------------- AVALIAÇÃO -------------------

# Prever no conjunto de teste
y_pred = model.predict(X_test)

# Avaliar o desempenho do modelo
mse = mean_squared_error(y_test, y_pred)
r2 = r2_score(y_test, y_pred)

print(f"Erro Quadrático Médio (MSE): {mse}")
print(f"R2 Score: {r2}")

# ------------------- PREVISÕES -------------------

# Dados futuros simulados para previsão
future_data = {
    'cpu_usage': [70, 80, 65],  # Exemplos de % de CPU
    'disk_usage': [200, 250, 150],  # Ocupação em MG
    'memory_available': [30, 20, 40],  # % RAM disponível
    'uptime': [50, 70, 90],  # Em horas
    'tasks_interval': [20, 30, 25]  # Número de tarefas
}

future_df = pd.DataFrame(future_data)

# Fazer previsões para novos dados
future_predictions = model.predict(future_df)

print("\nPrevisões de Tempo de Resposta para Dados Futuros:")
for i, prediction in enumerate(future_predictions):
    print(f"Node {i+1}: Tempo de Resposta Estimado = {prediction:.2f} segundos")

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
plt.show()

# ------------------- REDISTRIBUIÇÃO -------------------

# Redistribuir tarefas com base nas previsões
# Regras: Se previsão > 3s, redistribuir tarefas para o node com menor CPU
for i, row in future_df.iterrows():
    if future_predictions[i] > 3:
        print(f"Redistribuir tarefas do Node {i+1}: Previsão de tempo de resposta alta ({future_predictions[i]:.2f}s).")
    else:
        print(f"Node {i+1} opera dentro do limite ({future_predictions[i]:.2f}s).")
