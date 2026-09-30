package handlers_test

// Konformansi galat kelima jalur simpan (paket 8, tiket 10 -> K8): procedure
// tidak dipanggil (keputusan o), jadi yang ditegakkan adalah galat Go sendiri -
// Oracle, validasi, keadaan data - sampai ke layar berkata-kata, tanpa markup,
// tanpa teks mentah basis data; teks aslinya tetap di log server.

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/mastercontractretrolife/backend/handlers"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

const galatOracleBerHTML = `<span style="color:red">Data gagal disimpan</span> ORA-01722: invalid number`

func tangkapLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var b bytes.Buffer
	log.SetOutput(&b)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	return &b
}

func TestKelimaJalurSimpanGagalTerangTanpaMarkup(t *testing.T) {
	kasus := []struct {
		nama, penulis, jalur, badan string
	}{
		{"tahun", "SisipTahun", "/tahun", tahunBaru},
		{"kontrak", "SisipKontrak", "/tahun/UJI-T1/kontrak", kontrakBaru},
		{"reinsurer", "SisipReinsurer", "/kontrak/UJI-K1/reinsurer",
			`{"reinsurerName":"x","reinsurerId":"UJI-L01","pctShare":"40","komisi":"25","ovrComm":"0"}`},
		{"security", "SisipSecurity", "/reinsurer/UJI-R1/security", `{"reinsurerName":"x","reinsurerId":"UJI-L01","pctShare":"10"}`},
		{"business", "SisipBusiness", "/kontrak/UJI-K1/business",
			`{"bizCode":"UJI-B01","bizName":"x","riRateId":"UJI-RATE","riRate":"UJI R"}`},
	}
	for _, k := range kasus {
		u := serverBusiness(t)
		u.g.MasterRe = []models.MasterReinsurer{{ID: "UJI-L01", ClientName: "UJI REASURANSI"}}
		u.g.Reinsurer["UJI-R1"] = models.Reinsurer{ID: "UJI-R1", TreatyYearID: "UJI-T1", TreatyContractID: "UJI-K1"}
		u.g.GagalTulis[k.penulis] = errors.New(galatOracleBerHTML)
		catatan := tangkapLog(t)
		kode, badan := u.kirim(t, "POST", handlers.Prefix+k.jalur, k.badan)
		if kode != http.StatusInternalServerError || !strings.Contains(badan, `"galat":"failed to process`) {
			t.Errorf("%s: penyimpanan yang gagal tampak berhasil / tanpa kalimat: %d %s", k.nama, kode, badan)
		}
		if strings.Contains(badan, "<") || strings.Contains(badan, "u003c") || strings.Contains(badan, "ORA-") {
			t.Errorf("%s: markup / teks mentah Oracle bocor ke API: %s", k.nama, badan)
		}
		if !strings.Contains(catatan.String(), "ORA-01722") {
			t.Errorf("%s: pesan asli tidak tercatat di log server: %q", k.nama, catatan.String())
		}
	}
}

func TestMasterTidakTerbacaMenyebutObjekTanpaTeksOracle(t *testing.T) {
	u := server(t, true)
	u.g.GalatMaster = errors.New(galatOracleBerHTML)
	catatan := tangkapLog(t)
	kode, badan := u.minta(t, "GET", handlers.Prefix+"/master-reinsurer?cari=UJI", true)
	if kode != http.StatusServiceUnavailable || !strings.Contains(badan, "AGENT") {
		t.Errorf("master: %d %s", kode, badan)
	}
	if strings.Contains(badan, "<") || strings.Contains(badan, "u003c") || strings.Contains(badan, "ORA-") {
		t.Errorf("teks mentah Oracle bocor: %s", badan)
	}
	if !strings.Contains(catatan.String(), "ORA-01722") {
		t.Errorf("sebab asli tidak tercatat: %q", catatan.String())
	}
}

// Nilai masukan yang berisi markup tidak kembali sebagai markup.
func TestMarkupDariMasukanTidakKembaliKeLayar(t *testing.T) {
	kode, badan := server(t, true).kirim(t, "POST", handlers.Prefix+"/tahun",
		`{"treatyYear":"2027","underwritingYear":"2027","startDate":"<b>x</b>","endDate":"2027-12-31"}`)
	if kode != http.StatusUnprocessableEntity || strings.Contains(badan, "<b>") || strings.Contains(badan, "u003c") {
		t.Errorf("markup masukan: %d %s", kode, badan)
	}
}

// ⛔ Uji statik: setiap rute menjawab lewat `tulis`/`tulisDaftar` (=> jawabGalat)
// - tidak ada rute yang menulis jawaban sendiri dan menelan galat.
func TestSetiapRuteMenjawabLewatJawabGalat(t *testing.T) {
	polaRute := regexp.MustCompile(`pasang\("`)
	polaJawab := regexp.MustCompile(`\btulis(Daftar)?\(w, `)
	for _, berkas := range []string{"rute_mcrl.go", "rute_tulis.go"} {
		isi, err := os.ReadFile(berkas)
		if err != nil {
			t.Fatal(err)
		}
		rute, jawab := len(polaRute.FindAll(isi, -1)), len(polaJawab.FindAll(isi, -1))
		if rute == 0 || rute != jawab {
			t.Errorf("%s: %d rute, %d jawaban lewat tulis/tulisDaftar", berkas, rute, jawab)
		}
	}
}
