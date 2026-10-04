//go:build db

package repository_test

// Uji seam repository terhadap Oracle SUNGGUHAN (skema uji, `make test-db`):
// pulang-pergi halaman lewat katalog, generasi tertutup ditolak, nomor polis
// sekali dan unik, dan pembatalan transaksi yang tidak menyisakan baris.
//
// Tanpa instance Oracle (ORACLE_DSN / skema uji), seluruh test di sini
// MELEWATI dengan pesan. ⛔ Fixture berawalan UJI-.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	intidb "nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
	"nusantarare/uji/skemauji"
)

func pasang(t *testing.T) (*sql.DB, string, context.Context, *intidb.DB) {
	t.Helper()
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	// T_GENERAL_POLIS tabel dasar bersama nbfacin 182 (WO 04-10-2026): 320
	// hanya ALTER ... ADD ( atasnya - tabel dasar harus ada sebelum Pasang.
	siapkanGeneralPolisDasar(t, ctx, sqlDB, skema)
	if err := skemauji.Pasang(ctx, sqlDB, skema); err != nil {
		t.Fatalf("memasang skema uji: %v", err)
	}
	t.Cleanup(func() { _ = skemauji.Bongkar(ctx, sqlDB, skema) })
	// `BacaHalaman` membaca SuggestList dari tabel warisan
	// HISTORYAKSEPTASIPRODUCTION (K4) - tidak dibuat migrasi mana pun, jadi
	// tanpa tiruan setiap uji yang membuka halaman gagal ORA-00942.
	siapkanRiwayatProduksi(t, ctx, sqlDB, skema)
	repo, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	return sqlDB, skema, ctx, repo
}

func dalamTx(t *testing.T, ctx context.Context, d *intidb.DB, f func(tx *intidb.Tx) error) error {
	t.Helper()
	tx, err := d.Mulai(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := f(tx); err != nil {
		return err
	}
	return tx.Commit()
}

func TestPulangPergiHalamanLewatKatalog(t *testing.T) {
	_, _, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-1"
	h := models.HalamanBaru()
	h.Setel("PositionNote", models.PosisiAdmin)
	h.Setel("PolicyTreatyIn.PremiOgp", "830.82191781") // 8 desimal = skala penuh NUMBER(38,8)
	h.Setel("PolicyTreatyIn.RiCommOgp", "12.5")
	h.Setel("PolicyTreatyIn.StartDate", "2026-10-01")
	h.Setel("PolicyTreatyIn.StatementDate", "2026-10-03 09:15:00")
	h.Setel("PolicyTreatyIn.BizCode", "006")
	h.Setel("Quotation.BusinessOldId", "01")
	h.Setel("Quotation.BusinessFac", "T")
	h.SetelDaftar(models.TabelCeding.Daftar, []models.Baris{{"CedingCo": "UJI-C1", "CedingCoName": "UJI CEDING SATU"}})
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{{"TreatyType": "UJI-10015", "SharePercentage": "33.3333"}, {"TreatyType": "UJI-10218"}})
	h.SetelDaftar(models.DaftarAngsuran, []models.Baris{{"InstallmentNo": "1", "DueDate": "2026-11-01", "Premium": "1.5"}})
	h.SetelDaftar(models.JalurAnak(models.DaftarAngsuran, 1, "InstallmentList"), []models.Baris{{"InstallmentNo": "1", "PremiumAfterTax": "0.75"}})
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		if err := g.SisipKasus(ctx, tx, id, "UJI-AKUN", "UJI NAMA"); err != nil {
			return err
		}
		return g.SimpanHalaman(ctx, tx, id, h)
	}); err != nil {
		t.Fatal(err)
	}
	b, err := g.BacaHalaman(ctx, nil, id)
	if err != nil {
		t.Fatal(err)
	}
	for j, harap := range map[string]string{
		"PolicyTreatyIn.PremiOgp":                  "830.82191781", // skala penuh; 38 digit: TestUangPresisiPenuhTanpaPembulatanRepository (AC 23)
		"PolicyTreatyIn.RiCommOgp":                 "12.5",
		"PolicyTreatyIn.StartDate":                 "2026-10-01",
		"PolicyTreatyIn.StatementDate":             "2026-10-03 09:15:00",
		"PolicyTreatyIn.BizCode":                   "006", // nol di depan bertahan
		"Quotation.BusinessOldId":                  "01",
		"PolicyTreatyIn.QuotationData.BusinessFac": "T",
	} {
		if got := b.Ambil(j); got != harap {
			t.Errorf("%s = %q, harap %q", j, got, harap)
		}
	}
	sp := b.AmbilDaftar(models.DaftarSpreading)
	if len(sp) != 2 || sp[0]["SharePercentage"] != "33.3333" || sp[1]["TreatyType"] != "UJI-10218" {
		t.Errorf("spreading %+v", sp)
	}
	rinci := b.AmbilDaftar(models.JalurAnak(models.DaftarAngsuran, 1, "InstallmentList"))
	if len(rinci) != 1 || rinci[0]["PremiumAfterTax"] != "0.75" {
		t.Errorf("rincian angsuran %+v", rinci)
	}
	// T_POLIS_CEDING di bawah T_POLIS_QUOTATION (QUOTATION_ID, diagram O39) -
	// CEDING_CO_ID dan CEDING_CO_NAME terpisah (AC 28)
	if c := b.AmbilDaftar(models.TabelCeding.Daftar); len(c) != 1 || c[0]["CedingCo"] != "UJI-C1" || c[0]["CedingCoName"] != "UJI CEDING SATU" {
		t.Errorf("ceding %+v", c)
	}
	// NB: baris dihapus, NOURUT dinomori ulang (ID-12) - simpan kedua juga
	// menulis ulang quotation SESUDAH anaknya dihapus (FK ceding -> quotation)
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{{"TreatyType": "UJI-10218"}})
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SimpanHalaman(ctx, tx, id, h) }); err != nil {
		t.Fatal(err)
	}
	b, _ = g.BacaHalaman(ctx, nil, id)
	if sp := b.AmbilDaftar(models.DaftarSpreading); len(sp) != 1 || sp[0]["TreatyType"] != "UJI-10218" {
		t.Errorf("sesudah hapus baris: %+v", sp)
	}
}

// RD `GetListOpportunity`: pencarian hanya filter G (`.TextNoQuotation
// Contains Param.Search`, tanpa beda huruf besar/kecil) atas pengenal kasus;
// nama bisnis bukan saringan (filter C `.Name` tidak dibangun).
func TestDaftarKasusPortalMenurutGetListOpportunity(t *testing.T) { // P8, temuan tinjauan P9
	_, _, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-CARI-1"
	h := models.HalamanBaru()
	h.Setel("PositionNote", models.PosisiAdmin)
	h.Setel("Quotation.BusinessName", "UJI-BISNIS-CARI")
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		if err := g.SisipKasus(ctx, tx, id, "UJI-AKUN", "UJI NAMA"); err != nil {
			return err
		}
		return g.SimpanHalaman(ctx, tx, id, h)
	}); err != nil {
		t.Fatal(err)
	}
	for cari, harap := range map[string]int{"uji-nb-cari": 1, "UJI-BISNIS-CARI": 0} {
		r, err := g.DaftarKasus(ctx, models.SaringanKasus{Cari: cari})
		if err != nil {
			t.Fatal(err)
		}
		n := 0
		for _, b := range r {
			if b.ID == id {
				n++
			}
		}
		if n != harap {
			t.Errorf("cari %q: %d baris, harap %d", cari, n, harap)
		}
	}
}

func TestNomorPolisSekaliDanUnik(t *testing.T) { // spec AC 31, 74; spec-penyimpanan AC 1
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	for i := 1; i <= 2; i++ {
		id := fmt.Sprintf("UJI-NB-NO-%d", i)
		if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SisipKasus(ctx, tx, id, "UJI-AKUN", "UJI") }); err != nil {
			t.Fatal(err)
		}
	}
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SetelNomorPolis(ctx, tx, "UJI-NB-NO-1", "UJI-QR.T1.10.2026.00001") }); err != nil {
		t.Fatal(err)
	}
	err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SetelNomorPolis(ctx, tx, "UJI-NB-NO-1", "UJI-QR.T1.10.2026.00002") })
	if !errors.Is(err, repository.ErrNomorPolisSudahAda) {
		t.Fatalf("nomor kedua untuk berkas yang sama: %v", err)
	}
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SetelNomorPolis(ctx, tx, "UJI-NB-NO-2", "UJI-QR.T1.10.2026.00001") }); err == nil {
		t.Fatal("dua berkas bernomor sama (PRODKE sama) harus ditolak indeks unik")
	}
	// Spec-penyimpanan AC 1: "Test yang menemukan keduanya tersimpan gagal" -
	// kolom dibaca LANGSUNG: berkas kedua tetap tanpa nomor.
	q := fmt.Sprintf(`SELECT NOPOLIS FROM %s.T_GENERAL_POLIS WHERE ID = :1`, skema)
	if got := kolomTeks(t, ctx, sqlDB, q, "UJI-NB-NO-2"); got != "<NULL>" {
		t.Fatalf("berkas kedua tersimpan bernomor %s padahal ditolak", got)
	}
	// Kuncinya PASANGAN (NOPOLIS, PRODKE): nomor sama pada PRODKE lain diterima.
	if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`UPDATE %s.T_GENERAL_POLIS SET PRODKE = 1 WHERE ID = 'UJI-NB-NO-2'`, skema)); err != nil {
		t.Fatal(err)
	}
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SetelNomorPolis(ctx, tx, "UJI-NB-NO-2", "UJI-QR.T1.10.2026.00001") }); err != nil {
		t.Fatalf("nomor sama dengan PRODKE berbeda harus diterima: %v", err)
	}
}

func TestGenerasiTertutupDitolakDanPembatalanUtuh(t *testing.T) { // ID-10, spec AC 29, spec-penyimpanan AC 6
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-TUTUP"
	awal := models.HalamanBaru()
	awal.Setel("PositionNote", models.PosisiAdmin)
	awal.Setel("PolicyTreatyIn.PremiOgp", "100")
	awal.SetelDaftar(models.DaftarSpreading, []models.Baris{{"TreatyType": "UJI-AWAL"}})
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		if err := g.SisipKasus(ctx, tx, id, "UJI-AKUN", "UJI"); err != nil {
			return err
		}
		return g.SimpanHalaman(ctx, tx, id, awal)
	}); err != nil {
		t.Fatal(err)
	}
	// pembatalan: kasus kedua disisipkan lalu transaksi digagalkan
	err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		if err := g.SisipKasus(ctx, tx, "UJI-NB-BATAL", "UJI-AKUN", "UJI"); err != nil {
			return err
		}
		return errors.New("UJI-gagal di tengah")
	})
	if err == nil {
		t.Fatal("harap galat")
	}
	if _, err := g.Keadaan(ctx, nil, "UJI-NB-BATAL"); !errors.Is(err, repository.ErrKasusTidakAda) {
		t.Fatalf("baris tersisa sesudah pembatalan: %v", err)
	}
	// Generasi DITUTUP oleh lahirnya penerus (ID-10): baris lain yang
	// OLD_POLIS_ID-nya menunjuk generasi ini - tanpa kolom penanda.
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SisipKasus(ctx, tx, "UJI-NB-PENERUS", "UJI-AKUN", "UJI") }); err != nil {
		t.Fatal(err)
	}
	if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`UPDATE %s.T_GENERAL_POLIS SET OLD_POLIS_ID = :1, PRODKE = 1 WHERE ID = 'UJI-NB-PENERUS'`, skema), id); err != nil {
		t.Fatal(err)
	}
	if k, err := g.Keadaan(ctx, nil, id); err != nil || !k.GenerasiTertutup {
		t.Fatalf("keadaan generasi berpenerus: %+v %v", k, err)
	}
	ubah := models.HalamanBaru()
	ubah.Setel("PositionNote", models.PosisiAdmin)
	ubah.Setel("PolicyTreatyIn.PremiOgp", "999")
	ubah.SetelDaftar(models.DaftarSpreading, []models.Baris{{"TreatyType": "UJI-UBAH"}})
	err = dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SimpanHalaman(ctx, tx, id, ubah) })
	if !errors.Is(err, repository.ErrGenerasiTertutup) {
		t.Fatalf("generasi tertutup: %v", err)
	}
	// Spec-penyimpanan AC 6: "Test yang menemukan perubahan tersimpan gagal" -
	// induk dan anak dibaca LANGSUNG dari kolom.
	if p := kolomTeks(t, ctx, sqlDB, fmt.Sprintf(`SELECT %s FROM %s.T_GENERAL_POLIS WHERE ID = :1`,
		fmt.Sprintf(intidb.FmtDesimal, "PREMI_OGP"), skema), id); p != "100" {
		t.Errorf("PREMI_OGP generasi tertutup = %s, harap 100", p)
	}
	if tt := kolomTeks(t, ctx, sqlDB, fmt.Sprintf(`SELECT LISTAGG(TREATY_TYPE, ',') WITHIN GROUP (ORDER BY NOURUT)
		FROM %s.T_POLIS_SPREADING WHERE POLIS_ID = :1`, skema), id); tt != "UJI-AWAL" {
		t.Errorf("spreading generasi tertutup = %s, harap UJI-AWAL", tt)
	}
}

// K4, AC 39-44: catatan SuggestList ditulis ke tabel WARISAN
// HISTORYAKSEPTASIPRODUCTION dan dibaca kembali berurut NOURUT. Tabel itu
// TIDAK dibuat migrasi modul ini - `pasang` memakai tabel DBA bila ada, selain
// itu tiruan uji (`siapkanRiwayatProduksi`, U1).
func TestRiwayatProduksiPulangPergi(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-USUL"
	panjang := ""
	for i := 0; i < 3995; i++ {
		panjang += "x"
	}
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		if err := g.SisipKasus(ctx, tx, id, "UJI-AKUN", "UJI NAMA"); err != nil {
			return err
		}
		return g.CatatUsulan(ctx, tx, models.KunciInstans(id), []models.UsulanProduksi{
			{PIC: "UJI A", TglInp: "2026-10-03 15:30:00", Approval: "Reject", Keterangan: panjang, AksesLogin: "UJI-A", Type: "T"},
		})
	}); err != nil {
		t.Fatal(err)
	}
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		return g.CatatUsulan(ctx, tx, models.KunciInstans(id), []models.UsulanProduksi{{PIC: "UJI B", Approval: "Accept", Keterangan: "UJI-2", AksesLogin: "UJI-B"}})
	}); err != nil {
		t.Fatal(err)
	}
	b, err := g.BacaHalaman(ctx, nil, id)
	if err != nil {
		t.Fatal(err)
	}
	c := b.AmbilDaftar(models.DaftarUsulan)
	if len(c) != 2 || c[0]["IsApproved"] != "0" || len(c[0]["Suggest"]) != 3990 || c[0]["Date"] != "2026-10-03 15:30:00" ||
		c[0]["OperatorName"] != "UJI A" || c[1]["IsApproved"] != "1" || c[1]["OperatorID"] != "UJI-B" {
		t.Fatalf("catatan dibaca kembali: %+v", c)
	}
	var no []string
	rows, err := sqlDB.QueryContext(ctx, fmt.Sprintf(`SELECT TO_CHAR(NOURUT) FROM %s.HISTORYAKSEPTASIPRODUCTION WHERE IDPEGA = :1 ORDER BY TO_NUMBER(NOURUT)`, skema), models.KunciInstans(id))
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var n string
		_ = rows.Scan(&n)
		no = append(no, n)
	}
	// K4 [penyimpangan sadar — disetujui WO 04-10-2026] (F2): NOURUT j = pxListSubscript j
	// SuggestList yang dibangun ulang (c[0] = NOURUT 1, c[1] = NOURUT 2) - sama
	// dengan CARI2 `.pxListSubscript` SaveViewSuggest.
	if fmt.Sprint(no) != "[1 2]" {
		t.Fatalf("NOURUT berikutnya per IDPEGA: %v", no)
	}
}
