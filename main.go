package main

import (
    "context"
    "fmt"
    "log"
    // "os"
    "time"

	"github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoInstance struct {
    Client *mongo.Client
    DB     *mongo.Database
}

type Task struct {
	// ID   primitive.ObjectID `bson:"_id"`
	Name          string `bson:"name"`
	Mobile_Number string `bson:"mobile_number"`
	User_type     string `bson:"user_type"`
}
var MI MongoInstance

func init() {

 client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(ctx, readpref.Primary())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Database connected!")

    MI = MongoInstance{
        Client: client,
        DB:     client.Database("bbc"),
    }
}

func GetPost(c *fiber.Ctx) error {
    
    catchphraseCollection := MI.DB.Collection("post")
    fmt.Println(catchphraseCollection);
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

    var catchphrase Task
    objId, err := primitive.ObjectIDFromHex(c.Params("id"))
    fmt.Println(objId)
    findResult := catchphraseCollection.FindOne(ctx, bson.D{{}})
    if err := findResult.Err(); err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Catchphrase Not found",
            "error":   err,
        })
    }

    err = findResult.Decode(&catchphrase)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Catchphrase Not found!!",
            "error":   err,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "data":    catchphrase,
        "success": true,
    })
}

func AddPost(c *fiber.Ctx) error {
    catchphraseCollection := MI.DB.Collection("post")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    catchphrase := new(Task)

    if err := c.BodyParser(catchphrase); err != nil {
        log.Println(err)
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "Failed to parse body",
            "error":   err,
        })
    }

    result, err := catchphraseCollection.InsertOne(ctx, catchphrase)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Catchphrase failed to insert",
            "error":   err,
        })
    }
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "data":    result,
        "success": true,
        "message": "Catchphrase inserted successfully",
    })

}


func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World 👋!")
    })
    app.Get("/post/:id", GetPost)
    app.Post("/addpost", AddPost)

    // config.ConnectDB()
    app.Listen(":3000")
}
