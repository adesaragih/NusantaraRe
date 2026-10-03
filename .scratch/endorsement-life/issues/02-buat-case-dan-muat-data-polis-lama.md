# 02: Buat case endorsement + muat & salin data polis lama

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 01 (gerbang kelayakan — case hanya dibuat setelah kelimanya lolos)

## Hasil & nilai pengguna

Sebagai **inputor Life**, setelah polis dinyatakan layak saya menekan "buat endorsement" dan seluruh
data polis lama — termasuk ratusan barisan peserta — **otomatis tersalin** ke endorsement baru,
sehingga saya tidak mengetik ulang apa pun dan dapat langsung bekerja pada salinannya.
*(User story 9, 12 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas endorsement; rujukan ke polis lama (ID pega lama + nomor polis lama); baris detail beserta `EDMStatus` |
| `internal/repository` | Ambil `IDPEGA` polis terakhir dari `JSON_POLIS`; baca polis NB beserta detailnya |
| `internal/services` | Orkestrasi: gerbang → buat case → muat → salin → tandai `Old` |
| `internal/handlers` | Endpoint buat endorsement |
| `frontend/` | Tombol buat endorsement; grid detail hasil salinan |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `MappingEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `MAPPINGEDMLIFE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/MappingEDMLife.xml` (194.433 byte) | **mesin inti, 14 langkah** |
| `CreateCaseEMDL` | `DATA-PORTAL` / `CREATECASEEMDL` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/CreateCaseEMDL.xml` (119.680 byte) | pembuat case |
| `GetProdkeNopolis` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETPRODKENOPOLIS` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/GetProdkeNopolis.xml` | ambil `IDPEGA` polis terakhir |
| `InputEDMLife` (Flow) | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` / `RULE-OBJ-FLOW` | `Endorsement Life/Flow/InputEDMLife.xml` (68.374 byte) | alur case |

`[terverifikasi]` **Peta 14 langkah `MappingEDMLife`** — **nol `<pyStepsBlockName>` di berkas ini**:

| Step | Langkah | Catatan |
| ---: | --- | --- |
| 1 | *(gerbang)* | "Jika ada edm blm resolve exit act" — `@LengthOfPageList(ListEdm.pxResults)>0` (baris 456) |
| **3** | `Call svcAddWorkObject` | **membuat case EDM baru** |
| **6** | `RDB-List` → `GetProdkeNopolis` (baris 1172) | "Ambil idpega prodke terakhir dari tabel json_polis" |
| **7** | `Property-Set` | **"Set Old ID Pega dan Old Policy No"** |
| **8** | `Obj-Open-By-Handle` | **"Open IDPEGA NB Life"** — memuat polis new business |
| **9** | `Property-Set` | **"Mapping detail dari Life"** |
| **10** | `Property-Set` | **"Copy page dari Life ke EDM"** |
| 11.1 | `Property-Set` | **Set `"Old"` untuk detail lama** (baris 2608) |
| 11.2 | `Property-Remove` | buang baris ber-`EDMStatus == "Delete"` |
| 12 | `Property-Set` | set `.EditInput = 1` (baris 2956) — **kunci field**, lihat tiket 08 |
| 13–14 | `Obj-Save`, **`Commit`** | **AKTIF** |

`[terverifikasi]` Kunci pembacaan: `SELECT IDPEGA … WHERE NOPOLIS={TempWork.PolicyNo} AND PRODKE IS
NOT NULL ORDER BY **PRODKE DESC**`.

## ADR terkait

**ADR-0001** (batas konteks — endorsement membaca polis NB sebagai hulu), **ADR-0003** (uang
non-float pada detail tersalin), **ADR-0007** (jejak audit), **ADR-0015** (batas transaksi).

## Acceptance criteria

- [ ] Case endorsement dibuat **hanya setelah kelima gerbang lolos**; gerbang yang gagal **tidak**
      meninggalkan case separuh. *(AC 6 spec)*
- [ ] Data polis new business **termuat otomatis**; seluruh baris peserta tersalin ke endorsement.
      *(AC 7 spec)*
- [ ] Baris warisan polis lama bertanda **`Old`**. *(AC 7 spec)*
- [ ] **ID Pega lama dan nomor polis lama tersimpan** pada case endorsement — keduanya dipakai tiket
      03 (popup) dan tiket 04 (penomoran). *(AC 8 spec)*
- [ ] Baris polis new business **tidak pernah dihapus fisik** oleh proses apa pun. *(AC 10 spec)*
- [ ] Nilai uang tersalin **tanpa berubah** — tidak lewat `float`, tidak ada pembulatan diam.
      *(AC 21–22 spec; **ADR-0003**)*
- [ ] Case yang sudah dibuat **terlihat oleh pengguna lain** sejak saat itu, sehingga gerbang "satu
      endorsement terbuka per polis" (tiket 01) menolak percobaan kedua. `[keputusan work owner]`
      — commit dini di Pega **dipertahankan dengan sengaja** karena gerbang itu membutuhkannya.
- [ ] Menyalin polis dengan ratusan peserta selesai dalam satu permintaan tanpa memotong daftar.

### Penyalinan versi ⚠️ BARU 2026-09-16 — spec §16

- [ ] ⚠️ Submit **menyalin** header, seluruh peserta, seluruh spreading, dan seluruh spreading retro
      ke **versi baru**, masing-masing dengan **identitas baru**. *(AC 59 spec; penyimpangan sadar 9)*
- [ ] ⚠️ Setiap peserta hasil salin membawa **`PARENT_ID`** yang menunjuk peserta versi sebelumnya.
      *(AC 60 spec; penyimpangan sadar 10)*
- [ ] ⚠️ **Rekap mata uang TIDAK disalin** — ia **dihitung ulang** dari peserta versi baru. Test yang
      menemukan baris rekap tersalin **gagal**. *(AC 63 spec)*
- [ ] ⚠️ **Riwayat penawaran tidak disalin** ke versi endorsement. *(AC 64 spec)*
- [ ] ⚠️ **`OldData` tidak disimpan** — data lama dibaca dari **versi sebelumnya** (`PRODKE` lebih
      kecil, lewat `PARENT_ID`). Test yang menemukan tabel/kolom penyimpan `OldData` **gagal**.
      *(AC 69 spec)*
- [ ] ⚠️ Pencocokan peserta **tidak** bergantung urutan baris maupun `CERTIFICATE_NO`: menyisipkan
      dan menghapus peserta sehingga urutan bergeser **tetap** menghasilkan pasangan yang benar.
      *(AC 62 spec; penyimpangan sadar 10 — inilah kelas bug yang diperbaiki)*
- [ ] Versi baru ber-`PRODKE` **satu lebih besar** dari versi berjalan. *(AC 56 spec)*

## Catatan — kode mati yang tidak dimigrasikan

`[keputusan work owner]` **`CreateCaseEMDL` step 7–11 tidak direplikasi** — penyalinan data lama di
sana **mati**: `<pyStepsBlockName>//` di baris 1084 ("-- Set data dari json_polis"), 1273 (`Java`),
1382 ("-- Copy List old"), 1532 ("-- Copy Ke new"), 1671. Deskripsi langkah mati diawali `--`.
Yang hidup adalah jalur `MappingEDMLife` di atas.

⚠️ **OQ-066 berlaku.** Penanda `<pyStepsBlockName>` **tidak dapat dipercaya sendirian** di modul
Life — penetapan "mati" di atas bersandar pada **keputusan work owner**, dengan penanda sebagai
pendukung. Sebelum menyimpulkan langkah lain hidup atau mati, konfirmasikan ke work owner.

## Blocker

**Tidak ada.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
