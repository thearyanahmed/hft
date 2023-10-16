# HelloFresh DevOps Test

## Running the application

Run the following application to run the program.

- `git clone git@github.com:hellofreshdevtests/HFtest-platform-engineering-thearyanahmed.git`
- `cd HFtest-platform-engineering-thearyanahmed && git checkout dev`
- `cp .env.example .env`
- `make start`

The application can run using just go as well, make sure to export the .env values and run `go run cmd/pkg/main.go`

## Build process

by calling `make build` , it will create  a `go binary` using [Dockerfile.prod](Dockerfile.prod). For development, [Dockerfile](Dockerfile) is used, which comes with a watcher out of the box, for any code changes.

The `Dockerfile.prod` uses scratch image, resutling in a smaller binary (_8.68MB_ avg).

**Running With Kubernetes**,

After running `make build`, simply calling `make deploy` will run the necessary files to serve the application.

For simplicity, the deployment uses a load balancer service component, which will expose an
external IP (`127.0.0.1` for local machine).

We can call our service through this ip eg:

```bash
curl http://localhost:8001/configs
```

## Architecture

The service comes as a package inside `/pkg` directory. The service is called through `/cmd/pkg/main.go`. `main.go` calls the necessary functions in a specific order to serve using this service.

The service (/pkg) follows a service repository pattern. Everything is exposed via interfaces, thus can be switched with alternatives just in case.

For this demonstration, persistant data storage was not used rather an [in memory datastore](pkg/repository/in_memory_repository.go) was used.

The router expects a struct that will satisfy the necessary interface and binds endpoints with handlers using that service.

The service on the other hand expects a datastore to satisfy the necessary interface.

The endpoints invokes handlers -> handlers calls the service -> the service calls -> repository.

The handlers calls different methods of the service layer based on it's needs and takes care of parsing/validating the request and sending a response.

Every endpoint returns a json response. The [presenter](pkg/presenter/presenter.go) setups somes common functions to respond with.

The [serializer](pkg/serializer/) validates the requests.
Note: `store_config_request` and `update_config_request` are both same for our usecase, `update_config_request` extends from `store_config_request`, reducing code duplication.

**kubernetes** takes the built docker image and uses a deployment and service components.

## Available endpoints

Postman collection can be found [here](docs/hellofresh-test-rest-api.postman_collection.json).

**List**

`/configs` (GET method) endpoint returns the list of configs, limited to 100, thinking of resource usage / packet size over network, memory in runtime etc.

```bash,
# request
curl --location 'localhost:8001/configs'

# response
[
    {
        "name": "Config_01",
        "metadata": {
            "allergens": {
                "eggs": "true",
                "nuts": "false",
                "seafood": "false"
            },
            "calories": "100",
            "carbohydrates": {
                "dietary-fiber": "4g",
                "sugars": "1g"
            },
            "fats": {
                "saturated-fat": "0g",
                "trans-fat": "1g"
            }
        }
    }
    ... {}
]

# will return [] (empty list if there are no configs stored).
```

**Create**

`/configs` (POST method) validates and data, checks if it exists, if not, then it'll
store the new config map in storage.

*Note* content type needs to be `application/x-www-form-urlencoded`.

```bash
## request
curl --location 'localhost:8001/configs' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'name=Config_01' \
--data-urlencode 'metadata={"calories": "100" ,"fats":{"saturated-fat":"0g","trans-fat":"1g"},"carbohydrates":{"dietary-fiber":"4g","sugars":"1g"},"allergens":{"seafood":"false","nuts": "false","eggs":"true"}}
'

# response
# returns the newly stored config
{
    "name": "Config_01",
    "metadata": {
        "allergens": {
            "eggs": "true",
            "nuts": "false",
            "seafood": "false"
        },
        "calories": "100",
        "carbohydrates": {
            "dietary-fiber": "4g",
            "sugars": "1g"
        },
        "fats": {
            "saturated-fat": "0g",
            "trans-fat": "1g"
        }
    }
}
```

**Update**

`/configs/{name}` (PUT method) updates a given config. At first, it looks up by `$name`, if it exists, it updates with given data.

*Note* content type needs to be `application/x-www-form-urlencoded`.

```bash
# request
curl --location --request PUT 'localhost:8001/configs/Config_01' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'name=Config_01' \
--data-urlencode 'metadata={"no-calories":12345,"fats":{"saturated-fat":"0g","trans-fat":"1g"},"carbohydrates":{"dietary-fiber":"4g","sugars":"1g"},"allergens":{"nuts":"false","seafood":"false","eggs":"true"}}
'

# response
{
    "name": "Config_01",
    "metadata": {
        "allergens": {
            "eggs": "true",
            "nuts": "false",
            "seafood": "false"
        },
        "carbohydrates": {
            "dietary-fiber": "4g",
            "sugars": "1g"
        },
        "fats": {
            "saturated-fat": "0g",
            "trans-fat": "1g"
        },
        "no-calories": 12345
    }
}
```

**Delete**

`/configs/{name}` (DELETE method) looks up an entry by `$name`, given it's exists, it deletes from the stored record list.

```bash
## request
curl --location --request DELETE 'localhost:8001/configs/Config_02' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data ''

## response
{
    "message": "config deleted successfully"
}
```

**Find**

`/configs/{name}` (GET method) looks up an entry by `$name`, and returns the given config map if found.

```bash
## request
curl --location 'localhost:8001/configs/Config_01'

## response
{
    "name": "Config_01",
    "metadata": {
        "allergens": {
            "eggs": "true",
            "nuts": "false",
            "seafood": "false"
        },
        "carbohydrates": {
            "dietary-fiber": "4g",
            "sugars": "1g"
        },
        "fats": {
            "saturated-fat": "0g",
            "trans-fat": "1g"
        },
        "no-calories": 12345
    }
}
```

**Search**

`/search` (GET method) takes the query parameter and returns a list (limited to 100) of configs.

```bash
curl --location 'localhost:8001/search?metadata.allergens.eggs=true&metadata.calories=100&name=Config_02'

## response
[
    {
        "name": "Config_02",
        "metadata": {
            "allergens": {
                "eggs": "true",
                "nuts": "false",
                "seafood": "false"
            },
            "calories": "100",
            "carbohydrates": {
                "dietary-fiber": "4g",
                "sugars": "1g"
            },
            "fats": {
                "saturated-fat": "0g",
                "trans-fat": "1g"
            }
        }
    }
]
```

## Common erorrs

- Validation errors will result in `BadRequest`
- Mapformed request body / service errors will result in `UnprocessableEntity`. Service level errors should be on the 500 series, but for simplicy, it is kept in that way.

Hello and thanks for taking the time to try this out.

The goal of this test is to assert (to some degree) your coding, testing, automation and documentation skills. You're given a simple problem, so you can focus on showcasing your techniques.

## Constrains

- The `name` of the config has been used an unique key. Attempting to create multiple config with same name will result in failure.
- Logging has been kept minimal.
- Kubernetes deployment does not contain readiness probe or startup probes.
- Nor does it contain maxSurge or maxUnavailable values (for simplicity).
- List limit has been hardcoded to 100. It can come from query parameter + config value.
- Simple happy and sad paths were covered for test cases.
- An nginx ingress.yaml has been added but not being used.

## Problem definition

The aim of test is to create a simple HTTP service that stores and returns configurations that satisfy certain conditions.
Since we love automating things, the service should be automatically deployed to kubernetes.

_Note: While we love open source here at HelloFresh, please do not create a public repo with your test in! This challenge is only shared with people interviewing, and for obvious reasons we'd like it to remain this way._

## Instructions

1. Clone this repository.
2. Create a new `dev` branch.
3. Solve the task and commit your code. Commit often, we like to see small commits that build up to the end result of your test, instead of one final commit with all the code.
4. Do a pull request from the `dev` branch to the `master` branch. More on that right below.
5. Reply to the thread you are having with our HR department so we can start reviewing your code.

In your pull request, make sure to write about your approach in the description. One or more of our engineers will then perform a code review.
We will ask questions which we expect you to be able to answer. Code review is an important part of our process;
this gives you as well as us a better understanding of what working together might be like.

We believe it will take 4 to 8 hours to develop this task, however, feel free to invest as much time as you want.

### Endpoints

Your application **MUST** conform to the following endpoint structure and return the HTTP status codes appropriate to each operation.

Following are the endpoints that should be implemented:

| Name   | Method      | URL
| ---    | ---         | ---
| List   | `GET`       | `/configs`
| Create | `POST`      | `/configs`
| Get    | `GET`       | `/configs/{name}`
| Update | `PUT/PATCH` | `/configs/{name}`
| Delete | `DELETE`    | `/configs/{name}`
| Query  | `GET`       | `/search?metadata.key=value`

#### Query

The query endpoint **MUST** return all configs that satisfy the query argument.

Query example-1:

```sh
curl http://config-service/search?metadata.monitoring.enabled=true
```

Response example:

```json
[
  {
    "name": "datacenter-1",
    "metadata": {
      "monitoring": {
        "enabled": "true"
      },
      "limits": {
        "cpu": {
          "enabled": "false",
          "value": "300m"
        }
      }
    }
  },
  {
    "name": "datacenter-2",
    "metadata": {
      "monitoring": {
        "enabled": "true"
      },
      "limits": {
        "cpu": {
          "enabled": "true",
          "value": "250m"
        }
      }
    }
  },
]
```

Query example-2:

```sh
curl http://config-service/search?metadata.key7.key10=true
```

Response example-2:

```json
[
  {
    "name": "burger-nutrition",
    "metadata": {
      "key0": 230,
      "key1": {
        "key2": "0g",
        "key3": "1g"
      },
      "key4": {
          "key5": "4g",
          "key6": "1g"
      },
      "key7": {
        "key9": "false",
        "key8": "false",
        "key10": "true"
      }
    }
  }
]
```

#### Schema

- **Config**
  - Name (string)
  - Metadata (nested key:value pairs where both key and value are strings of arbitrary length)

### Configuration

Your application **MUST** serve the API on the port defined by the environment variable `SERVE_PORT`.
The application **MUST** fail if the environment variable is not defined.

### Deployment

The application **MUST** be deployable on a kubernetes cluster. Please provide manifest files and a script that deploys the application on a minikube cluster.
The application **MUST** be accessible from outside the minikube cluster.

## Rules

- Applicants are requested to use either Python or GoLang.
- The API **MUST** return valid JSON and **MUST** follow the endpoints set out above.
- You **SHOULD** write testable code and demonstrate unit testing it.
- You can use any testing, mocking libraries provided that you state the reasoning and it's simple to install and run.
- You **SHOULD** document your code and scripts.
