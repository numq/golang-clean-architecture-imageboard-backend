# Imageboard Clean Architecture Back-End

*This is the part of project which includes back-end and mobile client.* 

**Also check:**
[Android Compose Clean Architecture Client](https://github.com/numq/android-clean-architecture-imageboard-client)

## Project description:
  **Simple anonymous imageboard which allows to create threads, posts and reply to users by quoting their posts.**

## Techs:

*Common:*
- **Golang**
- **Docker**
- **Yaml**

*Data:*
- **Protobuf**
- **gRPC**
- **MongoDB**

## Structure:
- `Config`
- `Data`\
*Repository implementation, data sources initialization.*
- `Domain`\
*Abstract entities & repositories.*
- `Infrastructure`\
*Server & handlers.*
- `Mapping`\
*In fact, creating a contract between models and entities.*\
*It's very important, therefore it's a separate package.*
- `Proto`
- `UseCase`\
*Most likely: use case for each repository method.*

### Summary:
**Well, I use `Golang` when I need to get simple and obvious product. I really like it, because I can create what I need without tons of libs and avoiding extra dependencies. Also Golang is great to fit with *Clean Architecture* because of abstraction through composition. I think that the stucture of this project is pretty good and I will use it again.**

**I like `MongoDB` because it's flexible and not very complicated to compose effective database query inside repository. Also it doesn't need to create a special abstraction on top of database methods which is definitely an advantage.**

**What about `gRPC`? I think that it's good to be used as a CRUD alternative to REST, but I would prefer Web Sockets for real-time (I mean "actually real-time", e.g. chatting) operations to avoid polling strategy, although gRPC is good for it, but I just think that each task needs it's own tool.**

**`Protobuf` is good to use for crossplatform products, but it's a bit tricky to set up the projects, because the implementations are different (especially in the case of the JVM). Sometimes it matters.**

**`Docker` is simply the best way to deploy an application, especially if you want to use database replication (which is really painful without it).**

**`Yaml` is just the simple way to create separate config file.**
