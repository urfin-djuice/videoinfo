package main

import (
	"encoding/json"
	"github.com/zelenin/go-mediainfo"
	"testing"

	"gotest.tools/assert"
)

func equalJSON(t *testing.T, s1, s2 string) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		t.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		t.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	assert.DeepEqual(t, o1, o2)
}


func TestInformer_GetInfo(t *testing.T) {
	testTable := map[string]string{
		"./testfiles/small.mp4": `{"audio":[{"name":"mp4a-40-2","bit_rate":83051,"duration":"5.568s"}],"video":[{"name":"avc1","width":560,"height":320,"bit_rate":465642,"duration":"5.533s"}]}`,
		"./testfiles/small.webm": `{"audio":[{"name":"A_VORBIS","bit_rate":160000,"duration":"5.568s"}],"video":[{"name":"V_VP8","width":560,"height":320,"bit_rate":147762,"duration":"5.567s"}]}`,
	}

	for file, eq := range testTable {
		mi, err := mediainfo.Open(file)
		if err != nil {
			t.Errorf("Fail to open file %s", file)
		}
		informer := newInformer(mi)
		informerResult := informer.GetInfo()

		result, err := json.Marshal(informerResult)
		equalJSON(t, eq, string(result))
	}
}
