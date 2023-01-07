package controllers

import (
	"github.com/antidoid/flightwatch/initializers"
	"github.com/antidoid/flightwatch/models"

	"github.com/gofiber/fiber/v2"
)

func CreateSearch(c *fiber.Ctx) error {
    // Get the data off the request body
    type Body struct {
        Origin string `json:"origin"`
        Destination string `json:"destination"`
        StartAt string `json:"startat"`
        EndAt string `json:"endat"`
        Contact string `json:"contact"`
        WayToContact string `json:"waytocontact"`
    }

    body := new(Body)
    if err := c.BodyParser(body); err != nil {
        return err
    }

    // Create a Search
    search := models.Search{
        Origin: body.Origin,
        Destination: body.Destination,
        StartAt: body.StartAt,
        EndAt: body.EndAt,
        Contact: body.Contact,
        WayToContact: body.WayToContact,
    }

    result := initializers.DB.Create(&search)
    if result == nil {
        c.Status(400)
        return nil
    }

    // Return the Search
    return c.JSON(fiber.Map{
        "search": search,
    })
}

func GetSearchs(c *fiber.Ctx) error {
    // Find all the searchs
    var searchs []models.Search
    initializers.DB.Find(&searchs)

    if len(searchs) == 0 {
        c.Status(404)
        return c.SendString("No Searchs available yet")
    }

    // return them as JSON array
    c.Status(200)
    return c.JSON(fiber.Map{
        "searchs": searchs,
    })
}

func GetSearch(c *fiber.Ctx) error {
    // Find the Search
    var search models.Search
    id := c.Params("id")
    initializers.DB.First(&search, id)

    if search.ID == 0 {
        c.Status(404)
        return c.SendString("Unable to find that search")
    }

    // Respond with it
    c.Status(200)
    return c.JSON(fiber.Map{
        "search": search,
    })
}

func UpdateSearch(c *fiber.Ctx) error {
    // Find the Search
    id := c.Params("id")

    var search models.Search
    initializers.DB.First(&search, id)

    if search.ID == 0 {
        c.Status(404)
        return c.SendString("Unable to find that search")
    }

    // Get data of request body
    type Body struct {
        Origin string `json:"origin"`
        Destination string `json:"destination"`
        StartAt string `json:"startat"`
        EndAt string `json:"endat"`
        Contact string `json:"contact"`
        WayToContact string `json:"waytocontact"`
    }

    body := new(Body)
    if err := c.BodyParser(body); err != nil {
        return err
    }

    // Update it
    initializers.DB.Model(&search).Updates(models.Search{
        Origin: body.Origin,
        Destination: body.Destination,
        StartAt: body.StartAt,
        EndAt: body.EndAt,
        Contact: body.Contact,
        WayToContact: body.WayToContact,
    })

    // return the updated search
    c.Status(200)
    return c.JSON(fiber.Map{
        "search": search,
    })
}

func DeleteSearch(c *fiber.Ctx) error {
    // Find the Search
    id := c.Params("id")

    var search models.Search
    initializers.DB.First(&search, id)

    if search.ID == 0 {
        c.Status(404)
        return c.SendString("Unable to find that search")
    }

    // Delete it
    initializers.DB.Delete(&models.Search{}, id)

    // return the deleted search
    c.Status(200)
    return c.JSON(fiber.Map{
        "search": search,
    })
}

