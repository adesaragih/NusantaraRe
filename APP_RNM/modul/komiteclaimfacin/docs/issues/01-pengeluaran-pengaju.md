# 01: Pengeluaran pengaju — jenjang pertama saja

**Status:** ready-for-agent
**Blocked by:** 00
**Menutup:** AC 9 · 10 · 11 · 12 · 13 *(5 AC)* — US 13

## Hasil & nilai pengguna

Hari ini Seorang penilai dapat **menyetujui klaim yang ia ajukan sendiri**, sebab tidak ada yang memeriksanya.

Sesudah tiket ini, Bila **pemegang jenjang pertama adalah pengaju**, ia **dikeluarkan dari susunan jenjang** kasus itu, dan **jumlah jenjang berkurang satu**.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Pengeluaran pengaju | rule penyiapan wewenang penyetujuan mengeluarkan anggota **indeks 1** bila ia akun yang sedang membuka kasus, lalu **menghitung ulang** jumlah jenjang |
| ⛔ Lingkupnya | ⛔ **indeks ditanam `1`** — hanya anggota pertama yang diperiksa |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K9** — ⚠️ **DITIRU APA ADANYA** — ⛔ **BERBEDA dari rekomendasi asisten**, yang mengusulkan berlaku untuk semua jenjang. ⛔ **Disengaja**

## Yang harus diuji

- [ ] Pengaju di jenjang **pertama** ⇒ **dikeluarkan**, jumlah jenjang **berkurang satu**
- [ ] Jenjang berikutnya **naik menjadi jenjang pertama**
- [ ] ⛔⛔ **Pengaju di jenjang KEDUA ke bawah TETAP di daftar, dan TETAP dapat menyetujui klaim yang ia ajukan sendiri** — ⭐ **lubang yang dibawa masuk dengan sadar, bukan kelalaian**
- [ ] Susunan jenjang **kosong** sesudah pengeluaran ⇒ kasus **tidak dibentuk**, pengajuan **dikembalikan** berikut alasannya

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi daftar jabatan | ⚠️ menahan penentuan siapa pemegang jenjang pertama |

## Seam & verifikasi

**Seam:** lapisan layanan komite, pada pembentukan kasus.
⭐⭐ **UJI DUA ARAH — dan arah kedua ADALAH inti tiket ini:**
   *(a)* pengaju di jenjang **pertama** ⇒ ⭐ **dikeluarkan**, jumlah jenjang berkurang.
   *(b)* pengaju di jenjang **kedua** ⇒ ⛔ **TETAP DI DAFTAR**, dan **dapat menyetujui**.
⚠️ **Arah (b) bukan uji fungsi — ia MENGUNCI keputusan K9** supaya pengembang berikutnya tidak membacanya sebagai bug lalu 'memperbaikinya' diam-diam.
6. Pengaju satu-satunya anggota ⇒ ⭐ kasus **tidak dibentuk**, pengajuan dikembalikan.
