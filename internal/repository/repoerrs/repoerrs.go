package repoerrs

import "errors"

var (
	ErrDataDeleted      = errors.New("данные помечены как удалённые")
	ErrNotEnoughBalance = errors.New("недостаточно средств")
)
