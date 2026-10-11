# Modul `diseaselife` — Disease Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`PANDUAN-TIM-PER-MODUL.md` bab 5) — keputusan work owner 08-10-2026 (prompt
`PROMPT-DISEASELIFE-COVERLIFE.md`, K0, D1-D4, K5 final). Panduan PATOKAN: section Pega `InboxDisease` (kelas
`ASM-FW-GISFW-Int-DISEASE_LIFE`, judul "DISEASE" b371) di `D:\NUSARE DEV\Menu Disease\InboxDisease.xml` (satu-satunya
XML, 6.451 baris, dibaca utuh; `bNNNN` = baris XML). Templat: modul `causeoflosslife` (modul berjatah, migrasi di folder
modul, slot menu) dan `benefitlife` (layar Edit / Save / Cancel, saring + urut grid). Paritas:
`docs/PARITAS-LAYAR-DAN-AKSI.md`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Migrasi modul ini tinggal
di `backend/migrations/` (rentang inti 900-949 sudah penuh). Jatahnya **dipinjam**: rentang dari premiumlistlife
(`050-089` → `050-079`), slot dari claimlife (`950-951` → `950-950`) — keputusan work owner 08-10-2026 K0, **perlu
persetujuan tim inti (CODEOWNERS)**; dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `diseaselife` |
| Folder korpus | `Disease Life` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-DISEASELIFE` |
| Status | dimigrasi |
| Rentang migrasi | `080-084` |
| Slot menu | `951-951` |
| Prefix rute API | `/api/disease-life` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Keputusan work owner 08-10-2026 (final, dikutip dari prompt)

| # | Keputusan | Penerapan |
| --- | --- | --- |
| K0 | Nomor: premiumlistlife `050-079`, modul ini `080-084`; claimlife slot `950`, modul ini slot `951`; periksa dulu nomor itu kosong | Diperiksa sebelum menulis (Glob + isi MODUL.md): nol berkas `06[5-9]_*` / `07?_*` / `08?_*` / `950_*` / `951_*` / `956_*` / `957_*` di repo; nol MODUL.md lain yang menyatakan `080-089` / `951` / `957`. Tabel jatah diubah di `premiumlistlife`, `claimlife`, `treatycontractout` + dua modul baru. Penjaga `rentang_test.go` / `menu_test.go` |
| D1 | `DISEASE_LIFE` TIDAK di-RENAME, kolom TIDAK diubah. Migrasi hanya: 1. `SEQ_DISEASE_LIFE` (awal = ID angka tertinggi + 1, NOCACHE NOCYCLE, pola PERSIS 923), ID baru = `TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL)`; 2. `PK_DISEASE_LIFE` (berhenti dengan galat jelas bila ID NULL / kembar, tanpa menghapus data); 3. baris menu 951 | 080 (blok sequence-dari-kueri, `migrasi.BacaSequenceDariKueri`), 081 (blok berpelindung ALL_CONSTRAINTS; galat ALTER sendiri = penghenti: ORA-02437 kembar / ORA-01449 NULL), 951. `models.BentukID`, `services.ErrIDTerpakai` → **409** bila NEXTVAL menabrak ID (Pega masih dapat menulis) |
| D2 | Baris uji 102051 `TEST123 / Sakit` dihapus WO SEBELUM `-migrate` lewat `docs/sql/disease_hapus_baris_uji.sql` (cadangan dulu, DELETE persis satu baris, rollback bila bukan 1, COMMIT). `M_DISEASE_LIFE` + `PEGA_M_DISEASE_LIFE` TIDAK disentuh | Berkas WO ditulis (satu-satunya SQL tulis data modul ini); **menunggu WO**. Pertanyaan terbuka di `docs/PR-DISEASELIFE.md` |
| D3 | ICD Code dan Disease wajib; huruf besar hanya pada medan yang memanggil `SetUpperCase_DT`; ICD Code unik tanpa beda huruf (di luar XML); grid paging + saringan di server, urut ID | `SetUpperCase_DT` dipanggil KEDUA medan (ICD Code b1265 / b1386, Disease b1538 / b1654) → keduanya huruf besar. **Temuan**: XML menandai ICD Code TIDAK wajib (b1195 / b1237 / b1244) - wajib diterapkan sebagai keputusan WO (PARITAS). Grid: saring ICD Code / Disease + OFFSET / FETCH di server, bawaan ID menurun (b5012) |
| D4 | Menu `diseaselife`, LABEL `Disease Life`, MASTER TREATY URUTAN 15, DIMIGRASI '1', di slot 951, aman diulang | `951_menu_diseaselife.sql` (INSERT datar `NOT EXISTS`); penjaga `modulLuarKorpus` + `barisLahirDiSlot`; `MODUL_LUAR_KORPUS.diseaseLife` |
| K5 | Hak menu = salinan pemegang `causeoflosslife` beserta HAK-nya, aman diulang | `docs/sql/dis_d_hak.sql` (LANGKAH-WO (d)) |

Prasyarat P1-P3 (cadangan PALING ATAS; penulis berhenti + kunci; acuan dua kali), D2, bukti:
`docs/LANGKAH-WO-DISEASELIFE.md`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-DISEASELIFE.md`, `PARITAS-LAYAR-DAN-AKSI.md`, `LANGKAH-WO-DISEASELIFE.md` + `sql/dis_*.sql` + `sql/disease_hapus_baris_uji.sql`, `PR-DISEASELIFE.md` |
| `backend/` | `models/` `repository/` (`dsl.go`, `dsl_tabel.go`) `services/` `handlers/` `tiruan/` `migrations/` `modul.go` |
| `frontend/` | `pages/DiseaseLife.tsx` `labels.ts` `api.ts` `aturan.ts` `diseaselife.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Asumsi — bukti XML atau bawaan work owner

XML hanya rule **section**: isi `AddToList_Act`, `EditList_DT`, `NewData_DT`, `SetUpperCase_DT`, dan report definition
`BrowseDiseaseLife_RD` TIDAK ada (hanya dirujuk).

| # | Aturan yang berlaku | Status dan bukti |
| --- | --- | --- |
| A1 | ID tidak dapat diisi: Add = ID dari SEQ_DISEASE_LIFE (server), Edit = ID baris | **[terverifikasi]** pxTextInput `pyDisabled` true b1070, `pyDisabledNew` always b1069; isian `id` di badan = 400 |
| A2 | Save = Add bila form tanpa ID, Edit sesudah Edit baris | **[penyimpangan sadar - menunggu WO]** satu Save b2102 → `AddToList_Act` b2121; isi aktivitas tidak ada (pola benefitlife T2 / causeoflosslife T1) |
| A3 | Grid: bawaan ID menurun, urut pilihan ID / ICD Code, saring, 10 baris, halaman bernomor | **[terverifikasi]** `pySortType` DESC b5012 / `pySortOrder` 1 b5017; `pyColumnSorting` true b5014 (ID) / b5036 (ICD Code), false b5058 (Disease) / b5080 (Edit); `pyGridFiltering` true b5148; `pyPageSize` 10 b5169; Numeric b5140 |
| A4 | ICD Code dan Disease huruf besar | **[terverifikasi pemanggilnya, isi tidak ada]** `SetUpperCase_DT` pada perubahan kedua medan (b1265, b1386, b1538, b1654); arti nama diterapkan (preseden benefitlife T1) |
| A5 | Tanpa Delete, Upload, ringkasan / detail | **[terverifikasi]** `pyGridDeleteActivityExists` false b460 / b754 / b1929 / b3262 / b3619 / b5297; tidak ada tombol lain |

## Migrasi

Rentang `080-084` (terpakai `080-081`), slot menu `951-951` (terpakai `951`). Nama berkas yang sudah dijalankan tidak
pernah diubah (`T_MIGRASI` mencatat nama).

| No | Berkas | Isi |
| --- | --- | --- |
| 080 | `080_seq_disease_life.sql` | `SEQ_DISEASE_LIFE` START WITH ID angka tertinggi + 1 (dihitung di basis data), NOCACHE NOCYCLE, berpelindung ALL_SEQUENCES |
| 081 | `081_disease_life_pk.sql` | `PK_DISEASE_LIFE (ID)` berpelindung ALL_CONSTRAINTS; berhenti ORA-02437 / ORA-01449 bila ID kembar / NULL |
| 951 | `951_menu_diseaselife.sql` | INSERT baris `M_NAV_MENU` (D4), `NOT EXISTS` |

Urutan pelari di skema baru: 080-081 → 900 … 949 → 951 (uji `TestUrutanPelari080Lalu900Lalu951`; versi Oracle
`TestDBUrutanPelari080Lalu900Lalu951`, `-tags=db`, belum dijalankan). Skema uji punya tiruan `DISEASE_LIFE` (tanpa PK)
SEBELUM pelari (`uji/skemauji/disease_tiruan.go`). Karena 080-081 mengubah objek Pega, `DISEASE_LIFE` TIDAK dinyatakan
"Tabel warisan" (preseden causeoflosslife); bentuknya tercatat di `docs/STRUKTUR-TABEL-DISEASELIFE.md`.

## Pembaca lain

`claimlife` (pencarian diagnosa `repository/penyakit.go`, `services/penyakit.go`, `CariDiagnosa.tsx`) membaca
`DISEASE_LIFE` (`ID`, `ICD_CODE`, `DISEASE`, berbatas `pyMaxRecords` 500): nama dan kolom TETAP, kueri TIDAK berubah;
PK hanya menambah indeks unik pada ID. Tabel klaim menyimpan TEKS ICD / Disease, bukan ID - menghapus baris uji 102051
tidak mematahkan data klaim (fakta WO). `claimlife` hanya mendapat perubahan jatah K0 di `MODUL.md`.

## Butir terbuka

- **Bukti dan acuan** (`docs/sql/dis_bukti.sql`, `dis_p2_acuan.sql`) dan langkah D2 menunggu hasil WO - executor tidak
  terhubung ke Oracle (pembacaan DB ditolak pengaman sesi; tidak dicoba).
- Nasib `M_DISEASE_LIFE` (3 baris JSON lama) + `PEGA_M_DISEASE_LIFE` + `M_DISEASE_LIFE_SEQ`, dan prosedur
  `PEGA_DISEASE_LIFE` yang rumus ID-nya bertabrakan - `docs/PR-DISEASELIFE.md` bab "Pertanyaan terbuka".

## Menjalankan uji modul ini saja

```powershell
go test ./modul/diseaselife/...
go vet -tags db ./modul/diseaselife/...   # uji seam Oracle (skema uji) - butuh ORACLE_DSN skema uji
npx vitest run modul/diseaselife
```
