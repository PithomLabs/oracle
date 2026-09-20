AI Agents Failed to Prove Fermat's Last Theorem. Then They Got a Shared To-Do List
#
ai
#
agents
#
architecture
#
engineering
On September 4, Anthropic published something that sounds like a headline from a decade in the future: the first complete, computer-checked proof of Fermat's Last Theorem, written by a team of Claude agents working largely autonomously over 11 days. Thirteen million lines of Lean. Nearly 30,000 intermediate theorems. About six billion output tokens.

I want to talk about a detail that most coverage will bury, because it is the only part that matters if you build software with agents instead of reading about them. The first attempts failed. Not because the model was too weak. The agents had early success, then lost track of the project's state and stopped collaborating effectively. What fixed it was not a smarter model. It was a shared directed acyclic graph acting as the team's memory.

If you have ever run two AI agents on the same codebase and watched them trample each other's work, you already understand this failure. You just have not seen it dramatized at the scale of one of the hardest proofs in mathematics.

What actually happened, in numbers
First the facts, because they are dramatic enough on their own. Fermat scribbled his claim around 1637: no positive integers a, b, c satisfy aⁿ + bⁿ = cⁿ for any n greater than 2. Andrew Wiles proved it in 1995 after a 129-page proof, and even that is underselling the drama. He presented the proof in June 1993, a reviewer's question exposed a critical gap two months into verification, and Wiles spent a year, first alone and then with his former student Richard Taylor, fixing it.

Formalizing that proof, meaning rewriting it so a proof assistant like Lean can verify every step algorithmically, has been a community project since 2024, led by Kevin Buzzard at Imperial College London. The blueprint for just the initial phase runs 86 pages. It was scoped as a multi-year effort.

Then Tianyi Peng, an Anthropic researcher whose group at Columbia University builds AI formalization tools, tested whether Claude could make progress on it. The result, per Anthropic's post:

11 days of largely autonomous work produced the first end-to-end, computer-checked proof of FLT.
13 million lines of Lean, over 5x the size of Mathlib, the community proof library it builds on. Anthropic notes this is partly because Mathlib is concise and well-reviewed while their proof is longer than it needs to be.
30,300 theorems proved, 29,500 used in the final proof. Dozens of agents defined concepts, proved intermediate results, and chained them upward.
About 6 billion output tokens from an internal research model Anthropic describes as roughly comparable to Claude Fable 5.1.
The proof is real by Lean's standards. It uses only Lean's three standard axioms, contains no omitted steps, and a comparator tool confirmed the theorem statement matches Mathlib's own statement of FLT. Kevin Buzzard reviewed it.
The human role shrinks to something almost funny when you read the excerpts. Mathematical input was limited to occasional high-level nudges from Peng: "Jacobian as a scheme sounds high priority." "Push the Mazur theorem to be done soon." That is a manager setting priorities, not an engineer writing proofs.

The agents' own logs capture the finish line. At 02:00:57 UTC on August 18, one agent noted: "The FLT root reads PROVED on prove2me at 02:00:57Z Aug-18. Historic moment for this campaign."

The failure first, because that is the lesson
Here is the part I keep re-reading. From Anthropic's post:

A number of Claude's initial attempts failed: while agents had some early success, they quickly lost track of the project's state and stopped collaborating effectively. Their failed efforts contributed ~7% of the non-boilerplate lines in the final proof.

Parse that carefully. The agents were capable enough to make progress. The model did not change between the failed attempts and the successful one. What changed was the coordination layer. Early agents drifted because the project's state lived in their context windows, and context windows degrade. One agent's mental model of what was proven diverged from another's, and collaboration collapsed into noise.

The fix was Prove2Me, an open platform Peng and his Columbia collaborators designed, and it did three things:

Maintained a directed acyclic graph of theorem statements that agents consulted to decide what to prove next. This mitigated memory degradation and let many agents work in parallel without stepping on each other.
Separated theorem statements from their proofs into different files, with the links maintained independently. This sped up Lean compilation and cut resource consumption.
Kept a natural-language description of every node, which enabled search and reuse. Later agents found existing results instead of re-deriving them, producing a simpler proof path.
Notice what is on that list. Nothing about prompt engineering. Nothing about a bigger context window or a cleverer model. It is data modeling. State externalized into a structure every agent can read and trust.

Why this maps directly to your day job
I run a small fleet of scheduled AI agents on my own server: a trend scanner, a researcher, and writer pipelines that fire three times a day. My infrastructure is a rounding error next to Anthropic's, but I have hit exactly this failure mode, and maybe you have too.

Agents that hold shared state in their context are time bombs. My earliest pipeline runs would lose a lesson between the research step and the writing step unless it was written to a file the next step read. Once I moved coordination into persisted checklists and structured state, the same model stopped making the same mistakes. Anthropic's failed FLT attempts are that lesson at six billion tokens of scale. The failure was never intelligence. It was forgetting.

The DAG is the oldest idea in software, and it works. A directed acyclic graph of tasks with explicit dependencies is how Make, Bazel, and every CI pipeline have scheduled work for decades. Prove2Me's contribution is almost embarrassing in its simplicity: give agents a build system for theorems. Each node is a well-defined deliverable with a completion criterion, provable or not, so parallel agents have an unambiguous claim on work. If you are designing a multi-agent system today, your first artifact should be the dependency graph, not the agent prompts.

Verifiers are what let you trust agent output at scale. Nobody at Anthropic reviewed 13 million lines by hand. They did not have to. Lean's kernel checked every step, and a comparator matched the final statement against Mathlib's. That is the same trust model as a compiler and a test suite: you do not review the code, you review the gates. The more of your agent pipeline's output that flows through a mechanical gate, the more autonomy you can afford to hand over. Weak gates, weak agents, no matter how good the model is.

The human job becomes prioritization. "Jacobian as a scheme sounds high priority" is the entire management style that produced a historic proof. When the coordination layer is solid and the verifier is strict, the operator's leverage comes from choosing what matters next, not from supervising execution. That is a strange and useful reframe for anyone worried about what their job looks like in an agentic workflow.

Recovery value is real. Anthropic says the failed attempts still contributed about 7% of the non-boilerplate lines in the final proof. Failed agent runs are not pure waste if the work lands in shared, searchable state where later agents can salvage it. Another argument for externalized state: it turns even dead runs into partial assets.

The part that should concern the mathematics world less and the rest of us more
Buzzard, who has spent two years organizing humans to formalize this exact theorem, wrote after reviewing the proof:

If the automatic formalization of FLT is possible now, then we have taken a big step towards automatic formalization of the modern mathematical literature. Such autoformalization techniques will lead to new tools, rooting out errors in the current mathematical corpus and lightening the load of referees.

The verification history he is reacting to is grim. Thomas Hales' 1998 proof of the Kepler conjecture spent four years in review before a 12-referee panel settled on "99% certain," and Hales then led a twenty-person project to formalize it properly. Perelman's Poincare proof took roughly four years and three 300-page expositions to accept. Machine-checkable proofs end that class of agony.

And the barrier to entry is collapsing fast. Anthropic ran a side experiment where three consumer Claude Max subscription plans, collaborating through Prove2Me with no special access, formalized Vinogradov's Three Primes Theorem in three days. A result that previously required an academic collaboration was done by what is functionally three hobbyist accounts and a good to-do list.

For the rest of us, the transferable claim is this: a multi-agent system just delivered end-to-end on one of the hardest well-defined problems in existence, and the decisive variable was the collaboration architecture, not model capability. Anthropic's own framing agrees. The proof follows a simplified exposition by Darmon, Diamond, and Taylor rather than inventing new mathematics, which is exactly the shape of most software work: the destination is known, the path is long, and the bottleneck is coordination.

What I would copy into any multi-agent setup this week
A practical checklist, distilled from how Prove2Me actually worked:

Externalize the plan as a DAG. Every task gets a node, an explicit dependency list, and a machine-checkable done condition. Agents pull from the graph; nothing lives only in a prompt or a chat history.
Separate interface from implementation. Prove2Me split theorem statements from proofs to speed compilation. The software analogue: contracts, types, and interfaces in one place agents read cheaply, implementations behind them. Faster feedback loops for every worker.
Index everything in natural language. A short description of each node made past work searchable and reusable. If your agents cannot search what previous agents did, you are paying full price for every run.
Make the verifier the boss. Route agent output through the strictest mechanical gate available: compiler, type checker, test suite, linter, schema validation. Autonomy should be a function of gate strength.
Log state changes to shared storage, not to the conversation. Context is working memory at best. Anything that matters past this run goes to disk.
Spend your own attention on priorities, not supervision. If you are micromanaging agents, the coordination layer is broken. Fix the graph, not the agents.
None of this requires Anthropic's budget. The uncomfortable, excellent news of this result is that it mostly requires discipline about state that engineering teams have understood since the invention of the build tool.

I write about AI engineering, agents, and backend systems every week. Subscribe, it is free, and it tells me this kind of deep-dive is worth the evening it takes.

Have you run multi-agent workflows on real work? Did you hit the lost-state wall, and what fixed it for you? I am collecting war stories for a follow-up piece, so tell me in the comments.

Sources
Anthropic: Formalizing Fermat's Last Theorem
The proof on GitHub
Imperial College London FLT project and its blueprint
Prove2Me platform
Lean community: The Fermat's Last Theorem Project
