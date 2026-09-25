# 02: Kontrak treaty — batas proteksi per mata uang dan lebar layer yang dihitung

**Status:** ready-for-agent

**Blocked by:** 01 (kontrak lahir di bawah tahun treaty; pola jalur simpan ditetapkan di sana)

## Hasil & nilai pengguna

Sebagai **underwriter**, saya menetapkan **batas bawah dan batas atas** proteksi per mata uang pada
sebuah kontrak treaty, dan **lebar layer dihitung sistem** dari kedua batas itu — sehingga angkanya
tidak pernah menyimpang dari batasnya sendiri. *(User story 5–11 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas kontrak; nilai batas sebagai desimal presisi arbitrer |
| `internal/repository` | Pemanggilan procedure penulis kontrak; pemeriksaan `o_message` |
| `internal/services` | Hitung lebar layer; wajib-isi; pemeriksaan batas bawah ≤ batas atas |
| `internal/handlers` | Endpoint CRUD kontrak |
| `frontend/` | Grid kontrak + batas proteksi; lebar layer **hanya tampil**, tidak dapat diketik |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterTreatyContract_Life_SQL` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` / `ASM!SAVEMASTERTREATYCONTRACT_LIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/SaveMasterTreatyContract_Life_SQL.xml` | `POOLDATA.INSERTTREATYCONTRACT_LIFE` |
| `SaveTreatyLimit_Act` | `ASM-FW-GISFW-…` / `SAVETREATYLIMIT_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveTreatyLimit_Act.xml` | orkestrator simpan + wajib-isi |
| `NewInputTreatyLimit_Life` | `ASM-FW-GISFW-…` / `NEWINPUTTREATYLIMIT_LIFE` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/NewInputTreatyLimit_Life.xml` | baris baru |
| `SetValueRetroLimit_TreatyYearLife` | `@BASECLASS` / `SETVALUERETROLIMIT_TREATYYEARLIFE` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SetValueRetroLimit_TreatyYearLife.xml` | ⚠️ **penyalur konteks**, bukan rumus |
| `BrowseTreatyContract_Life_RD` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` / `BROWSETREATYCONTRACT_LIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseTreatyContract_Life_RD.xml` | daftar |
| `InboxRetroLimitReinsurers` | `DATA-PORTAL` / `INBOXRETROLIMITREINSURERS` / `RULE-HTML-HARNESS` | `Master Contract Retro Life/Harness/InboxRetroLimitReinsurers.xml` (575.179 byte) | layar grid |

`[fakta bisnis — work owner]` **Arti keenam kolom limit:**

| Kolom | Arti | Tipe `[data DBA]` |
| --- | --- | --- |
| `B_IDR` / `B_USD` | **batas bawah** (prefix `B_` = bawah) | `NUMBER` |
| `IDR` / `USD` | **batas atas** | `NUMBER` |
| `IDR_SELISIH` / `USD_SELISIH` | **atas − bawah** = **lebar layer** | `NUMBER` |

`[fakta bisnis — work owner]` **Layer tersusun menaik** — batas atas satu layer menjadi batas bawah
layer berikutnya: `QS → 2nd QS → Surplus → 2nd Surplus`.

`[terverifikasi]` Wajib-isi (`SaveTreatyLimit_Act`): `REINSTYPEID`, `TREATYSTARTDATE`,
`TREATYENDDATE`, `B_IDR`, `IDR`, `B_USD`. ⚠️ **`USD` TIDAK wajib** — `[data DBA]` kolomnya nullable,
jadi ini **aturan sah**, bukan kelalaian.

⚠️ `[data DBA]` `INSERTTREATYCONTRACT_LIFE` **menulis `IDR_SELISIH`/`USD_SELISIH` apa adanya dari
parameter** — basis data **tidak menghitungnya**. Tidak ada yang menjaganya sinkron dengan kedua
batas.

## ADR terkait

**ADR-0003** (uang non-float — `[data DBA]` kolom bertipe `NUMBER` tanpa presisi, jadi desimal
presisi arbitrer aman), **ADR-0006** (identitas dari basis data), **ADR-0007** (jejak audit).

## Acceptance criteria

- [ ] Kontrak lahir **di bawah** satu tahun treaty; kontrak tanpa tahun induk **ditolak**.
      *(AC 5 spec)*
- [ ] Kontrak wajib memuat `REINSTYPEID`, tanggal mulai, tanggal akhir, **batas bawah IDR**, **batas
      atas IDR**, dan **batas bawah USD** — ditegakkan **di Go**. *(AC 6 spec)*
- [ ] **`USD` boleh kosong** — kontrak ber-IDR saja tersimpan tanpa keluhan. *(AC 7 spec;
      `[data DBA]` aturan sah)*
- [ ] ⚠️ **Lebar layer DIHITUNG** dari `batas atas − batas bawah`, per mata uang. **Tidak ada** jalur
      yang membiarkan pengguna mengetiknya, dan **tidak ada** jalur yang mengirimkannya mentah dari
      masukan. *(AC 8 spec; `[keputusan work owner]` — penyimpangan sadar 2)*
- [ ] Lebar layer yang tersimpan **selalu** sama dengan selisih kedua batas, **termasuk setelah salah
      satu batas diubah**. *(AC 9 spec)*
- [ ] Batas bawah yang **lebih besar** dari batas atas **ditolak**, dengan pesan yang menyebut mata
      uangnya. *(AC 10 spec)*
- [ ] Susunan layer menaik dapat direkam: batas atas satu layer boleh menjadi batas bawah layer
      berikutnya **tanpa** dianggap tumpang tindih. *(AC 11 spec)*
- [ ] Seluruh nilai batas diperlakukan sebagai **desimal presisi arbitrer**; **tidak ada** yang
      melewati `float` di lapisan mana pun maupun di JSON API. *(AC 12 spec; **ADR-0003**)*
- [ ] Nilai batas yang ditulis dan dibaca kembali **identik** — tidak ada pembulatan diam, termasuk
      pada nilai berpecahan panjang. *(AC 13 spec)*
- [ ] Menyimpan kontrak yang sudah ada = **upsert**, bukan baris kedua; identitas baru dibuat basis
      data. *(AC 45, 4 spec)*
- [ ] `o_message` diperiksa; kegagalan ditampilkan. *(AC 46 spec — HTML dan konformansi di
      **tiket 10**)*

## Blocker

**Tidak ada.**

## Catatan

⚠️ **`REINSTYPEID` diterima apa adanya di tiket ini.** Enumerasi yang sah (`Flag = 1`, lima nilai),
penolakan nilai di luar itu, dan **normalisasi** (`REINSTYPEID` hanya di kontrak) ditambahkan
**tiket 04**. Sampai tiket 04 selesai, kontrak menyimpan nilai apa pun yang dikirim.

⚠️ **Nama menyesatkan.** `SetValueRetroLimit_TreatyYearLife` terbaca seperti rumus limit; `[terverifikasi]`
ia **hanya menyalin konteks tahun treaty** ke tiga halaman input dan mengatur tampilan
(`OutputParam.DATASHOW` / `STSSAVE`). **Tidak ada perhitungan.** Jangan tiru namanya.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — procedure tidak di-mock.

```
go test ./internal/...
cd frontend && npm test
make check
```
