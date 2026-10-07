# 01: Kolom khas endorsemen pada tabel yang sudah ada

> ## ⛔ KOREKSI 06-10-2026 — tabel NB nyata (log: `../KOREKSI-DOKUMEN-2026-10-06.md`)
>
> | Bunyi lama (dikutip) | Bunyi baru | Bukti |
> | --- | --- | --- |
> | *"**21 kolom khas endorsemen** — terisi pada endorsemen, ⛔ **kosong** pada polis baru"* | ⛔ **Bertentangan dengan NB nyata.** Dari daftar 21 (spec-penyimpanan bab 5.5), NB menulis sebagian besar: `PRODKE` (= 0, `NOT NULL DEFAULT 0`), `GROSS_PREMIUM`/OGP/ONP (rantai `CountNetPremi`), `QUARTAL`, `STATEMENT_TYPE`, `ID_NEW_BISNIS` (= ID kasus), `SPLIT_RNM_SHARE_PCT`. Nama nyata: `EDM_NO` → **`NOENDORS`**, `PROD_KE` → **`PRODKE`**. Yang terbukti kosong di NB hanya **`NOENDORS`** dan **`OLD_POLIS_ID`** (`PRODKE` = 0). ⇒ daftar kolom khas AC 32 **wajib disusun ulang** dari katalog NB — butir asisten utama | `modul/nbtreatyin/backend/migrations/320_t_general_polis_treaty.sql` baris 33–35; `backend/models/katalog.go` baris 137, 144, 160, 172–184; `models/hitung.go` baris 186–188; `models/layar.go` baris 128, 334; `models/produksi.go` baris 63 (`NOENDORS = NULL`), 199; `models/nonprop.go` (`SplitRNMSharePct`) |
> | *"**4 kolom khas polis baru**"* (`IS_EDM_INPUT_ON_NB` · `IS_APPROVEDTO_DEPT_HEAD` · `HAS_FAC_OUT` · `SHARE_CURRENCY`) | **3 kolom** — `IS_APPROVEDTO_DEPT_HEAD` **tidak ada** (P36, tidak dibangun) | 320 baris 49–50, 74; `modul/nbtreatyin/docs/spec.md` AC 64; `PERTANYAAN-untuk-Product-dan-Underwriting.md` baris 362 |
> | *"bentuk kesembilan tabel dasar"* · *"Sepuluh tabel dasar dibuat tiket NB 16 dan 19"* | tabel dasar = `T_WORK_POLIS` (premiumlistlife) + `T_GENERAL_POLIS_TREATY` + delapan anak `T_POLIS_*` (321–328, **`T_POLIS_SURVEY` ikut**) — dibuat **migrasi NB 320–328** | `modul/nbtreatyin/backend/migrations/`; 328 baris 1–3 |
> | AC 55 *"cacah tabel … sesuai ketetapan"* (10 / 14) | cacah 10 / 14 = daftar resmi diagram (memuat `HISTORYAKSEPTASIPRODUCTION`, tanpa `T_POLIS_SURVEY`). ⚠️ Korpus EDM **nol rujukan** `HISTORYAKSEPTASIPRODUCTION`; `T_POLIS_SURVEY` di luar diagram ⇒ cacah **belum pasti** — butir WO | `Diagram-Skema-Tabel-NusantaraRe.xlsx` sheet *EDM Treaty In Prop* K146–K156, *NonProp* K171–K185; sapuan korpus 163 xml (bab 11 `../KEADAAN-EDM-TREATY-IN.md`) |

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

- [ ] **AC 32** — ~~dua puluh satu kolom khas endorsemen **kosong** pada baris polis baru~~ ⛔ *(koreksi 06-10)* daftar disusun ulang dari katalog NB; yang pasti kosong di NB: `NOENDORS`, `OLD_POLIS_ID` (`PRODKE` = 0)
- [ ] **AC 33** — ~~empat~~ **tiga** kolom khas polis baru **kosong** pada endorsemen, dan **tidak** dihapus dari skema *(koreksi 06-10: `IS_APPROVEDTO_DEPT_HEAD` tidak ada)*
- [ ] **AC 34** — keadaan layar **tidak tersimpan** di tabel mana pun
- [ ] **AC 54** — bentuk ~~kesembilan~~ tabel dasar NB 320–328 **tidak berubah** *(koreksi 06-10)*
- [ ] **AC 55** — cacah tabel pada kedua bentuk endorsemen sesuai ketetapan
