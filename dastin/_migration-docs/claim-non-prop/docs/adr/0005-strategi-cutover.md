---
status: accepted
label: DECIDED
---

# Claim dan Komite satu unit cutover, dengan shadow-run sebagai bukti paritas

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Claim Non Prop dan Komite Claim Non Prop tidak dapat dipisahkan pada saat peralihan: keduanya menulis ke record yang sama dan saling mengunci (`pxAddChildWork` dari sisi Claim, penulisan balik ke `AdjustmentList` dari sisi Komite), sehingga memotong di antara keduanya berarti menciptakan penguncian lintas sistem atas baris Oracle yang sama. Kami menetapkan keduanya beralih sebagai satu unit; sebelum peralihan dijalankan shadow-run, yaitu sistem baru menghitung ulang klaim produksi yang sudah selesai dan hasilnya dibandingkan angka per angka dengan sistem lama; klaim yang sedang berjalan diselesaikan di sistem lama dan sistem baru hanya melayani klaim baru.

## Consequences

**Tidak ada satu pun bukti di XML yang mendukung maupun menolak strategi ini** — XML hanya membuktikan keterkopelannya, bukan cara memindahkannya. Ini murni keputusan kami.

Yang terbukti dari XML adalah titik koplingnya, dan itu tercatat lengkap di `BLUEPRINT.md` §8 sebagai Boundary Contract Draft: dua pemanggilan `pxAddChildWork`, kunci relasi `pxCoveredInsKeys`/`pxCoverInsKey`, dan `Adjustment.IndexObject` sebagai penunjuk posisi numerik ke `AdjustmentList` induk.

`Adjustment.IndexObject` adalah risiko tersendiri: ia menyimpan posisi dalam daftar, bukan kunci surrogate. Bila urutan `AdjustmentList` berubah, kaitan antara keputusan Komite dan Adjustment yang dimaksud menjadi salah. Sistem baru wajib mengganti ini dengan kunci yang stabil, dan migrasi data harus memetakan posisi lama ke kunci baru.

## Tambahan 18 September 2026 — prasyarat baseline shadow-run

Shadow-run mengandaikan baseline yang sehat. Dua hal dapat merusaknya, dan keduanya harus diselesaikan **sebelum** perbandingan dimulai:

1. **Baris ganda akibat hitung ulang yang tidak idempoten** (ADR-0011). Case terdampak dibersihkan atau dikeluarkan dari perbandingan, dan didaftar terpisah. Diukur lewat REQ-014.
2. **Riwayat persetujuan yang tidak dapat direkonstruksi** karena hanya tersimpan sebagai teks komentar (ADR-0009, FINDING-002 bagian 7). Case terdampak tidak dapat dibandingkan pada dimensi "siapa menyetujui". Diukur lewat REQ-013.

Kasus uji utama: **`CLMNP-975`**, lihat ADR-0004.

## Tambahan 18 September 2026 — komposisi sampel dan aturan kegagalan

### Komposisi sampel ditetapkan di muka

Sampel shadow-run **ditetapkan sekarang, bukan dipilih saat menjalankan**. Perbandingan atas 500 klaim rupiah satu layer akan lulus dengan mudah dan tidak membuktikan apa pun.

Sampel wajib memuat, dengan proporsi minimum terhadap keseluruhan sampel:

| Kelompok | Minimum |
|---|---|
| Klaim multi-mata-uang (lebih dari satu `Currency` pada satu klaim) | 15% |
| Klaim multi-layer (lebih dari satu baris `SpreadingRisk` selain `"UR"`) | 20% |
| Klaim dengan reinstatement (`CNPReinstatement` bukan nol) | 10% |
| Klaim dengan lebih dari satu Adjustment | 15% |
| Klaim yang melibatkan baris `"UR"` | 20% |
| Kedelapan klaim bertambalan | **100% — seluruhnya wajib ikut** |

Kelompok boleh bertumpang tindih; yang tidak boleh adalah ada kelompok yang kosong. Bila populasi produksi tidak menyediakan cukup kasus untuk satu kelompok, itu sendiri temuan dan dicatat, bukan diam-diam diturunkan ambangnya.

### Aturan kegagalan

| Jenis nilai | Bila selisih melampaui ambang |
|---|---|
| **Nilai yang masuk jurnal akuntansi** | **Cutover berhenti.** Tanpa negosiasi, tanpa daftar pengecualian. |
| **Nilai antara** (alokasi per layer sebelum pembulatan akhir) | Boleh masuk daftar pengecualian, tetapi **setiap baris wajib punya penjelasan dan persetujuan bernama**. Pengecualian tanpa penjelasan tidak dihitung sebagai pengecualian — ia dihitung sebagai kegagalan. |

Aturan ini ditetapkan **sebelum** ada pihak yang punya kepentingan atas hasilnya. Itu alasan ia ditulis sekarang dan bukan nanti.

Angka toleransi relatif untuk nilai antara masih menunggu akuntansi (lihat ADR-0003 yang berstatus draft).

### Prasyarat baseline

Selain dua hal yang sudah dicatat di atas, bertambah satu yang baru terbukti dari DDL: `PEGA_JSON_OS_AKSEP_KLAIMTNP` **selalu `INSERT` dan tidak pernah `UPDATE`** — logika upsert-nya ada tetapi dikomentari seluruhnya. Bersama ketiadaan primary key pada `OS_AKSEPTASI_KLAIM`, baris ganda pada tabel itu mungkin terjadi dan harus diukur sebelum perbandingan dimulai.

## Tambahan — nilai IDR lama tidak dipakai sebagai pembanding

Konversi mata uang di sistem lama **tidak dapat direproduksi**: `POOLDATA.GETCURRENCYSTANDARD` mengabaikan parameter tanggalnya dan selalu memakai kurs terbaru sampai `sysdate` (`FINDING-001` bagian 7.1). Angka IDR yang dihitung tahun lalu tidak dapat dihasilkan ulang hari ini.

Karena itu, untuk nilai hasil konversi, **shadow-run membandingkan nilai mata uang asli, bukan nilai IDR-nya.** Nilai IDR lama diterima apa adanya sebagai nilai historis dan tidak dijadikan pembanding.

Pengecualian ini berlaku hanya untuk kelompok nilai hasil konversi. Nilai yang sejak awal berdenominasi rupiah tetap dibandingkan seperti biasa, dengan aturan dua ambang di atas.
