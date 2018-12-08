# gae-graphql

## Installation

### Create Datastore indexes
`gcloud datastore indexes create index.yaml`

### Deploy App Engine
`gcloud app deploy`

## How to use

### 1. Creating users
POST
```
mutation {
  bob: createUser(name: "Taro", email: "taro@gmail.com") {
    id
    name
  }
  mark: createUser(name: "Jiro", email: "jiro@gmail.com") {
    id
    name
  }
}
```

### 2. Creating blogs
POST
```
mutation {
  a: createBlog(userId: "5153049148391424", title: "title1", content: "Hello World!") {
    id
    title
    content
  }
  b: createBlog(userId: "5153049148391424", title: "title2", content: "good morning") {
    id
    title
    content
  }
  c: createBlog(userId: "5635703144710144", title: "title3", content: "GraphQL is nice") {
    id
    title
    content
  }
}
```

### 3. Query posts with limit and offset
POST
```
{
  blogs(limit: 1, offset: 1) {
    totalCount
    nodes {
      id
      title
      content
      createdAt
    }
  }
}
```

### 4. Query user by user id
POST
```
{
  user(id: "5715161717407744") {
    name
    email
    blogs {
      totalCount
      nodes {
        title
        content
      }
    }
  }
}

```