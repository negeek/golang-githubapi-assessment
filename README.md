# GitHub API Data Fetching and Service

This service, written in Golang, interacts with GitHub's public APIs to fetch repository information and commits. It stores the data in a PostgreSQL database and continuously monitors the repositories for changes.

## Project Setup

1. **Install Dependencies:**
   ```sh
   make install
   ```

2. **Set Environment:**
   ```sh
   export ENVIRONMENT=dev
   ```

3. **Set Database URL:**
   ```sh
   export POSTGRESQL_URL="your_postgres_connection_string"
   ```

4. **Run Database Migrations:**
   ```sh
   make migrate_up
   ```

5. **Start the Service:**
   ```sh
   make run
   ```

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

- **URL:** `http://localhost:8080/api/v1/github/repo/{repo}/top/{n}/commit-authors/`
  - `{repo}`: Repository name.
  - `{n}`: Number of top authors to retrieve.

### Get Repository Commits

To get the commits for a repository, make a `GET` request to:

- **URL:** `http://localhost:8080/api/v1/github/repo/{repo}/commits/`
  - `{repo}`: Repository name.

## Project Approach

1. **Database Schema:**
   - Created three tables: `SetupData`, `Commits`, and `Repositories`.
   - Ensured unique constraints on the `repo` column in `SetupData` and `Repositories`, and on the `sha` column in `Commits`.
   - Indexed `date`, `sha`, and `repo` columns in the `Commits` table for efficient querying.
   - Indexed the `repo` column in the `Repositories` table for efficient querying.

2. **SetupData Table:**
   - Stores repositories to monitor. The application fetches commits from repositories listed here.

3. **Default Repository:**
   - By default, the service fetches commits from the Chromium repository. This is set by seeding the `SetupData` table upon application startup.

4. **Add Repository Endpoint:**
   - Allows the addition of repositories for monitoring via the `Add Repository` endpoint.

5. **Cron Job:**
   - The cron job runs hourly by default.
   - Can be instructed to run specific jobs immediately upon application startup.

6. **Repository Validation:**
   - Before fetching commits, the service checks if the repository exists in the `Repositories` table.
   - If not, it fetches repository metadata from GitHub and stores it in the `Repositories` table.

7. **Commit Fetching:**
   - Before making a request to GitHub, the service checks for the latest commit in the `Commits` table.
   - Uses this information to set the `since` query parameter and fetch only new commits.

8. **Efficient Data Fetching:**
   - Prevents fetching all commits repeatedly and handles unique constraint violations by fetching only new commits.