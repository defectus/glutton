package glutton

type NilNotifier struct{}

func (n *NilNotifier) Notify(*PayloadRecord) error {
	return nil
}

func (n *NilNotifier) Configure(*Settings) error {
	return nil
}
