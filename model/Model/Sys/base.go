package Sys

import (
	"GBT/model/Model/BaseModel"
	"fmt"
)

func TableName(name string) string {
	return fmt.Sprintf("%s%s%s", BaseModel.GetTablePrefix(),"sys_", name)
}