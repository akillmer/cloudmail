# cloudmail
A Google Cloud Function that relays messages, via SMTP, that get posted from a simple contact form. See it in action at https://mixplate.io/contact. This function also validates submissions via Recaptcha v2.

## Gmail SMTP requires less secure access
This function uses simple SMTP authorization with Gmail, but to allow that you need to enable _Less secure app access_ on the account. It's best to do this with a new account, not your primary one. To allow simple authorization, from within your Gmail account navigate to `Settings` → `Accounts and Import` → `Other Google Account settings` → `Security` → `Less Secure app access`.

## Setting up the Cloud Function
On the Google Cloud Platform it's simple to add a new Cloud Function, really just point and click. I keep the memory allocation to 128 MB but note that Google considers that a _testing_ environment. I typically just copy and paste into the provided inline editor. The function to execute is `SendMessage`.

There are several enviornment variables that you must add at this time:

`RECAPTCHA_SECRET`: provided by your Recaptcha account (Version 2)
`SMTP_ADDR`: for Gmail it's _smtp.gmail.com_
`SMTP_PORT`: for Gmail it's _587_
`SMTP_USER`: the email address that will be relaying messages for you
`SMTP_PW`: password for the above email account
`MAIL_TO`: where the messages will be relayed to, e.g. your primary inbox

## Posting to the function
For my Vue.js app I `POST` to the endpoint that Google will provide for your Cloud Function. The request's body results in a JSON string like this:

```json
{
    "recaptcha": ...,
    "name": "User's Name",
    "replyTo": "users@email.address",
    "message": "Hello, world."
}
```

`recaptcha` is the key that is automatically provided after the user confirms they're not a robot.

Note that this Cloud Function expects to decode JSON from the request's body. If you use a different method, such as query parameters, you must update the function to handle that.

Check out https://github.com/akillmer/portfolio/blob/master/src/views/Contact.vue to see how I implemented the client side.