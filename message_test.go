package g53

import (
	"github.com/cuityhj/g53/util"
	"testing"
)

func buildHeader(id uint16, setFlag []FlagField, counts []uint16, opcode Opcode, rcode Rcode) Header {
	h := Header{
		Id:      id,
		Opcode:  opcode,
		Rcode:   rcode,
		QDCount: counts[0],
		ANCount: counts[1],
		NSCount: counts[2],
		ARCount: counts[3],
	}

	for _, f := range setFlag {
		h.SetFlag(f, true)
	}

	return h
}

func matchMessageRaw(t *testing.T, rawData string, m *Message) {
	wire, _ := util.HexStrToBytes(rawData)
	buf := util.NewInputBuffer(wire)
	nm, err := MessageFromWire(buf)
	Assert(t, err == nil, "err should be nil")

	Equal(t, nm.Header, m.Header)
	matchQuestion(t, nm.Question, m.Question)
	matchSection(t, nm.GetSection(AnswerSection), m.GetSection(AnswerSection))
	matchSection(t, nm.GetSection(AuthSection), m.GetSection(AuthSection))
	matchSection(t, nm.GetSection(AdditionalSection), m.GetSection(AdditionalSection))

	render := NewMsgRender()
	nm.Rend(render)

	WireMatch(t, wire, render.Data())
}

func matchSection(t *testing.T, ns Section, s Section) {
	Equal(t, len(ns), len(s))
	for i := 0; i < len(ns); i++ {
		matchRRset(t, ns[i], s[i])
	}
}

func TestSimpleMessageFromToWire(t *testing.T) {
	qn, _ := NameFromString("test.example.com.")
	ra1, _ := AFromString("192.0.2.2")
	ra2, _ := AFromString("192.0.2.1")

	var answer Section
	answer = append(answer, &RRset{
		Name:   qn,
		Type:   RR_A,
		Class:  CLASS_IN,
		Ttl:    RRTTL(3600),
		Rdatas: []Rdata{ra1, ra2},
	})

	var authority Section
	ns, _ := NameFromString("example.com.")
	ra3, _ := NSFromString("ns1.example.com.")
	authority = append(authority, &RRset{
		Name:   ns,
		Type:   RR_NS,
		Class:  CLASS_IN,
		Ttl:    RRTTL(3600),
		Rdatas: []Rdata{ra3},
	})

	var additional Section
	glue, _ := NameFromString("ns1.example.com.")
	ra4, _ := AFromString("2.2.2.2")
	additional = append(additional, &RRset{
		Name:   glue,
		Type:   RR_A,
		Class:  CLASS_IN,
		Ttl:    RRTTL(3600),
		Rdatas: []Rdata{ra4},
	})

	matchMessageRaw(t, "04b0850000010002000100020474657374076578616d706c6503636f6d0000010001c00c0001000100000e100004c0000202c00c0001000100000e100004c0000201c0110002000100000e100006036e7331c011c04e0001000100000e100004020202020000291000000000000000", &Message{
		Header: buildHeader(uint16(1200), []FlagField{FLAG_QR, FLAG_AA, FLAG_RD}, []uint16{1, 2, 1, 2}, OP_QUERY, R_NOERROR),
		Question: &Question{
			Name:  qn,
			Type:  RR_A,
			Class: CLASS_IN,
		},
		Sections: [...]Section{answer, authority, additional},
		Edns: &EDNS{
			UdpSize:     uint16(4096),
			DnssecAware: false,
		},
	})
}

func benchmarkParseMessage(b *testing.B, raw string) {
	wire, _ := util.HexStrToBytes(raw)
	buf := util.NewInputBuffer(wire)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MessageFromWire(buf)
		buf.SetPosition(0)
	}
}

func BenchmarkParseKnetMessage(b *testing.B) {
	benchmarkParseMessage(b, "04b08180000100010004000d03777777046b6e657402636e0000010001c00c00010001000002580004caad0b0ac01000020001000000c1001404676e7331097a646e73636c6f7564036e657400c01000020001000000c10014046c6e7332097a646e73636c6f75640362697a00c01000020001000000c1001504676e7332097a646e73636c6f7564036e6574c015c01000020001000000c10015046c6e7331097a646e73636c6f756404696e666f00c039000100010000262c000401089801c0790001000100000599000401089901c09a00010001000007c800046f012189c09a00010001000007c8000477a7e9e9c09a00010001000007c80004b683170bc09a00010001000007c80004010865fdc09a001c0001000007c8001024018d00000400000000000000000001c0590001000100002fea000477a7e9ebc0590001000100002fea0004b683170cc0590001000100002fea0004010865fcc0590001000100002fea00046f01218ac059001c00010000249f001024018d000006000000000000000000010000291000000000000000")
}

func BenchmarkParseTestExample(b *testing.B) {
	benchmarkParseMessage(b, "04b0850000010002000100020474657374076578616d706c6503636f6d0000010001c00c0001000100000e100004c0000202c00c0001000100000e100004c0000201c0110002000100000e100006036e7331c011c04e0001000100000e100004020202020000291000000000000000")
}
