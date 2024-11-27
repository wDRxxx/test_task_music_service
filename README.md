# Online-music library service

## How to run service
1. Configure .env file
2. Optionally run database in docker 
     ```shell
    docker-compose up -d
    ```
3. run service
    
    ```shell
    go run ./cmd/api/main.go
    ```

Swagger documentation is available on /swagger/ route

Always edit date formats in swagger, it may break some requests

To update swagger documentation run 
```shell
swag init --generalInfo cmd/api/main.go --output ./docs
```