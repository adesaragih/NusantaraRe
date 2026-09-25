---
status: accepted
label: DECIDED
---

# Menjalankan ulang perhitungan alokasi harus menghasilkan keadaan yang identik

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Perhitungan alokasi bersifat **idempoten**: menjalankannya dua kali atas masukan yang sama menghasilkan keadaan yang sama persis. Ini diuji otomatis, bukan sekadar dijanjikan.

## Consequences

Alasannya tidak punya tandingan: bila menekan tombol hitung dua kali menghasilkan angka berbeda, tidak ada satu pun angka di sistem itu yang bisa dipertanggungjawabkan.

**Konsekuensi terhadap shadow-run, dan ini yang paling mahal.** Bila sistem lama ternyata tidak idempoten, ada data produksi dengan baris ganda. Itu bukan sekadar urusan migrasi data — ia **merusak dasar perbandingan**. Hitungan sistem baru tidak dapat dibandingkan terhadap baris yang sudah terlanjur ganda, karena selisihnya akan terbaca sebagai kesalahan sistem baru padahal berasal dari kerusakan baseline.

Karena itu, sebelum shadow-run dimulai: baseline dibersihkan lebih dulu, **atau** case yang terdampak dikeluarkan dari perbandingan dan didaftar terpisah. Salah satu harus dipilih; tidak boleh dibiarkan tercampur. Lihat juga ADR-0005.

Yang membuat risiko ini nyata di sistem lama: `CountLossAllocation_act` melakukan `<APPEND>` baris `"UR"`, dan **tidak ditemukan di XML** penghapusan PageList `SpreadingRisk` sebelum `APPEND` itu. Satu-satunya penghapusan yang menyasar `SpreadingRisk` di seluruh folder adalah `Property-Remove` atas properti tunggal `.SpreadingRisk(Param.idx).AdjClaimValue` di dua rule — bukan atas daftarnya. Selain itu `CountLossAllocation_act` langkah 10 memuat `Property-Remove` yang **parameter sasarannya kosong** di ekspor ini.

Besarnya kerusakan diukur, bukan ditebak — REQ-014 menghitung berapa case yang memiliki lebih dari satu baris `"UR"`. Keputusan idempoten ini **tidak menunggu** hasil pengukuran itu.
