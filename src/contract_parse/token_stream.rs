pub trait TokenStream {
    fn next(&mut self) -> Option<char>;
    fn look_ahead(&self, n: usize) -> Option<char>;
    fn get_pos(&self) -> usize;
    fn back(&mut self, n: usize);
    fn get_slice(&self, start: usize, end: usize) -> &str;

    fn next_identifier(&mut self) -> Option<&str> {
        let start_pos = self.get_pos();
        if let Some(x) = self.next() {
            match x {
                'a'..'z' | 'A'..'Z' | '_' => {},
                _ => return None,
            };
        } else {
            return None;
        }
        while match self.next() {
            Some(x) => match x {
                'a'..'z' | 'A'..'Z' | '_' | '0'..'9' => true,
                _ => false,
            },
            None => false,
        } {}
        self.back(1);
        let end_pos = self.get_pos();
        return Some(self.get_slice(start_pos, end_pos));
    }

    fn consume_spaces(&mut self) {
        while match self.next() {
            Some(' ') | Some('\t') => true,
            _ => false,
        } {}
        self.back(1);
    }

    fn next_line(&mut self) -> Option<&str> {
        if self.look_ahead(0) == None {
            return None;
        }
        let start_pos = self.get_pos();
        while match self.next() {
            Some('\n') | None => false,
            _ => true,
        }

    }
}

pub struct TokenStreamString<'a> {
    source: &'a str,
    source_chars: Vec<(usize, char)>,
    pos: usize,
}

impl TokenStream for TokenStreamString<'_> {
    fn next(&mut self) -> Option<char> {
        if self.pos < self.source_chars.len() {
            self.pos += 1;
            return Some(self.source_chars[self.pos - 1].1);
        }
        return None;
    }

    fn look_ahead(&self, n: usize) -> Option<char> {
        if self.pos + n < self.source_chars.len() {
            return Some(self.source_chars[self.pos + n].1);
        }
        return None;
    }

    #[inline]
    fn get_pos(&self) -> usize {
        self.pos
    }

    fn back(&mut self, n: usize) {
        if self.pos >= n {
            self.pos -= n;
        } else {
            self.pos = 0;
        }
    }

    fn get_slice(&self, start: usize, end: usize) -> &str {
        let s_byte = self.source_chars[start].0;
        let e_byte = self.source_chars[end].0;

        return &self.source[s_byte..e_byte];
    }
}

impl<'a> TokenStreamString<'a> {
    pub fn new(source: &'a str) -> TokenStreamString<'a> {
        TokenStreamString {
            source: source,
            source_chars: source.char_indices().collect(),
            pos: 0,
        }
    }
}
