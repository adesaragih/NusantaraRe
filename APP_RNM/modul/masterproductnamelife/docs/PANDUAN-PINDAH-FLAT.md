# Panduan peralihan — Master Product Name Life dari JSON ke tabel flat

> Untuk: work owner, DBA, dan operator yang menjalankan peralihan. Disusun 02-10-2026 (brief
> `PROMPT-PINDAH-FLAT-MASTER-PRODUCT-NAME-LIFE.md` T6; keputusan work owner K1–K7 di tiket 01 bab bertanggal 02-10-2026).

## Apa yang berubah

| Sebelum | Sesudah |
| --- | --- |
| Produk disimpan sebagai dokumen JSON di `M_PRODUCT_LIFE.JSONDATA` dan `M_PRODUCTINWARD_LIFE.JSONDATA` | Produk disimpan di **`M_PRODUCTNAME_LIFE`** (satu baris per produk, sisi umum + sisi inward) + tujuh tabel anak `M_PRODUCTNAME_LIFE_{LIEN,DOCCLAIM,PLAN,FINUW,UWLIMIT,OUTWARD,COMMENT}` |
| Aplikasi menulis kedua tabel JSON | Aplikasi menulis tabel flat saja. Kedua tabel JSON **tidak disentuh** — tetap ada sebagai cadangan |
| Tiga view DEV (`PRODUCT_LIFE`, `PRODUCTINWARD_LIFE`, `DOCUMENTCLAIM_LIFE`) membaca produk terbaru | Ketiga view **tidak dibangun ulang** (keputusan K7): sesudah peralihan view itu menampilkan isi JSON terakhir, tidak lagi diperbarui |

⚠️ **Pembaca hilir yang terdampak:** Claim Life `repository/ambangproduk.go` membaca view `PRODUCTINWARD_LIFE`
(ambang `MAXEXPIREDCLAIM`/`MAXDATARECEIVE`). Produk baru atau produk yang diubah sesudah peralihan **tidak terlihat**
di sana sampai OQ-FLAT-04 diputuskan (Claim Life dialihkan membaca `M_PRODUCTNAME_LIFE`, atau view dibangun ulang DBA).
*(02-10-2026: OQ-FLAT-04 ditunda — keputusan work owner *"OQ-FLAT-04 :  BIARKAN SAJA, NANTI PAS DEVELOP BAGIAN ITU AKAN DIPEERBAIKI!"* Celah ini diterima sampai bagian
Claim Life itu dikembangkan; bukan syarat peralihan.)*

## Syarat sebelum mulai

1. **Penulisan Pega ke layar Product Name Life dihentikan** (OQ-FLAT-03 — siapa dan kapan). Selama Pega masih menulis
   JSON, isi tabel flat tertinggal. Layarnya: aplikasi Pega LAMA, menu MASTER → **Master Product Name Life**, tombol
   `Save` (`SaveProductName_Act`) — bukan menu bernama sama di aplikasi baru.
2. ✅ Keputusan work owner 02-10-2026 atas OQ-FLAT-07/08/09 sudah dibangun (`LAPORAN-MIGRASI-FLAT.md`, putaran 20:46:
   gagal 0). Bila uji kering pada hari peralihan menemukan jenis normalisasi atau kegagalan BARU, berhenti dan minta
   keputusan lagi.
3. Lingkungan sasaran bukan produksi Pega: `IS_PEGA_PROD` harus `false`. Alat menolak `-jalankan` bila `true`.
4. Cadangan basis data menurut prosedur DBA.

Konfigurasi alat sama dengan `cmd/api`: `ORACLE_DSN`, `ORACLE_SCHEMA` (skema berisi kedua tabel JSON dan tabel flat),
`IS_PEGA_PROD`. Perintah dijalankan dari folder `APP_RNM/`.

## Langkah

| # | Langkah | Perintah / pemeriksaan | Lolos bila |
| ---: | --- | --- | --- |
| 1 | Hentikan penulisan Pega | (OQ-FLAT-03) | tidak ada simpan produk di Pega sejak titik ini |
| 2 | Buat tabel flat | `go run ./cmd/api -migrate` — menjalankan migrasi 140–148 (dan migrasi lain yang belum tercatat `T_MIGRASI`). *DEV 02-10-2026: 140–147 tercatat 15:39; yang tertunda tinggal `148_m_productname_life_outward_kolom` — jalankan `-migrate` sekali lagi* | `T_MIGRASI` mencatat `140_m_productname_life` … `148_m_productname_life_outward_kolom`; sebelumnya `147_m_productname_life_comment` |
| 3 | Uji kering | `go run ./modul/masterproductnamelife/backend/alat/pindahflat -terima-normalisasi="koma desimal,nol depan"` (mode `-uji`, hanya SELECT) | baris `gagal: 0`; cacah per tabel dan daftar K3/K4 sama dengan `LAPORAN-MIGRASI-FLAT.md` atau selisihnya dijelaskan |
| 4 | Pindahkan | `go run ./modul/masterproductnamelife/backend/alat/pindahflat -jalankan -terima-normalisasi="koma desimal,nol depan"` — kedua jenis diputuskan work owner 02-10-2026 (OQ-FLAT-07); jenis lain tetap menahan | baris terakhir `ditulis: true`; satu transaksi yang lebih dulu mengunci tabel induk (gagal di mana pun = nol tulisan) |
| 5 | Periksa agregat | `SELECT COUNT(*)` tiap tabel flat = baris "baris yang (akan) ditulis" laporan; `SELECT COUNT(*) FROM M_PRODUCT_LIFE` = 196 (atau cacah sumber saat itu) — tabel JSON tidak berubah | semua cacah sama |
| 6 | Pakai aplikasi versi flat | deploy biner `cmd/api` dari commit yang memuat tabel flat | layar Product Name Life menampilkan produk; simpan menulis tabel flat |

**Aman diulang:** langkah 3 boleh diulang kapan saja. Langkah 4 mengunci tabel induk (aplikasi yang berjalan menunggu),
menghapus lalu mengisi ulang produk **bersumber JSON** saja: produk yang hanya ada di tabel flat (dibuat aplikasi) dibiarkan
utuh, dan putaran **ditolak** bila produk bersumber JSON di tabel flat isinya sudah berbeda (diubah aplikasi sesudah
peralihan) — tulisan baru tidak pernah ditimpa.

**Isi laporan alat** (agregat — ID produk dan nama kolom, tidak pernah nilainya): cacah sumber; produk tanpa inward;
inward ber-ID lain (dicari lewat `PRODUCTID`); baris per tabel; K3 (nilai tidak sah yang menjadi NULL); K4 (objek
`OutwardList` kosong yang dibuang); normalisasi teks per kolom dan jenisnya; medan tanpa kolom yang terisi; yang dibuang
D2 (kunci internal Pega); kegagalan.

## Jalur mundur

| Keadaan | Jalur |
| --- | --- |
| **Sebelum** aplikasi versi flat menulis apa pun | kembali ke versi aplikasi sebelumnya — ia membaca dan menulis kedua tabel JSON, yang tidak pernah disentuh (isinya utuh). Tabel flat boleh dibiarkan (tidak dibaca versi lama) atau dibuang DBA dengan berkas `backend/migrations/14x_*_down.sql` (urutan 147 → 140) lalu baris `140_…`–`147_…` dihapus dari `T_MIGRASI` |
| **Sesudah** ada tulisan baru di tabel flat | butuh ekspor flat → JSON — **belum dibangun** (OQ-FLAT-02). Tanpanya, tulisan sesudah peralihan tidak dapat dikembalikan ke JSON |

⛔ **`cmd/api -migrate-down` BUKAN jalur mundur di skema sungguhan** (ralat atas brief T6): ia dipagari
`PastikanSkemaUji` — hanya berjalan bila `ORACLE_SKEMA_UJI=true` dan nama skema tidak memuat `POOLDATA` — dan ia
membongkar SELURUH migrasi yang tercatat, bukan hanya modul ini.
