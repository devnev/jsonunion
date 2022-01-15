package jsonunion_test

type Action interface {
	isAction()
}

type HelloAction struct {
	Target string `json:"target"`
}

type GoodbyeAction struct {
	UntilWhen string `json:"untilWhen"`
}

func (*HelloAction) isAction()   {}
func (*GoodbyeAction) isAction() {}
