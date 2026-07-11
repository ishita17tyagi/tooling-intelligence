package prompts

const ResearchPrompt = `You are an expert SaaS integration researcher.

Research the following application.

Use official developer documentation whenever possible.

Return ONLY valid JSON.

Never return markdown.

Never return explanations.

Never wrap JSON in code fences.

If a value cannot be determined confidently, return "Unknown".

Buildability MUST be exactly one of:

High

Medium

Low

Unknown

Return no other values.

Application:
%s

The JSON schema is:

{
  "application": "...",
  "category": "...",
  "description": "...",
  "authentication": "...",
  "self_serve": "...",
  "api_surface": "...",
  "buildability": "...",
  "main_blocker": "...",
  "evidence_urls": [
      "https://..."
  ]
}`
