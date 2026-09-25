---
status: selesai
---

# 16: Nilai kerugian per mata uang, dengan nilai IDR yang tidak dapat lahir tanpa kurs


> **SELESAI 19 September 2026 — dan dua cacat ditemukan saat mengerjakannya, bukan saat meninjaunya.**
>
> ### Cacat 1 — arah constraint terbalik
>
> Bentuk lama, dipasang di **23 constraint pada 10 tabel**:
>
> ```sql
> CHECK ((X_IDR IS NULL AND KURS IS NULL) OR (X_IDR IS NOT NULL AND KURS IS NOT NULL))
> ```
>
> Ia **menolak baris ber-`KURS` terisi tetapi `X_IDR` kosong** — dan itu persis kriteria penerimaan yang tiket `16` sudah tuliskan sejak awal: *"`KURS` terisi dan `*_IDR` kosong **diterima** — belum dikonversi, bukan salah."* Kriterianya benar; DDL-nya tidak menurutinya.
>
> ### Cacat 2 — saling mengunci, dan ini yang lebih berbahaya
>
> Lima tabel memasang **empat sampai lima** constraint berbentuk itu di atas **satu kolom `KURS` yang sama**. Akibatnya bukan penjumlahan, melainkan perkalian:
>
> > Begitu **satu** nilai IDR diisi, `KURS` menjadi `NOT NULL`. Cabang pertama setiap constraint lain gugur seketika — sehingga **keempat nilai IDR wajib terisi bersama**.
>
> Artinya: mustahil mencatat nilai kerugian dalam rupiah tanpa **sekaligus** mencatat rupiah untuk biaya penilaian, salvage, dan biaya lain. Padahal ADR-0029 baru saja menetapkan kesembilan nama yang hari ini tanpa padanan IDR **boleh berkolom kosong** — dan bentuk ini melarangnya.
>
> **Tidak ada yang pernah memutuskan itu.** Ia akibat bentuk, bukan niat. Dan ia tidak akan ketahuan sampai baris pertama dimasukkan.
>
> ### Yang berlaku sekarang — satu arah
>
> ```sql
> CHECK (X_IDR IS NULL OR KURS IS NOT NULL)
> ```
>
> *"Nilai rupiah tidak dapat lahir tanpa kurs"* (ADR-0029) — **dan tidak lebih dari itu.** Diperbaiki di 23 constraint, 10 berkas.
>
> ### Yang bertambah — asal-usul kurs, akhirnya ditegakkan
>
> ADR-0029 menuntut **lima** hal per nilai uang: nilai asli, kode mata uang, nilai rupiah, kurs yang dipakai, dan **asal-usul kurs itu**. Empat yang pertama sudah ditegakkan; yang kelima tidak pernah — `KURS_SUMBER` dan `KURS_TANGGAL` ada sebagai kolom, tetapi tidak ada yang mewajibkannya.
>
> Itu bukan kerapian. **FINDING-006**: `GETCURRENCYSTANDARD` mengembalikan `1` ketika kurs **tidak ditemukan**, dan angka `1` tidak dapat dibedakan dari kurs yang sah. Yang membedakannya **hanya asal-usulnya**. Kurs tanpa asal-usul mewarisi cacat itu utuh ke sistem baru.
>
> Ditambahkan di **10 tabel**: `KURS IS NULL OR (KURS_SUMBER IS NOT NULL AND KURS_TANGGAL IS NOT NULL)`.

*Asal: `T-04` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Nilai per mata uang tersimpan; nilai IDR tanpa kurs ditolak basis data.

`NILAI_KLAIM_MATA_UANG`, empat besaran uang berpasangan IDR, `UQ_NILAI_KLAIM_MATA_UANG_1`, empat `CK` pasangan IDR–kurs, `CK` domain keadaan.

**Tidak termasuk:** Penularan keadaan ke baris turunan — diuji di T-28 setelah tabel turunannya ada.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** DECIDED(ADR-0007, ADR-0014, ADR-0019). EVIDENCED: FINDING-006 — `GETCURRENCYSTANDARD` mengembalikan `1` bila kurs tidak ditemukan, dan `1` tidak dapat dibedakan dari kurs yang sah.

- [x] Baris dengan `NILAI_KERUGIAN_IDR` terisi sementara `KURS` kosong **ditolak** — berlaku untuk keempat besaran.
- [x] Baris dengan `KURS` terisi dan `*_IDR` kosong **diterima** — belum dikonversi, bukan salah. **Sebelumnya ditolak** oleh bentuk constraint lama; itulah cacat 1.
- [x] Baris kedua dengan (klaim, mata uang) sama **ditolak** — `UQ_NILAI_KLAIM_MATA_UANG_1`.
- [x] Nilai nol pada `*_IDR` **diterima** dan dapat dibedakan dari `NULL` lewat kueri — kolomnya nullable, tidak ada default.
- [x] Keempat nilai IDR **tidak saling mengunci**: mengisi satu tidak memaksa tiga lainnya. **Sebelumnya memaksa**; itulah cacat 2.
- [x] `KURS` terisi tanpa `KURS_SUMBER` **ditolak** — ADR-0029 butir kelima, FINDING-006.

**Ketidakpastian:** Tidak ada REQ yang menyentuhnya. Dua cacat yang ditemukan hari ini **bukan ketidakpastian** — keduanya salah tulis kami sendiri, dan keduanya sudah diperbaiki di tempatnya.
