package scrapy

import "testing"

var (
	// List supported types and values for stats collector
	values = map[string]interface{}{
		"integer": 0,
		"string":  "string",
	}
)

func TestNewSpiderStats(t *testing.T) {
	st := NewStats()

	if st.mutex == nil {
		t.Error("Error on create stats collector, mutex can not be empty or not initialized")
	}
}

func TestSpiderStats_String(t *testing.T) {
	st := NewStats()
	st.values = values

	v := st.String()
	if v != `{"integer":0,"string":"string"}` {
		t.Error(
			"Stats string representation incorrect",
			"expected", `{"integer":100,"string":"string"}`,
			"got", v,
		)
	}
}

func TestSpiderStats_SetValue(t *testing.T) {
	st := NewStats()

	for k, v := range values {
		st.SetValue(k, v)
	}

	for k, v := range values {
		if val, ok := st.values[k]; !ok || val != v {
			t.Error(
				"Incorrect value in stats collector for key", k,
				"expected", v,
				"got", val,
			)
		}
	}
}

func TestSpiderStats_GetValue(t *testing.T) {
	st := NewStats()
	st.values = values

	for k, v := range values {
		if val, ok := st.GetValue(k); !ok || val != v {
			t.Error(
				"Incorrect or does not exists value in stats collector for key", k,
				"expected", v,
				"got", val,
			)
		}
	}
}

func TestSpiderStats_Clear(t *testing.T) {
	st := NewStats()
	st.values = values

	if st.values == nil {
		t.Error("Empty values in stats collector")
	}

	st.Clear()

	if len(st.values) != 0 {
		t.Error(
			"Not empty values in stats collector",
			"expected", nil,
			"got", st.values,
		)
	}
}

func TestSpiderStats_IncValue(t *testing.T) {
	st := NewStats()
	st.values = values

	for i := range []int{1, 2, 3} {
		st.IncValue("integer")

		if v, ok := st.values["integer"]; !ok || v != i+1 {
			t.Error(
				"Wrong increment value by key",
				"expected", i+1,
				"got", v,
			)
		}
	}

	err := st.IncValue("string")
	if err == nil {
		t.Error("Can not increment string, function dont return any type of error")
	}
}
