# 04: Mesin status per baris + aturan turunan "klaim selesai"

**Status:** ready-for-agent

**Blocked by:** 03 (baris `AdjustmentList` + Save ke Outstanding)

## Hasil & nilai pengguna

Sebagai **pengguna mana pun**, saya melihat status tiap baris adjustment secara terpisah dan tahu
persis apa yang sudah diputuskan dan apa yang belum — dan sebuah klaim yang seluruh barisnya sudah
diputus tampak "selesai" sehingga antrean kerja saya bersih. *(User story 25 dan 28 di spec)*

Ini **inti spesifikasi**: tiket yang menetapkan bahwa yang diputuskan adalah baris, bukan klaim.

## Area codebase

`internal/models` (status baris sebagai nilai tertutup), `internal/services` (transisi + aturan
turunan), `internal/handlers` (endpoint transisi), `frontend/` (tampilan status klaim dan baris).

Status klaim **dihitung**, bukan disimpan sebagai kolom mandiri — meskipun sistem lama
menyimpannya sebagai kolom.

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/SaveOutStandingLife_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY` | penulis `0` |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` penulis `1` (6 `Property-Set`) dan `2` (2 `Property-Set`) |
| `Claim Life/Activity/RejectOSClaimLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `REJECTOSCLAIMLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` penulis `2` dari sisi Admin |
| `Claim Life/Activity/SetSTS_Reject.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SETSTS_REJECT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` penurunan status klaim → baris (`.STS_REJECT = Primary.STS_REJECT`) |

`[terverifikasi]` **Sensus penulis lengkap** ada di register **OQ-061**: tidak ada satu pun rule di
`Claim Life` maupun `Komite Claim Life` yang menulis `0` setelah `1` atau `2`.

## ADR terkait

**ADR-0011** (unit = baris; kefinalan; arti tunggal nilai `2`), **ADR-0001** (jalur balik Komite
bekerja pada tingkat baris).

## Acceptance criteria

- [ ] Baris yang sudah bernilai Aksep atau Ditolak **tidak dapat berubah lagi** melalui jalur mana
      pun. *(AC 2 spec)*
- [ ] Menolak sebuah baris **tidak** menutup klaim; klaim tetap dapat menerima baris baru.
      *(AC 4 spec)*
- [ ] Status "Ditolak" pada sebuah baris **selalu** berarti baris itu ditolak — **tidak pernah**
      berarti klaim selesai, apa pun sumber penolakannya.
- [ ] ⚠️ **Diselaraskan 2026-09-16:** pencerminan terjadi pada kolom `STS_REJECT` /
      `ACCEPTED_NO` di **`T_CLAIMLF_PREMIUMLIST_DETAIL`** dan **`T_GENERAL_CLAIM`** (spec §2b) — unit
      keputusannya tetap baris `T_CLAIMLF_ADJUSTMENT` (**ADR-0011**), yang kini menggantung pada
      **peserta**. *(AC 33 spec; penyimpangan sadar 2)*
- [ ] `PremiumListDetail` dan header klaim **selalu mencerminkan** baris adjustment terakhir, dan
      **tidak** ditulis sebagai status mandiri. *(AC 7 spec)*
- [ ] Klaim dilaporkan "selesai" **hanya** bila tidak ada baris berstatus Outstanding **dan** ada
      sekurangnya satu baris berstatus Aksep. *(AC 8 spec)*
- [ ] Bila tidak ada baris Outstanding dan tidak ada pula yang Aksep, klaim berada dalam keadaan
      **ditolak seluruhnya** — dan tetap dapat dilanjutkan dengan baris baru.
- [ ] Riwayat lengkap seluruh baris pada satu klaim dapat dilihat, sehingga putaran Komite terbaca.
      *(User story 27 spec)*

## Catatan penutupan (2026-09-14)

**Aturan turunan "klaim selesai" disetujui work owner** `[keputusan work owner]` dan tercatat di
`CONTEXT.md` (Lampiran, butir 1):

- Status klaim adalah **turunan**, **bukan** kolom tersimpan.
- Unit keputusan = **baris `AdjustmentList`** (**ADR-0011**).
- **Header klaim mengikuti adjustment terakhir.**

`[terverifikasi 2026-09-14]` **Di Pega**, `AdjustmentList` **tidak punya tabel fisik** (kelas
`ASM-FW-GISFW-Data-AdjustmentLife`, berawalan `Data-` = embedded), dan barisnya di-persist ke
`POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` yang kolomnya identik. Bukti:
`Claim Life/Activity/SaveOutStandingLife_Act.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` /
`SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`) — `pyStepsObjectName = .AdjustmentList`,
`pyStepsClassName = ASM-FW-GISFW-Data-AdjustmentLife`.

⚠️ **Diselaraskan 2026-09-16 — itu keadaan Pega, bukan keadaan sistem baru.** `[keputusan work
owner]` Di sistem baru baris adjustment punya **tabelnya sendiri**, `T_CLAIMLF_ADJUSTMENT`, yang
menggantung pada **peserta** (`T_CLAIMLF_PREMIUMLIST_DETAIL`) — bukan pada header klaim, dan **bukan**
pada `OS_AKSEPTASI_KLAIM_LIFE` (spec §2b, AC 33).

⚠️ **Koreksi 2026-09-16** `[keputusan work owner]`: `OS_AKSEPTASI_KLAIM_LIFE` **tetap di-`INSERT`
flat**, berdampingan dengan tabel relasional — hilir masih membaca dari sana. Yang **dibuang hanya
JSON**. Jadi baris adjustment tersimpan di **`T_CLAIMLF_ADJUSTMENT`** *dan* ikut terbawa ke rekam flat
itu; yang berubah adalah **di mana unit keputusan tinggal**, bukan hilangnya tabel akseptasi.
*(AC 32 spec)*

**Ini tidak mengubah ADR-0011** — unit keputusan tetap **baris**. Yang berubah hanya di mana baris
itu tinggal: dari tabel akseptasi bersama menjadi tabel adjungan per peserta.

`[data DBA]` `STS_REJECT` pada tabel **warisan** bertipe `NUMBER(38)`. ⚠️ Di skema baru tipenya
ditetapkan sendiri di tiket **14**; nilainya tetap `0`/`1`/`2` dan diisi **menurut aksi**, bukan
di-hardcode (AC 42 spec).

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
