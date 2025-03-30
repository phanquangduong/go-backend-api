package basic

import (
	"testing"
)

func TestAddOne(t *testing.T) {
	var (
		input  = 1
		output = 3
	)

	actual := AddOne(input)
	if actual != output {
		t.Errorf("AddOne(%d), output %d, actual = %d ", input, output, actual)
	}

	// assert.Equal(t, AddOne(2), 3, "AddOne(2) should be 3")
	// assert.NotEqual(t, 2, 3)
	// assert.Nil(t, nil, nil)
}

// func TestRequire(t *testing.T) {
// 	require.Equal(t, 2, 3)
// 	fmt.Println("Not executing")
// }

// func TestAssert(t *testing.T) {
// 	assert.Equal(t, 2, 3)
// 	fmt.Println("Executing")
// }
