pub enum ParserErrorType {
    UNEXPECTED_METACHAR,
}

pub struct ParserError {
    pub error: ParserErrorType,
    pub line: usize,
    pub col: usize,
}

pub fn parse() -> Result<(), ParserError> {
    return Ok(());
}
