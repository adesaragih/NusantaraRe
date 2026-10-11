package handlers

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// followingTiruanH - PenyimpanFollowing tiket 70 untuk uji HTTP.
type followingTiruanH struct{ tulis *[]string }

func (f followingTiruanH) CariFollowing(context.Context, string, string, int, int) ([]models.BarisFollowing, int, error) {
	return []models.BarisFollowing{{CaseID: "NB-9", PolicyNo: "RNM-F9", ReffNumber: "SLIP", InsuredName: "PT UJI", SourceOfBusiness: "SOB",
		QQName: "QQ", StartPeriod: "01-01-2026", EndPeriod: "01-01-2027"}}, 1, nil
}

func (f followingTiruanH) AdaPolisFollowing(_ context.Context, caseID, noPolis string) (bool, error) {
	return caseID == "NB-9" && noPolis == "RNM-F9", nil
}

func (f followingTiruanH) TulisFollowing(_ context.Context, _ *db.Tx, caseID, noPolis string, _ *models.NilaiFollowing) error {
	*f.tulis = append(*f.tulis, caseID+"|"+noPolis)
	return nil
}

// TestFollowingHandler - tiket 70 lewat HTTP: GET bentuk api.ts (HalamanFollowing / BarisFollowing, dua arah + urutan), 400
// opsi 3 / halaman; POST nopolis -> KasusNB; 409 sesudah Binding; 404 polis; 400 badan; 401.
func TestFollowingHandler(t *testing.T) {
	var tulis []string
	svc := layananKasus().DenganFollowing(followingTiruanH{&tulis}).
		DenganTahap(tahapTiruanH{models.TahapKasus{Position: "ReasFacInMarketing", FlagOnGoingPolicy: "0"}})
	jalur := "/api/nbfacin/kasus/UJI-NB-1/following"
	kode, isi := mintaPeran(t, svc, "GET", jalur+"?opsi=2&filter=RNM-F9&halaman=1", "", "UJI-MKT", "ReasFacInMarketing")
	if kode != 200 || !strings.HasPrefix(isi, `{"baris":[{"caseId":"NB-9","policyNo":"RNM-F9","reffNumber":"SLIP","insuredName":"PT UJI",`+
		`"sourceOfBusiness":"SOB","qqName":"QQ","startPeriod":"01-01-2026","endPeriod":"01-01-2027"}],"total":1,"halaman":1,"ukuran":20}`) {
		t.Errorf("GET %d %s", kode, isi)
	}
	if ts, gob := kunciInterfaceTS(t, "HalamanFollowing"), kunciTagJSON(halamanFollowingKabel{}); !reflect.DeepEqual(ts, gob) || len(ts) != 4 {
		t.Errorf("HalamanFollowing ts %v go %v", ts, gob)
	}
	if ts, gob := kunciInterfaceTS(t, "BarisFollowing"), kunciTagJSON(barisFollowingKabel{}); !reflect.DeepEqual(ts, gob) || len(ts) != 8 {
		t.Errorf("BarisFollowing ts %v go %v", ts, gob)
	}
	for _, q := range []string{"?opsi=3&filter=X", "?opsi=2&filter=X&halaman=dua", "?opsi=2&filter=X&halaman=0"} {
		if kode, isi := mintaPeran(t, svc, "GET", jalur+q, "", "UJI-MKT", "ReasFacInMarketing"); kode != 400 {
			t.Errorf("%s: %d %s", q, kode, isi)
		}
	}
	if kode, isi := mintaPeran(t, svc, "GET", jalur+"?opsi=2&filter=", "", "UJI-MKT", "ReasFacInMarketing"); kode != 200 ||
		isi != `{"baris":[],"total":0,"halaman":1,"ukuran":20}`+"\n" && isi != `{"baris":[],"total":0,"halaman":1,"ukuran":20}` {
		t.Errorf("filter kosong %d %q", kode, isi)
	}
	badan := `{"caseIdSumber":"NB-9","policyNo":"RNM-F9"}`
	if kode, isi := mintaPeran(t, svc, "POST", jalur+"/nopolis", badan, "UJI-MKT", "ReasFacInMarketing"); kode != 200 ||
		!strings.Contains(isi, `"caseId":"UJI-NB-1"`) || !strings.Contains(isi, `"bolehFollowing":true`) || len(tulis) != 1 {
		t.Errorf("nopolis %d %s %v", kode, isi, tulis)
	}
	if kode, _ := mintaPeran(t, svc, "POST", jalur+"/nopolis", `{"caseIdSumber":"NB-9","policyNo":"RNM-F8"}`, "UJI-MKT", "ReasFacInMarketing"); kode != 404 {
		t.Errorf("404: %d", kode)
	}
	if kode, _ := mintaPeran(t, svc, "POST", jalur+"/nopolis", `rusak`, "UJI-MKT", "ReasFacInMarketing"); kode != 400 {
		t.Errorf("400: %d", kode)
	}
	if kode, _ := mintaPeran(t, svc, "POST", jalur+"/nopolis", badan, "", ""); kode != 401 {
		t.Errorf("401: %d", kode)
	}
	terikat := layananKasus().DenganFollowing(followingTiruanH{&tulis}).
		DenganTahap(tahapTiruanH{models.TahapKasus{Position: "ReasFacInMarketing", FlagOnGoingPolicy: "2", ConfirmBinding: "1"}})
	if kode, isi := mintaPeran(t, terikat, "POST", jalur+"/salin-semua", badan, "UJI-MKT", "ReasFacInMarketing"); kode != 409 ||
		!strings.Contains(isi, "Following Policy tidak dapat diubah setelah Binding.") {
		t.Errorf("409: %d %s", kode, isi)
	}
	if _, isi := mintaPeran(t, terikat, "GET", "/api/nbfacin/kasus/UJI-NB-1", "", "UJI-MKT", "ReasFacInMarketing"); !strings.Contains(isi, `"bolehFollowing":false`) {
		t.Errorf("bolehFollowing %s", isi)
	}
}
