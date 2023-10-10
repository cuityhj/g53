package g53

import (
	"testing"

	"github.com/cuityhj/g53/util"
)

func parseMatchRender(t *testing.T, rawData string) {
	wire, _ := util.HexStrToBytes(rawData)
	buf := util.NewInputBuffer(wire)
	nm, err := MessageFromWire(buf)
	Assert(t, err == nil, "err should be nil")

	render := NewMsgRender()
	nm.Rend(render)
	WireMatch(t, wire, render.Data())
}

func TestRdataFromToWire(t *testing.T) {
	rawDatas := []string{
		//a
		"04b08500000100010001000103616161066e69757a756f036f72670000010001c00c0001000100000e10000403030303c0100002000100000e10001404636e7331097a646e73636c6f7564036e6574000000291000000000000000",

		//aaaa
		"04b08500000100010001000103626262066e69757a756f036f726700001c0001c00c001c000100000e10001024018d00000400000000000000000001c0100002000100000e10001404636e7331097a646e73636c6f7564036e6574000000291000000000000000",

		//dname
		"04b08500000100010001000103646464066e69757a756f036f72670000270001c00c0027000100000e1000100377777706676f6f676c6503636f6d00c0100002000100000e10001404636e7331097a646e73636c6f7564036e6574000000291000000000000000",

		//ptr
		"04b08500000100010001000103696969066e69757a756f036f726700000c0001c00c000c000100000e100016013101310131013107696e2d61646472046172706100c0100002000100000e10001404636e7331097a646e73636c6f7564036e6574000000291000000000000000",

		//mx
		"04b08500000100010001000103676767066e69757a756f036f726700000f0001c00c000f000100000e10000f000a037777770331363303636f6d00c0100002000100000e10001404636e7331097a646e73636c6f7564036e6574000000291000000000000000",

		//srv
		"04b08500000100010001000103686868066e69757a756f036f72670000210001c00c0021000100000e1000200001000000090d73797361646d696e732d626f78066e69757a756f036f726700c0400002000100000e10001404636e7331097a646e73636c6f7564036e6574000000291000000000000000",

		//soa
		"04b085000001000100010001066e69757a756f036f72670000060001c00c0006000100000e10003704636e7331097a646e73636c6f7564036e657400046d61696c046b6e657402636ec00c0000000100000e1000000e1000000e1000000e10c00c0002000100000e100002c0280000291000000000000000",

		//cname
		"04b08500000100010001000103636363066e69757a756f036f72670000050001c00c0005000100000e10000f0377777705626169647503636f6d00c0100002000100000e10001404636e7331097a646e73636c6f7564036e6574000000291000000000000000",
		//txt
		"04b08500000100010001000103656565066e69757a756f036f72670000100001c00c0010000100000e10001302446f03796f750477616e7402746f03646965c0100002000100000e10001404636e7331097a646e73636c6f7564036e6574000000291000000000000000",

		//spf
		"04b08500000100010001000103666666066e69757a756f036f72670000630001c00c0063000100000e10000f01490477616e7402746f046c697665c0100002000100000e10001404636e7331097a646e73636c6f7564036e6574000000291000000000000000",
	}

	for _, raw := range rawDatas {
		parseMatchRender(t, raw)
	}
}

func TestTxtParse(t *testing.T) {
	txts := []string{
		"\"good boy\" \"bad boy\"",
		"\"good \\\"boy\" \"bad boy\"",
		"   \"good\" \"v=1 boy\"  ",
		"\"good",
		"good boy",
		"good     boy",
		"\"good boy\"",
		"good \"boy\"",
	}
	type expectResult struct {
		err  error
		strs []string
	}
	expects := []expectResult{
		expectResult{
			err:  nil,
			strs: []string{"good boy", "bad boy"},
		},

		expectResult{
			err:  nil,
			strs: []string{"good \\\"boy", "bad boy"},
		},

		expectResult{
			err:  nil,
			strs: []string{"good", "v=1 boy"},
		},

		expectResult{
			err:  ErrQuoteInTxtIsNotInPair,
			strs: nil,
		},

		expectResult{
			err:  nil,
			strs: []string{"good", "boy"},
		},

		expectResult{
			err:  nil,
			strs: []string{"good", "boy"},
		},

		expectResult{
			err:  nil,
			strs: []string{"good boy"},
		},

		expectResult{
			err:  nil,
			strs: []string{"good", "\\\"boy\\\""},
		},
	}

	for i, txt := range txts {
		ss, err := txtStringParse(txt)
		Assert(t, err == expects[i].err, "")
		if expects[i].strs == nil {
			Assert(t, ss == nil, "")
		} else {
			Assert(t, len(ss) == len(expects[i].strs), "")
			for j, s := range ss {
				Assert(t, s == expects[i].strs[j], "")
			}
		}
	}
}
