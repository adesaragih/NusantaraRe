# 23: Berkas migrasi 78 tabel flat (rentang 180–219) — ⏸ DITAHAN

> ⚠️ **Disusun agent dari spec/DDL atas perintah work owner — bukan hasil `/to-tickets`.** Perintah: arahan work owner
> 02-10-2026 langkah 3 ("ditulis, tidak dijalankan"); **ditahan** oleh pilihan work owner butir 66 ("Flatten dulu,
> migrasi ditahan").

**What to build:** berkas migrasi `modul/nbfacin/backend/migrations/1NN_*.sql` (+ `_down.sql`) yang membuat 78 tabel
flat dari `D:\migrasi\RNM\OUTPUT\08-flat\DDL-tabel-flat-draf.sql` — **ditulis, tidak dijalankan**. ⛔ Menjalankan
migrasi ke Oracle mana pun (`-migrate`) hanya oleh work owner sendiri (butir 65).

**Blocked by:** keputusan tim inti atas presisi tabel flat (butir 66) — ~~tipe delapan kolom~~ diputuskan butir 68.1

**Status:** needs-info — ⏸ ditahan; tidak satu berkas migrasi pun ditulis

## Yang menahan

1. **Presisi.** `[terverifikasi]` DDL draf memakai **415** kolom `NUMBER` polos (156 di antaranya `ID`/`*_ID`) +
   **1** `NUMBER(3)` + **50** `NUMBER(5)`; penjaga inti `TestNolNumberTanpaPresisi` hanya mengizinkan `NUMBER(38,8)`,
   `(5)`, `(10)`, `(19)`. Uang `NUMBER(38,8)` **memotong** nilai yang terbukti ada: premi FIRE 20 desimal (fixture),
   dan angka berkoefisien **> 38 digit** di 115 contoh (`T_COVERAGELIST.NET_RATE` 185 nilai, lima kolom pembayaran — tiket
   22 bab *Pengukuran*). Penjaga inti **tidak diubah** (butir 66).
   **Usulan untuk tim inti** (butir 68.2): `docs/USULAN-PRESISI-TABEL-FLAT.md` — `NUMBER(38,8)` membulatkan **10.876**
   nilai di **82** kolom pada 115 contoh korpus (458 / 20 di 5 fixture).
2. ~~Tipe delapan kolom yang isinya bukan angka~~ → ✅ **butir 68.1**: teks apa adanya, penyimpangan sadar dari DDL draf
   (tiket 22 bab *Penyimpangan sadar*). Panjang `VARCHAR2`-nya ditetapkan saat berkas migrasi ditulis.
3. **Dua sumber skema** — sudah terjelaskan (tiket 22): Lintas-Siklus = DDL draf persis; `Claude outputs` himpunan
   bagiannya. Bila tetap dipakai bersamaan, keduanya ditulis, tidak dipilih (arahan work owner).

## Bila dibuka — rencana yang sudah diukur

- **Rentang `180-219` (40 nomor)** cukup bila satu berkas memuat beberapa tabel; 78 tabel × satu berkas per tabel
  **tidak** muat. ⛔ Bila dibutuhkan lebih, **jangan** ubah rentang di `MODUL.md` — laporkan (arahan work owner).
- Urutan berkas mengikuti urutan muat spec 11 butir 15 (kerja → akar → anak, sampai 8 tingkat); FK tunggal 64,
  induk ganda tanpa FK (12 menurut Daftar Relasi / 13 menurut Jalur Sumber — tiket 22 *Butir terbuka* 3).
- Tabel baru masuk `docs/STRUKTUR-TABEL-NB-FACIN.md` (`TestKolomDDLCocokDenganStruktur`); bukan tabel warisan.
- 24 kolom mata uang `DEFAULT 'UNKNOWN' NOT NULL` (K-069) apa adanya.
- **Amandemen rancangan butir 70** ikut ditulis (tiket 22 bab *Amandemen rancangan*): tabel `T_ADDITIONALSHIP`, kolom
  `T_COVERAGELIST.COVERAGE_INITIAL`, `T_CURRENCYLIST.CURRENCY_REF_ID`, `T_FR_CURRENCYLIST.POLICY_TSI` (tipe = keputusan
  presisi) — 79 tabel; dan **butir 72**: 48 kolom penunjuk teks mentah `VARCHAR2(50)` di 16 tabel (daftar lengkap:
  `loader/amandemen.go`, `amandemenPenunjuk`) — 1.390 kolom. ⚠️ Tipe kolom penunjuk tidak seragam di rancangan
  gabungan: kolom kunci lama ber-`NUMBER` (mis. `T_FR_ANEKALIST.IDX_LOCATION`), kolom penunjuk butir 72 `VARCHAR2(50)`
  (mis. `T_ANEKALIST.IDX_LOCATION`) — dicatat, tidak diseragamkan agent. ⛔ Kolom `IsCedingConfirm` di `POOLDATA.HISTORYAKSEPTASIPRODUCTION` **bukan** milik rentang ini:
  tabel lama, diubah DBA (P4).

## Comments
