# 02: Rantai generasi dan larangan percabangan

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Satu polis boleh di-endorse **tanpa batas**. Batas teknis 99 tidak dipertahankan; kolom nomor generasi dilebarkan. Diuji ke korpus: kenaikannya aritmetika *(`Prodke+1`)*, nol topeng lebar tetap.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 1.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** ⛔ `[work owner]` **batas berapa kali satu polis boleh di-endorse** — batas teknis 99, batas dagang belum ditetapkan~~ ⛔ **penahan gugur 23-09-2026**
**Bergantung pada tiket NB:** **16** *(kerangka penyimpanan dan kunci generasi)*
**Menutup:** AC **1–6** *(6 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-9 · ID-10 · ID-11 · ID-12 · ID-13

## Hasil & nilai pengguna

Endorsemen tersimpan sebagai **generasi berikutnya** dari polis yang sama, dan rantainya **lurus**:
satu generasi, satu penerus.

⛔⛔ **Ini memperbaiki bahaya nyata, bukan kehati-hatian teoretis.** `[terverifikasi]` Sistem lama
membaca nomor generasi terakhir **tanpa penguncian**, lalu menambah satu. Dua endorsemen serentak
atas polis yang sama membaca angka yang sama, dan **keduanya tersimpan**.
⭐ Di sistem baru yang kedua **gagal terang-terangan**.

## Yang dibangun

Baris generasi bernomor ≥ 1 dengan penunjuk ke generasi sebelumnya, beserta tiga aturan:

| Aturan | Yang dicegahnya |
| --- | --- |
| penunjuk generasi **unik** | rantai menjadi pohon — *"selisih mana yang benar"* tidak terjawab |
| kunci alami **unik** | dua endorsemen serentak mendapat nomor generasi yang sama |
| generasi lampau **tidak dapat disunting** | dasar selisih berubah sesudah disetujui |

Ditambah pembentukan nomor endorsemen: nomor polis induk, pemisah, lalu **dua digit** dengan nol di
depan bila kurang dari sepuluh.

## Batas — yang TIDAK termasuk

⛔ Aturan keutuhan nomor urut antar generasi — tiket **03**.
⛔ Perhitungan selisih — tiket **06**.
⛔ Jalur penomoran kedua yang ada di sistem lama — **tidak dimigrasi**.

## Cara mengujinya

Lewat seam `repository`. ⭐ **Tiga uji yang wajib berdiri sendiri:**

1. **rantai tiga generasi** — ⚠️ rantai dua generasi **tidak memadai**, sebab selisih terhadap
   generasi tepat sebelumnya baru terbukti mulai generasi ketiga;
2. **penolakan penerus kedua** atas generasi yang sama;
3. **dua endorsemen serentak** — yang kedua harus ditolak **oleh basis data**, bukan oleh
   pemeriksaan di aplikasi.

## Acceptance criteria

- [ ] **AC 1** — baris endorsemen tersimpan dengan nomor generasi ≥ 1
- [ ] **AC 2** — penunjuk generasi **terisi** dan menunjuk generasi sebelumnya
- [ ] **AC 3** — dua baris tidak boleh berbagi penunjuk generasi yang sama
- [ ] **AC 4** — dua endorsemen serentak: yang kedua **ditolak**
- [ ] **AC 5** — nomor endorsemen berbentuk nomor polis + pemisah + **dua digit**
- [ ] **AC 6** — menyunting generasi yang sudah punya penerus **ditolak**

## ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Belum ditetapkan berapa kali satu polis boleh di-endorse.** Batas teknisnya **99**
— nomor generasi dibentuk dua digit. Bila ada batas dagang yang lebih rendah, itu **validasi yang
perlu dibangun** di tiket ini; bila tidak ada, angka 99 tetap batas yang wajib diketahui sebelum
sistem baru memakainya.

⭐ **Lima dari enam AC tidak bergantung pada jawabannya** — hanya batas cacah generasi yang
menunggu.