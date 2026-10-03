# Modul `nbtreatyin` — NB Treaty In

Realisasi penutupan treaty inward: dari kontrak treaty yang disetujui menjadi polis treaty, melalui
tangga **Admin → Sec Head → Dept Head** (`docs/spec.md`). Satu folder, satu modul, satu pemilik: kode
backend, kode frontend, dan dokumen modul ini tinggal di sini (struktur tim satu folder per modul,
keputusan work owner 30-09-2026). Padanan Pega: Harness `SFAPortalOpportunities` dan
`Flow/InputRealizationTreatyIn` — seluruh rule terjangkau beserta nasibnya (dibangun/tidak + alasan)
di `docs/INVENTARIS-XML.md`; hasil per acceptance criteria di `docs/HASIL-IMPLEMENTASI.md`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `nbtreatyin` |
| Folder korpus | `NB Treaty In` |
| GROUPMENU | `TREATY` |
| Pemilik | `@PEMILIK-NBTREATYIN` |
| Status | dimigrasi |
| Rentang migrasi | `320-359` |
| Slot menu | `968-969` |
| Prefix rute API | `/api/nb-treaty-in` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | spec, tiket (`issues/`), grilling, `INVENTARIS-XML.md` (bangkitan `docs/alat/`), `STRUKTUR-TABEL-NB-TREATY-IN.md`, `HASIL-IMPLEMENTASI.md`, `PERMINTAAN-TIM-INTI.md` |
| `backend/` | `models/` (fungsi murni: halaman kerja, katalog, rantai uang, tangga, layar) `repository/` (seluruh SQL) `services/` (aturan, transaksi) `handlers/` (HTTP) `migrations/` `modul.go` |
| `frontend/` | `pages/` `components/` `api.ts` `labels.ts` `medan.ts` `menu.ts` `rute.tsx` `nbtreatyin.css` dan `*.test.ts` |

## Aturan (ringkas)

- **Sumber data kontrak** = view `POOLDATA.TREATYINDETAILJOINEDM` (pilih bisnis) dan tabel
  `TREATYINDETAIL` (grid popup) — bukan JSONDATA (P29). Nol penulisan JSON.
- **Penyimpanan** = `T_WORK_POLIS` (tabel kasus lintas-lini, migrasi premiumlistlife 050/057/059 —
  modul ini MEMBUTUHKAN premiumlistlife aktif lebih dulu) + `T_GENERAL_POLIS` dan anak `T_POLIS_*`
  (migrasi 320-327, tepat delapan tabel diagram grilling) menurut katalog `backend/models/katalog.go`;
  satu transaksi per tindakan.
- **Tangga** tiga posisi; berkas menunggu POSISI (workbasket), keanggotaan diperiksa menurut nama
  antrean (`inti.Pelaku.Peran`). Sec Head menyetujui → selalu Dept Head (AC 8, keputusan work owner —
  bertentangan dengan `CekLimitTreatyAcc_Act`, dicatat). Dept Head menerbitkan nomor polis
  (`inti/backend/penomor`, `KODE_PRODUKSI` NONLIFE).
- **Riwayat** `HISTORYAKSEPTASIPEGA` di transaksi submit; `OPERATORID` = identitas login, `USERNAME`
  = nama tampilan (`M_LOGIN_GO.NAME`).
- **Catatan usulan** (`PolicyTreatyIn.SuggestList`) = tabel lama `POOLDATA.HISTORYAKSEPTASIPRODUCTION`
  (keputusan work owner K4 03-10-2026; pengganti `SaveViewSuggest -> InsertViewSuggest_SQL`), ditulis
  di transaksi submit SETIAP jenjang dan dibaca balik untuk layar - `[penyimpangan sadar]` di
  `backend/models/usulan.go`. Tabel warisan: tidak dibuat, tidak diubah strukturnya.
- **Peran pengganti nama orang** = konstanta kode `backend/models/peran_tempat.go`
  (`PemetaanPeranTempat`, keputusan work owner K16 03-10-2026 - BUKAN tabel), KOSONG sampai IAM
  menjawab; tanpa baris = tempat tertunda (tiket 05). Peran pengguna dari `inti.Pelaku.Peran`.

## Migrasi

Rentang `320-359`: 320-327 - TEPAT delapan tabel diagram grilling (`Diagram-Skema-Tabel-NusantaraRe.xlsx`
sheet *NB Treaty In Prop* / *NonProp*), bangkitan `docs/alat/skema.py` dari katalog. Tidak ada tabel lain
(bab 0 butir 11 PROMPT putaran 2; K4, K16, K17): catatan usulan ke tabel warisan di bawah, pemetaan
peran-tempat konstanta kode, medan tak dikenal pemuat dokumen lama ke berkas laporan CSV.
Perbandingan kolom lawan diagram: `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`. Slot menu `968`: satu `UPDATE DIMIGRASI` baris modul ini, nol `INSERT`.

## Menjalankan uji modul ini saja

Dari akar repo:

```powershell
go test ./modul/nbtreatyin/...
npx vitest run modul/nbtreatyin
```

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per
butir. Judul yang tidak ada berarti modul ini tidak menyatakan apa pun untuk jenis itu.

### Tabel warisan: dibaca, tidak dibuat

Tabel yang `docs/STRUKTUR-TABEL-NB-TREATY-IN.md` gambarkan tetapi SENGAJA tidak dibuat migrasi mana pun
(`TestKolomDDLCocokDenganStruktur`, `TestTabelBukanMilikKitaTidakDibuat`). Mencabut satu baris =
kepemilikan tabel berpindah — keputusan work owner.

| Tabel | Alasan |
| --- | --- |
| `HISTORYAKSEPTASIPRODUCTION` | tabel warisan POOLDATA (15 kolom, sudah datar); catatan `SuggestList` ditulis dan dibaca di sini menurut diagram grilling (`SaveViewSuggest -> InsertViewSuggest_SQL`) dan keputusan work owner K4 03-10-2026 — tidak dibuat, tidak diubah strukturnya |
