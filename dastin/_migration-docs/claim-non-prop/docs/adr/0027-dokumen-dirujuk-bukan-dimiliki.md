---
status: accepted
label: DECIDED
---

# Dokumen klaim dirujuk, tidak disimpan ulang

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); `pengetahuan/DDL_Script_ClaimNonProp.xls` versi 2026-09-18 10:36 (48 objek); dan `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` (ditempel pengguna 2026-09-18).

Sistem baru **merujuk** dokumen klaim, tidak menyimpannya. Tidak ada migrasi berkas fisik.

## Consequences

`MEMORI_PEMAHAMAN.MD` §5.7: berkas sesungguhnya berada di **Google Cloud Storage**. `POOLDATA.T_STORAGE_IMAGE` hanya menyimpan **metadata dan URL beserta `EXPDATE`**, yang disegarkan lewat `GetUrlGoogleStorage_Act`.

Jadi kekhawatiran tentang volume berkas yang harus dipindahkan **gugur** — tidak ada berkas fisik di dalam basis data.

Yang tetap perlu dirancang: URL berumur terbatas (`EXPDATE`), sehingga sistem baru memerlukan mekanisme penyegaran yang setara dan tidak boleh menyimpan URL sebagai nilai tetap.

`POOLDATA.GET_TOKEN_STORAGE` dan class `Link-Attachment` termasuk lapisan rujukan ini, bukan lapisan penyimpanan.
