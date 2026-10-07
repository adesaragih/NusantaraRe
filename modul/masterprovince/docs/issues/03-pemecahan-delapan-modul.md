# 03: Pemecahan Master Data menjadi delapan modul menu (backend / inti sesi c3)

> ⚠️ **Disusun agent sesi c3 atas permintaan sesi `nusantarare-9d` — bukan hasil `/to-tickets`.** Keputusan di bawah
> diambil work owner lewat AskUserQuestion sesi 9d (04-10-2026) atas "rancangan 8-modul" sesi c3, diteruskan sesi 9d.
> Setiap berkas di `inti/` WAJIB ditinjau tim inti sebelum digabung. Tanpa commit / push (aturan work owner terbaru).

**Status:** backend + inti selesai 04-10-2026 (uji Go hijau, kecuali kegagalan lama di luar pekerjaan ini). Frontend
(`inti/frontend/master`, frontend delapan modul, perakit) dibangun sesi 9d. ⛔ Migrasi 912–920 ditulis, **belum
dijalankan**.

## Permintaan

Work owner (04-10-2026, dengan gambar): *"ini masing dibuatin didalam sub menu master data, bukan di jadiin 1 gini. jadi
saat master di klik, keluar sub menu isi nya itu"* · *"bukan di satuin begini, di pisah per sub modul, seperti gbr ke 2"*
(grup MASTER: Marketing Officer, Company Detail, Accounts — masing-masing satu modul). "Cara pecah": **"8 modul
terpisah"** (ditolak: 1 modul dengan 8 menu).

## Keputusan

| # | Pertanyaan | Jawaban work owner |
| --- | --- | --- |
| K1 | Nama / label / prefix | sesuai usulan: tabel di bawah |
| K2 | Nasib modul `masterdata` | **"Dihapus, tabel pindah ke Province"** — 880–882 pindah apa adanya (nama sama) ke `masterprovince` |
| K3 | Mesin generik | **"Pindah ke inti"** — `inti/backend/master` (c3), `inti/frontend/master` (9d) |
| K4 | Slot menu | **"Modul luar korpus tanpa slot"** — baris 912–919 lahir DIMIGRASI '1' |
| K5 | Hak akun | **"Ya, salin otomatis"** — 920 menyalin hak 'masterdata' ke delapan KODE |
| K6 | `/api/masterdata` | tidak ditanyakan; rekomendasi dipakai: dibuang sekaligus, backend + frontend dideploy bersama |

| Modul | LABEL | Master | Prefix | Rentang | Menu (inti) | URUTAN MASTER |
| --- | --- | --- | --- | --- | --- | --- |
| `masternation` | Nation | `nation` | `/api/master-nation` | 884-885 | 912 | 6 |
| `masterprovince` | Province | `province` | `/api/master-province` | 880-883 | 913 | 7 |
| `mastercity` | City | `city` | `/api/master-city` | 886-887 | 914 | 8 |
| `masterdistrict` | District | `district` | `/api/master-district` | 888-889 | 915 | 9 |
| `masterczone` | CZone | `czone` | `/api/master-czone` | 890-891 | 916 | 10 |
| `masteraccumulatedtype` | Accumulated Type | `accumulatedtype` | `/api/master-accumulated-type` | 892-893 | 917 | 11 |
| `masteraccumulation` | Accumulation | `accumulation` | `/api/master-accumulation` | 894-895 | 918 | 12 |
| `masterobjectitemtype` | Object Item Type | `objectitemtype` | `/api/master-object-item-type` | 896-897 | 919 | 13 |

URUTAN `[terverifikasi]` dari hasil bersih penjaga menu (900 + 901 + 906–909): MASTER sesudah 909 = Contract Retro Life 1,
Product Name Life 2, Marketing Officer 3, Company Detail 4, Accounts 5; Master Data (6) dibuang. 898-899 tidak dibagikan
(cadangan). Hanya `masterprovince` bermigrasi (880–882); tujuh lainnya belum butuh migrasi.

## Kontrak rute (per modul, prefix P)

Galat selain 2xx berbadan `{"galat": "<pesan>"}`: 400 masukan · 401 belum login · 403 tanpa menu modul ini (gerbang
`cmd/api`) · 404 ID / rujukan tidak ada · 409 ID sudah ada · 503 tanpa Oracle · 500.

| Rute | Jawaban |
| --- | --- |
| `GET P/meta` | `{kunci, judul, idOtomatis, kolom: [{kunci, kolom, lebar, wajib, turunan}], rujukan: [{kunci, judul, nilai, nama}]}` — kolom data lalu `createOp` / `tglCreate` / `updateOp` / `tglUpdate` (turunan); `rujukan` selalu larik |
| `GET P?q=&status=&halaman=` | `{baris: [{<kunci kolom>: teks, aktif: bool}], total, halaman, ukuran}` — status `""` / `aktif` / `nonaktif`; halaman ≥ 1; 20 per halaman |
| `POST P` | badan `{kunci: teks}` → 201 `{id}`; turunan / jejak yang dikirim diabaikan, kunci asing 400; Accumulation tanpa `id`, wajib `negara` |
| `PUT P/{id}` | badan `{kunci: teks}` → `{id}` |
| `PUT P/{id}/status` | badan `{"aktif": bool}` → `{id, aktif}` |
| `GET P/rujukan/{kolom}?q=&halaman=` | bentuk sama dengan daftar: baris master YANG DIRUJUK kolom itu, AKTIF saja; kolom tanpa rujukan master (BRANCHID) / asing → 404 |

`rujukan` meta (dulu `RUJUKAN` / `KOLOM_NAMA` frontend): province `nationId` → Nation (`id` / `note`); city `provinceId` →
Province (`id` / `note`); district `cityId` → City (`id` / `note`); czone `groupOf` → CZone (`code` / `description`);
accumulation `accumulation` → Accumulated Type (`id` / `accumulationType`) dan `provinceId` → Province (`id` / `note`).
Saran rujukan dilayani modul pemilik layar sendiri — pemegang menu City tidak butuh menu Province (nol `ruteDipinjam`).

## Yang dikerjakan (c3)

**Inti — WAJIB ditinjau tim inti:**

- `inti/backend/master/` — mesin bersama, dipindah dari `modul/masterdata/backend` (models / repository / services /
  handlers; isi sama kecuali): `Rujukan.Master`, `TabelMaster.KolomNama`, `TabelMaster.Saran()`; `services.Rujukan`;
  `handlers.Pasang(mux, prefix, kunci, …)` menggantikan `DaftarkanRute` bertabel; `master.Pendaftaran(nama, prefix,
  kunci, migrasi)` + tipe `Modul`. Uji: `TestMetaMaster`, `TestRujukanMaster` + uji lama.
- `inti/backend/migrations/912`–`919` (+ `_down`): baris menu delapan modul. `920_m_nav_menu_masterdata_pensiun`
  (+ `_down`): salin hak 'masterdata' → delapan KODE, buang hak, buang baris. `911_m_nav_menu_masterdata` (+ `_down`)
  DIBUANG.
- Penjaga: `menu_test.go` (`modulLuarKorpus` +8 −masterdata; `labelTampilDisetujui` +8, folder "Master <nama>";
  `langkahPensiunMenu` + skema tiruan mengenal `DELETE … M_NAV_MENU WHERE KODE` untuk KODE yang dipensiunkan saja dan
  pernyataan `M_LOGIN_GO_MENU`), `rentang_test.go` (`jatahModul.tanpaSlot`; 920 diizinkan; daftar urutan pelari),
  `modulmd_test.go` (`Slot menu` — sah hanya untuk modul luar korpus).
- `inti/backend/daftar/` bangkitan (`go generate`): −masterdata, +8. `inti/docs/STRUKTUR-TABEL-INTI.md`,
  `PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` (urutan deploy), komentar `inti/backend/migrasi`.

**Modul:** delapan `modul/master<nama>/backend/modul.go` + `MODUL.md` (frontend-nya sesi 9d). `masterprovince`:
`migrations/880`–`882` (+ `_down`, nama tetap), `docs/STRUKTUR-TABEL-MASTER-DATA.md`, `docs/issues/01`, `02`, ini,
pernyataan warisan `OBJECTITEMTYPE`, bab urutan deploy dan penomoran ulang. `modul/masterdata` dihapus (frontend-nya
oleh 9d). `nbfacin`: rujukan pemilik di `MODUL.md` dan `STRUKTUR-TABEL-NB-FACIN.md` (kode tidak berubah — ia membaca
tabelnya lewat SQL).

## Menyimpang dari pesan penerus

- ⚠️ Pesan 9d: *"Berkas 911 (inti) dan 996 sudah tercatat, jadi tetap ada. Tentukan pemilik 996 yang lolos penjaga."*
  Keduanya **DIBUANG**, bukan diberi pemilik. Alasan `[terverifikasi dari kode]`: pelari hanya menjalankan berkas yang
  ADA — maju (`migrasi.Jalankan` mengulang `Daftar` berkas) maupun mundur (`Bongkar`, sama) — jadi nama di `T_MIGRASI`
  tanpa berkas diabaikan, tidak digagalkan. 996 hanya menyalakan baris 'masterdata' dan tidak dapat dimiliki modul lain
  (penjaga slot: slot hanya menyentuh baris modulnya sendiri); 911 hanya membuat baris yang 920 buang. Di skema baru
  hasilnya sama: baris 'masterdata' tidak pernah lahir. Bila tim inti tetap ingin keduanya ada: 911 kembali ke
  `modulLuarKorpus`-pensiun dan 996 butuh pengecualian penjaga slot.

## Untuk work owner / DBA di DEV

1. `MODUL_AKTIF`: bila menyebut nama satu per satu, tambah kedelapan nama dan buang `masterdata`.
2. Pasang biner + bundel frontend BERSAMAAN.
3. `-migrate`: hanya 912–920 yang baru (880–882 sudah tercatat). `T_MIGRASI` TIDAK diubah.
4. Akun yang tidak memegang Master Data sebelum 920: beri menu lewat Kelola User.

⛔ Seluruhnya perubahan basis data / deploy: work owner / DBA, dengan persetujuan — bukan agent.

## Tinjauan dua sumbu (subagen Standards + Spec)

Diperbaiki: `920_down` mengembalikan baris 'masterdata' ber-`DIMIGRASI '0'` (bukan '1': modulnya sudah tidak ada, baris
hidup = menu mati); penjaga skema tiruan melewatkan pernyataan `M_LOGIN_GO_MENU` HANYA di langkah pensiun (920) — di
langkah menu lain ia "bentuk tidak dikenal" (uji mutasi: merah); uji gerbang `cmd/api` `TestGerbangMenuMaster` (pemegang
City: rute City + saran Province miliknya dilayani, rute Province 403, `/api/masterdata` 404); `PANDUAN-TIM-PER-MODUL.md`
bab 5 menyebut pilihan tanpa slot. Tidak diubah (penilaian, untuk tim inti): penjaga kini memuat nama modul (pola lama
`modulLuarKorpus` sejak 906); delapan KODE tertulis di 920, 920_down, `modulLuarKorpus`, `labelTampilDisetujui`, daftar
urutan pelari; mesin domain di `inti/` tanpa ADR (keputusan work owner "Pindah ke inti" tercatat di sini saja);
`Pendaftaran.NamaLama` tidak dipakai (satu nama lama → delapan nama baru); kontrak mesin inti tinggal di folder
`masterprovince`.

## Terbuka

- `920_down` mengembalikan hak 'masterdata' kepada pemegang SALAH SATU delapan menu — `[dugaan]` terdekat, rugi-informasi;
  pemegang asli tidak tersimpan sesudah 920. Di skema yang tak pernah menjalankan 911, barisnya tersisipkan mati.
- ADR untuk mesin master di `inti/` (tim inti).
- Uji perakit frontend (`daftar.menuTabel.test.ts`, `inti/frontend/uji/menuBersih.ts`) harus mengenal INSERT ber-DIMIGRASI
  '1' di 912–919 dan pernyataan 920 — sesi 9d.
