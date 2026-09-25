# Permintaan ke DBA — dua nilai versi instance

**Untuk**: DBA Oracle, instance tempat `POOLDATA` berjalan
**Dari**: tim migrasi Claim Non Prop
**Tanggal**: 19 September 2026
**Sifat**: **READ-ONLY, verifikasi.** Tidak ada pekerjaan kami yang tertahan olehnya.
**READ-ONLY.** Dua perintah, keduanya membaca metadata instance. Tidak menyentuh tabel, tidak menyentuh data nasabah, tidak ada DML.

---

## Yang diminta

```sql
SELECT version_full FROM v$instance;
```

```sql
SHOW PARAMETER COMPATIBLE
```

Bila `version_full` tidak tersedia pada versi instance ini, `SELECT banner FROM v$version;` sudah memadai sebagai gantinya.

**Bentuk jawaban yang diterima**: dua baris teks, disalin apa adanya. Tangkapan layar juga diterima. Tidak perlu diformat.

---

## Kenapa keduanya ditanyakan, bukan hanya versinya

Karena **keduanya dapat berbeda, dan yang menentukan justru yang kedua.**

Oracle membatasi panjang nama objek — tabel, kolom, index, constraint — pada **30 byte** sampai versi 12.1, dan melonggarkannya ke **128 byte** sejak 12.2. Tetapi pelonggaran itu **tidak otomatis mengikuti versi biner yang terpasang**: ia mengikuti parameter `COMPATIBLE`. Sebuah instance 19c yang berjalan dengan `COMPATIBLE = 12.1.0` tetap menolak nama 31 byte.

Jadi versi saja tidak menjawab pertanyaan kami, dan `COMPATIBLE` saja juga tidak. Kami perlu keduanya.

## Mengapa masih ditanyakan meski tidak menahan apa pun

**Rancangan kami tidak menunggu jawaban ini.** Seluruh pengenal skema baru dijaga di bawah **30 byte**, dan itu kami tetapkan sebagai aturan tetap — bukan sebagai pengamanan sementara sampai nilai `COMPATIBLE` diketahui. Tiga puluh byte sah di setiap versi Oracle; seratus dua puluh delapan hanya di sebagian. Memilih yang berlaku di mana-mana membuat skema ini benar tanpa bergantung pada nilai yang tidak kami pegang.

Jadi apa pun jawaban Anda, **tidak ada nama yang berubah**.

Yang kami perlukan dari jawabannya dua hal, dan keduanya soal pemasangan, bukan rancangan:

1. **Memastikan tidak ada kejutan saat DDL dijalankan.** Skrip kami belum pernah dijalankan di mana pun. Mengetahui versi dan `COMPATIBLE` lebih dulu lebih murah daripada menemukannya lewat galat.
2. **Mengetahui fitur mana yang tersedia saat pemasangan.** Beberapa hal yang kami pakai — misalnya Unified Auditing untuk kebijakan pengawasan tulis — berperilaku berbeda antar versi dan antar mode.

**Rujukan internal**: `REQ-032`, berstatus **VERIFIKASI** sejak 19 September 2026. Sebelumnya BLOCKER; diturunkan karena batas 30 byte dipilih sebagai aturan tetap.

## Konteks singkat

Skema baru untuk modul Claim Non Prop akan berdiri **sebagai skema tersendiri di instance yang sama dengan `POOLDATA`** — bukan instance baru, bukan mesin lain. Keputusan itu sudah final per 19 September 2026.

Permintaan ini tidak meminta perubahan apa pun pada instance. Ia hanya meminta dua nilai yang sudah ada di sana.

