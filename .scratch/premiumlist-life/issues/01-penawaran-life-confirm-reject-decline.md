# 01: Penawaran Life — Confirm / Reject / Decline, dan percabangan Offer / Premium

**Status:** ready-for-agent

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**

## Hasil & nilai pengguna

Sebagai **inputor Life**, saya ingin mencatat penawaran dari ceding dan menyatakan keputusan
**Confirm**, **Reject**, atau **Decline** atasnya, supaya setiap penawaran punya jejak keputusan yang
jelas dan alur berikutnya (berhenti di penawaran, atau lanjut ke premium list) ditentukan secara
sadar — bukan oleh aturan tersembunyi. *(User story 1–8 di spec)*

## Area codebase

`internal/handlers` (endpoint buat/ubah penawaran + endpoint keputusan), `internal/services`
(transisi tahap; penentuan Offer/Premium), `internal/repository` (tulis rekam offer JSON; baca balik
id offer), `frontend/` (layar Input Offer, tombol keputusan, tampilan tahap berjalan).

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InputPolicyHolder` | `ASM-FW-GISFW-WORK-LIFE` / `INPUTPOLICYHOLDER` / `RULE-OBJ-FLOW` | `PremiumList Life/InputPolicyHolder.xml` | titik masuk flow (106.880 byte) |
| `InputOfferLife_ACT` | `ASM-FW-GISFW-WORK-LIFE` / `INPUTOFFERLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/InputOfferLife_ACT.xml` | simpan penawaran (131.691 byte, **8 langkah**) |
| `SaveOfferJsonLife_SQL` | `ASM-FW-GISFW-INT-OFFERJSON` / `ASM!SAVEOFFERJSONLIFE_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/SaveOfferJsonLife_SQL.xml` | `POOLDATA.INSERTJSONOFFERLIFE(...)` + `COMMIT;` |
| `GetIdOffer_SQL` | `ASM-FW-GISFW-INT-OFFERJSON` / `ASM!GETIDOFFER_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetIdOffer_SQL.xml` | baca balik id offer |
| `IsLifeAccepted` | `ASM-FW-GISFW-WORK-LIFE` / `ISLIFEACCEPTED` / `RULE-OBJ-DECISIONTABLE` | `PremiumList Life/DecisionTable/IsLifeAccepted.xml` | nilai keluaran `Confirm` / `Decline` / `Reject` |
| `IsFlagOnGoingPolicy` | `ASM-FW-GISFW-WORK-LIFE` / `ISFLAGONGOINGPOLICY` / `RULE-OBJ-DECISIONTABLE` | `PremiumList Life/DecisionTable/IsFlagOnGoingPolicy.xml` | nilai keluaran `Decline` / `Offer` / `Premium` |
| `setNoOffer_Act`, `setCeding_act`, `setPolicyHolder_act`, `setSOB_act`, `setSecurityReinsurer_act` | `ASM-FW-GISFW-WORK-LIFE` / `…` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/` | pengisian field penawaran |

`[terverifikasi]` **Peta 8 langkah `InputOfferLife_ACT`** (nomor dari `<pyStepPageReference>`, yang
muncul **sesudah** `<pyStepsActivityName>` langkahnya):

| Step | Langkah | Catatan |
| ---: | --- | --- |
| ~~1~~ | ~~`Call ASMForceCaseClose`~~ "Jika status Decline / Closed" | **REMARK** (`<pyStepsBlockName>//`, baris 354) |
| 2 | `Property-Set` | |
| 3 | `Property-Set` "Pega to json_offer_life" | rakit payload |
| **4** | `RDB-List` "Insert to table json_offer_life" | → `SaveOfferJsonLife_SQL` |
| 5 | `RDB-List` "Get ID from json_offer_life" | → `GetIdOffer_SQL` |
| 6, 7 | `Property-Set` | |
| 8 | `Obj-Save` | |

`[keputusan work owner]` Keputusan `Confirm`/`Reject`/`Decline` dibuat **manual** oleh inputor/admin.
Kedua DecisionTable mengekspor **nol baris keputusan** — yang direplikasi adalah **akibat**
keputusan, bukan formula yang memilihnya.

`[keputusan work owner]` `IsFlagOnGoingPolicy`: **`1` = Offer** (berhenti di tahap penawaran),
**`2` = Premium** (lanjut Input Premium List Detail).

## ADR terkait

**ADR-0007** (jejak audit setiap transisi), **ADR-0003** (uang non-float — nilai penawaran),
**ADR-0009** (migrasi penuh), **ADR-0001** (batas konteks — penawaran adalah hulu Claim Life).

## Acceptance criteria

- [ ] `Confirm` pada tahap penawaran melanjutkan case ke penentuan Offer/Premium. *(AC 1 spec)*
- [ ] `Reject` pada tahap mana pun **mengembalikan** case ke layar Input Offer — bukan menutupnya,
      bukan memajukannya. *(AC 2 spec)*
- [ ] `Decline` pada tahap mana pun **menutup** case; case tertutup tidak dapat dilanjutkan maupun
      diputuskan ulang. *(AC 3 spec)*
- [ ] Keluaran **Offer** menghentikan siklus di tahap penawaran, dan penawaran **tetap tersimpan**
      serta dapat dibaca kembali. *(AC 4 spec)*
- [ ] Keluaran **Premium** membuka tahap Input Premium List Detail. *(AC 5 spec)*
- [ ] Tidak ada aturan otomatis yang menetapkan keputusan; keputusan selalu datang dari tindakan
      pengguna. *(AC 6 spec)*
- [ ] Setiap transisi tahap menulis jejak audit: siapa, kapan, dari tahap apa ke tahap apa.
      (**ADR-0007**)
- [ ] Rekam penawaran ditulis lewat `INSERTJSONOFFERLIFE` (upsert berkunci `IDPEGA` + `STATUS`),
      lalu id offer dibaca kembali — keputusan berikutnya memakai id itu.
- [ ] ⚠️ Pemanggilan `INSERTJSONOFFERLIFE` **commit sendiri**. Kegagalan sesudahnya tidak boleh
      menghapus rekam offer; pemanggilan ulang **aman** karena upsert. *(lihat catatan di bawah)*
- [ ] Tidak ada nilai uang pada penawaran yang melewati `float`, termasuk di JSON API.

### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §12

- [ ] ⚠️ Penawaran tersimpan **relasional**; `JSON_OFFER_LIFE` **tidak ditulis** dan
      `INSERTJSONOFFERLIFE` **tidak dipanggil**. AC di atas tentang "commit sendiri" karena itu
      **tidak berlaku**. *(AC 32, 33, 34 spec; penyimpangan sadar 1)*
- [ ] ⚠️ Riwayat penawaran/konfirmasi ceding tersimpan di **`T_VIEW_SUGGEST`**, anak **langsung**
      header polis — tidak ada simpul `OfferFacIn` di skema baru. Kolomnya: nomor urut, tanggal, PIC,
      hasil (`Accept`/`Reject`/`Decline`), komentar, dan tahap (`Offer`/`Bind`).
      *(AC 44 spec; `[terverifikasi]` `PremiumList Life/Activity/AddHistorySuggest.xml`,
      `ASM-FW-GISFW-WORK-LIFE!ADDHISTORYSUGGEST`)*
      *(AC 14 spec; **ADR-0003**)*

## Blocker

**Tidak ada pemblokir.** Satu catatan terbuka yang **tidak** memblokir: **OQ-067** — lihat di bawah.

## Catatan — `INSERTJSONOFFERLIFE` ada di tahap ini, bukan di rantai simpan premium list

⚠️ `[terverifikasi]` Aturan urutan procedure `[keputusan desain]` menempatkan `INSERTJSONOFFERLIFE`
**paling akhir**, sesudah `INSERTJSONPOLISLIFE`. Di korpus, ia **tidak berada di rantai itu**: satu-
satunya perujuk `SaveOfferJsonLife_SQL` di kedua modul adalah `InputOfferLife_ACT` step 4 — **satu
tahap lebih awal**. Rantai simpan premium list (`InsertJsonPolisLife_Act`) merujuk
`GetJsonProductLife`, `InsertPLSummary`, `InsertJsonPolis`, `SaveLifeinProduction_SQL`,
`GetNopolisByIDPega` — bukan `SaveOfferJsonLife_SQL`.

Tiket ini **mengikuti korpus**: rekam offer ditulis di tahap penawaran. Aturan urutan diterapkan pada
rantai simpan premium list di tiket **05a**/**05b**. Konfirmasi dicatat sebagai **OQ-067**.

## Catatan — kode mati yang tidak dimigrasikan

`[terverifikasi]` `Call ASMForceCaseClose` (step 1) **REMARK**. Penutupan case pada `Decline`
dikerjakan mesin alur sistem baru, bukan dengan memanggil penutup paksa.

⚠️ **OQ-066**: penanda `<pyStepsBlockName>` **tidak dapat dipercaya sendirian** di modul Life —
ada langkah aktif yang sebenarnya mati. Sebelum memigrasikan langkah mana pun dari modul ini,
konfirmasikan ke work owner; jangan menyimpulkan hidup/mati dari penanda saja.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```
