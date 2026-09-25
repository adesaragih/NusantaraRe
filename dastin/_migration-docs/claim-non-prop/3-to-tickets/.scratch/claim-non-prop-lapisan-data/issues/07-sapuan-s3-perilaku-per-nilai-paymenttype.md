---
status: selesai
---

# 07: Sapuan S3: perilaku per nilai `PaymentType`

> **SELESAI 2026-09-18.** Dijalankan 18 September 2026 atas empat lapisan; hasilnya `_migration-docs/claim-non-prop/SAPUAN-S3-S9.md` bagian S3.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-34` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Tabel nilai → perilaku yang dijaganya, disusun dari setiap kondisi dan ekspresi yang mengujinya.

Empat lapisan: Activity, Data Transform, Section/Harness/FlowAction, pemetaan RDB/REST.

**Tidak termasuk:** **Nama** tiap nilai — itu tetap **A6b**, pertanyaan untuk pemilik proses.

**Blocked by:**

- None (can start immediately)


**Dasar:** EVIDENCED: `MEMORI_PEMAHAMAN.MD` §7.4 menyebut `1` Final dan `2` Partial/interim, **disimpulkan dari pemakaian**. `Rule-Obj-FieldValue` **tidak ada** di ekspor.

- [ ] Setiap nilai yang muncul di penjaga terdaftar beserta berkas dan barisnya. Hasil nihil tetap dicatat beserta pola yang dipakai.

**Ketidakpastian:** Sebaran nilai yang sesungguhnya dipakai diukur **REQ-027**. Bila hasilnya menunjukkan himpunan tertutup, `CHECK`-nya dipasang saat itu — satu `ALTER`, bukan pembongkaran.

## Hasil sapuan

- **Ketujuh nilai `PaymentType` EVIDENCED**, beserta nilai kosong — delapan kemungkinan, bukan tiga.
- **Penolakan saya atas usulan §4.1 dicabut.** Pola sapuan lama hanya mencakup literal berkutip dan melewatkan perbandingan tanpa kutip di `HitServiceToKasir_Act` baris 7754 dan 7896. Dicatat sebagai **D9**.
- Nilai instruksi bayar menambahkan **tepat satu** besaran menurut `PaymentType`, bukan menjumlahkan seluruhnya. Dicatat sebagai **D15**.
- **Arti** tiap nilai tetap TIDAK DITEMUKAN — definisi propertinya ada di schema PegaRULES, yaitu **REQ-001**.
