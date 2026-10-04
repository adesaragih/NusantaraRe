# 04: Jenis reasuransi — enumerasi `Flag = 1` dan normalisasi `REINSTYPEID`

**Status:** ready-for-agent

**Blocked by:** 02 (jenis reasuransi melekat pada kontrak)

## Hasil & nilai pengguna

Sebagai **admin master**, saya memilih jenis reasuransi dari **daftar yang sah untuk lini life** dan
melihat **namanya**, bukan kodenya — sehingga tidak ada jenis dari lini lain yang menyelinap masuk
dan layar terbaca manusia. *(User story 30–31 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Jenis reasuransi; `REINSTYPENAME` ditandai **turunan** |
| `internal/repository` | Pembacaan enumerasi dengan penyaring `Flag = 1` |
| `internal/services` | Penolakan nilai di luar enumerasi; **normalisasi** — anak mewarisi dari induk |
| `internal/handlers` | Endpoint daftar jenis reasuransi |
| `frontend/` | Dropdown jenis reasuransi menampilkan nama |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `BrowseReinsuranceTypeLimit_RD` | `ASM-FW-GISFW-INT-REINSURANCETYPE` / `BROWSEREINSURANCETYPELIMIT_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseReinsuranceTypeLimit_RD.xml` | **sumber enumerasi**, penyaring `.Flag = 1` |
| `TreatyLimit_TypeProtect` | `@BASECLASS` / `TREATYLIMIT_TYPEPROTECT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/TreatyLimit_TypeProtect.xml` | ⚠️ **lookup**, bukan rumus |

`[terverifikasi]` `TreatyLimit_TypeProtect` mencocokkan `REINSTYPEID` dengan `.ID`, lalu menyalin
**`.Note`** ke `REINSTYPENAME`. Kolom yang tersedia di class: `.ID`, `.Code`, `.Note`, `.Type`,
`.Flag`, `.SOANote`, `.UserID`. Penyaring RD: **`.Flag = 1`** (baris 581/585 dan salinan indeks
887/893).

`[fakta bisnis — work owner]` **Lima jenis untuk lini life:**

| ID | Nama (`.Note`) |
| --- | --- |
| 10196 | QS |
| 10197 | 2ND QS |
| 10198 | SURPLUS |
| 10199 | 2ND SURPLUS |
| 10200 | OR |

`[fakta bisnis — work owner]` **`Flag = 1` = "for life"** — penanda lini, **bukan** aktif/nonaktif.
Sejalan dengan memo korpus `[terverifikasi]` `pyUsage: "Parameter Flag, 1 for life"`.
Kolom `.Code` dan `.Type` **tidak dipakai** sistem.

⚠️ `"OR"` dicatat **apa adanya**; artinya **tidak ditafsirkan**.

`[terverifikasi]` **Denormalisasi di Pega**: `REINSTYPEID`/`REINSTYPENAME` tersimpan **tiga kali**
(kontrak, reinsurer, business) dan `TREATYYEAR` **dua kali** (tahun, business), **tanpa penjaga
konsistensi** — dan `[data DBA]` procedure **tetap menulisnya** ke tabel anak.

## ADR terkait

**ADR-0001** (batas konteks — enumerasi ini dipakai lintas lini), **ADR-0009**.

## Acceptance criteria

- [ ] Daftar jenis reasuransi yang ditawarkan **hanya** yang **`Flag = 1`** — lima jenis, ID
      **10196**–**10200**. *(AC 38 spec)*
- [ ] Layar menampilkan **nama** dari kolom `.Note`, **bukan** kode. *(AC 39 spec)*
- [ ] `REINSTYPEID` di luar kelima nilai itu **ditolak** saat menyimpan kontrak. *(AC 40 spec)*
- [ ] ⚠️ **Normalisasi**: `REINSTYPEID` tersimpan **hanya pada kontrak**; reinsurer dan business
      mewarisinya lewat `TREATYCONTRACTID`. Test yang menemukan `REINSTYPEID` yang **dapat ditulis**
      pada tabel anak **gagal**. *(AC 41 spec; `[keputusan work owner]` — penyimpangan sadar 1)*
- [ ] ⚠️ `TREATYYEAR` tersimpan **hanya pada tahun treaty**; business mewarisinya.
      *(AC 42 spec; `[keputusan work owner]`)*
- [ ] `REINSTYPENAME` yang disimpan ditandai **turunan/cache** — **bukan sumber kebenaran** — dan
      **selalu dapat dibangun ulang** dari `REINSTYPEID`. *(AC 43 spec)*
- [ ] Mengubah jenis reasuransi pada kontrak **tidak** menuntut pembaruan manual di tabel anak —
      anak membacanya lewat induk.
- [ ] Enumerasi dibaca dari basis data, **tidak ditanam** sebagai konstanta di kode. Test yang
      menemukan kelima ID sebagai literal di lapisan services **gagal**.

## Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **OQ kecil** — **tidak memblokir**: nama **tabel fisik Oracle** untuk class
`ASM-FW-GISFW-INT-REINSURANCETYPE` belum diketahui (korpus hanya menyebut *class*), dan isi kolom
`.Code` serta `.Type` belum terbaca. `[fakta bisnis — work owner]` keduanya **tidak dipakai sistem**,
jadi tiket ini dapat selesai tanpa jawaban itu. Pemilik: **DBA**.

## Catatan

`[terverifikasi]` RD serupa (`BrowseReinsuranceType_RD`, class yang sama) ada di **sembilan modul
lain** — enumerasi ini **dipakai lintas lini**, bukan khas life. Versi modul lain difilter lewat
`Param.ID` / `Param.Name` / `Param.Note` / `Param.Type`, **bukan** `.Flag`.

⚠️ **Nama menyesatkan.** `TreatyLimit_TypeProtect` terbaca seperti logika "jenis proteksi limit";
`[terverifikasi]` ia **hanya lookup nama**. Jangan tiru namanya.

⚠️ **Normalisasi adalah keputusan skema baru, bukan tiruan.** `[data DBA]` Procedure **tetap
menerima dan menulis** `REINSTYPEID`/`REINSTYPENAME` ke tabel anak — jadi **lapisan repository yang
memutuskan apa yang dikirim**.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Status 01-10-2026 (paket 11)

**Status:** ✅ **dibangun** — master `REINSURANCETYPE` `.Flag = 1` (paket 1 `db7f636`, paket 3 `38c3d42`), dropdown `REINS TYPE` paket 9+10 (`a3c07bc`).
