package handlers

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// polisTiruanH - PembacaPolis tiket 71 untuk uji HTTP.
type polisTiruanH struct{}

func (polisTiruanH) VersiPolis(_ context.Context, noPolis string) ([]models.VersiPolis, error) {
	if noPolis != "RNM-1" {
		return nil, repository.ErrPolisTidakAda
	}
	return []models.VersiPolis{{ProdKe: "0", CaseID: "NB-1", SistemBaru: true}}, nil
}

func (polisTiruanH) RingkasanPolis(_ context.Context, noPolis, prodKe string) (models.RingkasanPolis, error) {
	if noPolis != "RNM-1" {
		return models.RingkasanPolis{}, repository.ErrPolisTidakAda
	}
	return models.RingkasanPolis{Kepala: models.KepalaPolis{CaseID: "NB-1", PolicyNo: "RNM-1", PolicyType: "1", StartPeriod: "01-01-2026"},
		Baris:  []models.BarisPolis{{ObjectNo: "1", Rate: "0.15", TSI100: apd.New(1000, 0), PctShare: apd.New(125, -1)}},
		Ceding: []models.CedingPolis{{CedingName: "PT UJI", PctShare: apd.New(100, 0)}}}, nil
}

// kunciBlokTS - kunci 4 spasi blok bersarang `nama: {` di interface api.ts (urut).
func kunciBlokTS(t *testing.T, iface, blok string) []string {
	t.Helper()
	isi, err := os.ReadFile(filepath.Join("..", "..", "frontend", "api.ts"))
	if err != nil {
		t.Fatal(err)
	}
	teks := strings.ReplaceAll(string(isi), "\r\n", "\n")
	teks = teks[strings.Index(teks, "export interface "+iface+" {\n"):]
	awal := strings.Index(teks, "\n  "+blok+": {\n")
	if awal < 0 {
		t.Fatalf("api.ts %s tanpa blok %s", iface, blok)
	}
	b := teks[awal+1:]
	b = b[:strings.Index(b, "\n  }")]
	var out []string
	for _, m := range regexp.MustCompile(`(?m)^    ([A-Za-z0-9_]+)\??: `).FindAllStringSubmatch(b, -1) {
		out = append(out, m[1])
	}
	return out
}

// TestPolisHandler - tiket 71: bentuk = api.ts (VersiPolis, RingkasanPolis kepala / baris dua arah + urutan; cari = HalamanFollowing),
// hanya login (akun tanpa peran workbasket), 400 / 401 / 404.
func TestPolisHandler(t *testing.T) {
	var tulis []string
	svc := layananKasus().DenganFollowing(followingTiruanH{&tulis}).DenganPolis(polisTiruanH{})
	if kode, isi := mintaPeran(t, svc, "GET", "/api/nbfacin/polis?opsi=5&filter=PT&halaman=1", "", "UJI-LAIN", ""); kode != 200 ||
		!strings.HasPrefix(isi, `{"baris":[{"caseId":"NB-9","policyNo":"RNM-F9",`) {
		t.Errorf("cari %d %s", kode, isi)
	}
	if kode, isi := mintaPeran(t, svc, "GET", "/api/nbfacin/polis/versi?nopolis=RNM-1", "", "UJI-LAIN", ""); kode != 200 ||
		strings.TrimSpace(isi) != `{"versi":[{"prodke":"0","caseId":"NB-1","sistemBaru":true}]}` {
		t.Errorf("versi %d %s", kode, isi)
	}
	kode, isi := mintaPeran(t, svc, "GET", "/api/nbfacin/polis/ringkasan?nopolis=RNM-1&prodke=0", "", "UJI-LAIN", "")
	if kode != 200 || !strings.Contains(isi, `"policyType":"1"`) || !strings.Contains(isi, `"tsi100":"1000","tsiRnm":"","rate":"0.15","premiRnm":""`) ||
		!strings.Contains(isi, `"pctShare":"12.5"`) || !strings.Contains(isi, `"ceding":[{"cedingName":"PT UJI","pctShare":"100"}]`) {
		t.Errorf("ringkasan %d %s", kode, isi)
	}
	if ts, gob := kunciInterfaceTS(t, "VersiPolis"), kunciTagJSON(versiPolisKabel{}); !reflect.DeepEqual(ts, gob) {
		t.Errorf("VersiPolis ts %v go %v", ts, gob)
	}
	if ts, gob := kunciInterfaceTS(t, "RingkasanPolis"), kunciTagJSON(ringkasanPolisKabel{}); !reflect.DeepEqual(ts, gob) {
		t.Errorf("RingkasanPolis ts %v go %v", ts, gob)
	}
	if ts, gob := kunciBlokTS(t, "RingkasanPolis", "kepala"), kunciTagJSON(kepalaPolisKabel{}); !reflect.DeepEqual(ts, gob) || len(ts) != 14 {
		t.Errorf("kepala ts %v go %v", ts, gob)
	}
	if ts, gob := kunciBlokTS(t, "RingkasanPolis", "baris"), kunciTagJSON(barisPolisKabel{}); !reflect.DeepEqual(ts, gob) || len(ts) != 13 {
		t.Errorf("baris ts %v go %v", ts, gob)
	}
	for jalur, mau := range map[string]int{"/api/nbfacin/polis/versi?nopolis=RNM-9": 404, "/api/nbfacin/polis/ringkasan?nopolis=RNM-1": 400,
		"/api/nbfacin/polis/versi?nopolis=": 400, "/api/nbfacin/polis?opsi=3&filter=X": 400, "/api/nbfacin/polis?opsi=2&filter=X&halaman=x": 400} {
		if kode, isi := mintaPeran(t, svc, "GET", jalur, "", "UJI-LAIN", ""); kode != mau {
			t.Errorf("%s: %d %s", jalur, kode, isi)
		}
	}
	if kode, _ := mintaPeran(t, svc, "GET", "/api/nbfacin/polis/versi?nopolis=RNM-1", "", "", ""); kode != 401 {
		t.Errorf("401: %d", kode)
	}
}
