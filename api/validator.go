package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/surojuvenkatesh20/bank-mgmt/utils"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currencyVal, ok := fieldLevel.Field().Interface().(string); ok {
		return utils.IsSupportedCurrency(currencyVal)
	}
	return false
}
