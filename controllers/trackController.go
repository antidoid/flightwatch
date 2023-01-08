package controllers

import (
	"github.com/antidoid/flightwatch/models"

	"github.com/gofiber/fiber/v2"
)

func CreateTrack(c *fiber.Ctx) error {
    c.Accepts("application/json")

    // Get data off the request body
    var track models.Track
    err := c.BodyParser(&track)
    if err != nil {
        return c.Status(fiber.StatusOK).JSON(fiber.Map{
            "message": "Error parsing incoming JSON " + err.Error(),
        })
    }

    // Create a track
    err = models.CreateTrack(track)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error creating a track in db, " + err.Error(),
        })
    }

    // Return the track
    return c.Status(fiber.StatusOK).JSON(track)
}

func GetTracks(c *fiber.Ctx) error {
    // Find all the tracks
    tracks, err := models.GetTracks()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error getting all the tracks, " + err.Error(),
        })
    }

    // return them as JSON array
    return c.Status(fiber.StatusOK).JSON(tracks)
}

func GetTrack(c *fiber.Ctx) error {
    // Find the track
    id := c.Params("id")
    track, err := models.GetTrack(id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error finding that track, " + err.Error(),
        })
    }

    // Respond with it
    return c.Status(fiber.StatusOK).JSON(track)
}

func UpdateTrack(c *fiber.Ctx) error {
    c.Accepts("application/json")

    // Find the track
    id := c.Params("id")

    track, err := models.GetTrack(id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error finding that track, " + err.Error(),
        })
    }

    // Get data of request body
    var newTrack models.Track
    err = c.BodyParser(&newTrack)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error parsing incoming JSON " + err.Error(),
        })
    }

    // Update it
    err = models.UpdateTrack(track, newTrack)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error updating the track in db " + err.Error(),
        })
    }

    // return the updated track
    updatedTrack, _ := models.GetTrack(id)
    return c.Status(fiber.StatusOK).JSON(updatedTrack)
}

func DeleteTrack(c *fiber.Ctx) error {
    // Find the track
    id := c.Params("id")

    track, err := models.GetTrack(id)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error finding that track, " + err.Error(),
        })
    }

    // Delete the track
    err = models.DeleteTrack(track)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error deleting the track in db,  " + err.Error(),
        })
    }

    // return the success message
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Successfully deleted track",
    })
}

