# 05: Reject Outstanding oleh Admin — tanpa Komite

**Status:** ready-for-agent

**Blocked by:** 04 (mesin status per baris)

## Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, saya dapat menolak baris adjustment yang saya input sendiri tanpa melalui
Komite — sehingga kesalahan input dapat saya batalkan cepat — dan penolakan itu **hanya membatalkan
baris itu**, klaimnya tetap hidup dan saya dapat menginput baris pengganti. *(User story 6 dan 7 di
spec)*

Ini jalur penolakan **kedua** di sistem ini. Jalur pertama (lewat Komite) datang di tiket 11.

## Area codebase

`internal/services` (aturan penolakan + gerbang status baris), `internal/handlers` (endpoint reject),
`frontend/` (kontrol "Reject Outstanding" pada baris adjustment).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/RejectOSClaimLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `REJECTOSCLAIMLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` menulis `2` ke **dua tingkat berbarengan**: `.STS_REJECT` dan `…PremiumListDetail(idx).STS_REJECT` |
| `Claim Life/Section/RejectOSClaimLife_Sec.xml` | layar pemicu — `[terverifikasi]` merujuk activity di atas 2× sebagai `<pyActivity>` | |
| `Claim Life/Section/AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | `[terverifikasi]` gerbang tombol "Reject Outstanding": `pyWorkPage.pyPosition =='ReasLifeAdmin' && pyWorkPage.ClaimData.PremiumListSummary.CLAIM_NO !='' && .STS_REJECT == 0` |

Gerbang itu menguji **`.STS_REJECT` tingkat baris** — bukti bahwa wewenang pun diukur per baris.

## ADR terkait

**ADR-0011** (dua sumber penolakan dengan makna setara pada tingkat baris; `RejectOSClaimLife_Act`
**bukan** dead rule), **ADR-0002** (peran).

## Acceptance criteria

- [ ] Penolakan oleh `ReasLifeAdmin` membuat **hanya baris itu** berstatus Ditolak; baris lain pada
      klaim yang sama tidak berubah. *(AC 3 spec)*
- [ ] Klaim **tidak** tertutup oleh penolakan itu, dan baris adjustment baru dapat diinput
      sesudahnya.
- [ ] Penolakan hanya mungkin pada baris yang **masih Outstanding** — baris yang sudah Aksep atau
      Ditolak menolak upaya itu.
- [ ] Penolakan hanya mungkin bila klaim sudah punya nomor (padanan `CLAIM_NO !=''`).
- [ ] Nilai status hasil penolakan Admin **sama** dengan hasil penolakan Komite; tidak ada nilai
      khusus yang membedakan keduanya.
- [ ] Pencerminan ke tingkat `PremiumListDetail` terjadi bersamaan, bukan menyusul. *(AC 7 spec)*
- [ ] ⚠️ **Diselaraskan 2026-09-16:** penolakan langsung oleh `ReasLifeAdmin` menulis `STS_REJECT`
      = **`2`** sebagai **nilai sebenarnya menurut aksi** — bukan nilai yang di-hardcode seperti di
      Pega. Test yang menemukan nilai di-hardcode **gagal**. *(AC 42 spec; penyimpangan sadar 4)*

## Catatan

`[keputusan work owner 2026-09-14]` Bahwa penolakan Admin membatalkan **baris saja** — bukan klaim —
berasal dari work owner. `[terverifikasi]` Korpus **tidak membedakan** kedua sumber penolakan: kedua
rule menulis `2` dalam bentuk yang sama persis. Perbedaan itu memang tidak seharusnya ada.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
