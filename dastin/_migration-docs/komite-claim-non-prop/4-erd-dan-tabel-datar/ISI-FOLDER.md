# Struktur — ERD dan tabel datar, modul Komite

Bentuk data, digambar. Semuanya **TURUNAN**: dibangkitkan alat di `alat/` dan tidak pernah disunting tangan. Bila salah satu berbeda dari sumbernya, **alatnya yang salah**.

**8 berkas** (termasuk berkas ini). Sejajar dengan [`claim-non-prop/4-erd-dan-tabel-datar/`](../../claim-non-prop/4-erd-dan-tabel-datar/ISI-FOLDER.md).

| Berkas | Isinya | Alat |
|---|---|---|
| `struktur-komite-lama.md` | struktur halaman kerja Komite sistem lama — **sumber** bagi ERD di folder ini | `buat-pohon-komite.py` |
| `datar-komite-lama.csv` | pohon properti — 688 simpul, satu baris per simpul | idem |
| `datar-komite-kelas.csv` | peta kelas — 112 pasangan halaman–kelas, 84 halaman | idem |
| `datar-komite-tulis-balik.csv` | kontrak tulis-balik ke klaim induk — 42 baris | idem |
| `ERD-KOMITE-LAMA.md` | ERD dan tabel datar — 14 tabel, 14 relasi, 60 kolom | `buat-skema-komite.py` |
| `datar-komite-tabel.csv` | **tabel datar** — 85 baris, satu baris per kolom | idem |
| `RELASI-KOMITE.csv` | 14 relasi | idem |
| `ERD-KOMITE-LAMA.html` + `Diagram-Skema-Tabel-Komite.xlsx` | ERD tergambar + workbook empat sheet | idem |

> **Ini SISTEM LAMA.** Rancangan skema baru hidup di [`SPEC-KOMITE-01.md`](../SPEC-KOMITE-01.md) bagian *Model data* — tujuh objek `KLAIMNP`, 56 kolom. Folder ini memerikan bentuk yang **ada**, termasuk cacatnya; folder itu menetapkan bentuk yang **akan berdiri**.

> **Penamaan tabel** mengikuti keputusan work owner 20 September 2026 yang tercatat di `Diagram-Skema-Tabel-NusantaraRe.xlsx`: akar `T_WORK_CLAIM`, badan 1:1 ber-SHARED PK `T_GENERAL_*`, lalu awalan lini — di sini `T_KOMITE_`. Dua nama dipakai ulang dari rancangan sisi Klaim (`T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`); keduanya ditebak di sana tanpa membuka folder ini, dan `ERD-KOMITE-LAMA.md` bagian 3 mencatat mana yang terkukuhkan dan mana yang terkoreksi.

---

Register keadaan modul: [`../INDEX.md`](../INDEX.md).
