# Modul `planlife` — Plan

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`PANDUAN-TIM-PER-MODUL.md` bab 5) — keputusan work owner 08-10-2026 (prompt
`PROMPT-PLANLIFE.md`, K1-K7 final). Panduan PATOKAN: section Pega `InboxProductType` (kelas
`ASM-FW-GISFW-Int-PRODUCT_TYPE_LIFE`, judul "Plan" b331 / "Product Type" b644) di
`D:\NUSARE DEV\Menu Plan\InboxProductType.xml` (satu-satunya XML, 7.676 baris; `bNNNN` = baris XML). Templat: modul
`benefitlife` (satu tabel, satu layar) + dua autocomplete. Paritas: `docs/PARITAS-LAYAR-DAN-AKSI.md`.

⛔ **Tanpa migrasi sendiri** (`—` di bawah): tabel = migrasi inti `946`-`948` (RENAME tabel Pega + satu tabel
`PRODUCT_TYPE_LIFE` ber-PK), baris menu = migrasi inti `949`.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `planlife` |
| Folder korpus | `Plan Life` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-PLANLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `—` |
| Slot menu | `—` |
| Prefix rute API | `/api/plan-life` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Keputusan work owner 08-10-2026 (final, dikutip dari prompt)

| # | Keputusan | Penerapan |
| --- | --- | --- |
| K1 | SATU tabel `PRODUCT_TYPE_LIFE (ID, COVERNAME, BUSINESS, BUSINESSID, BENEFIT, BENEFITID)`: DROP VIEW (hanya bila VIEW), RENAME `M_PRODUCT_TYPE_LIFE` (berpelindung), tambah lima kolom, isi dari JSONDATA (notasi titik persis view), ID NOT NULL + `PK_PRODUCT_TYPE_LIFE`, BERHENTI sebelum DROP JSONDATA bila ID NULL / ganda / ≠ JSONDATA.ID, DROP JSONDATA; nama 200, ID master 10, ID tetap VARCHAR2(6); tanpa FK | 946 (bentuk blok yang sudah ada - **nol perubahan `perintah_katalog.go`**), 947, 948 (5 blok: isi → periksa ID → periksa isi → PK → buang). Pemaksa berhenti = `SET m.COVERNAME = RPAD(CHR(88), 201, CHR(88))` → ORA-12899 (tanpa kutip; ID masih NULLABLE sehingga bentuk `ID = NULL` 944 tidak berlaku). PK sekaligus NOT NULL |
| K2 | Bukti sebelum migrate (`docs/sql/plan_bukti_k1.sql`, baca-saja) | **menunggu hasil WO**; uji Go `TestMigrasi948SatuTabel` membuktikan teks ekspresi = teks view |
| K3 | `PEGA_M_PRODUCT_TYPE_LIFE` INVALID diterima; sequence TETAP; ID dari NEXTVAL yang sudah ada = galat jelas | `models.BentukID`; `services.ErrIDTerpakai` → **409** (uji `TestSimpanIDDariSequenceSudahAda`, `TestRuteIDSequenceSudahAda`) |
| K4 | Business dari `BUSINESS` `GROUPPANEL = '009'` (diturunkan dari data), tampil OLDID + Note, simpan Note + ID; Benefit dari `BENEFIT_LIFE`; teks dicocokkan ulang server tanpa beda huruf, tidak cocok = 422 | `GET /pilihan/business`, `GET /pilihan/benefit`; `services.cocokBusiness` / `cocokBenefit` (ganda juga 422) |
| K5 | Plan Name wajib + unik (di luar XML); Business / Benefit wajib kecuali terbukti boleh kosong | ketiganya wajib - XML / aktivitas tidak membuktikan boleh kosong (PARITAS) |
| K6 | Menu `planlife`, LABEL `Plan`, MASTER TREATY URUTAN 13, DIMIGRASI '1' | `949_m_nav_menu_planlife.sql`; penjaga `modulLuarKorpus` + `labelTampilDisetujui`; `MODUL_LUAR_KORPUS.planLife` |
| K7 | Hak = gabungan `riratelife`, `ricommlife`, `ririsklife`, `benefitlife` | `docs/sql/plan_d_hak.sql` |

## RALAT

| # | Bunyi lama | Bunyi baru | Bukti |
| --- | --- | --- | --- |
| R1 (08-10-2026) | Petunjuk prompt: form memuat `.ID` (tampil) | Form `InboxProductType` TIDAK memuat medan `.ID`: medannya hanya Plan Name b812, Business b1079, Benefit b1432; `.ID` hanya parameter `EditProductTypeLife_Act` b4060. ID tidak ditampilkan (grid pun tanpa kolom ID); mode Edit ditandai teks "Editing …" | grep `.ID` di XML: b1265 (kolom RD business), b4060 / b4161 (parameter Edit) |
| R2 (08-10-2026) | Petunjuk prompt: tombol New = aksi umum | New b6667 (`NewProductTypeLife_act` b6686) tampil HANYA saat `InputParam.DATASHOW = 'IsEdit'` b6827 - berperan sebagai batal-edit; Save b6403 selalu tampil | b6539-b6831 |

Prasyarat P1-P3 + bukti K2: `docs/LANGKAH-WO-PLANLIFE.md`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-PLANLIFE.md`, `PARITAS-LAYAR-DAN-AKSI.md`, `LANGKAH-WO-PLANLIFE.md` + `sql/plan_*.sql`, `PR-PLANLIFE.md` |
| `backend/` | `models/` `repository/` (`plnl.go`, `plnl_tabel.go`) `services/` `handlers/` `tiruan/` `modul.go` (tanpa `migrations/`) |
| `frontend/` | `pages/PlanLife.tsx` `labels.ts` `api.ts` `aturan.ts` `planlife.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Asumsi — bukti XML atau bawaan work owner

XML hanya rule **section**: isi `SaveProductTypeLife_Act`, `EditProductTypeLife_Act`, `NewProductTypeLife_act`, dan RD
`BrowseProductTypeLife_RD` / `BrowseBusinessLife_RD` / `BrowseBenefitLife_RD` TIDAK ada (hanya dirujuk).

| # | Aturan yang berlaku | Status dan bukti |
| --- | --- | --- |
| A1 | Save = Add bila form tanpa ID, Edit sesudah Edit baris | **[penyimpangan sadar - menunggu WO]** satu Save b6422; isi aktivitas tidak ada (pola benefitlife T2) |
| A2 | Grid tanpa saring; urut Plan Name / Benefit, Business tidak; bawaan ID menaik | saring/urut **[terverifikasi]** b4380, b4247 / b4269 / b4291; bawaan **[penyimpangan sadar]** (`pySortType` NONE, `pyMaxSortOrder` 0 b2890 - tanpa urutan; ID menaik dipilih agar stabil) |
| A3 | Plan Name tidak diubah hurufnya | **[terverifikasi]** refresh b864-b875 tanpa data transform |
| A4 | Tanpa Delete, Upload, ringkasan/detail, saring | **[terverifikasi]** `pyGridDeleteActivityExists` false b419 / b2839; `pyGridFiltering` false b4380 |

## Migrasi

Nol migrasi modul. Migrasi inti `946_product_type_life_ganti_nama.sql`, `947_product_type_life_kolom.sql`,
`948_product_type_life_satu_tabel.sql`, `949_m_nav_menu_planlife.sql` (masing-masing + `_down`). Karena 946-948 mengubah
bentuk tabel Pega `M_PRODUCT_TYPE_LIFE`, `PRODUCT_TYPE_LIFE` TIDAK dinyatakan "Tabel warisan" (preseden benefitlife);
kolom yang dibuat migrasi tercatat di `docs/STRUKTUR-TABEL-PLANLIFE.md`.

## Butir terbuka

- **Bukti K2** (`docs/sql/plan_bukti_k1.sql`), acuan, dan lebar `BUSINESS.ID` (kueri E) menunggu hasil WO - executor
  tidak terhubung ke Oracle (pembacaan DB ditolak pengaman sesi; tidak dicoba).
- Kunci `px*` / `pyRuleHarness` tidak dipindah - hanya di cadangan P3.

## Menjalankan uji modul ini saja

```powershell
go test ./modul/planlife/...
go vet -tags db ./modul/planlife/...   # uji seam Oracle (skema uji) - butuh ORACLE_DSN skema uji
npx vitest run modul/planlife
```
