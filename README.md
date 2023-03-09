BUSHA MOVIE API

[Documentation](https://busha-movie-api-v1.herokuapp.com/docs)

### Getting Started (Development)

1. Install dev tools
```bash 
 $ make install_tools
```
2. Spin up docker-compose 
```bash 
 $ docker-compose up
```
3. go to localhost:{PORT}/docs to view api docs


### Migration
1. Go to [migration](./database/README.md) for more info

### Useful Links   
[swaggo](https://github.com/swaggo/swag#declarative-comments-format)



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