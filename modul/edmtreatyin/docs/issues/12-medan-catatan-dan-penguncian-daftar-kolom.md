# 12: Medan catatan dan penguncian daftar kolom

> ## ⭐ PENAHAN GUGUR — 23 September 2026 sore
>
> `[keputusan work owner]` *"Abaikan `JSON_DATAGUIDE`, ikuti dari data yang digunakan di Activity dan Section."* ⭐ Daftar medan disusun ulang dari **302 berkas aturan** — **394 medan unik**, 232 di antaranya tampil di layar. Hasilnya di **`DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`**. Data guide turun derajat jadi **penambal**, sah hanya untuk panjang maksimum per medan.
>
> ⭐ **`blocked` → `ready-for-agent`.** ⛔ **Nol butir `[data DBA]` tersisa di tiket ini.**
>
> Rinciannya: `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026 sore)*
~~**Blocked by:** ⛔ `[data DBA]` **panduan bentuk dokumen terbukti basi** — satu dokumen memuat~~ ⛔ **gugur 23-09-2026 sore**
95 jalur yang tidak ada di dalamnya
**Bergantung pada tiket NB:** **19** *(pemecah dokumen menjadi baris)*
**Rujukan:** `ID-27b` · `ID-27c`
**Menutup:** dua implementation decision yang tidak tercakup ronde tiket sebelumnya

> ⚠️ **Tiket ini lahir dari audit sesudah ronde, bukan dari ronde tiket.** Kedua butirnya
> disisipkan ke spec sesudah spec pertama selesai, sehingga tidak terbaca saat tiket 01–11 ditulis.
> Kekeliruan lingkup itu milik penyusun brief.

## Hasil dan nilai pengguna

Dua hal kecil yang, bila terlewat, menghasilkan kerugian yang tidak kelihatan sampai terlambat.

Yang pertama satu medan catatan bebas yang selama ini ikut tersimpan pada dokumen polis dan tidak
pernah masuk rancangan tabel. Bila ia hilang saat pemindahan, isinya tidak dapat dipulihkan —
tidak ada tempat lain yang menyimpannya.

Yang kedua bukan medan melainkan **disiplin**. Daftar kolom disusun dari medan yang benar-benar dipakai aturan — Activity dan Section — bukan dari panduan bentuk dokumen, yang terbukti tidak lengkap. Menyusunnya dari panduan berarti membangun pemecah dokumen yang diam-diam membuang medan yang tidak dikenalinya.

## Lingkup

### Bagian 1 · medan catatan *(`ID-27b`)*

- Medan catatan bebas tingkat atas, panjang **128**, **ikut dipindahkan**.
- ⛔ Dua medan tampilan yang bertetangga dengannya **tidak** dipindahkan — keduanya keadaan layar,
  bukan data dagang. Keputusan lama, tetap berlaku.
- Uji: dokumen yang membawa medan catatan terisi, dipindahkan lalu dibaca kembali, isinya sama
  persis termasuk spasi di ujung.

### Bagian 2 · penguncian daftar kolom *(`ID-27c`)*

- ⭐ **Dasar daftar kolom adalah `DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`** — 394 medan yang disapu dari 302 berkas aturan, bukan panduan bentuk dokumen.
- Panduan bentuk dokumen berkedudukan **penambal**: sah untuk panjang maksimum per medan, dan untuk 24 nama yang ditulis sistem sehingga tidak tersapu aturan.
- ⛔ Ia **tidak sah** dipakai untuk membuktikan bahwa sebuah medan tidak ada.
- Pemecah dokumen wajib punya **penampung medan tak dikenal**, dan penampung itu wajib **kosong**
  sebelum pekerjaan dinyatakan selesai.
- Uji: dokumen yang membawa medan di luar daftar tidak boleh diam-diam kehilangan medan itu —
  ia masuk penampung, dan keberadaan isi di penampung menggagalkan test.

## Batas — yang TIDAK termasuk

- ⛔ Bukan tempat menetapkan daftar kolom akhir. Itu tiket NB 19.
- ⛔ Bukan tempat memutuskan presisi fisik kolom uang. Itu tiket EDM 11.
- ⛔ Tidak membuat tabel apa pun.

## Uji

Lewat seam `repository`, seperti tiket lain di berkas ini.

⚠️ Uji pulang-pergi saja **tidak cukup** untuk bagian 2 — test wajib memeriksa **isi penampung
medan tak dikenal** secara langsung, karena dokumen yang kehilangan medan tetap lolos pulang-pergi.

## Butir tertahan

| Butir | Pemilik |
| --- | --- |
| ~~⛔ panduan bentuk dokumen perlu disegarkan~~ ✅ **gugur 23-09 sore** — panduan diabaikan, daftar medan disusun dari korpus | ~~`[data DBA]`~~ |

⛔ Butir ini **tidak ditutup di tiket ini**. Surat permintaannya sudah disusun.

## Area codebase

- Lapisan `repository`: pemecah dokumen dan penampung medan tak dikenal
- Pemuat pemindahan dokumen lama

## Disiplin berkas

⛔ Nol `CREATE TABLE`. ⛔ Nol nama orang. ⛔ Nol nomor polis harfiah. ⛔ Nol cuplikan data produksi.

---

*Ditulis 23 September 2026, sesudah audit pasca-ronde menemukan dua implementation decision yang
tidak tercakup tiket mana pun.*
