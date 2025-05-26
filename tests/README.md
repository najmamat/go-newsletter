# Postman API Testing for Go Newsletter

This guide explains how to set up and use Postman for testing the Go Newsletter API using a shared Postman collection.

## 1. Set Up a Postman Environment

Before importing the collection, you need to set up a Postman Environment. This environment will hold your specific Supabase credentials and API details, keeping them separate from the collection itself.

1.  In Postman, click the "Environments" tab on the left sidebar.
2.  Click the `+` button to create a new environment.
3.  Name it something descriptive, e.g., "GoNewsletter Dev".
4.  Add the following variables to the environment. You can get the Supabase URL and Anon Key from your Supabase project settings (Project Settings > API).

    | VARIABLE        | INITIAL VALUE                                     | TYPE    | NOTES                                                                 |
    | :-------------- | :------------------------------------------------ | :------ | :-------------------------------------------------------------------- |
    | `supabaseUrl`   | `https://<your-project-ref>.supabase.co`          | default | Your Supabase project URL.                                            |
    | `supabaseAnonKey` | `<your-supabase-anon-key>`                      | default | Your Supabase project's public anonymous key.                       |
    | `baseUrl`       | `http://localhost:8080/api/v1`                           | default | Base URL of your Go API (e.g., where your local server runs).         |
    | `bearerToken`           | (leave blank, will be auto-populated)             | secret  | Auto-filled by the 'Auth: Sign In (Get JWT)' request in the collection. |
    | `user_id`       | (leave blank, will be auto-populated)             | default | Auto-filled by the 'Auth: Sign In (Get JWT)' request.            |

5.  **Important**: After creating and configuring the environment, ensure you select it from the environment dropdown in the top-right corner of Postman before running any requests.

## 2. Import the Shared Postman Collection

Your team lead/member will provide a Postman collection file (e.g., `Go Newsletter API.postman_collection.json`). This file is also available in the `/tests` directory of this project.

1.  Open Postman.
2.  Click on "Import" (usually in the top left).
3.  Select the Postman collection JSON file (e.g., `go_newsletter_api.postman_collection.json` from the `tests` directory).
4.  This will import a collection likely named "Go Newsletter API" (or similar) with pre-configured requests.

## 3. Running Requests

### a. Sign In to Supabase and Get JWT

The imported collection contains a request specifically for authenticating with Supabase and retrieving a JWT. It is likely named "Auth: Sign In (Get JWT)" or similar.

**How this request is configured:**
*   **Method:** `POST`
*   **URL:** `{{supabaseUrl}}/auth/v1/token?grant_type=password` (Uses your environment variables)
*   **Headers:** Includes `apikey: {{supabaseAnonKey}}` and `Content-Type: application/json`.
*   **Body (raw JSON):** Contains placeholders for email and password.
*   **Post-Response Script:** Contains JavaScript to automatically extract the `access_token` and `user.id` from the Supabase response and save them to your `{{bearerToken}}` and `{{user_id}}` environment variables, respectively. The script also includes basic tests to verify the response structure.

**Steps to use it:**
1.  Find this "Auth: Sign In (Get JWT)" request in the imported collection.
2.  Go to its "Body" tab.
3.  Update the `email` and `password` fields with *your* valid test user credentials for your Supabase project:
    ```json
    {
        "email": "your-actual-test-email@example.com",
        "password": "your-actual-password"
    }
    ```
4.  Ensure your Go API server is running locally if you intend to test it immediately after (`make run` or `make dev`).
5.  Click "Send" for the Supabase sign-in request.
6.  If successful, the embedded post-response script will automatically populate your `{{bearerToken}}` and `{{user_id}}` environment variables. You can verify this by:
    *   Checking the "Test Results" tab for this request for success messages (e.g., "Response contains an access token").
    *   Clicking the "eye" icon next to your environment name in the top right of Postman to view current environment variable values.

### b. Testing Your API's Protected Endpoints

1.  Once the `{{bearerToken}}` environment variable is populated by the Supabase sign-in request, you can run other requests in the collection that target your Go API's protected endpoints (e.g., "GET Get Current Editor Profile", which calls `{{baseUrl}}/api/v1/me`).
2.  These requests are pre-configured in the shared collection to use `{{bearerToken}}` for Bearer Token authentication (found under the "Authorization" tab of the request).
3.  Simply open the desired request and click "Send".

## 4. Adding More Requests

*   The shared collection provides a good starting point. You can duplicate existing requests or add new ones to test other endpoints of your Go API as they are developed.
*   For new requests to protected endpoints, remember to:
    *   Set the URL to `{{baseUrl}}/api/v1/...your-endpoint...`.
    *   Go to the "Authorization" tab, select Type: `Bearer Token`, and set the Token field to `{{bearerToken}}`.

