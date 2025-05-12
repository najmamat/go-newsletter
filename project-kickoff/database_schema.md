# Go Newsletter API - Database Schema (Supabase/PostgreSQL)

This document outlines the proposed database schema for the Go Newsletter platform, hosted on Supabase (PostgreSQL).

## Table Overview

1.  **`profiles`**: Stores additional information for editors (including admin status), linked to Supabase's `auth.users` table.
2.  **`newsletters`**: Contains details about each newsletter.
3.  **`subscribers`**: Manages email subscriptions to newsletters.
4.  **`published_posts`**: Archives posts published to newsletters.

![Editor | Mermaid Chart-2025-05-12-163601](https://github.com/user-attachments/assets/026bc0be-6de5-4587-9edc-750238760394)

---

## Table Definitions

### 1. `profiles`

Stores public profile information for users (editors). This table is linked to the `auth.users` table provided by Supabase.

| Column Name      | Data Type     | Constraints                                                | Description                                     |
| ---------------- | ------------- | ---------------------------------------------------------- | ----------------------------------------------- |
| `id`             | `UUID`        | Primary Key, Foreign Key references `auth.users.id` ON DELETE CASCADE | User ID from Supabase Auth.                     |
| `full_name`      | `TEXT`        | Nullable                                                   | Full name of the editor.                        |
| `avatar_url`     | `TEXT`        | Nullable                                                   | URL to the editor's avatar image.              |
| `is_admin`       | `BOOLEAN`     | Not Null, Default `false`                                  | Flag to indicate if the user is an administrator. |
| `created_at`     | `TIMESTAMPTZ` | Not Null, Default `now()`                                  | Timestamp of profile creation.                  |
| `updated_at`     | `TIMESTAMPTZ` | Not Null, Default `now()`                                  | Timestamp of last profile update.               |

**Notes:**
*   Sensitive information like email and password hashes are managed by Supabase in the `auth.users` table.
*   You would typically create this table and then set up a trigger or use Row Level Security in Supabase to automatically create a profile entry when a new user signs up in `auth.users`.

---

### 2. `newsletters`

Stores information about each newsletter created by an editor.

| Column Name   | Data Type     | Constraints                                                               | Description                                     |
| ------------- | ------------- | ------------------------------------------------------------------------- | ----------------------------------------------- |
| `id`          | `UUID`        | Primary Key, Default `gen_random_uuid()`                                  | Unique identifier for the newsletter.           |
| `editor_id`   | `UUID`        | Not Null, Foreign Key references `auth.users.id` ON DELETE CASCADE          | ID of the editor who owns this newsletter.      |
| `name`        | `TEXT`        | Not Null                                                                  | Name of the newsletter.                         |
| `description` | `TEXT`        | Nullable                                                                  | Optional description of the newsletter.         |
| `created_at`  | `TIMESTAMPTZ` | Not Null, Default `now()`                                                 | Timestamp of newsletter creation.               |
| `updated_at`  | `TIMESTAMPTZ` | Not Null, Default `now()`                                                 | Timestamp of last newsletter update.            |

---

### 3. `subscribers`

Manages subscriptions to newsletters.

| Column Name         | Data Type     | Constraints                                                                    | Description                                                       |
| ------------------- | ------------- | ------------------------------------------------------------------------------ | ----------------------------------------------------------------- |
| `id`                | `UUID`        | Primary Key, Default `gen_random_uuid()`                                       | Unique identifier for the subscription record.                    |
| `newsletter_id`     | `UUID`        | Not Null, Foreign Key references `newsletters.id` ON DELETE CASCADE              | ID of the newsletter being subscribed to.                         |
| `email`             | `TEXT`        | Not Null                                                                       | Email address of the subscriber.                                  |
| `subscribed_at`     | `TIMESTAMPTZ` | Not Null, Default `now()`                                                      | Timestamp when the subscription was made.                         |
| `unsubscribed_at`   | `TIMESTAMPTZ` | Nullable                                                                       | Timestamp when the user unsubscribed, if applicable.              |
| `unsubscribe_token` | `TEXT`        | Not Null, Unique                                                               | Unique token for one-click unsubscription.                        |
| `is_confirmed`      | `BOOLEAN`     | Not Null, Default `false`                                                      | Flag to indicate if the subscription has been confirmed via email.|
| `confirmation_token`| `TEXT`        | Nullable, Unique                                                               | Unique token for email confirmation.                             |

**Constraints:**
*   `UNIQUE (newsletter_id, email)`: Ensures an email can only subscribe once to a specific newsletter.

**Notes:**
*   The `unsubscribe_token` should be generated upon subscription and included in every email.
*   The `confirmation_token` is used for the double opt-in process if implemented.

---

### 4. `published_posts`

Stores content of posts that have been published or are scheduled for publication to newsletters.

| Column Name     | Data Type     | Constraints                                                              | Description                                     |
| --------------- | ------------- | ------------------------------------------------------------------------ | ----------------------------------------------- |
| `id`            | `UUID`        | Primary Key, Default `gen_random_uuid()`                                 | Unique identifier for the published post.       |
| `newsletter_id` | `UUID`        | Not Null, Foreign Key references `newsletters.id` ON DELETE CASCADE        | ID of the newsletter this post belongs to.      |
| `editor_id`     | `UUID`        | Not Null, Foreign Key references `auth.users.id`                           | ID of the editor who published this post.       |
| `title`         | `TEXT`        | Not Null                                                                 | Title of the post.                              |
| `content_html`  | `TEXT`        | Not Null                                                                 | HTML content of the post.                       |
| `content_text`  | `TEXT`        | Nullable                                                                 | Plain text version of the post content.         |
| `status`        | `TEXT`        | Not Null, Default `'draft'`                                              | Status of the post (e.g., 'draft', 'scheduled', 'publishing', 'published', 'failed'). |
| `scheduled_at`  | `TIMESTAMPTZ` | Nullable                                                                 | If status is 'scheduled', this is the time it will be published (UTC). |
| `published_at`  | `TIMESTAMPTZ` | Nullable                                                                 | Timestamp when the post was actually published. Becomes non-null when status is 'published'. |
| `created_at`    | `TIMESTAMPTZ` | Not Null, Default `now()`                                                | Timestamp of post record creation.              |

**Notes:**
*   `editor_id` is included to track who specifically published the post, which might be useful if multiple editors could potentially manage one newsletter in the future (though current spec is 1 editor per newsletter).
*   Storing both HTML and plain text versions of content is good practice for email deliverability.
*   When a post is to be published immediately, `scheduled_at` can be null or same as `published_at`, and `status` would move quickly through 'publishing' to 'published'.
*   The `published_at` field should only be set when the newsletter is actually sent to subscribers.

---

This schema provides a foundation. You'll write SQL DDL statements (e.g., `CREATE TABLE ...`) to implement this in Supabase, or use the Supabase UI to create the tables. Remember to set up appropriate indexes for frequently queried columns (e.g., foreign keys, email addresses, tokens) to ensure good performance. 
