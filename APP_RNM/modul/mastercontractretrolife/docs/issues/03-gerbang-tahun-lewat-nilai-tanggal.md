# 03: Gerbang konsistensi tahun — dibandingkan dari nilai tanggal, bukan potongan teks

**Status:** ready-for-agent

**Blocked by:** 02 (gerbang ini berlaku saat menyimpan kontrak)

## Hasil & nilai pengguna

Sebagai **admin master**, saya ingin sistem menolak kontrak yang **tanggal mulainya jatuh di tahun
berbeda** dari tahun treaty induknya, sehingga kontrak tidak pernah nyasar tahun dan laporan per
tahun treaty tidak bocor. *(User story 9 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/services` | Aturan gerbang tahun — perbandingan tahun dari nilai tanggal |
| `internal/handlers` | Pesan penolakan yang menyebut kedua tahun |
| `frontend/` | Galat tampil di form kontrak, di dekat field tanggal mulai |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `SaveBusinessLife_Act` | `ASM-FW-GISFW-…` / `SAVEBUSINESSLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveBusinessLife_Act.xml` |
| `SaveSecurityLife_Act` | `ASM-FW-GISFW-…` / `SAVESECURITYLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveSecurityLife_Act.xml` |
| `SaveSecurityReinsurerLife_Act` | `ASM-FW-GISFW-…` / `SAVESECURITYREINSURERLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveSecurityReinsurerLife_Act.xml` |

`[terverifikasi]` **Ketiganya berbagi gerbang yang sama:**

```
@substring(InputRetrocessionLife.TREATYSTARTDATE,6,10) <> InputRetrocessionLifeTreatyType.TREATYYEAR_LIFE
```

yakni **memotong teks tanggal pada posisi 6–10** lalu membandingkannya dengan tahun treaty.

⚠️ `[data DBA]` **DDL membuktikan cara itu keliru sejak awal**: `TREATYSTARTDATE`, `TREATYENDDATE`,
`STARTDATE`, `ENDDATE`, dan `TGLUPDATE` seluruhnya bertipe **`DATE`** — bukan teks. Pemotongan posisi
tetap atas nilai `DATE` bergantung pada format tampilan, yang dapat berubah.

## ADR terkait

**ADR-0007** (jejak audit — penolakan tercatat), **ADR-0009** (migrasi penuh).

## Acceptance criteria

- [ ] ⚠️ **Gerbang tahun**: tahun pada **tanggal mulai kontrak** harus **sama** dengan tahun treaty
      induknya; yang berbeda **ditolak**. *(AC 14 spec; `[keputusan work owner]` — penyimpangan
      sadar 4)*
- [ ] ⚠️ Perbandingan memakai **tahun dari nilai bertipe tanggal**. Test yang menemukan **pemotongan
      teks posisi tetap** di jalur ini **gagal**. *(AC 14 spec)*
- [ ] Pesan penolakan menyebut **kedua tahun** — tahun pada tanggal mulai dan tahun treaty induk —
      sehingga pengguna tahu mana yang harus diperbaiki.
- [ ] Gerbang berlaku pada **penyimpanan kontrak**, dan tidak dapat dilewati lewat jalur API mana
      pun.
- [ ] Kontrak yang **tahunnya cocok** tersimpan tanpa keluhan — gerbang tidak menghasilkan positif
      palsu pada tanggal awal maupun akhir tahun (1 Januari dan 31 Desember diuji).
- [ ] Perubahan tanggal mulai pada kontrak yang sudah tersimpan **diuji ulang** oleh gerbang yang
      sama.

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Aturannya dipertahankan, caranya diperbaiki.** `[keputusan work owner]` Aturan "tahun tanggal
mulai = tahun treaty induk" **sah dan berguna** — ia mencegah kontrak nyasar tahun. Yang dibuang
hanyalah **cara Pega membandingkannya**.

`[terverifikasi]` Gerbang yang sama dipasang pada tiga activity simpan berbeda di Pega
(business, security, security reinsurer). Di sistem baru ia **satu aturan di satu tempat**, dipanggil
dari jalur simpan kontrak — bukan disalin tiga kali.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**, dengan **jam yang dapat dikendalikan** agar
kasus pergantian tahun dapat diuji.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Ralat bertanggal 30-09-2026 — sesi implementasi (paket 0)

> Sumber: `RALAT-DEV-30-09-2026.md` (K1–K8 katalog DEV, R1–R12 pembacaan ulang XML) dan `PARITAS-LAYAR-DAN-AKSI.md`. Kalimat di atas **tidak dihapus**; yang berlaku adalah ralat ini.

| Kalimat lama | Ralat |
| --- | --- |
| *"**Aturannya dipertahankan, caranya diperbaiki.**"* | ⛔ premis keliru: gerbang itu **MATI** di Pega (T1 ronde 2 — `PRE=false`) → **OQ-MCRL-01** (K7) |
| *"sistem menolak kontrak yang tanggal mulainya jatuh di tahun berbeda"* | bawaan sampai OQ-MCRL-01 dijawab: **ikut XML — tidak ditegakkan**. Pelajaran GILIRAN-11: dua gerbang mati pernah ditegakkan di modul lain dan harus dicabut |

**Status:** `[menunggu OQ-MCRL-01]` — tidak dibangun. Uji paket 3 membuktikan kontrak dengan tahun berbeda **tetap tersimpan**.

## Status 01-10-2026 (paket 11)

**Status:** ⏸️ tetap `[menunggu OQ-MCRL-01]` — gerbang tidak ditegakkan, di backend maupun layar (ikut Pega, K7).

## Status 01-10-2026 — OQ-MCRL-01 ditutup (`PROMPT-LANJUTAN-TIGA-MODUL-LIFE-KEPUTUSAN-OQ.md` §2)

Kalimat lama *"ralat bertanggal `[menunggu OQ-MCRL-01]`"* tidak berlaku lagi: **ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)** — gerbang tahun (tahun tanggal
mulai kontrak = tahun treaty) **ikut Pega, tidak ditegakkan** (langkahnya `PRE=false` di `SaveSecurityLife_Act` 4·5·6,
`SaveSecurityReinsurerLife_Act` 4·5·6, `SaveBusinessLife_Act` 6). Bawaan yang sudah dibangun (K7) dipertahankan; tiket ini final.
