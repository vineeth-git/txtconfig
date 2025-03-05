package txtconfig
import "fmt"
import "reflect"
import "regexp"
import "strings"
import "os"
import "strconv"
import "bufio"
import "errors"
var (
	ErrTargetNotPtr = errors.New("Target must be a struct pointer")
	ErrDataLoad = errors.New("Error loading data from file")
)

func camelToSnake(s string) string {
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}
func newTxtConfig(filename string) (map[string]string, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	data := make(map[string]string)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if !(strings.HasPrefix(line, "#")) && strings.Contains(line, "=") {
			pair := strings.Split(line, "=")
			pair[0] = strings.Trim(pair[0], " \t")
			pair[1] = strings.Trim(pair[1], " \t")
			data[pair[0]] = pair[1]
		}
	}
	return data, nil
}

func Load(filename string, def interface{}) error {
	s := reflect.ValueOf(def)
	if s.Kind() != reflect.Ptr {
		return ErrTargetNotPtr
	}
	cfg, err := newTxtConfig(filename)
	if err != nil {
		return ErrDataLoad
	}
	element := s.Elem()
	defType := element.Type()
	for i := 0; i < element.NumField(); i++ {
		field := element.Field(i)
		fieldMeta := defType.Field(i)
		key := fieldMeta.Tag.Get("key")
		value, ok := "", false
		if value, ok = cfg[key]; !ok {
			if value, ok = cfg[camelToSnake(fieldMeta.Name)]; !ok {
				value, ok = fieldMeta.Tag.Lookup("default")
			}
		}
		required := (fieldMeta.Tag.Get("required") == "true")
		if string(value) == "" && required {
			return errors.New(fmt.Sprint("Required field is empty: ", fieldMeta.Name))
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(value, 0, field.Type().Bits())
			if err != nil {
				return err
			}
			field.SetInt(val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(value, 0, field.Type().Bits())
			if err != nil {
				return err
			}
			field.SetUint(val)
		case reflect.Bool:
			val, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			field.SetBool(val)
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(value, field.Type().Bits())
			if err != nil {
				return err
			}
			field.SetFloat(val)
		}		
	}
	return nil
}