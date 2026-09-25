# 07: Penegakan peran di lapisan layanan + wewenang kirim-Komite per `Type`

**Status:** ready-for-agent

**Blocked by:** 05 (reject Outstanding oleh Admin) — gerbang perlu tindakan nyata untuk dijaga

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin setiap tindakan pada klaim hanya dapat dilakukan peran yang
berhak — dan penegakannya berada di **lapisan layanan**, bukan sekadar di visibilitas layar —
sehingga wewenang tidak dapat dilewati dengan memanggil API langsung. *(User story 15 di spec)*

## Area codebase

`internal/services` (penegakan wewenang sebelum setiap transisi), `internal/handlers` (identitas
pemanggil), `frontend/` (kontrol yang tampil/tersembunyi mengikuti peran — **sebagai kenyamanan,
bukan sebagai penegakan**).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Section/AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | tiga gerbang `<pyCondition>` di bawah |
| `Claim Life/Flow/Register_Flow.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` / `RULE-OBJ-FLOW` | `[terverifikasi]` 15 kemunculan `pyPosition` — penegakan sesungguhnya di sistem lama bertumpu pada penugasan tahap di sini |

`[terverifikasi]` Tiga gerbang pada `AdjustmentDetail_Section`:

| Kontrol | `<pyCondition>` |
| --- | --- |
| "Reject Outstanding" | `pyWorkPage.pyPosition =='ReasLifeAdmin' && …CLAIM_NO !='' && .STS_REJECT == 0` |
| Jalur Komite (`GetListKomiteLife`) | `pyWorkPage.pyPosition =='ReasLifeSPV' \|\| pyWorkPage.Type = 'TP' \|\| pyWorkPage.Type = 'TR'` |
| "Save to Outstanding" | `pyWorkPage.pyPosition =='ReasLifeSPV'` |

`[terverifikasi]` `Type` dibaca dari **dua salinan** di Pega — `pyWorkPage.Type` (gerbang wewenang)
dan `pyWorkPage.PolicyDataLife.Type` (validasi DOL) — disalin di
`Claim Life/Activity/LoadDataPeserta_Act.xml` (`ASM-FW-GCNMFW-WORK` / `LOADDATAPESERTA_ACT` /
`RULE-OBJ-ACTIVITY`): `pyWorkPage.Type = pyWorkPage.PolicyDataLife.Type`. **Sumber otoritatif =
`PolicyDataLife.Type`.**

## ADR terkait

**ADR-0002** (tiga peran; rangkap peran tidak diperbolehkan), **ADR-0012** (wewenang kirim-Komite
bergantung `Type`; **dibawa apa adanya sebagai paritas**, risiko RBAC diterima dan ditinjau saat
konteks Komite/IAM digarap).

## Acceptance criteria

- [ ] `ReasLifeMedicalAdvisor` **tidak dapat** mengubah status akseptasi baris mana pun. *(AC 9 spec)*
- [ ] Untuk klaim ber-`Type` `QP` atau `QR`, **hanya `ReasLifeSPV`** yang dapat mengirim ke Komite;
      upaya oleh peran lain ditolak. *(AC 10 spec)*
- [ ] Untuk klaim ber-`Type` `TP` atau `TR`, `ReasLifeAdmin` **dapat** mengirim ke Komite.
      *(AC 11 spec)*
- [ ] Penolakan wewenang terjadi **di lapisan layanan**, dan tetap terjadi meskipun kontrol UI-nya
      ditampilkan. *(AC 12 spec)*
- [ ] Wewenang dan validasi membaca **satu** field `Type` yang sama — tidak ada dua salinan yang
      dapat berbeda.
- [ ] Tidak ada nama orang ter-hardcode di lapisan mana pun; wewenang diukur dari peran akun.

## Catatan

⚠️ `[terverifikasi]` Penegakan di sistem lama **tidak seragam**: `pyPosition` tidak ada sama sekali
di `Section/InputOSClaimLife.xml`, `Section/RejectOSClaimLife_Sec.xml`, `Section/ClaimComite.xml`,
maupun `Harness/Committe_Life.xml`. **Ketidakseragaman itu tidak ditiru** — sistem baru menegakkan
di lapisan layanan untuk semua tindakan.

⚠️ Kelonggaran `TP`/`TR` **adalah risiko yang diterima sadar** (**ADR-0012**): wewenang bergantung
pada atribut data, bukan peran. Jangan "memperbaikinya" dalam tiket ini — itu keputusan terpisah.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
