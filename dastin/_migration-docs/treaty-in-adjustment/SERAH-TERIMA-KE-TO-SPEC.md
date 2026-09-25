# Serah terima — grilling Adjustment ke to-spec

**Tanggal:** 24 September 2026 · **Fase yang ditutup:** grilling · **Fase berikutnya:** to-spec

---

## ✅ TO-SPEC ADJUSTMENT TERBUKA

Keempat prasyarat cabang K terpenuhi:

| # | Prasyarat | Keadaan |
|---|---|---|
| 1 | butir GRILL terkunci | **terpenuhi** — `GRL-11` … `GRL-20`; ronde D membatalkan `GRL-04` dan `GRL-12` dan menggantinya |
| 2 | paket `REV` **diserahkan** | **terpenuhi** — bunyinya diubah `KTV-4`; tanggapan **ditagih sebelum to-ticket** |
| 3 | `EXP-1` ditutup | terpenuhi 26 Sep |
| 4 | tiap butir KONFIRMASI / REKOMENDASI / DITUNDA punya tempat | terpenuhi — `PEMILAHAN-SISA-GRILLING.md` |

**Gerbang 0 Prompt B sudah lewat.** To-spec dijalankan **mulai langkah 1**.

---

## 1. Apa yang berubah sejak peta terakhir

| | |
|---|---|
| putusan baru | `GRL-19` (`DOKUMEN_ADDENDUM`), `GRL-20` (materialitas masukan, sakelar dua arah) |
| putusan **batal** | `GRL-04`, `GRL-12` — **dicoret, tidak dihapus** |
| keputusan tanpa verifikasi | `KTV-1` … `KTV-4` — **nol terverifikasi** |
| invarian diusulkan | `INV-69`, `INV-70`, dengan uji negatif **dan positif** |
| `TDA-17` | diadili — **17 TDA, 17 nasib** |

---

## 2. Yang DITAGIH sebelum **to-ticket** — bukan sebelum to-spec

| # | Yang ditagih | Dari siapa | Bila jawabannya datang |
|---|---|---|---|
| 1 | **tanggapan paket `REV-1` … `REV-6`** | pemilik ADR induk | bila satu ditolak **dengan akibat pada bentuk data**, bagian to-spec itu dibuka kembali. Yang paling mungkin `REV-3` — keadaan **adalah** kolom |
| 2 | **`DB-20`** — apakah materialitas beku sejak lahir | bisnis | jejaknya **dipangkas**; satu perubahan aturan, bukan skema (`KTV-1`) |
| 3 | **`DB-16a`** — apakah revisi internal pernah disertai dokumen | bisnis | satu **nilai enum ditambahkan** (`KTV-3`) |
| 4 | **`DB-16b`** — apakah dokumen punya tanggal berlaku sendiri | bisnis | kolomnya **dicabut sebelum data masuk** (`KTV-2`) |

> **Butir 4 punya tenggat yang nyata:** *sebelum data masuk*. Sesudah migrasi berjalan, mencabut
> kolom berongkos; sebelum itu, tidak.

---

## 3. Sisa — dipisah menurut pembacanya

### 3.1 Yang MENAHAN MODEL — **kosong**

Tidak ada. Itu pernyataan, bukan kekosongan: keempat butir yang sempat disebut menahan **diuji ulang
dan tidak satu pun menentukan letak kolom**.

### 3.2 Yang hanya MENGUKUR KERUSAKAN

| # | Uji | Yang diukurnya | Menahan apa |
|---|---|---|---|
| `UA-3` | berapa baris warisan **melanggar** `INV-69`/`INV-70` | perlakuan atas yang melanggar | **migrasi**, bukan spesifikasi |
| `UA-19` | berapa nomor dokumen **berulang lintas cedant** | apakah keunikan global menolak data sah | constraint `DOKUMEN_ADDENDUM` |
| — | ongkos pengisian ulang nomor dokumen dari arsip kertas | apakah sepadan | **rencana peralihan** |

### 3.3 Yang milik modul INDUK, bukan modul ini

| # | Sisa | Kenapa tidak dapat diputuskan di sini |
|---|---|---|
| **`DEDUCTIBLE2`, `PREMIUM_EARNED`, `ROL_PCT`** | **tidak punya rumah di skema baru** | lubang pada §10 **induk**, sebabnya `L-8`. Ketiganya hanya muncul di pohon cermin, jadi §10 **tidak pernah diperlihatkan** kepadanya |
| angka presisi · bentuk induk polimorfik | gerbang sesi DDL induk | |

---

## 4. Yang sesi to-spec Adjustment tinggal menulis

`CABANG-K-PEMETAAN-TO-SPEC.md` kini memuat **artefak induk mana, bagian mana, bentuk mendaratnya
apa** untuk: `GRL-11` … `GRL-20`, `KTV-1` … `KTV-4`, `INV-69`, `INV-70`, dan seluruh bahan to-spec
`T-1`…`T-3`, `B-1`…`B-4`, `C-1`…`C-3`, `I-1`, `I-2`, `E3a`, `E3b`.

Dan dua hal yang **sengaja tidak mendarat di artefak to-spec**, dicatat supaya tidak dikira
terlupa: kolom dokumen yang **kosong untuk seluruh baris warisan** (rencana peralihan), dan ketiga
kolom `TDA-17` yang **milik induk**.

> **Peta yang tidak lengkap memaksa sesi to-spec MENCARI.** Gunanya supaya ia **tinggal menulis**.

---

**To-ticket tidak dimulai.** Ia hanya dimulai atas perintah pemilik proses, dan keempat butir §2
ditagih lebih dulu.
