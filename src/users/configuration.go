package users

func Build(repository Repository) Facade {
	return &usersFacade{repository: repository}
}
