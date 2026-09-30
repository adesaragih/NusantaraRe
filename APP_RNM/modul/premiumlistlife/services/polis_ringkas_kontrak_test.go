package services

// Kontrak `GET /api/polis-life/ringkas` - pl4, tiket 08. SISI BACKEND.
//
// ⛔ Pelajaran envelope `galat` (handlers/envelopegalat_test.go): backend
// menulis satu kunci, klien membaca kunci lain, dan tidak satu pun uji di
// satu sisi mana pun gagal. Kontrak lintas modul PremiumList → Claim Life
// karena itu dikunci di KEDUA sisi: di sini (Go membaca antarmuka
// TypeScript), dan di `frontend/src/modul/claimlife/policydatalife.kontrak.test.ts`
// (TypeScript membaca tag JSON Go).

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// kunciJSONStruct membaca kunci JSON sebuah struct lewat MARSHAL - yang
// benar-benar dikirim, bukan yang tertulis di tag.
func kunciJSONStruct(t *testing.T, v any) []string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	var k []string
	for x := range m {
		k = append(k, x)
	}
	sort.Strings(k)
	return k
}

// kunciAntarmukaTS membaca nama medan satu `export interface` TypeScript.
func kunciAntarmukaTS(t *testing.T, sumber, nama string) []string {
	t.Helper()
	i := strings.Index(sumber, "export interface "+nama+" {")
	if i < 0 {
		t.Fatalf("antarmuka %s tidak ditemukan di klien", nama)
	}
	j := strings.Index(sumber[i:], "\n}")
	blok := sumber[i : i+j]
	// Medan tingkat pertama saja: tepat dua spasi indentasi.
	pola := regexp.MustCompile(`(?m)^  ([A-Za-z][A-Za-z0-9]*)\??:`)
	var k []string
	for _, m := range pola.FindAllStringSubmatch(blok, -1) {
		k = append(k, m[1])
	}
	sort.Strings(k)
	return k
}

func TestPolicyDataLifeSamaDiKeduaSisi(t *testing.T) {
	isi, err := os.ReadFile(filepath.Join("..", "..", "..", "frontend", "src", "modul", "claimlife", "api.ts"))
	if err != nil {
		t.Fatalf("membaca klien: %v", err)
	}
	klien := kunciAntarmukaTS(t, string(isi), "PolicyDataLife")
	// `MedanTanpaSumber` diisi supaya slice nil tidak terkirim sebagai null
	// dan tetap tercacah sebagai kunci.
	server := kunciJSONStruct(t, PolicyDataLife{MedanTanpaSumber: []string{}})
	if len(server) < 10 {
		t.Fatalf("hanya %d kunci server terbaca; pembacanya yang rusak", len(server))
	}
	if !reflect.DeepEqual(klien, server) {
		t.Errorf("kunci PolicyDataLife berselisih:\n  server %v\n  klien  %v", server, klien)
	}
}
