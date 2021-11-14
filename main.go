package main

import (
	"RestApi/controllers" // new
	"RestApi/models"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)
  
  var identityKey = "id"

  type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}


func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"userID":   claims[identityKey],
			"name": user.(*models.User).UserName,
			"avatar": user.(*models.User).Avatar,
			"introduction": user.(*models.User).Introduction,
			"roles": user.(*models.User).Roles,
		},
	})
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
		fmt.Print("jancok g metu t  iki")
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH,OPTIONS,GET,PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func main() {
	r := gin.Default()

	models.ConnectDataBase()

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		// LoginResponse: func(c *gin.Context, code int, token string, t time.Time) {
		// 	c.JSON(code, gin.H{
		// 		"code": code,
		// 		"data": gin.H{
		// 			"token": token,
		// 		},
		// 	})
		// },
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
					"roles": v.Roles,
					"avatar": v.Avatar,
					"introduction": v.Introduction,
					"name": v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			aInterface := claims["roles"].([]interface{})
			aString := make([]string, len(aInterface))
			for i, v := range aInterface {
				aString[i] = v.(string)
			}
			return &models.User{
				UserName: claims[identityKey].(string),
				Roles: aString,
				Avatar: claims["avatar"].(string),
				Introduction: claims["introduction"].(string),
				Name: claims["name"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if (userID == "admin" && password == "admin123") {
				return &models.User{
					UserName:  userID,
					Roles:  []string{"admin"},
					Introduction: "Wu",
					Avatar: "Wu",
					Name: "Wu",
				}, nil
			}

			if (userID == "demo" && password == "demo123") {
				return &models.User{
					UserName:  userID,
					Roles:  []string{"demo"},
					Introduction: "Wu",
					Avatar: "Wu",
					Name: "Wu",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// if v, ok := data.(*models.User); ok && v.UserName == "admin" {
			// 	return true
			// }

			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "x-token"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
		  return true
		},
		MaxAge: 12 * time.Hour,
	}))

	r.Static("/admin", "./admin")
	r.POST("/login", authMiddleware.LoginHandler)
  
	r.GET("/resto", controllers.FindResto) // new
	r.GET("/resto/:id", controllers.FindRestoById) // new
	r.PATCH("/resto/:id", controllers.UpdateResto) // new
	r.POST("/resto", controllers.CreateResto)
	r.GET("/logscan", controllers.FindLogScan) // new

	qr := r.Group("/qr")
	qr.Use(authMiddleware.MiddlewareFunc()) 
	{
		qr.POST("/check", controllers.ScanQR)
	}
	r.GET("/", func(c *gin.Context) {
	  c.JSON(http.StatusOK, gin.H{"data": "hello world"})    
	})

	auth := r.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/info", helloHandler)
	}
  
	// r.Use(CORSMiddleware())


	


	r.Run()
  }