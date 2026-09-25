# 05: Jalur balik — tulis `STS_REJECT` ke dua tingkat baris

**Status:** ready-for-agent

**Blocked by:** 04b (rekam akseptasi — satu jalur simpan)

## Hasil & nilai pengguna

Sebagai **ReasLifeSPV di Claim — Life**, saya ingin hasil keputusan komite **memantul ke baris
`AdjustmentList` saya**, sehingga status di sistem saya selalu mencerminkan keputusan terakhir — dan
bila ditolak, saya dapat mengajukan ulang dengan baris baru. *(User story 18 di spec)*

Ini **kontrak keluar** konteks Komite. Yang **membaca dan menampilkannya** adalah Claim — Life
(**CL-11**); yang **menulisnya** adalah tiket ini.

## Area codebase

`internal/services` (penerapan hasil keputusan ke dua tingkat baris), `internal/repository`
(penulisan status baris), `internal/models` (bentuk baris `AdjustmentList` dan `PremiumListDetail`).

Tidak menyentuh `frontend/` konteks ini — tampilannya milik Claim — Life.

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, 515.675 byte | `[terverifikasi]` menulis `STS_REJECT` ke **dua tingkat baris** — `…PremiumListDetail(idx).STS_REJECT` dan `…AdjustmentList(idx).STS_REJECT` |

`[terverifikasi]` Sensus penulis: **6** `Property-Set` bernilai `1` dan **2** bernilai `2`.
Gerbang baris yang boleh diubah: `.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"`.
Gerbang tingkat final: baris **5695** (aksep) dan **8119** (tolak).

`[keputusan work owner]` Pemetaan hasil: **Setuju → `STS_REJECT = 1`**, **Tolak → `STS_REJECT = 2`**.
⚠️ Nama field menyesatkan — nilai `1` berarti **diaksep**.

## ADR terkait

**ADR-0011** (unit keputusan = baris `AdjustmentList`; `PremiumListDetail` dan header klaim adalah
**cerminan** baris terakhir), **ADR-0001** (Kontrak 2 — jalur balik; `AcceptStatus` adalah kosakata
konteks ini dan **dipetakan di batas**, tidak disimpan sebagai status kedua di Claim Life),
**ADR-0007** (penerapan hasil merekam siapa + kapan).

## Acceptance criteria

- [ ] Hasil **Setuju** di tingkat terakhir membuat baris yang diserahkan berstatus **Aksep**.
- [ ] Hasil **Tolak** di tingkat terakhir membuat baris berstatus **Ditolak**, dan klaim **tetap**
      dapat menerima baris baru di sisi Claim — Life.
- [ ] `STS_REJECT` tertulis pada **dua tingkat baris** — `PremiumListDetail` dan `AdjustmentList` —
      dengan nilai **yang sama**, dalam satu operasi. *(AC 15 spec)*
- [ ] Perubahan status **hanya** terjadi ketika putaran mencapai **tingkat terakhir**
      (`KomiteCount == KomiteLoop`). *(AC 13 spec)*
- [ ] Hanya baris yang **masih Outstanding**, sudah dipilih, dan **belum bernomor akseptasi** yang
      dapat diubah oleh hasil keputusan.
- [ ] Setiap penerapan hasil merekam **pelaku dan waktu**. *(**ADR-0007**)*
- [ ] `AcceptStatus` **tidak** diteruskan ke Claim — Life sebagai status kedua; ia dipetakan ke
      `STS_REJECT` **di batas ini**. *(**ADR-0001**)*

## Catatan — kepemilikan lintas konteks

`[terverifikasi]` **Penulisan ini milik Komite, bukan Claim Life** — pembuktiannya ada di korpus:
`KomitePostAdjustment.xml` berada di modul `Komite Claim Life` dengan class
`ASM-FW-GCNMFW-Work-KomiteLife`.

Konsekuensinya untuk **CL-11**: tiket itu **membaca dan menampilkan** hasil yang ditulis di sini, dan
**tidak menulis `STS_REJECT` sendiri**. Cakupan CL-11 sudah diselaraskan.

## Perintah verifikasi

```
go test ./internal/...
make check
```
