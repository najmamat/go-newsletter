# Extra Features Chosen

## 1. Newsletter Scheduling & Queues

**Description:**
This feature allows editors to compose newsletter posts and schedule them for publication at a future date and time. Instead of publishing immediately, the post is added to a queue.
A background scheduling mechanism within the API will monitor this queue and automatically trigger the publishing process (sending emails to subscribers and archiving the post) when the scheduled time arrives.

**Benefits:**
*   **Flexibility for Editors:** Editors can prepare content in advance and control when it's delivered, even if they are not actively online.
*   **Consistent Timing:** Allows for strategic timing of newsletter releases to maximize engagement.
*   **Batch Processing:** Can help in managing server load if many newsletters are published, though less critical for the initial scale.

**Key Components (High-Level):**
*   **Scheduling Interface:** Editors will need a way to specify the desired publication date/time when creating/editing a post.
*   **Scheduled Posts Storage:** The database will need to store the post content along with its scheduled publication time and current status (e.g., "scheduled", "processing", "sent").
*   **Scheduler/Worker Process:** A background component in the Go API responsible for:
    *   Periodically checking for posts due for publication.
    *   Initiating the publishing workflow for due posts.
    *   Handling potential retries or failures.
*   **Queue Mechanism:** Implicitly, a list or table of scheduled posts acts as a queue.
