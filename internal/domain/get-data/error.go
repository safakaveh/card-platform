package getdata

import "errors"

var (
	ErrOrderNameRequired = errors.New("نام سفارش الزامی است")
	ErrOrderNotFound     = errors.New("سفارش پیدا نشد")
	ErrCardNotFound      = errors.New("ردیف کارت پیدا نشد")
)
