| ID | Err code                    | Err message                                               | Response |
|----|-----------------------------|-----------------------------------------------------------|----------|
| 1  | EMAIL_ALREADY_EXISTS        | An account with this email already exists.                | 409      |
| 2  | INVALID_CREDENTIALS         | Email or password is incorrect.                           | 401      |
| 3  | UNAUTHORIZED                | Authentication is required to access this resource.       | 401      |
| 4  | FORBIDDEN                   | You do not have permission to perform this action.        | 403      |
| 5  | RESOURCE_NOT_FOUND          | The requested resource does not exist.                    | 404      |
| 6  | SUBSCRIPTION_ALREADY_EXISTS | You are already subscribed to this newsletter.            | 409      |
| 7  | INVALID_TOKEN               | The provided token is invalid or expired.                 | 401      |
| 8  | MISSING_REQUIRED_FIELDS     | One or more required fields are missing or empty.         | 400      |
| 9  | INTERNAL_SERVER_ERROR       | An unexpected error occurred on the server.               | 500      |
| 10 | SCHEDULE_IN_PAST            | Scheduled time must be in the future.                     | 400      |
| 11 | UNSUBSCRIBE_TOKEN_INVALID   | The unsubscribe token is invalid or already used.         | 400      |
| 12 | EMAIL_NOT_CONFIRMED         | Email must be confirmed before performing this action.    | 403      |
| 13 | SUBSCRIPTION_NOT_FOUND      | Subscription for this token was not found.                | 404      |
| 14 | INVALID_CONFIRMATION_TOKEN  | Confirmation token is invalid or expired.                 | 400      |
| 15 | PASSWORD_POLICY_VIOLATION   | Password does not meet minimum security requirements.     | 400      |