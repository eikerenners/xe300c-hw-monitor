use reqwest::header::{HeaderMap, HeaderValue};
use serde_json::Value;
use std::error::Error;

async fn fetch_data() -> Result<Value, Box<dyn Error>> {
    let client = reqwest::Client::new();
    let mut headers = HeaderMap::new();
    headers.insert("glinet", HeaderValue::from_static("1"));

    let res = client
        .post("http://127.0.0.1/rpc")
        .headers(headers)
        .body(r#"{"jsonrpc":"2.0","id":1,"method":"call","params":["","model","get_cells_info"]}"#)
        .send()
        .await?;

    let json: Value = res.json().await?; 
    Ok(json)
}



#[tokio::main]
async fn main() {
    match fetch_data().await {
        Ok(json) => println!("Response JSON: {:?}", json),
        Err(e) => eprintln!("Error: {}", e),
    }
}
