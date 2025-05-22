# Editor registration
## Proces
- User submits registration (email, name, password)
- System checks if user with given email already exists -> Else err 1
- System checks password strenght. Password must be min. 8 chars long, min. 1 number, min. 1 special character. -> else err 15
- System creates record in `profiles` table
- System generates access token
- System returns response with token

## Questions
- Should registration be confirmed via email?