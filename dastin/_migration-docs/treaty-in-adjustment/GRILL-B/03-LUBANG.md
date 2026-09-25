> Modul  : Treaty In Adjustment · Ronde B · 2026-09-27
> Status : TERBUKA
> Sifat  : TAMBAH-SAJA

# 03 · LUBANG RONDE B

## 1. Pemasangan dua sumbu lama ke satu sumbu ADR-0049 — bahan `C2`

Ini yang saya lewatkan pada pengajuan `C2` pertama, dan ia justru isi pertanyaannya.

### 1.1 Ketiga jenis ADR-0049, dikutip apa adanya

> **Satu sumbu.** Jenis addendum punya tiga nilai:
>
> | Jenis | Materialitas turunannya |
> |---|---|
> | perubahan estimasi | material |
> | penyesuaian ke nilai aktual | material |
> | administratif | non-material |
>
> **Kolom materialitas tidak ada** di model baru. Ia diturunkan dari jenis.

**Internal dan External tidak muncul sama sekali** di ADR-0049 — bukan sebagai jenis, bukan sebagai
catatan. Kosakatanya sepenuhnya berbeda dari kosakata sistem lama.

### 1.2 Pemasangan kelima kombinasi sah ke ketiga jenis

| Kombinasi lama | Nama layarnya | Jenis ADR-0049 yang sepadan | Dasarnya |
|---|---|---|---|
| **(3, 1)** | Addendum Premi + Material | **penyesuaian ke nilai aktual** | **terbaca** — `TreatyInSetEditPre` 3.5 menulis `ActualValue.EGNPI <- EGNPI`, yaitu menyalin ke "nilai aktual" secara harfiah |
| **(1, 2)** | Internal + Non Material | **administratif** | disimpulkan dari materialitasnya, bukan dari perilaku |
| **(2, 2)** | External + Non Material | **administratif** | idem — **dua kombinasi lama jatuh ke satu jenis baru** |
| **(1, 1)** | Internal + Material | **perubahan estimasi**? | **tidak ada dasar** — tidak ada apa pun di ekspor yang menyebut estimasi |
| **(2, 1)** | External + Material | **perubahan estimasi**? | idem |

**Kemiripan jumlah adalah kebetulan** (`METODE` §3.4): lima kombinasi sah dipasangkan ke tiga jenis,
dua di antaranya bergabung, dan dua sisanya dipasangkan ke sebuah jenis yang **tidak punya padanan
perilaku apa pun** di sistem lama.

### 1.3 Sesudah GRL-12, himpunan ADR-0049 sendiri ikut runtuh

ADR-0049 mendefinisikan **administratif** sebagai jenis yang materialitas turunannya
**non-material**. Sesudah GRL-12, non-material **sudah** menjadi turunan — versi tanpa baris
`NILAI_SELISIH`. Maka:

* **"administratif" berhenti menjadi nilai jenis.** Ia sama dengan non-material turunan, dan
  menyimpannya sebagai nilai enum berarti menyimpan satu fakta di dua tempat (ADR-0041).
* **Yang tersisa dari ADR-0049 hanya dua:** *perubahan estimasi* versus *penyesuaian ke nilai
  aktual*. Keduanya material, dan bedanya **jenis perubahannya**, bukan akibatnya.
* Dan pembedaan yang tersisa itu **persis bentuk (b)** yang saya tawarkan: penyesuaian premi versus
  lainnya — sebab (3,1) adalah satu-satunya kombinasi yang terbaca sebagai "penyesuaian ke nilai
  aktual".

**Jadi ADR-0049 sendiri, dibaca sesudah GRL-12, mengarah ke dua nilai — bukan tiga.** Apa pun
himpunan yang dipilih, ia **berbeda dari yang tertulis di ADR-0049**, sehingga butuh **`REV-5`**.

### 1.4 Dua argumen saya yang dicabut

| Argumen | Kenapa dicabut |
|---|---|
| *"membuangnya tidak dapat dibatalkan"* | **Salah.** `EDMState` warisan dibawa apa adanya sebagai nilai warisan dengan asal-usulnya — pola yang sama dengan GRL-06, GRL-09, GRL-12. **Tidak ada data lama yang hilang.** Yang dipertaruhkan hanya versi **baru** yang lahir sebelum `DB-16` terjawab: bila pembedaannya dibuang lalu ternyata bermakna, versi-versi itulah yang kehilangan klasifikasinya. **Itu** alasan sah memilih sisi yang dapat dibatalkan |
| *"`EDMEffective` disemai hanya untuk jenis 1"* | **Melanggar `METODE` §3.8.** Penyemaian itu hanya berarti bagi **pro rata**, dan pro rata **mati** (10 dari 53 langkah). Membenarkan sebuah pembedaan dengan perilaku kemampuan yang mati sama dengan menghitung kemampuan mati sebagai kemampuan berjalan. Dibuang dari daftar alasan, dan **dititipkan ke `E3a`** sebagai catatan: bila pro rata dibangun, penyemaian `EDMEffective` per jenis harus diputuskan bersamanya |

## 2. `DB-16` dipertajam, dan dikaitkan ke tiga butir lain

Bunyi lama terlalu longgar. Bunyi baru, konkret dan dapat dibantah (`METODE` §4.7):

> **`DB-16`.** *"Addendum eksternal selalu berupa dokumen yang disepakati dan ditandatangani
> cedant, dengan tanggal berlaku sendiri; revisi internal tidak pernah dikirim ke luar."*

**Bila External memang berarti dokumen dari luar, empat butir menyangkut hal yang sama** dan
sebaiknya ditanyakan bersama:

| Butir | Isinya | Hubungannya |
|---|---|---|
| **`DB-3`**, **`DB-4`** | satu dokumen addendum = satu kontrak; dua addendum tidak digabung | keduanya mengandaikan addendum **adalah dokumen** |
| **`DB-5`** | addendum tidak pernah berlaku di tengah periode | "tanggal berlaku sendiri" pada `DB-16` |
| **`DB-11`** | nomor addendum beredar di luar sistem | dokumen yang dikirim membawa nomornya |

Bila `DB-16` dibenarkan, keempatnya berbicara tentang **satu benda**: dokumen addendum eksternal
yang keluar dari sistem, bertanggal sendiri, bernomor, dan disepakati dua pihak. Bila dibantah,
keempatnya melemah bersama-sama.
