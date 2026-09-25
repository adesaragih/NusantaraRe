---
status: accepted
label: DECIDED
---

# Nilai hasil suntingan manual bertahan terhadap hitung ulang, dan selalu terlihat

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Ketika petugas mengganti nilai hasil perhitungan dengan nilai suntingan sendiri, nilai itu **dikunci** dan tidak ditimpa oleh perhitungan berikutnya. Layar menampilkan keduanya berdampingan — nilai hitungan yang digantikan dan nilai suntingan yang berlaku — beserta siapa yang menyunting dan kapan.

Alternatif yang ditolak: mengikuti perilaku sistem lama, yaitu suntingan selalu kalah terhadap hitung ulang.

## Consequences

Model data memerlukan konsep **nilai terkunci** sejak awal: setiap nilai yang dapat disunting menyimpan pasangan `nilai_hitungan` dan `nilai_disunting`, plus `pelaku` dan `waktu`. Menambahkan ini setelah tabel terbentuk berarti membongkar setiap baris yang sudah ada, karena nilai lama tidak dapat dipisahkan kembali menjadi dua asal yang berbeda.

Kombinasi terburuk yang dihindari keputusan ini adalah suntingan diam-diam yang dapat hilang diam-diam — persis keadaan yang terbaca di sistem lama dan dilaporkan di `FINDING-004-suntingan-manual-tertimpa.md`.

Keputusan ini **tidak menunggu** kepastian hipotesis F1/F2 tentang `.IsEditClaim`. Apa pun yang ditemukan di modul Komite nanti, sistem baru mengunci suntingan.

Konsekuensi migrasi data bergantung pada hipotesis mana yang benar, dan itu diuraikan di FINDING-004 bagian 4, bukan di sini.
