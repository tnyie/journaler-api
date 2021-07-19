![](./banner.svg)

| GROUP   | METHOD | PATH           | JSON RESP? | DESCRIPTION                                       |
|---------|--------|----------------|------------|---------------------------------------------------|
| USERS   |        | /users         |            |                                                   |
|         | GET    | /users         | true       | Get current user information                      |
|         | PUT    | /users         | false      | First time user setup                             |
|         |        |                |            |                                                   |
| JOURNAL |        | /journals      |            |                                                   |
|         | GET    | /journals      | true       | Get all owned journals                            |
|         | GET    | /journals/{id} | true       | Get info on specific journal and children/entries |
|         | POST   | /journals      | true       | Create journal                                    |
|         |        |                |            |                                                   |
| Entry   |        | /entries       |            |                                                   |
|         | GET    | /entries       | -          | Not implemented                                   |
|         | GET    | /entries/{id}  | true       | Get info on specific entry                        |
|         | POST   | /entries       | true       | Create entry                                      |