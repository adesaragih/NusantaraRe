//go:build db

package repository_test

// Uji seam repository lawan Oracle SUNGGUHAN yang menutup klaim AC tertahan
// K11 (PROMPT putaran 3 bab 3.3): setiap uji di sini menguji BUNYI AC-nya,
// bukan sekadar pulang-pergi beberapa medan.
//
//   - spec AC 23            uang berpresisi penuh: 38 digit utuh; pembulatan
//                           hanya di desimal ke-9 (NUMBER(38,8), AC 55), bukan 2
//   - spec AC 29, peny. 46  kegagalan di tengah urutan simpan = NOL baris di
//                           kesembilan tabel (T_WORK_POLIS + delapan diagram)
//   - peny. AC 49           polis proporsional: SETIAP kolom katalog kembali sama
//   - peny. AC 50           polis non-proporsional: SETIAP kolom, SETIAP layer
//   - peny. AC 51           urutan baris anak = urutan NOURUT (12 baris yang
//                           urutan teksnya berbeda dari urutan tulis)
//
// Nilai harapan adalah LITERAL masukan (kebenaran pulang-pergi), ditulis per
// golongan katalog di `nilaiUji` - bukan dihitung ulang dengan konversi
// repository. ⛔ Belum pernah dijalankan: K11 (skema uji) kosong.
// Pemanggil `skemauji.Buka()` tetap satu (`pasang`). Fixture berawalan UJI-.

import (
	"fmt"
	"strconv"
	"testing"

	intidb "nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

// uangPenuh - 30 digit bulat + 8 desimal = 38 digit, batas NUMBER(38,8);
// float64 hanya membawa ~16 digit, jadi setiap jalur float merusaknya.
const uangPenuh = "123456789012345678901234567890.12345678"

// nilaiUji - nilai UJI- satu kolom katalog menurut golongannya; `i` memvariasikan
// nilai antarkolom supaya kolom yang tertukar ketahuan.
func nilaiUji(k models.Kolom, i int) string {
	switch k.Golongan {
	case models.GolTeks:
		return "UJI " + k.Kolom
	case models.GolKode:
		return "00" + strconv.Itoa(i) // nol di depan bermakna (ID-16)
	case models.GolPenanda:
		return []string{"0", "1"}[i%2] // "0" tidak boleh menjadi NULL (ID-17)
	case models.GolUang:
		return []string{uangPenuh, "-0.5", "0.00000001", "1500000000.25"}[i%4]
	case models.GolPersen:
		return []string{"12.5", "0.00000001", "100", "33.33333333"}[i%4]
	case models.GolTanggal:
		return fmt.Sprintf("2026-%02d-%02d", 1+i%12, 1+i%28)
	case models.GolTanggalWaktu:
		return fmt.Sprintf("2026-10-03 %02d:%02d:%02d", i%24, (7*i)%60, (13*i)%60)
	case models.GolCacah:
		return []string{"9999999999", strconv.Itoa(i + 1)}[i%2] // NUMBER(10) penuh
	}
	panic("golongan tanpa nilai uji: " + string(k.Golongan))
}

// barisUji - n baris anak satu tabel katalog; kolom `kunci` (bila ada) diisi
// nilai pengenal baris supaya urutannya dapat diperiksa.
func barisUji(tb models.Tabel, n, geser int) []models.Baris {
	var out []models.Baris
	for r := 0; r < n; r++ {
		b := models.Baris{}
		for j, k := range tb.Kolom {
			b[k.Properti] = nilaiUji(k, geser+10*r+j)
		}
		out = append(out, b)
	}
	return out
}

// halamanPenuh - halaman yang mengisi SETIAP kolom katalog T_GENERAL_POLIS dan
// T_POLIS_QUOTATION, dua ceding, dua angsuran, dua spreading; non-proporsional
// menambah dua rincian per angsuran dan dua XOL bertiga layer.
func halamanPenuh(jenis string) *models.Halaman {
	h := models.HalamanBaru()
	for i, k := range models.TabelGeneralPolis.Kolom {
		h.Setel(k.Properti, nilaiUji(k, i))
	}
	for i, k := range models.TabelQuotation.Kolom {
		h.Setel(models.HalamanPolis+".QuotationData."+k.Properti, nilaiUji(k, 100+i))
	}
	h.Setel(models.HalamanPolis+".QuotationData.ProportionalType", jenis)
	h.SetelDaftar(models.TabelCeding.Daftar, barisUji(models.TabelCeding, 2, 200))
	h.SetelDaftar(models.DaftarAngsuran, barisUji(models.TabelAngsuran, 2, 300))
	h.SetelDaftar(models.DaftarSpreading, barisUji(models.TabelSpreading, 2, 400))
	if jenis == models.JenisNonProporsional {
		for i := 1; i <= 2; i++ {
			h.SetelDaftar(models.JalurAnak(models.DaftarAngsuran, i, models.TabelAngsuranRinci.Daftar),
				barisUji(models.TabelAngsuranRinci, 2, 500+100*i))
		}
		h.SetelDaftar(models.TabelXOL.Daftar, barisUji(models.TabelXOL, 2, 800))
		for i := 1; i <= 2; i++ {
			h.SetelDaftar(models.JalurAnak(models.TabelXOL.Daftar, i, models.TabelLayerXOL.Daftar),
				barisUji(models.TabelLayerXOL, 3, 900+100*i))
		}
	}
	return h
}

// samakanDaftar membandingkan baris yang terbaca dengan baris tulisan, SETIAP
// kolom katalog tabelnya.
func samakanDaftar(t *testing.T, jalur string, tb models.Tabel, harap, dapat []models.Baris) {
	t.Helper()
	if len(dapat) != len(harap) {
		t.Errorf("%s: %d baris terbaca, harap %d", jalur, len(dapat), len(harap))
		return
	}
	for r := range harap {
		for _, k := range tb.Kolom {
			if dapat[r][k.Properti] != harap[r][k.Properti] {
				t.Errorf("%s baris %d %s (%s) = %q, harap %q", jalur, r+1, k.Properti, k.Kolom, dapat[r][k.Properti], harap[r][k.Properti])
			}
		}
	}
}

// periksaPulangPergi membandingkan SETIAP kolom katalog halaman tulis dan baca.
func periksaPulangPergi(t *testing.T, h, b *models.Halaman) {
	t.Helper()
	for _, k := range models.TabelGeneralPolis.Kolom {
		if got, harap := b.Ambil(k.Properti), h.Ambil(k.Properti); got != harap {
			t.Errorf("T_GENERAL_POLIS.%s (%s) = %q, harap %q", k.Kolom, k.Properti, got, harap)
		}
	}
	for _, k := range models.TabelQuotation.Kolom {
		harap := h.Ambil(models.HalamanPolis + ".QuotationData." + k.Properti)
		for _, j := range []string{models.HalamanPolis + ".QuotationData." + k.Properti, models.HalamanQuotation + "." + k.Properti} {
			if got := b.Ambil(j); got != harap {
				t.Errorf("T_POLIS_QUOTATION.%s (%s) = %q, harap %q", k.Kolom, j, got, harap)
			}
		}
	}
	for _, tb := range []models.Tabel{models.TabelCeding, models.TabelAngsuran, models.TabelSpreading, models.TabelXOL} {
		samakanDaftar(t, tb.Daftar, tb, h.AmbilDaftar(tb.Daftar), b.AmbilDaftar(tb.Daftar))
	}
	for _, p := range []struct {
		induk models.Tabel
		cucu  models.Tabel
	}{{models.TabelAngsuran, models.TabelAngsuranRinci}, {models.TabelXOL, models.TabelLayerXOL}} {
		for i := range h.AmbilDaftar(p.induk.Daftar) {
			j := models.JalurAnak(p.induk.Daftar, i+1, p.cucu.Daftar)
			samakanDaftar(t, j, p.cucu, h.AmbilDaftar(j), b.AmbilDaftar(j))
		}
	}
}

// Spec-penyimpanan AC 49: "Polis proporsional yang disimpan lalu dibaca kembali
// menghasilkan nilai yang sama pada seluruh medan. Test yang menemukan satu
// medan berbeda gagal." - SETIAP kolom katalog kelima tabel bentuk proporsional.
func TestPolisProporsionalSeluruhMedanPulangPergi(t *testing.T) {
	_, _, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-PP-PROP"
	h := halamanPenuh(models.JenisProporsional)
	simpanKasusBaru(t, ctx, d, g, id, h)
	b, err := g.BacaHalaman(ctx, nil, id)
	if err != nil {
		t.Fatal(err)
	}
	periksaPulangPergi(t, h, b)
}

// Spec-penyimpanan AC 50: "Polis non-proporsional yang disimpan lalu dibaca
// kembali menghasilkan nilai yang sama, termasuk seluruh layer. Test yang
// menemukan layer hilang gagal." - SETIAP kolom ketujuh tabel bentuk
// non-proporsional, dua XOL x tiga layer, dua angsuran x dua rincian.
func TestPolisNonProporsionalSeluruhLayerPulangPergi(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-PP-XOL"
	h := halamanPenuh(models.JenisNonProporsional)
	simpanKasusBaru(t, ctx, d, g, id, h)
	b, err := g.BacaHalaman(ctx, nil, id)
	if err != nil {
		t.Fatal(err)
	}
	periksaPulangPergi(t, h, b)
	// layer dihitung LANGSUNG di tabel: 2 x 3, tidak satu pun hilang
	n := kolomTeks(t, ctx, sqlDB, fmt.Sprintf(`SELECT TO_CHAR(COUNT(*)) FROM %s.T_POLIS_XOL_LAYER l
		JOIN %s.T_POLIS_XOL x ON x.ID = l.XOL_ID WHERE x.POLIS_ID = :1`, skema, skema), id)
	if n != "6" {
		t.Errorf("T_POLIS_XOL_LAYER memuat %s baris polis ini, harap 6", n)
	}
}

// Spec-penyimpanan AC 51: "Urutan baris anak saat dibaca sama dengan urutan
// NOURUT. Test yang menemukan urutan acak gagal." Dua belas baris spreading
// dengan TreatyType yang urutan teksnya BERBEDA dari urutan tulis (kunci baris
// ID acak), lalu NOURUT dibaca langsung: baris ke-j = NOURUT j.
func TestUrutanBarisAnakMenurutNourut(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-URUT"
	tipe := []string{"UJI-12", "UJI-03", "UJI-11", "UJI-01", "UJI-10", "UJI-02", "UJI-09", "UJI-04", "UJI-08", "UJI-05", "UJI-07", "UJI-06"}
	var sp []models.Baris
	for _, x := range tipe {
		sp = append(sp, models.Baris{"TreatyType": x})
	}
	h := models.HalamanBaru()
	h.Setel("PositionNote", models.PosisiAdmin)
	h.SetelDaftar(models.DaftarSpreading, sp)
	simpanKasusBaru(t, ctx, d, g, id, h)
	b, err := g.BacaHalaman(ctx, nil, id)
	if err != nil {
		t.Fatal(err)
	}
	var dapat []string
	for _, r := range b.AmbilDaftar(models.DaftarSpreading) {
		dapat = append(dapat, r["TreatyType"])
	}
	if fmt.Sprint(dapat) != fmt.Sprint(tipe) {
		t.Errorf("urutan terbaca %v, harap urutan tulis %v", dapat, tipe)
	}
	rows, err := sqlDB.QueryContext(ctx, fmt.Sprintf(`SELECT TO_CHAR(NOURUT), TREATY_TYPE FROM %s.T_POLIS_SPREADING
		WHERE POLIS_ID = :1 ORDER BY NOURUT`, skema), id)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	j := 0
	for rows.Next() {
		var no, tt string
		if err := rows.Scan(&no, &tt); err != nil {
			t.Fatal(err)
		}
		if no != strconv.Itoa(j+1) || j >= len(tipe) || tt != tipe[j] {
			t.Errorf("NOURUT %s = %s, harap NOURUT %d = %s", no, tt, j+1, tipe[min(j, len(tipe)-1)])
		}
		j++
	}
	if j != len(tipe) {
		t.Errorf("%d baris spreading, harap %d", j, len(tipe))
	}
}

// Spec AC 23: "Nilai uang disimpan berpresisi penuh. Test yang menemukan
// pembulatan di lapisan repository gagal." Kolom dibaca LANGSUNG dan lewat
// `BacaHalaman`: 38 digit kembali utuh; 11 desimal dibulatkan Oracle di
// desimal ke-9 (NUMBER(38,8), KEPUTUSAN 23-09-2026; spec-penyimpanan AC 55),
// bukan ke 2 desimal; 8 desimal kecil tidak menjadi 0.
func TestUangPresisiPenuhTanpaPembulatanRepository(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-PRESISI"
	h := models.HalamanBaru()
	h.Setel("PositionNote", models.PosisiAdmin)
	harap := map[string][2]string{ // properti -> {masuk, kolom}
		"PolicyTreatyIn.PremiOgp":   {uangPenuh, uangPenuh},
		"PolicyTreatyIn.PremiOnp":   {"830.82191780804", "830.82191781"},
		"PolicyTreatyIn.NetPremium": {"0.00000001", "0.00000001"},
		"PolicyTreatyIn.ResultOgp1": {"-592629512.880000276", "-592629512.88000028"},
	}
	for p, v := range harap {
		h.Setel(p, v[0])
	}
	simpanKasusBaru(t, ctx, d, g, id, h)
	kolom := map[string]string{"PolicyTreatyIn.PremiOgp": "PREMI_OGP", "PolicyTreatyIn.PremiOnp": "PREMI_ONP",
		"PolicyTreatyIn.NetPremium": "NET_PREMIUM", "PolicyTreatyIn.ResultOgp1": "RESULT_OGP1"}
	b, err := g.BacaHalaman(ctx, nil, id)
	if err != nil {
		t.Fatal(err)
	}
	for p, v := range harap {
		got := kolomTeks(t, ctx, sqlDB, fmt.Sprintf(`SELECT %s FROM %s.T_GENERAL_POLIS WHERE ID = :1`,
			fmt.Sprintf(intidb.FmtDesimal, kolom[p]), skema), id)
		if rapikan(got) != v[1] {
			t.Errorf("kolom %s = %s, harap %s (masuk %s)", kolom[p], got, v[1], v[0])
		}
		if b.Ambil(p) != v[1] {
			t.Errorf("%s terbaca %q, harap %q (masuk %s)", p, b.Ambil(p), v[1], v[0])
		}
	}
}

// rapikan - TO_CHAR TM9 menulis ".5"/"-.5" untuk |x| < 1.
func rapikan(s string) string {
	switch {
	case len(s) > 0 && s[0] == '.':
		return "0" + s
	case len(s) > 1 && s[:2] == "-.":
		return "-0" + s[1:]
	}
	return s
}

// Spec AC 29: "Seluruh urutan penyimpanan berada dalam satu transaksi. Test
// yang menyuntikkan kegagalan di tengah dan menemukan ada baris tersisa gagal."
// dan spec-penyimpanan AC 46 ("lebih dari satu komit" = baris tersisa sesudah
// pembatalan). Urutan simpan pertama `services` (SisipKasus lalu SimpanHalaman
// bentuk non-proporsional penuh - kesembilan tabel tersentuh) digagalkan
// SESUDAH seluruh tulisan, sebelum komit: tidak satu pun baris boleh tersisa.
func TestKegagalanDiTengahTidakMenyisakanBarisDiTabelManaPun(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-BATAL-PENUH"
	err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		if err := g.SisipKasus(ctx, tx, id, "UJI-AKUN", "UJI NAMA"); err != nil {
			return err
		}
		if err := g.SimpanHalaman(ctx, tx, id, halamanPenuh(models.JenisNonProporsional)); err != nil {
			return err
		}
		return fmt.Errorf("UJI-gagal sesudah seluruh tulisan")
	})
	if err == nil || err.Error() != "UJI-gagal sesudah seluruh tulisan" {
		t.Fatalf("harap galat suntikan, dapat %v", err)
	}
	for tabel, syarat := range map[string]string{
		"T_WORK_POLIS":       "ID = :1",
		"T_GENERAL_POLIS":    "ID = :1",
		"T_POLIS_QUOTATION":  "POLIS_ID = :1",
		"T_POLIS_CEDING":     "QUOTATION_ID = :1",
		"T_POLIS_INSTALMENT": "POLIS_ID = :1",
		"T_POLIS_SPREADING":  "POLIS_ID = :1",
		"T_POLIS_XOL":        "POLIS_ID = :1",
		// cucu berkunci induk acak: seluruh tabel dihitung (skema uji dipasang
		// bersih per uji; uji ini satu-satunya penulisnya)
		"T_POLIS_INSTALMENT_DETAIL": "",
		"T_POLIS_XOL_LAYER":         "",
	} {
		q, args := fmt.Sprintf(`SELECT TO_CHAR(COUNT(*)) FROM %s.%s`, skema, tabel), []any{}
		if syarat != "" {
			q, args = q+" WHERE "+syarat, []any{id}
		}
		if n := kolomTeks(t, ctx, sqlDB, q, args...); n != "0" {
			t.Errorf("%s: %s baris tersisa sesudah pembatalan", tabel, n)
		}
	}
}
