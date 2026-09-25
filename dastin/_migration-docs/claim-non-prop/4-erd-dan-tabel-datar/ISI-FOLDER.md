# Struktur — ERD dan tabel datar

Bentuk data, digambar. Semuanya **TURUNAN**: dibangkitkan alat di `alat/` dan tidak pernah disunting tangan. Bila salah satu berbeda dari sumbernya, **alatnya yang salah**.

**9 berkas** (termasuk berkas ini). Dipindahkan ke sini 2026-09-20 dari akar folder modul.

| Berkas | Isinya |
|---|---|
| `struktur-claimdata-lama.md` | struktur `.ClaimData` sistem lama — **sumber** bagi ERD di folder ini |
| `Diagram-Skema-Tabel-NusantaraRe.xlsx` | rancangan lini lain — Claim Life, PremiumList, Claim Prop, Komite Prop. **Rujukan, tidak disunting** |
| `Diagram-Skema-Tabel-ClaimNonProp.xlsx` + `ERD-CLAIM-NON-PROP.html` | rancangan Claim Non Prop — 30 tabel, 27 relasi, hanya PK dan FK |
| `Diagram-Skema-Tabel-Gabungan.xlsx` + `ERD-GABUNGAN.html` | kedua berkas di atas disatukan — 51 tabel, 44 relasi pohon, 11 tali penghubung |
| `ERD-ORACLE.xlsx` + `ERD-ORACLE.md` | ERD **tabel Oracle** sistem lama dan skema baru, berikut bentuk datarnya per kolom |

> **Penamaan tabel berubah 20-09-2026.** Awalan `T_CLAIMNP_` dihapus; yang berlaku `T_CLAIM_`, mengikuti keputusan work owner yang tercatat di `Diagram-Skema-Tabel-NusantaraRe.xlsx`. Peta nama lama → baru, sebab tiap pemecahan, dan daftar lini pemakainya ada di [`PETA-NAMA-TABEL.md`](PETA-NAMA-TABEL.md). Nama di `ddl-usulan/` **tidak** berubah — penerjemahnya `peta-nama-tabel.tsv` dan sheet **PETA-NAMA-T** di `ERD-ORACLE.xlsx`.

---

Peta seluruh berkas modul dan folder barunya: [`../PETA-FOLDER.md`](../PETA-FOLDER.md).
Register keadaan: [`../STATUS.md`](../STATUS.md) · [`../KEADAAN-AKHIR.md`](../KEADAAN-AKHIR.md).
