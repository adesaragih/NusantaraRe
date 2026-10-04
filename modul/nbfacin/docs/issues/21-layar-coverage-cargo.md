# 21: Layar pertama NB — blok coverage kargo (port section) + isian uji premi

> ⚠️ **Disusun agent dari XML atas perintah work owner — bukan hasil `/to-tickets`.** Keputusan work owner 02-10-2026:
> butir 60 (layar = tiru section Pega) dan butir 61 ("Section + isian uji premi").

**What to build:** halaman awal modul `nbfacin` yang meniru **blok pertama** section
`InputCoverageCargo_FacIn` (MARINE CARGO, satu coverage) — label, urutan, kelompok, dan kontrol apa adanya — dengan
**penyimpangan sadar** (butir 61): medan masukan rumus premi dapat diisi, dan satu tombol tambahan "Hitung premi (alat
periksa)" memanggil `POST /api/nbfacin/premi` (tiket 20) lalu menampilkan Premi.

**Asal** `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\Section\InputCoverageCargo_FacIn.xml`
(ASM-FW-GISFW-DATA-COVERAGE!INPUTCOVERAGECARGO_FACIN), kelas `ASM-FW-GISFW-Data-Coverage`; rantai sisipan
`InputDtlCoverage_FacIn` → `InputCoverageCargo_FacIn`. Sel blok pertama (properti `pyValue`, label `pyLabelFor`, kontrol
`pyFormat`; semua **medan** `pyEditOptions = Read-only` di Pega, kedua tombol `Auto`):

| Layout | Sel | Label | Properti | Kontrol | Di port |
| --- | ---: | --- | --- | --- | --- |
| 1 | 3 | Keterangan Jaminan | `.CoverageNote` | TextArea | tampil (read-only) |
| 1 | 4 | Coverage Initial | `.CoverageInitial` | TextInput | tampil (read-only) |
| 1 | 5 | Select Coverage *(tombol)* | — | Button → `SetIndexMarineCargo_Act`, popup harness | nonaktif, "belum diport" |
| 6/8 | 10 | Name | `.Currency.Name` | Dropdown | **dapat diisi** sebagai kotak teks (butir 61; sumber opsi dropdown tidak ada di blok ini — penyimpangan kontrol) |
| 6/8 | 11 | Rate (%) | `.Rate` | Number | **dapat diisi** (butir 61) |
| 6/8 | 12 | LimitofLiability | `.LimitofLiability` | Number | tampil (read-only) |
| 13 | 15 | .CurrencyMaster | `.CurrencyMaster` | TextInput | tampil (read-only), label apa adanya |
| 13 | 16 | LimitofLiability | `.LimitofLiability` | Number | tidak tampil — lihat syarat |
| 13 | 17 | Get Currency Master *(tombol)* | — | Button → refresh `GetCurrMasterCargo` | nonaktif, "belum diport" |
| 13 | 18 | Diskon (%) | `.DiscountPercentage` | Number | tampil (read-only) |
| 19 | 21 | TSI *(wajib)* | `.TSI` | Number | **dapat diisi** (butir 61) |
| 19 | 22 | Premi | `.Premium` | Number | hasil hitung (read-only) |
| 19 | 23 | Min Premi | `.MinPremium` | Number | tampil (read-only) |
| 19 | 24 | Diskon | `.Discount` | Number | tampil (read-only) |

**Syarat tampil** `[terverifikasi]` (`pyUserData/pyCondition`): sel 12 tampil bila
`pyWorkPage.OfferFacIn.PolicyMasterNumber != <nomor polis master literal>`, sel 16 bila sama — **nilai literal tidak
disalin** (CLAUDE.md §4 butir 10, §10). Port: hanya cabang umum (sel 12); cabang khusus satu polis master tidak diport
(`[pertanyaan terbuka]` pemilik Product+Underwriting: aturan bisnis di balik satu nomor polis itu).

**Di luar tiket ini:** blok lain section (deductible sel 94–106, klausul 125–149), section pembungkus
(`InputDtlCoverage_FacIn`: grid coverage dan total per mata uang), varian `_IsUW` dan syarat read-only `IsUW`/`IsRenewal`,
master policy MARINE (`PolicyType`/`IsMOP` tidak ada di blok ini → selalu bukan master policy), sumber data kasus (tiket
17). Medan tanpa sumber tampil kosong — tidak diisi nilai rekaan.

**Blocked by:** 20

**Status:** ready-for-human — diport 02-10-2026, code review dua sumbu ditangani; menunggu tinjauan work owner

- [x] Halaman awal modul: blok di atas, label verbatim, urutan dan kelompok sesuai layout Pega
- [x] Name, Rate (%), TSI dapat diisi; "Hitung premi (alat periksa)" memanggil `POST /api/nbfacin/premi` (lini MARINE CARGO)
- [x] Premi tampil sebagai teks desimal dari backend (tanpa aritmetika uang di frontend); galat 422 tampil
- [x] Tombol Pega yang belum diport tampil nonaktif dengan keterangan
- [x] Penyimpangan butir 61 tertulis di layar dan di kode
- [x] `menu.ts`, `rute.tsx`, `labels.ts`, `api.ts` sesuai pola modul; uji frontend dan penjaga lapisan hijau

## Comments

### 2026-10-02 — diport (agent)

**Kode:** `frontend/pages/CoverageCargo.tsx` (halaman awal `nbfacin-coverage-cargo`), `labels.ts` (verbatim korpus) +
`labels.test.ts` (membuka section korpus; pemindai kedalaman `rowdata`; uji gigit label salah), `api.ts`
(`badanPremiCargo`: angka tetap teks) + `api.test.ts`, `menu.ts` (`kelompok` `NB FacIn` = `M_NAV_MENU.LABEL`),
`rute.tsx`. Keterangan Jaminan memakai `textarea` read-only lokal: `Area` inti tanpa mode read-only, dan komponen
bersama di luar jatah.

**Uji bersama disunting (butir 62):** `frontend/daftar.datar.test.ts`, `daftar.sinkron.test.ts` (4 → 5 modul aktif),
`daftar.menuTabel.test.ts` (contoh belum dimigrasi `nbfacin` → `nbtreatyin`). Suite frontend: semua hijau kecuali 8
kegagalan `treatycontractout` yang sudah ada di kode ter-commit sebelum pekerjaan ini (tidak disentuh).

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
