# Modul `treatydescription` — Treaty Description

Satu folder, satu modul, satu pemilik: kode backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur
tim satu folder per modul, keputusan work owner 30-09-2026).

**Modul di luar dua puluh folder korpus** (`APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 5) — perintah work owner
05-10-2026: *"INI MODUL BARU, SELECT * FROM TREATYDESC, PAHAMI TABEL ITU, AKU MAU BUAT CRUD"*, *"NAMANYA TREATY
DESCRIPTION YA"*. `POOLDATA.TREATYDESC` = master jenis klausul treaty (kelas Pega `ASM-FW-GISFW-Int-TREATYDESC`). Pega
membacanya lewat `BrowseTreatyDesc_RD` (grid "For Non XOL" / "For XOL" Treaty Contract Out) dan menulisnya lewat
prosedur `PEGA_TREATYDESC`; layar input Pega-nya **tidak ada di korpus** — tampilan modul ini rancangan sendiri
(keputusan work owner: "buat versi kamu"). Baris menunya dibuat migrasi inti `920`, dan `Folder korpus` di bawah =
label menu `M_NAV_MENU.LABEL` (golongan MASTER TREATY).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu — `—` = modul TANPA
migrasi sendiri (`tandaTanpaMigrasi`): seluruh nomor modul sudah terbagi dan modul lain tidak boleh disentuh (perintah
work owner 05-10-2026: "JANGAN ADA SENTUH MODUL LAIN!"), jadi tidak meminjam nomor modul mana pun.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `treatydescription` |
| Folder korpus | `Treaty Description` |
| GROUPMENU | `MASTER TREATY` |
| Pemilik | `@PEMILIK-TREATYDESCRIPTION` |
| Status | dimigrasi |
| Rentang migrasi | `—` |
| Slot menu | `—` |
| Prefix rute API | `/api/treaty-description` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-TREATYDESCRIPTION.md` — peta tabel warisan, pemakainya, dan sumber aturan |
| `backend/` | `models/` `repository/` `services/` `handlers/` `tiruan/` `modul.go` (tanpa `migrations/`) |
| `frontend/` | `pages/` `components/` `labels.ts` `api.ts` `aturan.ts` `treatydescription.css` `menu.ts` `rute.tsx` dan `*.test.ts` |

## Aturan (ringkas) — keputusan work owner 05-10-2026

- Tulis langsung ke `POOLDATA.TREATYDESC`: nol tabel baru, nol DDL. Add dan Edit; **tanpa hapus** — ID dipakai
  `PROPORTIONALARRG` dan rule Pega klaim (`TREATYDESCID = '10001'`); penonaktifan lewat Status. Prosedur
  `PEGA_TREATYDESC` / `PEGA_M_TREATYDESC` tidak dipanggil; `M_TREATYDESC` (JSON Pega) tidak disentuh.
- ID baru = `'1' || LPAD(TREATY_DESCRIPTION_SEQ.NEXTVAL, 4, '0')`, persis `PEGA_TREATYDESC` ("harusnya ada
  sequencenya" — katalog DEV: `TREATY_DESCRIPTION_SEQ`); nomor yang sudah dipakai dilompati; ID tidak pernah diubah.
- Description Name wajib, huruf besar, maks. 100 byte, tidak kembar (tanpa beda huruf, juga dengan baris nonaktif).
  Nama boleh diubah, termasuk `10001` TREATY LIMIT; salinan `TREATYDESCNAME` di `PROPORTIONALARRG` tidak ikut diubah.
- Type: Non XOL (`ISXOL` `0`) / XOL (`1`). Status: Active (`STATUSAKTIF` `1`) / Inactive (`0`); NULL (baris lama Pega)
  dibaca Active dan ditulis `1`/`0` saat Edit.
- HANYA `TREATYDESC` yang dibaca dan ditulis — tidak ada kolom atau data dari tabel lain di layar (perintah work owner
  05-10-2026: "hanya baca dari tabel yang saya kasih, jangan ber-experiment").
- Jenis baru tampil di Treaty Contract Out, tetapi klausulnya baru bisa disimpan di sana setelah modul itu diberi
  aturannya — modul ini tidak menyentuh modul lain.
- Hak menu Full / View only (`M_LOGIN_GO_MENU.HAK`): View only tanpa Add dan Edit.

## Migrasi

Nol DDL, nol migrasi sendiri. Baris menunya dibuat migrasi inti `920_m_nav_menu_treatydescription.sql` dengan
`DIMIGRASI '1'` langsung (tanpa slot menu modul).

## Menjalankan uji modul ini saja

Dari folder `APP_RNM/`:

```powershell
go test ./modul/treatydescription/...
npx vitest run modul/treatydescription/
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga`.

### Tabel warisan: dibaca, tidak dibuat

| Tabel | Alasan |
| --- | --- |
| `TREATYDESC` | tabel warisan POOLDATA (master jenis klausul treaty); modul ini menambah dan mengubah barisnya, tidak pernah membuat atau mengubah strukturnya (nol DDL) |
