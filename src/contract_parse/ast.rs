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
    pub fn new(name: *const str, binds: *const str, line: usize, col: usize) -> Param {
        Param {
            name,
            binds,
            mem_source: None,
            line,
            col,
        }
    }

    pub fn new_owned_mem(name: &str, binds: &str, line: usize, col: usize) -> Param {
        let mut p = Param {
            name,
            binds,
            mem_source: Some(format!("{}{}", name, binds)),
            line,
            col,
        };

        // it would be nice to paritially init these because it is an extra memcpy ... oh well ! maybe it will get optimized away
        p.name = &p.mem_source.as_ref().unwrap()[..name.len()];
        p.binds = &p.mem_source.as_ref().unwrap()[name.len()..];

        return p;
    }

    pub fn print(&self) {
        println!("{} = %{}", unsafe { &*self.name }, unsafe { &*self.binds });

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
    pub fn new(vibe_prose: String, line_start: usize, line_end: usize) -> VibeBlock {
        let mut block = VibeBlock {
            vibe_prose,
            meta_refs: vec![],
            line_start,
            line_end,
        };

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
