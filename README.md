# go-newsletter
A repo to collaborate on a school semestral project in GO.

## TODO List (Based on Work Division)

### ✅ Completed Tasks
1. **Project Foundation & Core Backend** - Basic Go project structure, HTTP server with Chi router, database connection setup
2. **Database Setup & Management** - Supabase project setup and PostgreSQL connection established

### 🚀 In Progress / Next Tasks

#### 3. Authentication & Authorization
- [ ] Integrate with Supabase Auth for Editor sign-up (email/password)
- [ ] Integrate with Supabase Auth for Editor sign-in and session management
- [ ] Implement JWT handling within the Go API (validation of Supabase-issued JWTs, middleware for protected routes)
- [ ] Implement password reset flow leveraging Supabase Auth
- [ ] Develop authorization logic to differentiate between standard Editors and Admin users
- [ ] Implement logic for populating and managing the `profiles` table upon user creation

#### 4. User (Editor & Admin) Management
- [ ] API endpoints for Editors to manage their own profiles
- [ ] Admin-specific API endpoints for user management (list all profiles, grant/revoke admin privileges)

#### 5. Newsletter Management
- [ ] Implement CRUD API endpoints for newsletters (Create, Read, Update, Delete)
- [ ] Database interactions for the `newsletters` table
- [ ] Admin-specific API endpoints for newsletter management

#### 6. Subscription Management
- [ ] API endpoint for public users to subscribe to a newsletter
- [ ] Email confirmation for new subscriptions
- [ ] API endpoint to verify confirmation token and activate subscription
- [ ] API endpoint for unsubscribing using a unique token
- [ ] API endpoint for Editors to list subscribers for their newsletters

#### 7. Publishing & Scheduling Workflow
- [ ] API endpoint for Editors to create/update post content
- [ ] Functionality for immediate publishing of a post
- [ ] Functionality for scheduled publishing of a post
- [ ] Implement Scheduler Service (background goroutine)
- [ ] API endpoints for Editors to manage scheduled posts

#### 8. Email Integration
- [ ] Research, select, and integrate external email service (Resend, SendGrid, AWS SES)
- [ ] Implement Go module/service for sending emails
- [ ] Develop email templates for various notifications

#### 9. Deployment & CI/CD
- [ ] Create comprehensive Dockerfile
- [ ] Set up CI/CD pipeline (GitHub Actions)
- [ ] Configure deployment to cloud platform
- [ ] Manage environment variables for different environments

#### 10. Testing
- [ ] Write unit tests for all services, handlers, and utility functions
- [ ] Develop integration tests for key API workflows
- [ ] Ensure tests are part of the CI/CD pipeline

## Links

### Documentation
[**Implementation Guide**](documentation/implementation-guide.md) - Complete guide for implementing new endpoints and understanding the architecture

[Architecture](https://github.com/najmamat/go-newsletter/blob/main/project-kickoff/architecture_1_supabase_explanation.md)

[Database design](https://github.com/najmamat/go-newsletter/blob/main/project-kickoff/database_schema.md)

[OpenAPI](api/openapi.yaml)

[Tasks and work division](https://github.com/najmamat/go-newsletter/blob/main/project-kickoff/work_division_supabase_monolith.md)

### Project Management
[Backlog](https://mattermost.hatrinh.cz/boards/team/94xsj4hkrpnbigytk8gudanzcr/bc8kuuqqrtp8txq4ksk5m778tww/vrxnoo8wrxjb59epes3btgptb5y)


