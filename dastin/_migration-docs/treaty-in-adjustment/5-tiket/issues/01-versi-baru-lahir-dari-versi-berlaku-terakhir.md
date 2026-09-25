---
status: aktif
---

# 01: Versi baru lahir dari versi berlaku terakhir, dan rujukan dasarnya disimpan eksplisit

*Asal: `DAFTAR-PEKERJAAN.md` `P-42` (diubah 24 Sep 2026) · `SPEC-MODEL-DATA.md` §10.2 · `GRL-10`, bahan to-spec `B-3` dan `B-4`.*

**What to build:** **PK** membuat versi penyesuaian atas sebuah kontrak, dan versi itu lahir `DRAFT` sambil
menyimpan **rujukan eksplisit** ke versi yang menjadi dasarnya. Dasarnya **bukan baris yang dipilih
di layar** dan **bukan versi `DISETUJUI` mana pun** — ia **versi berlaku terakhir** pada saat versi
baru dibuat.

Artefak: kolom `ID_VERSI_KONTRAK_DASAR` pada `VERSI_KONTRAK`, kunci asingnya, dan constraint yang
menjaga keterisiannya.

**Persyaratan:** `INV-04` · bahan to-spec `B-3` (wajib terisi pada versi penyesuaian, wajib kosong pada versi pertama) · `B-4` (yang ditunjuk harus versi berlaku terakhir, bukan `DITOLAK` maupun `DIBATALKAN`) · `ADR-0040`

**Tidak termasuk:** **Penomoran ulang baris warisan** — itu irisan 10. Baris warisan boleh punya dasar kosong
sampai irisan itu selesai.
**Turunan "versi berlaku"** — `GRL-11` menetapkan ia dihitung, bukan disimpan; tidak ada kolom yang
dibangun untuknya di sini.

**Jalur gagal:** Menunjuk versi berkeadaan `DITOLAK` sebagai dasar -> **ditolak saat simpan**, dan pesannya
menyebut keadaan versi yang ditunjuk · Versi pertama sebuah kontrak dengan dasar terisi -> ditolak ·
Versi penyesuaian dengan dasar kosong -> ditolak.

**Uji:** **Negatif:** simpan versi penyesuaian tanpa dasar; simpan versi pertama dengan dasar; tunjuk
versi `DITOLAK`; tunjuk versi `DIBATALKAN`.
**Positif — dan ia yang menangkap lingkup yang terlalu sempit:** kontrak yang punya **dua** versi
`DISETUJUI` dalam sejarahnya menerima versi baru yang menunjuk **yang terakhir**, bukan ditolak
karena ada dua.

**Menggantikan:** `P-42` bunyi lama: *"**PK** dapat membuat versi baru **dari versi yang sudah `DISETUJUI`**,
dan versi baru itu mulai dari `DRAFT` serta melewati keempat tingkat."*

Yang bergeser: dasarnya menjadi **versi berlaku terakhir**, dan rujukannya **disimpan eksplisit**.
Bunyi lama mengizinkan memilih versi `DISETUJUI` yang bukan terakhir — dan selisih yang dihitung
terhadap dasar yang salah **tidak menghasilkan galat**. Sistem lama: `OLDID` = ID baris yang
**dipilih di picker**, tanpa pembeda jenis dan tanpa saringan keadaan (`TDA-11`).

**Blocked by:** **14** — *ditambahkan 24 September 2026.* Tiket ini menambahkan **kolom** atau **tabel anak** pada `VERSI_KONTRAK`, dan tidak satu pun tiket di papan membuat tabelnya. Irisan `14` adalah **PEMBUAT PERTAMA** `KONTRAK` dan `VERSI_KONTRAK`.

**Dasar:**
```
EVIDENCED(TreatyInRevisi_post@ekspor-2026-09, TreatyLoadMasterJoinEdm@ekspor-2026-09)
        DECIDED(GRL-10, ADR-0040)
        DIASUMSIKAN-CLEAR(REV-3)
```

- [ ] `ID_VERSI_KONTRAK_DASAR` berdiri di `VERSI_KONTRAK` dengan kunci asingnya
- [ ] versi penyesuaian tanpa dasar **ditolak**, versi pertama dengan dasar **ditolak**
- [ ] menunjuk versi `DITOLAK` atau `DIBATALKAN` **ditolak**, pesannya menyebut keadaannya
- [ ] uji positif lulus: kontrak berdua versi `DISETUJUI` menerima versi baru
- [ ] `REV-3` tercatat di `ASUMSI-CLEAR.md` sebagai asumsi yang tiket ini bersandar padanya
