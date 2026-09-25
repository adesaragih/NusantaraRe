---
status: selesai
---

# 30: Rekap per klaim per mata uang, dengan Retensi Cedant yang tidak pernah tercampur


> **SELESAI 19 September 2026.** Kedua kriteria sudah terpenuhi bentuk yang ada; satu kolom bertambah.
>
> **Tiga subkueri terpisah adalah ADR-0010 sebagai SQL.** Di sistem lama Retensi Cedant adalah **baris di dalam daftar alokasi yang sama** — sehingga setiap loop atas `SpreadingRisk` ikut melihatnya **kecuali menyaringnya sendiri**, dan setiap penulis loop harus ingat melakukannya (`BLUEPRINT.md` §2.4). Di sini ia tabel tersendiri dan kolom tersendiri: menjumlahkan alokasi layer **tidak dapat** memuat Retensi Cedant — bukan karena ada penyaring yang harus diingat, melainkan **karena barisnya tidak ada di sana**.
>
> **`CACAH_BELUM_LENGKAP` bertambah.** Sebuah penjumlahan tidak menyatakan apakah baris yang terjumlah sudah selesai; dua kelompok berangka sama dapat berarti hal yang sangat berbeda. ADR-0019 ditegakkan di tingkat baris oleh `KEADAAN_BARIS`, dan **hilang begitu `SUM` dijalankan** — kecuali agregatnya ikut membawanya keluar.
>
> `SUM` atas nol baris menghasilkan `NULL`, bukan `0`, dan itu **dipertahankan**: `NULL` berarti tidak ada baris untuk dijumlah, berbeda dari nol yang berarti sudah dihitung dan hasilnya nol.

*Asal: `T-17` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Penjumlahan alokasi tidak pernah diam-diam memuat Retensi Cedant.

`V_REKAP_KLAIM_MATA_UANG`.

**Blocked by:**

- ~~`17`~~ *(selesai)* — Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang


**Dasar:** DECIDED(ADR-0010). EVIDENCED: `BLUEPRINT.md` §2.4 — di sistem lama setiap loop atas `SpreadingRisk` ikut melihat baris Retensi Cedant kecuali menyaringnya sendiri.

- [x] Retensi Cedant muncul sebagai **kolom tersendiri**, bukan bagian dari jumlah alokasi layer.
- [x] Menjumlahkan kolom alokasi saja menghasilkan angka tanpa Retensi Cedant — tanpa penyaring apa pun di sisi pemanggil.
- [x] Kelengkapan baris sumber terbaca dari agregatnya — `CACAH_BELUM_LENGKAP`, baru.

**Ketidakpastian:** Baris Retensi Cedant di sistem lama membawa `TotalClaim` yang **selalu nol** pada cabang mata uang sama (ADR-0010 tambahan). Sistem baru menghitungnya seperti Retensi Cedant biasa (AK-2b), jadi selisih terhadap sistem lama **diharapkan** dan masuk pengecualian bernama shadow-run.
