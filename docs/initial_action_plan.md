## Anglish Initial Action Plan

### 2025-09-05

First step is working on getting run-able code on a single machine -- Samoza OS distributed runtime needs to be built atop this later.

While waiting on more concrete spec. for the AIR, focus should be on input into LLMs:

1. Parsing Anglish description into series of LLM prompts.

2. Figuring out what LLMs give us best results for Anglish, where we can host them (locally? where?). Also look for existing research interfacing with LLMs in similar ways.

3. What are we outputting from LLMs? Existing languages we compile into AIR? See what gives us best ouput & is easy to verify.

4. Design list of Anglish Tools the LLMs will have access to.

5. Extracting verifiable rules from Anglish contracts to check LLM output over, tagged to vibe-blocks. -- First focus on just parsing into prompts before worrying about requirements.

n.b. these will be extracted both deterministically (eg. check that a =path is used in a vibe-block) & agenticly (eg. check =path uses correct vibe-policy during runtime)

6. Extracting & visualizing Space-Path Graph for verification & debugging.

#### Questions for meeting :

- How much of spatial features from Splash etc. are being carried forward into Samoza / Anglish.

In Anglish, scaling workers will happen only when explicitly specified by the user through exposed tasks native to the runtime. Offloads burden of managing runtime distribution to user rather than attempting to have system do so natively.

- What access to LLMs do we have available through the lab? Subscriptions to online tools / machines for running locally etc..

There is a citelab computer with a 3080 (I think it was); however using a home system with a 5060 Ti seems like a better approach for quickly iterating and for running larger models (16Gb VRAM).

- Are paths just protecting writes? Seems like reads are called into directly. If it's just being called directly from vibe lines, how do we verify this?

Spec should be changed with different path types to allow differentiating between "request & data read" vs. "write" (vs. other types? "control"? needs design work) paths.

#### Tech stack :

Rust,

[LangChain](https://github.com/Abraxas-365/langchain-rust),
