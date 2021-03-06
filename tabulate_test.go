package tabulate

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type MyStruct struct {
	name   string
	amount int
}

var testData = []*MyStruct{
	&MyStruct{"Apple", 15},
	&MyStruct{"Orange", 1},
}

func TestNoFormat(t *testing.T) {
	layout := &Layout{Format: NoFormat}
	table, err := Tabulate(testData, layout)
	if err != nil {
		t.Fatal(err)
	}

	expecting := ("  name" + "amount\n" + // 6 + 6
		" Apple" + "    15\n" +
		"Orange" + "     1\n")
	if table != expecting {
		t.Fatalf("Expecting %q, got %q", expecting, table)
	}

	layout.HideHeaders = true
	table, err = Tabulate(testData, layout)
	if err != nil {
		t.Fatal(err)
	}

	expecting = (" Apple" + "15\n" + // 6 + 2
		"Orange" + " 1\n")

	if table != expecting {
		t.Fatalf("Expecting %q, got %q", expecting, table)
	}
}

func TestPlainFormat(t *testing.T) {
	layout := &Layout{Format: PlainFormat}
	table, err := Tabulate(testData, layout)
	if err != nil {
		t.Fatal(err)
	}

	expecting := ("  name" + " amount\n" + // 6 + 7
		" Apple" + "     15\n" +
		"Orange" + "      1\n")
	if table != expecting {
		t.Fatalf("Expecting %q, got %q", expecting, table)
	}

	layout.HideHeaders = true
	table, err = Tabulate(testData, layout)
	if err != nil {
		t.Fatal(err)
	}

	expecting = (" Apple" + " 15\n" + // 6 + 3
		"Orange" + "  1\n")

	if table != expecting {
		t.Fatalf("Expecting %q, got %q", expecting, table)
	}
}

func TestSimpleFormat(t *testing.T) {
	// The Simple Format is used by default
	table, err := Tabulate(testData, &Layout{})
	if err != nil {
		t.Fatal(err)
	}

	expecting := ("  name" + " amount\n" + // 6 + 7
		"------" + " ------\n" +
		" Apple" + "     15\n" +
		"Orange" + "      1\n")
	if table != expecting {
		t.Fatalf("Expecting %q, got %q", expecting, table)
	}

	table, err = Tabulate(testData, &Layout{Format: SimpleFormat})
	if err != nil {
		t.Fatal(err)
	}
	if table != expecting {
		t.Fatalf("Expecting %q, got %q", expecting, table)
	}
}

func TestSimpleFormatCustomHeaders(t *testing.T) {
	layout := &Layout{
		Format:  SimpleFormat,
		Headers: []string{"produce", "stuff"},
	}

	table, err := Tabulate(testData, layout)
	if err != nil {
		t.Fatal(err)
	}

	expecting := ("produce" + " stuff\n" + // 6 + 6
		"-------" + " -----\n" +
		"  Apple" + "    15\n" +
		" Orange" + "     1\n")
	if table != expecting {
		t.Fatalf("Expecting %q, got %q", expecting, table)
	}
}

func TestGridFormat(t *testing.T) {
	table, err := Tabulate(testData, &Layout{Format: GridFormat})
	if err != nil {
		t.Fatal(err)
	}

	expecting := ("+--------" + "+--------+\n" + // 9 + 10
		"|   name " + "| amount |\n" +
		"+========" + "+========+\n" +
		"|  Apple " + "|     15 |\n" +
		"+--------" + "+--------+\n" +
		"| Orange " + "|      1 |\n" +
		"+--------" + "+--------+\n")

	if table != expecting {
		t.Fatalf("Expecting %q, got %q", expecting, table)
	}
}

func TestFancyGridFormat(t *testing.T) {
	table, err := Tabulate(testData, &Layout{Format: FancyGridFormat})
	if err != nil {
		t.Fatal(err)
	}

	expecting := ("╒════════" + "╤════════╕\n" + // 9 + 10
		"│   name " + "│ amount │\n" +
		"╞════════" + "╪════════╡\n" +
		"│  Apple " + "│     15 │\n" +
		"├────────" + "┼────────┤\n" +
		"│ Orange " + "│      1 │\n" +
		"╘════════" + "╧════════╛\n")

	if table != expecting {
		t.Fatalf("Expecting %q, got %q", expecting, table)
	}
}

func TestPipeFormat(t *testing.T) {
	table, err := Tabulate(testData, &Layout{Format: PipeFormat})
	if err != nil {
		t.Fatal(err)
	}

	expecting := ("" +
		"  name " + "| amount\n" +
		"------ " + "| ------\n" +
		" Apple " + "|     15\n" +
		"Orange " + "|      1\n")

	if table != expecting {
		t.Fatalf("Expecting %q, got %q", expecting, table)
	}
}

type FullName struct {
	first string
	last  string
}

// Name implements Stringer interface
func (n *FullName) String() string {
	return n.first + " " + n.last
}

type MyBiggerStruct struct {
	Name        *FullName
	Amount      int
	Location    string
	Done        bool
	SurfaceArea float64
}

func TestComplexFormat(t *testing.T) {
	records := []*MyBiggerStruct{
		&MyBiggerStruct{&FullName{"Roy", "Smith"}, 15, "Washington D.C.", true,
			0.3453},
		&MyBiggerStruct{&FullName{"Fred", "Flanders"}, 100, "Montreal", false,
			1.0},
		&MyBiggerStruct{&FullName{"Bobby", "Smith"}, -2, "San Fransisco",
			false, 124353.23333333},
		&MyBiggerStruct{&FullName{"Jolene", "Lee"}, 234, "Guyene", true,
			11.0000000000001},
	}

	table, err := Tabulate(records, &Layout{Format: SimpleFormat})
	if err != nil {
		panic(err)
	}

	expecting := ("" +
		"         Name Amount        Location  Done          SurfaceArea\n" +
		"------------- ------ --------------- ----- --------------------\n" +
		"    Roy Smith     15 Washington D.C.  true      0.3453         \n" +
		"Fred Flanders    100        Montreal false      1              \n" +
		"  Bobby Smith     -2   San Fransisco false 124353.23333333     \n" +
		"   Jolene Lee    234          Guyene  true     11.0000000000001\n")

	if table != expecting {
		t.Fatalf("Expecting %q, got %q", expecting, table)
	}
}

func TestTabulateStringMatrix(t *testing.T) {
	records := [][]string{
		[]string{"here", "there"},
		[]string{"1", "2"},
	}

	layout := &Layout{Format: SimpleFormat}
	layout.Headers = []string{"a", "b"}

	table, err := Tabulate(records, layout)
	require.Nil(t, err)

	expecting := ("" +
		"   a" + "     b\n" + // 4 + 6
		"----" + " -----\n" +
		"here" + " there\n" +
		"   1" + "     2\n")
	assert.Equal(t, expecting, table)
}

func TestTabulateNoHeader(t *testing.T) {
	records := [][]string{
		[]string{"here", "there"},
		[]string{"1", "2"},
	}

	layout := &Layout{Format: SimpleFormat, HideHeaders: true}

	table, err := Tabulate(records, layout)
	require.Nil(t, err)

	expecting := ("" +
		"here" + " there\n" +
		"   1" + "     2\n")
	assert.Equal(t, expecting, table)
}
