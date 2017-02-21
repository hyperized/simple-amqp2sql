# Simple AMQP to SQL demo

## Environment
This demo requires a basic MySQL and RabbitMQ setup.

## Docker stuff
* `docker run -d --hostname rabbitmq --name rabbitmq -p 15672:5672 -p 8080:15672 rabbitmq:3-management`

* `docker run --name mysql56 -e MYSQL_ROOT_PASSWORD=secret -e MYSQL_ROOT_HOST=172.17.0.1 -p 6603:3306 -d mysql/mysql-server:5.6`

### MySQL
* `docker exec -it mysql56 mysql -uroot -p`

* Run:

```
CREATE DATABASE mydb;

use mydb;

CREATE TABLE mycontent (
`id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
message varchar(255) not null,
content varchar(255) not null
)ENGINE=InnoDB;

explain mycontent;
```

### RabbitMQ
* Browse to: `http://localhost:8080/` (after a few seconds)

* Create a Queue: `myqueue`

* Start the app: `go run amqp2sql.go`

* Publish Message on `myqueue`: `{"message": "Hello", "content": "Aww yiss!"}`

* On the DB cli, do: `select * from mycontent;`