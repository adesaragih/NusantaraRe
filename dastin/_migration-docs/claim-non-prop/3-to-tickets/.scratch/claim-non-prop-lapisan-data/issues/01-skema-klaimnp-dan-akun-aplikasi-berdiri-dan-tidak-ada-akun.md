---
status: selesai
---

# 01: Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis

> **Ditinjau ulang 19 September 2026 sesudah gerbang penamaan dicabut.**
> Kedelapan pengenal berkas ini diperiksa terhadap aturan final (SPEC bagian 16): terpanjang `KLAIMNP_PERAN_PEMASANGAN` **24 byte**, di atas 30 byte **nol**. Tidak ada nama yang berubah.
> Dua nama dicatat terbuka sebagai bukan-kata-utuh — `KLAIMNP` dan sufiks `_APP` — beserta alasan keduanya tidak diganti. Lihat kepala `ddl-usulan/00_SKEMA_DAN_AKUN.sql`.

> **SELESAI 19 September 2026.** DDL-nya `ddl-usulan/00_SKEMA_DAN_AKUN.sql` — **usulan, untuk dibaca, belum pernah dijalankan**.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-01` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Skema, akun pemilik, akun aplikasi `KLAIMNP_APP`, dan akun hilir `KLAIMNP_HILIR` ada; akun pemilik hanya dipakai saat pemasangan.

Pembuatan skema dan ketiga akun; tablespace; peran dasar.

**Tidak termasuk:** `GRANT` dan `REVOKE` atas objek — belum ada objeknya. → T-20.

**Blocked by:**

- None (can start immediately)


> **DILEPAS 19 September 2026.** ADR-0028 naik ke `accepted` lewat konfirmasi manusia. Gerbang DDL terangkat.
>
> ~~**Gerbang penamaan tidak ikut terangkat.** Sampai REQ-032 kembali: tanpa nama constraint, tanpa nama index, tanpa singkatan.~~ **Dicabut 19 September 2026** — batas 30 byte jadi aturan tetap, dan nama di berkas ini diperiksa terhadap aturan final tanpa satu pun berubah.
>
> **Penegakan ADR-0017 lewat `GRANT` dan `REVOKE` masuk lingkup tiket ini sebagai bagian pekerjaannya**, bukan catatan operasional yang menyusul. Alasannya satu kalimat: **hak yang tidak tertulis tidak dapat diaudit.** Yang dituntut: akun aplikasi mendapat tepat hak yang dipakainya dan tidak lebih; akun hilir hanya `SELECT`; akun pemilik dipakai saat pemasangan lalu tidak dipakai lagi.

**Dasar:** DECIDED(ADR-0028) bagian *Satu pintu tulis*; DECIDED(ADR-0017). Premisnya EVIDENCED: `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` memberi `POOLDATA` hak `ALTER, DELETE, INSERT, UPDATE` atas tabel work Pega.

- [ ] Akun aplikasi dapat membuat sesi dan melihat skema.
- [ ] Akun hilir dapat membuat sesi dan **tidak** melihat satu objek pun (belum ada).
- [ ] Akun pemilik terpisah dari akun aplikasi — bukan akun yang sama dengan nama berbeda.

**Ketidakpastian:** Tidak ada. **REQ-032 turun jadi verifikasi** 19 September 2026: seluruh pengenal dijaga ≤ 30 byte **sebagai aturan tetap**, bukan sebagai pengamanan sementara. Bila versinya 12.2+, yang ada hanyalah kelonggaran yang sengaja tidak dipakai.

## Hasil

- **Dua tablespace, tiga akun, tiga peran.** Pengenal terpanjang `KLAIMNP_PERAN_PEMASANGAN` = **24 byte** — aman di 12.1 maupun 12.2+, dan itu berlaku di setiap versi Oracle.
- **Gerbang penamaan dipatuhi**: nol nama constraint, nol nama index, nol singkatan. Seluruh pengenal kata utuh. `KLAIMNP_INDEKS` adalah nama *tablespace*, bukan nama index.
- **Nol DML.**
- **Pemisahan pemilik dari aplikasi ditegakkan, bukan disepakati.** Akun aplikasi diberi `QUOTA 0` pada kedua tablespace dan **tidak diberi `CREATE TABLE`** — tanpa itu ia dapat memiliki objek yang tidak pernah masuk DDL dan tidak terlihat siapa pun.
- **Peran bawaan Oracle `CONNECT` dan `RESOURCE` tidak dipakai**, dan dicabut. Isinya berubah antarversi; memakainya berarti hak akun ditentukan versi basis data, bukan ditentukan dokumen.
- **`CREATE DATABASE LINK` dicabut dari ketiganya** (ADR-0016), dengan catatan bahwa instance ini sudah memuat satu link produksi — nasibnya K1, belum diputuskan.
- **Penguncian akun pemilik ditulis tetapi sengaja dibiarkan sebagai komentar**: menjalankannya di tengah pemasangan menghentikan pemasangan itu sendiri. Ia dijalankan sesudah `V00_HAK_AKSES.sql`.
- **Lima blok pemeriksaan penerimaan** ditulis sebagai `SELECT`, menjawab ketiga kriteria terima.

## Temuan yang lahir dari mengerjakannya

**Hak sistem berakhiran `ANY` mengatasi hak objek, dan itu belum pernah diperiksa.** `V00_HAK_AKSES.sql` mencabut hak `POOLDATA` atas tiap tabel satu per satu; satu akun yang memegang `INSERT ANY TABLE` membatalkan seluruhnya sekaligus, dan tidak ada baris DDL di skema baru yang dapat mencegahnya — pencabutannya di tingkat instance, oleh DBA.

Tidak disebut di ADR-0017, tidak di ADR-0028, tidak di `V00_HAK_AKSES.sql`, tidak di satu pun REQ sebelum hari ini. Didaftarkan **REQ-037**, BLOCKER. Batasnya ditulis di badan ADR-0017 dan di kepala `V00_HAK_AKSES.sql`, bukan hanya dilaporkan di sini.

## Perbaikan sesudah tinjauan, 19 September 2026

Empat, seluruhnya dari tinjauan atas `00_SKEMA_DAN_AKUN.sql`.

**1. `CREATE SESSION` diperiksa — aman.** Diberikan eksplisit ke ketiga peran (baris 119, 127, 130), **sebelum** `REVOKE CONNECT`. Hak itu datang lewat peran, bukan lewat `CONNECT`, jadi pencabutan `CONNECT` tidak menutup pintu masuk.

**2. Cacat yang ditemukan pemeriksaan itu, dan diperbaiki.** Kesembilan `REVOKE` ditulis sebagai perintah telanjang dengan catatan *"pada instance bersih ia tidak melakukan apa-apa"*. **Itu salah** — Oracle mengangkat `ORA-01951` dan `ORA-01952`, bukan mengabaikan. Pada instance bersih, yaitu keadaan yang justru diharapkan, baris pertama **menggugurkan seluruh pemasangan**. Diganti blok yang mencabut bila ada dan diam bila tidak, menelan hanya kedua galat itu.

**3. Penguncian akun pemilik keluar dari komentar, jadi `Z02_KUNCI_PEMILIK.sql`.** Alasannya sejarah proyek ini: dasar ADR-0024 adalah blok yang dikomentari, dan `PEGA_JSON_OS_AKSEP_KLAIMTNP` selalu `INSERT` tanpa pernah `UPDATE` karena logika upsert-nya dikomentari seluruhnya. Kode yang dimatikan dengan niat dinyalakan nanti tidak pernah dinyalakan. Sapuan atas seluruh berkas pemasangan sekarang menemukan **nol baris perintah yang dikomentari**.

**4. Pengawasan tulis ditambahkan** — `Z01_PENGAWASAN_TULIS.sql`. Skema tidak dapat **mencegah** tulisan dari pemegang hak `ANY`, tetapi dapat **memperlihatkannya**. Syaratnya memakai `SESSION_USER`, bukan `CURRENT_USER`: di dalam prosedur definer's rights, `CURRENT_USER` menjadi pemilik prosedur dan jejaknya hilang — persis mekanisme `PEGA_JSON_OS_AKSEP_KLAIMTNP`. Dijalankan **DBA**, dan `KLAIMNP` sengaja tidak diberi `AUDIT_ADMIN`: pengawasan yang dapat dimatikan oleh yang diawasi bukan pengawasan.

**Akibatnya bagi ADR-0017**: rumusan klaimnya diubah di badannya dari *mencegah* menjadi **"tidak dapat mencegah tulisan dari luar pintu, tetapi tidak ada tulisan dari luar pintu yang tidak meninggalkan jejak"**. Keputusannya tidak berubah.

**REQ-037 diperluas dari satu jalur menjadi enam**, dan satu di antaranya mengubah sifat REQ itu sendiri: pemegang `GRANT ANY PRIVILEGE` dapat memulihkan jalur tulis **besok**, sesudah pemeriksaan hari ini bersih. REQ mengukur satu titik waktu; ADR bicara keadaan yang bertahan. Pembedaan itu tertulis di badan REQ-nya.

**Catatan kerapian yang ditunda**: `V00_HAK_AKSES.sql` bukan view, dan awalan `V` padanya tidak konsisten dengan `V01`–`V09` yang memang view. Tidak diganti nama karena sudah dirujuk ADR-0028, tiket `01`, dan tiket `35`.
