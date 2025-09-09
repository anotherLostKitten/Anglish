pub trait TokenStream {
    fn has_next(&self) -> bool;
    fn next(&mut self) -> Option<char>;
    fn look_ahead(&self, n: usize) -> Option<char>;
    fn get_pos(&self) -> usize;
    fn seek_forward(&mut self, n: usize);
    fn seek_back(&mut self, n: usize);
    fn get_slice(&self, start: usize, end: usize) -> &str;
}

pub struct TokenStreamString<'a> {
    source: &'a str,
    source_chars: Vec<(usize, char)>,
    pos: usize,
}

impl TokenStream for TokenStreamString<'_> {
    #[inline]
    fn has_next(&self) -> bool {
        self.pos < self.source_chars.len()
    }
    fn next(&mut self) -> Option<char> {
        if self.has_next() {
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

    fn seek_forward(&mut self, n: usize) {
        self.pos = std::cmp::min(self.pos + n, self.source_chars.len());
    }
    fn seek_back(&mut self, n: usize) {
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
