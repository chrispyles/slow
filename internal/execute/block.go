package execute

type Block []Expression

func (b Block) Execute(e *Environment) (Value, error) {
	var val Value
	var err error
	for _, s := range b {
		val, err = s.Execute(e)
		if err != nil {
			return nil, err
		}
	}
	return val, nil
}
