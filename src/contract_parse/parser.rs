use super::token_stream;
use super::token_stream::TokenStream;

pub enum ParserErrorType {
    UnexpectedMetachar,
}

pub struct ParserError {
    pub error: ParserErrorType,
    pub line: usize,
    pub col: usize,
}

pub fn parse_from_str(source: &str) -> Result<(), ParserError> {
    let mut stream = token_stream::TokenStreamString::new(source);
    let first_ident = stream.next_identifier();
    if let Some(id) = first_ident {
        println!("ident found: {}", id);
    }

    return Ok(());
}
