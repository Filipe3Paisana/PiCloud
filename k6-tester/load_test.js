import http from 'k6/http';
import { check, sleep } from 'k6'; // Adiciona a importação de 'check'

export let options = {
    stages: [
        
        { duration: '1m', target: 5000 },  // Aumenta até 5000 usuários em 1 minuto
        
    ],
};

export default function () {
  const res = http.get('http://nginx-container:80');
  
  // Verifica se a resposta tem o status 200
  check(res, { 'status was 200': (r) => r.status == 200 });

  sleep(1); // Dorme por 1 segundo entre as requisições
}
