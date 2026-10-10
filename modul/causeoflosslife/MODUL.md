# Modul `causeoflosslife` — Cause Of Loss Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`PANDUAN-TIM-PER-MODUL.md` bab 5) — keputusan work owner 08-10-2026 (prompt
`PROMPT-CAUSEOFLOSSLIFE.md`, K0-K6 final). Panduan PATOKAN: section Pega `InboxCauseofLossLife` (kelas
`ASM-FW-GISFW-Int-CAUSEOFLOSS_LIFE`, judul "CAUSE OF LOSS" b343) di
`D:\NUSARE DEV\Menu Cause Of Loss\InboxCauseofLossLife.xml` (satu-satunya XML, 6.568 baris, dibaca utuh; `bNNNN` =
baris XML). Templat: modul `benefitlife` (satu tabel `ID` + satu nama, Edit / Save / Cancel). Paritas:
`docs/PARITAS-LAYAR-DAN-AKSI.md`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Berbeda dengan
benefitlife / planlife, migrasi modul ini tinggal di `backend/migrations/` (rentang inti 900-949 sudah penuh). Jatahnya
**dipinjam dari premiumlistlife** (K0: `050-099` → `050-089`, slot `954-955` → `954-954`) — keputusan work owner
08-10-2026, **perlu persetujuan tim inti (CODEOWNERS)**; dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `causeoflosslife` |
| Folder korpus | `Cause Of Loss Life` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-CAUSEOFLOSSLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `090-099` |
| Slot menu | `955-955` |
| Prefix rute API | `/api/cause-of-loss-life` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Keputusan work owner 08-10-2026 (final, dikutip dari prompt)

| # | Keputusan | Penerapan |
| --- | --- | --- |
| K0 | Nomor migrasi: premiumlistlife dikecilkan ke `050-089` / slot `954`; modul ini `090-099` / slot `955`; periksa dulu nomor itu kosong | Diperiksa sebelum menulis: nol berkas `06[5-9]_*`, `0[7-9]?_*`, `954_*`, `955_*`; nol `MODUL.md` lain yang menyatakan 090-099 / 955. Dua tabel jatah diubah (`modul/premiumlistlife/MODUL.md` + tabel di atas). Penjaga `rentang_test.go` / `menu_test.go` hijau untuk jatah baru |
| K1 | SATU tabel `CAUSEOFLOSS_LIFE (ID VARCHAR2(10) PK, CAUSEOFLOSS VARCHAR2(200) NULLABLE)`: DROP VIEW (berpelindung), RENAME `M_CAUSEOFLOSS_LIFE` (berpelindung, PK `SYS_C008825` ikut), ADD kolom, isi notasi titik PERSIS view, pemeriksaan (NULL = NULL sama), DROP JSONDATA | 090 (blok ALL_VIEWS + RENAME - bentuk yang sudah ada, **nol perubahan `perintah_katalog.go`**), 091, 092 (isi → periksa → buang). Pemaksa berhenti = `SET m.ID = NULL` (pola 944; ID NOT NULL) → ORA-01407, tanpa kutip |
| K2 | Bukti sebelum migrate (`docs/sql/col_bukti_k1.sql`, baca-saja): 4 baris, ekspresi = view per baris (100001 NULL di keduanya), nol nilai > 200 byte | **menunggu hasil WO**; uji Go `TestMigrasi092SatuTabel` / `TestPemeriksaanK14NullSamaDenganNull` membuktikan teks ekspresi = teks view dan logika NULL |
| K3 | `PEGA_M_CAUSEOFLOSS_LIFE` INVALID diterima; sequence TETAP (last 5, ID berikutnya 100005); NEXTVAL yang menabrak ID = galat jelas | `models.BentukID`; `services.ErrIDTerpakai` → **409** (uji `TestSimpanIDDariSequenceSudahAda`, `TestRuteIDSequenceSudahAda`) |
| K4 | Cause of Loss wajib (XML); huruf besar HANYA bila XML punya pengubahnya; tolak kembar + pangkas = di luar XML (PARITAS) | wajib b996 / b1044; XML TANPA pengubah huruf → huruf tidak diubah; pangkas + kembar tanpa beda huruf → 422 (PARITAS "Di luar XML") |
| K5 | Menu `causeoflosslife`, LABEL `Cause Of Loss Life`, MASTER TREATY URUTAN 14, DIMIGRASI '1', di slot 955, aman diulang | `955_menu_causeoflosslife.sql` (INSERT datar `NOT EXISTS`); penjaga `modulLuarKorpus` + `barisLahirDiSlot` (baris luar korpus PERTAMA yang lahir di slot modul); `MODUL_LUAR_KORPUS.causeOfLossLife` |
| K6 | Hak menu = salinan pemegang `planlife` beserta HAK-nya, aman diulang | `docs/sql/col_d_hak.sql` (LANGKAH-WO (d)) |

Prasyarat P1-P3 (penulis berhenti + kunci; acuan dua kali; cadangan ID + JSONDATA) + bukti K2:
`docs/LANGKAH-WO-CAUSEOFLOSSLIFE.md`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-CAUSEOFLOSSLIFE.md`, `PARITAS-LAYAR-DAN-AKSI.md`, `LANGKAH-WO-CAUSEOFLOSSLIFE.md` + `sql/col_*.sql`, `PR-CAUSEOFLOSSLIFE.md` |
| `backend/` | `models/` `repository/` (`coll.go`, `coll_tabel.go`) `services/` `handlers/` `tiruan/` `migrations/` `modul.go` |
| `frontend/` | `pages/CauseOfLossLife.tsx` `labels.ts` `api.ts` `aturan.ts` `causeoflosslife.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Asumsi — bukti XML atau bawaan work owner

XML hanya rule **section**: isi `AddToList_Act`, `EditList_DT`, `NewData_DT`, dan report definition
`BrowseCauseofLossLife_RD` TIDAK ada (hanya dirujuk).

| # | Aturan yang berlaku | Status dan bukti |
| --- | --- | --- |
| A1 | ID tidak dapat diisi: Add = ID dari sequence (server), Edit = ID baris | **[terverifikasi]** pxTextInput `pyDisabled` true b872, `pyDisabledNew` always b871; isian `id` di badan = 400 |
| A2 | Save = Add bila form tanpa ID, Edit sesudah Edit baris | **[penyimpangan sadar - menunggu WO]** satu Save b5366 → `AddToList_Act` b5385; isi aktivitas tidak ada (pola benefitlife T2) |
| A3 | Grid ID menaik tetap, tanpa saring, tanpa urut pilihan, 10 baris, halaman bernomor | **[terverifikasi]** `pySortType` ASC b3923 / `pySortOrder` 1 b3929; `pyColumnSorting` false b3925 / b3947 / b3969; `pyGridFiltering` false b4036; `pyPageSize` 10 b4057; Numeric b4028 |
| A4 | Huruf Cause of Loss tidak diubah | **[terverifikasi]** nol data transform / refresh pada medan b967-b1133 |
| A5 | Tanpa Delete, Upload, ringkasan/detail | **[terverifikasi]** `pyGridDeleteActivityExists` false b431 / b725 / b1605 / b1890 / b2847; tidak ada tombol lain |

## Migrasi

Rentang `090-099` (terpakai `090-092`), slot menu `955-955` (terpakai `955`). Nama berkas yang sudah dijalankan tidak
pernah diubah (`T_MIGRASI` mencatat nama).

| No | Berkas | Isi |
| --- | --- | --- |
| 090 | `090_causeofloss_life_ganti_nama.sql` | DROP VIEW (berpelindung ALL_VIEWS) + RENAME `M_CAUSEOFLOSS_LIFE` → `CAUSEOFLOSS_LIFE` |
| 091 | `091_causeofloss_life_kolom.sql` | ALTER ADD `CAUSEOFLOSS VARCHAR2(200)`, berdiri sendiri |
| 092 | `092_causeofloss_life_satu_tabel.sql` | isi dari JSONDATA → pemeriksaan K1.4 → DROP JSONDATA (aman diulang) |
| 955 | `955_menu_causeoflosslife.sql` | INSERT baris `M_NAV_MENU` (K5), `NOT EXISTS` |

Urutan pelari di skema baru: 090-092 → 900 … 949 → 955 (uji `TestUrutanPelari090Lalu900Lalu955`; versi Oracle
`TestDBUrutanPelari090Lalu900Lalu955`, `-tags=db`, belum dijalankan). Skema uji punya tiruan `M_CAUSEOFLOSS_LIFE` +
view + sequence SEBELUM pelari (`uji/skemauji/coll_tiruan.go`). Karena 090-092 mengubah bentuk tabel Pega,
`CAUSEOFLOSS_LIFE` TIDAK dinyatakan "Tabel warisan" (preseden benefitlife); kolom yang dibuat migrasi tercatat di
`docs/STRUKTUR-TABEL-CAUSEOFLOSSLIFE.md`.

## Pembaca lain

`masterproductnamelife` membaca `CAUSEOFLOSS_LIFE` (pemilih Cause Of Loss, `SELECT ID, CAUSEOFLOSS ... ORDER BY ID ASC`):
nama dan kolom TETAP, kueri tidak berubah. Diperbarui hanya `testdata/katalog-dev.json` (VIEW → TABLE), komentar, dan
DDL tiruan (lebar 10 / 200).

## Butir terbuka

- **Bukti K2** (`docs/sql/col_bukti_k1.sql`), acuan P2, dan keadaan menunggu hasil WO - executor tidak terhubung ke
  Oracle (pembacaan DB ditolak pengaman sesi; tidak dicoba).
- Baris 100001 (Cause of Loss kosong) tidak dipakai produk; dibiarkan (tidak diisi, tidak dihapus) - dilaporkan saja.
- Kunci `pxObjClass` / `pyRuleHarness` tidak dipindah - hanya di cadangan P3.

## Menjalankan uji modul ini saja

```powershell
go test ./modul/causeoflosslife/...
go vet -tags db ./modul/causeoflosslife/...   # uji seam Oracle (skema uji) - butuh ORACLE_DSN skema uji
npx vitest run modul/causeoflosslife
```
