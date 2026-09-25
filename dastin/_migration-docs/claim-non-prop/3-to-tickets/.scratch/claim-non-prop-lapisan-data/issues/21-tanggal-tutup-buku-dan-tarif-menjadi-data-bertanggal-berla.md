---
status: selesai
---

# 21: Tanggal tutup buku dan tarif menjadi data bertanggal berlaku


> **SELESAI 19 September 2026 — dan baris awalnya menghasilkan berkas DML pertama di proyek ini.**
>
> ⚠ **`ddl-usulan/Z00_ISIAN_AWAL.sql` berisi `INSERT`. Ia satu-satunya berkas DML di seluruh `ddl-usulan/`**, dan dipisahkan justru supaya kekecualian itu **terlihat** — bukan terselip di kaki sebuah berkas DDL tempat tidak ada yang mencarinya. Ke-22 berkas tabel, kesembilan view, dan berkas skema semuanya tetap menyatakan *"TIDAK ADA DML DI BERKAS INI"*, dan pernyataan itu tetap benar.
>
> Empat baris: tutup buku `25` berlingkup global; brokerage 2,5%; PPh 2%; PPN 2,2%.
>
> **Kenapa baris ini di berkas, bukan diketik orang saat pemasangan.** Karena keempat nilai itu **sekarang tertanam di kode** sistem lama, dan seluruh maksud ADR-0025 adalah memindahkannya dari kode ke data. Baris yang diketik orang tidak dapat ditinjau, tidak dapat dibandingkan, dan tidak meninggalkan jejak siapa yang memilih angkanya.
>
> Angka `25` punya **dua sumber kebenaran** di sistem lama — tertanam di `HitServiceToKasir_Act`, **dan** dibaca `PROC_GENERATE_SEQUENCE_NUMBER` dari `POOLDATA.TANGGAL_CLOSING`. Dua sumber untuk satu aturan adalah cacat yang tabel ini tutup, **tetapi hanya bila tabelnya benar-benar terisi**.
>
> `BERLAKU_SEJAK = 1900-01-01` bukan tanggal yang berarti — ia menyatakan *"sejak sebelum data tertua"*, sehingga tidak ada baris lama yang jatuh di luar masa berlaku tarif mana pun. `DIBUAT_OLEH = 'MIGRASI'` karena **tidak ada orang yang memutuskannya**: nilainya dibaca dari kode lama apa adanya, dan menuliskan nama seseorang akan mengarang pelaku untuk keputusan yang tidak pernah diambil siapa pun.
>
> **Ratifikasi tidak menahan.** `ASK-AKUNTANSI` butir 5 memverifikasi apakah ketiga tarif masih berlaku. Bila ternyata tidak, yang berubah **isi baris** — pekerjaan entri, bukan pekerjaan kode. Itu seluruh maksud ADR-0025.

*Asal: `T-10` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Satu sumber kebenaran untuk tanggal tutup buku dan tarif; keduanya terisi baris awalnya.

`TUTUP_BUKU`, `TARIF_BERLAKU`; baris awal: tutup buku `25` berlingkup global; tarif brokerage 2,5%, PPh 2%, PPN 2,2%, seluruhnya berlingkup global dan berlaku sejak sebelum data tertua.

**Tidak termasuk:** Riwayat tarif — **tidak ada bukti tarif pernah berubah** (AK-5), jadi tidak ada riwayat yang dibuat-buat. Faktor `102,2` **tidak disimpan**; ia diturunkan sebagai `(100 + tarif PPN)`.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** DECIDED(ADR-0025). DECIDED-TEKNIS(AK-5, AK-6.3). EVIDENCED: angka `25` tertanam di `HitServiceToKasir_Act`; `PROC_GENERATE_SEQUENCE_NUMBER` membacanya dari `POOLDATA.TANGGAL_CLOSING` — **dua sumber kebenaran untuk satu aturan**.

- [x] Dua baris dengan (lingkup, tanggal berlaku) sama **ditolak** — `UQ_TUTUP_BUKU_1`, `UQ_TARIF_BERLAKU_1`.
- [x] Tanggal tutup buku di luar 1–31 **ditolak** — `CK_TUTUP_BUKU_1`.
- [x] Tarif terbaca lewat satu kueri tanpa membaca kode mana pun.
- [x] Keduanya **terisi baris awalnya** — `ddl-usulan/Z00_ISIAN_AWAL.sql`, 1 baris tutup buku dan 3 baris tarif.

**Ketidakpastian:** **A8** dan **ASK-AKUNTANSI no. 5** — apakah ketiga tarif masih berlaku **belum dikonfirmasi**. Yang dibuat tempatnya; nilainya berlabel DECIDED-TEKNIS menunggu ratifikasi, dan **ratifikasi tidak menahan tiket ini**.
