# 05: Validasi wajib-isi — lima pemeriksaan, seragam di kedua sisi

**Status:** ready-for-agent

**Blocked by:** 03 (sisi inward — pemegang polis diperiksa dari sana), 04 (sumber bisnis dan ceding
berasal dari pemilih master)

## Hasil & nilai pengguna

Sebagai **admin master**, saya ingin produk **ditolak** bila field penentu belum terisi — dengan
pesan yang menyebut apa yang kurang — dan saya ingin aturannya **sama** dari jalur mana pun saya
menyimpan, sehingga tidak ada pintu yang lebih longgar. *(User story 1, 10 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/services` | **Satu** aturan wajib-isi, dipanggil kedua jalur simpan |
| `internal/handlers` | Pesan penolakan menyebut field yang kurang |
| `frontend/` | Galat tampil di dekat field yang bersangkutan |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `SaveProductName_Act` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SAVEPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SaveProductName_Act.xml` |
| `SaveInwardProductName_Act` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SAVEINWARDPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SaveInwardProductName_Act.xml` |

`[terverifikasi]` **Lima pemeriksaan di `SaveProductName_Act`:**

| Step | Pesan | Precondition korpus |
| ---: | --- | --- |
| 1 | — | `ProductName.TYPE=="" \|\| ProductName.GRUP==""` |
| 2 | "error Product Name" | `ProductName.PRODUCTNAME==""` |
| 3 | "error Ceding" | `ProductName.CEDING==""` |
| 4 | "error Policy Holder" | `ProductNameInward.POLICYHODERNAME==""` |
| 5 | "error SOB" | `ProductName.SOBNAME==""` |

⚠️ `[terverifikasi]` Di **`SaveInwardProductName_Act`**, dua pemeriksaan padanannya **ter-remark**:
step 2 "error Type" (`<pyStepsBlockName>//` baris 414) dan step 3 "error Grup" (baris 601).
**Jalur inward lebih longgar dari jalur utama.**

⚠️ `[data DBA]` Basis data **tidak menegakkan apa pun** — seluruh kolom hasil migrasi **nullable**
(tiket 01). **Wajib-isi seluruhnya ditegakkan di Go.**

## ADR terkait

**ADR-0007** (jejak audit penolakan), **ADR-0009** (migrasi penuh).

## Acceptance criteria

- [ ] Produk **ditolak** bila **tipe** atau **grup** kosong. *(AC 11 spec)*
- [ ] Produk **ditolak** bila **nama produk** kosong. *(AC 12 spec)*
- [ ] Produk **ditolak** bila **ceding** kosong. *(AC 13 spec)*
- [ ] Produk **ditolak** bila **pemegang polis** kosong. *(AC 14 spec)*
- [ ] Produk **ditolak** bila **sumber bisnis** kosong. *(AC 15 spec)*
- [ ] ⚠️ Kelima pemeriksaan berlaku **SAMA pada kedua sisi produk** — **tidak ada** sisi yang lebih
      longgar. Test yang menemukan jalur simpan dengan pemeriksaan lebih sedikit **gagal**.
      *(AC 16 spec)*
- [ ] Wajib-isi ditegakkan **di Go**, bukan oleh basis data — seluruh kolom hasil migrasi nullable.
      Test yang mengandalkan basis data untuk menolak nilai kosong **gagal**. *(AC 17 spec;
      `[data DBA]`)*
- [ ] Pesan penolakan **menyebut field** yang kurang, bukan galat umum.
- [ ] Bila **beberapa** field kurang sekaligus, **semuanya** dilaporkan — bukan hanya yang pertama.
- [ ] Aturan wajib-isi berada di **satu tempat**, dipanggil kedua jalur — bukan disalin dua kali.

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Pega lebih longgar di satu sisi, dan itu tidak ditiru.** `[terverifikasi]` Dua pemeriksaan pada
jalur inward ter-remark. Sistem baru **menyeragamkan** — dua jalur atas satu produk tidak boleh
berbeda ketatnya.

⚠️ **OQ-066 diterapkan:** remark itu **tidak** saya simpulkan sendiri sebagai kelalaian; penyeragaman
adalah **keputusan** yang tercatat di spec §8, bukan tafsir dari penanda korpus.

⚠️ **Perubahan perilaku yang disadari.** Bersama pembuangan salah ketik `PoductName` (tiket 02),
penegakan ini akan **mulai menolak** penyimpanan yang dahulu lolos. Itu disengaja.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Ralat bertanggal 01-10-2026 — sesi implementasi (paket 5)

> Sumber: `../RALAT-DEV-30-09-2026.md` (R7, R11, R13) dan `../PARITAS-LAYAR-DAN-AKSI.md` §7. Kalimat di atas **tidak dihapus**.

| Kalimat lama | Ralat |
| --- | --- |
| tabel *"Step 1 — `ProductName.TYPE=="" ‖ ProductName.GRUP==""`"*; AC *"Produk **ditolak** bila **tipe** atau **grup** kosong"* | **R7**: langkah 1 b359 adalah `Property-Set` ber-**PRE=false** — bukan pemeriksaan; medan `TYPE` b3143 / `GRUP` b6108 mati (`1=2`). AC ini **dicabut**: menegakkannya membuat setiap simpan gagal |
| Pesan *"error Product Name"*, *"error Ceding"*, *"error Policy Holder"*, *"error SOB"* | **R13**: itu `pyStepsDescription`. Pesan layar VERBATIM: `Product Name Empty` (b431), `Ceding Empty` (b452), `Policy Holder Empty` (b473), `SOB Empty` (b494) |
| *"Di **`SaveInwardProductName_Act`**, dua pemeriksaan padanannya **ter-remark** … Jalur inward lebih longgar"*; AC *"berlaku SAMA pada kedua sisi"* | **R11**: jalur inward tak terjangkau (`Inward` b75322 `1=2`). Satu jalur simpan (`POST`/`PUT /produk`) memanggil satu aturan `periksaWajibIsi` — uji: baru dan ubah ditolak dengan pesan yang sama |
| AC *"Bila **beberapa** field kurang sekaligus, **semuanya** dilaporkan"* | dibangun (keputusan tertulis): Pega berhenti di pesan pertama (transisi `1==1` → 6); di sini keempatnya dilaporkan, urutan langkah 2–5, dipisah `; `, bersama penolakan lain (angka, tanggal, pilihan master). Spasi = kosong |
| AC *"Galat tampil di dekat field"* | amplop galat bersama hanya membawa `galat` (teks); layar memetakan pesan VERBATIM ke medannya (paket 10) |
