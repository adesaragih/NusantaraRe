---
status: selesai
---

# 31: Lima PageList lama tersedia sebagai view, tanpa menyimpan hasil yang dapat menyimpang

> **DILEPAS DARI KETERGANTUNGAN KOMITE, 19 September 2026.** Folder Komite **tertutup untuk batch ini**.
>
> **26 penanda yatim golongan 3b** mengikuti kebijakan **E17**: view menampilkannya **apa adanya** sebagai data yang dimigrasi dan **ditandai tak berpemilik**. Tidak dibuang — membuang butuh kepastian yang tidak kita punya. Tidak diberi perilaku — perilakunya tidak terbaca.


> **SELESAI 19 September 2026.** Kelima view ada dan menghasilkan angka dari sumbernya; satu kolom bertambah di masing-masing.
>
> **`CACAH_BELUM_LENGKAP`** — alasannya sama dengan tiket `30`, dan di sini berlaku untuk lima view sekaligus. Baris yang belum lengkap **tetap ikut dijumlah**, dan itu disengaja: nilai dalam mata uang aslinya sudah benar sejak awal, yang belum ada hanya padanan rupiahnya. Membuangnya akan menghasilkan angka yang salah diam-diam — **cacat yang sama dari arah berlawanan**.
>
> **Tidak ada nilai IDR yang dijumlahkan di view mana pun.** Bila kelak ada yang menambahkannya, kolom ini penjaga yang membuat kesalahan itu terbaca alih-alih tersembunyi.
>
> **Dua catatan sistem lama masuk ke badan berkasnya**, karena keduanya menjelaskan kenapa view lebih baik daripada tabel di tempat ini:
>
> | View | Perilaku lama yang tidak dapat terulang |
> |---|---|
> | `V_PENYEBARAN_AGREGAT` | `CountSpreadingXOL` **menghapus lalu membangun ulang** isinya setiap kali dijalankan. Sebagai view, tidak ada yang tersimpan untuk menyimpang (ADR-0013) |
> | `V_TOTAL_NILAI_PERTANGGUNGAN` | satu rutin dedup Java **membuang baris** dari `.ClaimData.TotalInterestInsured` — data klaim, tanpa pencatatan, tanpa penanda, tanpa jejak (**D40**). Sebagai view, tidak ada baris yang dapat dibuang |
>
> **Penyaring quota share tetap dua nilai**, `'10028'` dan `'10004'`, dan ketidakpastiannya ditulis di badan `V07`: apakah hanya keduanya yang berarti quota share **TIDAK DITEMUKAN**, dan daftar penuhnya ada di `POOLDATA.REINSURANCETYPE` yang **tidak punya satu pun constraint**. Ditulis sebagai literal **di satu tempat** supaya bila nilai ketiga ditemukan, ada tepat satu baris yang perlu diubah.

*Asal: `T-18` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** `ListTotalEstimation`, `SpreadingAdjustment`, `SpreadingAdjustmentQS`, `TotalInterestInsured`, `ListClaimAcceptation` terbaca tanpa tabel penyimpan.

`V_TOTAL_ESTIMASI`, `V_PENYEBARAN_AGREGAT`, `V_PENYEBARAN_AGREGAT_QS`, `V_TOTAL_NILAI_PERTANGGUNGAN`, `V_REKAP_AKSEPTASI`.

**Blocked by:**

- ~~`17`~~ *(selesai)* — Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang
- ~~`18`~~ *(selesai)* — Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran
- ~~`19`~~ *(selesai)* — Masukan mesin alokasi dan penyebaran tersimpan, bukan hanya keluarannya
- ~~`20`~~ *(selesai)* — Objek, kronologi, dan dokumen klaim


**Dasar:** DECIDED(ADR-0013). EVIDENCED: `MEMORI_PEMAHAMAN.MD` §6.4 — kunci agregasi `SpreadingAdjustment` adalah `(Currency, TreatyName)`, yaitu **satuan kasar**; dan `CountSpreadingXOL` **menghapus lalu membangun ulang** isinya setiap kali dijalankan.

- [x] Setiap view menghasilkan angka yang sama dengan penjumlahan langsung atas tabel sumbernya.
- [x] Tidak ada tabel yang menyimpan angka-angka ini.
- [x] Kelengkapan baris sumber terbaca dari tiap agregat — `CACAH_BELUM_LENGKAP`, baru di kelimanya.

**Ketidakpastian:** Penyaring quota share memakai `JENIS_REASURANSI_ID IN ('10028','10004')` — EVIDENCED dari `CountLossAllocation_act` Langkah 10. Apakah hanya dua nilai itu yang berarti quota share **TIDAK DITEMUKAN**; daftar penuhnya ada di `POOLDATA.REINSURANCETYPE`, yang **tidak punya satu pun constraint**.
