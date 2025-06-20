# Gator RSS Aggregator

## Prerequisites

- Will need Postgres v15 installed
- Will need Golang v1.23.0 or later installed

## Installation

- Clone the repo
- install using `go install gator`
- Create json file "~/.gatorconfig.json" with the following format:
{"db_url":"postgres://username:@localhost:5432/gator?sslmode=disable","current_user_name":""}

