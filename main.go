package main

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/routes"
	"bwastartup/user"
)

func main() {
	db := helper.SetupDB()
	// Migrate Table From Entity
	db.AutoMigrate(&user.User{}, &campaign.Campaign{}, &campaign.CampaignImage{})

	router := routes.SetupRoutes(db)
	router.Run()
}
