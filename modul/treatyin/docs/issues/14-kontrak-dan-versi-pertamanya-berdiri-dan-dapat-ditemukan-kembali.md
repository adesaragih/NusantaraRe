---
status: aktif
golongan: pelestarian
---

# 14: Kontrak dan versi pertamanya berdiri, dan dapat ditemukan kembali dengan pengenalnya

*Asal: `DAFTAR-PEKERJAAN.md` `P-01` (**pecahan identitas** — lihat `Tidak termasuk`) · `SPEC-MODEL-DATA.md` §10.1, §10.2 · `ADR-0040`.*

**What to build:** **PK** membuat sebuah kontrak baru, mengisi kepala versinya, menyimpannya, dan
**menemukannya kembali dengan pengenal yang sistem berikan**. Identitasnya terpecah dua sebagaimana
`ADR-0040` menetapkan: `KONTRAK` memegang lapisan beku — cedant, asal bisnis, sifat proporsi,
periode — dan `VERSI_KONTRAK` memegang segala yang dapat berbeda antar versi.

Artefak: tabel `KONTRAK` dan `VERSI_KONTRAK` seluruhnya, kunci utamanya, kunci asing di antara
keduanya, dan sequence pengenalnya.

**PEMBUAT PERTAMA** untuk `KONTRAK` dan `VERSI_KONTRAK`.

**Kenapa begini:** **Tidak ada satu pun tiket di papan yang membuat `VERSI_KONTRAK`.** Tiket `01`,
`02`, dan `05` menambahkan **kolom** padanya; tiket `04` menambahkan tabel yang merujuknya. Keempatnya
mengandaikan tabelnya sudah ada, dan pembuatnya tidak pernah ditiketkan — sebab ia milik `P-01`,
kemampuan Treaty In, yang ronde tiket sebelumnya memang di luar lingkup. Irisan ini menutup lubang
itu, dan karena itu **seluruh papan menggantung padanya**.

**Persyaratan:** `INV-01` (setiap tabel berkunci utama) · `INV-02` (pengenal dari `SEQUENCE`, tidak
pernah dari cap waktu maupun teks) · `INV-03` (pengenal tidak dipakai ulang, sequence tanpa `CYCLE`) ·
`INV-04` (`NOMOR_URUT_VERSI` unik di dalam satu `KONTRAK`) · `INV-18` (perilaku hapus tiap kunci asing
ditetapkan sadar) · `INV-29` (`SIFAT_PROPORSI` dua nilai) · `INV-53` (`TANGGAL_MULAI` ≤
`TANGGAL_BERAKHIR`, keduanya inklusif) · `ADR-0040` · nama dan tipe kolom **mengikat** pada
`2-to-spec/KAMUS-KOLOM.md` (urutan wewenang butir 5)

**Tidak termasuk:** **Daftar nilai sah `KEADAAN_SIKLUS_HIDUP` dan mesin perpindahannya.** Kolomnya
berdiri di sini — tabelnya tidak dapat berdiri tanpanya — tetapi `INV-20`, `INV-22`, `INV-23`,
`INV-24`, dan `INV-25` **tegas bukan bagian irisan ini**. Ia batch 2, sebab `REV-3` merevisi ADR-0055
§4 yang menjadi sumber daftar keadaan itu. **Ini pecahan yang disengaja**, dan bunyi lengkap `P-01`
memuat keduanya.
**Peringatan kunci alami ganda** — irisan `16`. **Pencarian lewat nomor warisan** — irisan `17`.
**Pembekuan kunci alami** — irisan `18`.
**Mesin pro rata** — `GRL-15` memutuskan ia **sengaja tidak dibangun**; atribut *"berlaku sejak"*
dibawa, mesinnya tidak. Itu **pernyataan keputusan**, bukan pekerjaan tertunda.

**Jalur gagal:** Menyimpan versi tanpa kontrak induk -> ditolak kunci asing · Dua versi bernomor urut
sama pada satu kontrak -> **ditolak** `INV-04`, pesannya menyebut nomor yang bentrok · Kontrak
bertanggal berakhir lebih awal daripada tanggal mulai -> ditolak `INV-53` · Pengenal yang diminta
dari cap waktu atau dari teks -> **tidak ada jalurnya**; satu-satunya sumber pengenal adalah sequence.

**Uji:** **Negatif:** simpan versi yatim; simpan dua versi bernomor urut sama; simpan kontrak
bertanggal terbalik; minta pengenal berulang sesudah sequence berputar penuh (`INV-03`).
**Positif — dan ia yang menangkap pemisahan identitas yang terlalu ketat:** satu kontrak dengan
**tiga** versi berturut-turut **diterima**, dan ketiganya berbagi lapisan beku yang sama tanpa
menyalinnya. Constraint yang menolak versi kedua lulus setiap uji negatif yang pernah ditulis untuk
irisan ini.

**Menggantikan:** `P-01` bunyi lama tidak bergeser; yang bergeser **cakupan irisannya**, bukan
kemampuannya. Sistem lama menyimpan kontrak dan addendum di **satu baris `M_TREATY_IN`** dengan
seluruh halaman clipboard sebagai satu kolom `JSONDATA`, dan pengenalnya `ID` bertipe teks yang
dibentuk dari pola bernomor revisi (`TDA-11`). Identitas yang terpecah dua **tidak pernah ada** di
sana.

**Blocked by:** None (can start immediately)

**Dasar:**
```
EVIDENCED(M_TREATY_IN@Table/, TREATYINDETAIL@Table/, SaveTreatyInDetail_Act@ekspor-2026-09)
        DECIDED(ADR-0040, ADR-0042, KTV-A, KTV-C)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(T-6)
```

- [ ] `KONTRAK` dan `VERSI_KONTRAK` berdiri dengan **seluruh** kolom yang `2-to-spec/KAMUS-KOLOM.md` sebutkan, bernama dan bertipe persis
- [ ] pengenal keduanya datang dari `SEQUENCE` tanpa `CYCLE`; tidak ada jalur lain yang dapat memberi pengenal
- [ ] `INV-04` menolak nomor urut versi ganda di dalam satu kontrak, pesannya menyebut nomornya
- [ ] perilaku hapus kunci asing `VERSI_KONTRAK` → `KONTRAK` **ditetapkan sadar dan tertulis**, bukan dibiarkan bawaan (`INV-18`)
- [ ] uji positif lulus: satu kontrak dengan tiga versi diterima, lapisan bekunya tidak disalin
- [ ] kolom `KEADAAN_SIKLUS_HIDUP` **ada**, dan ketiadaan constraint daftar nilainya **tertulis sebagai pernyataan keputusan** di dalam berkas DDL-nya — bukan dibiarkan terbaca sebagai kelalaian
- [ ] `KTV-A` dan `T-6` tercatat di `ASUMSI-CLEAR.md` sebagai asumsi yang tiket ini bersandar padanya
