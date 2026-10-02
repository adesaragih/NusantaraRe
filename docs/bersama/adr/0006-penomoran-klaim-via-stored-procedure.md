---
status: accepted
tanggal: 2026-09-14
sumber: grilling Ronde 1 Q7, keputusan work owner
---

# Penomoran klaim Life memanggil `PROC_GENERATE_SEQUENCE_NUMBER`; dua SQL lama tidak dimigrasi

Nomor klaim Life **tidak dihitung ulang di aplikasi**. Sistem baru memanggil stored procedure
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` lewat jalur yang sama seperti sekarang. Dua rule penomoran
lama — `Generate_NoKlaim_Life` dan `Generate_NoKlaim_LifeRetro` — **sudah tidak dipakai
(di-remark)** dan **tidak dimigrasi**.

## Bukti

`[terverifikasi work owner 2026-09-14]` Kedua SQL lama sudah di-remark; penomoran sekarang lewat
`GetSequenceNumber_SQL` → `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`.

`[terverifikasi]` **Dikuatkan bukti korpus** — asimetri indeks rujukan rule di
`Claim Life/Activity/SaveOutStandingLife_Act.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`, 642.787 byte):

| RequestType | muncul sebagai `<RequestType>` | muncul di `pxRuleReferences` (`<pyRuleName>`) |
| --- | ---: | ---: |
| `Generate_NoKlaim_Life` | 1 | **0** |
| `Generate_NoKlaim_LifeRetro` | 1 | **0** |
| `GetSequenceNumber_SQL` | 1 | **2** |

Keduanya tertinggal sebagai parameter langkah tetapi **tidak terindeks sebagai rujukan aktif**.

```
f="Claim Life/Activity/SaveOutStandingLife_Act.xml"
grep -c '<pyRuleName>.*Generate_NoKlaim_Life<' "$f"        # 0
grep -c '<pyRuleName>.*GetSequenceNumber_SQL<' "$f"        # 2
```

**Koreksi artefak FASE A:** `discovery/flows/Claim Life.md` §3 mendaftar kedua rule lama sebagai
dipanggil — itu berdasarkan kehadiran `<RequestType>` saja. Status sesungguhnya: **tidak aktif**.
Koreksi ini sudah direkam di register OQ-002.

## Considered Options

- **Panggil `PROC_GENERATE_SEQUENCE_NUMBER`** — dipilih
- Replikasi logika penomoran di aplikasi — ditolak: **isi procedure tidak ada di korpus**
  (OQ-002), sehingga replikasi berarti menebak; dan penomoran yang bercabang di dua tempat berisiko
  menghasilkan nomor bentrok selama masa transisi
- Migrasikan dua SQL lama — ditolak: sudah tidak dipakai

## Consequences

- Sistem baru **bergantung pada Oracle** untuk penomoran klaim — ketergantungan runtime yang harus
  disadari saat merencanakan lingkungan dan pengujian.
- Format nomor klaim **tidak dinyatakan** di aplikasi; ia ditentukan procedure.
- Bila kelak penomoran dipindah ke aplikasi, itu keputusan baru yang menggantikan ADR ini.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-002** (dipersempit untuk Claim — Life, 2026-09-14) | **Kontrak `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`**: format nomor yang dihasilkan, dan apakah sequence di-reset per tahun / per lini / global. **Pemilik: DBA.** 66 stored procedure lain di korpus tetap terbuka |
| **OQ-013** | Apakah procedure melakukan `COMMIT` sendiri — batas transaksi di sisi database |

Pembanding yang menunjukkan kontrak semacam ini **bisa** terbaca bila SQL-nya bukan procedure:
format nomor polis treaty terbaca penuh di `NB Treaty In/RDBList/GenerateNoPolicy.xml`
(`ASM-FW-GISFW-INT-POLISTREATYIN` / `ASM!GENERATENOPOLICY` / `RULE-CONNECT-SQL`).
