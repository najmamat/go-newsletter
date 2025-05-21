# GO Newsletter API – Requirements (EN)

#### Architecture

- **REST / GraphQL** – The API must follow REST or GraphQL architecture.
- **API service** – Must serve both web and mobile clients.

#### Database

- **Editors DB** – Store editors in PostgreSQL.
- **PostgreSQL** – Primary general-purpose database.
- **Subscribers DB** – Store newsletter subscribers in Firebase.
- **No ORM** – ORMs are explicitly forbidden.
- **Newsletter DB** – Newsletters can be stored in either PostgreSQL or Firebase.

#### Deployment

- **Cloud Deployment** – App must be deployed to any cloud platform (e.g., Railway, Render, Heroku, GCP, AWS).

#### Documentation

- **API Documentation** – Must provide detailed documentation for client use and maintenance.

#### Naming Conventions

- **Firebase Naming** – Use format `strv-vse-go-newsletter-[last_name]-[first_name]`.
- **Cloud Naming** – Same convention as Firebase.

#### Newsletter Management

- **Newsletter CRUD Auth** – Only authenticated and authorized editors can manage newsletters.
- **Create Newsletter** – Editors can create newsletters (title required, description optional).
- **Update Newsletter** – Editors can update newsletter attributes (e.g., rename).
- **Delete Newsletter** – Editors can delete newsletters (consider what happens to links, subscribers, archives).

#### Publishing

- **Publishing Auth** – Only authenticated and authorized editors can publish.
- **Publish Posts** – Editors can send newsletter posts using external services (e.g., Resend, SendGrid, AWS SES).
- **Message Archiving** – Published posts must be stored in DB.
- **Schedule Posts** – Editors can schedule posts for future publication at a specific date and time.
- **Scheduled Post Processing** – The system must reliably process and send scheduled posts at their designated time.
- **View/Manage Scheduled Posts (Optional but Recommended)** – Editors should be able to view, modify, or cancel their scheduled (not yet published) posts.

#### Subscription Management

- **Subscribe via Link** – Subscribers register using a unique newsletter-specific link.
- **Email Confirmation** – Subscribers receive a confirmation email after subscribing.
- **Unsubscribe Link** – Every email must include a working unsubscribe link.

#### Transaction Handling

- **Transactional Context** – Operations must use DB transactions where needed for consistency and robustness.

#### User Management

- **Editor Sign Up – Password** – Registration via email and password.
- **Editor Sign Up – OAuth** – (Optional) Social login using OAuth.
- **Editor Sign In** – Email/password login using JWT-based stateless authentication.
- **Password Reset** – Editors can request password resets.
- **JWT Authorization** – Stateless authorization using JWT is required.
- **Auth Provider** – (Optional) Firebase may be used as the auth provider.