---
status: accepted
label: DECIDED
---

# Batas kepemilikan mengikuti nama class: -Work- dan -Data- dimiliki, -Int- tidak

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); `pengetahuan/DDL_Script_ClaimNonProp.xls` versi 2026-09-18 10:36 (48 objek); dan `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` (ditempel pengguna 2026-09-18).

Modul klaim **memiliki** klaim, akseptasi, alokasi, dan adjustment. Selebihnya **dibaca, tidak dimiliki**.

Aturannya tidak perlu ditimbang per tabel — `MEMORI_PEMAHAMAN.MD` §2.3 sudah memberinya: **integration class (`-Int-`) adalah pemetaan langsung ke tabel atau view Oracle.**

| Pola class | Kepemilikan |
|---|---|
| `…-Work-…` | **dimiliki** |
| `…-Data-…` | **dimiliki** |
| `…-Int-…` | **dibaca, tidak dimiliki** |

Dengan itu `V_POLIS`, `T_STORAGE_IMAGE`, `EMAILKOMITE`, dan `M_LINK_SERVICE` seluruhnya di luar kepemilikan — tanpa perlu memutuskan satu per satu.

## Consequences

**Tabel yang bukan milik klaim tidak dimigrasikan, meskipun ada di daftar tarikan.** Daftar itu untuk **memahami**, bukan untuk memindahkan. Ini memotong lingkup migrasi data secara langsung.

Bila modul lain belum siap pada saat cutover, klaim **membaca dari basis data lama lewat antarmuka yang disepakati** — bukan menyalin datanya. Konsekuensinya cutover dapat dilakukan per modul, tidak harus sekaligus.

Yang tetap perlu dipahami meski tidak dimiliki: 298 kolom di DDL tanpa pasangan properti Pega (`BLUEPRINT.md` §19.2) sebagian besar milik modul lain. Memahaminya perlu; memindahkannya tidak.
