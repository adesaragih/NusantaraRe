# Modul `masterproductnamelife` — Master Product Name Life

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di
sini (struktur tim satu folder per modul, keputusan work owner 30-09-2026). Commit Anda menyentuh
folder ini saja; berkas di luarnya milik tim inti (`.github/CODEOWNERS`).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `masterproductnamelife` |
| Folder korpus | `Master Product Name Life` |
| GROUPMENU | `MASTER` |
| Pemilik | `@PEMILIK-MASTERPRODUCTNAMELIFE` |
| Status | dimigrasi |
| Rentang migrasi | `140-179` |
| Slot menu | `960-961` |
| Prefix rute API | `/api/master-product-name-life` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda; akun sebenarnya diisi work owner di `.github/CODEOWNERS`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | spec, tiket (`issues/`), grilling, catatan — dulu `.scratch/master-product-name-life/` (dipindah dengan `git mv`, isi byte-identik) |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `migrations/` `modul.go` — paket Go `nusantarare/modul/masterproductnamelife/backend/...` |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `bentuk.ts` `masterproductnamelife.css` `menu.ts` `rute.tsx` dan berkas `*.test.ts` |

## Migrasi

Rentang `140-179` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `960-961` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026) saat modul mendapat layar pertamanya, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nomor selalu tiga digit.

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
| `M_PRODUCT_LIFE` | tabel warisan POOLDATA yang Master Product Name Life tulis dan baca sebagai JSON seperti Pega tanpa membuatnya (P1, OQ-MPNL-01) |
| `M_PRODUCTINWARD_LIFE` | tabel warisan POOLDATA yang Master Product Name Life tulis dan baca sebagai JSON seperti Pega tanpa membuatnya (P1, P4) |
| `M_ATTACHMENTPRODUCTNAME` | tabel warisan POOLDATA tempat Master Product Name Life merekam lampiran tanpa membuatnya (P5) |

### Pesan verbatim yang bukan nama orang

Konstanta teks yang cocok dengan pola nama orang tetapi BUKAN nama orang (`TestNolNamaOrangDiKode`).
Nilainya DIBACA dari sumber konstantanya, tidak diketik ulang.

| Paket | Konstanta | Alasan |
| --- | --- | --- |
| `backend/services` | `labelPolicyHolder` | label medan VERBATIM `Policy Holder` (`InboxProductName.xml` b17097) - dipakai kalimat penolakan, bukan nama orang. |

## Kontrak modul

| Hal | Isi |
| --- | --- |
| Rute | `GET /produk`, `GET /produk/{id}`, `POST /produk`, `PUT /produk/{id}`, `POST /produk/generate`, `GET /master/{jenis}`, `GET /master-plan`, `GET /rate` (503, OQ-MPNL-03), `GET`/`POST /produk/{id}/lampiran`, `POST …/{lid}/ulangi`, `GET …/{lid}/unduh`, `GET …/unduh-semua`, `GET …/{lid}/office` (503, OQ-MPNL-11), `DELETE …/{lid}` — rincian `docs/PARITAS-LAYAR-DAN-AKSI.md` §9 |
| Tabel ditulis | `M_PRODUCT_LIFE` (`JSONDATA` + kolom datar `RIRISKID`, `RIRISK` — katalog DEV), `M_PRODUCTINWARD_LIFE` (`JSONDATA`), `M_ATTACHMENTPRODUCTNAME`, `T_STORAGE_IMAGE` (pelaksana stub), outbox bersama `T_LOG_SERVICE_RNM` (`MODUL = 'MASTERPRODUCTNAMELIFE'`) |
| Tabel dibaca saja | `AGENT`, `CLIENT`, `CURRENCY`, `RIRISK_LIFE_SUMMARY`, `CAUSEOFLOSS_LIFE`, `PRODUCT_TYPE_LIFE`, `T_FOLDER_IMAGE`, `TREATYCONTRACT_LIFE`, `TREATYYEAR_LIFE` |
| Pembaca hilir | view `PRODUCT_LIFE`, `PRODUCTINWARD_LIFE`, `DOCUMENTCLAIM_LIFE` (Claim Life) — setiap kunci yang dibacanya dijamin ada di `JSONDATA` (`docs/dba-view-produk-life.md`) |
| Prosedur | `PEGA_M_PRODUCT_LIFE`, `PEGA_M_PRODUCT_INWARD_LIFE` **tidak** dipanggil (isinya ditiru di Go, satu transaksi, nol `COMMIT` di teks SQL) |

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/masterproductnamelife/...
go test -tags db -p 1 ./modul/masterproductnamelife/...    # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx vitest run modul/masterproductnamelife
```
