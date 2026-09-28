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

## Implementasi — 28-09-2026 (giliran 10)

### Pembacaan ulang XML — langkah 5 "Reject" (gerbang b8119 `AcceptStatus==2 && KomiteCount == KomiteLoop`)

| Sub | Isi | Nasib |
| --- | --- | --- |
| 5.1 b5785 | "Set Reject Komite berjenjang": perulangan **seluruh** `…AdjustmentList(…).KomiteList` — `KomiteAproval = 2`, `KomiteComment = ""`, `DateApprove = @CurrentDateTime()`, dan salinannya di `pyWorkPage.KomiteList`; **tanpa syarat** | ⛔ **tidak ditiru** — OQ-K-05 |
| 5.3 b6279 | "Set Nilai Akseptasi": `AdjustmentList.STS_REJECT = 2`, `PremiumListDetail.STS_REJECT = 2`, `PremiumListDetail.IsCheck = "false"` | ✅ |
| 5.4 b6450 | "Set Param" — bahan `UpdateOsAkseptasiClaimLife_sql` | ✅ lewat jalur tunggal 04b |
| 5.5 | "Tukar SecurityReinsurer dengan RetroName" `//` | ➖ mati |
| 5.6 b7911 | "Insert ke OS" | ✅ `RekamAkhirWarisan(…, "2", "")` |

⚠️ Penajaman: precondition `.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"` pada 4.15 dan
5.6 **tertulis tetapi mati** — `pyStepsPreCondition false` (pecahan baris 20/21 di jendela langkahnya).
Tanpa itu keduanya akan membatalkan dirinya sendiri, sebab 4.11/5.3 baru saja mengubah `STS_REJECT`.

### ⛔ OQ-K-05 `[terbuka — work owner]` — 5.1 menghapus riwayat tangga

Bila tingkat akhir menolak, 5.1 menimpa keputusan `1` dan komentar **setiap** tingkat sebelumnya dengan
`2` dan teks kosong. Menirunya menghapus satu-satunya catatan siapa yang menyetujui sebelum penolakan —
bertentangan dengan ADR-0007 dan dengan riwayat tangga tiket 09. Yang ditulis sistem ini: keputusan
tingkat akhir itu saja (tiket 02); tingkat lain tetap apa yang mereka putuskan. Bila paritas persis
dikehendaki, itu satu `UPDATE` bersyarat — keputusannya milik work owner.

### ⚠️ OQ-K-05b `[terbuka — work owner]` — Tolak di tingkat TENGAH

Langkah 5 bergerbang tingkat akhir, dan `IsKomiteLoop` palsu sesudah Tolak — jadi Tolak di tingkat
tengah **menghentikan tangga tanpa menyentuh baris klaim**: baris tetap `Outstanding`, dan
`KOMITE_ID`-nya tetap menunjuk kasus komite yang sudah berhenti. Di sistem ini akibatnya konkret: jalur
penyerahan Claim Life menolak baris yang sudah ber-`KOMITE_ID` (`ErrBarisSudahDiserahkan`), sehingga
baris itu tidak dapat diserahkan ulang **dan** tidak pernah ditolak. Perilaku korpusnya ditiru apa
adanya; jalan keluarnya (tolak tengah = tolak baris? lepaskan `KOMITE_ID`?) keputusan work owner.

### Yang dibangun

- `penyelesaiAkhirOracle.Tolak` (menggantikan penolak bawaan yang gagal terang): gerbang anggota
  berjalan; **dua tingkat baris** bernilai sama dalam satu operasi lewat `PerbaruiStatusBaris`
  (`Outstanding → Ditolak`, penjaga `STS_REJECT = 0` = hanya baris yang masih Outstanding); nomor dan
  tanggal akseptasi **kosong**; `CabutPenandaDipilih` (`IS_CHECK = "false"`, peserta dapat dipilih ulang
  — AC 2); `CerminkanHeader`; rekam akhir jalur tunggal status `2`; jejak — seluruhnya di transaksi
  keputusan akhir.
- `AcceptStatus` dipetakan di batas ini ke `KodeDitolak`; ia tidak disimpan sebagai status kedua di
  Claim Life (ADR-0001).

### AC — keadaan

| AC | Keadaan |
| --- | --- |
| Setuju akhir → Aksep | ✅ (04a) |
| Tolak akhir → Ditolak; klaim tetap dapat menerima baris baru | ✅ `IS_CHECK = "false"`, baris baru lewat jalur Claim Life |
| dua tingkat baris, nilai sama, satu operasi | ✅ |
| hanya di tingkat akhir | ✅ — dan OQ-K-05b untuk tingkat tengah |
| hanya baris Outstanding, dipilih, belum bernomor | ✅ penjaga `STS_REJECT = 0` + gerbang A2 saat penyerahan |
| pelaku dan waktu terekam | ✅ jejak |
| `AcceptStatus` tidak diteruskan | ✅ |

### Angka

Go **580 PASS · 0 FAIL** tingkat atas; vet (+`-tags db`), gofmt bersih · vitest **357** · tsc bersih.
