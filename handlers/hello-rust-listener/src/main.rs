use std::io::{BufRead, BufReader, BufWriter, Write};
use std::net::{TcpListener, TcpStream};
use serde_json::{json, Value};
use std::thread;

fn main() {
    let listener = TcpListener::bind("127.0.0.1:9123").unwrap();
    for stream in listener.incoming() {
        thread::spawn(move || {
            let stream = stream.unwrap();
            let mut reader = BufReader::new(&stream);
            let writer = BufWriter::new(&stream);
            let mut buffer = String::new();

            match reader.read_line(&mut buffer) {
                Ok(_) => handle_request(writer, buffer),
                Err(e) => panic!("encountered IO error: {}", e),
            };            
        });
    }
}

fn handle_request(mut writer: BufWriter<&TcpStream>, request: String) {
    let req: Value = serde_json::from_str(&request[..]).unwrap();

    let response = json!({
        "status": 200,
        "header": {},
        "body": format!("hello {}", req["Params"]["name"][0].as_str().unwrap())
    }); 

    let s = format!("{}\r\n", response.to_string());
    let d = s.as_bytes();

    writer.write(d).unwrap();
}