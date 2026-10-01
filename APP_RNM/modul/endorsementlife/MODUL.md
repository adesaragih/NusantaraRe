# Modul `endorsementlife` — Endorsement Life

⚠️ **Kerangka — belum dimigrasi.** Folder ini dibuat struktur tim satu folder per modul (keputusan work
owner 30-09-2026) supaya pemilik, rentang migrasi, dan slot menu modul ini TETAP sejak awal — satu
modul, satu folder, satu pemilik. Belum ada kode: tanpa `backend/modul.go` modul ini tidak terdaftar
(daftar Go bangkitan `inti/backend/daftar`, `import.meta.glob` frontend), dan kelompoknya di sidebar
tetap "belum dimigrasi" (`M_NAV_MENU.DIMIGRASI = '0'`). Cara memulainya:
`APP_RNM/PANDUAN-TIM-PER-MODUL.md` (akar repo) bab 4.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `endorsementlife` |
| Folder korpus | `Endorsement Life` |
| GROUPMENU | `TREATY` |
| Pemilik | `@PEMILIK-ENDORSEMENTLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `480-519` |
| Slot menu | `976-977` |
| Prefix rute API | `/api/endorsement-life` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `modul.go` — paket Go `nusantarare/modul/endorsementlife/backend/...` |
| `frontend/` | `pages/` `labels.ts` `api.ts` `tampilan.ts` `endorsementlife.css` `menu.ts` `rute.tsx` dan berkas `*.test.ts` |
| `docs/` | spec, tiket (`issues/`), grilling, PARITAS, RALAT, OQ, STRUKTUR — dulu `.scratch/endorsement-life/` |

## Migrasi

| Nomor | Isi |
| --- | --- |
| `480` | tiket 00 (ralat E1): `T_PREMIUM_LIST.PROD_KE NUMBER(5) DEFAULT 1` + isi mundur, index `IDX_PL_NOPOLIS_PRODKE`, `IDX_PL_OLD_POLICY_NO`, `IDX_PLD_PL_NUMBER` — nol kolom baru |
| `481` | sequence `SEQ_WORK_EDM_LIFE` — pengenal kasus `EDMLF-<n>` |
| `482` | index unik berfungsi `UX_PL_EDM_TERBUKA` — satu kasus endorsement terbuka per polis (AC 3) |
| `976` | slot menu: `UPDATE M_NAV_MENU SET DIMIGRASI = '1'` baris `endorsementlife`, nol `INSERT` |

⛔ **Nol tabel baru** (spec §16, AC 55): endorsement adalah versi baru di tabel PremiumList Life. Bentuk
SQL slot menu: `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6. Nomor selalu tiga digit.

## Tabel yang disentuh

| Tabel | Pemilik | Perlakuan |
| --- | --- | --- |
| `T_PREMIUM_LIST`, `T_PREMIUM_LIST_DETAIL`, `T_PREMIUM_LIST_SPREADING`, `T_PREMIUM_LIST_SPREADING_RETRO`, `T_PREMIUM_LIST_SUMMARY`, `T_VIEW_SUGGEST` | PremiumList Life (051–056) | dibaca dan ditulis — kasus EDM = baris versi `T_PREMIUM_LIST` ber-`ID` `EDMLF-<n>` |
| `JSON_POLIS` | warisan `POOLDATA` | dibaca saja — sumber polis lama sistem lama |
| `M_LIFE_PREMIUM_DETAIL` | warisan `POOLDATA` | **dibaca** saja (sumber versi new business warisan, `_INDEX4`); penulisan E2 (`SaveMasterLPDet`) menunggu OQ-EDM-016 (RALAT R29) |
| `M_LIFE_PREMIUM_SUMMARY` | warisan `POOLDATA` | dibaca (rekap versi warisan) dan **ditulis** saat `Confirm` — isi prosedur `PEGA_M_LIFE_PREMIUM_SUMMARY` ditiru, prosedur tidak dipanggil (E2) |
| `T_LOG_SERVICE_RNM` | inti (outbox) | ditulis — kegagalan efek keluar `MODUL = ENDORSEMENTLIFE` (tiket 10) |
| `T_CLAIMLF_JEJAK` | inti (jejak) | ditulis — buat kasus, `Save`, `Add CSV Data`, keputusan |
| `M_LINK_SERVICE` | inti (layanan) | dibaca saat jalan — kunci `Production` / `convertJsonNusareToProduction`; alamat tidak pernah dikutip |
| `ARASAPAS.DETAIL_INVOICE` | Arasapas | dibaca saja, satu repository (gerbang ke-5) |

## Rute API

Prefix `/api/endorsement-life` (`backend/handlers/rute_edm.go` `Prefix`). Setiap rute menuntut identitas
pelaku (401 tanpa), menjawab 503 bila Oracle tidak dikonfigurasi, dan galatnya berbadan `{"galat": "..."}`.

| Metode dan jalur | Layar / tombol Pega |
| --- | --- |
| `GET /inbox?halaman=` | `InboxEndorsementLife` grid b8284 (RD `InboxEDMLife`) |
| `GET /kasus/{id}` | kepala `InputEDMLife` |
| `GET /kasus/{id}/peserta?halaman=` | grid peserta `InputEDMLife` b11899 / b17500 |
| `GET /kasus/{id}/peserta/{pid}` | rincian `PL_Detail_Sec` (`pyEditAction` `PL_DetailAction`) + `RetroDetailLife` |
| `GET /kasus/{id}/polis-lama?halaman=` | popup `View Old Policy` (`ViewOldPolicy_EDM`, `_QP`, `_TP`, `_TR`) |
| `POST /kelayakan` | `SetErrorBatalEndorsement_Act` — lima gerbang, tanpa tulis |
| `POST /kasus` | `Submit` b4226 → `MappingEDMLife` (kasus `EDMLF-<n>` + salinan versi berjalan + jejak) |
| `POST /kasus/{id}/unggah` | `Upload CSV` b8973 — tinjau CSV (multipart `berkas`): cacah, penolakan berbaris/berkolom, judul diabaikan; nol tulis |
| `POST /kasus/{id}/csv` | `Add CSV Data` b10405 → `SaveCSVEDMLife` — baris `New`, tanpa batas baris, satu penolakan menolak seluruh berkas (422 berdaftar pesan) |
| `POST /kasus/{id}/putuskan` | `Submit` b37494 / b38109 — `{"status":"1"\|"2"\|"7","comment":""}`: riwayat `T_VIEW_SUGGEST`; Confirm = nomor `<polis>/NN`, anti-dobel `(NO_POLIS, PROD_KE)`, versi resmi, rekap warisan; Decline = `Resolved-Rejected`; satu transaksi |
| `POST /kasus/{id}/simpan` | `Save` b37202 → `SetPremi_EDM` — centang `.EdmBatal` / `DELETE ALL` (`{"pilih":[],"semua":false,"kecuali":[]}`), jurnal balik 32 kolom, rekap mata uang; simpan kedua 409 |
