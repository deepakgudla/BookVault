# GraphQL

- graphql is a query language and runtime for APIs that allows clients to request exactly the data they need

- unlike REST APIs where you have multiple endpoints, graphql uses a single endpoint with a flexible query syntax

- Benefits
    - single endpoint 
        - one url handles all requests
    - flexible queries
        - clients specify exactly what data they need
    - type safety
        - schema defines exact data structure
    - real time
        - has built in sub support
    - intospection
        - api is self documenting  

| REST | Graphql |
| ---- | ------- |
| has multiple endpoints like \n \users \profile | single endpoint \graphql |
| fixed response structure | client defines response shape |
| over or under-fetching common | get exactly what you ask for |
| multiple requests for related data | single request for nested data |

- resolvers
    - resolvers are like handlers or controllers in rest but has some key differences
    - resolvers are responsible for providing the data
    - most times you gonna use a single resolver to resolve all the fields
- resolvers work as such
    - define a schema
```json
type Todo {
    id: ID!
    text: String!
    done: Boolean! // ! means required field
}
```
- generate resolver interface from schema
    - for the above schema **glqgen** implements interface as following
```json
type QueryResolver interface {
    Todos(ctx context.Context) ([]*model.Todo, error)
    Todos(ctx context.Context, id string) (*model.Todo, error)

}
```
- implement
```json
query {
    user(id: 1) {
        name 
        posts {
            title 
            createdAt
        }
    }
}
```