package repository

// Pembuatan kasus endorsement dan penyalinan versi polis lama -
// `MappingEDMLife` (14 langkah, nol `//`).
//
// ⛔ PENYALINAN HIMPUNAN (`INSERT … SELECT`): peserta satu polis dapat
// berjumlah puluhan ribu; tidak satu pun dibaca ke Go. Identitas baru
// DETERMINISTIK - `STANDARD_HASH(<kasus> || '/<jenis>/' || <ID lama>, 'MD5')`,
// 32 heksa = `VARCHAR2(32)` - sehingga spreading dan spreading retro dapat
// menunjuk peserta/spreading versi baru tanpa tabel pemetaan.
//
// ⛔ Baris `Delete` versi lama TIDAK disalin (`MappingEDMLife` 11.2 b2737,
// RALAT R10); sisanya `EDM_STATUS = 'Old'` (11.1 b2609).
//
// Dibaca sesudah: edm_baca.go, edm_kolom.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
)

// ErrKasusTerbukaGanda - index unik `UX_PL_EDM_TERBUKA` (migrasi 482) menolak
// kasus terbuka kedua atas polis yang sama (pembuatan bersamaan).
var ErrKasusTerbukaGanda = errors.New("repository: polis sudah punya kasus endorsement terbuka")

// PengenalKasusBaru menerbitkan `EDMLF-<n>` dari `SEQ_WORK_EDM_LIFE`.
func (g *Gudang) PengenalKasusBaru(ctx context.Context, tx *db.Tx) (string, error) {
	urut, err := g.db.NomorBerikut(ctx, tx, urutanKasus)
	if err != nil {
		return "", err
	}
	return models.RakitPengenalKasus(urut)
}

// sqlKasusTerbuka - gerbang 3: kasus EDM terbuka atas polis itu
// (`IDX_PL_OLD_POLICY_NO`). `Resolved-Rejected` tidak memblokir (OQ-EDM-002).
func sqlKasusTerbuka(polis string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s p
	  WHERE p.OLD_POLICY_NO = :1 AND p.ID LIKE :2 AND p.STATUSS IS NULL AND ROWNUM = 1`, polis)
}

// AdaKasusTerbuka - gerbang 3 (`FilterProteksiEDMLife`).
func (g *Gudang) AdaKasusTerbuka(ctx context.Context, tx *db.Tx, nomorPolis string) (bool, error) {
	n, err := g.nama(tabelPolis)
	if err != nil {
		return false, err
	}
	q := sqlKasusTerbuka(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var c int
	if err := g.pakai(tx).QueryRowContext(ctx, q, nomorPolis, polaPengenalKasus).Scan(&c); err != nil {
		return false, fmt.Errorf("repository: membaca kasus terbuka polis %q: %w", nomorPolis, err)
	}
	return c > 0, nil
}

// --- kepala versi lama ------------------------------------------------------

// polaJalurJSON - jalur `JSON_VALUE` yang sah dirakit (konstanta korpus).
var polaJalurJSON = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_.]*$`)

// sqlKepalaAplikasi membaca 24 kolom kepala versi sistem baru.
func sqlKepalaAplikasi(polis string) string {
	var b strings.Builder
	b.WriteString("SELECT ")
	for i, k := range models.KolomKepalaSalin {
		if i > 0 {
			b.WriteString(", ")
		}
		if k.Tanggal {
			b.WriteString("TO_CHAR(p." + k.Kolom + ", 'YYYY-MM-DD')")
		} else {
			b.WriteString("p." + k.Kolom)
		}
	}
	fmt.Fprintf(&b, " FROM %s p WHERE p.ID = :1", polis)
	return b.String()
}

// sqlKepalaWarisan membaca 24 properti kepala dari `JSON_POLIS.DATA_JSON` -
// padanan `Obj-Open-By-Handle` work NB (`MappingEDMLife` 8 b1547) lalu
// penetapan langkah 9. `NULL ON ERROR`: properti yang tidak ada = kosong.
func sqlKepalaWarisan(jsonPolis string) (string, error) {
	var b strings.Builder
	b.WriteString("SELECT ")
	for i, k := range models.KolomKepalaSalin {
		if !polaJalurJSON.MatchString(k.Properti) {
			return "", fmt.Errorf("repository: jalur JSON %q tidak sah", k.Properti)
		}
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("JSON_VALUE(j.DATA_JSON, '$." + k.Properti + "' NULL ON ERROR)")
	}
	fmt.Fprintf(&b, " FROM %s j WHERE j.IDPEGA = :1 AND j.NOPOLIS = :2 FETCH FIRST 1 ROWS ONLY", jsonPolis)
	return b.String(), nil
}

// polaTanggalPega - tanggal Pega `YYYYMMDD` atau `YYYYMMDDTHHMMSS.mmm GMT`.
var polaTanggalPega = regexp.MustCompile(`^(\d{4})(\d{2})(\d{2})(T|$)`)
var polaTanggalISO = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})`)

// TanggalPega mengubah nilai tanggal `DATA_JSON` menjadi `YYYY-MM-DD`; bentuk
// lain = kosong (dan `ok` false - pemanggil mencatatnya, tidak menebak).
func TanggalPega(v string) (string, bool) {
	t := strings.TrimSpace(v)
	if t == "" {
		return "", true
	}
	if m := polaTanggalPega.FindStringSubmatch(t); m != nil {
		return m[1] + "-" + m[2] + "-" + m[3], true
	}
	if m := polaTanggalISO.FindStringSubmatch(t); m != nil {
		return m[1] + "-" + m[2] + "-" + m[3], true
	}
	return "", false
}

// KepalaSumber membaca kepala versi lama, berkunci nama kolom. Tanggal warisan
// yang tidak terbaca dikembalikan di `rusak` (nama kolom), nilainya kosong.
func (g *Gudang) KepalaSumber(ctx context.Context, tx *db.Tx, v models.Versi, nomorPolis string) (map[string]string, []string, error) {
	var q string
	var args []any
	switch v.Jenis {
	case models.SumberAplikasi:
		n, err := g.nama(tabelPolis)
		if err != nil {
			return nil, nil, err
		}
		q, args = sqlKepalaAplikasi(n[0]), []any{v.ID}
	case models.SumberWarisan:
		n, err := g.nama(tabelPolisWarisan)
		if err != nil {
			return nil, nil, err
		}
		if q, err = sqlKepalaWarisan(n[0]); err != nil {
			return nil, nil, err
		}
		args = []any{v.ID, nomorPolis}
	default:
		return nil, nil, fmt.Errorf("repository: sumber versi %q tidak dikenal", v.Jenis)
	}
	if err := db.PeriksaSQL(q); err != nil {
		return nil, nil, err
	}
	mentah := make([]sql.NullString, len(models.KolomKepalaSalin))
	ptr := make([]any, len(mentah))
	for i := range mentah {
		ptr[i] = &mentah[i]
	}
	err := g.pakai(tx).QueryRowContext(ctx, q, args...).Scan(ptr...)
	if err == sql.ErrNoRows {
		return nil, nil, ErrTidakAda
	}
	if err != nil {
		return nil, nil, fmt.Errorf("repository: membaca kepala versi %s %q: %w", v.Jenis, v.ID, err)
	}
	kepala := make(map[string]string, len(mentah))
	var rusak []string
	for i, k := range models.KolomKepalaSalin {
		s := strings.TrimSpace(mentah[i].String)
		if k.Tanggal && v.Jenis == models.SumberWarisan {
			iso, ok := TanggalPega(s)
			if !ok {
				rusak = append(rusak, k.Kolom)
			}
			s = iso
		}
		kepala[k.Kolom] = s
	}
	return kepala, rusak, nil
}

// --- kepala kasus -------------------------------------------------------------

// KasusTulis adalah isi baris `T_PREMIUM_LIST` kasus baru.
type KasusTulis struct {
	ID         string
	NomorPolis string
	EdmType    string
	EdmDate    string
	EdmNote    string
	// ProdKe - versi yang AKAN dicapai kasus (versi berjalan + 1); ditetapkan
	// ulang saat diresmikan.
	ProdKe  int
	Pembuat string
	Kepala  map[string]string
}

// sqlSisipKasus - kepala kasus. `NO_POLIS` SENGAJA kosong sampai diresmikan:
// pembaca versi berjalan (Claim Life `polis_ringkas.go`, `sqlVersiEDM`)
// mencari lewat kolom itu, dan kasus terbuka atau ditolak bukan versi polis.
func sqlSisipKasus(polis string) string {
	kolom := []string{"ID", "ID_PEGA", "TGL_INPUT", "CREATE_OP_NAME", "OLD_POLICY_NO", "EDM_TYPE", "EDM_DATE", "EDM_NOTE", "PROD_KE"}
	nilai := []string{":1", ":2", "SYSDATE", ":3", ":4", ":5", "TO_DATE(:6, 'YYYY-MM-DD')", ":7", ":8"}
	for i, k := range models.KolomKepalaSalin {
		kolom = append(kolom, k.Kolom)
		p := fmt.Sprintf(":%d", 9+i)
		if k.Tanggal {
			p = "TO_DATE(" + p + ", 'YYYY-MM-DD')"
		}
		nilai = append(nilai, p)
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, polis, strings.Join(kolom, ", "), strings.Join(nilai, ", "))
}

// SisipKasus menulis kepala kasus baru.
func (g *Gudang) SisipKasus(ctx context.Context, tx *db.Tx, k KasusTulis) error {
	n, err := g.nama(tabelPolis)
	if err != nil {
		return err
	}
	q := sqlSisipKasus(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	args := []any{k.ID, k.ID, k.Pembuat, k.NomorPolis, k.EdmType, db.KosongJadiNil(k.EdmDate), db.KosongJadiNil(k.EdmNote), k.ProdKe}
	for _, kol := range models.KolomKepalaSalin {
		args = append(args, db.KosongJadiNil(k.Kepala[kol.Kolom]))
	}
	hasil, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		if strings.Contains(err.Error(), "ORA-00001") && strings.Contains(strings.ToUpper(err.Error()), "UX_PL_EDM_TERBUKA") {
			return fmt.Errorf("%w: %s", ErrKasusTerbukaGanda, k.NomorPolis)
		}
		return fmt.Errorf("repository: menyisipkan kasus endorsement %q: %w", k.ID, err)
	}
	return db.PastikanSatuBaris(hasil, "kasus endorsement")
}

// --- penyalinan peserta -------------------------------------------------------

// idBaru - ekspresi identitas deterministik `<kasus>/<jenis>/<ID lama>`.
func idBaru(penampung, jenis, kolomLama string) string {
	return "RAWTOHEX(STANDARD_HASH(" + penampung + " || '/" + jenis + "/' || " + kolomLama + ", 'MD5'))"
}

// sqlSalinPesertaAplikasi - `MappingEDMLife` 10-11 atas versi sistem baru.
// `PARENT_ID` = peserta versi lama (penyimpangan sadar 10). Penampung:
// :1 kasus (hash), :2 kasus (induk), :3 kasus (ID_PEGA), :4 'Old',
// :5 versi lama, :6 'Delete'. Kunci baca `IDX_PLD_PL` (PREMIUM_LIST_ID).
func sqlSalinPesertaAplikasi(peserta string) string {
	sumber := make([]string, len(kolomNilaiPeserta))
	for i, k := range kolomNilaiPeserta {
		sumber[i] = "d." + k
	}
	return fmt.Sprintf(`INSERT INTO %s (ID, PREMIUM_LIST_ID, PARENT_ID, ID_PEGA, PL_NUMBER, EDM_STATUS, %s)
	SELECT %s, :2, d.ID, :3, d.PL_NUMBER, :4, %s
	  FROM %s d
	 WHERE d.PREMIUM_LIST_ID = :5 AND (d.EDM_STATUS IS NULL OR d.EDM_STATUS <> :6)`,
		peserta, strings.Join(kolomNilaiPeserta, ", "), idBaru(":1", "D", "d.ID"), strings.Join(sumber, ", "), peserta)
}

// sqlSalinSpreading - spreading peserta yang ikut tersalin (spec AC 59), di
// bawah peserta versi baru. Kunci baca `IDX_PLD_PL` lalu `IDX_PLS_DETAIL`.
func sqlSalinSpreading(spreading, peserta string) string {
	sumber := make([]string, len(kolomNilaiSpreading))
	for i, k := range kolomNilaiSpreading {
		sumber[i] = "s." + k
	}
	return fmt.Sprintf(`INSERT INTO %s (ID, DETAIL_ID, %s)
	SELECT %s, %s, %s
	  FROM %s s JOIN %s d ON d.ID = s.DETAIL_ID
	 WHERE d.PREMIUM_LIST_ID = :3 AND (d.EDM_STATUS IS NULL OR d.EDM_STATUS <> :4)`,
		spreading, strings.Join(kolomNilaiSpreading, ", "), idBaru(":1", "S", "s.ID"), idBaru(":2", "D", "s.DETAIL_ID"),
		strings.Join(sumber, ", "), spreading, peserta)
}

// sqlSalinSpreadingRetro - spreading retro di bawah spreading versi baru.
func sqlSalinSpreadingRetro(retro, spreading, peserta string) string {
	sumber := make([]string, len(kolomNilaiSpreadingRetro))
	for i, k := range kolomNilaiSpreadingRetro {
		sumber[i] = "r." + k
	}
	return fmt.Sprintf(`INSERT INTO %s (ID, SPREADING_ID, %s)
	SELECT %s, %s, %s
	  FROM %s r JOIN %s s ON s.ID = r.SPREADING_ID JOIN %s d ON d.ID = s.DETAIL_ID
	 WHERE d.PREMIUM_LIST_ID = :3 AND (d.EDM_STATUS IS NULL OR d.EDM_STATUS <> :4)`,
		retro, strings.Join(kolomNilaiSpreadingRetro, ", "), idBaru(":1", "R", "r.ID"), idBaru(":2", "S", "r.SPREADING_ID"),
		strings.Join(sumber, ", "), retro, spreading, peserta)
}

// sqlSalinPesertaWarisan - versi SISTEM LAMA (new business ATAU endorsement): peserta `PL_NUMBER` +
// `IDPEGA` versi itu, TANPA baris `Delete` - `MappingEDMLife` 10 b2339 salin halaman (PRE=false, selalu),
// 11.1 b2583 baris selain `Delete` → `"Old"` (b2609; prakondisi b2691 `.EDMStatus=="Delete"` T=3 = lewati),
// 11.2 b2737 buang `.EDMStatus=="Delete"` (b2831). Peserta new business (status NULL) selalu lolos.
// ⚠️ Dibandingkan sesudah `TRIM`, seperti penyaring hidup Claim Life (`TRIM(EDMSTATUS)`): Pega membandingkan
// nilai halaman apa adanya, tetapi baris warisan ber-`'Delete '` berspasi ada - tanpa `TRIM` peserta yang
// sudah dihapus hidup kembali sebagai `Old`, sementara Claim Life menganggapnya mati.
// Penyaring status EDM warisan ini diizinkan sejak penjaga Claim Life dipersempit ke `modul/claimlife/`
// (K5 keputusan work owner 01-10-2026, OQ-EDM-016, `6047ca8`) - sebelumnya versi endorsement sistem lama
// ditolak.
//
// ⛔ E4: kunci `PL_NUMBER = :5` → index `M_LIFE_PREMIUM_DETAIL_INDEX4`
// (`PL_NUMBER`), `IDPEGA` dan status penyaring sesudahnya; nol pemindaian penuh.
// `PARENT_ID` kosong (OQ-EDM-007).
func sqlSalinPesertaWarisan(peserta, warisan string) string {
	sumber := make([]string, len(kolomNilaiPeserta))
	for i, k := range kolomNilaiPeserta {
		sumber[i] = ekspresiWarisan(k)
	}
	return fmt.Sprintf(`INSERT INTO %s (ID, PREMIUM_LIST_ID, PARENT_ID, ID_PEGA, PL_NUMBER, EDM_STATUS, %s)
	SELECT %s, :2, NULL, :3, m.PL_NUMBER, :4, %s
	  FROM %s m
	 WHERE m.PL_NUMBER = :5 AND m.IDPEGA = :6 AND NVL(TRIM(m.EDMSTATUS), '-') <> :7`,
		peserta, strings.Join(kolomNilaiPeserta, ", "), idBaru(":1", "W", "m.ID"), strings.Join(sumber, ", "), warisan)
}

// Salinan adalah cacah baris yang tersalin ke versi kasus.
type Salinan struct {
	Peserta, Spreading, SpreadingRetro int
}

// SalinVersi menyalin peserta (dan untuk sumber aplikasi: spreading +
// spreading retro) versi lama ke kasus, di transaksi pemanggil.
func (g *Gudang) SalinVersi(ctx context.Context, tx *db.Tx, kasusID string, v models.Versi, nomorPolis string) (Salinan, error) {
	n, err := g.nama(tabelPeserta, tabelSpreading, tabelSpreadingRetro, tabelPesertaWarisanEDM)
	if err != nil {
		return Salinan{}, err
	}
	jalankan := func(nama, q string, args ...any) (int, error) {
		if err := db.PeriksaSQL(q); err != nil {
			return 0, err
		}
		h, err := tx.ExecContext(ctx, q, args...)
		if err != nil {
			return 0, fmt.Errorf("repository: menyalin %s kasus %q: %w", nama, kasusID, err)
		}
		c, err := h.RowsAffected()
		return int(c), err
	}
	var s Salinan
	switch v.Jenis {
	case models.SumberAplikasi:
		if s.Peserta, err = jalankan("peserta", sqlSalinPesertaAplikasi(n[0]),
			kasusID, kasusID, kasusID, models.StatusOld, v.ID, models.StatusDelete); err != nil {
			return Salinan{}, err
		}
		if s.Spreading, err = jalankan("spreading", sqlSalinSpreading(n[1], n[0]),
			kasusID, kasusID, v.ID, models.StatusDelete); err != nil {
			return Salinan{}, err
		}
		if s.SpreadingRetro, err = jalankan("spreading retro", sqlSalinSpreadingRetro(n[2], n[1], n[0]),
			kasusID, kasusID, v.ID, models.StatusDelete); err != nil {
			return Salinan{}, err
		}
	case models.SumberWarisan:
		if s.Peserta, err = jalankan("peserta warisan", sqlSalinPesertaWarisan(n[0], n[3]),
			kasusID, kasusID, kasusID, models.StatusOld, nomorPolis, v.ID, models.StatusDelete); err != nil {
			return Salinan{}, err
		}
	default:
		return Salinan{}, fmt.Errorf("repository: sumber versi %q tidak dikenal", v.Jenis)
	}
	return s, nil
}
