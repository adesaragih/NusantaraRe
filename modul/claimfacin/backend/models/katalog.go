package models

// Untuk apa berkas ini: KATALOG KOLOM - satu-satunya pemetaan properti Pega -> kolom Oracle -> golongan tipe untuk
// penyimpanan relasional Claim Fac In (keputusan work owner 09-10-2026 OQ-CFI-01 "tabel Prop + kolom item").
//
// Susunan (docs/STRUKTUR-TABEL-CLAIM-FACIN.md), POHON tiga tingkat objek -> item objek -> estimasi / adjustment:
//
//	T_GENERAL_CLAIM (kepala, tabel bersama Claim Life; kolom khas FACIN ALTER ADD 560)
//	└ T_CLAIM_OBJECT          (baru)  ClaimData.ObjectList                      ID stabil, CLAIM_ID CASCADE
//	  └ T_CLAIM_OBJECT_ITEM   (baru)  .ObjectItemList                           ID stabil, OBJECT_ID CASCADE
//	    ├ T_CLAIM_ESTIMATION          .EstimationList                           CLAIM_ID + OBJECT_ITEM_ID (ADD)
//	    ├ T_CLAIM_SPREADING           .SpreadingList (JENIS POLIS) / .SpreadingClaim (JENIS KLAIM)
//	    ├ T_CLAIM_BREAK_QS            .SpreadingAdjustment (tingkat item)
//	    └ T_CLAIM_ADJUSTMENT          .Adjustment                               ID stabil (komite menunjuknya)
//	      ├ T_CLAIM_ADJ_SPREADING     .SpreadingAdjustment                      ADJUSTMENT_ID
//	      ├ T_CLAIM_ADJ_QUOTA_SHARE   .SpreadingQuotaShare                      ADJUSTMENT_ID
//	      └ T_CLAIM_FAC_RETRO         .FacRetroList                             CLAIM_ID + ADJUSTMENT_ID (ADD, K5)
//
// Tabel `T_CLAIM_*` milik Claim Prop dipakai ulang; baris FAC mengisi CLAIM_ID DAN induk tingkatnya sekaligus (nol
// MODIFY - rencana induk ganda STRUKTUR 20-09 ditolak OQ-CFI-01). Karena UNIQUE (CLAIM_ID, NOURUT) tabel Claim Prop,
// NOURUT baris FAC berurut SE-KLAIM (lintas item), bukan per item.
//
// ⛔ Medan salinan polis (`LengkapiObjek`) dan medan TURUNAN (`HitungTurunan`) tidak punya kolom.

// Golongan adalah kategori tipe logis sebuah kolom.
type Golongan string

const (
	GolTeks         Golongan = "teks"
	GolKode         Golongan = "kode"    // teks - nol di depan bermakna
	GolPenanda      Golongan = "penanda" // teks - "" berbeda dari "0"
	GolUang         Golongan = "uang"    // NUMBER(38,10)
	GolPersen       Golongan = "persen"  // NUMBER(38,10) - 12.5 = 12,5 persen
	GolTanggal      Golongan = "tanggal" // DATE, tanggal saja
	GolTanggalWaktu Golongan = "tanggal-waktu"
)

// Desimal - golongan bertipe NUMBER(38,10).
func (g Golongan) Desimal() bool { return g == GolUang || g == GolPersen }

// Tanggal - golongan bertipe DATE.
func (g Golongan) Tanggal() bool { return g == GolTanggal || g == GolTanggalWaktu }

// Kolom memetakan satu properti ke satu kolom.
type Kolom struct {
	// Properti - jalur relatif pyWorkPage (tabel kepala) atau nama properti anggota daftar (tabel anak).
	Properti string
	Kolom    string
	Golongan Golongan
	// Panjang - panjang VARCHAR2; 0 untuk golongan non-teks.
	Panjang int
	// Baru - kolom lahir di migrasi modul ini (ALTER ADD / CREATE); false = kolom tabel bersama yang dipakai ulang.
	Baru bool
	// MilikKomite - kolom yang hanya ditulis keputusan komite (tahap 2, `SaveAccept_ACT` / `KomitePostAdjustment`):
	// dibaca dan diisi saat baris lahir, tetapi TIDAK ikut UPDATE simpan halaman klaim - tanpa itu simpan klaim yang
	// membaca halaman sebelum komite memutuskan menimpa keputusan itu (temuan review 10-10-2026).
	MilikKomite bool
}

// Tabel - satu simpul pohon katalog.
type Tabel struct {
	Nama string
	// Daftar - nama daftar di dalam baris induk (`ObjectItemList`); simpul akar anak kepala memakai jalur halaman
	// penuh (`ClaimData.ObjectList`).
	Daftar string
	// KolomInduk - kolom penunjuk induk langsung (CLAIM_ID / OBJECT_ID / OBJECT_ITEM_ID / ADJUSTMENT_ID).
	KolomInduk string
	// DenganKlaim - baris juga mengisi CLAIM_ID (tabel Claim Prop: CLAIM_ID NOT NULL + UNIQUE (CLAIM_ID, NOURUT));
	// NOURUT berurut se-klaim.
	DenganKlaim bool
	// Jenis - nilai kolom JENIS; "" = tabel tanpa pembeda.
	Jenis string
	// Stabil - ID baris stabil (properti `PropID` membawa kunci tersimpan); baris lain ditulis ulang setiap simpan.
	Stabil bool
	Kolom  []Kolom
	Anak   []*Tabel
}

func kTeks(p, k string, n int) Kolom {
	return Kolom{Properti: p, Kolom: k, Golongan: GolTeks, Panjang: n}
}
func kKode(p, k string, n int) Kolom {
	return Kolom{Properti: p, Kolom: k, Golongan: GolKode, Panjang: n}
}
func kPenanda(p, k string) Kolom {
	return Kolom{Properti: p, Kolom: k, Golongan: GolPenanda, Panjang: 16}
}
func kUang(p, k string) Kolom   { return Kolom{Properti: p, Kolom: k, Golongan: GolUang, Panjang: 0} }
func kPersen(p, k string) Kolom { return Kolom{Properti: p, Kolom: k, Golongan: GolPersen, Panjang: 0} }
func kTgl(p, k string) Kolom    { return Kolom{Properti: p, Kolom: k, Golongan: GolTanggal, Panjang: 0} }
func kTglWaktu(p, k string) Kolom {
	return Kolom{Properti: p, Kolom: k, Golongan: GolTanggalWaktu, Panjang: 0}
}

// baru - kolom yang lahir di migrasi modul ini.
func baru(k Kolom) Kolom { k.Baru = true; return k }

// milikKomite - kolom yang hanya ditulis komite (lihat Kolom.MilikKomite).
func milikKomite(k Kolom) Kolom { k.MilikKomite = true; return k }

// CD - awalan jalur `pyWorkPage.ClaimData`.
const CD = "ClaimData."

// PropID - properti tersembunyi baris tabel ber-ID stabil (objek, item, adjustment) yang membawa ID tersimpan.
const PropID = "ID"

// PropKomiteID - properti baris adjustment yang membawa KOMITE_ID tersimpan (baca saja dari layar).
const PropKomiteID = "KomiteID"
