use std::env;
use serde_json::{Value};
use serde_json::json;

fn main() {
    let args: Vec<String> = env::args().collect();
    let request: Value = serde_json::from_str(&args[1][..]).unwrap();

    let response = json!({
        "status": 200,
        "header": {},
        "body": format!("hello {}", request["Params"]["name"][0].as_str().unwrap())
    });    

    println!("{}", response.to_string());
    
}


