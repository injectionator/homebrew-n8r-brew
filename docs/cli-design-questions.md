# n8r CLI — Open Design Questions

These questions need to be answered to inform the next phase of CLI development.

## What is an "injectionator"?

Before we can build commands like `n8r inspect`, we need to define what the core object is that a user owns or interacts with. Possibilities include:

- A security challenge or lab environment they've configured?
- A runtime or agent they've deployed?
- A set of test results or scan reports?
- A project containing prompt injection test cases?
- A configuration profile for their security posture?

## What should `n8r inspect` show?

Once we know what the core object is, `n8r inspect` should surface its key properties. For example:

- **If it's a runtime/agent:** status (running/stopped), version, endpoint URL, last activity
- **If it's a project:** name, number of test cases, last run date, pass/fail summary
- **If it's a lab environment:** active challenges, completion status, score
- **If it's a config profile:** rules, policies, detection settings

## What other commands should exist?

Depending on the core object, the CLI might support:

- `n8r inspect` — view your injectionator's current state
- `n8r run` — execute a test or challenge
- `n8r list` — list available challenges/tests/projects
- `n8r logs` — view recent activity
- `n8r config` — manage settings

## What API endpoints are needed?

Each CLI command needs a corresponding API endpoint on injectionator.com. These need to be designed alongside the CLI commands. All endpoints should:

- Accept Bearer token authentication (using the token from `n8r login`)
- Return JSON responses
- Be documented for the CLI to consume

## Who are the users?

- Alpha testers only (gated by Clerk account + cohort membership)?
- What level of access do they have?
- Are there different tiers of access?
