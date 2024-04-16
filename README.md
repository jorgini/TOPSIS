# Decison Maker web-application
### Author: Беликов Георгий

## About
Decision maker is program implementation of multi-criteria decision maker methods like TOPSIS and SMART.
MCDM methods is very useful toolkit for solving difficult problems with some alternatives and fixed criteria.

Decision maker supports user jwt authentication with refresh token for saving user history and solving problems in
group mode with others.

## Launch
Before launch application you need to check if docker is installed because Decision Maker supplied in multi-container
application Docker Compose.

For start web-application needed write below command in console:

```bash
make build && make run
```

If application launch for first time you need enter below command:

```bash
make build && make start && make migrate
```