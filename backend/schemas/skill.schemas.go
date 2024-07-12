package skillschemas

type Skill struct {
	Key         string
	Name        string
	Description string
	Logo        string
	Tags        []string
}

type UpdateSkill struct {
	Name        string
	Description string
	Logo        string
	Tags        []string
}
