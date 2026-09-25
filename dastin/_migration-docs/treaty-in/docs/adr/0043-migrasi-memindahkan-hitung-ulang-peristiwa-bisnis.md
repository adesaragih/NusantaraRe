# ADR-0043 — Migrasi memindahkan; menghitung ulang adalah peristiwa bisnis tersendiri

**Status:** diterima, 23 September 2026
**Berlaku untuk:** migrasi data dan penerimaan kalkulator baru

## Konteks

Godaan terbesar sebuah migrasi adalah menjalankan kalkulator baru sambil memindahkan data, karena
terasa efisien. Akibatnya kedua peristiwa tercampur dan tidak ada yang bisa menjawab apakah
perbedaan angka berarti migrasi gagal atau perbaikan berhasil.

## Keputusan

Keduanya adalah peristiwa terpisah, dan kriteria keberhasilannya **berlawanan**:

| | Kriteria |
|---|---|
| **Migrasi** — memindahkan angka apa adanya | **Setiap angka identik.** Bukan mirip, bukan dalam toleransi. Selisih satu rupiah berarti migrasi gagal dan diulang. |
| **Menghitung ulang** dengan kalkulator baru | Hasilnya **berbeda**, dan perbedaannya **bisa dijelaskan** oleh cacat yang sudah dikenali. Perbedaan yang tidak bisa dijelaskan berarti kalkulator baru yang salah. |

Migrasi menulis baris warisan apa adanya. **Sentuhan pertama** (ADR-0042) yang menghitung ulang,
dan hitung ulang itu adalah peristiwa yang dibukukan, bukan hasil migrasi.

## Uji paritas

Kalkulator baru dijalankan atas seluruh data yang sudah dimigrasi, dan hasilnya **dikeluarkan
sebagai laporan selisih**. Tidak dituliskan ke mana pun.

Satu run itu punya dua kegunaan sekaligus:

1. **Uji terima kalkulator baru** — setiap selisih harus jatuh ke salah satu kelas cacat yang sudah
   dikenali. Selisih yang tidak jatuh ke kelas mana pun adalah cacat **baru** di kalkulator baru.
2. **Pengukuran kesalahan historis** — jumlahnya dalam rupiah, per kelas cacat, per tahun treaty.
   Itu angka yang dibutuhkan daftar eskalasi, dan datang gratis dari run yang memang harus jalan.

**Syarat yang tidak boleh dilewati:** tuliskan lebih dulu, per cacat yang sudah dikenali, **bentuk
selisih apa** yang seharusnya ia hasilkan — arahnya, besarannya, dan populasi terdampak. Tulis itu
**sebelum** run dijalankan. Tanpa itu, pembacanya akan merasionalisasi apa pun yang muncul, dan uji
terima berubah jadi latihan menjelaskan. Daftar harapan itu adalah keluaran tersendiri.

## Ukuran keberhasilan migrasi

1. Jumlah kontrak dan addendum cocok.
2. Setiap ID lama punya pasangan di sistem baru, **dan** setiap ID baru punya asal. Migrasi yang
   menciptakan baris dari ketiadaan sama buruknya dengan yang menghilangkan baris.
3. Nilai kunci dipindahkan **tanpa dihitung ulang** lalu dibandingkan — kriteria identik di atas.
4. **Jumlah baris anak per kontrak cocok**: layer, detail, penyebaran, installment, periode
   pelaporan. Kontrak yang headernya pindah utuh tetapi kehilangan satu baris layer akan lolos
   ketiga ukuran pertama, dan itu kegagalan yang paling sering terjadi pada migrasi dari dokumen
   JSON ke tabel relasional.
