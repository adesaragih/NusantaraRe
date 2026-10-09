# 00: PREFACTOR — dua tabel baru dan kolom induk kedua, expand–contract

**Status:** ready-for-agent
**Blocked by:** None (can start immediately)
**Menutup:** AC 1 · 2 · 3 · 83 · 84 · 85 · 86 · 103 · 104 *(9 AC)* — US 6 · 7

> ⛔ **RALAT 10-10-2026.** Judul lamanya dikutip utuh, tidak dihapus: *"00: PREFACTOR — dua tabel baru dan kolom induk
> kedua, expand–contract"* → **bentuk tiket ini ditolak** keputusan work owner **OQ-CFI-01 (09-10-2026)**: nol
> expand–contract, nol `MODIFY`, nol `CHECK` — **hanya `ADD`**. Tahap 1 dibangun 10-10-2026 (`MODUL.md` bab Migrasi);
> rinciannya di RALAT bab *Bentuk wajib* dan *Acceptance criteria* di bawah.

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

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⛔ **Jangan menulisnya sebagai satu langkah.**
> ⚠️ Memasang aturan penjaga bersamaan dengan kolomnya **akan menolak setiap baris lini PROP yang sudah ada**, sebab
> kolom baru mereka kosong."* → logika expand–isi–contract **tidak berlaku**, sebab **tidak ada aturan penjaga** yang
> dipasang: OQ-CFI-01 menolak induk ganda. Yang dibangun — **semuanya `ADD`, nol `MODIFY` / `DROP`**:
>
> | Migrasi | Isi |
> | --- | --- |
> | `560` | 19 kolom kepala FAC `ADD` ke `T_GENERAL_CLAIM` (nullable) |
> | `561` · `562` | tabel baru `T_CLAIM_OBJECT` (`CLAIM_ID` CASCADE) · `T_CLAIM_OBJECT_ITEM` (`OBJECT_ID` CASCADE, `CLAIM_ID` ikut diisi) |
> | `563`–`566` | `ADD OBJECT_ITEM_ID` (nullable, FK CASCADE ke item) + kolom khas FAC di `T_CLAIM_ESTIMATION`, `T_CLAIM_SPREADING` (+ `JENIS`), `T_CLAIM_BREAK_QS`, `T_CLAIM_ADJUSTMENT` |
> | `567` | `ADD ADJUSTMENT_ID` (nullable, FK CASCADE ke adjustment) di `T_CLAIM_FAC_RETRO` |
>
> Baris FAC mengisi **`CLAIM_ID` dan `OBJECT_ITEM_ID`** (atau `ADJUSTMENT_ID` untuk retro); `CLAIM_ID` tetap NOT NULL.
> Baris PROP lama tetap sah karena kolom baru boleh kosong dan tidak ada `CHECK`. Tabel *"Tiga perubahan"* di atas:
> perubahan 2 dan 3 bukan *"kolom induk kedua"* yang menggantikan `CLAIM_ID`, melainkan kolom **tambahan**.

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

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"**AC 103 · 104** — keutuhan penunjuk antar tabel
> terjaga"* dan *"⭐ Aturan penjaga **tepat satu induk terisi** aktif pada kelima tabel **hanya** sesudah tahap Isi
> selesai"* →
>
> - **AC 103** menyatakan sistem lama **nol pemeriksaan** penunjuk dan test yang menuntut pemeriksaan itu **ditiru**
>   gagal; **AC 104** ber-`[terbuka]` (butir 21 spec — perilaku bila penunjuk salah alamat belum diputuskan). Tiket
>   tidak dapat "menutup" AC `[terbuka]`. Yang dibangun: induk tiap tingkat lewat **kunci asing ber-ID stabil**
>   (`SEQ_T_CLAIM`, CASCADE, `561`–`567`), dan `T_GENERAL_KOMITE.ADJUSTMENT_ID` menunjuk ID stabil adjustment
>   (katalog `Stabil: true`), bukan indeks posisi.
> - Butir penjaga *"tepat satu induk terisi"* **gugur** — nol `CHECK` (OQ-CFI-01).

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **1** | Kolom tabel data kutipan — `[data DBA]` | ⛔ tidak menahan tahap Expand |
| ⚠️ **11** | **29 medan bergantung jenis objek** — menjadi kolom, atau dibaca lewat rujukan ke polis | ⚠️ **menahan bentuk akhir tabel objek** |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"Kolom tabel data kutipan — `[data DBA]`"* dan
> *"**29 medan bergantung jenis objek** — menjadi kolom, atau dibaca lewat rujukan ke polis"* →
>
> - ⚠️ **Dua ruang nomor tercampur di tabel ini.** Butir **1** = register **spec** (Bab 19 / *Butir `[terbuka]`*);
>   butir **11** = register **`STRUKTUR-TABEL-CLAIM-FACIN.md` §6** — butir 11 spec adalah hal lain (mesin tiket).
> - **Butir 1 spec**: tidak lagi memblokir — data polis dibaca dari `JSON_POLIS.DATA_JSON`, bukan dari
>   `T_QUOTATIONDATA` (RALAT spec AC 15).
> - **Butir 11 STRUKTUR**: dibaca lewat rujukan ke polis — `T_CLAIM_OBJECT.KUNCI_POLIS`, halaman polis dibaca ulang
>   dari `JSON_POLIS` (migrasi `561`). Tidak menahan apa pun lagi.

## Perintah verifikasi

1. Sisipkan satu baris lini FAC dan satu baris lini PROP pada tiap tabel berkunci ganda —
   ⭐ **keduanya diterima**.
2. Coba sisipkan baris dengan **kedua** kolom induk terisi — ⛔ **ditolak** *(sesudah Contract)*.
3. Coba sisipkan baris dengan **kedua** kolom induk kosong — ⛔ **ditolak** *(sesudah Contract)*.
4. Hapus satu objek pertanggungan — ⭐ item objek dan anaknya **ikut terhapus**.
5. Jalankan tahap Contract **sebelum** tahap Isi selesai — ⛔ **harus gagal**, dan kegagalan itu
   adalah **bukti urutannya benar**.

> ⛔ **RALAT 10-10-2026.** Langkah lamanya dikutip utuh, tidak dihapus: *"Coba sisipkan baris dengan **kedua** kolom
> induk terisi — ⛔ **ditolak** *(sesudah Contract)*."* → langkah **2, 3, dan 5 tidak berlaku** (nol `CHECK`, nol tahap
> Contract — OQ-CFI-01); baris FAC justru **wajib** mengisi kedua kolom. Langkah 1 dan 4 tetap: baris FAC dan PROP
> sama-sama diterima, dan menghapus objek menghapus item serta anaknya (CASCADE `562`–`567`, cucu adjustment lewat
> FK Claim Prop `529` / `530`).
