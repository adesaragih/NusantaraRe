package repository

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLKerugian - tiket 42: catatan dibaca lewat case urut property lalu SEQ_NO, CoinsData LEFT JOIN, uang lewat
// FmtDesimal; CoinsData dihapus sebelum catatan, catatan sebelum property.
func TestSQLKerugian(t *testing.T) {
	baca := sqlBacaKerugian(tabelObjekUji())
	for _, harus := range []string{"FROM UJI.C c", "JOIN UJI.P p ON p.ID = c.PARENT_ID", "LEFT JOIN UJI.D d ON d.PARENT_ID = c.ID",
		"WHERE l.PARENT_ID = :1", "ORDER BY c.PARENT_ID, c.SEQ_NO", fmt.Sprintf(db.FmtDesimal, "c.CLAIM")} {
		if !strings.Contains(baca, harus) {
			t.Errorf("baca kerugian tanpa %q", harus)
		}
	}
	posisi := map[string]int{}
	for i, q := range sqlHapusObjek(tabelObjekUji()) {
		posisi[strings.Fields(q)[2]] = i
	}
	if !(posisi["UJI.D"] < posisi["UJI.C"] && posisi["UJI.C"] < posisi["UJI.P"]) {
		t.Errorf("urutan hapus CoinsData -> catatan -> property: %v", posisi)
	}
}

// TestSisipKerugianMenurutKolom - bind sqlSisipKerugian dipasangkan lewat KUNCI; uang (AMOUNT, CLAIM,
// PREVENTION_OF_LOSS) tepat lewat TO_NUMBER bertopeng; CURRENCY tidak pernah NULL.
func TestSisipKerugianMenurutKolom(t *testing.T) {
	c := models.CatatanKerugian{DateOfLoss: "DATE_OF_LOSS", LossObject: "LOSS_OBJECT", Currency: "CURRENCY", Amount: desimalUji("1"),
		Claim: desimalUji("2.5"), PreventionOfLoss: desimalUji("3"), CauseOfLoss: "CAUSE_OF_LOSS", Remarks: "REMARKS", Detail: "DETAIL"}
	desimal := map[string]string{"AMOUNT": "1", "CLAIM": "2.5", "PREVENTION_OF_LOSS": "3"}
	m := regexp.MustCompile(`\(([^)]*)\) VALUES \((.*)\)$`).FindStringSubmatch(sqlSisipKerugian("UJI.C"))
	kolom := strings.Split(m[1], ", ")
	nilai := regexp.MustCompile(`TO_NUMBER\(:\d+[^)]*\)|:\d+`).FindAllString(m[2], -1)
	if len(kolom) != 13 || len(nilai) != 13 {
		t.Fatalf("%d kolom / %d nilai, mau 13", len(kolom), len(nilai))
	}
	uang := map[string]bool{"AMOUNT": true, "CLAIM": true, "PREVENTION_OF_LOSS": true}
	for i, k := range kolom {
		if !strings.Contains(nilai[i], fmt.Sprintf(":%d", i+1)) || strings.HasPrefix(nilai[i], "TO_NUMBER(") != uang[k] {
			t.Errorf("kolom %s nilai %q", k, nilai[i])
		}
	}
	arg := argKerugian(c)
	for i, k := range kolom[4:] {
		mau := any(k)
		if d, ada := desimal[k]; ada {
			mau = d
		}
		if arg[i] != mau {
			t.Errorf("bind kolom %s = %v", k, arg[i])
		}
	}
	for i, a := range argKerugian(models.CatatanKerugian{Currency: "IDR"}) {
		if (kolom[4+i] == "CURRENCY") != (a != nil) {
			t.Errorf("%s kosong = %v", kolom[4+i], a)
		}
	}
}

// TestSisipLokasiLossRatio - loss ratio hasil hitung ditulis ke T_LOCATIONLIST: amount lewat TO_NUMBER (:5, :7),
// percent teks (:6, :8).
func TestSisipLokasiLossRatio(t *testing.T) {
	q := sqlSisipLokasi("UJI.L")
	for _, harus := range []string{"LOSS_RATIO1_YEAR_AMOUNT, LOSS_RATIO1_YEAR_PERCENT, LOSS_RATIO35_YEAR_AMOUNT, LOSS_RATIO35_YEAR_PERCENT)",
		fmt.Sprintf(fmtAngkaMasuk, ":5") + ", :6, " + fmt.Sprintf(fmtAngkaMasuk, ":7") + ", :8)"} {
		if !strings.Contains(q, harus) {
			t.Errorf("sisip lokasi tanpa %q: %s", harus, q)
		}
	}
}

// TestBacaKerugianMemakaiKolomBernama - setiap kunci v("...") / FmtDesimal di kerugian.go ada di kolomBacaKerugian
// dan sebaliknya.
func TestBacaKerugianMemakaiKolomBernama(t *testing.T) {
	b, err := os.ReadFile("kerugian.go")
	if err != nil {
		t.Fatal(err)
	}
	ada, dibaca := map[string]bool{}, map[string]bool{}
	for _, k := range kolomBacaKerugian {
		ada[k] = true
	}
	for _, m := range regexp.MustCompile(`v\("([^"]+)"\)`).FindAllStringSubmatch(string(b), -1) {
		dibaca[m[1]] = true
	}
	for _, m := range regexp.MustCompile(`kolomDesimal\{angkaKeluar\("([a-z]\.[A-Z0-9_]+)"\)`).FindAllStringSubmatch(string(b), -1) {
		dibaca[fmt.Sprintf(db.FmtDesimal, m[1])] = true
	}
	for k := range dibaca {
		if !ada[k] {
			t.Errorf("kunci %q dibaca tetapi tidak dipilih", k)
		}
	}
	for _, k := range kolomBacaKerugian {
		if !dibaca[k] {
			t.Errorf("kolom %s dipilih tetapi tidak dibaca", k)
		}
	}
}
