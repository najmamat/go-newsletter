# How to Set Up Postman from OpenAPI Specification in GitHub

This guide outlines the steps to connect Postman to an OpenAPI specification hosted in a GitHub repository, set up an environment, and automate the extraction of authentication tokens. This allows for efficient API testing and collaboration.

## Prerequisites

*   A Postman account and the Postman desktop application.
*   A GitHub account with access to the repository containing the OpenAPI specification file (e.g., `openapi.yaml`).
*   The OpenAPI specification should define authentication endpoints that return a JWT and user information upon successful login.

## Steps

### 1. Connect Postman to GitHub

To import and sync your API definition, you first need to connect Postman to your GitHub account and grant it access to the relevant repository.

1.  In Postman, navigate to **APIs** in the left sidebar.
2.  Click on **Import** or **Create an API**. If creating a new API, give it a name.
3.  You'll be prompted to connect to a repository. Choose **GitHub**.
4.  Follow the on-screen instructions to authorize Postman to access your GitHub account.
5.  Select your GitHub organization/account and the repository that contains your OpenAPI specification.
6.  Choose the branch where your OpenAPI file is located (e.g., `main`, `develop`).
7.  Specify the path to your OpenAPI file within the repository (e.g., `api/openapi.yaml`).
8.  Postman will then import the API definition.

### 2. Set Up a Postman Environment

An environment in Postman allows you to store and manage variables that can be used across your requests, such as API base URLs, authentication tokens, and user IDs.

1.  In Postman, click the **Environments** tab on the left sidebar.
2.  Click the `+` button to create a new environment.
3.  Name it descriptively, for example, "[Project Name] Dev".
4.  Add the following variables. The `INITIAL VALUE` for `jwt` and `user_id` can be left blank as they will be auto-populated by a script later. You'll also likely need a `baseUrl` for your API.

    | VARIABLE    | INITIAL VALUE                 | TYPE    | NOTES                                                                   |
    | :---------- | :---------------------------- | :------ | :---------------------------------------------------------------------- |
    | `baseUrl`   | `http://localhost:8080`       | default | Your API's base URL (update as needed).                                 |
    | `bearerToken`       | (leave blank)                 | secret  | Will be auto-populated by the Sign-In request.                          |
    | `user_id`   | (leave blank)                 | default | Will be auto-populated by the Sign-In request.                          |
    | `supabaseUrl` | (if using Supabase directly)  | default | e.g., `https://<your-project-ref>.supabase.co` (optional, adjust as needed) |
    | `supabaseAnonKey` | (if using Supabase directly) | default | Your Supabase anon key (optional, adjust as needed)                     |

    *   **`default`**: Visible and included in exported environments.
    *   **`secret`**: Masked in the UI and **not** included in exported environments.

5.  **Important**: After creating the environment, ensure it's selected from the environment dropdown in the top-right corner of Postman before running requests.

### 3. Generate/Import Collection from OpenAPI Specification

Once your API is imported from GitHub, Postman can generate a collection from it, or you might import a pre-existing collection if your team uses one. If you've connected Postman to your GitHub repo containing the OpenAPI spec as described in Step 1, Postman often offers to generate a collection automatically.

1.  If Postman hasn't automatically created a collection, or you need to refresh it:
    *   Go to the **APIs** tab and select your imported API.
    *   Look for an option like **Generate Collection** or ensure your API definition is linked to a collection. Postman usually allows you to configure how the collection is generated (e.g., based on specific versions of your API definition).
2.  This collection will contain requests based on the paths and operations defined in your `openapi.yaml` file.

### 4. Create/Configure a Sign-In HTTP Request

You'll need a request to authenticate with your API and obtain the JWT and user ID. If your OpenAPI spec included an authentication endpoint, a request for it might already be in the generated collection. If not, or if it needs customization:

1.  Find the sign-in request in your collection (e.g., "POST /auth/signin") or create a new HTTP request.
    *   **Method:** `POST`
    *   **URL:** `{{baseUrl}}/auth/signin` (or the relevant sign-in path from your API). If authenticating directly against Supabase, this might be `{{supabaseUrl}}/auth/v1/token?grant_type=password`. Adjust the URL and any necessary headers (like `apikey: {{supabaseAnonKey}}` for Supabase) according to your API's authentication mechanism.
2.  Go to the **Body** tab.
3.  Select **raw** and choose **JSON** from the dropdown.
4.  Provide the necessary credentials. You can use environment variables or placeholder text that users will replace:
    ```json
    {
        "email": "your-test-email@example.com",
        "password": "your-test-password"
    }
    ```

### 5. Add Post-Response Script to Populate Environment

This script will run after the sign-in request receives a response. It extracts the `access_token` and `user.id` and saves them to your active Postman environment.

1.  Open the sign-in request.
2.  Go to the **Tests** tab (this is where post-response scripts are written).
3.  Add the following JavaScript code:

    ```javascript
    // Parse the JSON response
    var responseJSON = pm.response.json();

    // Test to check if the response contains an access token
    pm.test("Response contains an access_token", function () {
        pm.expect(responseJSON).to.have.property('access_token');
    });

    // Store the access_token in the 'jwt' environment variable
    if (responseJSON.access_token) {
        pm.environment.set("bearerToken", responseJSON.access_token);
        console.log("JWT set in environment.");
    } else {
        console.log("Access token not found in response.");
    }

    // Test to check if the user object and user ID are present in the response
    pm.test("User object with id is present in the response", function () {
        pm.expect(responseJSON).to.have.property('user').that.is.an('object');
        pm.expect(responseJSON.user).to.have.property('id');
    });

    // Store the user.id in the 'user_id' environment variable
    if (responseJSON.user && responseJSON.user.id) {
        pm.environment.set("user_id", responseJSON.user.id);
        console.log("User ID set in environment.");
    } else {
        console.log("User ID not found in response user object.");
    }
    ```

4.  **Save the request.**

### Using the Setup

1.  Ensure your API server is running.
2.  Select the correct Postman environment (e.g., "[Project Name] Dev").
3.  Run the Sign-In request.
4.  If successful, the `jwt` and `user_id` variables in your environment will be populated. You can verify this by clicking the "eye" icon next to your environment name.
5.  Other requests in your collection that require authentication can now be configured to use the `{{bearerToken}}` variable as a Bearer Token in their "Authorization" tab.