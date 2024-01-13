# JCAT
A WIP CLI to catogorize issues in a given Jira project

## Secrets
Export the following values in your terminal session before running jcat
```sh
export JIRA_API_KEY=<Your Jira API Key> \
export OPENAI_API_KEY=<Your OpenAI API Key>
```

## Usage
- Queries the IT Jira project for issues in the last 90 days and categorizes issues.
```sh
jcat
```
- Queries the IT Jira project for issues in the last 90 days and categorizes issues.
```sh
jcat -project IT
```
- Queries the IT Jira project key for issues in the last 10 days and categorizes issues.
```sh
jcat -project IT -days 10
```
- Queries with custom JQL and categorizes issues.
```sh
jcat -jql "project = IT AND createdDate >= 2024-01-08"
```