package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/tochukaso/golang-study/controller"
	"github.com/tochukaso/golang-study/env"
	"github.com/tochukaso/golang-study/graph"
	"github.com/tochukaso/golang-study/graph/generated"
	"github.com/tochukaso/golang-study/middleware"
	csrf "github.com/utrack/gin-csrf"
)

func main() {
	url := "https://www.asoview.com/purchase/scheduled-ticket/input/?ticketTypeCode=ticket0000009288&channelCode=EwWA0nCyCe"

	for i := 0; i < 1000; i++ {
		resp, err := http.Get(url)
		if err != nil {
			os.Exit(2)
		}
		defer resp.Body.Close()

		byteArray, _ := ioutil.ReadAll(resp.Body)
		s := string(byteArray)
		if strings.Index(s, "このページを再読み込みする") > 0 {
			fmt.Println("ng : ", i)
		} else {
			fmt.Println("next", s)
		}
	}

}

func startGin() {
	go runGraphQL()
	setLogger()
	engine := gin.Default()
	setCookiePolicy(engine)
	engine.Use(csrf.Middleware(csrf.Options{
		Secret: "token",
		ErrorFunc: func(c *gin.Context) {
			controller.RenderHTML(c, http.StatusForbidden, "token_error.tmpl", gin.H{})
			c.Abort()
		},
	}))

	engine.Use(middleware.RecordUaAndTime)
	engine.LoadHTMLGlob("templates/*.tmpl")

	addControllers(engine)

	loadTemplates(engine)
	engine.Run(":8081")
}

func runGraphQL() {
	port := "8082"

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func loadTemplates(engine *gin.Engine) {
	engine.Static("/static", "static")
}

func addControllers(engine *gin.Engine) {
	addLoginController(engine)
	addProductController(engine)
	addProductUploadController(engine)
	addUserController(engine)
	addMailController(engine)
	addGraphQL(engine)
}

func addLoginController(engine *gin.Engine) {
	engine.GET("/", controller.ShowLogin)
	engine.GET("/login/", controller.ShowLogin)
	engine.POST("/login/", controller.AttemptLogin)
	engine.GET("/logout/", controller.Logout)
}

func addProductController(engine *gin.Engine) {
	controller.InitProduct()
	group := engine.Group("/product")
	group.Use(controller.SessionCheck)
	group.GET("/", controller.ShowProducts)
	group.GET("/detail/:id", controller.GetProduct)
	group.GET("/download", controller.DownloadProduct)

	group.GET("/new", func(c *gin.Context) {
		controller.RenderHTML(c, http.StatusOK, "product_detail.tmpl", gin.H{})
	})
	group.POST("/", controller.PutProduct)
	group.POST("/delete", controller.DeleteProduct)

	apiGroup := engine.Group("/api/product")
	apiGroup.GET("/", controller.ShowProductsJSON)
}

func addProductUploadController(engine *gin.Engine) {
	controller.InitProduct()
	group := engine.Group("/product")
	group.Use(controller.SessionCheck)
	group.GET("/upload", func(c *gin.Context) {
		controller.RenderHTML(c, http.StatusOK, "product_upload.tmpl", gin.H{})
	})
	group.POST("/upload", controller.UploadProduct)
}

func addUserController(engine *gin.Engine) {
	controller.InitUser()
	group := engine.Group("/user")
	group.Use(controller.SessionCheck)
	group.GET("/", controller.ShowUsers)
	group.GET("/detail/:id", controller.GetUser)

	group.GET("/new", func(c *gin.Context) {
		controller.RenderHTML(c, http.StatusOK, "user_detail.tmpl", gin.H{})
	})
	group.POST("/", controller.PutUser)
	group.POST("/delete", controller.DeleteUser)
}

func addMailController(engine *gin.Engine) {
	controller.InitMailTemplate()
	group := engine.Group("/mail")
	group.Use(controller.SessionCheck)
	group.GET("/", controller.ShowMailTemplates)
	group.GET("/detail/:code", controller.GetMailTemplate)

	group.POST("/", controller.PutMailTemplate)
}

func addGraphQL(engine *gin.Engine) {
	engine.POST("/gq/query", graphqlHandler())
	engine.GET("/gq/", playgroundHandler())
}

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/gq/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func setLogger() {
	logFilePath := env.GetEnv().LogFilePath
	f, _ := os.Create(logFilePath)

	//gin.DefaultWriter = io.MultiWriter(f)
	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	log.SetOutput(f)
}

func setCookiePolicy(engine *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))
	store.Options(controller.MakeSessionOption())
	engine.Use(sessions.Sessions("sid", store))
}
