# Konbini Documentation

This is a simple documentation for Konbini. Things might not be 1-1 with the current
implementation in `main`.

## Table of Content

- [Routes](#routes)
  - [Sign up / Create an account](#sign-up-create-an-account)
  - [Sign in / Get access and refresh tokens](#sign-in-get-access-and-refresh-tokens)
  - [Get new access token](#get-new-access-token)
  - [Verify email](#verify-email)
  - [Resend verification email](#resend-verification-email)
  - [Forgot password](#forgot-password)
  - [Reset password](#reset-password)
  - [Prepare bento](#prepare-bento)
  - [Order bento](#order-bento)
  - [Rename bento](#rename-bento)
- [Custom tags](#custom-tags)
  - [Error msg tag `errormsg`](#error-msg-tag-errormsg)

## Routes

Documentation on the available routes, depracated routes, and upcoming routes.
Here you will find the route path and method along with the different request bodies
and response bodies.

### Sign up / Create an account

This route handles requests to create a new account. An account is required to prepare new bentos
and manage existing bentos.

```
POST /auth/signup HTTP/1.1
Host: konbini.juancwu.dev
Content-Type: application/json

{
    "email": "your_email@mail.com",
    "password": "strong password",
    "name": "Your Name"
}
```

### Sign in / Get access and refresh tokens

This route handles requests to sign into an account. This route will response with access and refrehs tokens.

```
POST /auth/signin HTTP/1.1
Host: konbini.juancwu.dev
Content-Type: application/json

{
    "email": "your_email@mail.com",
    "password": "strong password"
}
```

### Get new access token

In case the access token had expired, you can get a new one using this route as long as the refresh token is still valid.

```
PATCH /auth/refresh HTTP/1.1
Host: konbini.juancwu.dev
Authorization: Bearer <access_token>
```

### Verify email

Verify the email of an account using the code that was sent to the email.
This route's method is `GET` because this route is also sent to the email to facilitate one click verify.

```
GET /auth/email/verify?code HTTP/1.1
Host: konbini.juancwu.dev

Query:
    code: required, len=20
```

### Resend verification email

It resends an email to verify if the account using the given email has not been verified yet.

```
POST /auth/email/resend HTTP/1.1
Host: konbini.juancwu.dev
Content-Type: application/json

{
    "email": "your@mail.com"
}
```

### Update/Change email

This route will update/change a user's email. The user must provide a valid access token.

```
PATCH /auth/email/update HTTP/1.1
Host: konbini.juancwu.dev
Authorization: Bearer <token>

JSON Body:
    new_email: string

200 OK:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string

400 Bad Request:
    Content-Type: application/json
    JSON Body:
        errors: []string?
        message: string
        request_id: string

401 Unauthorized:
500 Internal Server Error:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string
```

### Forgot password

This route requests a new password reset code for the given email. The code will be sent
to the user's email and it expires in 3 minutes. You must use the [Reset password](#reset-password) route
to finalize the password reset.

```
GET /auth/forgot/password?email HTTP/1.1
Host: konbini.juancwu.dev

Query:
    email: required
```

### Delete account

Request your account to be deleted. This route, if successful, will delete the account immediately.
The user must provide a valid access token.

```
DELETE /auth/account HTTP/1.1
Host: konbini.juancwu.dev
Authorization: Bearer <token>

200 OK:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string

401 Unauthorized:
500 Internal Server Error:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string
```

### Reset password

This routes is to finalize the password reset process. It requires the email and the code that
its gotten from [Forgot password](#forgot-password) route.

```
Details not implemented yet.
```

### Prepare bento

It prepares a new bento and stores it in the database. You will need an access token.

```
POST /bento/prepare HTTP/1.1
Host: konbini.juancwu.dev
Content-Type: application/json
Authorization: Bearer <token>

JSON Body:
    name: required,min=3,max=50,ascii
    pub_key: required
    ingridients?: []{ name: string, value: string }

200 OK: Bento prepared but failed to add ingridients (if provided)
    Content-Type: application/json
    JSON Body:
        message: string
        bento_id: string

201 Created: Bento prepared and ingridients added (if provided)
    Content-Type: application/json
    JSON Body:
        message: string
        bento_id: string

400 Bad Request:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string
        errors?: []string

403 Unauthorized:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string

500 Internal Server Error:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string
```

### Order bento

This will get you the bento with all the ingridients in it. There must be an existing bento before ordering it.

```
GET /bento/order/:bento_id HTTP/1.1
Host: konbini.juancwu.dev
Authorization: Bearer <token>

Query:
    signature: signature using the RSA private key for the bento.
    challenge: a random message to sign with the RSA private key.

200 OK:
    Content-Type: application/json
    JSON Body:
        message: string
        ingridients: []{ name: string, value: string }

400 Bad Request:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string

404 Not Found:
    No content

500 Internal Server Error:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string
```

### Rename ingridient

This route will rename an ingridient if it exists.

```
PATCH /bento/ingridient/rename HTTP/1.1
Host: konbini.juancwu.dev
Authorization: Bearer <token>

JSON Body:
    bento_id: string
    new_name: string
    old_name: string

200 OK:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string

400 Bad Request:
    Content-Type: application/json
    JSON Body:
        errors: []string?
        message: string
        request_id: string

401 Unauthorized:
500 Internal Server Error:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string
```

### Re-season ingridient

This route will "re-season", change the value of an ingridient.

```
PATCH /bento/ingridient/reseason HTTP/1.1
Host: konbini.juancwu.dev
Authorization: Bearer <token>

JSON Body:
    bento_id: string
    name: string
    value: string

200 OK:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string

400 Bad Request:
    Content-Type: application/json
    JSON Body:
        errors: []string?
        message: string
        request_id: string

401 Unauthorized:
500 Internal Server Error:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string
```

### Delete ingridient

This route will delete an existing ingridient from a bento.

```
DELETE /bento/ingridient HTTP/1.1
Host: konbini.juancwu.dev
Authorization: Bearer <token>

JSON Body:
    bento_id: string
    name: string

200 OK:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string

400 Bad Request:
    Content-Type: application/json
    JSON Body:
        errors: []string?
        message: string
        request_id: string

401 Unauthorized:
500 Internal Server Error:
    Content-Type: application/json
    JSON Body:
        message: string
        request_id: string
```

## Custom tags

This section covers all the custom tags with explanation on how to use them and where they are in the project.

### Error msg tag `errormsg`

The `errormsg` tag can be used to define the error message of a validation error using golang's validator `github.com/go-playground/validator/v10`.
Since you can have multiple tags with errors, you can define the error message specifically for a tag
just like how you define multiple validation tags.

In case the validation tag is not found and no global error message is provided then an empty string is returned.

Example usage:

```go
type User struct {
    Name string `validate:"required,min=3" errormsg:"required=Name field is requried,min=Name must be at least 3 characters long"`
}
```

Or if you want to just use a global error message:

```go
type User struct {
    // required does not have any msg so it will just try to get __default
    Name string `validate:required,min=3" errormsg:"required,__default=Name field has error(s)"`
    Surname string `validate:required,min=3" errormsg:"Name field has error(s)"`
}
```

Or combine two or more tags with `|` to use the same error message:

```go

type User struct {
    Name string `validate:required,min=3,max=10" errormsg:"required=Some message,min|max=Too short/long"` // that's what she said...
}
```
