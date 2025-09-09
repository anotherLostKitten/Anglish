use std::mem::{MaybeUninit, transmute};

pub struct Contract {
    pub spaces: Vec<SpaceDecl>,
    pub agents: Vec<AgentDecl>,
    pub paths: Vec<PathDecl>,
}

impl Contract {
    pub fn new() -> Contract {
        Contract {
            spaces: vec![],
            agents: vec![],
            paths: vec![],
        }
    }
}

pub struct SpaceDecl {
    // parent: Option<&'a SpaceDecl<'a>>, // todo these lifetimes are weird again
    pub ident: String,
    pub space_type: SpaceType,
    pub params: Vec<Param>,
    pub vibe_desc: VibeBlock,

    // todo space content -- tasks etc.

    pub line_start: usize,
    pub line_end: usize,
}

pub enum SpaceType {
    UI,
    IO,
    DATA,
    FUNC,
    AGENTIC,
    Unknown,
}

pub struct AgentDecl {
    agent_type: AgentType,
    // todo!
}

pub enum AgentType {
    AF,
    DF,
}

pub struct PathDecl {
    path_type: PathType,
    // todo!
}

pub enum PathType {
    READ,
    WRITE,
}

pub struct TaskDecl {
    // todo!
}

pub struct DataDecl {
    //todo!
}

pub struct Param {
    name: *const str,
    binds: *const str,

    mem_source: Option<String>, // in decls, params own their own string memory; in vibe blocks, it is in the vibe prose.

    pub line: usize,
    pub col: usize,
}

impl Param {
    pub fn new(name: &str, binds: &str, line: usize, col: usize) -> Param {
        Param {
            name: name,
            binds: binds,
            mem_source: None,
            line: line,
            col: col,
        }
    }

    pub fn new_owned_mem(name: &str, binds: &str, line: usize, col: usize) -> Param {
        let mut p = Param {
            name: name,
            binds: binds,
            mem_source: Some(format!("{}{}", name, binds)),
            line: line,
            col: col,
        };

        // it would be nice to paritially init these because it is an extra memcpy ... oh well ! maybe it will get optimized away
        p.name = &p.mem_source.as_ref().unwrap()[..name.len()];
        p.binds = &p.mem_source.as_ref().unwrap()[name.len()..];

        return p;
    }
}

pub struct VibeBlock {
    vibe_prose: String,
    pub meta_refs: Vec<MetaRef>,

    line_start: usize,
    line_end: usize,
}

pub struct MetaRef {
    ident: *const str,
    ref_info: MetaRefData,

    line: usize,
    col: usize,
}

enum MetaRefData {
    DataRef,
    UseImportSpace(*const str),
    UseImportAgent(*const str),
    TaskRef(Vec<Param>),
    PathRef,
}

impl VibeBlock {
    // pub fn new(source: &str) -> VibeBlock {
    //     let mut block = VibeBlock {
    //         vibe_prose: String::from(source),
    //         meta_refs: vec![],
    //     };
    //     block.meta_refs.push(MetaRef {
    //         ident: block.vibe_prose.as_str(),
    //         ref_info: MetaRefData::DataRef,
    //     });
    //     return block;
    // }

    pub fn get_prose(&self) -> &str {
        self.vibe_prose.as_str()
    }
}

impl MetaRef {
    pub fn get_ident(&self) -> &str {
        unsafe {&*self.ident}
    }
}
