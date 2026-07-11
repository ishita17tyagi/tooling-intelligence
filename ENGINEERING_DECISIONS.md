# Engineering Decisions

This document captures the reasoning behind the architecture, trade-offs and implementation choices made while building **Tooling Intelligence**.

It is intended to provide context beyond the source code and explain the engineering thought process.

---

# Why this project?

Large Language Models make it surprisingly easy to generate structured data, but building a reliable system around them is a very different engineering problem.

I built Tooling Intelligence to explore that gap.

The objective wasn't simply to ask an LLM questions about SaaS products—it was to design an end-to-end pipeline that could research applications, validate the generated information, normalize inconsistent outputs, quantify confidence and finally produce artifacts that are easy to review and share.

Throughout the project I deliberately prioritized:

- determinism over clever prompting
- explainability over opaque confidence
- modularity over tightly coupled logic
- reproducibility over one-off execution

Those principles shaped every architectural decision in this project.

---

# Architecture Philosophy

I deliberately separated the pipeline into independent stages.

```
Research

↓

Normalize

↓

Verify

↓

Analyze

↓

Generate Report
```

Each stage has a single responsibility.

This makes the pipeline easier to test, maintain and extend.

---

# Why Go?

At the beginning of the project I considered Python because of its mature AI ecosystem and rich collection of libraries for evaluation and orchestration.

However, after evaluating the actual requirements, I realized the core of this project wasn't machine learning—it was software engineering.

Most of the work involved:

- HTTP services
- JSON processing
- batch execution
- report generation
- deterministic validation
- file processing

Go provided a better fit for these requirements.

Its standard library covered almost everything the project needed while keeping the implementation lightweight and dependency-free.

Additional advantages included:

- single static binaries
- fast compilation
- simple deployment
- predictable concurrency model
- strong type safety

Choosing Go also aligned with the type of backend systems I enjoy building, making the project both productive and enjoyable to develop.

---

# Why deterministic verification?

One design decision I spent considerable time thinking about was verification.

A common approach is to introduce another LLM—or an evaluation framework such as RAGAS—to score the quality of generated responses.

While powerful, that approach introduces another probabilistic layer into the pipeline.

It also significantly increases:

- latency
- implementation complexity
- API usage
- token consumption
- operational cost

For a project whose outputs followed a well-defined schema, I concluded that deterministic verification offered a better engineering trade-off.

Instead of asking another model whether a response looked correct, the verifier checks objective properties such as:

- structural completeness
- evidence quality
- documentation domains
- required fields
- confidence deductions

This makes every confidence score explainable while keeping the pipeline inexpensive to execute.

Working within free-tier APIs reinforced this decision. Limited request quotas encouraged me to treat tokens as an engineering resource rather than something to consume freely, leading to a solution that is both simpler and more efficient.

---

# Why normalization?

LLMs are inherently inconsistent.

For example:

```
OAuth

OAuth2

OAuth 2.0

OAuth 2

Bearer OAuth
```

all describe essentially the same authentication mechanism.

Without normalization, analytics become noisy and misleading.

Adding a normalization stage significantly improved the quality of generated insights.

---

# Error Handling Philosophy

The batch pipeline never aborts because a single application fails.

Instead:

```
Slack

✓

GitHub

✓

Stripe

✓

Shopify

✗

NotebookLM

✓
```

Partial failures are captured and processing continues.

This mirrors production batch systems where resilience is more valuable than perfect execution.

---

# Testing Strategy

Development followed an incremental approach.

Each phase introduced one new capability and was verified before proceeding.

Examples include:

- Gemini connectivity
- JSON parsing
- normalization
- verification
- analytics
- report generation

Only after one phase became stable was the next introduced.

This reduced debugging complexity considerably.

---

# Challenges

The most interesting challenge wasn't implementing the pipeline—it was building it within practical resource constraints. 😅

During development I relied primarily on free-tier Gemini API quotas, which meant every request mattered.

Instead of repeatedly invoking the model throughout the pipeline, I redesigned the architecture so that a single research call produced all the information needed for downstream processing.

Everything after that point—normalization, verification, analytics and report generation—runs entirely offline.

This separation produced several benefits:

- fewer API calls
- lower token consumption
- deterministic post-processing
- faster iterations during development
- easier debugging

Ironically, the limitations of the free tier resulted in a better overall architecture by encouraging careful resource usage instead of relying on repeated LLM calls.

---

# # Future Directions

If I continue evolving Tooling Intelligence, these are the areas I'd explore next.

## MCP Integration

Research through Model Context Protocol servers instead of relying solely on prompt-based retrieval.

---

## Parallel Worker Pool

Concurrent research with bounded worker pools.

---

## Retry Queue

Automatic retry of transient API failures.

---

## Multi-provider Support

OpenAI

Claude

Gemini

Groq

Local models

---

## Interactive Dashboard

Replace the static HTML report with an interactive dashboard supporting:

- filtering
- searching
- sorting
- charts

---

## Deployment

Containerized deployment using Docker and Kubernetes.

---

## Persistent Storage

Replace generated JSON artifacts with a database-backed architecture.

---

# Lessons Learned

The biggest lesson from this project wasn't learning how to call Gemini.

It was understanding where LLMs fit within a larger software system.

When I started, I assumed the interesting part would be prompt engineering.

By the end of the project, I realized prompting represented only a small fraction of the engineering effort.

The majority of the work went into designing reliable systems around the model:

- validation
- normalization
- observability
- fault tolerance
- confidence scoring
- analytics
- reproducibility

That shift in perspective fundamentally changed how I think about AI applications.

Large Language Models are powerful components, but they should rarely become the architecture themselves.

Instead, they work best as one stage within a carefully engineered pipeline where deterministic software handles everything it can, and the model is only responsible for tasks that genuinely require reasoning.

This project reinforced an engineering principle I expect to carry into future work:

> Treat LLMs as collaborators, not as sources of truth.

Good AI systems are rarely built by asking better questions.

They're built by designing better systems around the answers.