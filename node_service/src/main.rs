use actix_web::{web, App, HttpServer, Responder, HttpResponse};
use std::fs;
use std::path::PathBuf;
use actix_web::web::Bytes;

async fn store_fragment(
    fragment_id: web::Path<String>,
    data: Bytes,
) -> impl Responder {
    let fragment_id = fragment_id.into_inner();
    let mut file_path = PathBuf::from("fragments");
    file_path.push(&fragment_id);

    fs::create_dir_all("fragments").unwrap();
    match fs::write(&file_path, data) {
        Ok(_) => HttpResponse::Ok().body("Fragmento armazenado"),
        Err(e) => {
            println!("Erro ao armazenar fragmento: {}", e);
            HttpResponse::InternalServerError().body("Erro ao armazenar fragmento")
        }
    }
}

async fn retrieve_fragment(
    fragment_id: web::Path<String>,
) -> impl Responder {
    let fragment_id = fragment_id.into_inner();
    let mut file_path = PathBuf::from("fragments");
    file_path.push(&fragment_id);

    match fs::read(&file_path) {
        Ok(data) => HttpResponse::Ok().body(data),
        Err(_) => HttpResponse::NotFound().body("Fragmento não encontrado"),
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Node iniciado na porta 8081");
    HttpServer::new(|| {
        App::new()
            .route("/store/{fragment_id}", web::post().to(store_fragment))
            .route("/retrieve/{fragment_id}", web::get().to(retrieve_fragment))
    })
    .bind(("0.0.0.0", 8081))?
    .run()
    .await
}
