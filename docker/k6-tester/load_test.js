import http from 'k6/http';
import { check, sleep } from 'k6'; 

export let options = {
    stages: [
        
        { duration: '1m', target: 2000 },
        
    ],
};

export default function () {
  const res = http.get('https://picloud.pt'); 
  
  
  check(res, { 'status was 200': (r) => r.status == 200 });

  sleep(1); // Dorme por 1 segundo entre as requisições
}