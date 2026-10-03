# 18: Tipe kolom, konversi masuk, dan presisi uang

> ## ⭐ PENAHAN GUGUR — 23 September 2026 sore
>
> `[keputusan work owner]` *"Selesaikan, jangan jadi permasalahan."* ⭐ Presisi dinaikkan ke **`NUMBER(38,8)`** — **30 digit di depan koma**, delapan di belakang, batas tertinggi Oracle. Dasarnya sapuan korpus: ambang dagang nyata sudah **tepat di batas** 12 digit *(`181500000000.00` · `150000000000.00`)*, dan ada sentinel **17 digit** *(`99999999999999999.99`)*. Penjumlahan lintas mata uang dapat melewati keduanya. ⭐ `NUMBER` di Oracle berpanjang **berubah-ubah** — hanya digit bermakna yang tersimpan, sehingga pelebaran ini **tidak memakan ruang tambahan**. Butir ditutup oleh bukti, bukan oleh DBA.
>
> ⭐ **`blocked` → `ready-for-agent`.** ⛔ **Nol butir `[data DBA]` tersisa di tiket ini.**
>
> Rinciannya: `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`.

---


**Status:** sebagian *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026 sore)*)*
~~**Blocked by:** **16** · ⛔ `[data DBA]` **presisi fisik belum diuji terhadap nilai terbesar** — dua belas digit di depan koma belum dibuktikan cukup~~ ⛔ **gugur 23-09-2026 sore**
**Menutup:** NB AC **12–25** *(14 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-14..ID-20

## Hasil & nilai pengguna

Nilai yang masuk dari dokumen lama seluruhnya bertipe teks — termasuk uang, tanggal, dan penanda.
Tiket ini menetapkan bagaimana teks itu menjadi kolom bertipe, **sekali saat masuk**, bukan setiap
kali dibaca.

⛔⛔ **Yang paling mudah salah, dan akibatnya besar:** kode `"006"` yang disimpan sebagai bilangan
lalu dibaca balik menjadi `"6"` akan **lolos uji pulang-pergi** tetapi memecahkan penggolong jenis
usaha 36 baris — dan kegagalannya diam.

## Yang dibangun

| Golongan | Ketetapan |
| --- | --- |
| **uang dan persen** | angka presisi tetap; ⛔ **tidak pernah** bilangan mengambang |
| **kode** | ⭐ **tetap teks** — nol di depan membawa makna |
| **penanda** | ⭐ **tetap teks** — kosong adalah keadaan sah yang **berbeda** dari nol |
| **tanggal** | dua format masuk: delapan digit, dan cap waktu bersufiks zona |
| **teks kosong** | menjadi **kosong**, bukan nol |

⚠️ Pembandingan dua nilai uang memakai **toleransi atau bentuk terbulatkan**, bukan kesamaan
persis — data produksi terbukti membawa galat pecahan.

## Batas — yang TIDAK termasuk

⛔ `CREATE TABLE` — presisi fisik dicocokkan DBA **di dalam** tiket ini.
⛔ Pemecahan dokumen menjadi baris — tiket **19**.

## Cara mengujinya

Lewat seam `repository`, dan ⚠️ **sebagian test memeriksa nilai kolom langsung** — pulang-pergi
saja tidak cukup, sebab tulis dan baca yang sama-sama salah simetris tetap hijau.

⭐ Uji yang wajib: kode `"006"` disimpan lalu **dibaca dari kolomnya**, dan penggolong jenis usaha
dijalankan atasnya — hasil bawaan berarti gagal.

## Acceptance criteria

- [ ] 🟡 **AC 12** — kode tiga digit berawalan nol tersimpan dan terbaca utuh
- [ ] 🟡 **AC 13** — kode dua digit berawalan nol tersimpan utuh
- [x] **AC 14** — penggolong jenis usaha menemukan barisnya, bukan nilai bawaan
- [ ] 🟡 **AC 15** — penanda kosong tersimpan kosong, **bukan** nol dan bukan tak-bernilai
- [x] **AC 16** — penanda bernilai nol menghasilkan keputusan ditolak; nilai lain disetujui
- [x] **AC 17** — teks kosong pada medan uang tersimpan tak-bernilai
- [x] **AC 18** — teks kosong pada medan tanggal tersimpan tak-bernilai
- [x] **AC 19** — nilai uang berdesimal sembilan **dibulatkan pada desimal kedelapan**, bukan dipotong ke dua
- [x] **AC 20** — kolom uang berskala **delapan desimal**, ⭐ **tiga puluh digit di depan koma** *(`NUMBER(38,8)`; semula ~~dua belas~~ — dinaikkan 23-09-2026 sore)*
- [x] **AC 21** — tanggal delapan digit terurai benar *(putaran 2, pemuat tiket 22: `models.BacaTanggalLama` — `TestBacaTanggalLama`)*
- [x] **AC 22** — cap waktu bersufiks zona terurai benar *(putaran 2: `20170930T170000.000 GMT` → `2017-10-01 00:00:00` Asia/Jakarta — `TestBacaTanggalLama`; tanggal ambigu tidak ditebak, K15 — `TestTanggalAmbiguTidakDitebak`)*
- [x] **AC 23** — pengurutan menurut tanggal menghasilkan urutan kronologis, bukan leksikal
- [x] **AC 24** — nol kolom uang bertipe mengambang
- [ ] 🟡 **AC 25** — pembandingan uang memakai toleransi, bukan kesamaan persis

## ⛔ Kenapa tiket ini `blocked`

~~`[data DBA]` **Dua belas digit di depan koma belum diuji terhadap nilai terbesar.**~~ ⛔ **BUTIR GUGUR 23-09-2026 sore.** `[keputusan work owner]` *"Selesaikan, jangan jadi permasalahan."* Presisi dinaikkan ke **`NUMBER(38,8)`** — tiga puluh digit di depan koma. Dasarnya sapuan korpus: ambang dagang nyata **tepat di batas** dua belas digit *(`181500000000.00`)*, dan ada sentinel **tujuh belas digit** *(`99999999999999999.99`)*. ⭐ `NUMBER` di Oracle berpanjang berubah-ubah, jadi pelebaran ini **tidak memakan ruang tambahan**.

⚠️ **Dan satu pertentangan di dalam spec, dicatat di sini karena spec tidak boleh disunting:**
~~`[terverifikasi]` spec EDM **AC 49** menuntut sembilan desimal, bertentangan dengan AC 58.~~ ✅ **DITUTUP 23-09-2026 sore.** AC 49 diselaraskan ke **delapan**; bunyi lamanya dikutip di spec EDM. ⭐ Pertentangan ini ditemukan **pelaksana ronde tiket**, dan laporannya terbukti tepat..

## ⭐ Penerapan KEPUTUSAN-RONDE-12 dan catatan — 2026-10-03

- **Butir 8** — nomor generasi `PRODKE NUMBER(10) DEFAULT 0` (bilangan bulat). ID-14: `INSTALLMENT_NO`
  juga bilangan bulat `NUMBER(10)`.
- **AC 15** — Oracle menyimpan `''` sebagai NULL; penanda kosong dibaca kembali sebagai `""` (setara di
  halaman, tidak setara di SQL) — sebagian.
- **AC 21-22** (format dokumen lama `YYYYMMDD`, cap waktu ` GMT`) milik pemuat dokumen lama — tiket 22,
  belum dibangun.
  ⛔ **RALAT putaran 2 (03-10-2026):** bunyi lama *"belum dibangun"* → **dibangun** (tiket 22,
  `models.BacaTanggalLama`).
