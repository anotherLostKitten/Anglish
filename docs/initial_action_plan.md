## Anglish Initial Action Plan

### 2025-09-05

First step is working on getting run-able code on a single machine -- Samoza OS distributed runtime needs to be built atop this later.

While waiting on more concrete spec. for the AIR, focus should be on input into LLMs:

1. Parsing Anglish description into series of LLM prompts.

2. Figuring out what LLMs give us best results for Anglish, where we can host them (locally? where?). Also look for existing research interfacing with LLMs in similar ways.

3. What are we outputting from LLMs? Existing languages we compile into AIR? See what gives us best ouput & is easy to verify.

4. Design list of Anglish Tools the LLMs will have access to.

5. Extracting verifiable rules from Anglish contracts to check LLM output over.

6. Extracting & visualizing Space-Path Graph for verification & debugging.

#### Questions for meeting :

- How much of spatial features from Splash etc. are being carried forward into Samoza / Anglish.

- What access to LLMs do we have available through the lab? Subscriptions to online tools / machines for running locally etc..

- Are paths just protecting writes? Seems like reads are called into directly. If it's just being called directly from vibe lines, how do we verify this?

#### Tech stack :

Rust,

[LangChain](https://github.com/Abraxas-365/langchain-rust),
