package service

import (
	"context"

	"clean-arch/app/model/postgre"
	"clean-arch/app/repository/postgre"
	"github.com/gofiber/fiber/v2"
)

func CreateRoleService(c *fiber.Ctx) error {
	var body struct {
		Name string `json:"name"`
		Desc string `json:"description"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	r := &postgre.Role{Name: body.Name, Desc: body.Desc}
	if err := repository.CreateRole(context.Background(), r); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"roleId": r.ID})
}

func GetRoleByNameService(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"name required"}) }
	r, err := repository.GetRoleByName(context.Background(), name)
	if err != nil { return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()}) }
	if r == nil { return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error":"not found"}) }
	return c.JSON(r)
}
