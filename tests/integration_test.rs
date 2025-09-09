use anglish::agent::{
    ollama_agent::{init_provider, AgentChatParameters, OllamaAgent},
    Agent,
};
use log::info;

#[tokio::test]
#[test_log::test]
#[ignore = "needs ollama running"]
async fn test_response_from_ollama_agent() -> Result<(), ()> {
    unsafe {
        std::env::set_var("OLLAMA_BASE_URL", "http://localhost:11434");
        std::env::set_var("OLLAMA_MODEL", "llama3.2:1b");
        std::env::set_var("OLLAMA_KEEP_ALIVE", "10");
    }
    let provider = init_provider();

    assert!(provider.is_ok());

    let chat_parameters = AgentChatParameters {
        temperature: 0.7,
        max_tokens: 50,
        top_p: 1.0,
        request_id: "example-ollama-chat-001".to_string(),
    };

    let provider_ref = &provider.unwrap();

    let agent = OllamaAgent::new(
        provider_ref,
        "You are a helpful chatbot!".to_string(),
        chat_parameters,
    );

    let possible_response = agent
        .prompt("Hello!".to_string(), "example-user".to_string())
        .await;

    match possible_response {
        Err(e) => panic!("{:?}", e),
        Ok(response) => info!("{}", response.message.content),
    };
    Ok(())
}
