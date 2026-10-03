# 23: Berkas migrasi 78 tabel flat (rentang 180–219) — ⏸ DITAHAN

> ⚠️ **Disusun agent dari spec/DDL atas perintah work owner — bukan hasil `/to-tickets`.** Perintah: arahan work owner
> 02-10-2026 langkah 3 ("ditulis, tidak dijalankan"); **ditahan** oleh pilihan work owner butir 66 ("Flatten dulu,
> migrasi ditahan").

**What to build:** berkas migrasi `modul/nbfacin/backend/migrations/1NN_*.sql` (+ `_down.sql`) yang membuat 78 tabel
flat dari `D:\migrasi\RNM\OUTPUT\08-flat\DDL-tabel-flat-draf.sql` — **ditulis, tidak dijalankan**. ⛔ Menjalankan
migrasi ke Oracle mana pun (`-migrate`) hanya oleh work owner sendiri (butir 65).

**Blocked by:** keputusan tim inti atas presisi tabel flat (butir 66) — ~~tipe delapan kolom~~ diputuskan butir 68.1

**Status:** needs-info — ⏸ ditahan; tidak satu berkas migrasi pun ditulis

~~**Ditunggu oleh:** tiket 29 `POST /api/nbfacin/opportunity` (Create opportunity) — keputusan work owner 03-10-2026
butir 74.1 "Tunggu tabel flat (tiket 23)": tanpa tabel sementara; nomor case NB lanjut dari nomor terakhir Pega (butir
74.2).~~ ⛔ **Diralat butir 76** (03-10-2026): 74.1 diambil atas premis keliru (`T_WORK_POLIS` sudah ada, milik
premiumlistlife). Tiket 29 kini **tidak** menunggu tiket ini — case NB ditulis ke `T_WORK_POLIS` yang ada + tabel
`T_NB_OPPORTUNITY` (migrasi 180) + `SEQ_WORK_POLIS_NB` (181). Tiket ini tetap ditahan (presisi), dengan rancangan
`T_WORK_POLIS` diselaraskan di bab di bawah.

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
  `loader/amandemen.go`, `amandemenPenunjuk`) — 1.390 kolom (butir 76: skema loader 1.391 dengan `T_WORK_POLIS.LINI`, kolom yang sudah ada — tidak dibuat ulang). ⚠️ Tipe kolom penunjuk tidak seragam di rancangan
  gabungan: kolom kunci lama ber-`NUMBER` (mis. `T_FR_ANEKALIST.IDX_LOCATION`), kolom penunjuk butir 72 `VARCHAR2(50)`
  (mis. `T_ANEKALIST.IDX_LOCATION`) — dicatat, tidak diseragamkan agent. ⛔ Kolom `IsCedingConfirm` di `POOLDATA.HISTORYAKSEPTASIPRODUCTION` **bukan** milik rentang ini:
  tabel lama, diubah DBA (P4).

## Penyelarasan dengan `T_WORK_POLIS` yang ada — butir 76 (keputusan work owner 03-10-2026)

`[terverifikasi]` `T_WORK_POLIS` **sudah dibuat** premiumlistlife: `050_t_work_polis.sql` (ID VARCHAR2(32) NOT NULL PK, LINI,
POSITION, STATUS VARCHAR2(255)), `057` (+FLAG_ONGOING_POLICY, SEQ_WORK_POLIS), `059` (+STATUS_WORK, COVER_KEY + FK diri,
CREATE_OP VARCHAR2(64), CREATE_OP_NAME VARCHAR2(128), TGL_CREATE, TGL_UPDATE DATE; STATUS dibuang), `063` (STATUS →
STATUS_WORK). K-064: tabel yang **sama**; Fac In menyambung. ⛔ Tiket ini **tidak** membuat `T_WORK_POLIS`; kolom yang belum
ada ditambah lewat `ALTER TABLE … ADD` (nullable, K-064 konsekuensi 2). `T_GENERAL_POLIS` belum dibuat modul mana pun
(`[terverifikasi]` grep seluruh `*.sql` 03-10-2026: nol) — tetap dibuat tiket ini, dengan `ID VARCHAR2(32)`.

Pemetaan 18 kolom rancangan (`loader/skema_gen.go` = DDL draf) — dihitung dua cara: 1 sama + 4 digabung + 13 tambah = 18,
dan `TestAmandemenWorkPolis` (19 kolom skema = 18 + `LINI`):

| Kolom rancangan | Tipe rancangan | Kolom yang ada | Perlakuan (butir 76) |
| --- | --- | --- | --- |
| `ID` | NUMBER | `ID` VARCHAR2(32) NOT NULL | **sama**, tipe mengikuti yang ada (76.1); isi = `NO_WORK` (pyID, mis. `NB-184351`) |
| `IDPEGA` | VARCHAR2(50) | — | tambah |
| `JENIS_WORK` | VARCHAR2(10) | — | tambah |
| `NO_WORK` | VARCHAR2(20) | — | tambah (isinya sama dengan `ID`; tetap kolom tambahan, 76.1) |
| `POSISI` | VARCHAR2(100) | `POSITION` VARCHAR2(255) | **digabung** (76.4) |
| `NOURUT` | NUMBER(3) | — | tambah — ⛔ tipe tersangkut presisi (penjaga tidak mengizinkan `NUMBER(3)`) |
| `PUTARAN` | NUMBER | — | tambah — ⛔ presisi |
| `STATUS_PROSES` | VARCHAR2(20) | `STATUS_WORK` VARCHAR2(255) | **digabung** (76.4) |
| `STS_KONVERSI` | NUMBER | — | tambah — ⛔ presisi |
| `TGL_KONVERSI` | DATE | — | tambah |
| `TGL_INPUT` | DATE | `TGL_CREATE` DATE | **digabung** (76.4) |
| `USERNAME` | VARCHAR2(50) | `CREATE_OP` VARCHAR2(64) | **digabung** (76.4) |
| `DATE_TO_UW` | VARCHAR2(30) | — | tambah |
| `ID_NEW_BISNIS` | VARCHAR2(20) | — | tambah |
| `IS_FLAG_REJECT` | VARCHAR2(10) | — | tambah |
| `PIC` | VARCHAR2(100) | — | tambah |
| `POLICY_STATUS` | VARCHAR2(20) | — | tambah |
| `SUBMIT_TO_UW` | VARCHAR2(30) | — | tambah |

Kolom yang ada di luar rancangan: `LINI` — diisi `'FAC'` (76.2); `FLAG_ONGOING_POLICY`, `COVER_KEY`, `CREATE_OP_NAME`,
`TGL_UPDATE` — milik premiumlistlife, **tidak** diisi loader (pengisiannya untuk data lama: `belum terverifikasi`).
Penggabungan tidak menyempitkan satu kolom pun. **Rembetan 76.1:** `T_GENERAL_POLIS.ID` (PK bersama) dan `PARENT_ID` 10 tabel
berinduk `T_GENERAL_POLIS` (`T_CARGOLIST`, `T_CEDINGCEDANTLIST`, `T_CURRENCYLIST`, `T_FACRETRODETAILS`, `T_FACRETROLIST`,
`T_LOCATIONLIST`, `T_PERSONLIST`, `T_QUOTATIONDATA`, `T_SCORINGRISK`, `T_VEHICLELIST`) menjadi `VARCHAR2(32)` — dihitung
dua cara (himpunan tabel = 10, baris jalur = 10; nol tabel berinduk campuran). Diterapkan di `loader/amandemen.go`
(`selaraskanWorkPolis`) + `aturan.go`; mutasi 67/67 tertangkap.

⚠️ **Terbuka:** `ALTER` atas tabel milik premiumlistlife dari rentang 180–219 — koordinasi dengan pemiliknya; kotak masuk
PremiumList membaca **seluruh** `T_WORK_POLIS` tanpa saringan `LINI` (lihat register butir 76, risiko R1).

## Comments
