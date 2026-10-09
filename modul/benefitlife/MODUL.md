# Modul `benefitlife` — Benefit

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`PANDUAN-TIM-PER-MODUL.md` bab 5) — keputusan work owner 08-10-2026 (prompt
`PROMPT-BENEFITLIFE.md`, K1-K6 final). Panduan PATOKAN: section Pega `InboxBenefit` (kelas
`ASM-FW-GISFW-Int-BENEFIT_LIFE`, judul "INSURANCE BENEFIT") di `D:\NUSARE DEV\Menu Benefit\InboxBenefit.xml` (satu-satunya
XML; nomor `bNNN` = baris XML). Templat: modul saudara `ririsklife` (versi tanpa ringkasan / rincian / unggah).
Paritas setiap medan / tombol / aktivitas: `docs/PARITAS-LAYAR-DAN-AKSI.md`.

⛔ **Tanpa migrasi sendiri** (`—` di bawah): tabel = migrasi inti `942`-`944` (RENAME tabel Pega + satu tabel
`BENEFIT_LIFE`), baris menu = migrasi inti `945`.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `benefitlife` |
| Folder korpus | `Benefit Life` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-BENEFITLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `—` |
| Slot menu | `—` |
| Prefix rute API | `/api/benefit-life` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Keputusan work owner 08-10-2026 (final, dikutip dari prompt)

| # | Keputusan | Penerapan |
| --- | --- | --- |
| K1 | SATU tabel `BENEFIT_LIFE (ID VARCHAR2(10) PK, BENEFIT VARCHAR2(n))`: DROP VIEW (berpelindung, hanya bila VIEW), `ALTER TABLE M_BENEFIT_LIFE RENAME TO BENEFIT_LIFE` (berpelindung; PK `SYS_C009031` ikut), tambah `BENEFIT`, isi dari JSONDATA, `DROP COLUMN JSONDATA CASCADE CONSTRAINTS`; lebar = pola ririsklife (200) kecuali XML menetapkan lain | 942 (bentuk blok ALL_VIEWS + RENAME yang SUDAH diterima untuk 935/938 - **nol perubahan `perintah_katalog.go`**), 943 (`BENEFIT VARCHAR2(200)`: XML pxTextArea b1175 tanpa batas), 944 |
| K2 | `BENEFIT` dari kunci JSON `Benefit` (peka huruf); bukti SELECT baca-saja 10 nilai identik view; BERHENTI sebelum DROP JSONDATA bila BENEFIT NULL padahal JSON punya Benefit | 944 blok 1 = ekspresi notasi titik PERSIS teks view (`m.JSONDATA.Benefit` / view `a.JSONDATA.Benefit`); blok 2 = UPDATE yang sengaja gagal ORA-01407 (`SET m.ID = NULL` atas baris yang JSON-nya ber-Benefit tetapi BENEFIT NULL / berbeda) sebelum blok 3 (DROP). Bukti: `docs/sql/benefit_bukti_k2.sql` - **menunggu hasil WO** |
| K3 | Prosedur `PEGA_M_BENEFIT_LIFE` INVALID - diterima; sequence `M_BENEFIT_LIFE_SEQ` TETAP; NEXTVAL yang menghasilkan ID yang sudah ada = galat yang jelas | `models.BentukID` (`'1' \|\| LPAD(seq, 5, '0')`); `services.ErrIDTerpakai` → **409** berkalimat (uji `TestSimpanIDDariSequenceSudahAda`, `TestSimpanBalapanPK`, `TestRuteIDSequenceSudahAda`) |
| K4 | Menu `benefitlife`, LABEL `Benefit`, MASTER TREATY URUTAN 12, DIMIGRASI '1', STATUS_AKTIF seperti `ririsklife` | `945_m_nav_menu_benefitlife.sql`; penjaga `modulLuarKorpus` + `labelTampilDisetujui` (`Benefit Life` → kode, tampil `Benefit`); `MODUL_LUAR_KORPUS.benefitLife` |
| K5 | Hak menu = gabungan `riratelife`, `ricommlife`, `ririsklife` (PENUH bila salah satunya PENUH), aman diulang | `docs/sql/benefit_d_hak.sql` (LANGKAH-WO (d)) |
| K6 | Aturan simpan HANYA dari XML; tambahan pola MASTER TREATY ditulis di PARITAS | Benefit wajib (b1137) + huruf besar (`SetUpperCase_DT` b1211); di luar XML: pangkas spasi, batas 200 byte - PARITAS bab "Di luar XML". Nol penolakan kembar |

Prasyarat P1-P3 (penulis berhenti + kunci; acuan dua kali; cadangan ID + JSONDATA): `docs/LANGKAH-WO-BENEFITLIFE.md`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-BENEFITLIFE.md`, `PARITAS-LAYAR-DAN-AKSI.md`, `LANGKAH-WO-BENEFITLIFE.md` + `sql/benefit_*.sql`, `PR-BENEFITLIFE.md` |
| `backend/` | `models/` `repository/` (`bnfl.go`, `bnfl_tabel.go`) `services/` `handlers/` `tiruan/` `modul.go` (tanpa `migrations/`) |
| `frontend/` | `pages/BenefitLife.tsx` `labels.ts` `api.ts` `aturan.ts` `benefitlife.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Asumsi — bukti XML atau bawaan work owner

XML hanya rule **section**: isi `AddToList_Act`, `EditList_DT`, `NewData_DT`, `SetUpperCase_DT`, dan report definition
`BrowseBenefitLife_RD` TIDAK ada (hanya dirujuk). Ada bukti → **[terverifikasi]**; tidak ada → **[penyimpangan sadar -
menunggu WO]**.

| # | Aturan yang berlaku | Status dan bukti |
| --- | --- | --- |
| A1 | Number / ID tidak dapat diisi: Add = ID dari sequence (server), Edit = ID baris | **[terverifikasi]** pxTextInput `pyDisabled` true / `pyDisabledNew` always b1012-b1013; isian `id` di badan = 400 |
| A2 | `.Number` = ID, tanpa kolom NUMBER | **[keputusan work owner]** fakta WO: Number (5 baris) selalu = ID |
| A3 | Benefit huruf besar sebelum disimpan | **[terverifikasi sebagian]** `SetUpperCase_DT` b1211/b1327 pada perubahan textarea; isi DT tidak ada - arti namanya diterapkan (pertanyaan terbuka T1) |
| A4 | Grid 10 baris, halaman bernomor, ID menurun bawaan, ID dapat dibalik arahnya, Benefit tidak dapat diurutkan, saring ID / Benefit | **[terverifikasi]** b4526, b4498, b4391/b4397, b4493, b4415, b4505; saring "memuat" tanpa beda huruf = pola ririsklife |
| A5 | Save = Add bila form tanpa ID, Edit bila form diisi Edit | **[penyimpangan sadar - menunggu WO]** satu tombol Save b1772 → `AddToList_Act` b1791 untuk keduanya; isi aktivitas tidak ada |
| A6 | Tanpa Delete, Upload CSV, ringkasan/detail | **[terverifikasi]** `pyGridDeleteActivityExists` false b3292; tidak ada tombol lain di XML |

## Migrasi

Nol migrasi modul. Migrasi inti `942_benefit_life_ganti_nama.sql`, `943_benefit_life_kolom.sql`,
`944_benefit_life_satu_tabel.sql`, `945_m_nav_menu_benefitlife.sql` (masing-masing + `_down`). Karena 942-944 mengubah
bentuk tabel Pega `M_BENEFIT_LIFE`, `BENEFIT_LIFE` TIDAK dinyatakan "Tabel warisan" (preseden ririsklife 935-940); kolom
yang dibuat migrasi tercatat di `docs/STRUKTUR-TABEL-BENEFITLIFE.md`.

## Butir terbuka

- **Bukti K2** (`docs/sql/benefit_bukti_k2.sql`) dan semua angka acuan menunggu hasil WO - executor tidak terhubung ke
  Oracle (koneksi ditolak pengaman sesi; tidak dicoba ulang).
- **Kunci Number, pxObjClass, pyRuleHarness, px\*** tidak dipindah (bukan kolom view) - hanya di cadangan P3.

## Menjalankan uji modul ini saja

Dari folder akar repo:

```powershell
go test ./modul/benefitlife/...
go vet -tags db ./modul/benefitlife/...   # uji seam Oracle (skema uji) - butuh ORACLE_DSN skema uji
npx vitest run modul/benefitlife
```
