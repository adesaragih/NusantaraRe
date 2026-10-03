# 24: Rantai generasi endorsemen dan pembatalan

> ## ⛔⛔ DIGANTIKAN — 23 September 2026
>
> **Tiket ini tidak dikerjakan.** Ia digantikan oleh
> **`.scratch\edm-treaty-in\issues\01 · 02 · 03 · 05`**.
>
> ### Kenapa
>
> Brief ronde tiket sebelumnya berlingkup **dua spec sekaligus**, sehingga pekerjaan **endorsemen**
> jatuh ke folder tiket **polis baru**. Ronde berikutnya membaginya ulang ke folder yang benar,
> lebih halus, dan — yang menentukan — **setiap tiket menyebut tiket NB mana yang membuat
> tabelnya**. Tanpa gate itu, tiket endorsemen bisa dikerjakan sebelum tabelnya ada.
>
> ⚠️ **Kekeliruan lingkup ini milik penyusun brief, bukan pelaksana ronde ini.**
>
> ⭐ Berkas ini **tidak dihapus** supaya jejaknya tidak hilang — tanpa ini, nomor tiket NB akan
> tampak melompat dari 23 ke 29 tanpa sebab.
>
> ### Yang berlaku
>
> | | |
> | --- | --- |
> | Tiket penyimpanan **NB** | `nb-treaty-in\issues\` **16–23** |
> | Tiket penyimpanan **EDM** | `edm-treaty-in\issues\` **01–11** |
>
> ⛔ Isi di bawah dibiarkan apa adanya sebagai catatan, **bukan sebagai pekerjaan.**

---


**Status:** ~~ready-for-agent~~ — **DIGANTIKAN**
**Blocked by:** **16** · **17**
**Menutup:** EDM AC **1–15** *(15 AC)*
**Sumber:** `edm-treaty-in\spec-penyimpanan-relasional.md` ID-8..ID-19 · ID-41

## Hasil & nilai pengguna

Endorsemen tersimpan sebagai **generasi berikutnya** dari polis yang sama — bukan sebagai jenis
berkas tersendiri. Rantainya lurus: satu generasi, satu penerus.

⛔⛔ **Ini memperbaiki bahaya nyata.** `[terverifikasi]` Sistem lama membaca nomor generasi terakhir
**tanpa penguncian**, lalu menambah satu. Dua endorsemen serentak atas polis yang sama membaca angka
yang sama, dan **keduanya tersimpan**. ⭐ Di sistem baru yang kedua **gagal**.

## Yang dibangun

Baris generasi dengan nomor ≥ 1 dan penunjuk ke generasi sebelumnya, beserta empat aturan:

| Aturan | Yang dicegahnya |
| --- | --- |
| penunjuk generasi **unik** | rantai menjadi pohon; *"selisih mana yang benar"* tidak terjawab |
| kunci alami **unik** | dua endorsemen serentak mendapat nomor sama |
| generasi lampau **tidak dapat disunting** | dasar selisih berubah sesudah disetujui |
| ⭐ **aturan keutuhan** | generasi baru **wajib memuat setiap nomor urut** milik generasi sebelumnya |

⭐ **Pada endorsemen baris rincian TIDAK dapat dihapus**, dan nomor urut lama terbawa apa adanya —
⚠️ **berbeda dari polis baru**, tempat penomoran ulang dibolehkan.

⭐ **Pembatalan adalah jenis endorsemen**, dipilih di awal — ia menghasilkan generasi bernilai nol,
⛔ **bukan penghapusan baris**.

⭐ **Nilai lama tidak menjadi tabel.** Ia baris yang ditunjuk penunjuk generasi.

## Batas — yang TIDAK termasuk

⛔ Perhitungan selisih — tiket **25**.
⛔ Tabel proyeksi — tiket **26**.
⛔ Kolom khas endorsemen — tiket **27**.

## Cara mengujinya

Lewat seam `repository`. ⭐ **Tiga uji yang wajib berdiri sendiri:** rantai **tiga generasi** ·
penolakan penerus kedua · penolakan generasi yang kehilangan satu nomor urut.
⚠️ Rantai dua generasi **tidak memadai** — selisih terhadap generasi tepat sebelumnya baru terbukti
mulai generasi ketiga.

## Acceptance criteria

- [ ] **AC 1** — baris endorsemen tersimpan dengan nomor generasi ≥ 1
- [ ] **AC 2** — penunjuk generasi **terisi** dan menunjuk generasi sebelumnya
- [ ] **AC 3** — dua baris tidak boleh berbagi penunjuk generasi yang sama
- [ ] **AC 4** — dua endorsemen serentak: yang kedua **ditolak**
- [ ] **AC 5** — nomor endorsemen berbentuk nomor polis + pemisah + dua digit
- [ ] **AC 6** — menyunting generasi yang sudah punya penerus **ditolak**
- [ ] **AC 7** — selisih dihitung terhadap generasi **tepat sebelumnya**
- [ ] **AC 8** — generasi yang kehilangan satu nomor urut **ditolak**
- [ ] **AC 9** — nol tabel salinan nilai lama; ia dibaca lewat penunjuk generasi
- [ ] **AC 10** — pada endorsemen, baris rincian **tidak dapat dihapus**
- [ ] **AC 11** — nomor urut lama **terbawa apa adanya**
- [ ] **AC 12** — baris baru mendapat nomor urut **maksimum + 1**
- [ ] **AC 13** — pasangan yang bergeser diperlakukan sebagai **anomali**
- [ ] **AC 14** — pembatalan menghasilkan generasi bernilai **nol**, bukan penghapusan
- [ ] **AC 15** — jenis endorsemen tersimpan di kolomnya sendiri, dipilih di awal
