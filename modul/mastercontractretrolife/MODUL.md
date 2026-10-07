# Modul `mastercontractretrolife` — Master Contract Retro Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di
sini (struktur tim satu folder per modul, keputusan work owner 30-09-2026). Commit Anda menyentuh
folder ini saja; berkas di luarnya milik tim inti (`.github/CODEOWNERS`).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `mastercontractretrolife` |
| Folder korpus | `Master Contract Retro Life` |
| Nama tampilan | `Contract Retro Life` — keputusan work owner 03-10-2026 (kata "Master" dihapus, nama tampilan saja): label menu `M_NAV_MENU.LABEL` (slot menu 959), judul halaman, kartu Beranda. Kode modul, folder, rute API, dan `MODUL_AKTIF` tetap. |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-MASTERCONTRACTRETROLIFE` |
| Status | dimigrasi |
| Rentang migrasi | `100-139` |
| Slot menu | `958-959` |
| Prefix rute API | `/api/master-contract-retro-life` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `modul.go` — paket Go `nusantarare/modul/mastercontractretrolife/backend/...` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `tampilan.ts` `mcrl.css` `menu.ts` `rute.tsx` dan berkas `*.test.ts` |
| `docs/` | spec, tiket (`issues/`), grilling, PARITAS, RALAT, OQ, STRUKTUR — dulu `.scratch/master-contract-retro-life/` |

## Rute API

Prefix `/api/master-contract-retro-life` (`backend/handlers/rute_mcrl.go` `Prefix`). Setiap rute menuntut identitas
pelaku (401 tanpa), menjawab 503 bila Oracle tidak dikonfigurasi, dan galatnya berbadan `{"galat": "..."}`.
Nol kontrak lintas modul: master rujukan dibaca langsung dari tabelnya.

| Metode dan jalur | Layar / tombol Pega |
| --- | --- |
| `GET /tahun` | grid halaman awal `MASTER CONTRACT RETRO LIFE` |
| `POST /tahun` · `PUT /tahun/{id}` | `Add` (label sel `End Period`) / `Edit` → `Save` |
| `GET /tahun/{id}/kontrak` · `POST /tahun/{id}/kontrak` · `PUT /kontrak/{id}` | `ReinsType` → panel `Reins Type`; `Add`/`Edit` → `Save` |
| `GET /kontrak/{id}/reinsurer` · `POST /kontrak/{id}/reinsurer` · `PUT /reinsurer/{id}` | `Reinsurer List` (+ `Total Share -->>`) |
| `GET /reinsurer/{id}/security` · `POST /reinsurer/{id}/security` · `PUT /security/{id}` | `Security Reinsurer` (+ eksposur) |
| `GET /kontrak/{id}/business` · `POST /kontrak/{id}/business` · `PUT /business/{id}` | `Business List` |
| `GET /business/{id}/salin-semua` · `POST /business/{id}/salin-semua` | `Copy to all Reinstype` — pratinjau, lalu konfirmasi berbadan `{"sasaran": [...]}` |
| `GET /{kontrak\|reinsurer\|security\|business}/{id}/dampak-hapus` · `DELETE /{…}/{id}` | `Delete` — popup berdampak, lalu hapus berbadan `{"dampak": {...}}`; nol rute hapus tahun |
| `GET /jenis-reasuransi` | dropdown `REINS TYPE` (master `REINSURANCETYPE` `.Flag = 1`) |
| `GET /master-reinsurer?cari=` | autocomplete `REINSURER NAME` / `SECURITY REINSURER NAME` |
| `GET /master-business?cari=` | autocomplete `BUSINESS NAME` |
| `GET /ringkasan-rate?cari=` · `GET /rate?idusedby=` | autocomplete `R/I RATE` (view `RATE_LIFE_SUMMARY`) · `View Rate` (`Rate List`, view `RATE_LIFE`) — baca saja sejak K1 keputusan work owner 01-10-2026 (OQ-MCRL-13) *(RALAT 07-10-2026: ringkasan rate kini tabel `M_RATE_LIFE_SUMMARY` berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID - keputusan work owner 07-10-2026, `modul/riratelife/MODUL.md` RALAT R6; modul ini membacanya `SELECT ID, USEDBY`, tetap baca-saja)* |
| `GET /laporan/total-share-bukan-100?tahun=` | tanpa layar (tiket 11, OQ-MCRL-07) |

## Migrasi

Rentang `100-139` **tetap kosong**: K1 (`docs/RALAT-DEV-30-09-2026.md`, preseden tco4) — modul ini menulis
dan membaca lima tabel warisan `POOLDATA`, nol tabel baru, nol DDL (`TestMCRLNolMigrasiDiRentang`,
`TestMCRLNolDDL`). Peta tabelnya: `docs/STRUKTUR-TABEL-MASTER-CONTRACT-RETRO-LIFE.md`.

Slot menu `958-959`: `backend/migrations/958_menu_mastercontractretrolife.sql` (+ `_down`) — satu
`UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1'` baris modul ini, nol `INSERT` (menu datar 30-09-2026,
`APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md` bab 6). ⛔ **`-migrate` dijalankan work owner**, bukan sesi
pengembang; sampai 958 dijalankan, tombol menu tetap "belum dimigrasi" di basis data yang sudah berjalan.

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per
butir. Judul yang tidak ada berarti modul ini tidak menyatakan apa pun untuk jenis itu. Nilai di dalam
`` ` `` dibaca apa adanya.

### Tabel warisan: dibaca, tidak dibuat

Tabel yang dokumen STRUKTUR modul ini gambarkan tetapi SENGAJA tidak dibuat migrasi mana pun
(`TestKolomDDLCocokDenganStruktur`, `TestTabelBukanMilikKitaTidakDibuat`). Mencabut satu baris =
kepemilikan tabel berpindah — keputusan work owner.

| Tabel | Alasan |
| --- | --- |
| `TREATYYEAR_LIFE` | tabel warisan POOLDATA yang Master Contract Retro Life tulis dan baca tanpa membuatnya (K1, preseden tco4) |
| `TREATYCONTRACT_LIFE` | tabel warisan POOLDATA yang Master Contract Retro Life tulis dan baca tanpa membuatnya (K1, preseden tco4) |
| `TREATYREINSURER_LIFE` | tabel warisan POOLDATA yang Master Contract Retro Life tulis dan baca tanpa membuatnya (K1, preseden tco4) |
| `TREATYSECURITYREINSURER_LIFE` | tabel warisan POOLDATA yang Master Contract Retro Life tulis dan baca tanpa membuatnya (K1, preseden tco4) |
| `TREATYBUSINESS_LIFE` | tabel warisan POOLDATA yang Master Contract Retro Life tulis dan baca tanpa membuatnya (K1, preseden tco4) |

## Gelombang 2 brief rumpun Life (01-10-2026)

| Butir | Keadaan | Bukti |
| --- | --- | --- |
| A1 OQ-MCRL-13 rate | izin work owner **belum** tercatat — `GET /ringkasan-rate`, `GET /rate` tetap **503** berkalimat *(ralat 01-10-2026: izin K1 tercatat; kedua rute 200 — lihat bab Keputusan OQ 01-10-2026)* | `0f71e05`, `docs/OQ-MASTER-CONTRACT-RETRO-LIFE.md` |
| A2 `bNNN` dihitung ulang | PARITAS dan RALAT dapat diulang dengan perintah brief §1.2 (`sed` pemecah tag lalu `grep -n`) | `8eb57c9` |
| A3 OQ-MCRL-07 | tetap rute API tanpa layar (`GET /laporan/total-share-bukan-100`) | bab Rute API |
| A4 uji manual | **sebagian**: backend + `npm run dev` + seluruh rute baca kelima layar terhadap DEV, nol tulis; klik di peramban menunggu work owner | `docs/LAPORAN-UJI-MANUAL.md` |

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/mastercontractretrolife/...
go test -tags db -p 1 ./modul/mastercontractretrolife/...    # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx vitest run modul/mastercontractretrolife
```

`npx tsc --noEmit` memeriksa seluruh aplikasi sekaligus — tipe tidak dapat diperiksa per folder.
Sebelum pull request, jalankan gerbang lengkap (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` di akar repo, bab 8).

## Keputusan OQ work owner 01-10-2026 (`PROMPT-LANJUTAN-TIGA-MODUL-LIFE-KEPUTUSAN-OQ.md`)

| Butir | Keadaan | Bukti |
| --- | --- | --- |
| K1 OQ-MCRL-13 + OQ-MCRL-05 rate | view `RATE_LIFE_SUMMARY` (autocomplete `R/I RATE`) dan `RATE_LIFE` (`Rate List`) dibaca **saja**, kolom RD saja; business baru dapat disimpan (RIRATEID pilihan baru wajib ada di view ringkasan). DEV baca-saja: `GET /ringkasan-rate` 200 (100 saran), `GET /rate` 200 (1 dan 59 baris, ±0,35 detik), nol tulisan | `repository/mcrl_master.go`, uji `TestRateDibacaKolomRDSaja`, `TestPeriksaBacaSajaMenolakTulisanKeView`, `TestMCRLMasterHanyaDibacaSelect` *(RALAT 07-10-2026: ringkasan rate kini tabel `M_RATE_LIFE_SUMMARY` berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID - keputusan work owner 07-10-2026, `modul/riratelife/MODUL.md` RALAT R6; modul ini membacanya `SELECT ID, USEDBY`, tetap baca-saja)* |
| §2 dua belas OQ | OQ-MCRL-01 (gerbang tahun ikut Pega), 02, 03 (0–100 ditegakkan), 04, 06, 07 (laporan hanya API), 08–12, 14 **ditutup** dengan bawaan yang dibangun; konfirmasi menyusul OQ-MCRL-02 (DBA), OQ-MCRL-04 (pemilik ekspor Pega). OQ terbuka: nol | `docs/OQ-MASTER-CONTRACT-RETRO-LIFE.md` bab keputusan 01-10-2026; tiket 03, 07, 11 |

