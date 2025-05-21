# Work Division: Go Newsletter Platform (Monolith with Supabase)

**Team:** Matouš (Dev Lead), Adam, Jakub, Honza (Programmers)

This document outlines a proposed division of tasks for the development of the Go Newsletter Platform, based on a monolithic architecture using Go and Supabase. Matouš, as Dev Lead, will oversee the technical direction, guide the team, and contribute to development, particularly on foundational and complex tasks. Other programmers will take on tasks from the list below based on team discussion, priority, and expertise.

## Guiding Principles for Task Division:

*   **Clear Ownership for Tasks**: While tasks are listed, individuals or pairs can take clear ownership of specific tasks for a duration.
*   **Balanced Workload**: Aim to distribute tasks to ensure progress across all areas.
*   **Interdependencies & Collaboration**: Highlight tasks requiring close collaboration.
*   **Iterative Development**: Tasks can be broken down further into smaller user stories or sub-tasks for sprints/iterations.

## Task List

### 1. Project Foundation & Core Backend
*   **Responsibility**: Primarily Dev Lead (Matouš), with support/contribution from Programmers.
*   **Tasks**:
    *   Initialize Go project structure (directories, `go.mod`).
    *   Set up basic HTTP server and routing (e.g., using `net/http` or a lightweight router like `gorilla/mux` or `chi`).
    *   Implement robust configuration management (env vars, Supabase keys, external service keys).
    *   Establish a logging framework and conventions.
    *   Define common utilities, helper functions, and shared internal libraries/packages.
    *   Lead API design (REST with OpenAPI or GraphQL, based on team decision) and ensure consistency.
    *   Implement the main API request handling/middleware layer (e.g., request parsing, response formatting, error handling).
    *   Oversee and contribute to the creation and maintenance of API documentation (e.g., OpenAPI YAML).

### 2. Database Setup & Management
*   **Responsibility**: Collaborative (All team members, coordinated by Dev Lead).
*   **Tasks**:
    *   Set up the Supabase project.
    *   Define and implement the PostgreSQL schema using SQL DDL for all tables (`profiles` including `is_admin`, `newsletters`, `subscribers`, `published_posts` including `scheduled_at`, `status`).
    *   Establish a strategy for database schema migrations (if not solely relying on Supabase UI for initial setup).
    *   Ensure appropriate indexes are created for performance.

### 3. Authentication & Authorization
*   **Responsibility**: Programmer(s), guided by Dev Lead.
*   **Tasks**:
    *   Integrate with Supabase Auth for Editor sign-up (email/password).
    *   Integrate with Supabase Auth for Editor sign-in and session management.
    *   Implement JWT handling within the Go API (validation of Supabase-issued JWTs, middleware for protected routes).
    *   Implement password reset flow leveraging Supabase Auth.
    *   Develop authorization logic to differentiate between standard Editors and Admin users based on the `is_admin` flag in the `profiles` table.
    *   Implement logic for populating and managing the `profiles` table upon user creation in Supabase Auth.

### 4. User (Editor & Admin) Management
*   **Responsibility**: Programmer(s), guided by Dev Lead.
*   **Tasks**:
    *   API endpoints for Editors to manage their own profiles (e.g., update `full_name`, `avatar_url`).
    *   Admin-specific API endpoints for user management:
        *   List all user profiles.
        *   Grant/Revoke admin privileges for a user (updating `is_admin` in `profiles`).

### 5. Newsletter Management
*   **Responsibility**: Programmer(s), guided by Dev Lead.
*   **Tasks**:
    *   Implement CRUD API endpoints for newsletters (Create, Read, Update, Delete) for authenticated Editors (as owners).
    *   Database interactions for the `newsletters` table (raw SQL).
    *   Admin-specific API endpoints for newsletter management:
        *   List all newsletters in the system.
        *   Delete any newsletter.

### 6. Subscription Management
*   **Responsibility**: Programmer(s), guided by Dev Lead.
*   **Tasks**:
    *   Implement API endpoint for public users to subscribe to a newsletter (given newsletter ID and email).
    *   Implement email confirmation for new subscriptions (generate unique token, store it, send confirmation email).
    *   Implement API endpoint to verify confirmation token and activate subscription.
    *   Implement API endpoint for unsubscribing using a unique token (included in every email).
    *   Database interactions for the `subscribers` table (raw SQL).
    *   Implement API endpoint for Editors to list subscribers for their own newsletters.
    *   Ensure unsubscribe links are secure and correctly processed.

### 7. Publishing & Scheduling Workflow
*   **Responsibility**: Programmer(s), guided by Dev Lead.
*   **Tasks**:
    *   API endpoint for Editors to create/update post content (`title`, `content_html`, `content_text`).
    *   Functionality for immediate publishing of a post:
        *   Store post in `published_posts` table with `status = 'published'`, `published_at = NOW()`.
        *   Trigger email sending to subscribers.
    *   Functionality for scheduled publishing of a post:
        *   Allow `scheduled_at` (future time) in post creation/update API.
        *   Store post in `published_posts` table with `status = 'scheduled'`, and the specified `scheduled_at`.
    *   Implement Scheduler Service (e.g., a background goroutine):
        *   Periodically query `published_posts` for items where `status = 'scheduled'` and `scheduled_at <= NOW()`.
        *   For due posts, update `status` to `'publishing'` (to prevent reprocessing).
        *   Trigger the actual publishing process (retrieve subscribers, send emails via Email Service).
        *   Update `status` to `'published'` and set `published_at` upon success, or to `'failed'` on error.
    *   API endpoints for Editors to manage their scheduled (not yet published) posts:
        *   List scheduled posts for a newsletter.
        *   Get details of a specific scheduled post.
        *   Update a scheduled post (content or `scheduled_at` time).
        *   Cancel (delete) a scheduled post.
    *   Database interactions for `published_posts` table (raw SQL).

### 8. Email Integration
*   **Responsibility**: Programmer(s), guided by Dev Lead.
*   **Tasks**:
    *   Research, select, and integrate an external email service (e.g., Resend, SendGrid, AWS SES).
    *   Implement a Go module/service for sending emails (abstracting the chosen provider).
    *   Develop email templates for: subscription confirmation, password resets, published newsletter posts, and any other notifications.

### 9. Deployment & CI/CD
*   **Responsibility**: Programmer(s) (potentially Honza taking a lead here, as previously suggested), guided by Dev Lead.
*   **Tasks**:
    *   Create a Dockerfile for the Go application.
    *   Set up a CI/CD pipeline (e.g., using GitHub Actions) for automated builds, tests, and deployments.
    *   Configure deployment to the chosen cloud platform (e.g., Railway, Render, Heroku).
    *   Manage environment variables and application configuration for different environments (dev, staging, prod).

### 10. Testing
*   **Responsibility**: All Programmers, overseen by Dev Lead.
*   **Tasks**:
    *   Write unit tests for all services, handlers, and utility functions.
    *   Develop integration tests for key API workflows (e.g., auth, newsletter CRUD, subscription flow, publishing, scheduling).
    *   Ensure tests are part of the CI/CD pipeline.

## Cross-Cutting Concerns (All Team Members):

*   **Database Schema Design & Evolution**: While initial setup is a task, ongoing schema changes must be discussed and managed by the team.
*   **API Contract Adherence**: All backend work must align with the defined OpenAPI/GraphQL specification. The Dev Lead ensures this.
*   **Code Quality & Reviews**: Implement a mandatory code review process. Each PR should be reviewed by at least one other team member, ideally the Dev Lead for critical parts.
*   **Consistent Error Handling & Logging**: Adhere to established conventions throughout the application.
*   **Security Best Practices**: Apply security considerations in all relevant areas (input validation, SQL injection prevention, XSS protection if applicable, secure token handling, etc.).
*   **Documentation**: Maintain and update API documentation, internal code comments, and any changes to design documents.

## Initial Steps & Collaboration Points:

1.  **Team Kick-off & Plan Refinement**: Discuss this task list. Make adjustments based on team consensus.
2.  **Technology & Tooling Choices**: Finalize choices for router, database interaction libraries (if any, beyond `database/sql`), email service, etc.
3.  **API Design Finalization (Dev Lead)**: Solidify the initial OpenAPI/GraphQL schema. This is critical before extensive parallel development.
4.  **Supabase Project Setup (Team)**: Create the project, share access/keys securely.
5.  **Core Project Structure Setup (Dev Lead)**: Initialize the Go project, basic server, config, logging.
6.  **Task Prioritization & Assignment**: Break down initial high-priority tasks (e.g., auth, basic newsletter CRUD) and let programmers pick them up.
7.  **Regular Sync-ups**: Implement daily or frequent short stand-ups to discuss progress, blockers, and plan next steps.

This task-oriented division should provide flexibility while ensuring all aspects of the project are covered. Effective communication and collaboration will be key. 