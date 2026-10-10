//go:build db

package repository_test

// Uji seam repository pemuat dokumen lama (tiket 22) terhadap Oracle
// SUNGGUHAN (skema uji, `make test-db`; K11 kosong = belum pernah
// dijalankan). Dokumen fiktif UJI- dipecah `models.PecahDokumenLama`, lalu
// ditulis lewat antarmuka penyimpanan yang SAMA dengan jalur biasa (ID-3,
// AC 56), dan nilainya dibaca LANGSUNG dari kolom (spec-penyimpanan §7:
// pulang-pergi saja tidak cukup).

import (
	"context"
	"fmt"
	"testing"

	intidb "nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

const dokumenUjiLamaProp = `{
 "pxObjClass": "ASM-FW-GISFW-Data-PolicyTreatyIn",
 "PolicyNo": "UJI-QP.T1.10.2017.90001",
 "PremiOgp": "592629512.880000276",
 "StartDate": "20171001",
 "EndDate": "",
 "StatementDate": "20170930T170000.000 GMT",
 "CedingCoName": "UJI-CEDING A; UJI-CEDING B; ",
 "QuotationData": {"ProportionalType": "Proportional", "GroupPanel": "006", "BusinessOldId": "01",
  "CedingCoList": [{"CedingCo": "UJI-C1", "CedingCoName": "UJI-CEDING A"}, {"CedingCo": "UJI-C2", "CedingCoName": "UJI-CEDING B"}]},
 "ListInstallment": [{"InstallmentNo": "1", "DueDate": "20171101", "Premium": "148157378.220000069"}]
}`

const dokumenUjiLamaNonProp = `{
 "pxObjClass": "ASM-FW-GISFW-Data-PolicyTreatyIn",
 "PolicyNo": "UJI-QR.T1.01.2018.90002",
 "StartDate": "20180101",
 "EndDate": "20181231",
 "QuotationData": {"ProportionalType": "NonProportional"},
 "ListInstallment": [{"InstallmentNo": "1", "Premium": "3000.5",
   "InstallmentList": [{"InstallmentNo": "1", "DueDate": "20180131", "Premium": "1500.25"}]}],
 "TreatyXOLList": [{"GrossPremi": "3000.5", "Deduction": "200.5",
   "ValueList": [{"Layer": "1", "LayerType": "UJI-LT", "Deduction": "100.25"}, {"Layer": "2", "LayerType": "UJI-LT", "Deduction": "100.25"}]}]
}`

// muatLama = urutan tulis `services.Pemuat` untuk satu dokumen, satu transaksi.
func muatLama(t *testing.T, ctx context.Context, d *intidb.DB, g *repository.Gudang, b models.BarisJSONPolis) models.HasilPecah {
	t.Helper()
	h, err := models.PecahDokumenLama(b)
	if err != nil || len(h.Galat) > 0 {
		t.Fatalf("pecah: %v %+v", err, h.Galat)
	}
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		if err := g.SisipKasus(ctx, tx, h.ID, "", ""); err != nil {
			return err
		}
		if err := g.SimpanHalaman(ctx, tx, h.ID, h.Halaman); err != nil {
			return err
		}
		if err := g.SetelNomorPolis(ctx, tx, h.ID, h.NoPolis); err != nil {
			return err
		}
		if err := g.SetelKolomDatarLama(ctx, tx, h.ID, h.Datar); err != nil {
			return err
		}
		return g.TutupKasus(ctx, tx, h.ID, models.AssignmentAdmin, models.StatusSelesai)
	}); err != nil {
		t.Fatal(err)
	}
	return h
}

func TestPemuatLamaMenulisLewatAntarmukaSama(t *testing.T) { // AC 21, 22, 55, 56, 68, 69
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	b := models.BarisJSONPolis{IDPega: "ASM-FW-GISFW-WORK-NB UJI-990001", NoPolis: "UJI-QP.T1.10.2017.90001",
		ProdKe: "0", TglInput: "2017-10-02 08:00:00", Username: "UJI-AKUN", DataJSON: []byte(dokumenUjiLamaProp)}
	h := muatLama(t, ctx, d, g, b)

	var premi, mulai, akhir, statement, idpega, user, nopol, prodke string
	q := fmt.Sprintf(`SELECT %s, TO_CHAR(START_DATE, 'YYYY-MM-DD'), TO_CHAR(END_DATE, 'YYYY-MM-DD'),
	        TO_CHAR(STATEMENT_DATE, 'YYYY-MM-DD HH24:MI:SS'), IDPEGA, USERNAME, NOPOLIS, TO_CHAR(PRODKE)
	   FROM %s.T_GENERAL_POLIS_TREATY WHERE ID = :1`, fmt.Sprintf(intidb.FmtDesimal, "PREMI_OGP"), skema)
	if err := sqlDB.QueryRowContext(ctx, q, h.ID).Scan(&premi, &mulai, &akhir, &statement, &idpega, &user, &nopol, &prodke); err != nil {
		t.Fatal(err)
	}
	for nama, pasangan := range map[string][2]string{
		"PREMI_OGP (AC 19, 55; F20 skala 10: utuh)": {premi, "592629512.880000276"},
		"START_DATE (AC 21)":                        {mulai, "2017-10-01"},
		"END_DATE (AC 69)":                          {akhir, "2017-10-01"},
		"STATEMENT_DATE (AC 22)":                    {statement, "2017-10-01 00:00:00"},
		"IDPEGA (ID-21)":                            {idpega, b.IDPega},
		"USERNAME (ID-21)":                          {user, "UJI-AKUN"},
		"NOPOLIS":                                   {nopol, "UJI-QP.T1.10.2017.90001"},
		"PRODKE":                                    {prodke, "0"},
	} {
		if pasangan[0] != pasangan[1] {
			t.Errorf("%s = %q, harap %q", nama, pasangan[0], pasangan[1])
		}
	}
	k, err := g.Keadaan(ctx, nil, h.ID)
	if err != nil || k.StatusWork != models.StatusSelesai || k.Position != "" {
		t.Errorf("kasus lama harus tertutup Resolved-Completed: %+v %v", k, err)
	}
	baca, err := g.BacaHalaman(ctx, nil, h.ID) // AC 68: berkas lama dapat dibuka
	if err != nil {
		t.Fatal(err)
	}
	if baca.Ambil("PolicyTreatyIn.CedingCoName") != "UJI-CEDING A; UJI-CEDING B; " ||
		len(baca.AmbilDaftar(models.HalamanPolis+".QuotationData.CedingCoList")) != 2 ||
		baca.Ambil("Quotation.GroupPanel") != "006" {
		t.Errorf("halaman terbaca %+v", baca)
	}
	if ada, err := g.AdaKasus(ctx, nil, h.ID); err != nil || !ada {
		t.Errorf("AdaKasus %v %v - dokumen yang baru dimuat harus ada", ada, err)
	}
}

// dokumenUjiLamaUsulan - dokumen lama dengan dua catatan SuggestList (F3).
// AKSES_LOGIN ditulis NULL: baris SuggestList tidak punya anggota operator
// (dataguide; AddToListCommentsPolicyTreatyIn_DT) - tidak dikarang.
const dokumenUjiLamaUsulan = `{
 "pxObjClass": "ASM-FW-GISFW-Data-PolicyTreatyIn",
 "PolicyNo": "UJI-QP.T1.10.2017.90003",
 "StartDate": "20171001",
 "QuotationData": {"ProportionalType": "Proportional", "BusinessFac": "T", "BusinessCode": "UJI-B01"},
 "SuggestList": [
  {"Date": "20171002T020000.000 GMT", "IsApproved": "1", "OperatorName": "UJI-PENGGUNA A", "Suggest": "UJI-catatan satu"},
  {"Date": "20171003T100000.000 GMT", "IsApproved": "0", "OperatorName": "UJI-PENGGUNA B", "Suggest": "UJI-catatan dua"}
 ]
}`

// F3 (WO 04-10-2026): SuggestList dokumen lama disalin ke riwayat produksi
// lewat `CatatUsulan` jalur biasa, dengan penjaga dobel menurut IDPEGA -
// salinan kedua (pemuat diulang) tidak menggandakan baris. Nilai harapan dari
// `SaveViewSuggest` langkah 2.1.2 (dihitung tangan, lihat models
// TestSuggestListLamaDisalinMenurutSaveViewSuggest), dibaca LANGSUNG dari kolom.
func TestPemuatLamaMenyalinSuggestListSekaliMenurutIDPega(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	b := models.BarisJSONPolis{IDPega: "ASM-FW-GISFW-WORK-NB UJI-990003", NoPolis: "UJI-QP.T1.10.2017.90003",
		ProdKe: "0", DataJSON: []byte(dokumenUjiLamaUsulan)}
	t.Cleanup(func() {
		_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`DELETE FROM %s.HISTORYAKSEPTASIPRODUCTION WHERE IDPEGA = :1`, skema), b.IDPega)
	})
	h := muatLama(t, ctx, d, g, b)
	// langkah services `salinUsulanLama`: salinan kedua melihat baris pertama (proteksi dobel WO 07-10-2026)
	for i, harapAda := range []bool{false, true} {
		if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
			ada, err := g.AdaRiwayatIDPega(ctx, tx, b.IDPega)
			if err != nil {
				return err
			}
			if ada != harapAda {
				t.Errorf("salinan ke-%d: ada %v, harap %v (proteksi dobel IDPEGA)", i+1, ada, harapAda)
			}
			if ada {
				return nil
			}
			return g.CatatUsulan(ctx, tx, b.IDPega, h.Usulan)
		}); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := sqlDB.QueryContext(ctx, fmt.Sprintf(`SELECT TO_CHAR(NOURUT), TYPE_POLIS, POSISI, PIC,
	        TO_CHAR(TGL_INP, 'YYYY-MM-DD HH24:MI:SS'), TYPE, PUTARAN, APPROVAL, KETERANGAN,
	        NVL(AKSES_LOGIN, '<NULL>'), BUSINESS_CODE
	   FROM %s.HISTORYAKSEPTASIPRODUCTION WHERE IDPEGA = :1 ORDER BY TO_NUMBER(NOURUT)`, skema), b.IDPega)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var dapat []string
	for rows.Next() {
		var v [11]string
		if err := rows.Scan(&v[0], &v[1], &v[2], &v[3], &v[4], &v[5], &v[6], &v[7], &v[8], &v[9], &v[10]); err != nil {
			t.Fatal(err)
		}
		dapat = append(dapat, fmt.Sprint(v))
	}
	harap := []string{
		fmt.Sprint([11]string{"1", "UJI", "Policy", "UJI-PENGGUNA A", "2017-10-02 09:00:00", "T", "2", "Accept", "UJI-catatan satu", "<NULL>", "UJI-B01"}),
		fmt.Sprint([11]string{"2", "UJI", "Policy", "UJI-PENGGUNA B", "2017-10-03 17:00:00", "T", "2", "Reject", "UJI-catatan dua", "<NULL>", "UJI-B01"}),
	}
	if fmt.Sprint(dapat) != fmt.Sprint(harap) {
		t.Errorf("riwayat produksi\n dapat %v\n harap %v", dapat, harap)
	}
}

func TestPemuatLamaNonProporsionalBersarang(t *testing.T) { // AC 50, 53
	_, _, ctx, d := pasang(t)
	g := repository.Baru(d)
	b := models.BarisJSONPolis{IDPega: "ASM-FW-GISFW-WORK-NB UJI-990002", NoPolis: "UJI-QR.T1.01.2018.90002",
		ProdKe: "0", DataJSON: []byte(dokumenUjiLamaNonProp)}
	h := muatLama(t, ctx, d, g, b)
	baca, err := g.BacaHalaman(ctx, nil, h.ID)
	if err != nil {
		t.Fatal(err)
	}
	if r := baca.AmbilDaftar(models.JalurAnak(models.DaftarAngsuran, 1, "InstallmentList")); len(r) != 1 || r[0]["DueDate"] != "2018-01-31" {
		t.Errorf("rincian angsuran %+v", r)
	}
	if l := baca.AmbilDaftar(models.JalurAnak(models.HalamanPolis+".TreatyXOLList", 1, "ValueList")); len(l) != 2 || l[1]["Layer"] != "2" {
		t.Errorf("layer %+v", l)
	}
}
