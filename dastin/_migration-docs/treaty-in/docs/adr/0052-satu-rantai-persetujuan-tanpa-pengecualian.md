# ADR-0052 — Satu rantai persetujuan empat tingkat, tabel perutean kosong dari pengecualian

**Status:** diterima **tanpa verifikasi**, 23 September 2026
**Berlaku untuk:** Treaty In

## Konteks

Sistem lama punya dua penyimpangan dari rantai empat tingkat: sebuah jalur alternatif lewat
penerima tugas kelompok, dan sebuah jalur revisi yang memendekkan empat tingkat menjadi dua.
Pemilihan jalur revisi dilakukan lewat tombol, tanpa aturan apa pun yang mengaturnya.

Uji H, yang akan merekonstruksi jalur yang benar-benar ditempuh, **tidak dijalankan**.

## Keputusan

1. **Satu rantai bawaan empat tingkat.**
2. Tabel peruteannya **kosong dari pengecualian**.
3. Jalur alternatif kelompok **tidak dibawa**.
4. Jalur revisi **tidak dipendekkan**. Berlaku aturan: *jalur untuk suatu perubahan tidak boleh
   lebih pendek daripada jalur yang diperlukan untuk nilai hasilnya.*

## Dasar — arah pembatalan yang tidak setangkup

Tidak ada bukti jalur alternatif itu pernah disahkan sebagai kebijakan, dan pemendekan jalur revisi
dipilih lewat tombol oleh orang yang sama yang mengisi kontraknya.

| Pilihan | Bila ternyata salah |
|---|---|
| **Membuang keduanya** | dibatalkan dengan **menambah satu baris tabel** |
| Mempertahankan keduanya | tidak dapat dibatalkan tanpa **mengaudit ulang seluruh kontrak yang sudah melewatinya** |

## Syarat pembalikan

Satu baris di tabel perutean — karena aturan perutean berbentuk data, bukan kode (ADR-0038).

Ini **bukan** salah satu keputusan yang mengubah bentuk model.
