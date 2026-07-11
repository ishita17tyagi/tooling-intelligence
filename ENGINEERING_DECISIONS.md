# Engineering Decisions

This document captures the reasoning behind the architecture, trade-offs and implementation choices made while building **Tooling Intelligence**.

It is intended to provide context beyond the source code and explain the engineering thought process.

---

# Why this project?

The assignment required researching SaaS applications and producing structured integration metadata.

Instead of building a simple prompt wrapper around an LLM, I wanted to build something that resembles a production engineering pipeline.

The focus became:

- deterministic outputs
- reproducibility
- robustness
- explainability
- modularity

rather than only generating JSON from prompts.

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

I initially explored implementing the project in Python.

However, I found myself spending more time debugging environments and dependency issues than solving the engineering problem.

Since the project mostly consisted of:

- HTTP APIs
- JSON processing
- File generation
- Batch execution

Go was a better fit.

Advantages:

- single static binary
- fast compilation
- minimal dependency management
- excellent standard library
- straightforward concurrency model

---

# Why deterministic verification?

Many similar projects rely entirely on another LLM or evaluation library to judge LLM outputs.

Instead, I chose deterministic verification.

Each record is scored using:

- structural validation
- evidence validation
- completeness validation

This produces explainable confidence scores and avoids introducing another probabilistic layer into the system.

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

The largest practical challenge was operating within the limits of free-tier APIs.

Google Gemini imposes request quotas and temporary availability limits.

Rather than treating these constraints as blockers, they influenced the design.

The implementation minimizes unnecessary API calls by:

- generating structured output in one request
- separating online research from offline analytics
- generating reports entirely from saved artifacts
- avoiding repeated LLM verification

This reduced token usage while making the system faster and easier to reproduce.

---

# What I'd improve next

If this project evolved beyond the assignment, I would add:

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

The biggest takeaway from this project wasn't using Gemini.

It was learning how to engineer reliable systems around LLMs.

Prompting is only a small part of an AI application.

The majority of engineering effort goes into:

- validation
- normalization
- observability
- fault tolerance
- analytics
- reproducibility

Designing these supporting systems proved to be significantly more valuable than simply obtaining an LLM response.

That perspective shaped the final architecture of Tooling Intelligence.