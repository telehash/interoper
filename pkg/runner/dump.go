package runner

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func (report *Report) Dump() error {
	os.MkdirAll("interop/reports", 0755)

	data, err := json.Marshal(report)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join("interop/reports", report.SUT+"-"+time.Now().Format("20060102-1504")+".dump"), data, 0644)
}
