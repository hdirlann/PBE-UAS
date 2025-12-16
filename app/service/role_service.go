package service

import (
	"context"

	"clean-arch/app/model"
	"clean-arch/app/repository"
	"github.com/gofiber/fiber/v2"
)

// CreateRoleService
// @Summary Create role
// @Tags Roles
// @Description Create a new role.
// @Accept json
// @Produce json
// @Param body body object true "Role body" example({"name":"role-lab","description":"Peran lab"})
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /roles [post]
func CreateRoleService(c *fiber.Ctx) error {
	var body struct {
		Name string `json:"name"`
		Desc string `json:"description"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	r := &model.Role{Name: body.Name, Desc: body.Desc}
	if err := repository.CreateRole(context.Background(), r); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"roleId": r.ID})
}

// GetRoleByNameService
// @Summary Get role by name
// @Tags Roles
// @Description Get role by name.
// @Produce json
// @Param name path string true "Role name"
// @Success 200 {object} model.Role
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /roles/{name} [get]
func GetRoleByNameService(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name required"})
	}
	r, err := repository.GetRoleByName(context.Background(), name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if r == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(r)
}
