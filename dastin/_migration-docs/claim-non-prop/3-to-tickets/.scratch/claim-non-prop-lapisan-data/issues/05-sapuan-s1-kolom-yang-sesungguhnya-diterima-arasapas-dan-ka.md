---
status: selesai
---

# 05: Sapuan S1: kolom yang sesungguhnya diterima Arasapas dan kasir

> **SELESAI 19 September 2026.** Sisi Arasapas dan sisi kasir keduanya terpetakan; hasilnya `SAPUAN-CELAH-DAN-JAVA.md` §5.1.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

> **DIBUKA KEMBALI 2026-09-18 — sapuan ulang menemukan cakupannya kurang, bukan polanya salah.**
> Pola literal tiga-bentuk tidak mengubah satu pun dari 17 parameter `InputParamOs.*`. Yang kurang adalah **sisi kasir**: `HitServiceToKasir_Act` memakai 25 medan bernama urut `TempKasir.CARI1`…`CARI24` dan `CARIDATETIME`, dan tidak satu pun pernah diinventarisasi. Judul tiket ini menyebut **“Arasapas dan kasir”**; yang terpetakan baru Arasapas.
> Lima berkas juga tidak disebut laporan S1: `SaveToOS`, `GetSelisihActual_Act`, `CloseClaimTNonProp`, dan dua layar uji `Hitung_Test` yang mengikat medan langsung ke `InputParamOs.GrossValue`.
> Dicatat sebagai **D36**. Rinci di `SAPUAN-S10-S16.md` §1.

*Asal: `T-32` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** 17 parameter `InputParamOs.*` dan 25 parameter `TempKasir.CARI*` terbaca, beserta asal tiap satunya.

**Blocked by:**

- None (can start immediately)


**Dasar:** EVIDENCED: `Activity\SaveDataToOSAksep_Act.xml`, `Activity\HitServiceToKasir_Act.xml`. Hasil lengkap di `SPEC-MODEL-DATA.md` bagian 18.

- [ ] Terpenuhi. Temuan yang mengubah T-16: keempat nilai uang dikirim sebagai **selisih**, bukan nilai penuh.

**Ketidakpastian:** `InsertOSKlaimCNP` **TIDAK DITEMUKAN** sebagai nama berkas; `KonversiKlaim_Act.xml` ada dengan **nol** pasangan properti; `ConnectREST\KonversiKlaimNonLife.xml` belum dibaca isinya. → T-37.

## Hasil

- **Sisi Arasapas**: 17 parameter `InputParamOs.*` — tidak berubah oleh sapu ulang pola tiga-bentuk.
- **Sisi kasir**: **21 medan bernama dari 25 slot** `TempKasir.CARI*`. Nama bisnisnya terbaca dari 480 baris Java yang belum pernah dibaca sebagai kode. `CARI19`–`CARI21` hanya perakit tanggal; `CARIDATETIME` penampung sementara. **D46**.
- **`Nett`, `Deductible`, `KaliDeduct` dikirim tanpa kutip** sebagai angka JSON, dari `String` hasil `.toString()`. Nilai kosong menghasilkan JSON rusak. **D41**.
- **Tujuh nilai tertanam di kode**, termasuk `StsSyariah = 0` dan `CompanyName = "NUSARE"` — jalur pembayaran tidak pernah menyatakan syariah. **D42**.
- Satu hal **tidak** disimpulkan: `CARI20 = @toDecimal(@month(...)) + 1`. Perilaku `@month()` tidak ada di ekspor — diajukan sebagai **REQ-036**, bukan ditebak.
