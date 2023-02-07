package id_generator

type IdGenerator interface {
	Generate() (string, error)
}
