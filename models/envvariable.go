package models

import "gorm.io/gorm"

type EnvironmentVariable struct {
	gorm.Model
	Name  string
	Value interface{}
}

func CreateEnvironmentVariable(name string, value interface{}) EnvironmentVariable {
	return EnvironmentVariable{
		Name:  name,
		Value: value,
	}
}
func (env EnvironmentVariable) GetName() string {
	return env.Name
}
func (env EnvironmentVariable) GetValue() interface{} {
	return env.Value
}
func (env *EnvironmentVariable) SetName(name string) {
	env.Name = name
}
func (env *EnvironmentVariable) SetValue(value interface{}) {
	env.Value = value
}
