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

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

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
| `M_PRODUCT_LIFE` | tabel warisan POOLDATA berisi produk JSON seperti Pega; sejak 02-10-2026 (OQ-MPNL-01 flat, K5) hanya CADANGAN dan sumber alat pindah `backend/alat/pindahflat` - dibaca, tidak dibuat, tidak diubah migrasi mana pun |
| `M_PRODUCTINWARD_LIFE` | tabel warisan POOLDATA berisi sisi inward JSON seperti Pega; sejak 02-10-2026 hanya CADANGAN dan sumber alat pindah - dibaca, tidak dibuat, tidak diubah migrasi mana pun |
| `M_ATTACHMENTPRODUCTNAME` | tabel warisan POOLDATA tempat Master Product Name Life merekam lampiran tanpa membuatnya (P5) |

### Kaskade ON DELETE CASCADE

Kaskade HANYA pada berkas migrasi modul ini yang berawalan di bawah; berkas lain modul ini (induk `140_`, slot menu `960_`)
tanpa `ON DELETE CASCADE` (`TestKaskadeHanyaPadaRelasiTerdaftar`). Mendaftarkan yang baru menuntut bukti.

| Awalan berkas | Relasi |
| --- | --- |
| `141_` | `M_PRODUCTNAME_LIFE_LIEN` → `M_PRODUCTNAME_LIFE`: baris `LienClause` hidup DI DALAM halaman produk (grid b12201); simpan menulis ulang seluruh anak satu produk |
| `142_` | `M_PRODUCTNAME_LIFE_DOCCLAIM` → induk: `DocumentClaim` (grid b14601), alasan sama |
| `143_` | `M_PRODUCTNAME_LIFE_PLAN` → induk: `PlanList` (grid b31557), alasan sama |
| `144_` | `M_PRODUCTNAME_LIFE_FINUW` → induk: `FinancialUnderwritingList` (grid b37148), alasan sama |
| `145_` | `M_PRODUCTNAME_LIFE_UWLIMIT` → induk: `UnderwritingLimitList` (grid b42075), alasan sama |
| `146_` | `M_PRODUCTNAME_LIFE_OUTWARD` → induk: `OutwardList` (`GetReinsTypeOR_Life`), alasan sama |
| `147_` | `M_PRODUCTNAME_LIFE_COMMENT` → induk: `CommentList` (`AddCommentList_Act`), alasan sama |

### Pesan verbatim yang bukan nama orang

Konstanta teks yang cocok dengan pola nama orang tetapi BUKAN nama orang (`TestNolNamaOrangDiKode`).
Nilainya DIBACA dari sumber konstantanya, tidak diketik ulang.

| Paket | Konstanta | Alasan |
| --- | --- | --- |
| `backend/services` | `labelPolicyHolder` | label medan VERBATIM `Policy Holder` (`InboxProductName.xml` b17097) - dipakai kalimat penolakan, bukan nama orang. |

## Kontrak modul

| Hal | Isi |
| --- | --- |
| Rute | `GET /produk`, `GET /produk/{id}`, `POST /produk`, `PUT /produk/{id}`, `POST /produk/generate`, `GET /master/{jenis}`, `GET /master-plan`, `GET /rate` (503, OQ-MPNL-03), `GET`/`POST /produk/{id}/lampiran`, `POST …/{lid}/ulangi`, `GET …/{lid}/unduh`, `GET …/unduh-semua`, `GET …/{lid}/office` (503, OQ-MPNL-11), `DELETE …/{lid}` — rincian `docs/PARITAS-LAYAR-DAN-AKSI.md` §9 *(ralat 01-10-2026, K1 keputusan work owner 01-10-2026: `GET /rate` dan `GET /master/ri-rate` = 200, view rate baca saja)* |
| Tabel ditulis | ⭐ sejak 02-10-2026 (K5, tiket 01): tabel flat `M_PRODUCTNAME_LIFE` + `M_PRODUCTNAME_LIFE_{LIEN,DOCCLAIM,PLAN,FINUW,UWLIMIT,OUTWARD,COMMENT}` (migrasi 140–147), `M_ATTACHMENTPRODUCTNAME`, `T_STORAGE_IMAGE` (pelaksana stub), outbox bersama `T_LOG_SERVICE_RNM` (`MODUL = 'MASTERPRODUCTNAMELIFE'`). *(Dulu: `M_PRODUCT_LIFE` `JSONDATA` + `RIRISKID`/`RIRISK`, `M_PRODUCTINWARD_LIFE` `JSONDATA` — kini hanya dibaca alat pindah `backend/alat/pindahflat` dan pemeriksa identitas.)* |
| Tabel dibaca saja | `AGENT`, `CLIENT`, `CURRENCY`, `RIRISK_LIFE_SUMMARY`, `CAUSEOFLOSS_LIFE`, `PRODUCT_TYPE_LIFE`, `T_FOLDER_IMAGE`, `TREATYCONTRACT_LIFE`, `TREATYYEAR_LIFE` |
| Pembaca hilir | ⚠️ ketiga view `PRODUCT_LIFE`, `PRODUCTINWARD_LIFE`, `DOCUMENTCLAIM_LIFE` **tidak** dibangun ulang (K7, keputusan work owner 02-10-2026) — tetap membaca `JSONDATA` tabel warisan yang berhenti diperbarui sesudah peralihan. Claim Life `ambangproduk.go` membaca `PRODUCTINWARD_LIFE`: **OQ-FLAT-04** |
| Prosedur | `PEGA_M_PRODUCT_LIFE`, `PEGA_M_PRODUCT_INWARD_LIFE` **tidak** dipanggil (isinya ditiru di Go: satu baris induk + tujuh anak, satu transaksi, nol `COMMIT` di teks SQL) |

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/masterproductnamelife/...
go test -tags db -p 1 ./modul/masterproductnamelife/...    # tanpa ORACLE_DSN: uji db SKIP dengan pesan
npx vitest run modul/masterproductnamelife
```

## Keputusan OQ work owner 01-10-2026 (`PROMPT-LANJUTAN-TIGA-MODUL-LIFE-KEPUTUSAN-OQ.md`)

| Butir | Keadaan | Bukti |
| --- | --- | --- |
| K1 OQ-MPNL-03 R/I Rate dan View Rate | view `RATE_LIFE_SUMMARY` (`Choose R/I Rate`) dan `RATE_LIFE` (`View Rate`) dibaca **saja**, kolom RD saja; baris `PLAN LIST` baru dapat diberi R/I Rate (pilihan baru wajib ada di view, nama dari master). DEV baca-saja: `GET /master/ri-rate` 200 (346 baris), `GET /rate` 200 (59 baris, 0,35 detik), nol tulisan | `repository/mpnl_master.go`, uji `TestMPNLRateDibacaKolomRDSaja`, `TestMPNLSetiapSQLMasterAdalahSelect`, `TestPlanRIRateBaruDariViewRingkasan` |
| §2 delapan OQ | OQ-MPNL-01 (JSON seperti Pega; tiket 01 tetap ditangguhkan), 05, 06, 07, 10, 11, 13, 14 **ditutup** dengan bawaan yang dibangun; konfirmasi menyusul OQ-MPNL-05 (pemilik ekspor Pega). OQ terbuka: nol | `docs/OQ-MASTER-PRODUCT-NAME-LIFE.md` bab keputusan 01-10-2026; tiket 01, 06 |
