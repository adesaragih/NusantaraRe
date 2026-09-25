# ADR-0051 — Anggap ada setidaknya satu konsumen hilir yang tidak dikenal

**Status:** diterima **tanpa verifikasi**, 23 September 2026
**Berlaku untuk:** Treaty In

## Konteks

Tidak ada satu pun pembaca tabel terbitan di dalam ekspor Pega — konsumennya, bila ada, berada di
luar. Uji F dan G, yang akan menyebutkan siapa mereka, **tidak dijalankan**.

Pencocokan antar-sistem yang terlihat di sistem lama memakai **tujuh karakter pertama** nomor
offer, yang terikat pada format ID `kode_situs || nomor urut enam digit`.

## Keputusan

1. **Anggap ada setidaknya satu konsumen yang tidak dikenal.**
2. Bentuk terbitan hilir dipertahankan sebagai **turunan** dari bentuk kanonik baru — view, bukan
   salinan (ADR-0023).
3. View itu **memaparkan ID lama dalam format lama**, supaya pemotong tujuh karakter tetap bekerja.
4. Setiap kontrak yang dimigrasi **menyimpan ID lamanya secara permanen**, bukan sementara. Orang,
   surat, dan dokumen merujuk padanya.

## Dasar — ongkos yang tidak setangkup

Ini bukan taruhan atas kemungkinan, melainkan atas **ongkos salahnya ke dua arah**:

| Bila salah | Akibatnya |
|---|---|
| Menganggap ada konsumen, padahal tidak | memelihara sebuah view yang tidak dibaca siapa pun — **murah** |
| Menganggap tidak ada, padahal ada | integrasi putus, dan baru ketahuan **di produksi** — **mahal** |

Bila ongkosnya tidak setangkup, ambil sisi yang murah tanpa menunggu bukti.

## Syarat pembalikan

Sensus konsumen menunjukkan tidak ada pembaca sama sekali — view boleh dipensiunkan.
**Pemensiunan mudah, pemulihan tidak**, dan itu sebabnya urutannya begini dan bukan sebaliknya.
