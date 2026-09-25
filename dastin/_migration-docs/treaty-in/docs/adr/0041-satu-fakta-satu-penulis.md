# ADR-0041 — Satu fakta, satu penulis, selalu

**Status:** diterima, 23 September 2026
**Berlaku untuk:** masa berdampingan Pega dan sistem baru

## Konteks

Gelombang rilis ditetapkan sebagai **urutan merilis**, bukan urutan membangun. Maka ada masa di
mana sebagian fungsi sudah dilayani sistem baru sementara sisanya masih dilayani Pega, atas kontrak
yang sama.

Pilihan "kedua sistem sama-sama menulis, dibagi per kontrak" ditolak: Pega menulis dokumen JSON
sebagai bentuk kanoniknya, sistem baru menulis skema relasional sebagai bentuk kanoniknya. Dua
bentuk kanonik yang sama-sama ditulis melanggar ADR-0023. Itu bukan pilihan berisiko — ia pilihan
yang melanggar keputusan yang sudah mengikat.

## Keputusan

> Pada setiap saat, **tepat satu** sistem berwenang menulis sebuah fakta. Sistem yang lain boleh
> **membacanya**, tidak pernah menulisnya.

Penegakannya **bukan kebijakan, melainkan izin**: hak tulis Pega dicabut di basis data atas objek
yang fungsinya sudah pindah.

Urutan operasional tiap gelombang:

1. cabut hak tulis Pega atas objek yang akan pindah;
2. **verifikasi tidak ada penulisan lagi** — aturan Pega yang gagal di sini justru menemukan
   penulis yang tidak kita ketahui. Langkah ini adalah **alat penemuan**, bukan sekadar pemeriksaan;
3. pindahkan datanya;
4. verifikasi hasil pemindahan (ADR-0043);
5. buka hak tulis sistem baru.

## Konsekuensi

- Masalah ketiadaan kolom waktu **lenyap**: tidak pernah ada dua catatan atas fakta yang sama untuk
  diperbandingkan umurnya.
- Tabrakan urutan ID dan nomor offer tidak mungkin terjadi, karena hanya satu sistem yang
  menerbitkan.
- Batas gelombang menjadi **batas kewenangan tulis** yang bisa diperiksa, bukan dipercaya.
- Sistem baru harus menyediakan **bentuk baca yang cocok dengan Pega** untuk fungsi yang belum
  pindah — view atas bentuk kanonik baru, bukan salinan.
- View itu **harus memaparkan ID lama dalam format lama**. Pencapaian dan produksi dicocokkan lewat
  pemotongan tujuh karakter pertama nomor offer; begitu ID berganti bentuk, pencocokan itu putus.
  Karena itu setiap kontrak yang dimigrasi **menyimpan ID lamanya secara permanen**, bukan
  sementara: orang, dokumen, surat, dan sistem lain merujuk padanya.
