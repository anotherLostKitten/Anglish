#[allow(unused_imports)]

mod contract_parse;

fn main() {
    // let bazinga = contract_parse::ast::VibeBlock::new(">>> test test 123");
    // println!("{}", bazinga.meta_refs[0].get_ident());

    let _ = contract_parse::parser::parse_from_str("
@space_mode_THE_MOST_999:AGENTIC ( in = %test1 ,out=%test2)
>>> i am vibing
>>> vibe coding!
>>>vibes...
>>> (0=0)/ <-- he is waving !
");
}
