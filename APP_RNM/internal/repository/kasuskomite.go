package repository

// Penulis kasus Komite - butir af, A2.
//
// Untuk apa berkas ini: melahirkan satu kasus komite beserta tangganya, dalam
// SATU transaksi bersama penautan baris adjustment.
//
// Meniru `Claim Life/Activity/CreateKMTLife_Act.xml`:
//
//	 866 · 972  `KomiteList(<APPEND>).KomiteID = .OPERATOR_ID`
//	 912 · 993  `KomiteList(<LAST>).KomiteAproval = 0`
//	 932 · 1014 `KomiteList(<LAST>).KomiteEmail = .EMAIL`
//	 952 · 1041 `KomiteList(<LAST>).IDKomite = .JABATAN`
//	1322        `childPageKomite.KomiteLoop = SizeOfPropertyList(KomiteList)`
//	1398        `childPageKomite.KomiteCount = 1`
//	1461        `childPageKomite.CLMNO = pyWorkPage.pyID`
//
// ⛔ `pxAddChildWork` (1543) melahirkan work object anaknya; di sini itu satu
// baris `T_WORK_CLAIM` berawalan `KMTLF-` dengan `COVER_KEY` menunjuk klaim
// induknya.
//
// Dibaca sesudah: roster.go.

import (
	"context"
	"fmt"
	"time"
)

// AnggotaTangga adalah satu anggota komite yang akan ditulis ke tangganya.
//
// ⚠️ `OperatorID` dan `Email` DATA ORANG - dibaca saat jalan dari roster,
// ditulis ke tangganya, dan tidak ke mana pun lagi.
type AnggotaTangga struct {
	Urut       int
	OperatorID string
	Jabatan    string
	Email      string
}

// approvalAwal adalah nilai `KOMITE_APPROVAL` saat tangga dibentuk.
//
// `[terverifikasi]` `CreateKMTLife_Act.xml` 912 dan 993: `= 0`. TEKS, bukan
// bilangan (ADR-U-0022).
const approvalAwal = "0"

// BuatKasusKomite melahirkan kasus komite beserta tangganya.
//
// Mengembalikan pengenal kasus (`KMTLF-xxxxxx`), yang sekaligus menjadi
// `T_GENERAL_KOMITE.ID` - shared primary key.
func (r *PohonKlaim) BuatKasusKomite(ctx context.Context, tx *Tx,
	klaimID, adjustmentID, lini, tipe string, anggota []AnggotaTangga,
	saat time.Time) (string, error) {

	if len(anggota) == 0 {
		return "", fmt.Errorf("repository: tangga komite kosong; " +
			"kasus tanpa anggota tidak dapat diputuskan siapa pun")
	}
	kasusID, err := r.PengenalWorkBerikut(ctx, tx, AwalanKomite)
	if err != nil {
		return "", err
	}

	// 1. Work object anaknya - `pxAddChildWork`.
	work, err := r.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return "", err
	}
	qWork := fmt.Sprintf(`INSERT INTO %s
		(ID, COVER_KEY, LINI, TYPE, TGL_UPDATE)
		VALUES (:1,:2,:3,:4,:5)`, work)
	if err := PeriksaSQL(qWork); err != nil {
		return "", err
	}
	hasil, err := tx.tx.ExecContext(ctx, qWork, kasusID,
		kosongJadiNil(klaimID), kosongJadiNil(lini), kosongJadiNil(tipe), saat)
	if err != nil {
		return "", fmt.Errorf("repository: melahirkan work object komite: %w", err)
	}
	if err := pastikanSatuBaris(hasil, "kelahiran work object komite"); err != nil {
		return "", err
	}

	// 2. Kasus komitenya - shared PK, dan tinggi tangga dari cacah anggotanya.
	gen, err := r.db.Qualify("T_GENERAL_KOMITE")
	if err != nil {
		return "", err
	}
	qGen := fmt.Sprintf(`INSERT INTO %s
		(ID, ADJUSTMENT_ID, KOMITE_LOOP, KOMITE_COUNT)
		VALUES (:1,:2,:3,:4)`, gen)
	if err := PeriksaSQL(qGen); err != nil {
		return "", err
	}
	// ⛔ `KOMITE_COUNT = 1` seperti XML (1398), dan `KOMITE_LOOP` dari cacah
	// anggota (1322). `ACCEPT_STATUS` sengaja dibiarkan KOSONG: belum ada
	// keputusan, dan kosong berbeda dari nol (ADR-U-0027).
	hasil, err = tx.tx.ExecContext(ctx, qGen, kasusID, adjustmentID, len(anggota), 1)
	if err != nil {
		return "", fmt.Errorf("repository: melahirkan kasus komite: %w", err)
	}
	if err := pastikanSatuBaris(hasil, "kelahiran kasus komite"); err != nil {
		return "", err
	}

	// 3. Tangganya, satu baris per anggota.
	list, err := r.db.Qualify("T_KOMITE_KOMITELIST")
	if err != nil {
		return "", err
	}
	qList := fmt.Sprintf(`INSERT INTO %s
		(ID, DATA_KOMITE_ID, KOMITE_URUT, KOMITE_OPERATORID, KOMITE_JABATAN,
		 KOMITE_EMAIL, KOMITE_APPROVAL)
		VALUES (:1,:2,:3,:4,:5,:6,:7)`, list)
	if err := PeriksaSQL(qList); err != nil {
		return "", err
	}
	for _, a := range anggota {
		id, err := r.nomorBerikut(ctx, tx, "SEQ_KOMITE_KOMITELIST")
		if err != nil {
			return "", err
		}
		hasil, err := tx.tx.ExecContext(ctx, qList, id, kasusID, a.Urut,
			kosongJadiNil(a.OperatorID), kosongJadiNil(a.Jabatan),
			kosongJadiNil(a.Email), approvalAwal)
		if err != nil {
			return "", fmt.Errorf("repository: menulis anggota tangga komite: %w", err)
		}
		if err := pastikanSatuBaris(hasil, "penulisan anggota tangga komite"); err != nil {
			return "", err
		}
	}
	return kasusID, nil
}
