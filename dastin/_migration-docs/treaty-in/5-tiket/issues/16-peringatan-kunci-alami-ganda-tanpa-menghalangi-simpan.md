---
status: aktif
golongan: pelestarian
---

# 16: Pengisi kontrak diperingatkan bahwa kunci alaminya sama dengan kontrak lain, tanpa dihalangi menyimpan

*Asal: `DAFTAR-PEKERJAAN.md` `P-02` · `SPEC-MODEL-DATA.md` §10.1 · `ADR-0040` §2.*

**What to build:** **PK** yang menyimpan kontrak berkunci alami sama dengan kontrak lain — cedant +
asal bisnis + tanggal mulai + sifat proporsi — **tetap dapat menyimpannya**, dan melihat peringatan
yang **menyebut kontrak pembandingnya**.

**Kenapa begini:** `ADR-0040` §2 memutuskan kunci alami kontrak **memperingatkan, tidak melarang** —
dan bedanya bukan selera. Kontrak yang benar-benar kembar memang terjadi: pembaruan tahunan yang
dicatat ulang, atau dua perjanjian terpisah dengan cedant yang sama pada tanggal yang sama.
Melarangnya berarti **menolak data yang sah**, yaitu kekeliruan lingkup yang sudah tiga kali terjadi
di modul ini. Dan peringatan tanpa menyebut pembandingnya adalah peringatan yang tidak dapat
ditindaklanjuti: yang membacanya harus dapat membuka kontrak yang dimaksud.

**Persyaratan:** `ADR-0040` §2 · `SPEC-MODEL-DATA.md` §10.1 *("kunci alami dan lingkupnya —
`ID_CEDANT` + `ID_ASAL_BISNIS` + `TANGGAL_MULAI` + sifat proporsi")* · **ketiadaan nomor `INV`**-nya
**disengaja** dan dinyatakan di §10.1 — ia bukan constraint

**Tidak termasuk:** **`UNIQUE` atas kunci alami kontrak** — tegas tidak dipasang, dan ketiadaannya
adalah **pernyataan keputusan** yang sudah berdiri di `ddl-usulan/Z00_KUNCI_ALAMI.sql`. Siapa pun yang
menambahkannya membalikkan `ADR-0040` §2 tanpa membukanya.
**Pembekuan kunci alami sesudah kontrak lahir** — irisan `18`. Yang di sini soal **lahirnya**, yang di
sana soal **mengubahnya**.

**Jalur gagal:** Menyimpan kontrak kembar -> **BERHASIL**, dan peringatannya muncul menyebut pengenal
kontrak pembandingnya · Peringatan muncul tanpa menyebut pembandingnya -> **kriteria selesai tidak
terpenuhi** · Penyimpanan ditolak -> **kriteria selesai tidak terpenuhi**, sebab itu perilaku yang
`ADR-0040` §2 larang.

**Uji:** **Negatif — dan di irisan ini "negatif" berarti sesuatu yang lain:** yang harus gagal bukan
penyimpanannya melainkan **kebisuan**. Simpan kontrak kembar dan periksa bahwa **peringatannya ada**;
ketiadaan peringatan adalah kegagalan.
**Positif:** dua kontrak yang **berbeda pada satu ruas kunci saja** — misalnya sifat proporsi —
disimpan **tanpa** peringatan. Tanpa uji ini, peringatan yang terlalu lebar akan berbunyi untuk
hampir setiap kontrak dan orang berhenti membacanya.

**Menggantikan:** `P-02` tidak menggantikan aturan sistem lama, sebab **sistem lama tidak punya
pemeriksaan kembar sama sekali** pada kunci ini. Yang ada hanya `CheckDuplicateOffer`, yang bekerja
atas **penawaran** dan berbunyi *"Have similarities in SoB, Ceding, and Period, please check for
duplicates"* — sumbu yang **mirip tetapi tidak sama**: ia tidak memuat sifat proporsi. Golongan
**PELESTARIAN** dipertahankan karena perilakunya — memperingatkan, bukan menolak — memang
dilestarikan; yang berubah **sumbunya**.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(CheckDuplicateOffer@ekspor-2026-09 - "Have similarities in SoB, Ceding, and Period")
        DECIDED(ADR-0040)
```

- [ ] menyimpan kontrak berkunci alami kembar **berhasil**, dan tidak ada constraint yang menolaknya
- [ ] peringatannya menyebut **pengenal kontrak pembandingnya**, bukan hanya menyatakan ada kembar
- [ ] uji positif lulus: beda pada satu ruas kunci saja -> **tidak ada peringatan**
- [ ] ketiadaan `UNIQUE` atas kunci alami kontrak tertulis sebagai **pernyataan keputusan**, dan sapuan atas `ddl-usulan/` membuktikannya masih berupa pernyataan — bukan sudah dipasang diam-diam
- [ ] perbedaan sumbu terhadap `CheckDuplicateOffer` **tercatat**: sistem lama tidak memuat sifat proporsi
