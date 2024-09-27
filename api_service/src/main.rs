use actix_web::{web, App, HttpServer, Responder, HttpResponse, Error};
use serde::{Deserialize, Serialize};
use std::env;

// Estrutura para os dados de registro
#[derive(Deserialize)]
struct RegisterData {
    username: String,
    password: String,
}

// Função para a rota de registro
async fn register(data: web::Json<RegisterData>) -> Result<HttpResponse, Error> {
    // Lógica de registro (futuramente pode incluir banco de dados)
    println!("Registrando usuário: {}", data.username); // Para debugar

    Ok(HttpResponse::Ok().body(format!("Usuário {} registrado com sucesso!", data.username)))
}

// Função principal (inicializa o servidor)
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Obtém a porta das variáveis de ambiente ou usa 8080 por padrão
    let port = env::var("PORT").unwrap_or_else(|_| "8080".to_string());
    let port: u16 = port.parse().expect("PORT deve ser um número");

    HttpServer::new(|| {
        App::new()
            .route("/", web::get().to(index))              // Rota inicial
            .route("/register", web::post().to(register))  // Rota de registro
    })
    .bind(("0.0.0.0", port))?
    .run()
    .await
}

// Função para a rota inicial
async fn index() -> impl Responder {
    HttpResponse::Ok().body("API do PiCloud")
}