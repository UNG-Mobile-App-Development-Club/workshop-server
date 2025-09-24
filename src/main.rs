use axum::extract::{Path, State};
use axum::http::StatusCode;
use axum::routing::{get, post};
use axum::{Json, Router};
use axum_extra::response::ErasedJson;
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use std::time::SystemTime;
use std::usize;
use tokio::sync::RwLock;

// Define our todo_ struct, deriving Serialize and Deserialize from serde
#[derive(Serialize, Deserialize, Debug)]
struct Todo {
    id: Option<usize>,
    content: String,
    created_at: Option<SystemTime>,
}
// Annotate Todo_ type as being a Vec of Todos
type Todos = Vec<Todo>;

// Atomically Reference Counted (ARC) wrapper for sharing Todos list across threads,
// using RwLock to control access to reading and writing data
type TodoStore = Arc<RwLock<Todos>>;

#[tokio::main]
async fn main() {
    // In-memory todos list
    let todos = vec![
        Todo {
            id: Some(0),
            content: "Learn Rust".to_string(),
            created_at: Some(SystemTime::now()),
        },
        Todo {
            id: Some(1),
            content: "Build a web app".to_string(),
            created_at: Some(SystemTime::now()),
        },
    ];

    // Initialize tracing (logging)
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::DEBUG)
        .init();

    // Define routes with HTTP function (get, post)
    let app = Router::new()
        .route("/todos", get(get_todos))
        .route("/todos/{id}", post(post_todo).get(get_todo))
        .layer(tower_http::trace::TraceLayer::new_for_http()) // Add tracing layer
        .with_state(Arc::new(RwLock::new(todos)));

    // Listen and Serve on 127.0.0.1:8080
    let listener = tokio::net::TcpListener::bind("127.0.0.1:8080")
        .await
        .unwrap();
    tracing::info!("Server listening on 127.0.0.1:8080");
    axum::serve(listener, app).await.unwrap();
}

async fn get_todos<'a>(State(todos): State<TodoStore>) -> ErasedJson {
    // Using ErasedJson over Json(T) to avoid cloning data
    ErasedJson::new(&*todos.read().await)
}

async fn post_todo(State(todos): State<TodoStore>, Json(mut todo): Json<Todo>) -> ErasedJson {
    // Acquire write lock on todos list
    let mut todos = todos.write().await;
    todo.id = Some(todos.len() + 1);
    todo.created_at = Some(SystemTime::now());
    // Use ErasedJson over Json(T) to avoid cloning
    let response = ErasedJson::new(&todo);
    todos.push(todo);
    response
}

async fn get_todo(
    Path(id): Path<usize>,
    State(todos): State<TodoStore>,
) -> Result<ErasedJson, StatusCode> {
    // Path(id): Path<usize> automatically extracts the id from the URL path, returning BadRequest if it fails
    todos
        .read()
        .await
        .get(id)
        .map(|todo| ErasedJson::new(todo))
        .ok_or(StatusCode::NOT_FOUND)
}
