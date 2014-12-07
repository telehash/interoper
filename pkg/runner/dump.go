package runner

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func (ctx *Context) Dump() error {
	os.MkdirAll("dumps", 0755)

	data, err := json.Marshal(ctx.Report)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join("dumps", ctx.SUT+"-"+time.Now().Format("20060102-1504")+".dump"), data, 0644)
}
