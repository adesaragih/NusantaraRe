---
status: selesai
---

# 10: Sapuan S5, S6, S7, S8, dan dua sumber yang belum pernah dibaca

> **SELESAI 2026-09-18.** Dijalankan 18 September 2026; hasilnya `SAPUAN-S3-S9.md` bagian S5, S6, S7, S8, dan dua sumber.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-37` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Empat pertanyaan terbuka terjawab atau tercatat nihil beserta polanya.

- **S5** penulis `AlokasiXOLPaid` — termasuk sebagai sasaran `Page-Copy`, `Page-New`, `RDB-List` ke page, dan hasil Report Definition, bukan hanya `Property-Set`.
- **S6** asal `"Previously Calculated UR"` — **jangan cari literal utuh**; cari potongan `"Previously"` dan `"Calculated"` di dalam ekspresi penyambungan, karena nilainya mungkin dirakit.
- **S7** `Property-Remove` di `CountLossAllocation_act` langkah 10 yang parameter sasarannya kosong — baca langkah itu utuh beserta tetangganya.
- **S8** `IsTreatyIn`, `IsReject`, `IsCloseFile`, `IsAnyAcceptation` — keempatnya belum berlaku pernyataan "bersih".
- **Dua sumber**: `pengetahuan/DDL_Script_ClaimNonProp2.xls` (**belum pernah dibaca**) dan `excludeXML/GetBase64Attachment.xml`.

**Tidak termasuk:** H1/H2, I1/I2, F1/F2 — berada di modul Komite.

**Blocked by:**

- None (can start immediately)


**Dasar:** `_selesai/OPEN-QUESTIONS.md` C5 dan C9; `PENGETAHUAN.md` §15.

- [ ] Tiap sasaran menghasilkan jawaban atau pernyataan nihil **beserta lapisan yang disapu dan pola yang dipakai**, sesuai `PENGETAHUAN.md` §14: yang sah hanya *"tidak ada X di lapisan yang sudah saya baca"*. Temuan yang mengubah putusan lama masuk `_selesai/OPEN-QUESTIONS.md` sebagai D9 dan seterusnya. Hasilnya dibandingkan terhadap `BLUEPRINT.md` §7.5 — tambalan per-case baru adalah temuan tersendiri (ADR-0012).

**Ketidakpastian:** S6 menentukan apakah `RETENSI_CEDANT` perlu **keadaan ketiga**; bila ya, kunci (klaim, mata uang) **belum cukup** dan T-05 berubah.

## Hasil sapuan

- **S5** — `AlokasiXOLPaid` **tanpa penulis** pada empat lapisan (EVIDENCED-NIHIL). Yang justru terbaca: daftar kolom muatan selisih jalur Komite. **D11**.
- **S6** — `"Previously Calculated UR"` memang **dirakit**, lalu diuji utuh. Kunci (klaim, mata uang) **belum cukup**; menaikkan pentingnya REQ-033. **D12**.
- **S7** — langkah 10 memang kosong: sasaran, deskripsi, prasyarat, dan parameter wajibnya seluruhnya tidak terisi. Ia sisa, bukan perilaku. Migrasi tidak perlu menirunya.
- **S8** — pernyataan “bersih” **tidak** berlaku seragam. `IsTreatyIn` dibaca empat Activity dengan tipe tak konsisten; `IsReject` dan `IsCloseFile` ditulis dan dibaca; `IsAnyAcceptation` **dibaca sepuluh kali dan ditulis nol kali**. **D10**.
- **`DDL_Script_ClaimNonProp2.xls` dibaca** — 39 objek, empat di antaranya menyentuh tiket berjalan: `DIRECTTOKASIR_LOG` (**D14**), `JSON_KLAIM` (**D13**, mengubah sifat REQ-018), `CLAIMXOL2`, `CLAIMREJECTED`.
- **`excludeXML/GetBase64Attachment.xml` TIDAK DIBACA** — ia di dalam folder `Komite Claim Non Prop`, yang tertutup. Saudaranya di folder terbuka dibaca sebagai gantinya, dan perbedaannya tidak diketahui.
