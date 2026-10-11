package handlers

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
)

type kotakMasukTiruanH struct{}

func (kotakMasukTiruanH) HitungKotakMasuk(context.Context, repository.SaringanKotakMasuk) (map[string]int, error) {
	return map[string]int{"ReasFacInTeamLeader": 3}, nil
}
func (kotakMasukTiruanH) DaftarMenunggu(context.Context, repository.SaringanKotakMasuk) ([]models.BarisKotakMasuk, error) {
	return []models.BarisKotakMasuk{{CaseID: "UJI-NB-4", FlagOnGoingPolicy: "0", ApprovalTerakhir: "Accept", PICAkhir: "UJI PIC",
		SubmitToUW: "20261006T020000.000 GMT"}}, nil
}
func (kotakMasukTiruanH) TeamKotakMasuk(context.Context, string) ([]string, []string, error) {
	return nil, []string{"1"}, nil
}

func (kotakMasukTiruanH) NamaWorkbasket(context.Context, []string) (map[string]string, error) {
	return map[string]string{"ReasFacInTeamLeader": "Team Leader"}, nil
}

// TestKotakMasukHandler - bentuk persis kontrak api.ts (larik), 403 workbasket tidak dipegang, 401 tanpa sesi.
func TestKotakMasukHandler(t *testing.T) {
	svc := services.Baru(nil).DenganKotakMasuk(kotakMasukTiruanH{})
	if kode, isi := mintaPeran(t, svc, "GET", "/api/nbfacin/kotak-masuk", "", "UJI-TL", "ReasFacInTeamLeader"); kode != 200 ||
		isi != `[{"workbasket":"ReasFacInTeamLeader","nama":"Team Leader","jumlah":3}]` {
		t.Errorf("kotak masuk: %d %s", kode, isi)
	}
	if kode, isi := mintaPeran(t, svc, "GET", "/api/nbfacin/kotak-masuk/kasus?workbasket=ReasFacInTeamLeader", "", "UJI-TL", "ReasFacInTeamLeader"); kode != 200 ||
		isi != `[{"caseId":"UJI-NB-4","insuredName":"","businessName":"","sourceOfBusiness":"","cedingCompany":"","inceptionDate":"",`+
			`"shareRnm":"","marketing":"","status":"Offer","createdDate":"","offeringDate":"","validityDate":"","statusInbox":"NEW NB",`+
			`"aksiTerakhir":"Accept by UJI PIC","tglUpdate":"","submitToUw":"2026-10-06 09:00:00"}]` {
		t.Errorf("daftar: %d %s", kode, isi)
	}
	if kode, _ := mintaPeran(t, svc, "GET", "/api/nbfacin/kotak-masuk/kasus?workbasket=ReasFacInMarketing", "", "UJI-TL", "ReasFacInTeamLeader"); kode != 403 {
		t.Errorf("403: %d", kode)
	}
	if kode, isi := mintaPeran(t, svc, "GET", "/api/nbfacin/kotak-masuk", "", "UJI", ""); kode != 200 || isi != `[]` {
		t.Errorf("tanpa workbasket: %d %s", kode, isi)
	}
	if kode, _ := mintaPeran(t, svc, "GET", "/api/nbfacin/kotak-masuk", "", "", ""); kode != 401 {
		t.Errorf("401: %d", kode)
	}
}

// kunciInterfaceTS - kunci medan tingkat atas satu `export interface` di frontend/api.ts (baris `  nama: tipe`, dua spasi;
// komentar JSDoc dan medan bersarang dilewati).
func kunciInterfaceTS(t *testing.T, nama string) []string {
	t.Helper()
	isi, err := os.ReadFile(filepath.Join("..", "..", "frontend", "api.ts"))
	if err != nil {
		t.Fatal(err)
	}
	teks := strings.ReplaceAll(string(isi), "\r\n", "\n")
	awal := strings.Index(teks, "export interface "+nama+" {\n")
	if awal < 0 {
		t.Fatalf("api.ts tanpa interface %s", nama)
	}
	blok := teks[awal:]
	blok = blok[:strings.Index(blok, "\n}")]
	var out []string
	for _, m := range regexp.MustCompile(`(?m)^  ([A-Za-z0-9_]+)\??: `).FindAllStringSubmatch(blok, -1) {
		out = append(out, m[1])
	}
	return out
}

// kunciTagJSON - tag json medan struct, urut deklarasi.
func kunciTagJSON(v any) []string {
	var out []string
	ty := reflect.TypeOf(v)
	for i := 0; i < ty.NumField(); i++ {
		out = append(out, strings.Split(ty.Field(i).Tag.Get("json"), ",")[0])
	}
	return out
}

// TestKunciBarisInboxSamaDenganApiTs - dua arah (handler \ api.ts dan api.ts \ handler) + urutan: barisInbox =
// BarisKotakMasukNbFacIn (16 kunci, terakhir submitToUw tiket 58). Instrumen diuji atas jawaban yang sudah diketahui.
func TestKunciBarisInboxSamaDenganApiTs(t *testing.T) {
	ts, gob := kunciInterfaceTS(t, "BarisKotakMasukNbFacIn"), kunciTagJSON(barisInbox{})
	if len(ts) != 16 || ts[0] != "caseId" || ts[15] != "submitToUw" || len(gob) != 16 || gob[15] != "submitToUw" {
		t.Fatalf("instrumen: ts %v, go %v", ts, gob)
	}
	ada := func(d []string) map[string]bool {
		m := map[string]bool{}
		for _, k := range d {
			m[k] = true
		}
		return m
	}
	mTS, mGo := ada(ts), ada(gob)
	for _, k := range gob {
		if !mTS[k] {
			t.Errorf("kunci handler %q tidak ada di api.ts", k)
		}
	}
	for _, k := range ts {
		if !mGo[k] {
			t.Errorf("kunci api.ts %q tidak ada di handler", k)
		}
	}
	if !reflect.DeepEqual(ts, gob) {
		t.Errorf("urutan beda:\n ts %v\n go %v", ts, gob)
	}
}
