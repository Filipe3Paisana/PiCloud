#[macro_use]
extern crate rocket;

use rocket::serde::{json::Json, Deserialize};

#[derive(Deserialize)]
struct CreateUser {
    name: String,
    email: String,
}

#[post("/create_user", format = "json", data = "<user>")]
async fn create_user(user: Json<CreateUser>) -> String {
    format!("Usuário criado na API 1: {} - {}", user.name, user.email)
}

#[launch]
fn rocket() -> _ {
    rocket::build().mount("/", routes![create_user])
}
