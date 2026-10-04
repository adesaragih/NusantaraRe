package repository

// Sub-tab Object Item (tiket 39): .Property.PropertyItemList -> T_PROPERTYITEMLIST (migrasi 188), banyak baris per
// property, urut SEQ_NO; ikut baca/ganti objek (objek.go). Pilihan Object Item Type (V_JN_OBJ_ITEM) dan Currency
// (CURRENCY).
//
// ⛔ Uang/persen tidak pernah float (ADR-0003/0016/0034): *apd.Decimal di aplikasi; ikat/baca lewat desimal.go.

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelJenisItem - view warisan POOLDATA, baca saja (kelas ASM-FW-GISFW-Int-V_JN_OBJ_ITEM; `FROM V_JN_OBJ_ITEM`
	// `[terverifikasi]` `RDBList\GetObjectItembyName_SQL.xml`).
	TabelJenisItem = "V_JN_OBJ_ITEM"
	// TabelCurrency - tabel warisan POOLDATA, baca saja (kelas ASM-FW-GISFW-Int-CURRENCY; `from currency`
	// `[terverifikasi]` `RDBList\GetAllCurrency.xml`).
	TabelCurrency = "CURRENCY"
	// BatasJenisItem - pyMaxRecords RD BrowseV_JN_OBJ_ITEM. BatasMataUang - pyMaxRecords RD BrowseCurrency_RD.
	BatasJenisItem = 10000
	BatasMataUang  = 500
	// JenisItemAktif - filter RD `.ISACTIVE = 1` (SQL Pega GetObjectItembyName_SQL: `ISACTIVE ='1'`).
	JenisItemAktif = "1"
	// MataUangDikecualikan - filter RD BrowseCurrency_RD `.Currency != "ITL"`.
	MataUangDikecualikan = "ITL"
)

// kolomBacaItem - kolom sqlBacaItem BERNAMA (dibaca lewat kunci) -> medan ItemObjek.
var kolomBacaItem = []string{"TO_CHAR(i.PARENT_ID)", "i.ITEM_TYPE_ID", "i.ITEM_TYPE", "i.PROPERTI_ITEM_NOTE", "i.PROPERTY_YEAR",
	"i.UNIT", "i.CONDITION", "i.CURRENCY", angkaKeluar("i.TSI_OBJECT_ITEM"), "i.YEAR", "i.NO_OF_TREE",
	"i.AREA_HECTAR", "i.REMARK", "i.IS_ADJUSTABLE_FLAG", angkaKeluar("i.PCT_ADJUST2"), angkaKeluar("i.PCT_ADJUST_OTHER"),
	// tiket 43: kunci coverage dan total item
	"TO_CHAR(i.ID)", angkaKeluar("i.TOTAL_GROSS_PREMI"), angkaKeluar("i.TOTAL_NET_RATE")}

// sqlBacaItem - seluruh item case :1, urut property lalu SEQ_NO.
func sqlBacaItem(t tabelObjek) string {
	return "SELECT " + strings.Join(kolomBacaItem, ", ") + `
FROM ` + t.item + ` i
JOIN ` + t.prop + ` p ON p.ID = i.PARENT_ID
JOIN ` + t.loc + ` l ON l.ID = p.PARENT_ID
WHERE l.PARENT_ID = :1
ORDER BY i.PARENT_ID, i.SEQ_NO`
}

// sqlSisipItem - :1 ID, :2 PARENT_ID, :3 SEQ_NO, :4 ROW_UID, :5 PROPERTY_ITEM_NO, :6..:20 medan (urut kolom).
func sqlSisipItem(item string) string {
	angka := func(n int) string { return angkaMasuk(":" + strconv.Itoa(n)) }
	return "INSERT INTO " + item + " (ID, PARENT_ID, SEQ_NO, ROW_UID, PROPERTY_ITEM_NO, ITEM_TYPE_ID, ITEM_TYPE," +
		" PROPERTI_ITEM_NOTE, PROPERTY_YEAR, UNIT, CONDITION, CURRENCY, TSI_OBJECT_ITEM, YEAR, NO_OF_TREE, AREA_HECTAR, REMARK," +
		" IS_ADJUSTABLE_FLAG, PCT_ADJUST2, PCT_ADJUST_OTHER, TOTAL_GROSS_PREMI, TOTAL_NET_RATE) VALUES (:1, :2, :3, :4, :5, :6, :7," +
		" :8, :9, :10, :11, :12, " + angka(13) + ", :14, :15, :16, :17, :18, " + angka(19) + ", " + angka(20) + ", " + angka(21) +
		", " + angka(22) + ")"
}

// argItem - bind :6..:20 sqlSisipItem, urutan sama dengan daftar kolomnya. Kosong -> NULL; CURRENCY
// tidak pernah kosong (diperiksa services, A133).
func argItem(it models.ItemObjek) []any {
	k := db.KosongJadiNil
	return []any{k(it.ItemTypeID), k(it.ItemType), k(it.Note), k(it.PropertyYear), k(it.Unit), k(it.Condition),
		it.Currency, ikatDesimal(it.TSI), k(it.YearOfPlanting), k(it.NoOfTree), k(it.AreaHectar), k(it.Remark),
		teksBool(it.IsAdjustable), ikatDesimal(it.PctAdjust2), ikatDesimal(it.PctAdjustOther), ikatDesimal(it.TotalGrossPremi),
		ikatDesimal(it.TotalNetRate)}
}

// bacaItem - item case `id` per T_PROPERTY.ID (teks).
func (r *ObjekOracle) bacaItem(ctx context.Context, t tabelObjek, id string) (map[string][]models.ItemObjek, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaItem(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca item objek: %w", err)
	}
	defer baris.Close()
	hasil := map[string][]models.ItemObjek{}
	// letakItem - posisi item di `hasil` dan T_PROPERTYITEMLIST.ID-nya, untuk memasang coverage (tiket 43).
	type letakItem struct {
		induk  string
		indeks int
		id     string
	}
	var letak []letakItem
	for baris.Next() {
		teks := make(map[string]*sql.NullString, len(kolomBacaItem))
		tujuan := make([]any, len(kolomBacaItem))
		for i, k := range kolomBacaItem {
			teks[k] = &sql.NullString{}
			tujuan[i] = teks[k]
		}
		if err := baris.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: item objek: %w", err)
		}
		v := func(k string) string { return teks[k].String }
		induk := v("TO_CHAR(i.PARENT_ID)")
		it := models.ItemObjek{ItemTypeID: v("i.ITEM_TYPE_ID"), ItemType: v("i.ITEM_TYPE"), Note: v("i.PROPERTI_ITEM_NOTE"),
			PropertyYear: v("i.PROPERTY_YEAR"), Unit: v("i.UNIT"), Condition: v("i.CONDITION"), Currency: v("i.CURRENCY"),
			YearOfPlanting: v("i.YEAR"), NoOfTree: v("i.NO_OF_TREE"), AreaHectar: v("i.AREA_HECTAR"), Remark: v("i.REMARK"),
			IsAdjustable: v("i.IS_ADJUSTABLE_FLAG") == teksBenar}
		if err := bacaDesimalKe(TabelPropertyItemList+" induk "+induk, teks,
			kolomDesimal{angkaKeluar("i.TSI_OBJECT_ITEM"), "TSI_OBJECT_ITEM", &it.TSI},
			kolomDesimal{angkaKeluar("i.PCT_ADJUST2"), "PCT_ADJUST2", &it.PctAdjust2},
			kolomDesimal{angkaKeluar("i.PCT_ADJUST_OTHER"), "PCT_ADJUST_OTHER", &it.PctAdjustOther},
			kolomDesimal{angkaKeluar("i.TOTAL_GROSS_PREMI"), "TOTAL_GROSS_PREMI", &it.TotalGrossPremi},
			kolomDesimal{angkaKeluar("i.TOTAL_NET_RATE"), "TOTAL_NET_RATE", &it.TotalNetRate}); err != nil {
			return nil, err
		}
		it.Coverages = []models.CoverageObjek{}
		hasil[induk] = append(hasil[induk], it)
		letak = append(letak, letakItem{induk, len(hasil[induk]) - 1, v("TO_CHAR(i.ID)")})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: item objek: %w", err)
	}
	baris.Close()
	cov, err := r.bacaCoverage(ctx, t, id)
	if err != nil {
		return nil, err
	}
	for _, l := range letak {
		if d, ada := cov[l.id]; ada && l.id != "" {
			hasil[l.induk][l.indeks].Coverages = d
		}
	}
	return hasil, nil
}

// sisipItem - item satu property, SEQ_NO dan PROPERTY_ITEM_NO = 1..n (K-6 tiket 39).
func (r *ObjekOracle) sisipItem(ctx context.Context, tx *db.Tx, t tabelObjek, idProperty string, item []models.ItemObjek) error {
	for j, it := range item {
		idItem, err := r.db.NomorBerikut(ctx, tx, sequencePropertyItemList)
		if err != nil {
			return err
		}
		uid, err := uidAcak()
		if err != nil {
			return err
		}
		arg := append([]any{idItem, idProperty, j + 1, uid, strconv.Itoa(j + 1)}, argItem(it)...)
		h, err := jalankan(ctx, tx, sqlSisipItem(t.item), "menyisipkan "+TabelPropertyItemList, arg...)
		if err != nil {
			return err
		}
		if err := db.PastikanSatuBaris(h, TabelPropertyItemList); err != nil {
			return err
		}
		if err := r.sisipCoverage(ctx, tx, t, idItem, it.Currency, it.Coverages); err != nil {
			return err
		}
	}
	return nil
}

// PembacaJenisItem - pilihan Object Item Type.
type PembacaJenisItem interface {
	DaftarJenisItem(ctx context.Context) ([]models.JenisItem, error)
}

// PembacaMataUang - pilihan Currency (juga dipakai memeriksa mata uang item).
type PembacaMataUang interface {
	DaftarMataUang(ctx context.Context) ([]string, error)
}

// PilihanItemOracle - PembacaJenisItem + PembacaMataUang atas Oracle.
type PilihanItemOracle struct{ db *db.DB }

// NewPilihanItemOracle merakit pembaca V_JN_OBJ_ITEM dan CURRENCY.
func NewPilihanItemOracle(d *db.DB) *PilihanItemOracle { return &PilihanItemOracle{db: d} }

// sqlJenisItem - RD BrowseV_JN_OBJ_ITEM `[terverifikasi]`: logika `A AND B AND C` - A `.JN_OBJ_ITEM = Param`, C
// `.KELOMPOK = Param` (dropdown sel tidak mengisi keduanya -> filter dibuang, pola RD tiket 38), B `.ISACTIVE = 1`;
// DISTINCT; urut JN_OBJ_ITEM ASC; maks 10000. KETERANGAN ikut dibaca (kolom RD yang sama) - A134. Tipe kolom view
// `belum terverifikasi` (DDL tidak ada): MJOI_KODE lewat TO_CHAR, ISACTIVE dibandingkan teks seperti SQL Pega.
func sqlJenisItem(v string) string {
	return "SELECT DISTINCT TO_CHAR(MJOI_KODE), JN_OBJ_ITEM, KETERANGAN FROM " + v +
		" WHERE ISACTIVE = :1 ORDER BY JN_OBJ_ITEM, 1 FETCH FIRST :2 ROWS ONLY"
}

// sqlMataUang - RD BrowseCurrency_RD `[terverifikasi]`: `.Currency != "ITL"` (filter berparameter lain dibuang);
// tanpa DISTINCT; maks 500. Urutan tidak ada di RD -> urut CURRENCY (A135).
func sqlMataUang(c string) string {
	return "SELECT CURRENCY FROM " + c + " WHERE CURRENCY <> :1 ORDER BY CURRENCY FETCH FIRST :2 ROWS ONLY"
}

// DaftarJenisItem - lihat PembacaJenisItem.
func (r *PilihanItemOracle) DaftarJenisItem(ctx context.Context) ([]models.JenisItem, error) {
	q, err := r.db.Qualify(TabelJenisItem)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlJenisItem(q), JenisItemAktif, BatasJenisItem)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelJenisItem, err)
	}
	defer baris.Close()
	hasil := []models.JenisItem{}
	for baris.Next() {
		var kode, nama, ket sql.NullString
		if err := baris.Scan(&kode, &nama, &ket); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelJenisItem, err)
		}
		hasil = append(hasil, models.JenisItem{Kode: kode.String, Nama: nama.String, Keterangan: ket.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelJenisItem, err)
	}
	return hasil, nil
}

// DaftarMataUang - lihat PembacaMataUang.
func (r *PilihanItemOracle) DaftarMataUang(ctx context.Context) ([]string, error) {
	q, err := r.db.Qualify(TabelCurrency)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlMataUang(q), MataUangDikecualikan, BatasMataUang)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelCurrency, err)
	}
	defer baris.Close()
	hasil := []string{}
	for baris.Next() {
		var c sql.NullString
		if err := baris.Scan(&c); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelCurrency, err)
		}
		hasil = append(hasil, c.String)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelCurrency, err)
	}
	return hasil, nil
}
