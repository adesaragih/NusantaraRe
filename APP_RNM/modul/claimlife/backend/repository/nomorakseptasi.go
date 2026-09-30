package repository

// Sumber urut nomor akseptasi dan pemeriksa keunikannya - audit A0.
//
// Untuk apa berkas ini: DUA hal yang hanya Oracle dapat jawab - urut berikut
// dari sequence, dan apakah sebuah nomor sudah pernah dipakai.
//
// ⛔ PERAKITAN nomornya TIDAK di sini. Bentuk
// `'RNML-A'||{BusinessCode}||'.'||{MM}||'.'||{YY}||'.'||LPAD(urut,5,'0')`
// (`Claim Life/RDBList/Generate_NoAccept_Life.xml` baris 85) adalah aturan
// dagang, dan aturan dagang hidup di `services`. Arah ketergantungan
// `services -> repository` tidak boleh dibalik demi kenyamanan.
//
// ⛔ Yang tetap di Oracle hanya `NEXTVAL`: urut yang aman dari dua pemanggil
// serentak memang milik sequence, dan menirunya di aplikasi berarti menulis
// ulang kunci baris.
//
// Dibaca sesudah: pengenalwork.go.

import (
	"context"
	"fmt"

	"nusantarare/inti/backend/db"
)

// SequenceNomorAkseptasi adalah sequence penerbit urut nomor akseptasi.
//
// `[data DBA]` `POOLDATA.ACCEPTATIONNOLIFE_SEQ` - `INCREMENT BY 1`, `NOCACHE`,
// `NOCYCLE`.
//
// ⚠️ BUKAN sequence nomor KLAIM. Menyatukan keduanya membuat dua seri nomor
// saling memakan urut, dan nomor akseptasi tercetak di dokumen.
const SequenceNomorAkseptasi = "ACCEPTATIONNOLIFE_SEQ"

// UrutAkseptasiBerikut mengambil satu urut dari sequence.
func (r *PohonKlaim) UrutAkseptasiBerikut(ctx context.Context, tx *db.Tx) (string, error) {
	return r.db.NomorBerikut(ctx, tx, SequenceNomorAkseptasi)
}

// NomorAkseptasiDipakai meniru `GetAcceptedNoCL`.
//
// `[terverifikasi]` `SaveAdjustment_Act` langkah 1.6.1:
// `SELECT NO_ACCEPTATION FROM <tabel datar warisan> WHERE NO_ACCEPTATION = …`
// - nama tabelnya sendiri datang dari `namaTabelLama` lewat `Qualify`,
// dan tidak ditulis telanjang di mana pun (ADR-U-0033).
//
// ⚠️ Sequence menjamin URUT tidak berulang, tetapi TIDAK menjamin nomor
// RAKITAN unik: nomor lama dari migrasi tidak lahir dari sequence ini sama
// sekali, dan periode yang sama dapat bertemu urut yang sama bila sequence
// pernah di-reset.
func (r *PohonKlaim) NomorAkseptasiDipakai(ctx context.Context, tx *db.Tx,
	nomor string) (bool, error) {

	tabel, err := r.db.Qualify(namaTabelLama)
	if err != nil {
		return false, err
	}
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE NO_ACCEPTATION = :1`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := tx.QueryRowContext(ctx, q, nomor).Scan(&n); err != nil {
		return false, fmt.Errorf(
			"repository: memeriksa keunikan nomor akseptasi: %w", err)
	}
	return n > 0, nil
}
