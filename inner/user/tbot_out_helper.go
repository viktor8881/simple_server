package user

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
	"simpleserver/generated"
)

func InListUserByEmail(payload string, in interface{}) error {
	inDto, ok := in.(*generated.ListUserByEmailRequest)
	if !ok {
		return fmt.Errorf("unexpected type of in: %T", in)
	}

	inDto.Email = payload
	return nil
}

func OutListUserHelper(logger *zap.Logger) func(c telebot.Context, in any) error {
	return func(c telebot.Context, in any) error {
		outDto, ok := in.(*generated.ListUserResponse)
		if !ok {
			logger.Error("error: unexpected type of in: %T\n", zap.Any("in", in))
			return c.Send("Invalid data format")
		}

		var str string
		for _, user := range outDto.Users {
			str += fmt.Sprintf("ID: %s, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
		}

		if str == "" {
			str = "No users found"
		}

		return c.Send(str)
	}
}
