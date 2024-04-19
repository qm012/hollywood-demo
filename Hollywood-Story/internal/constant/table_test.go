package constant

import "testing"

func TestTableRandomAttributes(t *testing.T) {
	for i := 0; i < 1; i++ {
		table := TableRandomAttributes()
		t.Log("\n", TableRandomAttributes().Text)
		attributeValues := table.AttributeValueByRandom(4)
		for _, value := range attributeValues {
			t.Log(value)
		}
	}
}

func TestGetValuesByElementMap(t *testing.T) {
	values := getValuesByElementMap()
	t.Log(values)
	t.Log("len(values)=", len(values))
}
