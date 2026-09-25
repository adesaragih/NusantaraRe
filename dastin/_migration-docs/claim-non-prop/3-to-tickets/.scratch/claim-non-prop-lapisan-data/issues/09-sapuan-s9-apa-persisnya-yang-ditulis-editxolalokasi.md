---
status: selesai
---

# 09: Sapuan S9: apa persisnya yang ditulis `EditXOLAlokasi`

> **SELESAI 2026-09-18.** Dijalankan 18 September 2026 atas empat lapisan; hasilnya `SAPUAN-S3-S9.md` bagian S9.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-36` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Daftar properti yang benar-benar dapat disunting petugas.

Seluruh `Property-Set` di rule itu, per properti. Empat lapisan.

**Blocked by:**

- None (can start immediately)


**Dasar:** EVIDENCED: `BLUEPRINT.md` §2.5 — `EditXOLAlokasi` menulis `.IsEditClaim = 1` pada class `ASM-FW-GISFW-Data-SpreadingRisk`, yang di sistem lama memuat baris layer **dan** baris Retensi Cedant.

- [ ] Daftar properti beserta baris; dan pernyataan tegas properti mana yang **tidak** tersentuh.

**Ketidakpastian:** **Menahan T-05.** Memasang `_HITUNG`/`_SUNTING` di kolom yang tidak pernah disunting adalah beban mati; tidak memasangnya di kolom yang disunting berarti membongkar tabel kemudian (ADR-0008).

## Hasil sapuan

- `EditXOLAlokasi` menulis **tiga** hal saja: pesan galat lokal, penanda `.IsEditClaim = 1`, dan satu catatan kronologi. **Tidak satu pun properti bernilai uang.**
- Ia **gerbang kata sandi**, bukan penyunting. Layar dan harness-nya memuat dua medan saja, `CARI1` dan `CARI2`.
- **Pertanyaan `_HITUNG`/`_SUNTING` tidak terjawab** — bukan karena sapuannya gagal, melainkan karena tertuju pada rule yang keliru. Sasaran penggantinya belum ditetapkan; dicatat di daftar yang tetap tertutup.
- Temuan tambahan: `.IsEditClaim` ditulis tiga kali dan **dibaca nol kali**. Dicatat sebagai **D10**.
