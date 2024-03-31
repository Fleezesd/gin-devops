package models

func CasbinAddPolicies(rules [][]string) (bool, error) {
	return Enforcer.AddPolicies(rules)
}
func CasbinAddOnePolicy(sub, obj, act string) (bool, error) {
	return Enforcer.AddPolicy(sub, obj, act)
}
