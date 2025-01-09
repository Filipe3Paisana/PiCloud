import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
    stages: [
        { duration: '2m', target: 1000 }, // Sobe para 100 usuários em 2 minutos
        { duration: '5m', target: 1000 }, // Mantém 100 usuários por 5 minutos
        { duration: '2m', target: 5000 }, // Escala para 500 usuários em 2 minutos
        { duration: '10m', target: 5000 }, // Mantém 500 usuários por 10 minutos
        { duration: '2m', target: 0 },   // Reduz para 0 usuários em 2 minutos
    ],
    thresholds: {
        http_req_duration: ['p(95)<2000'], // 95% das requisições devem ser menores que 2s
        http_req_failed: ['rate<0.01'],    // Taxa de falhas menor que 1%x
    },
};

export default function () {
    const BASE_URL = 'https://picloud.pt/'; // Substitua pelo seu domínio

    const responses = http.batch([
        ['GET', `${BASE_URL}/path1`, null, { tags: { name: 'Path1' } }],
        ['GET', `${BASE_URL}/path2`, null, { tags: { name: 'Path2' } }],
        ['POST', `${BASE_URL}/path3`, { key: 'value' }, { tags: { name: 'Path3' } }],
    ]);

    responses.forEach((res, i) => {
        if (res.status !== 200) {
            console.error(`Erro na requisição ${i + 1}: ${res.status}`);
        }
    });

    sleep(1); // Simula tempo entre as ações dos usuários
}
