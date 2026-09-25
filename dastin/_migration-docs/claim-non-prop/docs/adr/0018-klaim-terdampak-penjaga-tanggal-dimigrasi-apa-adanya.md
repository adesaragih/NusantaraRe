---
status: accepted
label: DECIDED
---

# Klaim yang terdampak penjaga tanggal yang salah dimigrasi apa adanya, dan didaftar

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Berlaku hanya untuk keadaan sistem pada kedua ekspor itu.

Klaim yang hasilnya berbeda bila penjaga `CheckDateDOL_Act` dijalankan dengan perbandingan yang benar **dimigrasi apa adanya**. Tidak ada koreksi otomatis, tidak ada perubahan status.

Sebagai gantinya: penjaga yang benar dijalankan atas seluruh data lama, hasilnya menjadi daftar, dan keputusan per klaim diambil manusia.

## Consequences

Klaim yang sudah dibayar tidak boleh berubah statusnya karena sistem baru menghitung ulang sebuah penjaga — itu bukan koreksi, itu mengubah angka yang sudah diterima akuntansi. Membiarkannya tanpa daftar berarti masalahnya ikut pindah tanpa ada yang tahu.

Tiga syarat yang mengikat:

1. **Daftar dihasilkan sebelum cutover.** Ia prasyarat, bukan laporan pasca-migrasi.
2. **Dua daftar terpisah, tidak digabung.** Yang seharusnya tertolak tetapi lolos; dan yang seharusnya lolos tetapi tertolak. Kelompok kedua lebih sensitif — itu klaim yang mungkin ditolak secara keliru.
3. **Keputusan per klaim dicatat.** Tidak ada koreksi massal.

**Keterjangkauan data — diperiksa dan hasilnya baik.** Sempat dikhawatirkan daftar ini tidak dapat dibuat karena `.EndDateTreaty` tersimpan IN-BLOB. Pemeriksaan menunjukkan sisi lainnya terjangkau: `.ClaimData.EndDateTreaty` diisi dari `pyWorkPage.TreatyInMaster.Termination`, dan **`POOLDATA.TREATYINDETAIL.TERMINATION` adalah kolom bertipe `DATE`**. `POOLDATA.TREATYCONTRACT.TREATYENDDATE` juga `DATE`. Sisi kiri perbandingan, `DATEOFLOSS`, adalah kolom `VARCHAR2(8)` di tabel work.

Jadi **REQ-020 dapat dijalankan dengan SQL biasa** — tanpa membongkar blob — dengan `TO_DATE(DATEOFLOSS,'YYYYMMDD')` dibandingkan terhadap kolom `DATE` itu. Yang masih perlu dipastikan hanyalah kunci penghubungnya; `MASTERID` termasuk kolom yang di-expose di tabel work, tetapi pasangannya di `TREATYINDETAIL` belum diverifikasi.
