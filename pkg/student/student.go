package student

type Student struct {
	name  string
	age   int
	grade int
}

func NewStudent() *Student {
	return new(Student)
}

func (std *Student) GetName() string {
	return std.name
}

func (std *Student) GetAge() int {
	return std.age
}

func (std *Student) GetGrade() int {
	return std.grade
}

func (std *Student) PutName(name string) {
	std.name = name
}

func (std *Student) PutAge(age int) {
	std.age = age
}

func (std *Student) PutGrade(grade int) {
	std.grade = grade
}
