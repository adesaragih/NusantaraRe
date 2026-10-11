# 01: Pengeluaran pengaju — jenjang pertama saja

**Status:** ready-for-agent
**Blocked by:** 00
**Menutup:** AC 9 · 10 · 11 · 12 · 13 *(5 AC)* — US 13

> ⛔ **RALAT 10-10-2026 — TIKET GUGUR, jangan dibangun.** Kepala lamanya dikutip utuh, tidak dihapus: *"**Status:**
> ready-for-agent"* dan *"**Menutup:** AC 9 · 10 · 11 · 12 · 13 *(5 AC)* — US 13"* → **K9 salah baca** (prompt tahap 2
> §3 "Gugur"). `ApprovalKomite_Act` L6.2 (`Property-Remove`, when `pyWorkPage.KomiteList(1).KomiteID == .OPERATOR_ID`)
> hanya membuang **calon perluasan** (`GetKomite.pxResults`, roster DEGREE > 1 dari L4 / L5) yang sama dengan **anggota
> tingkat 1**, supaya satu baris roster tidak masuk tangga dua kali. Ia **tidak** membandingkan pengaju / pembuka kasus
> dan **tidak** mengeluarkan anggota tingkat 1. Korpus Komite Claim FacIn tidak memuat larangan menyetujui klaim sendiri
> di jenjang mana pun, jadi tidak ada yang dibangun; tanpa larangan rangkap (pola Komite Prop 09-10-2026). L6.2 tetap
> dibangun sebagai bagian perluasan KCF-02 (tiket 00). AC 9–13 dan US 13 gugur (RALAT spec AC 9–13 dan ID-5).

## Hasil & nilai pengguna

Hari ini Seorang penilai dapat **menyetujui klaim yang ia ajukan sendiri**, sebab tidak ada yang memeriksanya.

Sesudah tiket ini, Bila **pemegang jenjang pertama adalah pengaju**, ia **dikeluarkan dari susunan jenjang** kasus itu, dan **jumlah jenjang berkurang satu**.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Pengeluaran pengaju | rule penyiapan wewenang penyetujuan mengeluarkan anggota **indeks 1** bila ia akun yang sedang membuka kasus, lalu **menghitung ulang** jumlah jenjang |
| ⛔ Lingkupnya | ⛔ **indeks ditanam `1`** — hanya anggota pertama yang diperiksa |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"Pengeluaran pengaju | rule penyiapan wewenang
> penyetujuan mengeluarkan anggota **indeks 1** bila ia akun yang sedang membuka kasus, lalu **menghitung ulang** jumlah
> jenjang"* → yang dibuang L6.2 adalah **calon** di daftar hasil browse roster, bukan anggota indeks 1; indeks `1` adalah
> pembanding (anggota tingkat 1), bukan sasaran. Hitung ulang `KomiteLoop` (L8) terjadi sesudah calon ditambahkan (L7),
> bukan sesudah ada yang dikeluarkan. Sejak KCF-01 anggota tingkat 1 = `ReasClaimSPVA` dan calon DEGREE > 1 berisi
> workbasket lain, sehingga L6.2 praktis tidak pernah membuang baris; langkahnya tetap dibangun sebagai penjaga baris
> ganda.

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
