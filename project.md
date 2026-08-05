.
├── LICENSE
├── Makefile
├── README.md
├── bin
│   ├── api
│   └── notifier
├── cmd
│   ├── api
│   │   └── main.go
│   └── notifier
│       └── main.go
├── db
│   └── migrations
│       ├── 001_create_users_table.down.sql
│       ├── 001_create_users_table.up.sql
│       ├── 002_create_refresh_tokens_table.down.sql
│       ├── 002_create_refresh_tokens_table.up.sql
│       ├── 003_create_categories_table.down.sql
│       ├── 003_create_categories_table.up.sql
│       ├── 004_create_products_table.down.sql
│       ├── 004_create_products_table.up.sql
│       ├── 005_create_product_images_table.down.sql
│       ├── 005_create_product_images_table.up.sql
│       ├── 006_create_carts_table.down.sql
│       ├── 006_create_carts_table.up.sql
│       ├── 007_create_cart_items_table.down.sql
│       ├── 007_create_cart_items_table.up.sql
│       ├── 008_create_orders_table.down.sql
│       ├── 008_create_orders_table.up.sql
│       ├── 009_create_order_items_table.down.sql
│       └── 009_create_order_items_table.up.sql
├── docker
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── init
│   │   └── aws
│   │       └── init-localstack.sh
│   └── nginx
│       └── nginx.conf
├── docs
│   ├── docs.go
│   ├── rapidoc.html
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── database.go
│   ├── dto
│   │   ├── auth.go
│   │   ├── order.go
│   │   └── product.go
│   ├── events
│   │   ├── publisher.go
│   │   └── watermill.go
│   ├── interfaces
│   │   └── upload.go
│   ├── logger
│   │   └── logger.go
│   ├── models
│   │   ├── order.go
│   │   ├── product.go
│   │   └── user.go
│   ├── providers
│   │   ├── aws.go
│   │   ├── local_provider.go
│   │   └── s3.go
│   ├── server
│   │   ├── auth_handlers.go
│   │   ├── cart_handler.go
│   │   ├── middlewares.go
│   │   ├── order_handler.go
│   │   ├── product_handlers.go
│   │   └── server.go
│   ├── services
│   │   ├── auth_service.go
│   │   ├── cart_service.go
│   │   ├── order_service.go
│   │   ├── product_service.go
│   │   ├── upload_service.go
│   │   └── user_service.go
│   └── utils
│       ├── jwt.go
│       ├── password.go
│       └── response.go
├── notifications
│   └── email.go
├── project.md
└── uploads
    └── products
        ├── 1
        │   └── khaleja-wallpaper.jpg
        ├── 2
        │   ├── Screenshot 2026-06-03 193254.png
        │   └── khaleja-wallpaper.jpg
        ├── 3
        │   ├── 0c508af8-f9f8-4d30-8cbf-0301288dd39a.png
        │   ├── Screenshot 2026-06-03 193254.png
        │   └── khaleja-wallpaper.jpg
        ├── 4
        │   ├── Screenshot 2026-06-03 193254.png
        │   └── khaleja-wallpaper.jpg
        └── 6
            └── Screenshot 2026-06-03 193254.png

32 directories, 76 files
