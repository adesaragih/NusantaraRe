// Package models memuat bentuk data modul Bordereaux - kelas Pega `ASM-FW-GISFW-Int-BORDEREAUX` dan 29 kelas
// detailnya (folder korpus `D:\XML\RNM_BRD\Bordereaux`, keputusan work owner 04-10-2026).
package models

import "fmt"

// Type berkas Bordereaux - radio `TYPE` section `InputBordereaux` (nilai tersimpan di BORDEREAUX.TYPE).
const (
	TypePremium     = "PREMIUM"
	TypeClaim       = "CLAIM"
	TypeSubrogation = "SUBROGATION"
)

// DaftarType - urutan radio Type.
var DaftarType = []string{TypePremium, TypeClaim, TypeSubrogation}

// BusinessBonding - satu-satunya Business untuk SUBROGATION (`SetSubrogation_ACT` langkah 1).
const BusinessBonding = "BONDING"

// Jenis kolom tabel detail - menentukan cara nilai CSV dibaca (perbaikan bug Pega: menurut TIPE KOLOM TABEL, bukan
// tanda tiap activity Mapping*).
type Jenis int

const (
	// Teks - VARCHAR2: disimpan apa adanya (dipangkas spasi tepinya).
	Teks Jenis = iota
	// Angka - NUMBER: format Indonesia `1.234.567,89`; `%` dibuang; `-` dan kosong = NULL.
	Angka
	// Tanggal - DATE: `dd/MM/yyyy`.
	Tanggal
)

// KolomCSV - satu kolom berkas CSV dan kolom tabel tujuannya.
type KolomCSV struct {
	// Posisi kolom di CSV, mulai 1 (sama dengan urutan langkah activity Mapping* Pega).
	Posisi int
	// Kolom tabel detail tujuan.
	Kolom string
	// Judul - judul kolom di header templat (spasi diringkas); dipakai pesan validasi dan kepala grid.
	Judul string
	Jenis Jenis
	// Panjang - batas karakter kolom teks.
	Panjang int
	// Presisi dan Skala kolom angka; 0 = NUMBER tanpa batas.
	Presisi, Skala int
}

// KombinasiBdx - satu pasangan Type x Business yang punya Upload CSV, tabel detail, dan templat.
type KombinasiBdx struct {
	// Kode - akhiran kode slot Template Manager (`bordereaux.<Kode>`).
	Kode     string
	Type     string
	Business string
	// Tabel detail tujuan, mis. BORDEREAUX_PREMI_FIRE.
	Tabel string
	// KolomInduk - kolom penunjuk BORDEREAUX.BDX_ID di tabel detail: `BDX_ID` atau `ID_BDX` (tabelnya tidak seragam;
	// rule hapus Pega selalu memakai BDX_ID sehingga gagal di tabel ber-ID_BDX - diperbaiki di sini).
	KolomInduk string
	// BerkasTemplat - nama berkas templat bawaan di `backend/templat/`.
	BerkasTemplat string
	// DaftarJSON - nama page list kombinasi ini di JSON Pega lama (`M_BORDEREAUX.DATA_JSON`), dibaca Copy Old Data.
	DaftarJSON string
	Kolom      []KolomCSV
}

// CariKombinasi - kombinasi untuk Type dan Business.
func CariKombinasi(tipe, bisnis string) (KombinasiBdx, error) {
	for _, k := range Kombinasi {
		if k.Type == tipe && k.Business == bisnis {
			return k, nil
		}
	}
	return KombinasiBdx{}, fmt.Errorf("Type %s and Business %s have no bordereaux detail", tipe, bisnis)
}

// DaftarBusiness - Business untuk satu Type, urutan radio Pega.
func DaftarBusiness(tipe string) []string {
	var hasil []string
	for _, k := range Kombinasi {
		if k.Type == tipe {
			hasil = append(hasil, k.Business)
		}
	}
	return hasil
}
