package router

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/aikwen/aifriend-go/internal/api/handler"
	"github.com/aikwen/aifriend-go/internal/api/mw"
)

// SetupRouter
func SetupRouter(h *handler.Handler, accessSecret string, env string) *gin.Engine {
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())

	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered any) {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
	}))


	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true

	corsConfig.AllowOrigins = []string{
		"http://localhost:5173", // Vue Vite 默认端口
		"http://127.0.0.1:5173",
	}

	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(corsConfig))
	r.Static("/media", "./media")


	if env == "dev" {
		r.LoadHTMLFiles("web/index.html")
		r.Static("/assets", "./web/assets")
		r.StaticFile("/favicon.ico", "./web/favicon.ico")

		r.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})

		r.NoRoute(func(c *gin.Context) {
			if c.Request.Method == "GET" && !strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.HTML(http.StatusOK, "index.html", nil)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"result": "接口不存在"})
			}
		})
	} else {

		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "AIFriend API is running"})
		})

		r.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"result": "接口不存在"})
		})
	}

	accountGroup := r.Group("/api/user/account")
	{
		// 公开接口区
		accountGroup.POST("/login/", h.Login)
		accountGroup.POST("/register/", h.Register)
		accountGroup.POST("/refresh_token/", h.Refresh)

		// protect 接口区
		authRequiredGroup := accountGroup.Group("/")
		authRequiredGroup.Use(mw.JWTAuthMiddleware(accessSecret))
		{
			authRequiredGroup.POST("/logout/", h.Logout)
			authRequiredGroup.GET("/get_user_info/", h.GetUserInfo)
		}
	}

	characterGroup := r.Group("/api/create/character")
	{
		characterGroup.GET("/get_list/", h.GetCharacterList)

		// protect 接口区域
		authRequiredGroup := characterGroup.Group("/")
		authRequiredGroup.Use(mw.JWTAuthMiddleware(accessSecret))
		{
			authRequiredGroup.POST("/create/", h.CreateCharacter)
			authRequiredGroup.GET("/get_single/", h.GetCharacter)
			authRequiredGroup.POST("/update/", h.UpdateCharacter)
			authRequiredGroup.POST("/remove/", h.DeleteCharacter)
		}
	}

	profileGroup := r.Group("/api/user/profile")
	profileGroup.Use(mw.JWTAuthMiddleware(accessSecret))
	{
		profileGroup.POST("/update/", h.UpdateUserInfo)
	}


	homepageGroup := r.Group("/api/homepage")
	{
		homepageGroup.GET("/index/", h.SearchCharacters)
	}

	return r
}