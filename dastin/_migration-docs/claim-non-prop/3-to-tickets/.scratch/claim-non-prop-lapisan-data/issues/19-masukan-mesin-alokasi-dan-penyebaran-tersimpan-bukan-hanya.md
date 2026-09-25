---
status: selesai
---

# 19: Masukan mesin alokasi dan penyebaran tersimpan, bukan hanya keluarannya


> **SELESAI 19 September 2026.**
>
> **Kriteria butir 1 DIUBAH, bukan sekadar diberi catatan penyimpangan.** Bunyi lama: *"nomor urut ganda dalam satu klaim ditolak, **di ketiga tabel**"*. Bunyi baru memisahkan **tabel yang punya kunci alami** dari yang tidak.
>
> Sebabnya: kriteria yang dibiarkan salah akan dibaca orang berikutnya sebagai **pekerjaan yang belum selesai**, dan ia akan menambahkan `NOMOR_URUT` ke `ESTIMASI_AWAL` justru untuk memenuhinya — melemahkan tabel demi memenuhi kalimat.
>
> `PEMBAGIAN_KERUGIAN` dan `PENYEBARAN` memakai `UQ (ID_KLAIM, NOMOR_URUT)`. **`ESTIMASI_AWAL` tidak punya `NOMOR_URUT` sama sekali**, dan kolomnya **tidak ditambahkan**.
>
> Sebabnya: tabel itu punya **kunci alami yang sesungguhnya** — `(klaim, jenis reasuransi, mata uang)`. Kunci alami **lebih kuat** daripada nomor urut: nomor urut hanya melarang dua baris **bernomor sama**; kunci alami melarang dua baris **berarti sama**.
>
> Menambahkan `NOMOR_URUT` di sana justru akan **melemahkan** tabelnya — ia mengizinkan estimasi kedua untuk jenis dan mata uang yang sama asal nomornya berbeda. Yang dijanjikan tiket, baris ganda ditolak, **terpenuhi dan terpenuhi lebih ketat**; yang tidak terpenuhi hanya **bentuk** kuncinya.
>
> **Butir 2 dilengkapi**: `IX_PENYEBARAN_1` atas `(ID_KLAIM, JENIS_REASURANSI, MATA_UANG)`. Kolomnya sudah ada di baris itu — itu yang membuat janji *"tanpa membaca tabel lain"* benar; index ini yang membuatnya murah.

*Asal: `T-08` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Pembagian kerugian, penyebaran, dan estimasi awal tersimpan sebagai data.

`PEMBAGIAN_KERUGIAN`, `PENYEBARAN`, `ESTIMASI_AWAL`.

**Tidak termasuk:** `SpreadingAdjustment` dan `SpreadingAdjustmentQS` — keduanya **view**, bukan tabel. → T-18.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** EVIDENCED: `MEMORI_PEMAHAMAN.MD` §6.2 Langkah 2 — `CNPSpreadLoss` adalah **masukan** mesin alokasi. DECIDED(ADR-0011, ADR-0013): hitung ulang idempoten mensyaratkan seluruh masukan tersimpan. `CONTEXT.md` **Penyebaran**.

- [x] **Nomor urut ganda dalam satu klaim ditolak, pada tabel yang tidak punya kunci alami** — `PEMBAGIAN_KERUGIAN` dan `PENYEBARAN`, lewat `UQ (ID_KLAIM, NOMOR_URUT)`.
- [x] **Pada tabel yang punya kunci alami, kunci itu yang ditegakkan, bukan nomor urut** — `ESTIMASI_AWAL`, lewat `UQ_ESTIMASI_AWAL_1 (ID_KLAIM, JENIS_REASURANSI, MATA_UANG)`.
- [x] Baris penyebaran dapat dijumlahkan per jenis reasuransi per mata uang tanpa membaca tabel lain — kolomnya di baris itu, `IX_PENYEBARAN_1` menopangnya.

**Ketidakpastian:** Tidak ada REQ yang menyentuhnya.
