# Laporan uji manual — Master Product Name Life (baca-saja terhadap DEV)

> 01-10-2026, lanjutan 1 §3 (`PROMPT-LANJUTAN-MASTER-PRODUCT-NAME-LIFE-1.md`), sesudah L1 (`6fd539c`) dan L2–L7 (`6209638`).
> ⛔ **Nol penulisan ke DEV:** tidak satu pun `Save`, `Attach`, `Delete`, `Retry`, `Generate`, `Copy` ditekan. Nilai data, nama orang,
> dan nama berkas **tidak** disalin ke dokumen ini — yang dicatat: tombol, layar yang terbuka, dan cocok/tidaknya nilai layar dengan `JSONDATA`.

## Susunan

| Hal | Isi |
| --- | --- |
| Backend | biner `cmd/api` (HEAD `6209638`), env `.env` DEV dimuat ke prosesnya saja (nilai tidak dicetak), **`MODUL_AKTIF=masterproductnamelife`** (nol pekerja latar modul lain), `AUTH_STUB=true`, `HTTP_ADDR` port sendiri; tanpa `-migrate` |
| Frontend | `npx vite --port 5175 --strictPort`, `DEV_PROXY_TARGET` ke backend di atas. ⚠️ Port **5174** yang diminta brief sedang dipakai sesi lain — sesi itu tidak diganggu |
| Peramban | Chrome tanpa kepala lewat CDP; skrip menolak mengklik tombol penulis (daftar terlarang di skrip, dilempar sebagai galat bila disentuh) |
| Produk contoh | **`100003`** — satu-satunya produk DEV ber-lampiran (2), 4 baris `PLAN LIST` ber-R/I Rate, 7 komentar, 2 baris `OutwardList`. Dipilih dengan `GET` atas ke-196 produk (cacah saja) |

## Hasil

| Tombol (bNNN) | Layar yang terbuka | Hasil |
| --- | --- | --- |
| menu **Master Product Name Life** | grid `InboxProductName` mode daftar | 10 baris per halaman (`pyRDLPageSize`), tombol `Add` + `View` per baris + penomoran; total 196 produk = katalog DEV |
| **`View`** b74753 (baris `100003`) | form mode lihat | 16 medan dibandingkan dengan `GET …/produk/100003`: `Product Code`, `Product Name`, `Ceding`, `SOB`, `Deduction (%)`, `R/I Risk Name`, `Treaty Name`, `Treaty Number`, `Policy Holder`, `Insured`, `Begin Date`, `Expired Date`, `Ceding's Limit`, `Max Sum Insured`, `Max Sum Reasured`, `Currency` — **16 cocok, 0 beda**. 46 medan baca-saja (`ro = IsView=='true'`). Tombol tampil: `Close`, `Edit`, `Copy`, `Generate` (tanpa `Save`, tanpa `Choose*`), `PLAN LIST` `Add` + `View Rate`/`Delete` per baris (vis `ALWAYS` / `.RIRATE!=''`) |
| (grid bersarang) | — | `PLAN LIST` 4 baris = `JSONDATA` 4; komentar 7 baris = 7; `LIEN CLAUSE`, `DOCUMENT CLAIM`, `FINANCIAL UNDERWRITING`, `UNDERWRITING LIMIT` kosong = `JSONDATA` kosong |
| panel lampiran (wadah b64133) | daftar `GetAttachmentProdName_Sql` | 2 baris, status **Uploaded** keduanya (objek `T_STORAGE_IMAGE` ada); tautan nama berkas, `Delete` per baris, `Add attachment`, `Refresh`, `Download All`. **`View Office Online`** b69247 tampil untuk berkas `.pptx` dan **tidak** untuk `.jpg` — sesuai syarat `.pyFileMimeType` |
| **`Refresh`** b65223 | — | daftar dibaca ulang (`GET …/lampiran`) |
| **`View Rate`** b34067 | dialog `ViewRate` (`Outward List`, `Cancel`/`Submit`) | jawaban 503 berkalimat yang menyebut OQ-MPNL-03 tampil — L8 tidak diizinkan (register OQ tanpa izin bertanggal), tetap 503 |
| **`Edit`** b59443 | dialog `EditProductName_Confirm` | pertanyaan `Do you want to Edit the data?` + `Cancel` / `Edit`; sesudah `Edit`: `Save` tampil (tidak ditekan), keenam `Choose*`, `Choose R/I Rate`, ikon grid; medan baca-saja tinggal 1 (`Product Code`, R14) |

## Jaringan

22 permintaan ke `/api/` selama uji, **seluruhnya `GET`**: `/api/menu`, `/api/modul-aktif`, `/api/klaim-life` (sidebar), `…/produk`,
`…/produk/100003`, `…/produk/100003/lampiran`, `…/master/penyebab`, `…/rate`. **Nol** `POST`/`PUT`/`DELETE`.

## Temuan sampingan

- DEV memuat `PAYMENT = 5` (contoh produk `100023`) — kode di luar 1–4 `GenerateUpload_Act` b1141; dropdown menampilkannya sebagai nilai
  di luar daftar (tidak dibuang). Bukti tambahan untuk OQ-MPNL-05 (kode `Single`).
- Tidak teramati di DEV: unduh berkas lampiran lama — berkas Pega tinggal di penyimpanan asal, bukan di folder stub (OQ-MPNL-10); tidak diklik.
