# JCAT
A WIP CLI to catogorize issues in a given Jira project

## Usage
- Queries the IT Jira project for issues in the last 90 days and categorizes issues.
```sh
jcat IT
```
- Queries the IT Jira project key for issues in the last 10 days and categorizes issues.
```sh
jcat IT -days 10
```
- Queries with custom JQL and categorizes issues.
```sh
jcat -jql "project = IT AND createdDate >= 2024-01-08"
```