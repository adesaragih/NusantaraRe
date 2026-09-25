# Serah terima — to-spec Treaty In Adjustment

**Tanggal:** 24 September 2026 · **Fase:** to-spec, selesai · **Berikutnya:** menunggu perintah

> **Modul ini tidak punya spesifikasi sendiri** (`GRL-01`). Hasilnya mendarat di **artefak induk**.
>
> **Ketujuh diff disetujui 24 September 2026. LIMA DITERAPKAN, DUA DIPARKIR** — keadaan per diff di
> `USULAN-DIFF-KE-INDUK.md`. **Bukan "selesai"**: selama dua masih diparkir, hasil to-spec belum
> seluruhnya ada di spesifikasi.

## 1. Gerbang selesai — diperiksa satu per satu

| Gerbang | Keadaan |
|---|---|
| penyebut tiga pohon cermin dihitung ulang; tidak ada jalur tanpa nasib; hitungan menutup | **YA** — **279** dalam lingkup, 388 di luar, menutup ke **667**. `DITUNDA = 0` |
| keempat pohon punya baris keputusan, termasuk yang sengaja hilang | **YA** — `OLDDATA`, `ActualValue`, `ValueBeforeProrate` dibuang dengan dasarnya; `ValueDifference` menjadi tabel selisih |
| tabel selisih berkunci `KUNCI_PADANAN` dengan mata uang di dalamnya; `ReinstatementPct` diadili | **YA** — dan `ReinstatementPct` menjadi **temuan**, `TDA-18` calon |
| `DOKUMEN_ADDENDUM` ditulis; akibat migrasinya masuk rencana peralihan | **YA** — §3.1 `STRUKTUR-ADDENDUM.md` |
| `SIFAT_MATERIAL_ADDENDUM`, `INV-69`, `INV-70`, `N-7a`…`N-7c+` mendarat di artefak induk | ~~sebagai USULAN~~ → **SELURUHNYA DITERAPKAN** 24 Sep 2026 — `D-4` dan `D-5` pagi, `D-1` sore sesudah `F-2` ditutup |
| empat kolom `TDA-17` ada; tiga sisanya bertanda "lubang induk, bukan dibuang" | **YA** — §5 `STRUKTUR-ADDENDUM.md`, penandanya ditulis supaya tidak dapat disalahbaca |
| dua prosedur terpetakan habis; nol `CREATE PROCEDURE`; nol kolom blob | **YA** — `PEMETAAN-PROCEDURE-ADDENDUM.md`; DDL induk diperiksa: **0** dan **0** |
| kalimat pembaca arsip tertulis | **YA** — §6 |
| serah terima memisahkan dua daftar | **YA** — §3 dan §4 |

## 2. Yang dihasilkan

| Berkas | Isi |
|---|---|
| `PENELUSURAN-POHON-CERMIN.md` | 279 jalur, empat nasib, `SISI` per akar dari selisih himpunan berkas |
| `STRUKTUR-ADDENDUM.md` | **delta** — keempat pohon, tabel selisih, `DOKUMEN_ADDENDUM`, materialitas, `TDA-17`, kalimat arsip |
| `PEMETAAN-PROCEDURE-ADDENDUM.md` | dua prosedur, tiga golongan, dua cacat |
| `USULAN-DIFF-KE-INDUK.md` | **tujuh diff** menunggu persetujuan |
| `tools/telusur-pohon-cermin.py` | perkakas, disimpan sebagai berkas |

## 3. Yang MENAHAN to-ticket

| # | Penahan | Siapa mencabutnya |
|---|---|---|
| **1** | ~~DUA dari tujuh diff masih diparkir~~ → **SATU.** Enam sudah mendarat 24 Sep 2026: `D-3`, `D-2`, `D-4`, `D-5`, `D-6`, dan — sore hari, ketika sesi to-spec induk **menutup `F-2`** — **`D-1`**. Yang tersisa **`D-7`**, menunggu **satu kalimat izin**; dasarnya sudah lewat, sebab `F-1` ternyata sudah diterapkan | **pemilik proses** (`D-7`) |
| 2 | tanggapan paket **`REV-1` … `REV-6`** | pemilik ADR induk |
| 3 | **`DB-20`** — titik beku materialitas | bisnis |
| 4 | **`DB-16a`**, **`DB-16b`** — dokumen: dikirim ke luar, dan tanggal berlakunya | bisnis |
| 5 | angka **presisi** dan bentuk **induk polimorfik** | gerbang sesi DDL induk |
| 6 | **tiga kolom `TDA-17` tanpa rumah** — lubang §10 induk | sesi to-spec **induk** |

> **Butir 4 punya tenggat nyata:** `KTV-2` mensyaratkan kolom tanggal berlaku **dicabut sebelum data
> masuk**. Sesudah migrasi berjalan, mencabut kolom berongkos.

## 4. Yang hanya MENGUKUR KERUSAKAN

| # | Uji | Yang diukurnya |
|---|---|---|
| `UA-3` | berapa baris warisan **melanggar** `INV-69`/`INV-70` | perlakuan atas yang melanggar — keputusan migrasi |
| `UA-18` | berapa daftar memuat **dua baris berkunci padanan sama** | apakah pemadanan selisih memasangkan baris yang salah |
| `UA-19` | berapa nomor dokumen **berulang** | apakah keunikan global menolak data sah |
| **Uji AN** | baris ganda `TREATYINDETAIL` | akibat prosedur tanpa cabang `ELSE` |
| — | ongkos pengisian ulang nomor dokumen dari arsip kertas | apakah sepadan |

## 5. Dua temuan baru — keduanya `SISI` IRISAN, keduanya diusulkan ke INDUK

| # | Temuan | Keadaan |
|---|---|---|
| **`TDA-18`** *(calon)* | **`ReinstatementPct` disalin, tidak dikurangi.** Sapuan seluruh korpus, kalibrasi lulus: 10 penugasan, **nol pengurangan**. Maka perubahan persentase reinstatement **tidak pernah menghasilkan baris selisih** | **belum melewati ronde TDA** |
| **`TDA-17`** | tiga kolom tanpa rumah, sebabnya `L-8` | diadili `GRILL-D`; **tiganya milik induk** |

> **Keduanya tidak diperbaiki diam-diam**, dan keduanya menyentuh `GRL-20`: versi yang **hanya**
> mengubah reinstatement akan terbaca **Non Material** di sistem lama dan **Material** di sistem
> baru. **Itu perubahan perilaku yang disengaja** — sistem baru menangkap yang lama lewatkan.

---

**To-ticket tidak dimulai.** Ia hanya dimulai atas perintah pemilik proses, dan butir 1 §3 ditagih
lebih dulu — tanpa diff yang disetujui, **tidak ada hasil to-spec yang benar-benar ada di
spesifikasi**.
