---
status: accepted
label: DECIDED
---

# Aturan penguraian teks ke angka, dan perlakuan atas kurs yang tidak ada

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Berlaku hanya untuk keadaan sistem pada kedua ekspor itu.

Nilai uang di sistem lama tersimpan sebagai teks di dalam `DATA_JSON`, dan bahkan view `CLAIMXOL` mengeluarkannya sebagai `varchar2`. Aturan penguraiannya ditetapkan **tertulis dan di muka**, sebelum uji coba migrasi mana pun dijalankan.

## Aturan penguraian

| Hal | Perlakuan |
|---|---|
| Pemisah ribuan | ditolak bila ada; angka yang sah tidak memakainya |
| Tanda desimal | titik saja; koma sebagai desimal ditolak dan masuk daftar |
| Spasi di awal/akhir | dipangkas, tidak dianggap galat |
| Notasi ilmiah (`1.2E5`) | ditolak dan masuk daftar |
| Posisi tanda minus | hanya di depan; `123-` ditolak |
| Presisi melebihi ADR-0003 | ditolak dan masuk daftar, **tidak dibulatkan** |

## Tiga nilai kosong yang tidak boleh disamakan

Ini yang paling mudah terlewat, dan kerusakannya tidak dapat dipulihkan setelah migrasi:

| Bentuk | Arti yang dipegang | Perlakuan |
|---|---|---|
| `""` (string kosong) | **tidak pernah diisi** | dipetakan ke "belum diisi", bukan nol |
| `null` / kunci tidak ada di JSON | **tidak ada** | dipetakan ke null |
| `"0"` | **diisi nol** | dipetakan ke nol |

Menyamakan ketiganya menghancurkan informasi yang tidak dapat dipulihkan. Apakah ketiga makna itu benar-benar dibedakan oleh sistem lama **belum terbukti dari XML** — bila ternyata tidak dapat dibedakan maknanya, itu pertanyaan tersendiri yang dicatat di `_selesai/OPEN-QUESTIONS.md`, bukan diputuskan sepihak di sini.

## Kurs yang tidak ada

Konversi tanpa kurs **tidak menghasilkan angka**. Nilai IDR dibiarkan kosong dan ditandai "menunggu kurs". Sistem lama memilih `RETURN 1` — lihat `FINDING-006` — dan itu tidak diwarisi.

## Consequences

Ambang penghentian migrasi **berbasis nilai, bukan cacah baris**. Satu baris bermasalah senilai lima miliar lebih berat daripada lima ratus baris senilai seratus ribu.

- Setiap baris yang tidak dapat diurai diselesaikan satu per satu, berapa pun jumlahnya. Tidak ada baris yang dibuang karena "cuma sedikit".
- Migrasi dihentikan bila **nilai total yang belum terselesaikan** melewati ambang yang disepakati akuntansi, atau bila **ada satu baris** yang nilainya melewati ambang material.

Kedua angka ambang itu belum ditetapkan dan menunggu akuntansi — `ASK-AKUNTANSI.md` pertanyaan 4.

## Tambahan 18 September 2026 — aturan diturunkan dari kotoran yang benar-benar ada

Daftar aturan di atas semula disusun dari kemungkinan umum. Sapuan seluruh 279 berkas menggantinya dengan yang terbukti.

| | Jumlah |
|---|---|
| Seluruh pemanggilan `toDecimal` | **127** |
| Didahului pembersihan `replaceAll` | **9** — seluruhnya atas `.Deductible2` |
| **Tanpa pembersihan apa pun** | **118** |

**Koma sebagai pemisah desimal bukan hipotesis.** `@toDecimal(@replaceAll(.Deductible2,",","."))` di `Activity\CountLossAllocation_act.xml` membuktikan seseorang pernah menemuinya di produksi dan menambalnya. Tidak ada yang menulis pembersihan untuk masalah yang tidak pernah terjadi.

**Pembersihan lain yang ditemukan, dan sasarannya bukan angka**:

| Ekspresi | Sasaran |
|---|---|
| `@replaceAll(pyWorkPage.ClaimData.InsuredName,",","")` | menghapus koma dari **nama**, bukan dari angka — diduga demi keamanan muatan |
| `@replaceAll(.AcceptedNo,".","")` | menghapus titik dari **nomor dokumen** |
| `@pxReplaceAllViaRegex(.NoAccount,"[^0-9]","")` | menyisakan hanya digit pada **nomor rekening** |

Ketiganya memperlihatkan pola yang sama: pembersihan dilakukan **per tempat, saat masalahnya muncul**, bukan sebagai aturan yang berlaku menyeluruh.

**Yang paling perlu diperhatikan**: dua dari 118 pemanggilan tanpa pembersihan adalah `@toDecimal(TempKasir.CARI20)` dan `@toDecimal(TempKasir.CARI21)` — muatan menuju sistem kasir, yaitu nilai yang masuk jurnal menurut `ASK-AKUNTANSI.md` bagian 0.2.

Aturan penguraian di sistem baru karena itu **berlaku di satu tempat untuk semua nilai**, bukan ditambahkan per ekspresi saat kegagalannya ditemukan.

## Tambahan — keadaan "menunggu kurs" menular

Nilai yang belum terkonversi menandai **seluruh nilai turunannya**. Apa pun yang dihitung dari nilai bertanda "menunggu kurs" ikut bertanda sama.

Nilai bertanda itu **tidak boleh** masuk penjumlahan, tidak boleh dibandingkan terhadap ambang persetujuan, dan tidak boleh muncul di laporan — baik sebagai nol maupun sebagai nilai apa adanya.

Tanpa aturan penularan ini, keputusan "gagalkan, jangan mengarang angka" hanya berlaku satu tingkat: nilai pertama ditandai, lalu penjumlahan di tingkat berikutnya memperlakukannya sebagai nol dan cacatnya kembali muncul dalam bentuk yang lebih sulit dilihat — persis pola `RETURN 1` yang sedang kita tinggalkan (`FINDING-006`).

Ini melengkapi `ADR-0019` lapis 2: keadaan di tingkat baris mencatat *bahwa* perhitungan gagal; aturan penularan ini menentukan *apa yang terjadi pada baris-baris sesudahnya*.
