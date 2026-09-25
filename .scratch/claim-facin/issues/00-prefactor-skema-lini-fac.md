# 00: PREFACTOR — dua tabel baru dan kolom induk kedua, expand–contract

**Status:** ready-for-agent
**Blocked by:** None (can start immediately)
**Menutup:** AC 1 · 2 · 3 · 83 · 84 · 85 · 86 · 103 · 104 *(9 AC)* — US 6 · 7

## Hasil & nilai pengguna

Hari ini **data klaim fakultatif masuk tidak punya tempat**. `[terverifikasi]` Skema klaim non-life
sudah berjalan untuk lini PROP — tetapi lini FAC menyimpan datanya **tiga tingkat**: objek
pertanggungan → item objek → baris penyesuaian, dan **dua tingkat teratas belum ada tabelnya**.

Sesudah tiket ini, penilai klaim fakultatif dapat menyimpan objek pertanggungan dan itemnya sebagai
**data yang dapat ditanyakan**, dan baris estimasi, pembagian, serta penyesuaian dapat menggantung
pada tingkat yang benar.

⚠️ **Tiket ini memblokir seluruh tiket lain.** Tidak ada perilaku lain yang dapat selesai
sebelumnya, karena semuanya menulis ke tabel yang belum lengkap.

## ⚠️⚠️ Ini BUKAN membuat skema dari nol — ia menambah ke skema yang SUDAH BERISI DATA

⛔ **Bedakan dari tiket 00 Claim Prop**, yang membuat skema dari nol. ⭐ Di sini **skemanya sudah
berjalan** dan **sudah berisi baris lini PROP**. `[keputusan work owner]` **K1 · K2** menjadikan
tabelnya **dipakai bersama**.

**Tiga perubahan, dan hanya yang pertama aman:**

| # | Perubahan | Sifatnya |
| --- | --- | --- |
| ⭐ **1** | **Dua tabel baru** — objek pertanggungan dan item objek *(khas lini FAC)* | ⭐ tambah murni, **aman** |
| ⚠️ **2** | **Lima tabel yang sudah berisi data** mendapat **kolom induk kedua** | ⛔ **berdampak luas** |
| ⚠️ **3** | **Tabel retro fakultatif** mendapat induk **tingkat penyesuaian** | ⛔ **berdampak luas** |

## ⛔⛔ Bentuk wajib: EXPAND – ISI – CONTRACT

⛔ **Jangan menulisnya sebagai satu langkah.** ⚠️ Memasang aturan penjaga bersamaan dengan kolomnya
**akan menolak setiap baris lini PROP yang sudah ada**, sebab kolom baru mereka kosong.

| Tahap | Yang dikerjakan | Keadaan sesudahnya |
| --- | --- | --- |
| ⭐ **Expand** | Tambah kolom induk kedua, **boleh kosong**, ⛔ **tanpa aturan penjaga** | ⭐ baris lama **tetap sah**; baris baru lini FAC dapat ditulis |
| ⭐ **Isi** | Isi kolom induk yang sesuai untuk **baris lini FAC**; baris lini PROP tetap memakai kolom lamanya | ⭐ tiap baris punya **tepat satu** induk terisi |
| ⭐ **Contract** | ⛔ **Baru** pasang aturan penjaga *"tepat satu induk terisi"* | ⭐ bentuknya terkunci |

⭐ **Ketiganya boleh menjadi tiga tiket terpisah** bila lebih jelas begitu — ⛔ **tetapi wajib
berurut dan saling memblokir**, dan tahap **Contract** wajib diblokir oleh **seluruh** tahap Isi.

⭐ **Preseden pola dua induk sudah ada** dan `[keputusan work owner]` **2026-09-18**: tabel usulan
pandangan memakai dua induk boleh-kosong dengan penjaga *"tepat satu terisi"*.

## Area codebase

- Lapisan penyimpanan klaim non-life — **dipakai bersama lini FAC dan PROP**
- Lapisan layanan klaim: penulisan objek pertanggungan dan item objek
- Migrasi skema — ⭐ **tiga langkah berurut**, bukan satu

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Objek pertanggungan dan itemnya | daftar objek pada data klaim · daftar item di dalamnya |
| Tingkat menggantungnya estimasi dan pembagian | sarang item objek |
| Retro fakultatif tingkat penyesuaian | daftar retro di dalam baris penyesuaian |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0014** — wewenang ditegakkan di lapisan layanan *(menyentuh tiket wewenang, bukan tiket ini)*
- **ADR-0007** — jejak audit

## Acceptance criteria

- [ ] **AC 1** — objek kerja klaim fakultatif **terpisah** dari objek kerja lini PROP
- [ ] **AC 2** — data tersimpan **tiga tingkat**; item objek tanpa objek induk **ditolak**, baris
      penyesuaian tanpa item objek **ditolak**
- [ ] **AC 3** — objek dan item objek **tetap ada** sesudah kasus ditutup lalu dibuka kembali
- [ ] **AC 83 · 84 · 85 · 86** — aturan baca yang mengikat pembangunan dihormati
- [ ] **AC 103 · 104** — keutuhan penunjuk antar tabel terjaga
- [ ] ⭐ **Baris lini PROP yang sudah ada tetap sah** sepanjang tahap Expand dan Isi
- [ ] ⭐ Aturan penjaga **tepat satu induk terisi** aktif pada kelima tabel **hanya** sesudah tahap Isi selesai

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | ⛔ tidak menahan tahap Expand |
| ⚠️ **11** | **29 medan bergantung jenis objek** — menjadi kolom, atau dibaca lewat rujukan ke polis | ⚠️ **menahan bentuk akhir tabel objek** |

## Perintah verifikasi

1. Sisipkan satu baris lini FAC dan satu baris lini PROP pada tiap tabel berkunci ganda —
   ⭐ **keduanya diterima**.
2. Coba sisipkan baris dengan **kedua** kolom induk terisi — ⛔ **ditolak** *(sesudah Contract)*.
3. Coba sisipkan baris dengan **kedua** kolom induk kosong — ⛔ **ditolak** *(sesudah Contract)*.
4. Hapus satu objek pertanggungan — ⭐ item objek dan anaknya **ikut terhapus**.
5. Jalankan tahap Contract **sebelum** tahap Isi selesai — ⛔ **harus gagal**, dan kegagalan itu
   adalah **bukti urutannya benar**.
