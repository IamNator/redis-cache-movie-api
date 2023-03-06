# Migrations


1. Creating a new migration

    ```bash
    $  migrate create -dir database/migrations -ext sql -seq <title>
    ```

For example, if you want to create a migration for creating a table called `users`, you can run the following command:

    ```bash
    $  migrate create -dir database/migrations -ext sql -seq create_users_table
    ```

2. Running migrations

    ```bash
    $  migrate -path database/migrations -database "postgres://localhost:5432/<database_name>?sslmode=disable" up
    ```

3. if this error occurs:
    ```bash
    $ error: Dirty database version <version_no> . Fix and force version.
    ```
   then run this command:
    ```bash
    $  migrate -path database/migrations -database "postgres://localhost:5432/<database_name>?sslmode=disable" force <version_no>
    ```
   then run the command in step 2 again.

for example, if the error is `error: Dirty database version 2 . Fix and force version.`, then run the following command:
```bash
$  migrate -path database/migrations -database "postgres://localhost:5432/<database_name>?sslmode=disable" force 2
```
then run the command in step 2 again.

4. Rolling back migrations

    ```bash
    $  migrate -path database/migrations -database "postgres://localhost:5432/<database_name>?sslmode=disable" down
    ```

for more information, please refer to [migrate](https://www.freecodecamp.org/news/database-migration-golang-migrate/)
