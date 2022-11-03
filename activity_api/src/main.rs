// #[macro_use] extern crate rocket;
// use serde::{Serialize, Deserialize};

// #[get("/latest")]
// // #[tokio::main]
// async fn latest() -> &'static str {
//     let access_token = format!("Bearer {token}", token={your_token_here}
//     let client = reqwest::Client::new();
//     let resp = client
//         .get("https://guc-spclient.spotify.com/presence-view/v1/buddylist")
//         .header("Authorization", access_token)
//         // .header(CONTENT_TYPE, "application/json")
//         // .header(ACCEPT, "application/json")
//         .send()
//         .await
//         .unwrap();
//
//     // let body = resp.text().await?;
//     println!("{:?}", resp);
//     return resp.text();
//     // Ok(())
// }
//
// #[launch]
// fn rocket() -> _ {
//     rocket::build().mount("/latest", routes![latest])
// }

// extern crate serde_json;
// use std::collections::HashMap;

// use rocket::serde::json::Json;
use serde::{Deserialize, Serialize};

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct FriendActivity {
    pub friends: Vec<Friend>,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Friend {
    pub timestamp: i64,
    pub user: User,
    pub track: Track,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct User {
    pub uri: String,
    pub name: String,
    pub image_url: String,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Track {
    pub uri: String,
    pub name: String,
    pub image_url: String,
    pub album: Album,
    pub artist: Artist,
    pub context: Context,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Album {
    pub uri: String,
    pub name: String,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Artist {
    pub uri: String,
    pub name: String,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Context {
    pub uri: String,
    pub name: String,
    pub index: i64,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Token {
    pub client_id: String,
    pub access_token: String,
    pub access_token_expiration_timestamp_ms: i64,
    pub is_anonymous: bool,
}

fn main() {
    let token = {your_token_here}
    match refresh_token(token) {
        Ok(access_token) => {
            println!("{}", access_token);
            match latest(access_token) {
                Ok(activity) => {
                    println!("{}", activity);
                }
                Err(e) => {
                    println!("Error: {}", e)
                }
            }
        }
        Err(e) => {
            println!("Error: {}", e)
        }
    };

    // let text = serde_json::to_string(&addr);
    // println!("{}", addr)
}

// use std::error::Error;

#[tokio::main]
async fn latest(access_token: String) -> Result<String, Box<dyn std::error::Error>> {
    let access_token = format!("Bearer {token}", token = access_token);
    let client = reqwest::Client::new();
    let resp = client
        .get("https://guc-spclient.spotify.com/presence-view/v1/buddylist")
        .header("Authorization", access_token)
        .header("Content-Type", "application/json")
        .header("Accept", "application/json")
        .send()
        .await?;
    let status = resp.status();
    let res: &str = &resp.text().await?;
    match status {
        reqwest::StatusCode::OK => {
            println!("No errors!")
        }
        reqwest::StatusCode::UNAUTHORIZED => {
            println!("Need to grab a new token");
        }
        other => {
            panic!("Uh oh! Something unexpected happened: {:?}", other);
        }
    };
    let json = match serde_json::from_str(res) {
        Ok(parsed) => parsed,
        Err(e) => println!("Error: {}", e),
    };
    Ok(serde_json::to_string(&json)?)
}

#[tokio::main]
async fn refresh_token(sp_dc: &str) -> Result<String, Box<dyn std::error::Error>> {
    // let access_token = format!("sp_dc={token}", token = sp_dc);
    let access_token = format!("sp_dc={token}", token={your_token_here}
    let client = reqwest::Client::new();
    let resp = client
        .get("https://open.spotify.com/get_access_token?reason=transport&productType=web_player")
        .header("Cookie", access_token)
        .send()
        .await?;
    let status = resp.status();
    let res: &str = &resp.text().await?;
    println!("{}", res);
    match status {
        reqwest::StatusCode::OK => {
            println!("No errors!")
        }
        reqwest::StatusCode::UNAUTHORIZED => {
            println!("Need to grab a new token");
        }
        other => {
            panic!("Uh oh! Something unexpected happened: {:?}", other);
        }
    };
    let json = match serde_json::from_str(res) {
        Ok(parsed) => parsed,
        Err(e) => println!("Error: {}", e),
    };
    Ok(serde_json::to_string(&json)?)
}
