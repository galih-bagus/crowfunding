package routes

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()

	api := router.Group("/api/v1")

	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.Login)
	api.POST("/check-email", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaignDetail)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	return router
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}

// func authAdminMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		session := sessions.Default(c)

// 		userIDSession := session.Get("userID")

// 		if userIDSession == nil {
// 			c.Redirect(http.StatusFound, "/login")
// 			return
// 		}
// 	}
// }

// func loadTemplates(templatesDir string) multitemplate.Renderer {
// 	r := multitemplate.NewRenderer()

// 	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	includes, err := filepath.Glob(templatesDir + "/**/*")
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	for _, include := range includes {
// 		layoutCopy := make([]string, len(layouts))
// 		copy(layoutCopy, layouts)
// 		files := append(layoutCopy, include)
// 		r.AddFromFiles(filepath.Base(include), files...)
// 	}
// 	return r
// }
