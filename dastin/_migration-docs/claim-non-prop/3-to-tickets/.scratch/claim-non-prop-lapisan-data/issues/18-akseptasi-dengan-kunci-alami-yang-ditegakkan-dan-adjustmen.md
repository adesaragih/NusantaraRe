---
status: selesai
menunggu-luar: [REQ-018]
---

# 18: Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran

> **DILEPAS DARI KETERGANTUNGAN KOMITE, 19 September 2026.** Folder Komite **tertutup untuk batch ini**. H1/H2 dan I1/I2 tinggal permanen sebagai hipotesis sejajar; tiket ini **tidak menunggunya**.
>
> **Dikerjakan di atas aturan yang sama-sama benar di kedua cabang:**
>
> - **I** — model **menyediakan tempat** bagi transisi `AcceptanceStatus` yang datang dari luar, **tanpa mengandaikan ia terjadi**. Keadaannya **empat, termasuk kosong**, sesuai temuan **C7**: properti ini dibandingkan sebagai angka (`0`,`1`,`2`) *dan* sebagai teks, dan bentuk teksnya memuat nilai kosong tanpa padanan numerik. Nilai kosong itu justru yang diuji `CloseClaimMD` (D2).
> - **H** — nilai Adjustment tetap dipasang **berpasangan mata uang** dan **tidak mengandaikan** salah satunya sudah rupiah.
>
> Tidak boleh ada rancangan yang mengandaikan salah satu cabang benar.


> **SELESAI 19 September 2026 — dan satu kolom harus dilonggarkan supaya keadaan keempat dapat disimpan.**
>
> `KEPUTUSAN_KOMITE` dipasang **`NOT NULL`**. Itu membuat **keadaan keempat mustahil disimpan**, dan keadaan keempat justru yang paling berakibat.
>
> S12 menemukan properti ini dibandingkan sebagai **angka** (`0`, `1`, `2`) **dan sebagai teks**, dan bentuk teksnya memuat **nilai kosong** yang tidak punya padanan numerik. Keadaannya **empat**, bukan tiga — itu sudah ditulis di C7 dan `SPEC-MODEL-DATA.md` §21.1, tetapi DDL-nya belum menurutinya.
>
> Dan yang kosong bukan keadaan sepele: penjaga `CloseClaimMD` menolak `AcceptanceStatus == "0"`, sementara Adjustment yang **belum pernah dikirim** ke komite bernilai **kosong**, bukan `"0"` — sehingga ia **lolos** (**D2**). Seluruh pembacaan jalur tutup-langsung bergantung pada perbedaan itu.
>
> Sekarang: `NULL` = **belum pernah dikirim ke komite**, berbeda dari `'0'` = **sudah dikirim, belum diputus**. ADR-0019 sebagai kolom, bukan sebagai kalimat.
>
> **Tanpa `CHECK` domain, dan itu tetap keputusan**: domainnya EXTERNAL — modul Komite yang menulisnya, dan folder itu tertutup. Domain tertutup di sini akan menolak nilai yang sah dari seberang.
>
> Kriteria lain sudah terpenuhi bentuk yang ada: `UQ_AKSEPTASI_1` (empat field layer + mata uang — **lebih ketat** daripada yang tiket sebut), `UQ_ADJUSTMENT_1`, `FK_ADJUSTMENT_2` menolak penghapusan rekening yang masih dirujuk, `CK_AKSEPTASI_1` domain tiga keadaan.

*Asal: `T-06` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Akseptasi kedua untuk klaim, layer, dan mata uang yang sama ditolak. Satu klaim dapat memiliki banyak Adjustment bernomor urut tetap.

`AKSEPTASI`, `ADJUSTMENT`, `REKENING_PENERIMA`; `UQ_AKSEPTASI_1`; `UQ_ADJUSTMENT_1`; `FK_ADJUSTMENT_2` beserta `IX_ADJUSTMENT_1`; `CK_AKSEPTASI_1`.

**Tidak termasuk:** `CHECK` domain untuk `KEPUTUSAN_KOMITE` — **sengaja tidak ada**, domainnya EXTERNAL. Tabel Komite dan keanggotaannya tidak dirancang.

**Blocked by:**

- **REQ-018** — **menunggui, tidak menahan** (dari luar papan): ia mengukur baris ganda pada **kunci lama yang lebih longgar**; kunci di tiket ini lebih ketat, jadi jawabannya batas bawah, bukan penentu
- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap
- ~~`07`~~ *(selesai)* — Sapuan S3: perilaku per nilai `PaymentType`


**Dasar:** DECIDED(ADR-0024, ADR-0015, ADR-0023). EVIDENCED: blok yang **dikomentari** di `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` memakai persis kombinasi itu; `OS_AKSEPTASI_KLAIM` **tanpa primary key maupun unique**, dan prosedurnya **selalu `INSERT`**.

- [x] `INSERT` kedua dengan (klaim, layer, mata uang) sama **ditolak** — `UQ_AKSEPTASI_1`, yang menegakkan **keempat** field layer, bukan satu.
- [x] Adjustment dengan nomor urut yang sudah dipakai pada klaim yang sama **ditolak** — `UQ_ADJUSTMENT_1`.
- [x] Menghapus rekening yang masih dirujuk Adjustment **ditolak** — `FK_ADJUSTMENT_2`.
- [x] `KEADAAN_AKSEPTASI` di luar tiga nilai domain **ditolak** — `CK_AKSEPTASI_1`.
- [x] `KEPUTUSAN_KOMITE` bernilai apa pun **diterima** — dan itu keputusan, bukan kelalaian.
- [x] `KEPUTUSAN_KOMITE` **kosong diterima**, dan dapat dibedakan dari `'0'`. **Sebelumnya `NOT NULL`** — keadaan keempat tidak dapat disimpan sama sekali.

**Ketidakpastian:** **I1/I2** — transisi status akseptasi `0 → 1/2` **TIDAK DITEMUKAN** di 279 berkas; ia ditulis satu kali saja sebagai `0`. **H1/H2** — apakah `.ValueAdjustment` sudah IDR di hulu belum diketahui, sehingga kolom nilai Adjustment dipasang berpasangan mata uang dan **tidak mengandaikan** salah satunya. **REQ-018** mengukur baris ganda pada kunci lama, yang lebih longgar dari kunci ini.
