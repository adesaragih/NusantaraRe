# 14: Special Acceptance dan tangga akseptasi putaran kedua

**What to build:** Dua cabang inti tangga akseptasi yang sempat **ditangguhkan** kini berjalan penuh —
jalur **Special Acceptance** dan **tangga putaran kedua**.

Keduanya pernah ditangguhkan karena rule yang dipanggilnya tidak ada di ekspor yang kami terima.
Rule-nya kini tiba, rujukannya sudah diverifikasi ulang ke korpus, dan aturannya berlaku: **rule ada →
cabang dipulihkan**.

⚠️ **Dua rule Special Acceptance BUKAN Activity.** Keduanya dieksekusi sebagai **`RequestType` pada
langkah RDB-List**, berkelas integrasi, berdampingan dengan penanda akses dan halaman browse. Artinya
yang perlu ditulis adalah **query dan pemanggilan bacanya** — bukan sebuah fungsi layanan. Menyebutnya
"activity" akan menyesatkan porting.

Cabang putaran kedua berbeda: ia **benar-benar** panggilan activity, dipanggil dari post-activity
validasi tanggal.

⚠️ **Jalur ini hanya membaca.** Tidak ada bagian tiket ini yang menulis ke basis data; jalur tulis
produksi masih menunggu dan berada di luar lingkup.

**Blocked by:** 11, 13

**Status:** ready-for-agent

- [ ] Jalur Special Acceptance berjalan lewat pemanggilan **baca** bergaya `RequestType`, bukan dipaksa menjadi panggilan fungsi layanan
- [ ] Kedua rule Special Acceptance diperlakukan sebagai **rule integrasi/SQL**, dengan komentar menyebut asal dan kelasnya
- [ ] Tangga putaran kedua berjalan sebagai panggilan activity dari post-activity validasi tanggal
- [ ] Keduanya tetap mematuhi aturan **satu keputusan = satu transisi** (tiket 11)
- [ ] Tidak ada operasi tulis di jalur ini
- [ ] Fixture menutupi keduanya; kasus uji membuktikan cabangnya **benar-benar tercapai**, bukan sekadar ada
- [ ] Status "ditangguhkan" dicabut di dokumentasi rancangan, dengan menyebut bukti rujukan yang memulihkannya
