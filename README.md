# Order service

Order service provides API for placing, taking and retrieving order. 


## Getting Started

This instruction will get you a copy of the project up and running on your local machine.

### Prerequisites

* [docker](https://docs.docker.com/engine/installation/)
* [docker compose](https://docs.docker.com/compose/install/)

### Installing

#### 1. clone repo

```
git clone https://github.com/MinggaoYin/order-service
```

#### 2. cd to project root dir

```
cd order-service
```

#### 3. Update GOOGLE_API_KEY in .env file under project root directory

```
GOOGLE_API_KEY=XXXXXXXXXXXXXXXXXXXXXXX
```

#### 4. Run start.sh to build and run container

```
./start.sh
```

#### 5. Open another terminal and ping server

```
curl localhost:8080/

{"status":"OK"}
```

#### 6. Import Postman collection and play around with the API

```
postman/order-service.postman_collection.json
```