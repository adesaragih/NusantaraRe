---
status: selesai
---

# 25: Arsip muatan keluar: apa yang benar-benar dikirim, tersimpan sebagai catatan


> **SELESAI 19 September 2026 — dan grant-nya bertentangan dengan tiketnya sendiri.**
>
> Butir 1 berbunyi *"baris arsip **tidak dapat** di-`UPDATE` maupun di-`DELETE` oleh akun aplikasi — tulis-sekali ditegakkan hak akses, bukan kesepakatan"*. Tiket `35` juga sudah menyebut pengecualian ini sejak ditulis. **DDL-nya memberi keempat hak**: `GRANT SELECT, INSERT, UPDATE, DELETE`.
>
> Diperbaiki: **`GRANT SELECT, INSERT`** saja.
>
> Sebabnya bukan kehati-hatian. Arsip ini menjawab pertanyaan **"apa yang benar-benar dikirim"**. Arsip yang dapat diubah menjawab pertanyaan lain — *"apa yang sekarang kita katakan telah dikirim"* — dan itu bukan catatan.
>
> Dan ia **bukan turunan yang dapat dihitung ulang**: `GetDataOS` menjumlahkan seluruh baris berkunci sama ber-`STS_REJECT = 0`, lalu `SaveDataToOSAksep_Act` mengurangkan hasilnya. **Baris adalah tambahan, bukan keadaan.** Satu `UPDATE` mengubah setiap selisih yang pernah dihitung sesudahnya.
>
> `IX_ARSIP_MUATAN_KELUAR_1` diperluas dengan **`DITOLAK`**, karena setiap penjumlahan selisih menyaringnya (butir 3).
>
> Butir 4 sudah terpenuhi: kelima besaran `NUMBER(38,2)`. **Ini satu-satunya tempat pembulatan tepi menjadi baris tersimpan** — di mana pun lagi skala kanonik 20 berlaku (ADR-0003).

*Asal: `T-39` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Setiap muatan yang dikirim ke Arasapas dan kasir tercatat apa adanya, dan muatan berikutnya dapat dihitung sebagai selisih terhadapnya.

`ARSIP_MUATAN_KELUAR`: tujuan, kunci (klaim, jenis reasuransi, mata uang), lima besaran uang **berskala 2**, kurs, penanda ditolak, pelaku, waktu kirim. `IX_ARSIP_MUATAN_KELUAR_1`, `CK_ARSIP_MUATAN_KELUAR_1`.

**Tidak termasuk:** Proses pengirimannya — lapisan aplikasi. Muatan dalam bentuk JSON — **tidak ada kolom JSON di sini**; kolomnya bertipe, karena daftar kolomnya sudah terbaca dari S1 (ADR-0028).

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap


**Dasar:** DECIDED-TEKNIS(AK-1) — *"bila muatan keluar perlu diarsipkan, arsipnya tabel tersendiri, bukan kolom di tabel bisnis."* Ini pemakaian pertamanya. EVIDENCED: `RDBList\GetDataOS.xml` menjumlahkan seluruh baris berkunci `(CASEID, TypeLoss, Currency)` dengan `STS_REJECT = 0`, dan `SaveDataToOSAksep_Act` mengurangkan hasilnya — **baris adalah tambahan, bukan keadaan**.

- [x] Baris arsip tidak dapat di-`UPDATE` maupun di-`DELETE` oleh akun aplikasi — `GRANT SELECT, INSERT` saja. **Sebelumnya keempat hak diberikan.**
- [x] Dua muatan atas kunci yang sama pada waktu berbeda **keduanya tersimpan** — tidak ada `UNIQUE` atas kuncinya, dan itu disengaja.
- [x] Muatan bertanda ditolak **tidak** ikut terjumlah — `DITOLAK` ada di kolom **dan** di `IX_ARSIP_MUATAN_KELUAR_1`.
- [x] Nilai tersimpan **berskala 2** — `NUMBER(38,2)` pada kelima besaran.

**Ketidakpastian:** Tiga hal yang ditulis terang supaya tidak dibaca sebagai pelanggaran ADR: view kompatibilitas **membaca** arsip sehingga ia tetap view, bukan salinan (ADR-0023 utuh); yang menulis arsip adalah proses pengiriman lewat akun aplikasi (ADR-0017 utuh); dan arsip berisi **apa yang benar-benar dikirim**, bukan apa yang seharusnya — ia catatan, bukan turunan yang dapat dihitung ulang.
