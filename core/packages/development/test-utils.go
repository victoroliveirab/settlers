//go:build test

package development

func (d *Instance) SetCardByIndex(index int, name string) {
	d.cards[index].Name = name
}
