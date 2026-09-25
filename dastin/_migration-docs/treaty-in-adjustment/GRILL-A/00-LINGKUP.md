> Modul  : Treaty In Adjustment · Ronde A · 2026-09-24
> Peran  : interogator
> Masukan: `PENGETAHUAN.md` (2026-09-23, dikoreksi 2026-09-24) · `METODE-GRILLING.md` · ADR induk 0036, 0040, 0048, 0049, 0052, 0055 · `SEAM-ADJUSTMENT.md` · `KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` · `TEMUAN-ADJUSTMENT-DITUNDA.md` · 379 XML `Treaty In Adjustment/` + 329 XML `Treaty In/` (ekspor 2026-09-03)
> Status : TERBUKA
> Sifat  : TAMBAH-SAJA

# 00 · LINGKUP RONDE A

## 1. Cabang yang digrilling

**Cabang A — batas lingkup, dan rekonsiliasi ADR modul induk.** Ronde ini menjawab dua hal yang
mendahului segalanya:

1. apakah spesifikasi Adjustment ditulis ke dalam spesifikasi Treaty In atau di sampingnya;
2. ADR induk mana yang berlaku apa adanya, mana yang berlaku dengan penyesuaian, dan mana yang
   bertentangan — **satu ADR per giliran**, dengan bentuk premis -> uji terhadap koreksi
   `PENGETAHUAN.md` -> klasifikasi.

Ditambah satu pertanyaan yang, menurut `METODE` §8.6 Q1, adalah satu-satunya yang dapat membalik
putusan "addendum adalah versi".

## 2. Bukti yang dipakai

| Bukti | Dipakai untuk |
|---|---|
| Ekspor XML kedua modul | seluruh pernyataan fakta tentang sistem lama |
| `PENGETAHUAN.md` | ringkasan AS-IS; **dikoreksi tiga kali dalam ronde ini**, lihat `01-TEMUAN` |
| ADR induk | premis yang direkonsiliasi |
| `METODE-GRILLING.md` | aturan kerja; Bagian VIII disidangkan di `02-SIDANG` |
| Prosedur `POOLDATA.PEGA_M_TREATY_IN_EDM` | perilaku sisip-atau-perbarui |

## 3. Bukti yang sengaja TIDAK dipakai

| Tidak dipakai | Alasan |
|---|---|
| Data produksi | tidak tersedia; seluruh pertanyaan yang menuntutnya dicatat sebagai `UA-x`, dan **tidak satu pun memblokir rancangan** (`METODE` §7.2) |
| Aturan **Field Value** `EDMState` / `EDMMaterialType` | tidak ikut ter-ekspor — `EXP-1` |
| Wawancara bisnis | belum dijalankan; pertanyaannya disusun sebagai **daftar untuk dibantah** (`METODE` §4.7) |
| Cabang B-K | belum dibuka; temuan yang menyentuhnya **dititipkan**, tidak diadili |

## 4. Aliran: yang terbuka dan yang beku

| Beku sebelum ronde ini | Dasarnya |
|---|---|
| Addendum adalah **VERSI**, tidak ada entitas `PENYESUAIAN` | `METODE` §8.3 |
| Bentuk sempit `NILAI_SELISIH`, menggantung pada `VERSI_KONTRAK` | `METODE` §8.3 |
| Empat aturan bisnis pemilik proses | `METODE` §8.5 |

| Terbuka saat ronde ini dibuka |
|---|
| Susunan spesifikasi (satu model atau dua) |
| Status keenam ADR induk terhadap modul ini |
| Apakah §8.6 Q1 membalik putusan "versi" |

## 5. Deret nomor ronde ini

* Temuan: `NA-01 … NA-06`
* Pertanyaan: `QA-1 … QA-5` (A1 s.d. A5)
* Putusan: `GRL-01 … GRL-04`

## 5a. Aturan bukti ronde ini

Ditambahkan 24 September 2026 sesudah `MA-04`:

> **Kemunculan di dalam `pyIncludedRuleXML` atau `pyRuleVersionsList` tidak dihitung dan tidak
> dipakai sebagai bukti.** Keduanya memuat salinan aturan lain — dan menurut §1.2 dapat memuat
> **versi lama** aturan yang sama — yang ikut terbungkus saat ekspor. Salinan bukan rujukan,
> apalagi panggilan (`METODE` §2.0a).

Setiap hitungan yang dilaporkan ronde ini menyebut **badan** dan **salinan** terpisah.

## 6. Aturan berhenti

Ronde A ditutup ketika keenam ADR induk sudah diklasifikasikan dan §8.6 Q1 terjawab. Yang tersisa
sesudah itu **bukan pertanyaan cabang A** dan dititipkan ke cabangnya — bukan digali di sini.

Anggaran seluruh grilling ditetapkan **±40 pertanyaan** sebagai batas atas sementara, ditinjau
ulang pada titik periksa sesudah ronde A ditutup (`METODE` §7.1).
