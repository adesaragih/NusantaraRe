---
status: accepted
label: DECIDED
---

# Keputusan alur disimpan sebagai field terstruktur; komentar tidak pernah dibaca mesin

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Setiap keputusan yang menggerakkan alur disimpan sebagai tiga field terpisah — **tindakan**, **pelaku**, **waktu** — bukan disimpulkan dari isi kolom komentar. Kolom komentar tetap ada sebagai catatan manusia, dan **tidak pernah** menjadi masukan bagi kondisi mana pun.

Sistem lama melakukan sebaliknya: `SethistoryKlaimTreaty` mengambil keputusan lewat `@contains(.CommentSuggest,"Accepted by <nama>")` pada enam pasang kondisi. Uraiannya di `FINDING-002-percabangan-identitas.md` bagian 6.

## Consequences

Tiga cacat sekaligus hilang, dan ketiganya nyata di sistem lama: satu huruf besar berbeda membuat kondisi gagal tanpa pesan galat; siapa pun yang boleh mengisi komentar dapat memenuhi kondisi alur; dan orang yang sama ditulis dalam dua ejaan sehingga sebagian cabang tidak pernah terpicu.

**Harga yang dibayar ada di migrasi data, bukan di rancangan.** Riwayat lama menyimpan keputusannya hanya di dalam teks, dan sebagian teks itu sudah ditimpa oleh langkah penerjemahan di sistem lama sendiri. Kolom "siapa menyetujui" pada sistem baru tidak dapat diisi lengkap dari data lama secara otomatis; sebagian akan kosong. Besarnya bagian yang kosong belum diketahui dan diukur lewat REQ-013.

Konsekuensi itu diterima secara sadar di muka. Ia bukan temuan yang muncul saat UAT.
