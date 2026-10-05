package models

// OQ-PL-12 DITUTUP (GILIRAN-17): kolom uang CSV yang kosong = 0, seperti
// `ValidasiUploadPL_act` langkah 2 (2.1 b2121 ... 2.32 b6691; tidak ter-remark).

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestIsiNolUangKosongSepertiLangkah2(t *testing.T) {
	baris := []BarisUnggah{{Nomor: 1, Nilai: map[string]string{
		"NET_PREMIUM": "", "GROSS_PREMIUM": "  ", "SUM_INSURED": "1000", "CERTIFICATE_NO": "",
	}}}
	IsiNolUangKosong(baris)
	n := baris[0].Nilai
	if n["NET_PREMIUM"] != "0" || n["GROSS_PREMIUM"] != "0" {
		t.Errorf("uang kosong: %q %q, mau 0", n["NET_PREMIUM"], n["GROSS_PREMIUM"])
	}
	if n["CLAIM_AMOUNT"] != "0" {
		t.Errorf("kolom uang yang tidak ada di berkas: %q, mau 0 (@PropertyHasValue salah)", n["CLAIM_AMOUNT"])
	}
	if n["SUM_INSURED"] != "1000" || n["CERTIFICATE_NO"] != "" {
		t.Errorf("nilai terisi atau kolom bukan uang ikut diubah: %v", n)
	}
	// Enam kolom wajib yang kosong kini lolos "HARUS ADA", seperti Pega.
	if h := validasiQR(baris); hasTolakKolom(h, "NET_PREMIUM") {
		t.Errorf("NET_PREMIUM kosong masih ditolak sesudah diisi 0: %+v", h)
	}
}

func hasTolakKolom(h HasilUnggah, kolom string) bool {
	for _, p := range h.Ditolak {
		if p.Kolom == kolom {
			return true
		}
	}
	return false
}

// Himpunannya dari korpus: `.X = 0` langkah 2 = `KolomUangUnggah`, urut sama.
func TestNolUangLangkah2DariKorpus(t *testing.T) {
	isi, err := os.ReadFile(`D:\XML\RNM_BRD\PremiumList Life\Activity\ValidasiUploadPL_act.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	pecahan := strings.ReplaceAll(string(isi), "><", ">\n<")
	var dari []string
	for _, m := range regexp.MustCompile(`<PropertiesName>\.([A-Z_]+)</PropertiesName>\s*<PropertiesValue>0</PropertiesValue>`).
		FindAllStringSubmatch(pecahan, -1) {
		dari = append(dari, m[1])
	}
	if strings.Join(dari, ",") != strings.Join(KolomUangUnggah, ",") {
		t.Errorf("langkah 2 mengisi 0 pada %v\nKolomUangUnggah %v", dari, KolomUangUnggah)
	}
}
