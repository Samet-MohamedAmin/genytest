package combo

import (
	"genytest/pkg/possiblevalues"
	"genytest/utilities"
)

type Combos []*Combo

func (cs Combos) CloneCombos() (dst []*Combo) {
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

func createCombos(allValues map[string][]string, keys []string) (cs Combos) {
	for _, pv := range allValues[keys[0]] {
		combo := Combo{}
		combo.Items = map[string]string{
			keys[0]: pv,
		}
		cs = append(cs, &combo)
	}
	return
}

func (cs Combos) combosCartesionProduct(allValues map[string][]string, keys []string) Combos {
	combos := cs
	for _, key := range keys[1:] {
		oldCombos := combos
		combos = Combos{}
		for _, value := range allValues[key] {
			newCombos := oldCombos.CloneCombos()
			for _, combo := range newCombos {
				combo.Items[key] = value
			}
			combos = append(combos, newCombos...)
		}
	}

	return combos
}

func (Combos) GenCombos() (cs Combos) {
	allValues, keys := possiblevalues.GetAllValues()

	cs = createCombos(allValues, keys)

	cs = cs.combosCartesionProduct(allValues, keys)

	for _, c := range cs {
		c.GenHash()
	}

	return
}
