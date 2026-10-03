package repository

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLItem - tiket 39: item dibaca lewat case (JOIN property -> lokasi) urut induk lalu SEQ_NO; uang/persen
// dibaca TO_CHAR TM9 ber-NLS titik dan ditulis TO_NUMBER bertopeng + NLS - tidak pernah bergantung NLS sesi.
func TestSQLItem(t *testing.T) {
	baca := sqlBacaItem(tabelObjekUji())
	for _, harus := range []string{"FROM UJI.I i", "JOIN UJI.P p ON p.ID = i.PARENT_ID", "JOIN UJI.L l ON l.ID = p.PARENT_ID",
		"WHERE l.PARENT_ID = :1", "ORDER BY i.PARENT_ID, i.SEQ_NO",
		"TO_CHAR(i.TSI_OBJECT_ITEM, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')"} {
		if !strings.Contains(baca, harus) {
			t.Errorf("baca item tanpa %q", harus)
		}
	}
	sisip := sqlSisipItem("UJI.I")
	m := regexp.MustCompile(`\(([^)]*)\) VALUES \((.*)\)$`).FindStringSubmatch(sisip)
	kolom := strings.Split(m[1], ", ")
	// VALUES dipisah per bind: TO_NUMBER(...) memuat koma di dalam kutip, jadi dipasangkan lewat nomor bind.
	nilai := regexp.MustCompile(`TO_NUMBER\(:\d+[^)]*\)|:\d+`).FindAllString(m[2], -1)
	if len(kolom) != 20 || len(nilai) != 20 {
		t.Fatalf("%d kolom / %d nilai, mau 20", len(kolom), len(nilai))
	}
	uang := map[string]bool{"TSI_OBJECT_ITEM": true, "PCT_ADJUST2": true, "PCT_ADJUST_OTHER": true}
	for i, k := range kolom {
		if bind := fmt.Sprintf(":%d", i+1); !strings.Contains(nilai[i], bind) {
			t.Errorf("kolom %s bind %q, mau %s", k, nilai[i], bind)
		}
		if adalah := strings.HasPrefix(nilai[i], "TO_NUMBER("); adalah != uang[k] {
			t.Errorf("kolom %s TO_NUMBER=%v, mau %v", k, adalah, uang[k])
		}
	}
	if !strings.Contains(sisip, "'FM999999999999999999999999999999D99999999', 'NLS_NUMERIC_CHARACTERS=''.,''')") {
		t.Error("topeng TO_NUMBER harus 30 digit bulat + 8 desimal (NUMBER(38,8)) dengan NLS titik")
	}
}

// TestArgItemMenurutKolom - bind :6..:20 dipasangkan lewat KUNCI (teks: nilai uji = nama kolomnya sendiri; desimal:
// angka berbeda per kolom); kosong -> NULL kecuali CURRENCY (wajib) dan IS_ADJUSTABLE_FLAG ("false").
func TestArgItemMenurutKolom(t *testing.T) {
	it := models.ItemObjek{ItemTypeID: "ITEM_TYPE_ID", ItemType: "ITEM_TYPE", Note: "PROPERTI_ITEM_NOTE", PropertyYear: "PROPERTY_YEAR",
		Unit: "UNIT", Condition: "CONDITION", Currency: "CURRENCY", TSI: desimalUji("1.5"), YearOfPlanting: "YEAR",
		NoOfTree: "NO_OF_TREE", AreaHectar: "AREA_HECTAR", Remark: "REMARK", PctAdjust2: desimalUji("2"), PctAdjustOther: desimalUji("3.25")}
	desimal := map[string]string{"TSI_OBJECT_ITEM": "1.5", "PCT_ADJUST2": "2", "PCT_ADJUST_OTHER": "3.25"}
	m := regexp.MustCompile(`\(([^)]*)\) VALUES`).FindStringSubmatch(sqlSisipItem("UJI.I"))
	kolom := strings.Split(m[1], ", ")[5:]
	arg := argItem(it)
	if len(arg) != len(kolom) || len(kolom) != 15 {
		t.Fatalf("%d bind, %d kolom, mau 15", len(arg), len(kolom))
	}
	for i, k := range kolom {
		mau := any(k)
		if k == "IS_ADJUSTABLE_FLAG" {
			mau = "false"
		}
		if d, ada := desimal[k]; ada {
			mau = d
		}
		if arg[i] != mau {
			t.Errorf("bind kolom %s = %v", k, arg[i])
		}
	}
	for i, a := range argItem(models.ItemObjek{Currency: "IDR", IsAdjustable: true}) {
		switch kolom[i] {
		case "CURRENCY":
			if a != "IDR" {
				t.Errorf("CURRENCY = %v", a)
			}
		case "IS_ADJUSTABLE_FLAG":
			if a != "true" {
				t.Errorf("IS_ADJUSTABLE_FLAG = %v", a)
			}
		default:
			if a != nil {
				t.Errorf("%s kosong = %v, mau NULL", kolom[i], a)
			}
		}
	}
}

// TestBacaItemMemakaiKolomBernama - setiap kunci v("...") / FmtDesimal di item.go ada di kolomBacaItem, dan setiap
// kolom terpilih dibaca.
func TestBacaItemMemakaiKolomBernama(t *testing.T) {
	b, err := os.ReadFile("item.go")
	if err != nil {
		t.Fatal(err)
	}
	ada, dibaca := map[string]bool{}, map[string]bool{}
	for _, k := range kolomBacaItem {
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
	for _, k := range kolomBacaItem {
		if !dibaca[k] {
			t.Errorf("kolom %s dipilih tetapi tidak dibaca", k)
		}
	}
	if len(kolomBacaItem) != 16 {
		t.Errorf("%d kolom, mau 16 (kunci induk + 15 medan)", len(kolomBacaItem))
	}
}

// TestDesimalTeks - TM9 Oracle membuang nol depan (".5"); dibaca ulang kanonik lewat apd, NULL = kosong, rusak = galat.
func TestDesimalTeks(t *testing.T) {
	for masuk, mau := range map[string]string{".5": "0.5", "1500000.25": "1500000.25", "100": "100", "": ""} {
		v := sql.NullString{String: masuk, Valid: masuk != ""}
		if got, err := bacaDesimal("UJI", "TSI_OBJECT_ITEM", &v); err != nil || utils.FormatDecimal(got) != mau {
			t.Errorf("%q -> %v (%v), mau %q", masuk, got, err, mau)
		}
	}
	if _, err := bacaDesimal("UJI", "TSI_OBJECT_ITEM", &sql.NullString{String: "1,5", Valid: true}); err == nil {
		t.Error("koma desimal harus galat, bukan diam-diam")
	}
}

// TestSQLPilihanItem - tiket 39: V_JN_OBJ_ITEM aktif, DISTINCT, urut JN_OBJ_ITEM; CURRENCY tanpa ITL, urut CURRENCY.
func TestSQLPilihanItem(t *testing.T) {
	if q := sqlJenisItem("UJI.V"); q != "SELECT DISTINCT TO_CHAR(MJOI_KODE), JN_OBJ_ITEM, KETERANGAN FROM UJI.V WHERE ISACTIVE = :1 ORDER BY JN_OBJ_ITEM, 1 FETCH FIRST :2 ROWS ONLY" {
		t.Errorf("jenis item %q", q)
	}
	if q := sqlMataUang("UJI.C"); q != "SELECT CURRENCY FROM UJI.C WHERE CURRENCY <> :1 ORDER BY CURRENCY FETCH FIRST :2 ROWS ONLY" {
		t.Errorf("mata uang %q", q)
	}
	if JenisItemAktif != "1" || MataUangDikecualikan != "ITL" || BatasJenisItem != 10000 || BatasMataUang != 500 {
		t.Error("parameter RD berubah")
	}
}
