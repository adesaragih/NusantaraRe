package handlers

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// badanBarisUji - badan hitung-baris / PUT kasus spreading S1 (satu coverage C1) dengan satu baris FACOUT 40 %.
const badanBarisUji = `{"percentShare":"50","lokasi":[{"items":[{"coverages":[{"oldId":"C1","spreading":[` +
	`{"treatyType":"10015","treatyName":"FACOUT","sharePercentage":"40"}]}]}]}]}`

const badanSalinUji = `{"percentShare":"50","template":[{"treatyType":"10015","treatyName":"FACOUT","sharePercentage":"100"}]}`

func layananSpreadingTahap(posisi string) (*services.Service, *spreadingTiruanH) {
	src := &spreadingTiruanH{tersimpan: kasusSpreadingUji()}
	return layananKasus().DenganSpreading(src).DenganTahap(tahapTiruanH{models.TahapKasus{Position: posisi, FlagOnGoingPolicy: "0"}}), src
}

// TestSpreadingTahapUWHandler - tiket 59 lewat HTTP: matriks tahap × workbasket untuk PUT spreading / salin / hitung-baris /
// hitung-share dan PUT general (endpoint non-spreading tetap 409 di uw); GET kasus bolehUbahSpreading / ubahBarisSpreading;
// hitung-baris tidak menulis; 401 / 400.
func TestSpreadingTahapUWHandler(t *testing.T) {
	jalur := "/api/nbfacin/kasus/S1"
	for _, c := range []struct {
		posisi, peran                       string
		put, salin, baris, share, general   int
		bolehUbah, ubahSpreading, ubahBaris bool
	}{
		{"Offer", "ReasFacInMarketing", 200, 409, 409, 200, 200, true, true, false},
		{"Offer", "ReasFacInUnderwriting", 403, 409, 409, 200, 403, false, false, false},
		{"ReasFacInTeamLeader", "ReasFacInTeamLeader", 200, 409, 409, 200, 200, true, true, false},
		{"ReasFacInUnderwriting", "ReasFacInUnderwriting", 200, 200, 200, 200, 409, false, true, true},
		{"ReasFacInJuniorUnderwritingB", "ReasFacInJuniorUnderwritingB", 200, 200, 200, 200, 409, false, true, true},
		{"ReasFacInJuniorUnderwritingA", "ReasFacInMarketing", 403, 403, 403, 200, 409, false, false, false},
		{"Decline", "ReasFacInMarketing", 409, 409, 409, 200, 409, false, false, false},
	} {
		svc, src := layananSpreadingTahap(c.posisi)
		got := [5]int{}
		got[0], _ = mintaPeran(t, svc, "PUT", jalur+"/spreading", badanBarisUji, "UJI", c.peran)
		// PUT di marketing / TL: baris badan (FACOUT 40) tidak tersimpan - kasus uji tanpa baris tetap tanpa baris; uw: tersimpan.
		// Dihitung SEGERA: tiruan berbagi slice dengan pemanggil, hitung-share uw sesudahnya mengubahnya di tempat.
		barisPUT := len(src.tersimpan.Lokasi[0].Items[0].Coverages[0].Spreading)
		got[1], _ = mintaPeran(t, svc, "POST", jalur+"/spreading/salin", badanSalinUji, "UJI", c.peran)
		tulis := src.tulis
		var isiBaris string
		got[2], isiBaris = mintaPeran(t, svc, "POST", jalur+"/spreading/hitung-baris", badanBarisUji, "UJI", c.peran)
		got[3], _ = mintaPeran(t, svc, "POST", jalur+"/spreading/hitung-share", `{"percentShare":"50"}`, "UJI", c.peran)
		got[4], _ = mintaPeran(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/general", `{"reffNumber":"UJI"}`, "UJI", c.peran)
		if mau := [5]int{c.put, c.salin, c.baris, c.share, c.general}; got != mau {
			t.Errorf("%s / %s: put, salin, hitung-baris, hitung-share, general = %v, mau %v", c.posisi, c.peran, got, mau)
		}
		if mauBaris := map[bool]int{true: 1, false: 0}[c.ubahBaris]; c.put == 200 && barisPUT != mauBaris {
			t.Errorf("%s: PUT menyimpan %d baris, mau %d", c.posisi, barisPUT, mauBaris)
		}
		if src.tulis != tulis {
			t.Errorf("%s: salin / hitung-baris / hitung-share menulis", c.posisi)
		}
		if c.baris == 409 && !strings.Contains(isiBaris, "hanya tersedia di tahap Underwriting") {
			t.Errorf("%s: pesan 409 %s", c.posisi, isiBaris)
		}
		if c.baris == 200 && !strings.Contains(isiBaris, `"sharePercentage":"40","tsiGrossSpreaded":"240","tsiSpreaded":"200"`) {
			t.Errorf("%s: hitung-baris %s", c.posisi, isiBaris)
		}
		_, isi := mintaPeran(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1", "", "UJI", c.peran)
		for kunci, mau := range map[string]bool{"bolehUbah": c.bolehUbah, "bolehUbahSpreading": c.ubahSpreading, "ubahBarisSpreading": c.ubahBaris} {
			if want := `"` + kunci + `":` + map[bool]string{true: "true", false: "false"}[mau]; !strings.Contains(isi, want) {
				t.Errorf("%s / %s GET kasus tanpa %s", c.posisi, c.peran, want)
			}
		}
	}
	svc, _ := layananSpreadingTahap("ReasFacInUnderwriting")
	if kode, _ := minta(t, svc, "POST", jalur+"/spreading/hitung-baris", badanBarisUji, ""); kode != 401 {
		t.Errorf("hitung-baris 401: %d", kode)
	}
	if kode, _ := minta(t, svc, "POST", jalur+"/spreading/salin", badanSalinUji, ""); kode != 401 {
		t.Errorf("salin 401: %d", kode)
	}
	if kode, isi := mintaPeran(t, svc, "POST", jalur+"/spreading/hitung-baris", `{"percentShare":"50"}`, "UJI", "ReasFacInUnderwriting"); kode != 400 ||
		!strings.Contains(isi, "lokasi") {
		t.Errorf("hitung-baris tanpa lokasi: %d %s", kode, isi)
	}
}

// TestKunciKasusSamaDenganApiTs - dua arah + urutan: kasusKabel = api.ts KasusNB (19 kunci; tiket 59 bolehUbahSpreading /
// ubahBarisSpreading, tiket 60 ubahAkumulasi, tiket 64d kunciPeriode, tiket 70 bolehFollowing, tiket 73 businessType / lini /
// tampilanAneka). Instrumen diuji atas jawaban yang sudah diketahui.
func TestKunciKasusSamaDenganApiTs(t *testing.T) {
	ts, gob := kunciInterfaceTS(t, "KasusNB"), kunciTagJSON(kasusKabel{})
	if len(gob) != 19 || gob[18] != "tampilanAneka" || len(ts) < 18 || ts[0] != "caseId" {
		t.Fatalf("instrumen: ts %v, go %v", ts, gob)
	}
	if !reflect.DeepEqual(ts, gob) {
		t.Errorf("kunci / urutan beda:\n ts %v\n go %v", ts, gob)
	}
}

// TestPutSpreadingUWJenisTreatyHandler - K59-8 lewat HTTP: PUT uw dengan jenis treaty tidak sah -> 400 verbatim tanpa tulisan;
// jenis sah -> 200.
func TestPutSpreadingUWJenisTreatyHandler(t *testing.T) {
	svc, src := layananSpreadingTahap("ReasFacInUnderwriting")
	tidakSah := strings.Replace(badanBarisUji, `"treatyType":"10015"`, `"treatyType":"X9"`, 1)
	kode, isi := mintaPeran(t, svc, "PUT", "/api/nbfacin/kasus/S1/spreading", tidakSah, "UJI", "ReasFacInUnderwriting")
	if kode != 400 || !strings.Contains(isi, "Invalid Treaty Type Spreading!") || src.tulis != 0 {
		t.Fatalf("tidak sah: %d %s tulis %d", kode, isi, src.tulis)
	}
	if kode, isi := mintaPeran(t, svc, "PUT", "/api/nbfacin/kasus/S1/spreading", badanBarisUji, "UJI", "ReasFacInUnderwriting"); kode != 200 || src.tulis != 1 {
		t.Errorf("sah: %d %s", kode, isi)
	}
}

// akumulasiTiruanH - PenyimpanAkumulasiCoverage dalam memori (tiket 60).
type akumulasiTiruanH struct{ tulis *[]models.AkumulasiCoverage }

func (a akumulasiTiruanH) UbahAkumulasiCoverage(_ context.Context, _ *db.Tx, _ string, baris []models.AkumulasiCoverage) error {
	*a.tulis = append(*a.tulis, baris...)
	return nil
}

// TestAkumulasiUWHandler - tiket 60 lewat HTTP: PUT …/akumulasi 200 (uw pemegang) / 403 / 409 (DepHead, marketing), PUT
// …/objek tetap 409 di uw, GET kasus ubahAkumulasi, 401, 400 badan rusak, 409 bentuk beda.
func TestAkumulasiUWHandler(t *testing.T) {
	badan := `{"baris":[{"items":[{"coverages":[{"oldId":"C1","accumulationCode":"AK1","accumulationDescription":"ALAMAT","premium":"999"}]}]}]}`
	for _, c := range []struct {
		posisi, peran        string
		akumulasi, objek     int
		ubahAkumulasi, tulis bool
	}{
		{"ReasFacInUnderwriting", "ReasFacInUnderwriting", 200, 409, true, true},
		{"ReasFacInJuniorUnderwritingB", "ReasFacInJuniorUnderwritingB", 200, 409, true, true},
		{"ReasFacInJuniorUnderwritingA", "ReasFacInMarketing", 403, 409, false, false},
		{"ReasFacInDepHeadUnderwriting", "ReasFacInDepHeadUnderwriting", 200, 409, true, true}, // K60-2
		{"Decline", "ReasFacInMarketing", 409, 409, false, false},
		{"Offer", "ReasFacInMarketing", 409, 200, false, false},
	} {
		var d []models.ObjekFire
		var tulis []models.AkumulasiCoverage
		svc := layananKasus().DenganObjek(objekTiruan{&d}).DenganSpreading(&spreadingTiruanH{tersimpan: kasusSpreadingUji()}).
			DenganAkumulasiCoverage(akumulasiTiruanH{&tulis}).DenganTahap(tahapTiruanH{models.TahapKasus{Position: c.posisi, FlagOnGoingPolicy: "0"}})
		kode, isi := mintaPeran(t, svc, "PUT", "/api/nbfacin/kasus/S1/akumulasi", badan, "UJI", c.peran)
		if kode != c.akumulasi || (len(tulis) == 1) != c.tulis {
			t.Errorf("%s / %s akumulasi: %d %s tulis %v", c.posisi, c.peran, kode, isi, tulis)
		}
		if c.tulis && (tulis[0] != models.AkumulasiCoverage{ID: "11", AccumulationCode: "AK1", AccumulationDescription: "ALAMAT"} || isi != `{"baris":[]}`) {
			t.Errorf("%s tulis %+v jawaban %s", c.posisi, tulis, isi)
		}
		if kode, _ := mintaPeran(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", `{"baris":[]}`, "UJI", c.peran); kode != c.objek {
			t.Errorf("%s PUT objek: %d", c.posisi, kode)
		}
		_, isi = mintaPeran(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1", "", "UJI", c.peran)
		if want := `"ubahAkumulasi":` + map[bool]string{true: "true", false: "false"}[c.ubahAkumulasi] + `,"kunciPeriode":false,"bolehFollowing":`; !strings.Contains(isi, want) {
			t.Errorf("%s GET kasus tanpa %s: %s", c.posisi, want, isi)
		}
	}
	var d []models.ObjekFire
	var tulis []models.AkumulasiCoverage
	svc := layananKasus().DenganObjek(objekTiruan{&d}).DenganSpreading(&spreadingTiruanH{tersimpan: kasusSpreadingUji()}).
		DenganAkumulasiCoverage(akumulasiTiruanH{&tulis}).DenganTahap(tahapTiruanH{models.TahapKasus{Position: "ReasFacInUnderwriting", FlagOnGoingPolicy: "0"}})
	if kode, _ := minta(t, svc, "PUT", "/api/nbfacin/kasus/S1/akumulasi", badan, ""); kode != 401 {
		t.Errorf("401: %d", kode)
	}
	if kode, _ := mintaPeran(t, svc, "PUT", "/api/nbfacin/kasus/S1/akumulasi", `{"baris":`, "UJI", "ReasFacInUnderwriting"); kode != 400 {
		t.Errorf("400: %d", kode)
	}
	if kode, isi := mintaPeran(t, svc, "PUT", "/api/nbfacin/kasus/S1/akumulasi", strings.Replace(badan, `"C1"`, `"C9"`, 1), "UJI",
		"ReasFacInUnderwriting"); kode != 409 || !strings.Contains(isi, "muat ulang") || len(tulis) != 0 {
		t.Errorf("bentuk: %d %s", kode, isi)
	}
}

// TestPeriksaTemplateSpreadingHandler - POST …/spreading/periksa-template lewat HTTP: {template} -> {"pesan": [...]} (kunci =
// api.ts `{ pesan: string[] }`); 409 marketing; 400 badan rusak / tanpa template; 401.
func TestPeriksaTemplateSpreadingHandler(t *testing.T) {
	jalur := "/api/nbfacin/kasus/S1/spreading/periksa-template"
	svc, src := layananSpreadingTahap("ReasFacInUnderwriting")
	kode, isi := mintaPeran(t, svc, "POST", jalur, `{"template":[{"treatyType":"S1","treatyName":"SPL","sharePercentage":"100"}]}`, "UJI", "ReasFacInUnderwriting")
	if kode != 200 || isi != `{"pesan":["Can not proceed spreading without QS"]}` || src.tulis != 0 {
		t.Errorf("SPL saja: %d %s", kode, isi)
	}
	if kode, isi := mintaPeran(t, svc, "POST", jalur, `{"template":[{"treatyType":"Q1","treatyName":"QS","sharePercentage":"50"}]}`, "UJI", "ReasFacInUnderwriting"); kode != 200 || isi != `{"pesan":[]}` {
		t.Errorf("QS: %d %s", kode, isi)
	}
	for nama, c := range map[string]struct {
		badan, peran string
		kode         int
	}{
		"400 rusak": {`{"template":`, "ReasFacInUnderwriting", 400}, "400 tanpa template": {`{}`, "ReasFacInUnderwriting", 400},
		"403": {`{"template":[]}`, "ReasFacInMarketing", 403},
	} {
		if kode, isi := mintaPeran(t, svc, "POST", jalur, c.badan, "UJI", c.peran); kode != c.kode {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	mkt, _ := layananSpreadingTahap("Offer")
	if kode, _ := mintaPeran(t, mkt, "POST", jalur, `{"template":[]}`, "UJI", "ReasFacInMarketing"); kode != 409 {
		t.Errorf("409 marketing: %d", kode)
	}
	if kode, _ := minta(t, svc, "POST", jalur, `{"template":[]}`, ""); kode != 401 {
		t.Errorf("401: %d", kode)
	}
}
