use std::env;

use async_trait::async_trait;
use ferrous_llm::{
    ollama::{ChatProvider, OllamaChatResponse, OllamaConfig, OllamaError, OllamaProvider},
    ChatRequest, ChatRequestBuilder,
};

use crate::agent::traits::Agent;

pub fn init_provider() -> Result<OllamaProvider, OllamaError> {
    let config = OllamaConfig::from_env()?;
    OllamaProvider::new(config)
}

pub struct AgentChatParameters {
    pub temperature: f32,
    pub max_tokens: u32,
    pub top_p: f32,
    pub request_id: String,
}

pub struct OllamaAgent<'a> {
    pub provider: &'a OllamaProvider,
    pub system_prompt: String,
    pub chat_parameters: AgentChatParameters,
}

impl<'a> OllamaAgent<'a> {
    pub fn new(
        provider: &'a OllamaProvider,
        system_prompt: String,
        chat_parameters: AgentChatParameters,
    ) -> OllamaAgent<'a> {
        OllamaAgent {
            provider,
            system_prompt,
            chat_parameters,
        }
    }
}

#[async_trait]
impl<'a> Agent for OllamaAgent<'a> {
    async fn prompt(
        &self,
        message: String,
        user_id: String,
    ) -> Result<OllamaChatResponse, OllamaError> {
        let request = ChatRequest::builder()
            .system_message(self.system_prompt.clone())
            .user_message(message)
            .temperature(self.chat_parameters.temperature)
            .max_tokens(self.chat_parameters.max_tokens)
            .top_p(self.chat_parameters.top_p)
            .request_id(self.chat_parameters.request_id.clone())
            .user_id(user_id)
            .build();
        (*self.provider).chat(request).await
    }
}
