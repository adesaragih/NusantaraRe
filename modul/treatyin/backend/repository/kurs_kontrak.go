package repository

// Kurs MILIK KONTRAK — migrasi 455 (8 Oktober 2026).
//
// Laporan pemakai (kontrak HEALTH QUOTA SHARE 2026): grid Rate of Exchange
// menampilkan TUJUH baris padahal yang diinput hanya DUA — *"biarkan apa yg
// di input user yang tampil … jangan di tambah tambahkan"*. Grid dulu
// membaca SEMUA baris `TREATYEXCHANGEYEARLY` tahun treaty-nya
// (`BacaKursTahunan`), dan tabel itu tidak tahu kontrak pemiliknya.
//
// ⭐ Kurs TETAP di `TREATYEXCHANGEYEARLY` (keputusan pemilik proses 4 & 6
// Oktober 2026). `T_TREATY_KURS` hanya mencatat baris mana milik kontrak
// mana. Kontrak TANPA catatan sama sekali (kontrak Pega lama, atau sebelum
// 455 terpasang) tetap membaca kurs tahun treaty-nya seperti dulu.

import (
	"context"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// TabelKursKontrak — penghubung kontrak ↔ baris `TREATYEXCHANGEYEARLY`.
const TabelKursKontrak = "T_TREATY_KURS"

// kursKontrakKosong — `IDKURS` penanda "grid kurs kontrak ini memang
// KOSONG" (semua baris di-Delete pemakai). Tanpa penanda, kontrak itu tidak
// punya catatan dan grid-nya jatuh lagi ke kurs SELURUH tahun — tampak
// bertambah sendiri (laporan 9 Oktober 2026). Tidak cocok dengan ID kurs
// mana pun, jadi tidak pernah ikut terbaca oleh join.
const kursKontrakKosong = "-"

// BacaKursKontrak — grid Rate of Exchange SATU kontrak: baris kurs yang
// terhubung ke kontrak ini, urut `URUTAN`. Hanya kontrak yang BELUM PERNAH
// punya catatan yang membaca kurs tahun treaty (`BacaKursTahunan`).
func (g *Gudang) BacaKursKontrak(ctx context.Context, masterID, tahunTreaty string) ([]models.BarisKursWarisan, error) {
	kol, err := g.kolomTerpasang(ctx, TabelKursKontrak)
	if err != nil {
		return nil, err
	}
	if kol == nil || strings.TrimSpace(masterID) == "" {
		return g.BacaKursTahunan(ctx, tahunTreaty)
	}
	hub, err := g.db.Qualify(TabelKursKontrak)
	if err != nil {
		return nil, err
	}
	kurs, err := g.db.Qualify(TabelKursTahunan)
	if err != nil {
		return nil, err
	}
	cacah := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE MASTERID = :1`, hub)
	if err := db.PeriksaSQL(cacah); err != nil {
		return nil, err
	}
	var n int
	if err := g.db.QueryRowContext(ctx, cacah, masterID).Scan(&n); err != nil {
		return nil, fmt.Errorf("repository: menghitung kurs kontrak %s: %w", masterID, err)
	}
	if n == 0 {
		return g.BacaKursTahunan(ctx, tahunTreaty)
	}
	// ⚠️ `TREATYEXCHANGEYEARLY.ID` TIDAK unik antartahun (10124 dipakai
	// baris 2026 IDR DAN baris 2025 JPY) — kuncinya (ID, TREATYYEAR), sama
	// dengan `UPDATE` di `tulisKurs`. Tanpa saringan tahun, grid kontrak 2026
	// ikut menampilkan kurs 2025 yang ID-nya kebetulan sama.
	q := fmt.Sprintf(`SELECT k.ID, k.CURRENCY, k.IDCURRENCY, k.TOIDR, k.STARTDATE, k.ENDDATE
		FROM %s h JOIN %s k ON k.ID = h.IDKURS
		WHERE h.MASTERID = :1 AND k.TREATYYEAR = :2 ORDER BY h.URUTAN`, hub, kurs)
	return g.bacaBarisKurs(ctx, q, masterID, strings.TrimSpace(tahunTreaty))
}

// kursMilikLain — ID kurs yang tercatat milik kontrak SELAIN `masterID`.
// `tulisKurs` tidak pernah meng-`UPDATE` baris ini. Sebelum 455 terpasang:
// kosong.
func (g *Gudang) kursMilikLain(ctx context.Context, tx *db.Tx, masterID string) (map[string]bool, error) {
	kol, err := g.kolomTerpasang(ctx, TabelKursKontrak)
	if err != nil || kol == nil {
		return map[string]bool{}, err
	}
	hub, err := g.db.Qualify(TabelKursKontrak)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT DISTINCT IDKURS FROM %s WHERE MASTERID <> :1`, hub)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := tx.QueryContext(ctx, q, masterID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca kurs milik kontrak lain: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := map[string]bool{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("repository: membaca kurs milik kontrak lain: %w", err)
		}
		out[strings.TrimSpace(id)] = true
	}
	return out, rows.Err()
}

// idKursBaru — situs || LPAD(TREATYEXCHANGE_SEQ.NEXTVAL, 4, '0'), dilewati
// selama ID itu SUDAH dipakai baris tahun MANA PUN. Sequence-nya tertinggal
// dari data Pega lama: 10124, 10125, 10130 … yang dibuatnya sudah dipakai
// kurs 2025 (JPY, KRW, PGK), dan ID kembar antartahun itulah yang membuat
// grid 2026 memuat kurs 2025 lalu Save menimpanya (9 Oktober 2026).
func (g *Gudang) idKursBaru(ctx context.Context, tx *db.Tx, situs string) (string, error) {
	nama, err := g.db.Qualify(TabelKursTahunan)
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE ID = :1`, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	for coba := 0; coba < 10000; coba++ {
		n, err := g.db.NomorBerikut(ctx, tx, seqKursTahunan)
		if err != nil {
			return "", err
		}
		id := situs + kiriNol(n, 4)
		var ada int
		if err := tx.QueryRowContext(ctx, q, id).Scan(&ada); err != nil {
			return "", fmt.Errorf("repository: memeriksa ID kurs %s: %w", id, err)
		}
		if ada == 0 {
			return id, nil
		}
	}
	return "", fmt.Errorf("repository: tidak menemukan ID kurs yang belum dipakai")
}

// tulisHubunganKurs — catatan baris kurs milik kontrak `masterID`, ditulis
// ULANG sesuai SELURUH grid yang disimpan (urutan grid; baris yang tidak
// diubah ikut sebagai `Tetap`). Baris yang di-Delete di grid hanya LEPAS
// dari kontrak ini; baris `TREATYEXCHANGEYEARLY`-nya tetap (tabel bersama —
// nol `DELETE` di sana). Grid kosong dicatat dengan `kursKontrakKosong`.
//
// Sebelum 455 terpasang: nol tulis, nol galat.
func (g *Gudang) tulisHubunganKurs(ctx context.Context, tx *db.Tx, masterID string, idKurs []string) error {
	kol, err := g.kolomTerpasang(ctx, TabelKursKontrak)
	if err != nil || kol == nil {
		return err
	}
	hub, err := g.db.Qualify(TabelKursKontrak)
	if err != nil {
		return err
	}
	hapus := fmt.Sprintf(`DELETE FROM %s WHERE MASTERID = :1`, hub)
	if err := db.PeriksaSQL(hapus); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, hapus, masterID); err != nil {
		return fmt.Errorf("repository: melepas kurs kontrak %s: %w", masterID, err)
	}
	sisip := fmt.Sprintf(`INSERT INTO %s (MASTERID, IDKURS, URUTAN) VALUES (:1, :2, :3)`, hub)
	if err := db.PeriksaSQL(sisip); err != nil {
		return err
	}
	sudah := map[string]bool{}
	for i, id := range idKurs {
		if id == "" || sudah[id] {
			continue
		}
		sudah[id] = true
		if _, err := tx.ExecContext(ctx, sisip, masterID, id, i+1); err != nil {
			return fmt.Errorf("repository: menghubungkan kurs %s ke kontrak %s: %w", id, masterID, err)
		}
	}
	if len(sudah) == 0 {
		if _, err := tx.ExecContext(ctx, sisip, masterID, kursKontrakKosong, 0); err != nil {
			return fmt.Errorf("repository: mencatat kurs kosong kontrak %s: %w", masterID, err)
		}
	}
	return nil
}
