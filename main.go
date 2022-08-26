package main

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/routes"
	"bwastartup/transaction"
	"bwastartup/user"
)

func main() {
	db := helper.SetupDB()
	// Migrate Table From Entity
	db.AutoMigrate(&user.User{}, &campaign.Campaign{}, &campaign.CampaignImage{}, &transaction.Transaction{})

	router := routes.SetupRoutes(db)
	router.Run()
}
