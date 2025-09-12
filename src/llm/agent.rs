use std::sync::Arc;

use langchain_rust::{
    agent::{AgentError, AgentExecutor, OpenAiToolAgent, OpenAiToolAgentBuilder},
    chain::options::ChainCallOptions,
    llm::{OpenAI, OpenAIConfig},
    memory::SimpleMemory,
    tools::Tool,
};

pub fn init_llm() -> OpenAI<OpenAIConfig> {
    OpenAI::default()
        .with_config(
            OpenAIConfig::default()
                .with_api_base(
                    std::env::var("OPENAI_BASE_URL").expect("OPENAI_BASE_URL must be set!"),
                )
                .with_api_key(
                    std::env::var("OPENAI_API_KEY").expect("OPENAI_API_KEY must be set!"),
                ),
        )
        .with_model(std::env::var("OPENAI_MODEL").expect("OPENAI_MODEL must be set!"))
}

pub fn build_agent(
    llm: &OpenAI<OpenAIConfig>,
    optional_temperature: Option<f32>,
    optional_max_tokens: Option<u32>,
    optional_tools: Option<&[Arc<dyn Tool>]>,
    optional_memory: Option<SimpleMemory>,
) -> Result<AgentExecutor<OpenAiToolAgent>, AgentError> {
    let temperature = optional_temperature.unwrap_or(1.0);
    let max_tokens = optional_max_tokens.unwrap_or(150);
    let options = ChainCallOptions::new()
        .with_temperature(temperature)
        .with_max_tokens(max_tokens);
    let agent_or_error = match optional_tools {
        Some(tools) => OpenAiToolAgentBuilder::new()
            .tools(tools)
            .options(options)
            .build(llm.clone()),
        None => OpenAiToolAgentBuilder::new()
            .options(options)
            .build(llm.clone()),
    };
    match agent_or_error {
        Ok(agent) => match optional_memory {
            Some(memory) => Ok(AgentExecutor::from_agent(agent).with_memory(memory.into())),
            None => Ok(AgentExecutor::from_agent(agent)),
        },
        Err(e) => Err(e),
    }
}

#[tokio::test]
#[test_log::test]
#[ignore = "needs vllm running"]
async fn test_initializing_provider() {
    unsafe {
        std::env::set_var("OPENAI_BASE_URL", "http://localhost:8080");
        std::env::set_var("OLLAMA_MODEL", "llama3.2:1b");
    }
}
