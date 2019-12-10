# Data Bunker API


**Data Bunker is advanced personal information tokenization and storage service build to comply with GDPR.**

---

## User Api


| Resource / HTTP method | POST (create)    | GET (read)    | PUT (update)     | DELETE (delete)  |
| ---------------------- | ---------------- | ------------- | ---------------- | ---------------- |
| /v1/user               | Create new user  | Error         | Error            | Error            |
| /v1/user/token/{token} | Error            | Get user      | Update user      | Delete user PII  |
| /v1/user/login/{login} | Error            | Get user      | Update user      | Delete user PII  |
| /v1/user/email/{email} | Error            | Get user      | Update user      | Delete user PII  |
| /v1/user/phone/{phone} | Error            | Get user      | Update user      | Delete user PII  |


## Create user record
### `POST /v1/user`

### Explanation

This API is used to create new user record and if the request is successful it returns new `{token}`.
On the database level, each records is encrypted with it's own key.


### POST Body Format

POST Body can contain regular form data or JSON. Data Bunker extracts `{login}`, `{phoen}` and `{email}` out of
POST data or from JSON first level and builds additional hashed index for user object. These fields, if
provided must be unique, otherwise you will ge an error. So, you can not create additional user object
with duplicate email.

The following content type supported:

* **application/json**
* **application/x-www-form-urlencoded**


### Example:

Create used by posting JSON:

```
curl -s http://localhost:3000/v1/user -XPOST \
  -H "X-Bunker-Token: $XTOKEN" \
  -H "Content-Type: application/json" \
  -d '{"firstName": "John","lastName":"Doe","email":"user@gmail.com"}'
{"status":"ok","token":"db80789b-0ad7-0690-035a-fd2c42531e87"}
```

Create user by POSTing user key/value fiels as post parameters:

```
curl -s http://localhost:3000/v1/user -XPOST \
  -H "X-Bunker-Token: $XTOKEN" \
  -d 'firstName=John' \
  -d 'lastName=Doe' \
  -d 'email=user2@gmail.com'
{"status":"ok","token":"db80789b-0ad7-0690-035a-fd2c42531e87"}
```

**NOTE**: Keep this user token privately as it provides user private information in your system.
For semi-trusted environments or 3rd party companies, use **shareable identity** instead.


---

## Get user record
### `GET /v1/user/{token,login,email,phone}/{indexValue}`

### Explanation
This API is used to get user PII records. You can lookup user token by **token**, **email**, **phone** or **login**.

### Example:

Fetch by user token:

```
curl --header "X-Bunker-Token: $XTOKEN" -XGET \
   https://localhost:3000/v1/user/token/DAD2474A-E9A7-4BA7-BFC2-C4506880198E
{"status":"ok","token":"DAD2474A-E9A7-4BA7-BFC2-C4506880198E","data":{"k1":[1,10,20],
"k2":{"f1":"t1","f3":{"a":"b"}},"login":"user1","name":"tom"}}
```

Fetch by "login" name:

```
curl --header "X-Bunker-Token: $XTOKEN" -XGET \
   https://localhost:3000/v1/user/login/user1
{"status":"ok","token":"DAD2474A-E9A7-4BA7-BFC2-C4506880198E","data":{"k1":[1,10,20],
"k2":{"f1":"t1","f3":{"a":"b"}},"login":"user1","name":"tom"}}
```


---

## Update user record
### `PUT /v1/user/{token,login,email,phone}/{indexValue}`

### Explanation

This API is used to update user record. You can update user by **token**, **email**, **phone** or **login**.
This call returns update status on success or error message on error.

### POST Body Format

POST Body can contain regular form POST data or JSON. When using JSON, you can remove the record by setting it's value to null.
For example {"key-to-delete":null}.

The following content type supported:

* **application/json**
* **application/x-www-form-urlencoded**


### Example:

The following command will change user name to "Alex". An audit event will be generated showing previous and new value.

```
curl --header "X-Bunker-Token: $XTOKEN" -d 'name=Alex' -XPUT \
   https://localhost:3000/v1/user/token/DAD2474A-E9A7-4BA7-BFC2-C4506880198E
```

---

## Delete user by record
### `DELETE /v1/user/{token,login,email,phone}/{indexValue}`

This command will remove all user records from the database, leaving only user token id.

```
curl -header "X-Bunker-Token: $XTOKEN" -XDELETE \
  https://localhost:3000/v1/user/token/DAD2474A-E9A7-4BA7-BFC2-C4506880198E
{"status":"ok","result":"done"}
```

## User App Api


| Resource / HTTP method            | POST (create)       | GET (read)        | PUT (update)  | DELETE |
| --------------------------------- | ------------------- | ----------------- | ------------- | ------ |
| /v1/userapp/token/:token/:appname | Create new user app | Get record        | Change record | Delete |
| /v1/userapp/token/:token          | Error               | Get user app list | Error         | Error  |
| /v1/userapp/list                  | Error               | Get all app list  | Error         | Error  |


## Create user app record
### `POST /v1/userapp/token/:token/:appname`

### Explanation

This API is used to create new user app record and if the request is successful it returns `{"status":"ok"}`.


---

## User Session Api

| Resource / HTTP method       | POST (create)      | GET (read)     | PUT (update)   | DELETE (delete) |
| ---------------------------- | ------------------ | -------------- | -------------- | --------------- |
| /v1/session/token/:token      | Create new session | Get sessions   | Error          | Error           |
| /v1/session/session/:session | Error              | Get session    | Error??        | Error??         |
| /v1/session/clientip/:ip     | Error              | Get sessions   | Error          | Error           |

## Create user session record
### `POST /v1/session/token/:token`

### Explanation

This API is used to create new user session and if the request is successful it returns new `{session}`.
You can now use this id in your logs instead of user IP and browser user-agent info, etc...

Our API supports generation of session tokens based on the following information:
user ip, mobile device info, user agent, etc...

You can send the data as JSON POST or as regular POST parameters when working with this API.

## Get user session record
### `GET /v1/session/session/:session`

### Explanation

This API returns session data.


## Get session records by user token.
### `GET /v1/session/token/:token`

### Explanation

This API returns an array of sessions of the same user.

## Get session records by ip address.
### `GET /v1/session/clientip/:ipaddress`

### Explanation

This API returns an array of user sessions by IP address. These sessions can be of different people.

---

## Passwordless tokens API

| Resource / HTTP method | POST (create)     | GET (read)    | PUT (update)     | DELETE (delete)  |
| ---------------------- | ----------------- | ------------- | ---------------- | ---------------- |
| /v1/xtoken/:token      | Create new record | Error         | Error            | Error            |
| /v1/xtoken/:token      | Error             | Get data      | Error            | Error            |

	router.POST("/v1/xtoken/:token", e.userNewToken)
	router.GET("/v1/xtoken/:xtoken", e.userCheckToken)


---

## Shareable token API

| Resource / HTTP method | POST (create)     | GET (read)    | PUT (update)     | DELETE (delete)  |
| ---------------------- | ----------------- | ------------- | ---------------- | ---------------- |
| /v1/shareable/:token   | Create new record | Error         | Error            | Error            |
| /v1/shareable/:token   | Error             | Get data      | Error            | Error            |


---

**TODO-FINISH**


## Temporary user access tokens

Sometimes, for example, when working with 3rd party partners or semi-trusted environments, you might
need to generate a user access token with a specific expiration time. Your partner can retrieve user
information during this specific time only.
Afterward, access will be blocked.

The following command will generate a token to access user email and name for 7 days:

```
curl --header "X-Bunker-Token: $XTOKEN" -d 'fields=email,name' -d 'expiration=7d' -d 'partner=sms' \
   https://bunker.company.com/gentokens/DAD2474A-E9A7-4BA7-BFC2-C4506880198E
```

Output:

```
476E41E7-72AD-448A-BB43-7ACDB8C53735
```


### 3rd party logging

Instead of maintaining internal logs, a lot of companies are using 3rd party logging facility like logz or coralogix or something else.
To improve adherence to GDPR, we build a special feature - generate specific session id for such 3rd party service.

When using these uuids in external systems, you basically **pseudonymise personal data**. In addition, in accordance with GDPR Article 5:
**Principles relating to processing of personal data**. Personal data shall be: (c) 
adequate, relevant and limited to what is necessary in relation to the purposes for which they are processed (‘**data minimisation**’);

Here is a command to do it:

```
curl -d 'ip=user@example.com' \
     -d 'user-agent=mozila' \
     -d 'partner=coralogix' \
     -d 'expiration=7d'\
   https://bunker.company.com/gensession/DAD2474A-E9A7-4BA7-BFC2-C4506880198E
```

It will generate a new token, that you can now pass to 3rd party system as a user id.


## User consent management

One of the GDPR requirements is the storage of user consent. For example, your customer must approve to receive email marketing information.

Using the GDPR language, your customer must give explicit consent to receive marketing information.

Consent must be freely given, specific, informed and unambiguous. From GDPR, Article 7, item 3:

* **The data subject shall have the right to withdraw his or her consent at any time.**
* **It shall be as easy to withdraw as to give consent.**

To comply with this requirement, we added support to manage user consent. We support the following APIs:


### List granted

```
curl --header "X-Bunker-Token: $XTOKEN" \
   https://bunker.company.com/consent/DAD2474A-E9A7-4BA7-BFC2-C4506880198E
```

### List all

```
curl --header "X-Bunker-Token: $XTOKEN" \
   https://bunker.company.com/consent/DAD2474A-E9A7-4BA7-BFC2-C4506880198E?all
```

### Cancel consent

```
curl --header "X-Bunker-Token: $XTOKEN" -XDELETE \
   https://bunker.company.com/consent/DAD2474A-E9A7-4BA7-BFC2-C4506880198E/<consent-id>
```

### User gives consent

**TODO**

### Easily cancel consent for email marketing

For example, for email marketing, users got distracted, when they need to login in order to unsubscribe from the newsletter.
To simplify this operation, users will be allowed to unsubscribe only using email address without full login operation.

### Unlock bunker

Run the following command with different keys:

```
bunker unlock **key**
```

Or you can provide multiple keys at once:

```
bunker unlock key1 key2 key3
```

### View lock status

```
bunker status | jq .lock
```

Result:

```
locked
```


## Audit API


It is not compliant, unless you have a real reason to share this specific personal sub-record. For example,
sending customer phone when notifying customer using 3rd party SMS gateway.



# SECTION IS NOT UPDATED BELLOW

## Data Bunker init

Upon initial init, the Data Bunker service will check if the system is initialized for the first time, and if yes,
it will generate root password, master key and derived keys out of it. Otherwise, an error will be printed.

```
bunker init
```

Output:

```
Root password: 123456
Key1: abcdefg
Key2: abcdefg
key3: abcdefg
Key4: abcdefg
Key5: abcdefg
```

**TODO**: Secret keys printed to output can be easily extracted in cloud environments for example in Kubernetes logs!