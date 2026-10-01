# Modul `endorsementlife` — Endorsement Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di
sini (struktur tim satu folder per modul, keputusan work owner 30-09-2026). Commit Anda menyentuh
folder ini saja; berkas di luarnya milik tim inti (`.github/CODEOWNERS`).

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

`Pemilik` adalah penanda; akun sebenarnya diisi work owner di `.github/CODEOWNERS`.

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
| `M_LIFE_PREMIUM_DETAIL`, `M_LIFE_PREMIUM_SUMMARY` | warisan `POOLDATA` | dibaca (sumber warisan) dan ditulis (E2, persis `SaveMasterLPDet`/`InsertPLSummary`); setiap kueri berkunci ber-index (E4) |
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
| `POST /kasus/{id}/simpan` | `Save` b37202 → `SetPremi_EDM` — centang `.EdmBatal` / `DELETE ALL` (`{"pilih":[],"semua":false,"kecuali":[]}`), jurnal balik 32 kolom, rekap mata uang; simpan kedua 409 |
