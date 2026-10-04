---
status: aktif
golongan: pelestarian
---

# 30: Dokumen dilampirkan pada kontrak, dan dokumennya dirujuk bukan disalin

*Asal: `DAFTAR-PEKERJAAN.md` `P-15` · `SPEC-MODEL-DATA.md` §10.19 · `ADR-0027` · `SPEC-INVARIAN.md` `INV-59`, `INV-67`.*

**What to build:** **PK** melampirkan dokumen pada sebuah versi kontrak — slip, wording, surat — dan
yang tersimpan adalah **rujukan ke dokumennya**, bukan salinan isinya.

Artefak: entitas `DOKUMEN_KONTRAK`, kunci asingnya ke `VERSI_KONTRAK`, dan `INV-67`.

**PEMBUAT PERTAMA** untuk `DOKUMEN_KONTRAK`.

**Kenapa begini:** Menyalin isi berkas ke dalam basis data melahirkan **dua sumber kebenaran untuk
satu dokumen**, dan yang kedua tidak pernah ikut berubah ketika yang pertama diganti — persis yang
`INV-59` larang. `ADR-0027` sudah memutuskan lampiran **dirujuk, tidak dimiliki**; sistem lama pun
demikian, lewat `M_ATTACHMENTTREATY_2`. Yang berubah hanya bahwa aturannya kini **tertulis dan
diperiksa**.

**Persyaratan:** `ADR-0027` · `INV-59` (tidak ada nilai yang disalin dari entitas lain; hilir diberi
rujukan) · `INV-67` (dokumen unik di dalam satu versi) · `INV-61` **tidak berlaku di sini** — ia
tentang arsip JSON, irisan `42`

**Tidak termasuk:** **`DOKUMEN_ADDENDUM`** — entitas yang **berbeda**, milik irisan `04`, dan
namanya sengaja dibedakan: `DOKUMEN_KONTRAK` berstatus *"dirujuk, tidak dimiliki"* untuk lampiran,
sementara `DOKUMEN_ADDENDUM` adalah **benda yang berdiri sendiri di atas versi**, dengan nomor yang
beredar di luar sistem. Menyatukannya membuat dua arti pada satu nama.
**Penyimpanan berkasnya** — di luar skema; yang di sini rujukannya.
**Arsip JSON sistem lama** — irisan `42`, dan ia benda lain lagi.

**Jalur gagal:** Dua baris merujuk dokumen yang sama pada satu versi -> **ditolak** `INV-67` ·
Sebuah kolom yang memuat isi berkas -> **tidak ada**, dan ketiadaannya diperiksa dengan sapuan atas
`KAMUS-KOLOM.md` · Nama dokumen disalin dari sumbernya sebagai teks yang kemudian basi -> dilarang
`INV-59`.

**Uji:** **Negatif:** lampirkan dokumen yang sama dua kali pada satu versi.
**Positif — dan ia yang menangkap kunci yang terlalu ketat:** **dokumen yang sama dilampirkan pada
dua versi berbeda** dari kontrak yang sama -> **diterima**. Wording yang tidak berubah antar versi
adalah keadaan yang lazim; `UNIQUE` yang ditulis atas dokumen saja — tanpa versinya — menolaknya, dan
lulus uji negatif di atas.

**Menggantikan:** `P-15` melestarikan `M_ATTACHMENTTREATY_2`. Yang bergeser: keunikannya
**diperiksa**, dan rujukan-bukan-salinan **dinyatakan sebagai invarian** alih-alih dibiarkan sebagai
kebiasaan.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(M_ATTACHMENTTREATY_2@Table/ - lampiran sebagai rujukan)
        DECIDED(ADR-0027, INV-59, INV-67, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `DOKUMEN_KONTRAK` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-67` terpasang atas **dokumen di dalam satu versi** — bukan atas dokumen saja
- [ ] uji positif lulus: satu dokumen melekat pada dua versi berbeda **diterima**
- [ ] sapuan membuktikan **tidak ada** kolom pembawa isi berkas
- [ ] perbedaannya dari `DOKUMEN_ADDENDUM` **tertulis di dalam tiket ini**, supaya dua nama tidak berarti satu benda
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
