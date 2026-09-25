---
status: accepted
label: DECIDED
---

# Data akseptasi disimpan dalam satu bentuk kanonik; sistem hilir diberi view, bukan salinan

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Tidak ada data yang disimpan dua kali. Bila sistem hilir memerlukan bentuk kolom, ia diberi **view atau projeksi** atas bentuk kanonik — bukan salinan yang ditulis terpisah.

Sistem lama menyimpannya dua kali. `POOLDATA.OS_AKSEPTASI_KLAIM` punya **66 kolom**, dan `PEGA_JSON_OS_AKSEP_KLAIMTNP` hanya mengisi **8** di antaranya (`CASEID`, `NOCLAIM`, `MASTERID`, `DATA_JSON`, `TANGGAL`, `NOPOLIS`, `STS_REJECT`, `STS_KONVERSI`). **59 kolom sisanya** — `GROSSVALUE`, `ADJUSTERFEE`, `KURSIDR`, `CNPREINSTATEMENT`, `TOTALXOL`, `SALVAGE`, dan seterusnya — diisi oleh sesuatu yang lain.

`DATA_JSON` adalah sumber kebenaran; kolom skalar adalah keluaran proses hilir.

## Consequences

**Temuan bypass Arasapas menghasilkan ramalan yang dapat diuji, dan itu memperkuat keputusan ini, bukan mengubahnya.** Untuk case yang dibuat `VINCENTVERNANDO_1`, konversi tidak berjalan (`FINDING-002` bagian 8.6). Maka barisnya seharusnya memiliki `DATA_JSON` terisi tetapi **kolom skalarnya kosong**.

Bila ternyata terisi, berarti ada penulis lain yang belum diketahui — dan seluruh pemahaman tentang siapa menulis apa ke tabel ini harus ditinjau ulang. Digabungkan ke **REQ-016**.

Siapa penulis dan siapa pembaca 59 kolom itu tetap **REQ-023**. Keputusan di atas **tidak menunggunya**: penyimpanan ganda ditolak apa pun jawabannya, karena dua salinan yang dapat menyimpang satu sama lain adalah cacat terlepas dari siapa yang menulisnya.
