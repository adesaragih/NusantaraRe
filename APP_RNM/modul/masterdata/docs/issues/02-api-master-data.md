# 02: Master Data — pendaftaran modul, migrasi, dan API (backend sesi c3)

> ⚠️ **Disusun agent sesi c3 atas permintaan sesi `nusantarare-0f` (rencana `01-rencana-master-data.md`) — bukan hasil
> `/to-tickets`.** Lingkup inti / frontend diizinkan pengguna sesi c3 (AskUserQuestion 04-10-2026: "Semua, tandai untuk
> tim inti"). Setiap berkas di `inti/` dan `frontend/` WAJIB ditinjau tim inti (CODEOWNERS) sebelum commit.

**Status:** backend selesai 04-10-2026 (uji hijau), jejak ubah MD-7 (migrasi 762) menyusul hari yang sama.
Frontend dibangun sesi 0f.
⛔ Migrasi 760 / 761 / 762 / 990 (modul ini) dan 904 (inti) ditulis, **belum dijalankan** — urutan pelari menurut nama:
904 → 760 → 761 → 762 → 990.

## Pendaftaran (inti, langkah 1)

- ⚠️ **Pelari migrasi (inti, 04-10-2026, sesudah bukti DEV dari sesi 0f: keenam objek masih `VIEW`)**:
  `inti/backend/migrasi/migrasi.go` pra-terbang membandingkan kolom objek bernama sama dengan setiap `CREATE TABLE` yang
  belum tercatat - dan `SYS.ALL_OBJECTS` / `SYS.ALL_TAB_COLUMNS` memuat VIEW. Tanpa perbaikan, 760 berhenti di
  pra-terbang ("kolom yang diminta migrasi tetapi tidak ada: [STS_AKTIF]") dan SELURUH migrasi tertunda ikut tertahan.
  Perbaikan: perbandingan dilewati hanya bila objeknya VIEW **dan** langkah yang sama membongkarnya (`DROP VIEW
  {skema}.X`) sebelum `CREATE TABLE {skema}.X` (`ViewDibongkarDulu`, `adalahView`); tabel yang sudah berdiri tetap
  dibandingkan. Uji: `TestViewDibongkarDulu` (migrasi) dan `TestViewDigantiTabelDikenaliPraTerbang` (penjaga, jawaban
  diketahui: enam view 760). Jalur Oracle `adalahView` belum teruji tanpa basis data. 760 TIDAK diubah.

- `inti/backend/migrations/904_m_nav_menu_masterdata.sql` (+ `_down`): baris `M_NAV_MENU` `masterdata` / `Master Data`,
  `MASTER`, URUTAN 4 (MASTER sudah berisi mastercontractretrolife 1, masterproductnamelife 2, treatycontractout 3),
  `DIMIGRASI '0'`; bentuk datar idempoten. Akses akun TIDAK diisi (Kelola User, M-6).
- Penjaga: `inti/backend/penjaga/menu_test.go`, `rentang_test.go` (904 diterapkan skema tiruan; hanya INSERT bentuk datar
  + DELETE mundurnya yang sah), `inti/frontend/uji/menuBersih.ts`.
- Frontend: `frontend/katalogKorpus.ts` (`MODUL_LUAR_KORPUS`, `KATALOG_MODUL`; `FOLDER_KORPUS` tetap 20), `Beranda.tsx`,
  uji `Beranda` / `Shell` / `daftar.menuTabel` / `daftar.sinkron`.
- `inti/backend/daftar/modul_masterdata_gen.go` (`go generate ./inti/backend/daftar`).

## Migrasi modul (rentang 760–799, slot 990)

| Berkas | Isi |
| --- | --- |
| `760_view_ke_tabel_flat.sql` | view PROVINCE, CITYINPUT, DISTRICTINPUT, ACCUMULATEDTYPE, CZONE, ACCUMULATION → tabel flat bernama sama (+ `STS_AKTIF VARCHAR2(1) DEFAULT '1' NOT NULL`); jalur mundur memulihkan view persis teks `DDL\<NAMA>.txt` |
| `761_t_master_status.sql` | `T_MASTER_STATUS (NAMA_TABEL, ID_BARIS, STS_AKTIF)` — status master bertabel warisan tanpa kolom status (NATION) |
| `762_jejak_ubah_master.sql` | jejak ubah (MD-7): `CREATE_OP VARCHAR2(64)`, `TGL_CREATE DATE`, `UPDATE_OP VARCHAR2(64)`, `TGL_UPDATE DATE` ke keenam tabel flat dan `T_MASTER_STATUS` (jejak NATION / OBJECTITEMTYPE); `T_MASTER_STATUS.ID_BARIS` → VARCHAR2(4000) |
| `990_menu_masterdata.sql` | slot menu: `DIMIGRASI = '1'` |

Peta kolom: `docs/STRUKTUR-TABEL-MASTER-DATA.md`.

## Kontrak API

- `GET /api/masterdata` → `{ tabel: [{ kunci, judul, idOtomatis, kolom: [{ kunci, kolom, lebar, wajib, turunan }] }] }` —
  delapan master, urutan tampil: `nation`, `province`, `city`, `district`, `czone`, `accumulatedtype`, `accumulation`,
  `objectitemtype`.
- `GET /api/masterdata/{tabel}?q=&status=&halaman=` → `{ baris: [{ <kunci kolom>: teks, …, aktif: boolean }], total,
  halaman, ukuran }`; `q` Contains tidak peka huruf atas ID + kolom nama; `status` = `aktif` / `nonaktif` / kosong;
  `halaman` ≥ 1, 50 baris per halaman, urut ID.
- `POST /api/masterdata/{tabel}` badan objek bernilai teks (kunci kolom; turunan diabaikan) → **201** `{ id }`. Master
  `accumulation`: ID dibuat backend, badan WAJIB memuat `negara` (awalan ID, mis. `INA`).
- `PUT /api/masterdata/{tabel}/{id}` badan sama (ID dari rute; `id` di badan diabaikan) → `{ id }`.
- `PUT /api/masterdata/{tabel}/{id}/status` badan `{ "aktif": true|false }` → `{ id, aktif }`. Tidak ada hapus.
- Galat: 400 isian (wajib, lebar, medan asing, rujukan tak ada, `q` > 255, status / halaman), 401 tanpa identitas
  (tulis), 404 master / ID tak ada, 409 ID ganda (akumulasi: Note + zip ganda — pesan prosedur), 503 tanpa basis data.

### Bentuk baris per master (kunci JSON; *turunan* = diisi backend)

| Master | Kunci |
| --- | --- |
| `nation` | `id` (wajib, ≤ 10), `oldId` (≤ 6), `note` (wajib, ≤ 100), `nationInitial` (≤ 20) |
| `province` | `id` (wajib), `nationId` (rujukan nation), `note` (wajib), *`nationName`* |
| `city` | `id` (wajib), `provinceId` (rujukan province), `note` (wajib), `branchId` (rujukan BRANCH), `email`, `moId`, `jabodetabekStatus` |
| `district` | `id` (wajib), `cityId` (rujukan city), `districtName` (wajib) |
| `czone` | `id` (wajib), `code` (wajib), `description`, `groupOf` (rujukan czone.code), *`groupOfName`* |
| `accumulatedtype` | `id` (wajib), `accumulationType` (wajib), `keyword`, `note`, `type` |
| `accumulation` | *`id`*, `accumulation` (rujukan accumulatedtype), *`accumulationName`*, `note` (wajib, disimpan HURUF BESAR), `keyword`, `scopeArea`, `cZone`, `cZoneId`, *`province`*, `provinceId` (rujukan province), `zipCode` (wajib), `accumulationType`, `syariahStatus` + masukan `negara` saat tambah |
| `objectitemtype` | `id` (wajib), `objectItemType` (wajib), `note`, `pctAdjustable1`, `pctAdjustable2`, `objectItemTypeIna`, `group`, `type` |

Lebar kolom selain nation = 4000 byte.

**Jejak ubah (MD-7) — setiap master:** `createOp`, `tglCreate`, `updateOp`, `tglUpdate` (turunan, baca-saja; ada di
`GET /api/masterdata` sesudah kolom data). Diisi backend dari pelaku: POST → pembuat + tanggal buat; PUT dan PUT status
→ pengubah + tanggal ubah. Tanggal berbentuk `YYYY-MM-DD HH24:MI:SS` (waktu basis data, SYSDATE). Kosong untuk baris yang
disalin 760 sebelum disentuh menu. NATION / OBJECTITEMTYPE: jejaknya di `T_MASTER_STATUS`; menambah ID yang barisnya
dulu dihapus di luar menu mengatur ulang baris sisanya (aktif, pembuat baru, pengubah kosong) — bukan galat PK. Dikirim balik di badan = diabaikan. Akun pelaku > 64 byte → 400.

## Penyaring aktif di nbfacin (langkah 4)

Saran / popup akumulasi nbfacin (`modul/nbfacin/backend/repository/akumulasi.go`) kini hanya baris AKTIF: province /
accumtype / czone `STS_AKTIF = '1'`; city / district lewat `CITYINPUT` / `DISTRICTINPUT` (view CITY / DISTRICT tidak
membawa status); nation lewat `T_MASTER_STATUS`; pencarian akumulasi (jalur RD, kota / kecamatan, nomor polis) hanya
akumulasi aktif. Object Item Type sudah menyaring `ISACTIVE = '1'` sejak tiket 39 (RD BrowseV_JN_OBJ_ITEM). Area (RW)
memakai `STS_AKTIF` RW sendiri.

## Keputusan agent — DISETUJUI work owner 04-10-2026

> Jawaban work owner atas laporan sesi 0f (yang menyebut M-1 … M-8 dan MD-1 … MD-9): "setuju". ~~Pengecualian: MD-7 (jejak ubah) memuat pilihan yang belum dijawab — tetap `[pertanyaan terbuka]`.~~ MD-7 dijawab
> kemudian: "Kolom diubah oleh/tanggal" — dibangun migrasi 762 (MD-10).

Cakupan persetujuan di berkas ini: MD-1 … MD-9.

| # | Keputusan | Dasar |
| --- | --- | --- |
| MD-1 | Kepemilikan (M-2): migrasi flat **dipindah** dari nbfacin 196 ke modul ini (760) — tabel milik Master Data, nbfacin membaca | 196 ditahan work owner (belum dijalankan) — pindah aman; kepemilikan satu tempat |
| MD-2 | NATION (tabel warisan) TIDAK diubah strukturnya; statusnya di tabel baru `T_MASTER_STATUS` (tanpa baris = aktif). OBJECTITEMTYPE memakai `ISACTIVE` sendiri (`'1'` aktif `[terverifikasi]` RDBList GetObjectItembyName_SQL / GetDataObjectItem; `'0'` saat dinonaktifkan `[dugaan]`) | `ALTER` tabel warisan = mengambil kepemilikannya (`TestTabelBukanMilikKitaTidakDibuat`) — keputusan work owner, bukan modul |
| MD-3 | ACCUMULATION: ID = `negara-ZIPCODE-lpad(ACCUMULATION_SEQ.NEXTVAL,6,'0')` (prosedur RDBMASTERACCUMULATION, ADR-0043: tanpa prosedur); `negara` = masukan wajib saat tambah `[dugaan]` arti (param prosedur `p_COUNTRY`, contoh korpus `INA`); sequence `ACCUMULATION_SEQ` diandaikan ada di skema `[dugaan]` (DDL-nya tidak ada); Note + zip ganda → 409; NOTE huruf besar; `PROVINCE` diturunkan dari `PROVINCE.NOTE` `[dugaan]` (view asal membaca JSON ProvinceName) | `DDL\RDBMASTERACCUMULATION.txt` `[terverifikasi]` |
| MD-4 | Rujukan yang diisi wajib ada (400); kolom turunan diisi backend dari definisi view asal (NATIONNAME, GROUPOFNAME, ACCUMULATIONNAME); `CZONE.TGLUPDATE` / `USERID` tidak ditulis (format belum terverifikasi) | M-5 rencana |
| MD-5 | Penyaring aktif nbfacin seperti di atas, termasuk ketiga jalur pencarian akumulasi | M-3 rencana: "Popup / saran di modul lain hanya menampilkan baris aktif" |
| MD-6 | Daftar 50 baris per halaman, urut ID, `q` atas ID + kolom nama | — |
| ~~MD-7~~ | ~~Jejak ubah (M-7) BELUM dibangun~~ → **dijawab work owner: "Kolom diubah oleh/tanggal"** (hanya perubahan terakhir; tabel riwayat umum tidak dipilih) — dibangun migrasi 762 | `T_CLAIMLF_JEJAK` khusus klaim, tidak dipakai |
| MD-10 | Nama / tipe kolom jejak = pola `T_WORK_*` (claimlife 001, premiumlistlife 059): `CREATE_OP VARCHAR2(64)` (= `M_LOGIN_GO.LOGIN_ID`), `TGL_CREATE DATE`, pasangan `UPDATE_OP VARCHAR2(64)`, `TGL_UPDATE DATE`; boleh kosong. Jejak NATION / OBJECTITEMTYPE di `T_MASTER_STATUS` (MD-2: warisan tidak di-ALTER), `ID_BARIS` dilebarkan ke 4000 = lebar OBJECTITEMTYPE.ID. Tanggal ditulis SYSDATE di SQL, dibaca TO_CHAR berformat tetap | **DISETUJUI** work owner 04-10-2026 ("setuju") |
| MD-8 | Tanpa aturan peran baru (M-6): tulis butuh identitas (401) seperti modul lain; akses menu lewat Kelola User | — |
| MD-9 | ID tidak dapat diubah (PUT mengabaikan `id` badan); tanpa hapus | M-3 rencana |
