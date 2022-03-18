package json

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/lipence/config"
)

type testDataST struct {
	CommonField1 string      `json:"commonField1"`
	StructField1 testDataST2 `json:"structField1"`
}

type testDataST2 struct {
	CommonField2 string `json:"commonField2"`
}

var testDataObj = &testDataST{
	CommonField1: "data1",
	StructField1: testDataST2{
		CommonField2: "data2",
	},
}

var cfgPath string

const cfgFileName = "testData.json"

func TestMain(m *testing.M) {
	var err error
	if cfgPath, err = os.MkdirTemp("", "gotest_lipence_config-json"); err != nil {
		panic(fmt.Errorf("failed to create temp config dir: %w", err))
	}
	var testFileData []byte
	if testFileData, err = json.Marshal(testDataObj); err != nil {
		panic(fmt.Errorf("failed to create temp config json: %w", err))
	}
	defer func() {
		if err = os.RemoveAll(cfgPath); err != nil {
			fmt.Println(fmt.Errorf("failed to delete temp config json: %w, (path: %s)", err, cfgPath))
		}
	}()
	if err = os.WriteFile(filepath.Join(cfgPath, cfgFileName), testFileData, 0644); err != nil {
		panic(fmt.Errorf("failed to create temp config json: %w", err))
	}
	m.Run()
}

func TestJSON(t *testing.T) {
	if err := config.Use(&Loader{}); err != nil {
		t.Fatal(err)
	}
	if err := config.LoadConfigs(filepath.Join(cfgPath, cfgFileName)); err != nil {
		t.Fatal(err)
	}
	t.Run("decode", func(t *testing.T) {
		var data = testDataST{}
		if err := config.Root().Decode(&data); err != nil {
			t.Error("field not found")
		} else if !reflect.DeepEqual(&data, testDataObj) {
			t.Error("invalid field value")
		}
	})
	t.Run("rootField", func(t *testing.T) {
		if data, ok := config.Lookup("commonField1"); !ok {
			t.Error("field not found")
		} else if str, err := data.String(); err != nil {
			t.Error("invalid field type")
		} else if str != testDataObj.CommonField1 {
			t.Error("invalid field value")
		}
	})
	t.Run("subField", func(t *testing.T) {
		if data, ok := config.Lookup("structField1", "commonField2"); !ok {
			t.Error("field not found")
		} else if str, err := data.String(); err != nil {
			t.Error("invalid field type")
		} else if str != testDataObj.StructField1.CommonField2 {
			t.Error("invalid field value")
		}
	})
}
