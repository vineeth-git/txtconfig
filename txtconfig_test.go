package txtconfig
import "testing"

type Config struct {
	FirstName string `required:"true"`
	LastName string
	EmpId string `default:"XX"`
	Address string `key:"permanent_address"`
	country string `required:"true"`
}

func TestPassingNonPtr(t *testing.T) {
	c := Config{}
	err := Load("", c)
	if err == nil || err.Error() != "Target must be a struct pointer" {
		t.Errorf("Expecting error")
	}
}

func TestFileLoadError(t *testing.T) {
	c := Config{}
	err := Load("testdata/cf.txt", &c)
	if err == nil {
		t.Errorf("Expecting file load error")
	}
}

func TestStringField(t *testing.T) {
	c := Config{}
	err := Load("testdata/config.txt", &c)
	if err == nil {
		t.Error("Required field is not validated")
	}
	if c.FirstName != "John" {
		t.Error("First name is not set using camel case conversion")
	}
	if c.EmpId != "XX" {
		t.Error("Default value is not set")
	}
	if c.Address != "India" {
		t.Error("value is not picked up using key tag")
	}
}