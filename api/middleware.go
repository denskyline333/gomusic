package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/denskyline333/gomusic/helper"
)

func JWTAuthMiddleware(ctx *fiber.Ctx) error {
	authorizationHeader := ctx.Get("Authorization")
	if authorizationHeader == "" {
		ctx.Status(401).JSON(fiber.Map{
			"ok":    false,
			"error": helper.NoTokenError.Error(),
		})
		return helper.NoTokenError
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 {
		ctx.Status(401).JSON(fiber.Map{
			"ok":    false,
			"error": helper.InvalidTokenError.Error(),
		})
		return helper.InvalidTokenError
	}

	if headerParts[0] != "Bearer" {
		ctx.Status(401).JSON(fiber.Map{
			"ok":    false,
			"error": helper.InvalidTokenError.Error(),
		})
		return helper.InvalidTokenError
	}

	claims, err := helper.ParseToken(headerParts[1])
	if err != nil {
		ctx.Status(401).JSON(fiber.Map{
			"ok":    false,
			"error": err.Error(),
		})
		return err
	}

	ctx.Context().SetUserValue("userID", claims["userId"])
	ctx.Next()
	return nil
}

func UserAudioMiddleware(ctx *fiber.Ctx) error {
	tokenUserID := ctx.Context().UserValue("userID").(string)
	paramsUserID := ctx.Params("userId")

	if (tokenUserID != paramsUserID) {
		return fmt.Errorf("You are not allowed to fetch/modify audio for this user")
	}

	ctx.Next()
	return nil
}
