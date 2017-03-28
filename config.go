package lamp

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// LoadConfig load configuration by filename
func LoadConfig(filename string, v interface{}) error {
	filepath := path.Join(os.Getenv("LAMP_CONFIG_PATH"), filename)
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
