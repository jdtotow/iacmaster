package models

import "gorm.io/gorm"

type EnvironmentVariable struct {
	gorm.Model
	Name  string
	Value string
}

func CreateEnvironmentVariable(name, value string) EnvironmentVariable {
	return EnvironmentVariable{
		Name:  name,
		Value: value,
	}
}
func (env EnvironmentVariable) GetName() string {
	return env.Name
}
func (env EnvironmentVariable) GetValue() string {
	return env.Value
}
func (env *EnvironmentVariable) SetName(name string) {
	env.Name = name
}
func (env *EnvironmentVariable) SetValue(value string) {
	env.Value = value
}
