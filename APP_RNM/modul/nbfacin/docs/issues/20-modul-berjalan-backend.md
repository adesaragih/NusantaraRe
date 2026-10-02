# 20: NB FacIn sebagai modul berjalan — backend (premi, akseptasi, tabel limit)

> ⚠️ **Disusun agent dari spec/XML atas perintah work owner — bukan hasil `/to-tickets`.** Perintah:
> `PROMPT-LANJUT-IMPLEMENT-NB-SAJA.md` (root repo, 01-10-2026) urutan 3; keputusan work owner 02-10-2026 (butir 58–60).

**What to build:** modul `nbfacin` terdaftar di aplikasi (`backend/modul.go`, `Pendaftaran()`) dan melayani dua
permintaan di atas mesin yang sudah diport: **hitung premi satu coverage** dan **satu langkah tangga akseptasi** yang
membaca **tabel limit dari Oracle** (`POOLDATA.M_LIMIT_*`), bukan dari CSV uji.

**Asal** `[terverifikasi]`:
- Skema tabel limit — **hanya** DDL `D:\migrasi\RNM\DDL\M_LIMIT_*.txt` (enam tabel). Kolom yang dibaca: bentuk A
  `JABATAN`, `TEAM_GROUP`, `LIMIT_BOTTOM`, `LIMIT_BOTTOM2`; bentuk B (`M_LIMIT_FINANCIALINS`) `JABATAN`,
  `LIMITBOND_BOTTOM`, `LIMITCREDITCL_BOTTOM`, `LIMITCREDITNCL_BOTTOM` — persis yang dibaca rule SQL Pega
  (`NB FacIn\RDBList\`, tiket 11/12). Semua `NUMBER(*,0)` (bilangan bulat) → dibaca teks, diurai desimal (tanpa `float`).
  *(Ralat 02-10-2026: semula tertulis "NUMBER tanpa presisi" — pola ekstraksi agent `\([0-9, ]*\)` membuang `(*,0)`.)*
  ⛔ `NAMA` dan `LOGIN` **tidak** dibaca.
- Pola modul: `docs/bersama/PANDUAN-TIM-PER-MODUL.md` bab 4, `modul/claimlife/backend/modul.go`.

**Keputusan work owner 02-10-2026:**
- **Butir 58** — endpoint akseptasi memakai **jabatan dari isian permintaan** (bukan dari sesi/IAM). ⚠️ **Tidak aman
  untuk produksi**: klien dapat mengaku jabatan mana pun. Ditandai di kode dan di respons; diganti begitu model peran
  (OQ RBAC) diputuskan.
- **Butir 59** — berkas slot menu `962_menu_nbfacin.sql` boleh ditulis (tidak dijalankan).
- **Butir 60** — layar pertama meniru section Pega (tiket 21).

**Di luar tiket ini:** pemuat data kasus NB dari Oracle (tiket 17: menunggu tabel flat; `FACINOFFER` adalah tabel
ringkasan keluaran, bukan sumber masukan premi); penulisan apa pun ke Oracle; migrasi tabel baru.

**Blocked by:** 11, 12, 18

**Status:** ready-for-human — diport 02-10-2026, code review dua sumbu ditangani; menunggu tinjauan work owner

- [x] `repository` membaca enam tabel limit lewat `db.Qualify` (ADR-U-0033), kolom persis di atas; teks SQL diuji
- [x] `services` menyusun `acceptance.TabelLimit` / daftar financial dari repository dan menjalankan satu langkah tangga
      (bentuk A, beralih ke B bila `ErrBentukB`) — tanpa salinan mesin
- [x] `POST /api/nbfacin/premi` — premi satu coverage + asal rumus; galat mesin dipetakan ke 422, bukan 500
- [x] `POST /api/nbfacin/akseptasi/langkah` — butir 58, bertanda tidak aman untuk produksi
- [x] `backend/modul.go` + daftar bangkitan; kontrak `PenilaiPredikatFacIn` dan `MesinPremiFacIn` **disediakan** lewat
      `Pendaftaran()` (menutup A32 untuk dua kontrak murni; tangga butuh tabel → tetap disambung pemakai)
- [x] Slot menu 962 + `MODUL.md` `Status`/prefix — dinyalakan bersama layar tiket 21 (uji menu dua arah)
- [x] Penjaga Go dan frontend hijau

## Comments

### 2026-10-02 — diport (agent)

**Kode:** `backend/modul.go` (`Pendaftaran()`: menyediakan `PenilaiPredikatFacIn`, `MesinPremiFacIn`),
`models/limit.go`, `repository/limit.go` (`LimitOracle.MuatLimit` — enam tabel lewat `db.Qualify`, kolom persis DDL,
`NUMBER` → teks → desimal), `services/layanan.go` (`HitungPremi`, `LangkahAkseptasi` lewat `kontrakfacin.Tangga`;
`ErrTidakDapatDiproses` → 422, `ErrTanpaDatabase` → 503), `handlers/handlers.go` (`POST /api/nbfacin/premi`,
`POST /api/nbfacin/akseptasi/langkah` — jawaban membawa peringatan butir 58; JSON medan tak dikenal → 400),
`migrations/962_menu_nbfacin.sql` + `_down.sql` (ditulis, **belum dijalankan**), daftar bangkitan
`inti/backend/daftar/modul_nbfacin_gen.go`, `MODUL.md` (`Status` dimigrasi, prefix, kontrak, bab "Tabel warisan"),
`docs/STRUKTUR-TABEL-NB-FACIN.md` (peta enam tabel warisan, tipe dari DDL). `premium.LiniDikenal` ditambah agar lini
dari luar ditolak sebelum `SatuanRate` panic.

**Uji:** repository (teks SQL tanpa `NAMA`/`LOGIN`/`*`, urai NUMBER dan NULL), services (tangga dari repository = hasil
tiket 11; galat dipilah), handlers (200/400/422/500/503). Uji DB nyata (`db`-tag, `ORACLE_DSN`) **tidak** ditulis —
membaca `POOLDATA` sungguhan butuh persetujuan; pemetaan diuji tanpa DB. Penjaga Go dan frontend hijau.

### 2026-10-02 — code review dua sumbu (agent)

**Spec** — dicocokkan reviewer ke XML: urutan/kelompok sel, label, properti, satu-satunya `pyRequired` (TSI), syarat sel
12/16 (literal tidak disalin), kolom SQL persis DDL, pemetaan 400/422/500/503, tanpa aritmetika uang di frontend.
Diperbaiki: tabel limit kosong kini **503 bersebab** (`ErrTabelLimitTakTersedia`), bukan 500 tanpa sebab; tipe kolom
**`NUMBER(*,0)`** (bilangan bulat), bukan "tanpa presisi" — pola ekstraksi agent membuang `(*,0)`; kontrol Name
(dropdown → kotak teks) dinyatakan sebagai penyimpangan; keterangan "belum diport" kini terlihat; uji paritas layar
`pages/CoverageCargo.test.ts` (urutan sel, empat kelompok, tepat tiga medan dapat diisi, tombol nonaktif) — uji mutasi
`../alat/mutasi_layar21.py` **5/5**.
**Standar** — tanpa pelanggaran keras di kode. Diperbaiki: `recover()` hanya memetakan **panic sikap predikat**
(`rules.PanikSikap`, tipe baru) ke 422, panic lain (bug program) tetap 500; galat 500 dicatat `log.Printf`; badan
permintaan dibatasi 64 KiB dan pesan galat urai tetap; `MODUL.md`/`modul.go` tidak lagi menyebut jalur "dirakit pemakai
lewat `kontrakfacin`" (terlarang §5); respons basi diabaikan sesudah isian berubah; uji keselarasan
`repository.TabelBentukA` ↔ `services.tabelDikenal`. **Dicatat, tidak diubah** (*smell*, keputusan menimbang):
`permintaanPremi` menyalin 21 medan kontrak; `premium.LiniDikenal` mengulang daftar `SatuanRate`; `Tabel`/`Jabatan`
`models` bertipe `string`; `Close()` manual di repository; *fixture* uji berulang; `toHaveLength(5)` di uji bersama
(milik tim inti, disunting minimal atas izin butir 62); jabatan dari isian (butir 58) ditinjau ulang saat RBAC diputuskan.

### 2026-10-02 — verifikasi independen (sesi `nusantarare-0f`) ditangani

- `HitungPremi` kini punya `recover`: panic dipilah lewat **tipe** — `premium.PanikLini` (lini di luar peta skala) dan
  `rules.PanikSikap` (predikat belum diport) → 422; panic lain (bug program) → 500 + `log.Printf`.
- Uji HTTP baru: 422 (panic sikap), 500 (panic bug), 503 (tabel limit kosong) — lewat `services.BaruDenganTangga`.
- Komentar galat di `services/layanan.go` dibetulkan (503, bukan 500).
- Penjaga identitas `rules/registry_test.go` tidak peka huruf besar-kecil, sama dengan generator; uji instrumen
  `TestPenjagaIdentitasMenggigit` (butir berjawaban diketahui), dan terbukti merah bila pencocokan dibuat peka huruf.
- `-migrate`: lihat butir 65 — hanya work owner, ke DEV/lokal; menjalankan migrasi tertunda semua modul.

