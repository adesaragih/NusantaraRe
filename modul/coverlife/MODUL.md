# Modul `coverlife` — Cover Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`PANDUAN-TIM-PER-MODUL.md` bab 5) — keputusan work owner 08-10-2026 (prompt
`PROMPT-DISEASELIFE-COVERLIFE.md`, K0, C1-C4, K5 final). Panduan PATOKAN: section Pega `InboxCoverLife` (kelas
`ASM-FW-GISFW-Int-COVER_LIFE`, judul "COVER" b368) di `D:\NUSARE DEV\Menu Cover\InboxCoverLife.xml` (satu-satunya XML,
6.685 baris, dibaca utuh; `bNNNN` = baris XML). Templat: modul `causeoflosslife` (satu tabel `ID` + nama, Edit / Save /
Cancel, migrasi modul berjatah). Paritas: `docs/PARITAS-LAYAR-DAN-AKSI.md`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Migrasi modul ini tinggal
di `backend/migrations/` (rentang inti 900-949 sudah penuh). Jatahnya **dipinjam**: rentang dari premiumlistlife
(`050-089` → `050-079`), slot dari treatycontractout (`956-957` → `956-956`) — keputusan work owner 08-10-2026 K0,
**perlu persetujuan tim inti (CODEOWNERS)**; dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `coverlife` |
| Folder korpus | `Cover Life` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-COVERLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `085-089` |
| Slot menu | `957-957` |
| Prefix rute API | `/api/cover-life` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Keputusan work owner 08-10-2026 (final, dikutip dari prompt)

| # | Keputusan | Penerapan |
| --- | --- | --- |
| K0 | Nomor: premiumlistlife `050-079`, modul ini `085-089`; treatycontractout slot `956`, modul ini slot `957`; periksa dulu nomor itu kosong | Diperiksa sebelum menulis (Glob + isi MODUL.md) - lihat `modul/diseaselife/MODUL.md` K0 (satu pemeriksaan untuk kedua modul). Tabel jatah diubah di `premiumlistlife`, `claimlife`, `treatycontractout` + dua modul baru |
| C1 | `M_COVER_LIFE` dijadikan FLAT di tempat, nama TETAP, **TANPA RENAME**: ADD (COVER, NOTE) berdiri sendiri; isi notasi titik PERSIS teks view SELAGI view + JSONDATA ada; periksa (NULL = NULL sama); DROP JSONDATA; DROP VIEW `COVER_LIFE` TERAKHIR (berpelindung) | 085 (ALTER ADD), 086 (isi → periksa → buang JSONDATA → DROP VIEW; **nol perubahan `perintah_katalog.go`**). Pemaksa berhenti = `SET m.ID = NULL` (pola 944 / 092; ID NOT NULL) → ORA-01407, tanpa kutip. Hasil: SATU tabel `M_COVER_LIFE (ID, COVER, NOTE)`, nol objek `COVER_LIFE`; uji penjaga modul `TestKodeTidakMenyebutViewLama` |
| C2 | `PEGA_M_COVER_LIFE` INVALID diterima; ID baru = rumus Pega `'1' \|\| LPAD(M_COVER_LIFE_SEQ.NEXTVAL, 5, '0')`; galat jelas bila bertabrakan | `models.BentukID`; `services.ErrIDTerpakai` → **409** (uji `TestSimpanIDDariSequenceSudahAda`, `TestRuteIDSequenceSudahAda`) |
| C3 | Cover wajib (XML), Note opsional; huruf besar hanya bila XML punya pengubahnya; Cover unik + pangkas = di luar XML | wajib b843 / b895; Note `pyRequired` false b1020 / b1070; XML TANPA pengubah huruf (`pyBehaviors` kosong b897 / b1072) → huruf tidak diubah; pangkas + kembar tanpa beda huruf → 422 (PARITAS "Di luar XML") |
| C4 | Menu `coverlife`, LABEL `Cover Life`, MASTER TREATY URUTAN 16, DIMIGRASI '1', di slot 957 | `957_menu_coverlife.sql` (INSERT datar `NOT EXISTS`); penjaga `modulLuarKorpus` + `barisLahirDiSlot`; `MODUL_LUAR_KORPUS.coverLife` |
| K5 | Hak menu = salinan pemegang `causeoflosslife` beserta HAK-nya, aman diulang | `docs/sql/cov_d_hak.sql` (LANGKAH-WO (d)) |

Prasyarat P1-P3 (cadangan PALING ATAS; penulis berhenti + kunci; acuan dua kali) + bukti: `docs/LANGKAH-WO-COVERLIFE.md`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-COVERLIFE.md`, `PARITAS-LAYAR-DAN-AKSI.md`, `LANGKAH-WO-COVERLIFE.md` + `sql/cov_*.sql`, `PR-COVERLIFE.md` |
| `backend/` | `models/` `repository/` (`cvl.go`, `cvl_tabel.go`) `services/` `handlers/` `tiruan/` `migrations/` `modul.go` |
| `frontend/` | `pages/CoverLife.tsx` `labels.ts` `api.ts` `aturan.ts` `coverlife.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Asumsi — bukti XML atau bawaan work owner

XML hanya rule **section**: isi `AddToList_Act`, `EditList_DT`, `NewData_DT`, dan report definition `BrowseCoverLife_RD`
TIDAK ada (hanya dirujuk).

| # | Aturan yang berlaku | Status dan bukti |
| --- | --- | --- |
| A1 | Form TANPA medan ID: Add = ID dari sequence (server), Edit = ID baris (ditampilkan sebagai teks status, bukan medan) | **[terverifikasi]** form S2 hanya Cover b848 dan Note b1026; ID hanya di grid (b3023 / b3449) dan `EditList_DT` b3811; isian `id` di badan = 400 |
| A2 | Save = Add bila form tanpa ID, Edit sesudah Edit baris | **[penyimpangan sadar - menunggu WO]** satu Save b5407 → `AddToList_Act` b5426; isi aktivitas tidak ada (pola causeoflosslife T1) |
| A3 | Grid ID menaik tetap, tanpa saring, tanpa urut pilihan, **50** baris, halaman bernomor; Note tidak tampil | **[terverifikasi]** `pySortType` ASC b3967 / `pySortOrder` 1 b3972; `pyColumnSorting` false b3969 / b3991 / b4013; `pyGridFiltering` false b4080; `pyPageSize` 50 b4101; Numeric b4072; kolom grid hanya ID b3023, Cover b3161, Edit b3277 |
| A4 | Huruf Cover / Note tidak diubah | **[terverifikasi]** nol data transform / refresh pada kedua medan (b897, b1072) |
| A5 | Tanpa Delete, Upload, ringkasan / detail | **[terverifikasi]** `pyGridDeleteActivityExists` false b457 / b749 / b1638 / b1925 / b2633 / b2882 / b4229 / b5236; tidak ada tombol lain |

## Migrasi

Rentang `085-089` (terpakai `085-086`), slot menu `957-957` (terpakai `957`). Nama berkas yang sudah dijalankan tidak
pernah diubah (`T_MIGRASI` mencatat nama).

| No | Berkas | Isi |
| --- | --- | --- |
| 085 | `085_cover_life_kolom.sql` | ALTER ADD `COVER VARCHAR2(200)`, `NOTE VARCHAR2(1000)`, berdiri sendiri |
| 086 | `086_cover_life_satu_tabel.sql` | isi dari JSONDATA → pemeriksaan C1.3 → DROP JSONDATA → DROP VIEW `COVER_LIFE` (aman diulang) |
| 957 | `957_menu_coverlife.sql` | INSERT baris `M_NAV_MENU` (C4), `NOT EXISTS` |

Urutan pelari di skema baru: 085-086 → 900 … 949 → 957 (uji `TestUrutanPelari085Lalu900Lalu957`; versi Oracle
`TestDBUrutanPelari085Lalu900Lalu957`, `-tags=db`, belum dijalankan). Skema uji punya tiruan `M_COVER_LIFE` + view +
sequence SEBELUM pelari (`uji/skemauji/cover_tiruan.go`). Karena 085-086 mengubah bentuk tabel Pega, `M_COVER_LIFE`
TIDAK dinyatakan "Tabel warisan" (preseden causeoflosslife); kolom yang dibuat migrasi tercatat di
`docs/STRUKTUR-TABEL-COVERLIFE.md`.

## Pembaca lain

Tidak ada pembaca lain di repo (fakta WO; dicari ulang: nol rujukan `COVER_LIFE` / `M_COVER_LIFE` di luar folder ini
dan skema uji). `COVERAGE`, `COVERAGE_FACIN`, `COVERAGETRAVEL*`, `COVERNOTE*` milik aplikasi lain - TIDAK disentuh.

## Butir terbuka

- **Bukti** (`docs/sql/cov_bukti.sql`), acuan P2, dan keadaan menunggu hasil WO - executor tidak terhubung ke Oracle
  (pembacaan DB ditolak pengaman sesi; tidak dicoba).
- Kunci `pxObjClass` / `pyRuleHarness` tidak dipindah - hanya di cadangan P3.

## Menjalankan uji modul ini saja

```powershell
go test ./modul/coverlife/...
go vet -tags db ./modul/coverlife/...   # uji seam Oracle (skema uji) - butuh ORACLE_DSN skema uji
npx vitest run modul/coverlife
```
