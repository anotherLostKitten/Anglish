pub mod ollama_agent;

use async_trait::async_trait;
use ferrous_llm::ollama::{OllamaChatResponse, OllamaError};

#[async_trait]
pub trait Agent {
    async fn prompt(
        &self,
        message: String,
        user_id: String,
    ) -> Result<OllamaChatResponse, OllamaError>;
}
