package repository

// Pindah JSON -> flat (keputusan work owner 06-10-2026 butir 3; alat `backend/alat/pindahflat`, pola
// masterproductnamelife): `M_RICOMM_LIFE` (JSON, 0 baris DEV) -> tabel flat `RICOMM_LIFE` (migrasi inti 924).
//
//	-uji       (bawaan) hanya SELECT; nol tulisan; laporan AGREGAT (ID baris dan nama kolom, tidak pernah nilainya).
//	-jalankan  ditolak bila IS_PEGA_PROD=true (ErrPindahDiProduksi), bila ada kegagalan, bila tabel flat memuat ID yang
//	           sama dengan isi BERBEDA (tulisan aplikasi tidak pernah ditimpa), dan bila ada normalisasi teks angka yang
//	           belum diterima (`-terima-normalisasi`). SATU transaksi yang lebih dulu MENGUNCI tabel flat; rencana
//	           dihitung ulang di dalamnya; baris yang sudah ada dan sama dilewati - aman diulang.
//
// M_RICOMM_LIFE tidak pernah ditulis.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/ricommlife/backend/models"
)

var (
	// ErrPindahDiProduksi - `-jalankan` di lingkungan IS_PEGA_PROD=true.
	ErrPindahDiProduksi = errors.New("repository: moving R/I Comm Life rows to the flat table is refused when IS_PEGA_PROD=true")
	// ErrPindahTidakLolos - rencana memuat kegagalan, baris flat berbeda, atau normalisasi yang belum diterima.
	ErrPindahTidakLolos = errors.New("repository: the move plan did not pass; nothing was written")
)

// PeriksaMode - `-jalankan` ditolak bila IS_PEGA_PROD=true; uji kering selalu boleh.
func PeriksaMode(jalankan, pegaProduksi bool) error {
	if jalankan && pegaProduksi {
		return ErrPindahDiProduksi
	}
	return nil
}

// BarisJSON - satu baris `M_RICOMM_LIFE`.
type BarisJSON struct{ ID, JSON string }

// LaporanPindah - laporan agregat satu putaran.
type LaporanPindah struct {
	Mode   string
	Sumber int
	// Ditulis* - baris yang (akan) disisipkan.
	AkanDitulis int
	// SudahSama - ID sumber yang sudah ada di tabel flat dengan isi sama (dilewati).
	SudahSama int
	// FlatSaja - baris tabel flat tanpa pasangan sumber (tulisan aplikasi) - dibiarkan.
	FlatSaja int
	// Berbeda - ID sumber yang sudah ada di tabel flat dengan isi BERBEDA (menahan -jalankan).
	Berbeda []string
	// Gagal - "ID KOLOM: jenis" (tanpa nilai).
	Gagal []string
	// Kosong - kolom -> cacah nilai kosong / tidak ada (menjadi NULL).
	Kosong map[string]int
	// Normalisasi - kolom -> cacah teks angka yang bentuknya berubah (mis. `05` -> `5`, `0,5` -> `0.5`).
	Normalisasi map[string]int
	// TerimaNormalisasi - operator menerima normalisasi (`-terima-normalisasi`).
	TerimaNormalisasi bool
	Ditulis           bool
}

// BolehDitulis - nol kegagalan, nol baris flat berbeda, dan normalisasi (bila ada) sudah diterima.
func (l LaporanPindah) BolehDitulis() bool {
	return len(l.Gagal) == 0 && len(l.Berbeda) == 0 && (len(l.Normalisasi) == 0 || l.TerimaNormalisasi)
}

func urutKunci(m map[string]int) string {
	var k []string
	for x, n := range m {
		k = append(k, fmt.Sprintf("%s=%d", x, n))
	}
	sort.Strings(k)
	if len(k) == 0 {
		return "-"
	}
	return strings.Join(k, ", ")
}

// Teks - laporan untuk operator (agregat).
func (l LaporanPindah) Teks() string {
	var b strings.Builder
	fmt.Fprintf(&b, "pindahflat R/I Comm Life (%s)\n", l.Mode)
	fmt.Fprintf(&b, "  sumber M_RICOMM_LIFE: %d\n", l.Sumber)
	fmt.Fprintf(&b, "  baris yang (akan) ditulis: %d\n", l.AkanDitulis)
	fmt.Fprintf(&b, "  sudah ada dan sama (dilewati): %d\n", l.SudahSama)
	fmt.Fprintf(&b, "  hanya di tabel flat (dibiarkan): %d\n", l.FlatSaja)
	fmt.Fprintf(&b, "  kosong menjadi NULL: %s\n", urutKunci(l.Kosong))
	fmt.Fprintf(&b, "  normalisasi angka: %s (diterima: %v)\n", urutKunci(l.Normalisasi), l.TerimaNormalisasi)
	fmt.Fprintf(&b, "  tabel flat berbeda: %d %s\n", len(l.Berbeda), strings.Join(l.Berbeda, " "))
	fmt.Fprintf(&b, "  gagal: %d\n", len(l.Gagal))
	for _, g := range l.Gagal {
		fmt.Fprintf(&b, "    %s\n", g)
	}
	fmt.Fprintf(&b, "  ditulis: %v\n", l.Ditulis)
	return b.String()
}

// konversiAngka - nilai JSON -> kanonik flat; jenis "bulat5" / "tahun" / "desimal". Kosong = "" tanpa galat.
func konversiAngka(kolom, nilai string) (kanonik string, normal bool, err error) {
	if nilai == "" {
		return "", false, nil
	}
	switch kolom {
	case JSONContract:
		kanonik, err = models.NormalContract(nilai)
	case JSONYear:
		kanonik, err = models.NormalYear(nilai)
	default:
		kanonik, err = models.NormalComm(nilai)
	}
	if err != nil {
		return "", false, err
	}
	return kanonik, kanonik != nilai, nil
}

// KonversiBaris - satu baris JSON -> baris flat; `gagal` berisi nama kolom dan jenisnya, tanpa nilai.
func KonversiBaris(b BarisJSON, lap *LaporanPindah) (models.Komisi, []string) {
	id := strings.TrimSpace(b.ID)
	k := models.Komisi{ID: id}
	var gagal []string
	if id == "" || len(id) > models.BatasID {
		gagal = append(gagal, fmt.Sprintf("ID: kosong atau lebih dari %d karakter", models.BatasID))
	}
	if strings.TrimSpace(b.JSON) != "" {
		if _, err := TerapkanKunci(b.JSON, nil); err != nil {
			return k, append(gagal, "JSONDATA: bukan objek JSON")
		}
	}
	ambil := func(kunci string) string {
		v, ada := TeksKunci(b.JSON, kunci)
		if !ada || v == "" {
			lap.Kosong[kunci]++
		}
		return v
	}
	k.IDUsedBy = ambil(JSONIDUsedBy)
	if len(k.IDUsedBy) > models.BatasID {
		gagal = append(gagal, fmt.Sprintf("IDUSEDBY: lebih dari %d karakter", models.BatasID))
	}
	k.UsedBy = ambil(JSONUsedBy)
	if len(k.UsedBy) > models.BatasNama {
		gagal = append(gagal, fmt.Sprintf("USEDBY: lebih dari %d byte", models.BatasNama))
	}
	for _, x := range []struct {
		kunci string
		ke    *string
	}{{JSONContract, &k.Contract}, {JSONYear, &k.Year}, {JSONComm, &k.Comm}} {
		nilai := ambil(x.kunci)
		kanonik, normal, err := konversiAngka(x.kunci, nilai)
		if err != nil {
			gagal = append(gagal, x.kunci+": tidak muat tipe kolom flat")
			continue
		}
		if normal {
			lap.Normalisasi[x.kunci]++
		}
		*x.ke = kanonik
	}
	return k, gagal
}

// RencanaPindah - laporan dan baris yang akan disisipkan, dari sumber JSON dan isi tabel flat saat ini. Murni.
func RencanaPindah(sumber []BarisJSON, flat []models.Komisi, terimaNormalisasi bool) (LaporanPindah, []models.Komisi) {
	lap := LaporanPindah{Sumber: len(sumber), Kosong: map[string]int{}, Normalisasi: map[string]int{},
		TerimaNormalisasi: terimaNormalisasi}
	ada := map[string]models.Komisi{}
	for _, f := range flat {
		ada[f.ID] = f
	}
	dipakai := map[string]bool{}
	var tulis []models.Komisi
	for _, b := range sumber {
		k, gagal := KonversiBaris(b, &lap)
		for _, g := range gagal {
			lap.Gagal = append(lap.Gagal, strings.TrimSpace(b.ID)+" "+g)
		}
		if len(gagal) > 0 {
			continue
		}
		if dipakai[k.ID] {
			lap.Gagal = append(lap.Gagal, k.ID+" ID: kembar di sumber")
			continue
		}
		dipakai[k.ID] = true
		if f, sudah := ada[k.ID]; sudah {
			if f == k {
				lap.SudahSama++
			} else {
				lap.Berbeda = append(lap.Berbeda, k.ID)
			}
			continue
		}
		tulis = append(tulis, k)
	}
	for id := range ada {
		if !dipakai[id] {
			lap.FlatSaja++
		}
	}
	lap.AkanDitulis = len(tulis)
	return lap, tulis
}

// bacaSumber - seluruh baris JSON lama.
func (g *Gudang) bacaSumber(ctx context.Context, tx *db.Tx, n nama) ([]BarisJSON, error) {
	q := SqlSumberJSON(n.tabelLama)
	if err := siap(TabelJSONLama, q); err != nil {
		return nil, err
	}
	rows, err := g.dari(tx).QueryContext(ctx, q)
	if err != nil {
		return nil, bungkus(err, "membaca M_RICOMM_LIFE")
	}
	defer func() { _ = rows.Close() }()
	var out []BarisJSON
	for rows.Next() {
		var id sql.NullString
		var p db.PindaiTeksPanjang
		if err := rows.Scan(&id, &p); err != nil {
			return nil, bungkus(err, "memindai M_RICOMM_LIFE")
		}
		out = append(out, BarisJSON{ID: id.String, JSON: p.Teks()})
	}
	return out, bungkus(rows.Err(), "membaca M_RICOMM_LIFE")
}

func mode(jalankan bool) string {
	if jalankan {
		return "jalankan"
	}
	return "uji"
}

// SqlSesiNLS - setelan sesi alat pindah: titik desimal. Pembacaan dan penulisan angka sudah tidak bergantung NLS
// (`fmtAngka` berargumen NLS, AngkaOracle di Go, `TO_NUMBER` atas teks angka murni); setelan ini lapis kedua untuk
// konversi implisit, dan karena itu WAJIB berlaku di koneksi yang sama dengan transaksinya (SesiPindah).
const SqlSesiNLS = `ALTER SESSION SET NLS_NUMERIC_CHARACTERS = '.,'`

// SesiPindah - SATU koneksi terkunci (`db.Koneksi`, `*sql.Conn`): ALTER SESSION lalu transaksi di koneksi yang SAMA.
// Pemanggil menutup keduanya (Rollback / Commit, lalu Close).
func (g *Gudang) SesiPindah(ctx context.Context) (*db.Koneksi, *db.Tx, error) {
	if err := siap("SESI", SqlSesiNLS); err != nil {
		return nil, nil, err
	}
	kon, err := g.db.Koneksi(ctx)
	if err != nil {
		return nil, nil, err
	}
	if _, err := kon.ExecContext(ctx, SqlSesiNLS); err != nil {
		_ = kon.Close()
		return nil, nil, bungkus(err, "menyetel sesi pindah")
	}
	tx, err := kon.Mulai(ctx)
	if err != nil {
		_ = kon.Close()
		return nil, nil, err
	}
	return kon, tx, nil
}

// PindahFlat - satu putaran alat pindah (lihat kepala berkas): SELURUHNYA di satu koneksi (SesiPindah).
func (g *Gudang) PindahFlat(ctx context.Context, jalankan, terimaNormalisasi bool) (LaporanPindah, error) {
	n, err := g.nama()
	if err != nil {
		return LaporanPindah{}, err
	}
	if err := PeriksaMode(jalankan, g.db.PegaProduksi()); err != nil {
		return LaporanPindah{Mode: "jalankan"}, err
	}
	kon, tx, err := g.SesiPindah(ctx)
	if err != nil {
		return LaporanPindah{Mode: mode(jalankan)}, err
	}
	defer func() { _ = kon.Close() }()
	defer func() { _ = tx.Rollback() }()
	rencana := func(tx *db.Tx) (LaporanPindah, []models.Komisi, error) {
		sumber, err := g.bacaSumber(ctx, tx, n)
		if err != nil {
			return LaporanPindah{}, nil, err
		}
		flat, err := g.bacaBaris(ctx, tx, TabelKomisi, SqlSemuaKomisi(n.tabelKomisi), 6)
		if err != nil {
			return LaporanPindah{}, nil, err
		}
		lap, tulis := RencanaPindah(sumber, keKomisi(flat), terimaNormalisasi)
		return lap, tulis, nil
	}
	if !jalankan {
		// Uji kering: hanya SELECT di dalam transaksi koneksi ini; ditutup Rollback.
		lap, _, err := rencana(tx)
		lap.Mode = "uji"
		return lap, err
	}
	if _, err := g.tulis(ctx, tx, TabelKomisi, SqlKunciKomisi(n.tabelKomisi), "mengunci RICOMM_LIFE"); err != nil {
		return LaporanPindah{Mode: "jalankan"}, err
	}
	lap, tulis, err := rencana(tx)
	lap.Mode = "jalankan"
	if err != nil {
		return lap, err
	}
	if !lap.BolehDitulis() {
		return lap, ErrPindahTidakLolos
	}
	for _, k := range tulis {
		if err := g.SisipKomisi(ctx, tx, k); err != nil {
			return lap, err
		}
	}
	if err := tx.Commit(); err != nil {
		return lap, fmt.Errorf("repository: menutup transaksi pindah: %w", err)
	}
	lap.Ditulis = true
	return lap, nil
}
