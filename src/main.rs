#[allow(unused_imports)]

mod contract_parse;

fn main() {
    let bazinga = contract_parse::ast::VibeBlock::new(">>> test test 123");
    println!("{}", bazinga.meta_refs[0].get_ident());
    contract_parse::parser::parse_from_str("ident_1 _idEEent2 ident3");
}
