# Pertanyaan terbuka — empat kode kategori lampiran tanpa nama

> ⭐ **TERJAWAB 8 Oktober 2026.** Sumbernya tabel yang RD Pega baca sendiri —
> `GetMasterTreatyCategory_SQL`: `SELECT id, note FROM POOLDATA.M_KATEGORIMASTERTREATY order by note`.
> Tabel itu tidak ada di sapuan di bawah. Isinya (terukur):
>
> | Kode | Nama |
> |---|---|
> | `00003` | Binding, signed share Email |
> | `00004` | Info Pack *(NOTE tersimpan dengan ekor CR LF)* |
> | `00008` | Letter of Acknowledgment / LOA |
> | `00009` | Claim Data |
>
> `repository.BacaKatalogKategoriLampiran` kini membaca katalog itu lebih dulu (data lampiran hanya
> melengkapi), sehingga kesebelas kategori `Dipastikan` dan tombol **Upload file** hidup
> (`modul/treatyin/docs/KEPUTUSAN-SASARAN-TULIS.md` §5.5). Isi di bawah dibiarkan sebagai riwayat.

**Diajukan kepada pemilik proses 4 Oktober 2026.**

⛔ **Berkas ini TIDAK memuat tebakan, dan itu seluruh pokoknya.** Empat kode dan empat nama
sama-sama diketahui; **pasangannya tidak**. Menebaknya menaruh berkas di kategori yang salah,
dan itu baru ketahuan bertahun kemudian — saat seseorang mencari dokumen yang "seharusnya ada
di sana".

---

## Yang TERBUKTI — tujuh pasangan, dibaca dari data

`POOLDATA.M_ATTACHMENTTREATY_2` menyimpan `CATEGORY_ID` **dan** `CATEGORY` berdampingan, jadi
ketujuh pasangan yang terpakai adalah fakta yang diambil, bukan daftar yang dihafal:

| Kode | Nama | Baris |
| --- | --- | ---: |
| `00000` | Others | 5 |
| `00001` | Analysed Email | 10 |
| `00002` | Approval Email | 14 |
| `00005` | Summary Treaty Leader | 10 |
| `00006` | Assessment Inward Treaty Form / Format Analisa Treaty | 2 |
| `00007` | Pega Proportional Calculation /Perhitungan Pega Proportional | 1 |
| `00010` | Offer Email | 1 |
| | **total** | **43** |

`TestKatalogKategoriTerbacaDariData` mengunci ketujuhnya, dan `TestSatuKodeSatuNama` menjaga satu
kode tidak diam-diam punya dua nama.

---

## Yang TIDAK diketahui — empat kode, empat nama, pasangan nihil

**Empat kode tanpa baris:** `00003` · `00004` · `00008` · `00009`
**Empat nama di layar lama:** *Binding, signed share Email* · *Claim Data* · *Info Pack* ·
*Letter of Acknowledgment / LOA*

⛔ **Keduanya ditulis menurut abjad di atas, dan urutannya BUKAN pasangan.**

### Sudah disapu, dan nihil

| Tempat | Hasil |
| --- | --- |
| `D:\XML_NURE\Treaty In` + `Treaty In Adjustment`, keempat kode | **nol** — satu-satunya `00009` ternyata potongan cap waktu `20210609T100009_887_GMT` di `TreatyInTabsNPValueDifference_NoProRate.xml` |
| idem, keempat nama | **nol berkas** untuk keempatnya; `LOA` sebagai KATA juga nol |
| `M_ATTACHMENTTREATY_2` | keempat kode **nol baris**, jadi pasangannya tidak dapat dibaca dari data |
| `CATEGORY_ATTACH_REAS` | milik FACIN/FACOUT |
| `M_ATTACHMENT_CATEGORY` | milik klaim kendaraan |
| `setCategoryAttachment_DT.xml` | ekspornya kosong |
| `SetkategoriDoc.xml` | `CARI40` disetel saat jalan |
| `Section/ShowAttachmentTreaty.xml` | nol kode `000xx` |

---

## Yang berlaku SEMENTARA

Keempat kode **tetap tampil** di panel dengan `Dipastikan = false` dan **tanpa nama**.

- **Tetap tampil**, sebab menyembunyikannya membuat kategori yang ada di sistem lama lenyap dari
  layar, dan berkas yang tersimpan di sana menjadi tidak terjangkau.
- **Tanpa nama**, sebab nama yang salah lebih buruk daripada nama yang kosong: yang kosong
  memicu pertanyaan, yang salah memicu keyakinan.

`TestKeempatKodeTanpaNamaTIDAKDitebakNamanya` menolak nama apa pun yang menempel pada keempat
kode itu — termasuk keempat nama yang belum berumah.

⭐ **Pertanyaannya dapat terjawab sendiri.** Begitu salah satu kode muncul di
`M_ATTACHMENTTREATY_2` lengkap dengan `CATEGORY`-nya, katalog membacanya dan penandanya lepas
tanpa perubahan kode. `TestKatalogKategoriTerbacaDariData` **merah** pada hari itu, menyebut
kode dan namanya, supaya berkas ini ikut diperbarui.

---

## Yang ditanyakan

1. **Kode mana untuk nama mana** di antara keempatnya?
2. Bila pemilik proses pun tidak tahu: apakah keempat kode itu **boleh dianggap tidak terpakai**
   dan disembunyikan — atau tetap tampil tanpa nama seperti hari ini?
3. Apakah ada **daftar kategori di luar korpus** (layar admin, tabel konfigurasi Pega yang tidak
   ikut diekspor) yang memuat pasangannya?

---

## ⚠️ Temuan sampingan yang perlu keputusannya sendiri

**37 dari 43 nama berkas yang SUDAH tersimpan melanggar aturan spanduk biru**
(*"Recommended safe substitute should be . or _"*) — spasi dan tanda kurung, misalnya
`Signed Schedule Marine Hull Quota Share 2025 (3).pdf`.

Aturan itu ditegakkan untuk unggahan **baru** (`services.NamaBerkasAman`). Ia **tidak**
diterapkan surut, dan tidak ada berkas lama yang ditolak atau diganti namanya — menolak 86% dari
yang sudah ada berarti menutup pintu yang sistem lama buka.

**Ditanyakan:** apakah nama lama perlu dirapikan saat pemindahan, atau dibiarkan apa adanya?
