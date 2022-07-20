### MeoCloud API Client in GO

This module is intended to provide an implementation in GO to consume the API of the Meocloud Storage for the Meo customers in Portugal.

### Note
For now, since this is my first contact with the language, I do not provide a way to fetch the access token to be used to try this module. You can take a look at the sister projects in PHP, [php-meocloud-cli](https://github.com/digfish/php-meocloud.cli) or [meocloudrepl](https://github.com/digfish/meocloudrepl) in Python, and following the directions, generate the access token and secret, this two more the consumer token and secret which can be obtained at https://meocloud.pt/my_apps. These four should be inscribed in a [.env file](https://www.dotenv.org/) with the following format:
```
CONSUMER_KEY=
CONSUMER_SECRET=
OAUTH_TOKEN=
OAUTH_TOKEN_SECRET=
```
I intend to write new functions to get the credentials later.

### What is implemented

|   Method          |    API                   |
|-------------------|--------------------------|
| account_info()    | GET Account/Info         |
| get_metadata()    | GET Metadata/meocloud/   |
| get_file()        | GET Files/meocloud/:name |
| send_file()       | PUT Files/meocloud/:name |
| delete_file()     | POST Fileops/Delete      |
| create_dir()      | POST Fileops/CreateFolder|
------------------------------------------------
 
