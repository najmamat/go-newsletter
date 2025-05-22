# Editor login
- Editor submits login credentials
- System validates input
- System verifies that a user with the given email exists
- System checks if the password matches
- If email does not exist or password is incorrect, process terminates with error 2`. 
  - **Error does not specify wheter email or password were incorect**
- If authentication succeeds, system generates access token
- System returns response with token