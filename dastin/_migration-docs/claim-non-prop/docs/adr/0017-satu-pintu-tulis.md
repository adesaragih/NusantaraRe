---
status: accepted
label: DECIDED
---

# Data aplikasi hanya diubah lewat aplikasi

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.
> Berlaku hanya untuk keadaan sistem pada kedua ekspor itu.

Tidak ada jalur tulis ke data aplikasi selain lewat aplikasi itu sendiri. Proses batch yang perlu mengubah data memanggil layanan yang sama dengan yang dipakai layar, sehingga aturan dan jejaknya berlaku sama.

Di sistem lama, schema `POOLDATA` memiliki akses ke tabel work Pega `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` — artinya data case dapat diubah oleh prosedur basis data, tanpa melewati aplikasi dan tanpa jejak di dalamnya.

## Consequences

Alasannya berdasar rekam jejak analisis ini sendiri. Sudah ditemukan **tiga** mekanisme yang hasilnya tidak dapat dilacak ke pelakunya:

| Mekanisme | Dokumen |
|---|---|
| Penggantian identitas `DARTO` menjadi `CHRISTINEANGELINA` | `FINDING-002` bagian 9 |
| Keputusan alur diambil dari teks komentar bebas | `FINDING-002` bagian 6 |
| Suntingan manual tertimpa hitung ulang tanpa peringatan | `FINDING-004` |

Pintu tulis kedua akan menambah yang keempat.

**Keputusan desainnya final; ukurannya belum, dan itu risiko lingkup.** Bila REQ-021 menunjukkan ada proses produksi yang menulis lewat jalur itu, proses-proses tersebut harus ditulis ulang **sebelum** cutover — pekerjaan yang belum masuk perkiraan mana pun.

Karena itu **REQ-021 ditandai BLOCKER perkiraan biaya, bukan sekadar BLOCKER teknis.** Bila hasilnya besar, ia dibawa ke manajemen sebagai pokok tersendiri, bukan diselipkan ke dalam laporan teknis.

## Batas yang baru terlihat — 19 September 2026, saat tiket 01 dikerjakan

ADR ini menyatakan tidak ada jalur tulis ke data aplikasi selain lewat aplikasi. Penegakannya dirancang lewat `GRANT` dan `REVOKE` atas objek, dan itulah isi `ddl-usulan/V00_HAK_AKSES.sql`.

**Penegakan itu berlaku pada tingkat objek saja, dan itu belum pernah dikatakan.**

Oracle mengenal hak sistem berakhiran `ANY` — `SELECT ANY TABLE`, `INSERT ANY TABLE`, `UPDATE ANY TABLE`, `DELETE ANY TABLE`, `ALTER ANY TABLE`. **Hak `ANY` mengatasi hak objek.** Satu akun yang memegangnya menulis ke tabel mana pun di skema mana pun, dan seluruh `REVOKE` di `V00_HAK_AKSES.sql` tidak menyentuhnya sedikit pun.

Yang membuat ini bukan kekhawatiran teoretis: premis ADR ini sendiri adalah bahwa `POOLDATA` sudah memegang hak tulis yang tidak ia perlukan. Pertanyaan apakah ia juga memegang hak `ANY` **belum pernah diajukan** — tidak di ADR ini, tidak di ADR-0028, tidak di `V00_HAK_AKSES.sql`, dan tidak di satu pun REQ sebelum hari ini.

Diajukan sebagai **REQ-037**, ditandai BLOCKER. Pencabutannya bukan pekerjaan pemilik skema; ia pekerjaan DBA di tingkat instance.

**Sampai REQ-037 kembali, klaim ADR ini dibaca dengan batas ini terpasang**: satu pintu tulis ditegakkan pada tingkat objek, dan belum diketahui apakah ada pintu di tingkat sistem yang melewatinya. Keputusannya tidak berubah; yang berubah adalah apa yang boleh diklaim sudah tertutup.

## Apa yang benar-benar dapat dijanjikan skema ini — rumusan diperbaiki 19 September 2026

**Keputusannya tidak berubah.** Yang berubah rumusan klaimnya, dan itu diubah di sini supaya tidak ada yang mengutip janji yang tidak dapat ditepati.

### Yang tidak dapat dijanjikan

*"Tidak ada jalur tulis selain lewat aplikasi"* — tanpa syarat — **tidak dapat dipertahankan** selama REQ-037 terbuka. Enam jalur melewati hak objek, dan tidak satu pun dapat ditutup dari dalam skema:

| Jalur | Kenapa skema tidak dapat menutupnya |
|---|---|
| Hak sistem berakhiran `ANY` | dicabut di tingkat instance, oleh DBA |
| Prosedur definer's rights | pemiliknya di skema lain; **pola ini sudah dipakai** — `PEGA_JSON_OS_AKSEP_KLAIMTNP` |
| Hibah ke `PUBLIC` | tidak terlihat saat memeriksa akun satu per satu |
| `CREATE ANY TRIGGER` | menaruh penulis di dalam tabel kita sendiri |
| Peran yang memuat hak `ANY` | `DBA`, `IMP_FULL_DATABASE`, dan peran buatan lokal |
| `GRANT ANY PRIVILEGE` | **memulihkan jalur tulis besok**, sesudah pemeriksaan hari ini bersih |

Jalur terakhir mengubah sifat persoalan: REQ-037 mengukur **satu titik waktu**, sementara ADR ini bicara tentang **keadaan yang bertahan**. Jawaban bersih hari ini tidak menjamin keadaan besok.

### Yang dapat dijanjikan, dan ditegakkan

> **Skema ini tidak dapat mencegah tulisan dari luar pintu. Tetapi tidak ada tulisan dari luar pintu yang tidak meninggalkan jejak.**

Ditegakkan `ddl-usulan/Z01_PENGAWASAN_TULIS.sql`: setiap `INSERT`, `UPDATE`, dan `DELETE` atas tabel skema ini yang datang dari akun selain `KLAIMNP_APP` tercatat beserta siapa, kapan, dan tabel mana.

Syaratnya memakai `SESSION_USER`, bukan `CURRENT_USER`. Pembedaan itu yang membuatnya bekerja: di dalam prosedur definer's rights, `CURRENT_USER` menjadi pemilik prosedur dan jejaknya hilang — yaitu persis jalur kedua di tabel atas. `SESSION_USER` tidak berubah oleh prosedur siapa pun.

Kebijakan itu **dijalankan DBA, bukan pemilik skema**, dan `KLAIMNP` sengaja tidak diberi `AUDIT_ADMIN`: pengawasan yang dapat dimatikan oleh yang diawasi bukan pengawasan.

### Apa yang tetap terbuka

Pengawasan memperlihatkan, tidak mencegah. Menutup jalurnya tetap pekerjaan DBA di tingkat instance, dan itu tetap **REQ-037, BLOCKER**. Yang berubah: klaim ADR ini kini dapat dipertahankan apa adanya sambil REQ-037 berjalan, alih-alih menunggu REQ-037 untuk berarti.
