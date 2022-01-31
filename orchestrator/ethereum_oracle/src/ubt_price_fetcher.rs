use serde_json::Value;
use std::env;

async fn fetch_ubt_price() -> Result<f32, Box<dyn std::error::Error>> {
    let mut token_price = 0f32;
    let coinmarketcap_result = fetch_from_coinmarket_cap().await;
    match coinmarketcap_result {
        Ok(v) => {
            token_price = v;
            println!("Token price is: {:?}", v);
        },
        Err(e) => { println!("Error: {:?}", e) };,
    }

    if token_price == 0f32 {
        let coingecko_result = fetch_from_coingecko().await;
        match coingecko_result {
            Ok(v) => {
                token_price = v;
                println!("Token price is: {:?}", v);
            },
            Err(e) => { println!("Error: {:?}", e) };,
        }
    }

    if token_price == 0f32 {
        let coinapi_result = fetch_from_coinpaprika().await;
        match coinapi_result {
            Ok(v) => {
                token_price = v;
                println!("Token price is: {:?}", v);
            },
            Err(e) => { println!("Error: {:?}", e) };,
        }
    }

    if token_price == 0f32 {
        let coinapi_result = fetch_from_coinapi().await;
        match coinapi_result {
            Ok(v) => {
                token_price = v;
                println!("Token price is: {:?}", v);
            },
            Err(e) => { println!("Error: {:?}", e) };,
        }
    }

    return Ok(token_price);
}

async fn fetch_from_coinmarket_cap() -> Result<f32, Box<dyn std::error::Error>> {
    println!("Fetching coinmarketcap price");
    
    let api_key_var_name = "COINMARKETCAP_API_TOKEN");
    let token = "";
    match env::var(api_key_var_name) {
        Ok(v) => token = v,
        Err(e) => panic!("${} is not set ({})", api_key_var_name, e)
    }

    let url = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=UBT&convert=EUR&CMC_PRO_API_KEY=";
    let full_url = format!("{}\n{}", url, token);

    let body: Value = perform_request_and_get_body(&full_url).await?;
    let price_str = body["data"]["UBT"]["quote"]["EUR"]["price"].to_string();
    let price_decimal = try_parse_price_string(&price_str);

    return Ok(price_decimal);
}

async fn fetch_from_coingecko() -> Result<f32, Box<dyn std::error::Error>> {
    println!("Fetching coingecko price");
    let url = "https://api.coingecko.com/api/v3/simple/price?ids=unibright&vs_currencies=EUR";

    let body: Value = perform_request_and_get_body(url).await?;

    let price_str = body["unibright"]["eur"].to_string();
    let price_decimal = try_parse_price_string(&price_str);

    return Ok(price_decimal);
}

async fn fetch_from_coinpaprika() -> Result<f32, Box<dyn std::error::Error>> {
    println!("Fetching coinpaprika price");

    let url = "https://api.coinpaprika.com/v1/price-converter?base_currency_id=ubt-unibright&quote_currency_id=eur-euro&amount=1";

    let body: Value = perform_request_and_get_body(url).await?;

    let price_str = body["price"].to_string();
    let price_decimal = try_parse_price_string(&price_str);

    return Ok(price_decimal);
}

async fn fetch_from_coinapi() -> Result<f32, Box<dyn std::error::Error>> {
    println!("Fetching coinapi price");
    let url = "https://rest.coinapi.io/v1/exchangerate/UBT/EUR";

    let api_key_var_name = "COINAPI_API_TOKEN");
    let token = "";
    match env::var(api_key_var_name) {
        Ok(v) => token = v,
        Err(e) => panic!("${} is not set ({})", api_key_var_name, e)
    }
    
    let client = reqwest::Client::new();
    let res = client
        .get(url)
        .header("X-CoinAPI-Key", token)
        .send()
        .await?;
    
    let body = res.text().await?;
    let v: Value = serde_json::from_str(&body)?;

    let price_str = v["rate"].to_string();
    let price_decimal = try_parse_price_string(&price_str);
    return Ok(price_decimal);
}

async fn perform_request_and_get_body(url: &str) -> Result<Value, Box<dyn std::error::Error>>  {
    let res = reqwest::get(url).await?;
    let body = res.text().await?;
    let v: Value = serde_json::from_str(&body)?;

    return Ok(v);
}

fn try_parse_price_string(price_string: &str) -> f32 {
    let price_decimal_result = price_string.parse();
    match price_decimal_result {
        Ok(res) => return res,
        Err(error) => println!("Problem parsing price string: {:?}, {:?}", price_string, error),
    };

    return 0f32;
}