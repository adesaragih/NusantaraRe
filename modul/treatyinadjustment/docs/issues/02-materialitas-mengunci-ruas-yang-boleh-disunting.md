---
status: aktif
---

# 02: Materialitas mengunci ruas mana yang boleh disunting, dan penolakannya ditegakkan di sisi simpan

*Asal: `DAFTAR-PEKERJAAN.md` `P-60` · `GRL-20` · `SPEC-MODEL-DATA.md` §10.2a.*

**What to build:** **PK** menyatakan materialitas sebuah versi — `MATERIAL` atau `TIDAK_MATERIAL` — dan
pernyataan itu **menentukan ruas mana yang boleh disunting**. Penolakannya berdiri **di sisi
simpan**, bukan hanya di layar.

Artefak: kolom `SIFAT_MATERIAL_ADDENDUM`, dan penegakan per-ruas yang menolak penyimpanan.

**PEMBUAT PERTAMA** untuk `SIFAT_MATERIAL_ADDENDUM`.

**Persyaratan:** `GRL-20` · `SPEC-MODEL-DATA.md` §10.2a · `ADR-0037` butir 3 (besaran yang boleh disepakati dinaikkan menjadi masukan)

**Tidak termasuk:** **`INV-69` dan `INV-70`** — penegakan atas baris `NILAI_SELISIH` — itu irisan 11, dan ia
menunggu `NILAI_SELISIH` lahir di irisan 06.
**Titik beku materialitas** — itu irisan 07.
**Nilai warisan** `EDMMaterialType` — dibawa apa adanya oleh jalur migrasi, bukan oleh irisan ini.

**Jalur gagal:** Versi `TIDAK_MATERIAL` yang penyimpanannya mengubah ruas terkunci -> **ditolak**, dan
pesannya menyebut **ruas mana** yang terkunci · Versi `MATERIAL` yang mengubah ruas yang sama ->
diterima.

**Uji:** **Negatif:** simpan versi `TIDAK_MATERIAL` yang mengubah ruas uang terkunci.
**Positif — dan ia yang menangkap penguncian yang terlalu lebar:** versi `TIDAK_MATERIAL` yang
mengubah ruas yang **memang boleh** berubah **diterima**. Ketiga kekeliruan lingkup di modul ini
lolos uji negatif; yang menangkapnya hanya uji positif.

**Menggantikan:** `TDA-10` — *materialitas hanya ditegakkan di layar; nol penegakan di sisi simpan.*
Sistem lama memasang kondisi penguncian di badan seksi dan **tidak memeriksa apa pun saat
menyimpan**, sehingga nilai yang seharusnya terkunci tetap dapat masuk lewat jalur lain.

**Blocked by:** **14** — *ditambahkan 24 September 2026.* Tiket ini menambahkan **kolom** atau **tabel anak** pada `VERSI_KONTRAK`, dan tidak satu pun tiket di papan membuat tabelnya. Irisan `14` adalah **PEMBUAT PERTAMA** `KONTRAK` dan `VERSI_KONTRAK`.

**Dasar:**
```
EVIDENCED(TreatyInEDMSetValue@ekspor-2026-09, ekspor-tambahan/EDMMaterialType.xml)
        DECIDED(GRL-20, ADR-0037)
        DIASUMSIKAN-CLEAR(DB-20)
```

- [ ] `SIFAT_MATERIAL_ADDENDUM` berdiri dengan dua nilainya, boleh kosong pada versi pertama
- [ ] daftar ruas yang dikunci tiap nilai **tertulis**, dan diturunkan dari kondisi layar lama — bukan disusun ulang
- [ ] simpan yang melanggar **ditolak**, pesannya menyebut ruas yang terkunci
- [ ] uji positif lulus: ruas yang boleh berubah tetap dapat berubah
- [ ] `DB-20` tercatat di `ASUMSI-CLEAR.md`
