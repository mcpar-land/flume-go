package core_test

import (
	"flume-go/core"
	"testing"
)

var EngineJson string = `
{

}
`

func TestCompileBlueprint(t *testing.T) {
	bp, err := core.BlueprintFromJson(EngineJson)
	if err != nil {
		t.Error(err)
	}
	t.Log(bp)
}
