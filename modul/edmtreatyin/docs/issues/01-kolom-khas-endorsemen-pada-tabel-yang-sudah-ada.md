# 01: Kolom khas endorsemen pada tabel yang sudah ada

**Status:** ready-for-agent
**Blocked by:** — *(dapat mulai sesudah tiket NB)*
**Bergantung pada tiket NB:** **16** *(kerangka penyimpanan)* · **19** *(pemecah dokumen)*
**Menutup:** AC **32–34** · AC **54–55** *(5 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-5 · ID-6 · ID-20 · ID-21 · ID-22

## Hasil & nilai pengguna

Endorsemen menempati **tabel yang sama** dengan polis baru. Yang membedakannya hanya **kolom mana
yang terisi** — bukan bentuk tabelnya.

⭐⭐ **Ini tiket pertama rantai endorsemen, dan sekaligus penjaganya:** ia membuktikan bahwa bentuk
tabel dasar **tidak berubah** oleh seluruh pekerjaan endorsemen.

## Yang dibangun

Pengisian dan pengosongan kolom menurut jenis berkas:

| Golongan | Perlakuan |
| --- | --- |
| **21 kolom khas endorsemen** | terisi pada endorsemen, ⛔ **kosong** pada polis baru |
| **4 kolom khas polis baru** | kosong pada endorsemen, ⛔ **tidak dihapus** dari skema — bukan kolom mati |
| **keadaan layar** | ⛔ **tidak disimpan** di tabel mana pun |

Ditambah dua pernyataan bentuk yang diuji langsung: bentuk kesembilan tabel dasar **tidak berubah**,
dan cacah tabel pada kedua bentuk endorsemen sesuai ketetapan.

⛔ **Nol tabel dasar dibuat tiket ini.** Sepuluh tabel dasar dibuat tiket NB **16** dan **19**.

## Batas — yang TIDAK termasuk

⛔ Pembuatan tabel dasar — tiket NB **16** · **19**.
⛔ Tabel proyeksi selisih — tiket **07**.
⛔ Tipe dan presisi kolom — tiket **11**.

## Cara mengujinya

Lewat seam `repository`. ⭐ **Uji berdampingan:** kasus yang sama disimpan sebagai polis baru dan
sebagai endorsemen — kolom yang terisi harus **berbeda**, dan kolom yang kosong **tetap ada**.

⚠️ Uji bentuk tabel dijalankan sebagai pemeriksaan skema, bukan lewat data.

## Acceptance criteria

- [ ] **AC 32** — dua puluh satu kolom khas endorsemen **kosong** pada baris polis baru
- [ ] **AC 33** — empat kolom khas polis baru **kosong** pada endorsemen, dan **tidak** dihapus dari skema
- [ ] **AC 34** — keadaan layar **tidak tersimpan** di tabel mana pun
- [ ] **AC 54** — bentuk kesembilan tabel dasar **tidak berubah**
- [ ] **AC 55** — cacah tabel pada kedua bentuk endorsemen sesuai ketetapan
