package repository

import (
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
)

// TestKolomWajibCSVSamaDenganRumusRekap - kolom wajib per Type = kolom yang dijumlah rumus rekap PREMIUM/BALANCE
// (keputusan work owner 10-10-2026: "cocokkan dengan rekap").
func TestKolomWajibCSVSamaDenganRumusRekap(t *testing.T) {
	pola := regexp.MustCompile(`NVL\(d\.(\w+), 0\)`)
	for _, tipe := range []string{models.TypeQR, models.TypeQP, models.TypeTP, models.TypeTR} {
		premi, balance := rumusPremiBalance(tipe)
		set := map[string]bool{}
		for _, m := range pola.FindAllStringSubmatch(premi+" "+balance, -1) {
			set[m[1]] = true
		}
		// PENGECUALIAN (keputusan work owner 10-10-2026: "kolom tax, prof_comm, claim tidak perlu"): TAX, PROF_COMM,
		// CLAIM ada di rumus rekap QR/QP tetapi sengaja TIDAK wajib di CSV (kosong = 0, rekap NVL).
		opsional := map[string]bool{"TAX": true, "PROF_COMM": true, "CLAIM": true}
		var rekap []string
		for k := range set {
			if !opsional[k] {
				rekap = append(rekap, k)
			}
		}
		if (tipe == models.TypeQR || tipe == models.TypeQP) && !(set["TAX"] && set["PROF_COMM"] && set["CLAIM"]) {
			t.Errorf("%s: rumus rekap tidak lagi memuat TAX/PROF_COMM/CLAIM - tinjau pengecualian", tipe)
		}
		for k := range opsional {
			if slices.Contains(models.KolomWajibCSV(tipe), k) {
				t.Errorf("%s: %s wajib, padahal tidak perlu (10-10-2026)", tipe, k)
			}
		}
		wajib := models.KolomWajibTipeCSV(tipe)
		sort.Strings(rekap)
		sort.Strings(wajib)
		if strings.Join(rekap, ",") != strings.Join(wajib, ",") {
			t.Errorf("%s: rekap %v, wajib CSV %v", tipe, rekap, wajib)
		}
		// Setiap komponen rekap wajib di kolom wajib keseluruhan Type ini.
		semua := models.KolomWajibCSV(tipe)
		for _, k := range rekap {
			if !slices.Contains(semua, k) {
				t.Errorf("%s: kolom rekap %s tidak wajib", tipe, k)
			}
		}
		// PENGECUALIAN (keputusan work owner 10-10-2026): NET_PREMIUM wajib untuk SEMUA Type walau di luar rumus
		// rekap - satu-satunya kolom wajib bernilai uang di luar rekap.
		if set["NET_PREMIUM"] || !slices.Contains(semua, "NET_PREMIUM") {
			t.Errorf("%s: NET_PREMIUM harus wajib dan di luar rumus rekap", tipe)
		}
	}
}

func TestSQLAcuanDanContohCSV(t *testing.T) {
	for nama, q := range map[string]string{"mata uang": sqlMataUangCSV(uPeserta), "sertifikat": sqlSertifikatCSV(uPeserta)} {
		if penampungUnik(t, nama, q) != 2 || !strings.Contains(q, "d.PREMIUM_LIST_ID = :1 AND (d.EDM_STATUS IS NULL OR d.EDM_STATUS = :2)") {
			t.Errorf("%s: %s", nama, q)
		}
	}
	if !strings.Contains(sqlSertifikatCSV(uPeserta), "SELECT DISTINCT UPPER(TRIM(d.CERTIFICATE_NO))") {
		t.Error("sertifikat tidak dinormalkan")
	}
	q := sqlContohPesertaCSV(uPeserta)
	if penampungUnik(t, "contoh", q) != 2 || !strings.Contains(q, "TO_CHAR(d.DOB, 'DD/MM/YYYY')") ||
		!strings.Contains(q, "TO_CHAR(d.GROSS_PREMIUM, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')") ||
		!strings.HasSuffix(q, "ORDER BY d.CERTIFICATE_NO, d.NAME_OF_INSURED, d.ID FETCH FIRST 1 ROWS ONLY") || strings.Contains(q, "*") {
		t.Errorf("contoh: %s", q)
	}
}
