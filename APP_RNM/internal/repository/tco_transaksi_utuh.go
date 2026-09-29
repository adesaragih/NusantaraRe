package repository

// Transaksi utuh satu kontrak - tiket 09 Treaty Contract Out.
//
// Untuk apa berkas ini: dua perkakas yang membuat penyimpanan lintas enam
// tabel dapat dijalankan dalam SATU transaksi dengan penulis per baris yang
// sudah ada (tiket 04-08), tanpa menyalin aturannya:
//
//  1. PEMBACA SADAR-TRANSAKSI. Selama `DenganTransaksiUtuhTCO`, setiap pembaca
//     modul (`bacaTCO`) membaca lewat transaksi itu - baris yang baru ditulis
//     permintaan yang sama (kontrak baru, reinsurer baru, induk klausul baru)
//     terlihat oleh langkah berikutnya. Di luar mode itu tetap pool biasa.
//  2. IDENTITAS SEMENTARA. Selama mode itu `IdentitasBerikutTCO` TIDAK
//     menyentuh sequence: ia memberi `S#########T` (lebar tetap, bukan angka,
//     tidak pernah bertabrakan dengan identitas '1' + digit). Baru sesudah
//     SELURUH baris lolos, `TetapkanIdentitasTCO` mengambil nomor sequence dan
//     mengganti identitas sementara. Kegagalan di baris ke-N karena itu tidak
//     menghabiskan satu nomor pun (AC tiket 09: "tidak menyisakan identitas
//     terpakai").
//
// `[data DBA]` keenam prosedur penulis tidak COMMIT sendiri - tidak ada titik
// potong transaksi di sisi basis data. Prosedur tidak dipanggil (keputusan o);
// commit sekali oleh `Service.DalamTransaksi` pemanggil.
//
// Dibaca sesudah: tco_identitas.go, tco_jejak.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// kolomReinsurerTCO - kolom `T_TREATYREINSURER` (migrasi 302), dipakai
// penggantian identitas reinsurer yang sudah punya security (FK tanpa ON UPDATE).
var kolomReinsurerTCO = []string{"ID", "TREATYYEAR", "TREATYGROUPID", "TREATYGROUPNAME", "REINSTYPEID",
	"REINSTYPENAME", "REINSURERID", "CLIENTID", "NAME", "RICOMM", "PCTSHARE", "IUDATE", "USERID", "STARTDATE",
	"ENDDATE", "STATUSON", "STDRATING", "OPERATORNAME", "TGLUPDATE"}

// kuerierTCO - bagian `*sql.DB` / `*sql.Tx` yang dipakai pembaca modul.
type kuerierTCO interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type kunciUtuhTCO struct{}

// utuhTCO adalah keadaan satu transaksi utuh.
type utuhTCO struct {
	tx       *Tx
	urut     int
	sementar []barisSementaraTCO
	tetap    bool
}

type barisSementaraTCO struct{ sequence, sementara string }

// ErrIdentitasUtuhTCO - penetapan identitas sementara gagal atau dipanggil keliru.
var ErrIdentitasUtuhTCO = errors.New("repository: penetapan identitas transaksi utuh gagal")

// polaIdentitasSementaraTCO - bentuk identitas sementara.
var polaIdentitasSementaraTCO = regexp.MustCompile(`S\d{9}T`)

// PolaIdentitasSementaraTCO - untuk menyamarkan identitas sementara di pesan galat.
func PolaIdentitasSementaraTCO() *regexp.Regexp { return polaIdentitasSementaraTCO }

// DenganTransaksiUtuhTCO menandai ctx: pembaca modul memakai tx, identitas
// baru sementara. tx nil (uji tanpa Oracle) mengembalikan ctx apa adanya.
func DenganTransaksiUtuhTCO(ctx context.Context, tx *Tx) context.Context {
	if tx == nil {
		return ctx
	}
	return context.WithValue(ctx, kunciUtuhTCO{}, &utuhTCO{tx: tx})
}

// DenganBacaTxTCO menandai ctx: pembaca modul memakai tx, identitas TETAP
// seperti biasa (tiket 10: hitung ulang dampak di dalam transaksi hapus).
func DenganBacaTxTCO(ctx context.Context, tx *Tx) context.Context {
	if tx == nil {
		return ctx
	}
	return context.WithValue(ctx, kunciUtuhTCO{}, &utuhTCO{tx: tx, tetap: true})
}

func utuhDari(ctx context.Context) *utuhTCO {
	u, _ := ctx.Value(kunciUtuhTCO{}).(*utuhTCO)
	return u
}

// bacaTCO memilih jalur baca: transaksi utuh bila ada, pool bila tidak.
func (d *DB) bacaTCO(ctx context.Context) kuerierTCO {
	if u := utuhDari(ctx); u != nil && u.tx != nil && u.tx.tx != nil {
		return u.tx.tx
	}
	return d.sql
}

// identitasSementaraTCO memberi identitas sementara bila ctx dalam mode utuh.
func identitasSementaraTCO(ctx context.Context, sequence string) (string, bool) {
	u := utuhDari(ctx)
	if u == nil || u.tetap {
		return "", false
	}
	u.urut++
	s := fmt.Sprintf("S%09dT", u.urut)
	u.sementar = append(u.sementar, barisSementaraTCO{sequence: sequence, sementara: s})
	return s, true
}

// tabelPerSequenceTCO - tabel yang identitasnya diterbitkan sequence itu.
var tabelPerSequenceTCO = map[string]string{
	SeqKontrakTCO:   TabelKontrakTCO,
	SeqReinsurerTCO: TabelReinsurerTCO,
	SeqSecurityTCO:  TabelSecurityTCO,
	SeqBusinessTCO:  TabelBusinessTCO,
	SeqKlausulTCO:   TabelKlausulTCO,
	SeqJejakTCO:     TabelJejakTCO,
}

func sqlGantiIdentitasTCO(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET ID = :1 WHERE ID = :2`, tabel)
}

// sqlSalinReinsurerTCO - salinan baris reinsurer dengan identitas tetap.
func sqlSalinReinsurerTCO(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (%s) SELECT :1, %s FROM %s WHERE ID = :2`, tabel,
		strings.Join(kolomReinsurerTCO, ", "), strings.Join(kolomReinsurerTCO[1:], ", "), tabel)
}

func sqlAlihSecurityTCO(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET REAS_ID = :1 WHERE REAS_ID = :2`, tabel)
}

func sqlBuangReinsurerSementaraTCO(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1`, tabel)
}

func sqlGantiBarisJejakTCO(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET BARIS_ID = :1 WHERE TABEL = :2 AND BARIS_ID = :3`, tabel)
}

// TetapkanIdentitasTCO mengambil nomor sequence untuk setiap identitas
// sementara - SEKALI, sesudah seluruh baris lolos - dan menggantinya di baris
// dan di jejaknya. Mengembalikan peta sementara -> tetap.
func (d *DB) TetapkanIdentitasTCO(ctx context.Context, tx *Tx) (map[string]string, error) {
	u := utuhDari(ctx)
	if u == nil || tx == nil || u.tx != tx {
		return nil, fmt.Errorf("%w: bukan transaksi utuh yang sama", ErrIdentitasUtuhTCO)
	}
	if u.tetap {
		return nil, fmt.Errorf("%w: identitas sudah ditetapkan", ErrIdentitasUtuhTCO)
	}
	u.tetap = true // sejak sini IdentitasBerikutTCO kembali ke sequence
	peta := make(map[string]string, len(u.sementar))
	for _, b := range u.sementar {
		tabel, dikenal := tabelPerSequenceTCO[b.sequence]
		if !dikenal {
			return nil, fmt.Errorf("%w: sequence %s tidak didukung transaksi utuh", ErrIdentitasUtuhTCO, b.sequence)
		}
		tetap, err := d.IdentitasBerikutTCO(ctx, tx, b.sequence)
		if err != nil {
			return nil, err
		}
		peta[b.sementara] = tetap
		if err := d.gantiIdentitasTCO(ctx, tx, tabel, b.sementara, tetap); err != nil {
			return nil, err
		}
	}
	// Jejak menyebut baris yang disentuhnya: ganti identitasnya juga.
	jejak, err := d.Qualify(TabelJejakTCO)
	if err != nil {
		return nil, err
	}
	q := sqlGantiBarisJejakTCO(jejak)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	for _, b := range u.sementar {
		tabel := tabelPerSequenceTCO[b.sequence]
		if tabel == TabelJejakTCO {
			continue
		}
		if _, err := tx.tx.ExecContext(ctx, q, peta[b.sementara], tabel, b.sementara); err != nil {
			return nil, fmt.Errorf("repository: mengganti identitas jejak %s: %w", b.sementara, err)
		}
	}
	return peta, nil
}

func (d *DB) gantiIdentitasTCO(ctx context.Context, tx *Tx, tabelLogis, sementara, tetap string) error {
	tabel, err := d.Qualify(tabelLogis)
	if err != nil {
		return err
	}
	if tabelLogis == TabelReinsurerTCO {
		return d.gantiIdentitasReinsurerTCO(ctx, tx, tabel, sementara, tetap)
	}
	q := sqlGantiIdentitasTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, tetap, sementara)
	if err != nil {
		return fmt.Errorf("repository: mengganti identitas %s %s: %w", tabelLogis, sementara, err)
	}
	return pastikanSatuBaris(hasil, "penggantian identitas "+tabelLogis)
}

// gantiIdentitasReinsurerTCO - salin, alihkan security, buang: FK security ->
// reinsurer tidak membolehkan PK induk diubah selagi anaknya menunjuknya.
func (d *DB) gantiIdentitasReinsurerTCO(ctx context.Context, tx *Tx, tabel, sementara, tetap string) error {
	security, err := d.Qualify(TabelSecurityTCO)
	if err != nil {
		return err
	}
	langkah := []struct {
		q    string
		args []any
		satu bool
	}{
		{sqlSalinReinsurerTCO(tabel), []any{tetap, sementara}, true},
		{sqlAlihSecurityTCO(security), []any{tetap, sementara}, false},
		{sqlBuangReinsurerSementaraTCO(tabel), []any{sementara}, true},
	}
	for _, l := range langkah {
		if err := PeriksaSQL(l.q); err != nil {
			return err
		}
		hasil, err := tx.tx.ExecContext(ctx, l.q, l.args...)
		if err != nil {
			return fmt.Errorf("repository: mengganti identitas reinsurer %s: %w", sementara, err)
		}
		if l.satu {
			if err := pastikanSatuBaris(hasil, "penggantian identitas reinsurer"); err != nil {
				return err
			}
		}
	}
	return nil
}
