package services

// Penegakan peran di lapisan layanan - tiket 07.
//
// Untuk apa berkas ini: SATU tempat yang menetapkan siapa boleh apa. Kontrol
// layar yang disembunyikan adalah kenyamanan; yang menolak adalah di sini,
// dan ia tetap menolak walau layarnya dilewati sama sekali.
//
// Dibaca sesudah: services.go.
//
// ⛔ Nol nama orang. Wewenang diukur dari PERAN akun (ADR-U-0030), dan nama
// tidak pernah menjadi dasar keputusan di lapisan mana pun.
//
// Istilah:
//   - peran : `pyPosition` di sistem lama - jabatan pemroses, bukan orangnya.

import (
	"errors"
	"fmt"

	"nusantarare/internal/models"
)

// Peran yang muncul di gerbang XML dan di Flow.
//
// `[terverifikasi]` `Flow/Register_Flow.xml`: `Assignment1` (Outstanding
// Claim) dan `Assignment2` (Input Register) dipegang `ReasLifeAdmin`;
// `Assignment3` (Medical Check) dan `Decision3` dipegang
// `ReasLifeMedicalAdvisor`; `Decision1` dipegang `ReasLifeSPV`.
//
// ⛔ Teks perannya hidup HANYA di sini. Tiket 03 dan 05 menamai IZIN-nya
// lebih dulu (`PeranSimpanOutstanding`, `PeranRejectOutstanding`), dan
// keduanya kini bernilai dari konstanta di bawah - bukan menaruh teks
// `"ReasLifeAdmin"` untuk kedua kalinya.
const (
	PeranAdmin          = "ReasLifeAdmin"
	PeranSPV            = "ReasLifeSPV"
	PeranMedicalAdvisor = "ReasLifeMedicalAdvisor"
)

// Izin - apa yang boleh dilakukan - dinamai terpisah dari PERAN, dan
// nilainya berasal dari peran di atas.
//
// ⛔ Ronde pertama tiket ini menulis "satu peran, satu nama" lalu tetap
// menaruh teks `"ReasLifeAdmin"` di tiga berkas. Yang benar: teksnya SATU
// tempat, sedangkan izin boleh bernama banyak - sebab "boleh menolak baris"
// dan "boleh mendaftarkan klaim" memang dua hal berbeda yang kebetulan
// dipegang peran yang sama hari ini.
const (
	// PeranInputRegister - `Assignment2` di Flow (tiket 02).
	PeranInputRegister = PeranAdmin
)

// ErrAksepBukanDariModulIni - Aksep ditulis Komite, bukan modul ini.
var ErrAksepBukanDariModulIni = errors.New(
	"services: status Aksep ditulis modul Komite, bukan Claim Life")

// WajibPeranPengubahStatus menolak pelaku yang tidak boleh menulis status KE.
//
// ⛔ Gerbangnya bergantung TUJUAN, bukan satu daftar peran yang datar. Ronde
// pertama memakai daftar `{SPV, Admin}` untuk tujuan apa pun - dan itu
// membuat SPV dapat MENOLAK baris lewat `Ubah`, melewati gerbang Admin di
// `Tolak` beserta syarat klaim-bernomornya. Satu daftar untuk dua keputusan
// yang berbeda adalah lubang, bukan kesederhanaan.
//
// `[terverifikasi]` sensus penulis status (tiket 04) dan gerbang XML:
//
//	Ditolak : `ReasLifeAdmin` - gerbang `AdjustmentDetail_Section` baris 15399
//	Aksep   : modul KOMITE (`KomitePostAdjustment`), bukan Claim Life
//
// ⚠️ `ReasLifeMedicalAdvisor` tidak menulis satu pun - tahap telaah medis
// MENELAAH, ia tidak memutuskan akseptasi (AC 9 spec).
func WajibPeranPengubahStatus(p Pelaku, ke models.StatusBaris) error {
	switch ke {
	case models.StatusDitolak:
		return WajibPeran(p, PeranRejectOutstanding)
	case models.StatusAksep:
		// ⛔ Ditolak bagi SIAPA PUN dari modul ini. Jalurnya lewat Komite, dan
		// pembacanya lahir di tiket 11. Melonggarkannya kepada Admin atau SPV
		// berarti mengarang wewenang yang korpus tidak berikan kepada mereka.
		return fmt.Errorf("%w: pakai jalur Komite (tiket 11)", ErrAksepBukanDariModulIni)
	default:
		return fmt.Errorf("%w: tujuan %v bukan keputusan akhir yang dapat ditulis di sini",
			ErrTransisiTidakSah, ke)
	}
}

// WajibWewenangKomite menolak pengiriman ke Komite yang tidak berwenang.
//
// `[terverifikasi]` `Section/AdjustmentDetail_Section.xml` baris 16091 berkas
// pecahan: `pyWorkPage.pyPosition =='ReasLifeSPV' || pyWorkPage.Type = 'TP'
// || pyWorkPage.Type = 'TR'`.
//
// ⚠️ Bacaan HARFIAHnya: untuk `TP`/`TR` gerbangnya terbuka tanpa memeriksa
// peran sama sekali. Tiket menuliskannya sebagai "`ReasLifeAdmin` dapat
// mengirim ke Komite"; keduanya sejalan selama Admin memang yang memegang
// tahapnya (`Assignment1`), dan selisihnya dicatat di tiket - bukan
// dipersempit diam-diam menjadi "hanya Admin", yang akan menolak SPV pada
// TP/TR padahal XML menerimanya.
//
// Identitas tetap wajib: gerbang yang terbuka bukan gerbang yang hilang.
func WajibWewenangKomite(p Pelaku, tipe string) error {
	if err := WajibIdentitas(p); err != nil {
		return err
	}
	if !TypeDikenal(tipe) {
		return fmt.Errorf("%w: %q; wewenang kirim-Komite tidak dapat ditentukan",
			ErrTypeTidakDikenal, tipe)
	}
	switch tipe {
	case TypeTP, TypeTR:
		// ⚠️ TERBUKA - dan tidak ada peran yang dibaca sama sekali. Itu bacaan
		// harfiah gerbangnya: `|| Type='TP' || Type='TR'` berdiri sendiri,
		// tanpa syarat peran. Akibatnya pemegang `ReasLifeMedicalAdvisor` pun
		// lolos di sini; yang menahannya adalah gerbang PENGUBAH STATUS di
		// hilir, bukan gerbang ini. Dicatat apa adanya, bukan dipersempit -
		// mempersempitnya akan menolak SPV pada TP/TR, yang XML terima.
		return nil
	case TypeQP, TypeQR:
		if p.PunyaPeran(PeranSimpanOutstanding) {
			return nil
		}
		return fmt.Errorf("%w: klaim ber-Type %s hanya dapat dikirim ke Komite oleh %s",
			ErrTanpaWewenang, tipe, PeranSimpanOutstanding)
	}
	return nil
}

// TypeDikenal menyatakan Type klaim termasuk keempat nilai yang ada cabangnya.
//
// ⛔ SATU tempat. Ronde pertama menulis `switch` keempat Type di dua berkas
// dengan urutan cabang yang berbeda - dan Type kelima kelak akan menuntut
// suntingan di keduanya, yang salah satunya pasti terlewat.
func TypeDikenal(tipe string) bool {
	switch tipe {
	case TypeQR, TypeQP, TypeTR, TypeTP:
		return true
	default:
		return false
	}
}
