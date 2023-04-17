package combo

import (
	"genytest/pkg/possiblevalues"
	"genytest/utilities"
)

type Combos []*Combo

func (cs Combos) CloneCombos() (dst Combos) {
	for _, c := range cs {
		dst = append(dst, c.Clone())
	}
	return
}

func (cs Combos) GetSubCombosHash() string {
	allHash := ""
	for _, c := range cs {
		allHash = allHash + c.Hash
	}

	return utilities.Hash(allHash)
}

func (cs Combos) VerifyDuplicate(c *Combo) bool {
	for _, combo := range cs {
		if combo.Hash == c.Hash {
			return true
		}
	}
	return false
}

func (cs *Combos) addToAllItems(key, value string) {
	for _, combo := range *cs {
		combo.Items[key] = value
	}
}

func (cs Combos) combosCartesionProduct(allValues map[string][]string, keys []string) Combos {
	combos := Combos{}

	// create first values
	for _, pv := range allValues[keys[0]] {
		combo := Combo{}
		combo.Items = map[string]string{keys[0]: pv}
		combos = append(combos, &combo)
	}

	// cartesion product
	// add value of each key to the existing items
	for _, key := range keys[1:] {
		oldCombos := combos
		combos = Combos{}
		for _, value := range allValues[key] {
			newCombos := oldCombos.CloneCombos()
			newCombos.addToAllItems(key, value)
			combos = append(combos, newCombos...)
		}
	}

	return combos
}

func (Combos) GenCombos() (cs Combos) {
	allValues, keys := possiblevalues.GetAllValues()

	cs = cs.combosCartesionProduct(allValues, keys)

	for _, c := range cs {
		c.GenHash()
	}

	return
}
