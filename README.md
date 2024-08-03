# GitHub API Data Fetching and Service

## Overview
This service, written in Golang, interacts with GitHub's public APIs to fetch repository information and commits. It stores the data in a PostgreSQL database and continuously monitors the repositories for changes.

## Technologies Used
- **Golang**: Programming language.
- **PostgreSQL**: An open-source object-relational database system.
- **Golang-Migrate**: An open-source migration tool.
- **Gorilla-Mux**: An open-source routing tool.
- **Docker**: A containerization platform for seamless application building and sharing.

## Installation

### Docker (Preferred)
1. Clone the repository.
2. Create a `.env` file by copying from `.env.example`.
3. Run the following commands:
   ```bash
   make build      # Build the application
   make migrate_up # Apply migrations and create necessary tables
   make run          # Start the application
   ```
4. Access the application at `http://localhost:8080`.

### Local Installation
1. **Golang**: Install Golang. 
2. **PostgreSQL**: Install PostgreSQL. 
3. **Install Dependencies**:
   - Install all dependencies:
     ```bash
     go mod download
     ```
   - Install golang-migrate:
     ```bash
     go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
     ```
   - Create and configure a `.env` file using `.env.example`.
   - Export the below variables.
      ```bash
      export POSTGRESQL_URL="your postresql url in env"
      export ENVIRONMENT=dev
      ```
   - Run migrations:
     ```bash
     migrate -database ${POSTGRESQL_URL} -path db/migrations up
     ```
   - Start the server:
     ```bash
     go run main.go
     ```
4. Access the application at `http://localhost:8080`.

## API Endpoints

### Add Repository

To start fetching commits for a repository, make a `POST` request to the following endpoint:

- **URL:** `http://localhost:8080/api/v1/github/setup/`
- **Body:**
  ```json
  {
    "repo": "owner_name/repo_name",
    "from_date": "YYYY-MM-DDTHH:MM:SSZ",
    "to_date": "YYYY-MM-DDTHH:MM:SSZ"
  }
  ```
  - `repo` (required): Full repository name in the format `owner_name/repo_name`.
  - `from_date` (optional): Start date for fetching commits in ISO 8601 format.
  - `to_date` (optional): End date for fetching commits in ISO 8601 format.

### Get Top N Commit Authors

To retrieve the top N commit authors for a repository, make a `GET` request to:

- **URL:** `http://localhost:8080/api/v1/github/repo/{owner_name}/{repo_name}/top/{n}/commit-authors/`
  - `{owner_name}`: Repository owner name.
  - `{repo_name}`: Repository name.
  - `{n}`: Number of top authors to retrieve.

### Get Repository Commits

To get the commits for a repository, make a `GET` request to:

- **URL:** `http://localhost:8080/api/v1/github/repo/{owner_name}/{repo_name}/commits/`
  - `{owner_name}`: Repository owner name.
  - `{repo_name}`: Repository name.


## Testing

Automated tests are available and can be executed with:

- **Docker**:
  ```bash
  make test
  ```
- **Local**:
  ```bash
  go test ./tests/v1 
  ```

## Discussion

### Approach and Assumptions
  - Designed the schema with three key tables: `SetupData`, `Commits`, and `Repositories`.
  - Implemented unique constraints on the `repo` column in both `SetupData` and `Repositories` tables, and on the `sha` column in the `Commits` table to prevent duplicates.
  - Optimized query performance by indexing the `date`, `sha`, and `repo` columns in the `Commits` table, and the `repo` column in the `Repositories` table.
  - Utilized the `SetupData` table to manage repositories under monitoring. Upon startup, the application is pre-configured to fetch commits from the Chromium repository by seeding this table.
  -Before fetching commits, repository existence is checked in the `Repositories` table. If absent, its metadata is retrieved from GitHub and stored in the `Repositories` table.
  - To prevent duplicate commit retrieval, the `sha` column in the `Commits` table is uniquely constrained. The latest commit date for each repository is queried and used as the `since` parameter in GitHub API requests, ensuring that only new commits are fetched.
  - Set up a cron job to run hourly for monitoring purposes.

### Future Improvements
  - **Enhanced Test Coverage**: Develop more comprehensive test cases to ensure robustness.
  - **Improved Error Handling**: Implement more granular and informative error handling to aid troubleshooting.
  - **Code Structure Optimization**: Refine the code structure for better maintainability and readability.
  - **Asynchronous Commit Fetching**: Improve efficiency by fetching commits for multiple repositories concurrently, rather than sequentially processing one repository at a time.

