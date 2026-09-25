---
status: accepted
label: DECIDED
---

# Pengenal polis dan klaim diberi nama menurut isinya; CASEID dipensiunkan

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Istilah `CASEID` **dipensiunkan** dari model internal dan masuk daftar `_Avoid_` di `CONTEXT.md`. Setiap pengenal diberi nama yang menyebut isinya dan pemiliknya.

Sebabnya: nama yang sama membawa arti berbeda di dua tempat. Di `V_POLIS`, `CASEID` = `DATA_JSON.IDNewBisnis`, identitas **polis**. Di `OS_AKSEPTASI_KLAIM`, `CASEID` berbentuk `'ASM-FW-GCNMFW-WORK CLMNP-…'`, identitas **case Pega**.

## Peta pengenal — lebih dari dua

Penamaan tidak ditetapkan sebelum peta ini lengkap, karena berhenti di dua akan mengulang kesalahan yang sama dalam bentuk lebih halus.

| Pengenal | Ditemukan di | Dugaan pemilik | Status |
|---|---|---|---|
| `CASEID` (V_POLIS) = `DATA_JSON.IDNewBisnis` | view `V_POLIS` | kita | `IDNewBisnis` **nol kemunculan** di 279 XML — hanya ada di view |
| `CASEID` (OS_AKSEPTASI_KLAIM) | tabel akseptasi | Pega | berbentuk `<class> <pyID>` |
| `NOPOLIS` | `JSON_KLAIM`, `JSON_POLIS`, `OS_AKSEPTASI_KLAIM`, `TREATYINPRODUCTION`, 2 prosedur | kita | kolom, 4 tabel |
| `.ClaimData.PolicyNo` / `.ClaimData.PolicyData.PolicyNo` | XML, 115 kemunculan | kita | apakah keduanya sama **belum terbukti** |
| `QuotationData.PolicyNoSinarmas` | XML | **perusahaan lain dalam grup** | bentuk keempat, menyebut nama perusahaan |
| `POLICY_CEDING` | tabel work Pega | **cedant** | nomor polis milik cedant |
| `DLANO_CEDING` / `DLANO_SOB` | `OS_AKSEPTASI_KLAIM` | cedant / SOB | nomor DLA, bukan nomor polis |

**Sedikitnya lima bentuk pengenal polis hidup berdampingan, dan sedikitnya dua di antaranya milik pihak lain.** Nomor polis cedant dan nomor polis kita adalah dua hal berbeda; `POLICY_CEDING` memperlihatkan keduanya memang disimpan berdampingan.

Pemetaan mana yang sama dan mana yang berbeda **belum selesai** dan menunggu DDL tabel work serta REQ-023.

## Consequences

**Penggantian nama berhenti di batas integrasi.** Payload REST ke Arasapas dan ke sistem kasir membawa `CASEID` apa adanya. Nama internal boleh berubah; **kontrak ke luar tidak**. Pemetaan nama internal ke nama payload ditulis eksplisit di Boundary Contract, supaya penggantian istilah tidak merambat ke payload dan memutus integrasi yang selama ini berjalan.

Satu nama untuk dua hal tidak menimbulkan galat, hanya salah paham — dan salah pahamnya baru muncul saat dua orang membicarakan hal berbeda dengan kata yang sama.
