BUSHA MOVIE API

[Documentation](https://busha-movie-api-v1.herokuapp.com/docs)

`` 
<host>/docs
 ``
 ---

### Getting Started (Development)

1. Install dev tools
```bash 
 $ make install_tools
```
2. Run Migrations [migration](./database/README.md)
```
 $ migrate -path database/migrations -database "postgres://postgres:password@localhost:5432/busha?sslmode=disable" up
```
3. Spin up docker-compose 
```bash 
 $ docker-compose up
```
4. Test if server is running
```bash
 $ curl localhost:9500/health
```
5. go to localhost:9500/docs to view api docs

---
### Useful Links   
[swaggo](https://github.com/swaggo/swag#declarative-comments-format)
---


### TODO

1. Setup CI/CD to run unit & integration tests
2. Automate generation of swagger docs
3. Add rate limiting to api
2. Optimize fetching of movies and characters from external api
3. Refactor code make use of more custom types and constants
4. Improve Error Codes Used (currently using 500 for most errors)
3. Write more tests for handlers, services and repositories


### Author
[Nator Verinumbe](github.com/iamnator)
