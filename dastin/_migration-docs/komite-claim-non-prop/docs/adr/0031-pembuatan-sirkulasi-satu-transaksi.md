> Modul  : Komite Claim Non Prop · Ronde 06 · 2026-09-20
> Peran  : juru catat
> Masukan: KETETAPAN.md K5-2, K5-3, K5-7, H-5 · GRILL-05/01-TEMUAN.md N-07, N-08 · GRILL-05/06-PUTUSAN.md T5-03
> Status : TERBUKA
> Sifat  : HIDUP

# ADR-0031 · Pembuatan sirkulasi dan akibatnya adalah satu transaksi

**Status:** Diterima · 2026-09-20 · sumber normatif `KETETAPAN.md` `K5-2`, didampingi
`K5-3` dan `K5-7`

## Konteks

Pada kedua activity pembuat sirkulasi, pemanggilan `pxAddChildWork` dijaga pra-syarat
`@hasMessages(myStepPage)` dengan aksi **lewati langkah** — bukan keluar activity.

Akibatnya, ketika pembuatan tidak jadi, langkah-langkah sesudahnya tetap berjalan seolah ia
berhasil. Keduanya membaca `pxCoveredInsKeys(<LAST>)` — case tercakup terakhir — dan menulis
potongannya ke klaim sebagai nomor komite, lalu menyimpan klaim, lalu mengirim surat. Bila
klaim pernah punya sirkulasi sebelumnya, yang tertulis adalah nomor **sirkulasi lain**. Bila
belum pernah, yang tertulis kosong.

Jalur kegagalan kedua lebih senyap lagi. `CreateChildKomiteCNP_Act` langkah 28 melompat ke
label `END` ketika daftar penyebaran risiko kosong; label itu menetapkan
`ClaimData.IsFlagError = "true"` dan menyimpan klaim. Nama properti itu muncul di **satu**
berkas dari 338 — berkas yang menulisnya. Tidak ada satu pun rule dalam ekspor yang
membacanya. Pengguna menekan tombol, tidak ada sirkulasi yang lahir, dan tidak ada apa pun
yang memberitahunya.

Persoalan ketiga menyertai keduanya: penanda maksud ditulis ke klaim **pada langkah yang
sama** yang membuat halaman anak, yaitu sebelum komite memutuskan apa pun. Klaim membawa
`IsCloseFile` atau `IsReject` sejak pengajuan — dan tetap membawanya ketika sirkulasinya
gagal lahir.

## Keputusan

Pembuatan sirkulasi, penulisan nomornya ke klaim, penyimpanan klaim, dan pengiriman
pemberitahuan **berhasil bersama atau gagal bersama**. Tidak ada langkah sesudah pembuatan
yang boleh berjalan ketika pembuatan tidak jadi.

Setiap jalur pembentukan yang berakhir tanpa sirkulasi berakhir dengan **kegagalan yang
terlihat** oleh pengguna dan tercatat. Penanda yang tidak dibaca siapa pun bukan penanganan
galat.

**Maksud dan akibat adalah dua penanda berbeda.** "Diajukan untuk ditutup" tercatat sejak
pengajuan; "ditutup" hanya lahir dari keputusan komite lewat fungsi akibat. Sirkulasi yang
gagal lahir meninggalkan maksudnya tercatat dan akibatnya tidak.

Tali antara klaim dan sirkulasi adalah relasi **dari sisi sirkulasi**, bukan nomor yang
disalin ke klaim.

## Konsekuensi

- Validasi menolak kasus-guna pembuatan, bukan menandai halaman — `H-5` sudah menetapkan
  ini; ADR ini memperluasnya ke seluruh urutan yang bergantung pada pembuatan.
- Pengguna yang gagal membentuk sirkulasi mendapat pesan, bukan kesunyian. Uji `P5-07` dan
  `P5-08` dirancang gagal.
- Klaim tidak lagi menyimpan nomor komite sebagai tali. Migrasi data **tidak boleh
  mempercayai** `.KomiteNo` yang tersimpan; berapa banyak klaim di produksi membawa nomor
  bukan miliknya adalah cacah yang belum diambil (`INVENTARIS-BUKTI.md` §2.5 baris 6).
- Penanda maksud dan penanda akibat menempati kolom berbeda, sehingga klaim yang pernah
  diajukan untuk ditutup tetapi tidak jadi ditutup dapat dibedakan dari klaim yang ditutup.

## Alternatif yang ditolak

**Menjaga tiap langkah sesudah pembuatan dengan pra-syarat "bila pembuatan berhasil".**
Ditolak: itu memindahkan kewajiban kepada penulis tiap langkah berikutnya, dan langkah
berikutnya akan selalu bertambah. Sistem lama menjalankan bentuk ini dengan satu penjaga,
dan penjaganya menjaga satu langkah.

**Membiarkan penulisan nomor tetap berjalan dan memperbaikinya belakangan lewat pekerjaan
rekonsiliasi.** Ditolak: klaim yang tersimpan dengan nomor milik sirkulasi lain sudah
mengirim surat atas nomor itu. Rekonsiliasi memperbaiki kolom, bukan surat.

**Menampilkan `IsFlagError` di layar dan menganggapnya selesai.** Ditolak: yang salah bukan
bahwa penandanya tak terlihat, melainkan bahwa jalur itu berakhir dengan penanda alih-alih
dengan penolakan.
