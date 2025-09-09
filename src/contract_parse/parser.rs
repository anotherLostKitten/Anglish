#![feature(substr_range)]

use super::token_stream;
use super::token_stream::TokenStream;
use super::ast;

pub enum ParserErrorType {
    UnexpectedMetachar,
    SpacesAtDeclLineStart,
    ExpectedOuterDecl,
    ExpectedSpaceDecl,
    NoSpacesInContract,
    SpaceDeclMissingName,
    UnknownSpaceType,
}

pub struct ParserError {
    pub error: ParserErrorType,
    pub line: usize,
    pub col: usize,
}

pub struct ParserErrorSet {
    errors: Vec<ParserError>,
}

impl ParserErrorSet {
    fn add_error(&mut self, error: ParserErrorType, line: usize, col: usize) {
        self.errors.push(ParserError {
            error: error,
            line: line,
            col: col,
        });
    }
}

struct ParseContext<T: TokenStream> {
    ts: T,

    errors: ParserErrorSet,

    cur_line: usize,
    cur_col: usize,
}

impl<T: TokenStream> ParseContext<T> {
    fn new(ts: T) -> ParseContext<T> {
        ParseContext {
            ts: ts,
            errors: ParserErrorSet {
                errors: vec![],
            },
            cur_line: 0,
            cur_col: 0,
        }
    }

    fn add_error(&mut self, error_type: ParserErrorType) {
        self.errors.add_error(error_type, self.cur_line, self.cur_col);
    }

    fn next_char(&mut self) -> Option<char> {
        match self.ts.next() {
            None => None,
            Some('\n') => {
                self.cur_line += 1;
                self.cur_col = 0;
                Some('\n')
            },
            c => {
                self.cur_col += 1;
                c
            },
        }
    }

    fn next_identifier(&mut self) -> Option<&str> {
        let start_pos = self.ts.get_pos();
        if let Some(x) = self.ts.next() {
            match x {
                'a'..'z' | 'A'..'Z' | '_' => {},
                _ => {
                    self.ts.seek_back(1);
                    return None;
                },
            };
        } else {
            return None;
        }
        while if let Some(c) = self.ts.next() {
            match c {
                'a'..'z' | 'A'..'Z' | '_' | '0'..'9' => true,
                _ => {
                    self.ts.seek_back(1);
                    false
                },
            }
        } else { false } {}
        let end_pos = self.ts.get_pos();
        self.cur_col += end_pos - start_pos;
        return Some(self.ts.get_slice(start_pos, end_pos));
    }

    fn consume_spaces(&mut self) {
        let old_pos = self.ts.get_pos();
        while match self.ts.next() {
            Some(' ') | Some('\t') => true,
            None => false,
            _ => {
                self.ts.seek_back(1);
                false
            },
        } {}

        self.cur_col += self.ts.get_pos() - old_pos;
    }

    fn consume_blank_lines(&mut self) {
        let mut look_ahead = 0;
        loop {
            match self.ts.look_ahead(look_ahead) {
                Some(' ') | Some('\t')  => {
                    look_ahead += 1;
                },
                Some('\n') => {
                    self.ts.seek_forward(look_ahead + 1);
                    look_ahead = 0;
                    self.cur_line += 1;
                    self.cur_col = 0;
                },
                None => {
                    self.ts.seek_forward(look_ahead);
                    self.cur_line += 1;
                    self.cur_col = 0;
                    return;
                },
                _ => return,
            }
        }
    }

    fn next_line(&mut self) -> Option<&str> {
        if !self.ts.has_next() {
            return None;
        }
        let start_pos = self.ts.get_pos();
        let mut end_pos = 0;
        loop {
            match self.ts.next() {
                None => { // cursor doesn't advance past file end, so if we `pos - 1` we ellide last character.
                    end_pos = self.ts.get_pos();
                    break;
                },
                Some('\n') => { // we dont go back so we can consume the '\n', but `pos - 1` so dont include '\n' in the &str.
                    end_pos = self.ts.get_pos() - 1;
                    break;
                },
                _ => {},
            }
        }
        self.cur_line += 1;
        self.cur_col = 0;
        return Some(self.ts.get_slice(start_pos, end_pos));
    }

    fn parse_contract(&mut self) -> Option<ast::Contract> {
        let mut root = ast::Contract::new();

        while self.ts.has_next() {
            self.consume_blank_lines();
            match self.ts.look_ahead(0) {
                None => {},
                Some('@') => {
                    if let Some(n) = self.parse_space_decl() {
                        root.spaces.push(n);
                    }
                },
                Some('#') => {
                    if let Some(n) = self.parse_agent_decl() {
                        root.agents.push(n);
                    }
                },
                Some('=') => {
                    if let Some(n) = self.parse_path_decl() {
                        root.paths.push(n);
                    }
                },
                Some(' ') | Some('\t') => {
                    self.add_error(ParserErrorType::SpacesAtDeclLineStart);
                    self.consume_spaces();
                },
                _ => {
                    self.add_error(ParserErrorType::ExpectedOuterDecl);
                    self.next_line();
                },
            }
        }

        if root.spaces.len() > 0 {
            return Some(root);
        } else {
            self.add_error(ParserErrorType::NoSpacesInContract);
            return None;
        }
    }

    fn parse_space_decl(&mut self) -> Option<ast::SpaceDecl> {
        let line_start = self.cur_line;
        let col_start = self.cur_col;
        let pos_start = self.ts.get_pos();

        if self.next_char() != Some('@') {
            self.add_error(ParserErrorType::ExpectedSpaceDecl);
            return None;
        }

        let space_name = if let Some(x) = self.next_identifier() { String::from(x) } else {
            self.add_error(ParserErrorType::SpaceDeclMissingName);
            return None;
        };

        self.consume_spaces();

        let space_type = if self.ts.look_ahead(0) == Some(':') {
            _ = self.next_char();
            self.consume_spaces();
            if let Some(s) = self.next_identifier() {
                let res = match s {
                    "UI" | "ui" => ast::SpaceType::UI,
                    "IO" | "io" => ast::SpaceType::IO,
                    "DATA" | "Data" | "data" => ast::SpaceType::DATA,
                    "FUNC" | "Func" | "func" => ast::SpaceType::FUNC,
                    "AGENTIC" | "Agentic" | "agentic" => ast::SpaceType::AGENTIC,
                    _ => {
                        self.add_error(ParserErrorType::UnknownSpaceType);
                        ast::SpaceType::Unknown
                    },
                };
                self.consume_spaces();
                res
            } else {
                self.add_error(ParserErrorType::UnknownSpaceType);
                ast::SpaceType::Unknown
            }
        }else {
            ast::SpaceType::Unknown
        };

        _ = self.next_line();

        let vibe_desc = if let Some(x) = self.parse_vibe_block() {x} else {
            return None;
        };

        let space = ast::SpaceDecl {
            ident: space_name,
            space_type: space_type,
            params: vec![],
            vibe_desc: vibe_desc,

            line_start: line_start,
            line_end: self.cur_line,
        };

        return Some(space);
    }

    fn parse_agent_decl(&mut self) -> Option<ast::AgentDecl> {
        todo!();
    }

    fn parse_path_decl(&mut self) -> Option<ast::PathDecl> {
        todo!();
    }

    fn parse_vibe_block(&mut self) -> Option<ast::VibeBlock> {
        todo!();
    }
}

pub fn parse_from_str(source: &str) -> (Option<ast::Contract>, ParserErrorSet) {
    let stream = token_stream::TokenStreamString::new(source);

    return parse_from_ts(stream);
}

fn parse_from_ts<T: TokenStream>(ts: T) -> (Option<ast::Contract>, ParserErrorSet) {
    let mut pctx = ParseContext::new(ts);

    let contract = pctx.parse_contract();

    return (contract, pctx.errors);
}
