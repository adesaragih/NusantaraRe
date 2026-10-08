package repository

// Jalur SIMPAN layar ADJUSTMENT (EDM) — tombol `Save`, `Submit`, `Actions`,
// dan `Decline offer` modul Treaty In Adjustment.
//
// ---------------------------------------------------------------------
// ⭐ MENGAPA DI MODUL INI (7 Oktober 2026)
// ---------------------------------------------------------------------
// Dokumen Adjustment mendarat di tabel `T_TREATY_*` yang SAMA dengan Treaty
// In (`MuatPenyesuaian`, `AkhiranSisiLama`), dan penulisnya hidup di sini.
// Penjaga arsitektur melarang modul mengimpor modul lain, jadi layar
// Adjustment memanggil rute `/api/treaty-in/penyesuaian/*` — pola yang
// sudah ia pakai untuk `/api/treaty-in/hitung/*`.
//
// ---------------------------------------------------------------------
// Sasaran tulis, SATU transaksi
// ---------------------------------------------------------------------
//
//	T_TREATY_*     sisi New (`MASTERID = ID`) dan — draf saja — sisi Old
//	               (`ID + AkhiranSisiLama`), lewat `MuatKontrakSebagian`
//	TREATY_IN_EDM  kepala — kolom prosedur `PEGA_M_TREATY_IN_EDM`
//	               (`STSINPUT='0'` → INSERT, selainnya UPDATE semua kolom
//	               kecuali ID); `EDMDATE = SYSDATE`
//
// ⛔ Dokumen JSON EDM TIDAK disentuh — prosedur Pega menulisnya, jalur ini
// tidak (keputusan pemilik proses: nol akses ke dokumen JSON). ⛔ Detail EDM
// saat Resolve Complete (`SaveTreatyInDetailEdm_Act`) tidak dibangun —
// pemilik proses: data ditarik dari tabel tiap tab, status dari kepalanya.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// TabelKepalaPenyesuaian - kepala penyesuaian warisan.
const TabelKepalaPenyesuaian = "TREATY_IN_EDM"

// ErrPenyesuaianSudahAda - draf berpengenal yang SUDAH tersimpan. Prosedur
// Pega menimpanya (`STSINPUT` ≠ 0 → UPDATE); di sini ditolak, sebab
// menimpa penyesuaian lain tanpa suara adalah kehilangan data.
var ErrPenyesuaianSudahAda = errors.New("repository: pengenal penyesuaian sudah ada")

// ErrBukanPengenalPenyesuaian - pengenal tanpa `/R`. Penjaga penghapusan:
// `Decline offer` EDM tidak boleh pernah menyentuh pendaratan KONTRAK.
var ErrBukanPengenalPenyesuaian = errors.New("repository: bukan pengenal penyesuaian")

// KolomKepalaPenyesuaian - kolom `TREATY_IN_EDM` ↔ properti `TreatyIn`,
// selain `ID` dan `EDMDATE`: `OLDID`, kesembilan belas kolom kepala Treaty
// In, lalu dua milik EDM.
var KolomKepalaPenyesuaian = append(append([][2]string{{"OLDID", "OLDID"}}, KolomKepalaTreatyIn...),
	[2]string{"EDMSTATE", "EDMState"},
	[2]string{"EDMMATERIALTYPE", "EDMMaterialType"},
)

// Bentuk `EDMDATE` tersimpan — `SYSDATE` ke `VARCHAR2` dengan
// `NLS_DATE_FORMAT` sesi Pega (`02-FEB-21`, 280 dari 280 baris). Bahasanya
// dikunci supaya nama bulan tidak bergantung sesi.
const ekspresiTanggalEDM = "TO_CHAR(SYSDATE, 'DD-MON-RR', 'NLS_DATE_LANGUAGE=AMERICAN')"

func pengenalPenyesuaian(id string) bool {
	return strings.Contains(id, "/R")
}

func wadahTeks(n int) ([]sql.NullString, []any) {
	sel := make([]sql.NullString, n)
	tuju := make([]any, n)
	for i := range sel {
		tuju[i] = &sel[i]
	}
	return sel, tuju
}

// BacaKepalaPenyesuaian - kepala satu penyesuaian di `TREATY_IN_EDM`,
// berkunci properti Pega. Kolom `NULL` tidak dimasukkan.
func (g *Gudang) BacaKepalaPenyesuaian(ctx context.Context, id string) (map[string]any, bool, error) {
	nama, err := g.db.Qualify(TabelKepalaPenyesuaian)
	if err != nil {
		return nil, false, err
	}
	kolom := make([]string, len(KolomKepalaPenyesuaian))
	for i, k := range KolomKepalaPenyesuaian {
		kolom[i] = k[0]
	}
	// ⚠️ Tabelnya tanpa kunci unik — baris pertama yang dipakai.
	q := fmt.Sprintf("SELECT %s FROM %s WHERE ID = :1 FETCH FIRST 1 ROWS ONLY", strings.Join(kolom, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, false, err
	}
	sel, tuju := wadahTeks(len(kolom))
	if err := g.db.QueryRowContext(ctx, q, id).Scan(tuju...); err == sql.ErrNoRows {
		return map[string]any{}, false, nil
	} else if err != nil {
		return nil, false, fmt.Errorf("repository: membaca kepala %s %s: %w", TabelKepalaPenyesuaian, id, err)
	}
	out := map[string]any{}
	for i, k := range KolomKepalaPenyesuaian {
		if sel[i].Valid {
			out[k[1]] = sel[i].String
		}
	}
	return out, true, nil
}

// SimpanPenyesuaian menulis seluruh rencana dalam SATU transaksi.
func (g *Gudang) SimpanPenyesuaian(ctx context.Context, r models.RencanaPenyesuaian) (err error) {
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if err = g.simpanPenyesuaianDalam(ctx, tx, r); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("repository: mengikat simpanan penyesuaian %s: %w", r.ID, err)
	}
	return nil
}

// SimpanPenyesuaianLaluBatalkanUntukUji - seluruh tulisan di dalam transaksi
// yang SELALU dibatalkan; mengembalikan cacah baris kedua sisi dan kepala
// yang sempat tertulis. Hanya untuk uji `db`.
func (g *Gudang) SimpanPenyesuaianLaluBatalkanUntukUji(ctx context.Context, r models.RencanaPenyesuaian) (baru, lama map[string]int, kepala map[string]any, err error) {
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	defer func() { _ = tx.Rollback() }()
	if err = g.simpanPenyesuaianDalam(ctx, tx, r); err != nil {
		return nil, nil, nil, err
	}
	if baru, err = g.CacahBarisKontrak(ctx, tx, r.ID); err != nil {
		return nil, nil, nil, err
	}
	if lama, err = g.CacahBarisKontrak(ctx, tx, r.ID+AkhiranSisiLama); err != nil {
		return nil, nil, nil, err
	}
	kepala, err = g.kepalaPenyesuaianDalam(ctx, tx, r.ID)
	return baru, lama, kepala, err
}

func (g *Gudang) simpanPenyesuaianDalam(ctx context.Context, tx *db.Tx, r models.RencanaPenyesuaian) error {
	id := strings.TrimSpace(r.ID)
	if !pengenalPenyesuaian(id) {
		return fmt.Errorf("%w: %q", ErrBukanPengenalPenyesuaian, id)
	}
	ada, err := g.adaKepalaPenyesuaian(ctx, tx, id)
	if err != nil {
		return err
	}
	if r.Draf && ada {
		return fmt.Errorf("%w: %s", ErrPenyesuaianSudahAda, id)
	}
	doc := r.Baru
	if doc == nil {
		doc = map[string]any{}
	}
	doc["ID"] = id
	// ⭐ `Sebagian`: tabel yang dokumen tidak bawa tidak disentuh — peta
	// layar Adjustment boleh lebih sempit dari peta ini.
	if _, err := g.MuatKontrakSebagian(ctx, tx, id, doc); err != nil {
		return err
	}
	if r.Lama != nil {
		if _, err := g.MuatKontrakSebagian(ctx, tx, id+AkhiranSisiLama, r.Lama); err != nil {
			return err
		}
	}
	return g.tulisKepalaPenyesuaian(ctx, tx, id, doc, !ada)
}

func (g *Gudang) adaKepalaPenyesuaian(ctx context.Context, tx *db.Tx, id string) (bool, error) {
	nama, err := g.db.Qualify(TabelKepalaPenyesuaian)
	if err != nil {
		return false, err
	}
	// `EDMCheckExistingData`: `select count(ID) … where id = {TreatyIn.ID}`.
	q := fmt.Sprintf("SELECT COUNT(ID) FROM %s WHERE ID = :1", nama)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := tx.QueryRowContext(ctx, q, id).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: memeriksa %s %s: %w", TabelKepalaPenyesuaian, id, err)
	}
	return n > 0, nil
}

func (g *Gudang) kepalaPenyesuaianDalam(ctx context.Context, tx *db.Tx, id string) (map[string]any, error) {
	nama, err := g.db.Qualify(TabelKepalaPenyesuaian)
	if err != nil {
		return nil, err
	}
	kolom := []string{"EDMDATE"}
	for _, k := range KolomKepalaPenyesuaian {
		kolom = append(kolom, k[0])
	}
	q := fmt.Sprintf("SELECT %s FROM %s WHERE ID = :1 FETCH FIRST 1 ROWS ONLY", strings.Join(kolom, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	sel, tuju := wadahTeks(len(kolom))
	if err := tx.QueryRowContext(ctx, q, id).Scan(tuju...); err != nil {
		return nil, fmt.Errorf("repository: membaca kepala %s %s: %w", TabelKepalaPenyesuaian, id, err)
	}
	out := map[string]any{"EDMDATE": sel[0].String}
	for i, k := range KolomKepalaPenyesuaian {
		if sel[i+1].Valid {
			out[k[1]] = sel[i+1].String
		}
	}
	return out, nil
}

// tulisKepalaPenyesuaian - `INSERT` (baris belum ada) atau `UPDATE` semua
// kolom kecuali `ID`, seperti `PEGA_M_TREATY_IN_EDM`.
func (g *Gudang) tulisKepalaPenyesuaian(ctx context.Context, tx *db.Tx, id string, doc map[string]any, baru bool) error {
	nama, err := g.db.Qualify(TabelKepalaPenyesuaian)
	if err != nil {
		return err
	}
	nilai := func(k string) any {
		v, ada := NilaiTeks(nilaiJalur(doc, k))
		if !ada {
			return nil
		}
		return v
	}
	args := []any{}
	var q string
	if baru {
		kolom := []string{"ID", "EDMDATE"}
		args = append(args, id)
		tanda := []string{":1", ekspresiTanggalEDM}
		for _, k := range KolomKepalaPenyesuaian {
			kolom = append(kolom, k[0])
			args = append(args, nilai(k[1]))
			tanda = append(tanda, fmt.Sprintf(":%d", len(args)))
		}
		q = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", nama, strings.Join(kolom, ", "), strings.Join(tanda, ", "))
	} else {
		set := []string{"EDMDATE = " + ekspresiTanggalEDM}
		for _, k := range KolomKepalaPenyesuaian {
			args = append(args, nilai(k[1]))
			set = append(set, fmt.Sprintf("%s = :%d", k[0], len(args)))
		}
		args = append(args, id)
		q = fmt.Sprintf("UPDATE %s SET %s WHERE ID = :%d", nama, strings.Join(set, ", "), len(args))
	}
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("repository: menulis kepala %s %s: %w", TabelKepalaPenyesuaian, id, err)
	}
	return nil
}

// HapusPenyesuaian - `Decline offer` EDM (`TreatyInDeclineConfirmation_
// postactEDM` [6]–[7]): baris kepala DIHAPUS FISIK, beserta pendaratan
// kedua sisinya — padanan dokumen JSON yang Pega hapus di langkah [6].
//
// ⛔ Hanya pengenal penyesuaian (`…/R…`): pendaratan KONTRAK Treaty In
// tidak pernah terhapus lewat tombol ini.
func (g *Gudang) HapusPenyesuaian(ctx context.Context, id string) (err error) {
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if err = g.hapusPenyesuaianDalam(ctx, tx, id); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("repository: mengikat penghapusan penyesuaian %s: %w", id, err)
	}
	return nil
}

// HapusPenyesuaianLaluBatalkanUntukUji - penghapusan di dalam transaksi yang
// SELALU dibatalkan; mengembalikan apakah kepalanya masih ada sesudahnya.
func (g *Gudang) HapusPenyesuaianLaluBatalkanUntukUji(ctx context.Context, id string) (masihAda bool, err error) {
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback() }()
	if err = g.hapusPenyesuaianDalam(ctx, tx, id); err != nil {
		return false, err
	}
	return g.adaKepalaPenyesuaian(ctx, tx, id)
}

func (g *Gudang) hapusPenyesuaianDalam(ctx context.Context, tx *db.Tx, id string) error {
	id = strings.TrimSpace(id)
	if !pengenalPenyesuaian(id) {
		return fmt.Errorf("%w: %q", ErrBukanPengenalPenyesuaian, id)
	}
	nama, err := g.db.Qualify(TabelKepalaPenyesuaian)
	if err != nil {
		return err
	}
	q := fmt.Sprintf("DELETE FROM %s WHERE ID = :1", nama)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("repository: menghapus %s %s: %w", TabelKepalaPenyesuaian, id, err)
	}
	for _, m := range []string{id, id + AkhiranSisiLama} {
		if _, err := g.KosongkanKontrak(ctx, tx, m); err != nil {
			return err
		}
	}
	return nil
}
