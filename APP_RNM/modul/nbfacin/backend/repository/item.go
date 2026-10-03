package repository

// Sub-tab Object Item (tiket 39): .Property.PropertyItemList -> T_PROPERTYITEMLIST (migrasi 188), banyak baris per
// property, urut SEQ_NO; ikut baca/ganti objek (objek.go). Pilihan Object Item Type (V_JN_OBJ_ITEM) dan Currency
// (CURRENCY).
//
// ⛔ Uang/persen tidak pernah float (ADR-0003/0016): dibaca TO_CHAR TM9 ber-NLS titik (db.FmtDesimal), diurai apd,
// ditulis TO_NUMBER dengan topeng dan NLS eksplisit - tanpa bergantung NLS sesi.

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
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

// fmtAngkaMasuk - teks desimal bertitik -> NUMBER(38,8): 30 digit bulat, 8 desimal; NLS titik eksplisit.
const fmtAngkaMasuk = `TO_NUMBER(%s, 'FM999999999999999999999999999999D99999999', 'NLS_NUMERIC_CHARACTERS=''.,''')`

// kolomBacaItem - kolom sqlBacaItem BERNAMA (dibaca lewat kunci) -> medan ItemObjek.
var kolomBacaItem = []string{"TO_CHAR(i.PARENT_ID)", "i.ITEM_TYPE_ID", "i.ITEM_TYPE", "i.PROPERTI_ITEM_NOTE", "i.PROPERTY_YEAR",
	"i.UNIT", "i.CONDITION", "i.CURRENCY", fmt.Sprintf(db.FmtDesimal, "i.TSI_OBJECT_ITEM"), "i.YEAR", "i.NO_OF_TREE",
	"i.AREA_HECTAR", "i.REMARK", "i.IS_ADJUSTABLE_FLAG", fmt.Sprintf(db.FmtDesimal, "i.PCT_ADJUST2"),
	fmt.Sprintf(db.FmtDesimal, "i.PCT_ADJUST_OTHER")}

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
	angka := func(n int) string { return fmt.Sprintf(fmtAngkaMasuk, ":"+strconv.Itoa(n)) }
	return "INSERT INTO " + item + " (ID, PARENT_ID, SEQ_NO, ROW_UID, PROPERTY_ITEM_NO, ITEM_TYPE_ID, ITEM_TYPE," +
		" PROPERTI_ITEM_NOTE, PROPERTY_YEAR, UNIT, CONDITION, CURRENCY, TSI_OBJECT_ITEM, YEAR, NO_OF_TREE, AREA_HECTAR, REMARK," +
		" IS_ADJUSTABLE_FLAG, PCT_ADJUST2, PCT_ADJUST_OTHER) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, " +
		angka(13) + ", :14, :15, :16, :17, :18, " + angka(19) + ", " + angka(20) + ")"
}

// argItem - bind :6..:20 sqlSisipItem, urutan sama dengan daftar kolomnya. Kosong -> NULL; CURRENCY
// tidak pernah kosong (diperiksa services, A133).
func argItem(it models.ItemObjek) []any {
	k := db.KosongJadiNil
	return []any{k(it.ItemTypeID), k(it.ItemType), k(it.Note), k(it.PropertyYear), k(it.Unit), k(it.Condition),
		it.Currency, k(it.TSI), k(it.YearOfPlanting), k(it.NoOfTree), k(it.AreaHectar), k(it.Remark),
		teksBool(it.IsAdjustable), k(it.PctAdjust2), k(it.PctAdjustOther)}
}

// desimalTeks - teks TM9 Oracle (mis. ".5") -> teks desimal kanonik ("0.5") lewat apd; galat urai diteruskan.
func desimalTeks(idBaris, kolom string, v *sql.NullString) (string, error) {
	d, err := db.UraiDesimal(idBaris, kolom, *v)
	if err != nil {
		return "", err
	}
	return utils.FormatDecimal(d), nil
}

// bacaItem - item case `id` per T_PROPERTY.ID (teks).
func (r *ObjekOracle) bacaItem(ctx context.Context, t tabelObjek, id string) (map[string][]models.ItemObjek, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaItem(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca item objek: %w", err)
	}
	defer baris.Close()
	hasil := map[string][]models.ItemObjek{}
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
		for _, d := range []struct {
			kolom string
			ke    *string
		}{{"TSI_OBJECT_ITEM", &it.TSI}, {"PCT_ADJUST2", &it.PctAdjust2}, {"PCT_ADJUST_OTHER", &it.PctAdjustOther}} {
			if *d.ke, err = desimalTeks(TabelPropertyItemList+" induk "+induk, d.kolom,
				teks[fmt.Sprintf(db.FmtDesimal, "i."+d.kolom)]); err != nil {
				return nil, err
			}
		}
		hasil[induk] = append(hasil[induk], it)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: item objek: %w", err)
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
