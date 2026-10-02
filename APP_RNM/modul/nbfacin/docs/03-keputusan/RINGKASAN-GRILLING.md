# Ringkasan Sesi Grilling — 15–16 September 2026

Tiga ronde, 12 pertanyaan, **7 keputusan baru** (K-007…K-013) dan **6 ADR**. Tiga blocker inti
tertutup dari korpus tanpa eskalasi.

---

## 1. Hasil terbesar: tiga blocker inti tertutup tanpa bertanya ke siapa pun

Empat pertanyaan yang ditetapkan sebagai prioritas — Q1 arti `ProposalAcceptStatus`, Q2 uang sebagai
string, Q3 satuan `.Rate`, Q14 dasar akseptasi renewal — **seluruhnya selesai**, dan hanya satu yang
memerlukan keputusan work owner.

### Q1 · Arti `ProposalAcceptStatus` — ✅ tertutup dari korpus

`[terverifikasi]` Ditemukan **dekoder literal eksplisit** di
`NB FacIn\Activity\SaveViewSuggest.xml`, hadir identik di ketiga folder:

```
@if(.Approval="1","Accept", @if(.Approval="2","Reject", @if(.Approval="3","Ask",
@if(.Approval="4","Banding", @if(.Approval="5","Reject Ceding", @if(.Approval="6","Ask Ceding",
@if(.Approval="7","Decline", @if(.Approval="9","Revise",""))))))))
```

| Nilai | Arti | Penguat |
| ---: | --- | --- |
| `1` | Accept | `When\IsUWAccepted` · SQL menulis `'ACCEPT'` |
| `2` | Reject | subjek email `"Rejected - ID Pega : "` · SQL `'REJECT'` |
| `3` | Ask | subjek email `"Ask- ID Pega : "` · SQL `'ASK'` |
| `4` | **Banding** | penulisnya menyalakan `IsBanding="true"` · SQL `'BANDING'` |
| `7` | Decline | subjek `"Decline Offer…"` · SQL `'DECLINE'` · **default** `IsUWAccepted` |
| `9` | Revise | `pyDeleteMemo`: *"added status 9, revise"* |

`5` dan `6` **ada** tetapi milik ranah ceding (`Reject Ceding` / `Ask Ceding`), di properti terpisah,
tidak pernah menulis `ProposalAcceptStatus`. `8` nol jejak di seluruh 6.071 berkas.

⚠️ Dugaan lama **terbantah**: nilai `4` memang **banding**, bukan fac out. `When\IsFacout` yang
mengujinya **salah nama**.

**Sisa yang masih terbuka:** baris keputusan `IsUWAccepted` tetap tidak terekspor — kosakatanya
pasti, logika pemilihnya belum. Ditangani K-008.

### Q3 · Satuan `.Rate` — ✅ tertutup, dan premisnya terbantah

`[terverifikasi]` Bukan per mille **atau** persen: **keduanya, ditentukan lini bisnis**.

| COB | Skala |
| --- | --- |
| FIRE, PA | **‰** |
| ANEKA, BONDING, GOLF, MARINE CARGO, MBU | **%** |

Bukti penutup: satu berkas (`CountGPWMarinePAMbu_Act.xml`) memuat tiga COB berdampingan dengan dua
satuan. Diperkuat label layar yang menyebut satuannya harfiah — `‰ Standard Rate` dengan karakter
U+2030 pada layar FIRE.

**Bahayanya berbalik arah:** risikonya bukan salah memilih satuan, melainkan **menyeragamkannya**.
Ditangani K-010.

### Q14 · Dasar akseptasi renewal — ✅ tertutup

`[terverifikasi]` **TSI penuh.** Dua lapis bukti bebas: cabang selisih digerbangi kesetaraan string
eksak `StatusBusiness=="3"`; dan `OldData` tidak pernah terisi di luar endorsement — satu-satunya
rule pengisinya hanya ada di folder EDM.

### Q2 · Uang sebagai string — ✅ sudah diputuskan piagam proyek

Bukan pertanyaan terbuka: `CLAUDE.md` §1 dan §4.1 sudah menetapkan **reproduksi apa adanya**, catat
sebagai temuan. Yang tersisa hanya eskalasi nilai ambang yang *seharusnya* benar.

---

## 2. Keputusan yang diambil

| # | Keputusan | ADR |
| --- | --- | --- |
| **K-007** | Rekonsiliasi **eksak**, nol selisih; bertahap sesuai ketersediaan ekspor produksi | 0001 |
| **K-008** | Satu flag fase: `panic` saat paralel run, `decline`+log saat produksi | 0002 |
| **K-009** | Mesin akseptasi ditulis sekarang; fixture **menjadi kontrak** data ke DBA | 0003 |
| **K-010** | `Money` dan `Ratio` tipe terpisah; skala melekat pada nilai | 0004 |
| **K-011** | Presisi pembulatan literal per-langkah, bukan registry | 0005 |
| **K-012** | Mata uang `Unknown` eksplisit; `panic` hanya lintas-MU dan tulis Oracle | 0006 |
| **K-013** | Urutan kerja; keputusan siapa menulis kode **ditunda** | — |

---

## 3. Dua koreksi yang diterima work owner selama sesi

Keduanya menghentikan perubahan perilaku yang hampir masuk diam-diam:

1. **Default `decline` adalah perilaku terekam, bukan celah pengetahuan.** Usulan menggantinya
   dengan `panic` di semua fase ditarik — `<pyDefaultResult>decline</pyDefaultResult>` terbaca
   langsung dari `IsUWAccepted`. Menggantinya melanggar `CLAUDE.md` §1. Diselesaikan lewat pemisahan
   fase (K-008).
2. **"Hilang dari ekspor" ≠ "usang"** (K-006, ronde sebelumnya). Enam cabang alur yang memanggil rule
   yang tidak ada di ekspor **ditangguhkan, bukan dibuang** — karena ekspor korpus berasal dari titik
   waktu berbeda antar folder.

---

## 4. Item kuesioner yang terkumpul

Tidak dipaksakan dijawab sekarang. Ini isi Prompt 4.

### Underwriting

| # | Pertanyaan |
| ---: | --- |
| U1 | **`2` Reject vs `7` Decline — apa beda operasionalnya?** Reject membatalkan binding (`ConfirmBinding=0`, `ReceivedRiSlip=false`) dan menjadi syarat jalur banding; decline tanpa efek samping dan merupakan hasil default. Apakah decline = penolakan final tanpa hak banding? |
| U2 | **`9` Revise hanya ditulis dari jalur binding** dan **tidak ada di siklus Endorsement**. Apakah revise memang mustahil pada endorsement? |
| U3 | **Apakah `5` / `6` pernah tertulis ke `ProposalAcceptStatus` di data produksi?** Bila tidak, kolom dapat dibatasi ke {1,2,3,4,7,9}. |
| U4 | **`When\IsFacout` menguji `ProposalAcceptStatus = 4` (banding).** Apakah banding memang memicu fac out, atau rule ini salah nama? |
| U5 | **Asimetri `SetBanding_ACT` di Endorsement** — meng-OR banding (`4`) dengan ask (`3`); di NB/RNW hanya `4`. Disengaja atau cacat? |
| U6 | **Asimetri renewal vs endorsement.** Renewal dengan kenaikan TSI kecil naik tangga atas **TSI penuh**; endorsement dengan dampak ekonomi sama hanya atas **delta**. Kebijakan atau cacat yang tertutup karena `OldData` memang tidak tersedia pada renewal? |
| U7 | **Nilai absolut pada selisih endorsement** — penurunan TSI diperlakukan seberat kenaikan. Disengaja? |

### Product / Aktuaria

| # | Pertanyaan |
| ---: | --- |
| P1 | **Tiga label rate bertentangan dengan rumusnya** — MBU berlabel `(‰)` padahal rumusnya %; view PA dan `LayerListDtl` berlabel `(%)` padahal rumusnya ‰. Label mana yang benar? (Tidak mengubah uang; rumus yang menghitung.) |
| P2 | Konfirmasi satuan rate yang dipakai **di slip** untuk MBU, PA, dan Layering. |

### Keamanan / IT

| # | Pertanyaan |
| ---: | --- |
| S1 | **`SendEmailPolicy` memuat alamat email tujuan dan BCC sebagai literal**, termasuk alamat perorangan. Melanggar `CLAUDE.md` §4.4 — harus menjadi konfigurasi. Siapa pemilik daftar penerimanya, dan bolehkah dipindahkan ke tabel/env? (Nilainya tidak disalin ke artefak mana pun.) |

### DBA — hasil pengukuran, bukan pertanyaan pendapat

| # | Item |
| ---: | --- |
| D1 | **Hitungan berapa banyak nilai produksi tiba tanpa mata uang.** Dihasilkan implementasi `Money.Currency = Unknown` (K-012). Angka ini mengubah pertanyaan dari dugaan menjadi ukuran, dan menentukan seberapa mendesak pengisian mata uang di sumber. |
| D2 | **Sampel nilai `RATE`, `MIN_RATE`, `MAX_RATE`** dari tabel sumber rate, dan kolom rate di tabel polis/produksi. Besaran tipikalnya membuktikan satuan secara langsung — menutup P1/P2 tanpa perlu pendapat. |

---

## 5. Yang tetap terbuka setelah sesi ini

| Butir | Pemilik | Catatan |
| --- | --- | --- |
| Baris keputusan `IsUWAccepted` | IT | ditangani sementara oleh K-008; ditutup oleh ekspor produksi tunggal |
| Enam cabang ditangguhkan (K-006) | IT | menunggu ekspor produksi tunggal |
| Isi `M_LIMIT_*` + ejaan `JABATAN` | DBA | fixture K-009 menjadi kontraknya |
| `ALL_SOURCE` 35 prosedur `POOLDATA` | DBA | menggerbangi paralel run end-to-end |
| DDL 112 tabel | DBA | menggerbangi tipe data |
| Isi 604 dropdown | Product | menggerbangi layar isian |
| Kamus `BusinessCode` (98) dan `BusinessOldId` (87) | Product | |
| ~~Arti `Type`, `EdmType`, `DeductibleType`, `TeamGroup`, `ProRateType`, `IsB2B`~~ | ~~Product~~ | ✅ **TERTUTUP K-029** (17 Sept 2026) — artinya terbaca dari property Pega di `DDL\`; `IsB2B` = flag penanda. Tidak menunggu Product lagi |
| Siapa menulis kode Go/React | Tim | K-013, ditunda |

Daftar permintaan yang harus dieksekusi **selagi Pega masih hidup** ada di
`../_EKSTRAKSI-PEGA-SELAGI-HIDUP.md` — lima butir, satu di antaranya datanya musnah bersama Pega.
