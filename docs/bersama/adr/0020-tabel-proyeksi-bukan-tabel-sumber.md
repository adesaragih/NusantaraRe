---
status: accepted
tanggal: 2026-09-23
sumber: spec penyimpanan EDM Treaty In - `.scratch/edm-treaty-in/`, keputusan work owner
---

# Tabel proyeksi selisih bukan tabel sumber

Angka selisih endorsemen disimpan pada tabel **proyeksi** supaya dapat dibaca lewat SQL tanpa
menjalankan perhitungan. Tabel itu **bukan** sumber kebenaran - ia salinan hasil hitung.

Tiga aturan berikut yang membuatnya bukan tabel sumber. Ketiganya wajib, tidak boleh diringkas.

| | Aturan |
| ---: | --- |
| 1 | Ia **hanya ditulis** oleh sistem baru, di dalam **transaksi yang sama** dengan generasinya |
| 2 | Ia **boleh dihapus total dan dibangun ulang** - **tetapi hanya baris hasil hitung sistem baru** |
| 3 | Bila isinya berbeda dari hasil hitung ulang, **tabelnya yang salah**, bukan operannya |

## Kolom penanda asal baris

Setiap baris membawa penanda asal: hasil **migrasi** atau hasil hitung **sistem baru**.

Baris hasil migrasi **beku**. Perintah bangun ulang tidak boleh menyentuhnya - angkanya berasal
dari sistem lama dan tidak dapat dihitung ulang, karena rumus lama punya varian yang sengaja tidak
ditiru.

**Tanpa aturan 2, satu perintah bangun ulang menimpa seluruh angka historis dan tidak dapat
dikembalikan.** Itu sebabnya aturan ini ditulis sebagai catatan keputusan, bukan sekadar komentar
di kode.

## Kenapa bukan view basis data

View pernah diusulkan lalu **dibatalkan**. Alasannya: rumus selisih hidup di lapisan layanan
*(ADR-0021)*. View berarti rumus yang sama ditulis **dua kali** - sekali di kode, sekali di SQL -
dan cepat atau lambat keduanya bercabang tanpa ada yang tahu.

**Jangan dihidupkan kembali diam-diam.**

## Akibat

1. Tabel proyeksi ditulis di dalam transaksi yang sama dengan barisnya - tidak menyusul.
2. Perintah bangun ulang menyaring pada penanda asal. Test yang menemukan baris hasil migrasi ikut
   terhapus **gagal**.
3. Kunci penyaring disertakan pada tabel proyeksi supaya pembaca SQL tidak perlu join balik.
