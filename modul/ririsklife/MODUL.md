# Modul `ririsklife` — R/I Risk

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`PANDUAN-TIM-PER-MODUL.md` bab 5) — keputusan work owner 08-10-2026 (prompt
`PROMPT-RIRISKLIFE.md`, K1-K4 final). Panduan PATOKAN: section Pega `InboxSummaryRIRisk` (kelas
`ASM-FW-GISFW-Int-RIRISK_LIFE_SUMMARY`, judul "R/I RISK SUMMARY") dan `InboxRIRisk` (kelas `ASM-FW-GISFW-Int-RI_RISK_LIFE`,
judul "R/I RISK DETAIL") di `D:\NUSARE DEV\Menu RI Risk\`. Templat: modul saudara `ricommlife`. Paritas setiap tombol /
aktivitas: `docs/PARITAS-LAYAR-DAN-AKSI.md`.

⛔ **Tanpa migrasi sendiri** (`—` di bawah): tabel = migrasi inti `935`-`940` (RENAME tabel Pega + satu tabel per jenis
data), baris menu = migrasi inti `941`.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `ririsklife` |
| Folder korpus | `R/I Risk Life` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-RIRISKLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `—` |
| Slot menu | `—` |
| Prefix rute API | `/api/ri-risk-life` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Keputusan work owner 08-10-2026 (final, dikutip dari prompt)

| # | Keputusan | Penerapan |
| --- | --- | --- |
| K1 | Nama akhir `RIRISK_LIFE_SUMMARY` / `RIRISK_LIFE`, tabel FLAT tanpa JSONDATA, BUKAN tabel baru: DROP VIEW (hanya bila VIEW, berpelindung) lalu `ALTER TABLE M_RIRISK_LIFE* RENAME TO …` (berpelindung: sumber ada DAN target belum TABLE), tambah kolom, isi ulang SEMUA kolom dari JSONDATA, `DROP COLUMN JSONDATA CASCADE CONSTRAINTS` | migrasi inti 935-940; DROP VIEW berpelindung `SYS.ALL_VIEWS` dan RENAME berpelindung `ALL_TAB_COLUMNS` tabel sumber (target TABLE yang sudah ada = ORA-00955, gagal keras) - DUA bentuk blok baru diterima penjaga (tinjauan tim inti, `docs/PR-RIRISKLIFE.md`) |
| K2 | Kolom akhir = kolom view lama (urutan sama); tipe warisan untuk kolom yang ada (`RISK` NUMBER tanpa skala); kolom baru pola ricommlife; konversi RISK tanpa NLS (bukti SELECT baca-saja); PK `RIRISK_LIFE(ID)` hanya bila pola ricommlife; indeks IDUSEDBY wajib tanpa kembar | 936 (USEDBY 200, MODIFIEDDATE 50, OPERATORID 200), 939 (AGE VARCHAR2(10)); 940: konversi TRANSLATE koma/titik → pemisah desimal sesi + TO_NUMBER; `PK_RIRISK_LIFE` (rincian ricommlife ber-PK); indeks hanya bila belum ada indeks berkolom pertama IDUSEDBY (`SYS.ALL_IND_COLUMNS` - INDEX4) |
| K3 | Prosedur `PEGA_M_RIRISK_LIFE*` INVALID - diterima; sequence TETAP; `M_RIRISK_LIFE_TEMP` tidak disentuh | `rirk_tabel.go` (sequence warisan); nol rujukan ke TEMP (uji `TestMigrasiRiRiskTanpaCreateTable`, `TestPeriksaTulis`) |
| K4 | Menu `ririsklife`, LABEL `R/I Risk`, MASTER TREATY URUTAN 11, DIMIGRASI '1', hak SUPERADMIN lewat LANGKAH-WO (d) | `941_m_nav_menu_ririsklife.sql`; penjaga `modulLuarKorpus` + `labelTampilDisetujui` (`R/I Risk Life` → kode, tampil `R/I Risk`); `MODUL_LUAR_KORPUS.riRiskLife` |

## RALAT

| # | Bunyi lama | Bunyi baru | Bukti |
| --- | --- | --- | --- |
| R1 (08-10-2026) | A4: *"Ringkasan TANPA Add / Edit-simpan; ringkasan baru hanya lewat Simpan Upload"* (Save / Cancel di wadah `1=2` S1511) | **Keputusan work owner 08-10-2026: Save / Cancel ringkasan DITAMPILKAN, sama dengan riratelife / ricommlife** - `POST` (201, ID site \|\| LPAD(M_RIRISK_LIFE_SUMMARY_SEQ, 6)), `PUT` (ubah nama; USEDBY rincian ikut, satu transaksi), MODIFIEDDATE / OPERATORID seperti ricommlife; nama wajib, dipangkas, tidak kembar tanpa beda huruf (indeks `IX_RIRISK_LIFE_SUMMARY_NAMA` 937 sudah ada); View only 403. Label MONTH tetap "MONTH" (A9 diputuskan). Perluasan `perintah_katalog.go` (ALL_VIEWS, ALL_IND_COLUMNS) dan RENAME berpelindung DITERIMA WO | uji `TestRuteRingkasan` (201/200/422/403/404), `TestAdd`, `TestValidasiNama`, `TestEditGantiNamaRincian`, `-tags=db` `TestDBSimpanRingkasan`, vitest `ririsklife.test.ts`; `docs/PARITAS-LAYAR-DAN-AKSI.md` |

Prasyarat P1-P3 (penulis berhenti + kunci + acuan dua kali; pemakaian skema GL; cadangan CSV): `docs/LANGKAH-WO-RIRISKLIFE.md`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-RIRISKLIFE.md`, `PARITAS-LAYAR-DAN-AKSI.md`, `LANGKAH-WO-RIRISKLIFE.md` + `sql/ririsk_*.sql`, `PR-RIRISKLIFE.md` |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `modul.go` (tanpa `migrations/`) |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `ririsklife.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Asumsi — bukti XML atau bawaan work owner

Kedua XML hanya rule **section**: isi `AddToList_Act`, `AddToListSummary_Act`, `SubmitRIRisk_Act`,
`DeleteSummaryDetail`, `setIDUsedBy_Act`, `EditRIRiskLife_Act`, `NewRIRiskLife_Act`, `clearInputFieldRIRisk_act`,
`ChangeDotToPoint_DT` TIDAK ada. Ada bukti → **[terverifikasi]**; tidak ada → pola ricommlife, **[penyimpangan sadar -
menunggu WO]**.

| # | Aturan yang berlaku | Status dan bukti |
| --- | --- | --- |
| A1 | Kombinasi (CONTRACT, YEAR, MONTH) kembar dalam satu ringkasan DITOLAK (Detail Save dan Upload; nol depan setara) | **[penyimpangan sadar - menunggu WO]** - pola ricommlife (CONTRACT, YEAR) + MONTH, karena baris DEV memakai YEAR ATAU MONTH; isi `AddToList_Act` tidak ada |
| A2 | CONTRACT wajib bulat ≤ 10 angka; YEAR / MONTH kosong atau bulat ≤ 4 angka; RISK wajib desimal TIDAK negatif, koma atau titik, ≤ 38 angka | wajib / tidak wajib **[terverifikasi]** D1933 / D2216 / D2403 / D2595, pyMax 4 D2220 / D2407; batas angka dan "tidak negatif" **[penyimpangan sadar - menunggu WO]** |
| A3 | Detail Save dan Simpan Upload memperbarui MODIFIEDDATE dan OPERATORID ringkasan (transaksi yang sama) | **[penyimpangan sadar - menunggu WO]** - pola ricommlife |
| A4 | Ringkasan Add / Edit lewat Save (Cancel saat Edit); Edit nama ikut mengganti USEDBY rinciannya (satu transaksi) | **[keputusan work owner 08-10-2026]** (RALAT R1) - di XML wadahnya `1=2` S1511; perilaku = ricommlife |
| A5 | CSV: kepala USEDBY, CONTRACT, YEAR, MONTH, RISK (urutan bebas); pemisah `;`/`,`; 4 MB, 10.000 baris, 100 nama; nama dicocokkan ke ringkasan atau dibuat baru | kolom **[terverifikasi]** S5495; sisanya pola ricommlife **[penyimpangan sadar - menunggu WO]** |
| A6 | R/I RISK DETAIL tanpa Delete per baris; View only tanpa Edit / Delete / Upload / form detail | tanpa Delete **[terverifikasi]** (D hanya EDIT D9784); View only = pola ricommlife |
| A7 | ID yang ternyata terpakai dilewati (≤ 100 nomor); nomor sequence melebihi lebar LPAD = galat | rumus **[terverifikasi]** (teks prosedur, fakta WO); penolakan = **[penyimpangan sadar]** (LPAD Oracle memotong diam-diam) |
| A8 | Delete ringkasan menghapus rinciannya dalam SATU transaksi, TANPA cek rujukan produk | berantai **[terverifikasi]** S9648 `DeleteSummaryDetail`; ⚠️ produk life (`RIRISKID` MPNL) yang menunjuk ringkasan itu tertinggal |
| A9 | Label medan MONTH = "MONTH" | **[keputusan work owner 08-10-2026]** - `pyLabelFieldValue` D2358 bertuliskan "YEAR" (salin-tempel Pega), preview D2381 dan grid D8277 "MONTH" |
| A10 | AGE tidak dibaca / ditulis | **[terverifikasi]** tidak ada medan AGE di D; kunci AGE nol baris DEV |

## Migrasi

Nol migrasi modul. Migrasi inti `935_ririsk_life_summary_ganti_nama.sql`, `936_ririsk_life_summary_kolom.sql`,
`937_ririsk_life_summary_satu_tabel.sql`, `938_ririsk_life_ganti_nama.sql`, `939_ririsk_life_kolom.sql`,
`940_ririsk_life_satu_tabel.sql`, `941_m_nav_menu_ririsklife.sql` (masing-masing + `_down`). Karena 935-940 mengubah
bentuk tabel Pega `M_RIRISK_LIFE*`, `RIRISK_LIFE_SUMMARY` / `RIRISK_LIFE` TIDAK dinyatakan "Tabel warisan" (preseden
`adjusterconsultant` 870, riratelife 927/929, ricommlife 931/933); kolom yang dibuat migrasi tercatat di
`docs/STRUKTUR-TABEL-RIRISKLIFE.md`.

## Butir terbuka

- **Bukti konversi RISK** (`docs/sql/ririsk_bukti_konversi.sql`) dan semua angka acuan menunggu hasil WO - executor
  tidak terhubung ke Oracle.
- **INDEX4**: kolomnya tidak tercatat di repo; 940 hanya membuat `IX_RIRISK_LIFE_IDUSEDBY` bila belum ada indeks
  berkolom pertama IDUSEDBY. LANGKAH-WO (a) kueri 8 mencatatnya.
- **Yatim**: 2 rincian (IDUSEDBY 1000085, 1000087) dan 2 ringkasan tanpa rincian - dipindah apa adanya, TIDAK dihapus;
  rincian yatim tidak tampil di layar mana pun (Detail dibuka dari ringkasan).
- **pxObjClass** ringkasan / rincian tidak dipindah (bukan kolom view) - hanya di cadangan CSV P3.

## Menjalankan uji modul ini saja

Dari folder akar repo:

```powershell
go test ./modul/ririsklife/...
go vet -tags db ./modul/ririsklife/...   # uji seam Oracle (skema uji) - butuh ORACLE_DSN skema uji
npx vitest run modul/ririsklife
```
