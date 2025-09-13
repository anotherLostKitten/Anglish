use rig::{
    agent::{Agent, AgentBuilder},
    client::{CompletionClient, ProviderClient},
    providers::openai::{Client, CompletionModel},
};

pub fn init_llm() -> AgentBuilder<CompletionModel> {
    Client::from_env()
        .completion_model(&std::env::var("OPENAI_MODEL").expect("OPENAI_MODEL must be set!"))
        .completions_api()
        .into_agent_builder()
}

pub fn build_agent(
    llm: AgentBuilder<CompletionModel>,
    optional_temperature: Option<f64>,
    optional_max_tokens: Option<u64>,
    optional_system_prompt: Option<&str>,
) -> Agent<CompletionModel> {
    let temperature = optional_temperature.unwrap_or(1.0);
    let max_tokens = optional_max_tokens.unwrap_or(1024);
    match optional_system_prompt {
        Some(system_prompt) => llm.preamble(system_prompt),
        None => llm,
    }
    .temperature(temperature)
    .max_tokens(max_tokens)
    .build()
}

#[cfg(test)]
mod tests {

    use rig::completion::Prompt;

    use crate::llm::agent::{build_agent, init_llm};

    #[tokio::test]
    #[test_log::test]
    #[ignore = "needs vllm running"]
    async fn test_simple_prompt() {
        unsafe {
            std::env::set_var("OPENAI_BASE_URL", "http://localhost:8000/v1");
            std::env::set_var("OPENAI_API_KEY", "EMPTY");
            std::env::set_var("OPENAI_MODEL", "google/gemma-3-4b-it");
        }

        let llm = init_llm();
        let agent = build_agent(llm, None, None, Some("You are a helpful assistant!"));
        match agent.prompt("Hello!").await {
            Ok(response) => println!("Result: {}", response),
            Err(e) => panic!("{e}"),
        }
    }
}
