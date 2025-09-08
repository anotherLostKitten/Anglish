pub struct VibeBlock {
    vibe_prose: String,
    pub meta_refs: Vec<MetaRef>,
}

pub struct MetaRef {
    pub ident: *const str,
    ref_info: MetaRefData,
}

enum MetaRefData {
    DataRef,
    UseImportSpace(*const str),
    UseImportAgent(*const str),
    TaskRef(Vec<Param>),
    PathRef,
}

struct Param {
    name: *const str,
    binds: *const str,
}

pub struct SpaceDecl<'a> {
    parent: Option<&'a SpaceDecl<'a>>,
    ident: *const str,
    space_type: SpaceType,
    params: Vec<Param>,
    vibe_desc: VibeBlock,
    decl_text: String,
}

enum SpaceType {
    UI,
    IO,
    Data,
    Func,
    Agentic,
}

impl VibeBlock {
    pub fn new(source: &str) -> VibeBlock {
        let mut block = VibeBlock {
            vibe_prose: String::from(source),
            meta_refs: vec![],
        };
        block.meta_refs.push(MetaRef {
            ident: block.vibe_prose.as_str(),
            ref_info: MetaRefData::DataRef,
        });
        return block;
    }

    pub fn get_prose(&self) -> &str {
        self.vibe_prose.as_str()
    }
}

impl MetaRef {
    pub fn get_ident(&self) -> &str {
        unsafe {&*self.ident}
    }
}
