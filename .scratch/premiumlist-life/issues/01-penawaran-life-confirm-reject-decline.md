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


## Pembacaan ulang XML — 28 September 2026 (peta konektor utuh)

`InputPolicyHolder.xml` dibaca sebagai **pohon**: tiap blok disusuri dari `<pyMOId>` ke
`<pyMOId>` berikutnya, lalu `pyTo` dan `pyFrom` di dalamnya.

⚠️ **`pyTo` MENDAHULUI `pyFrom` di DOM.** Membaca berpasangan dari atas tanpa menyadari itu
menghasilkan graf yang **seluruh panahnya terbalik** — dan graf terbalik tetap terlihat masuk akal.

### Dua belas konektor, seluruh `pyExpression` VERBATIM

```
Start1        --Transition1 -------------------> Assignment2   (Input Offer Life b1340)
Assignment2   --Transition3  [InputDataOfferLife b1735] -> Decision1
Decision1     --Transition4  [Confirm b1574] ---> Decision3
Decision1     --Transition5  [Decline b2235] ---> End1         (Resolved-Rejected b848)
Decision3     --Transition10 [Premium b1658] ---> ASSIGNMENT63 (Input Premium Detail b1069)
Decision3     --Transition11 [Offer   b1807] ---> END52        (Resolved-Completed b947)
ASSIGNMENT63  --TRANSITION54 [ShowLifePremiumDetail b1505] -> Decision2
Decision2     --Transition7  [Confirm b2090] ---> Utility1     (InsertJsonPolisLife b765)
Decision2     --Transition6  [Decline b2162] ---> End1         (Resolved-Rejected)
Decision2     --Transition9  [Reject  b2306] ---> Assignment2  (kembali ke Input Offer Life)
Utility1      --Transition8 -------------------> END52         (Resolved-Completed)
Assignment1   --Transition2  [ShowLifePremiumSummary b1881] -> END52
```

### ⛔ RALAT ATAS AC 2 TIKET INI — 28 September 2026

AC 2 berbunyi: *"`Reject` pada tahap **mana pun** mengembalikan case ke layar Input Offer"*.

**Konektornya membantah.** `Reject` muncul **tepat sekali** di seluruh flow — `Transition9`
b2306, pada `Decision2`, yaitu penggolong **sesudah Input Premium Detail**. `Decision1` —
penggolong sesudah tahap penawaran — hanya punya `Confirm` b1574 dan `Decline` b2235.

Jadi **menolak di tahap penawaran tidak punya jalur di sistem lama**. Orang yang berada di layar
Input Offer dan tidak ingin melanjutkan memakai **`Decline`**. Menyediakan `Reject` di sana
berarti membuat jalur yang tidak pernah ada, dan kasus yang menempuhnya akan mendarat di tempat yang
tidak dikenal sistem hilir.

Dikunci `TestRejectDiTahapPenawaranTidakPunyaJalur`, dan cacahnya dikunci
`TestNamaTahapDanStatusVERBATIMDariKorpus` yang **membaca berkas korpus langsung** dan menagih
`pyExpression Reject` muncul tepat satu kali.

### ⚠️ `[terbuka — work owner]` `Assignment1` nol konektor masuk

Kedua belas konektor disusuri satu per satu; **nol** di antaranya ber-`pyTo` `Assignment1`
*(Input Premium Summary b1215/b1232)*. Yang **keluar** ada — `Transition2` b1866. Artinya tahap
itu dicapai lewat **ticket** *(`Ticket1` b958/b1456)* atau lewat jalur yang tidak ikut diekspor.

Bentuk yang sama pernah ditemukan di Claim Life: `Assignment4`, yang nol `pyPosition`-nya.
Perpindahan **ke** tahap ini karena itu **tidak disediakan** — ia akan menjadi jalur karangan.

### ⭐ Dua hal yang tiket ini sebut benar, dan kini VERBATIM

- **AC 3** — `Decline` menutup kasus. Statusnya **`Resolved-Rejected`** b848, bukan status
  tersendiri. ⚠️ Nama statusnya menyebut *"Rejected"* sedangkan keputusan yang menuju ke sana
  bernama *"Decline"* — kosakata sistem lama, ditiru apa adanya.
- **AC 4** — keluaran `Offer` **menutup** kasus dengan `Resolved-Completed` b947
  *(`Transition11` b1807)*, bukan menunggu. Penawarannya tetap tersimpan; yang selesai
  **pekerjaannya**.

## Implementasi — tiket 01 bagian 1 (28 September 2026)

| Sisi | Isi |
| --- | --- |
| Model | `models/polis_penawaran.go` — peta konektor sebagai komentar, `TransisiPenawaran`, `LanjutanPenggolong`, `KasusPolisTertutup`; nama tahap dan status **VERBATIM** |
| Repository | `repository/polis_work.go` — `Keadaan`, `PindahTahap` *(syarat `POSITION` lama)*, `TutupKasus` *(syarat `STATUS IS NULL`)* |
| Services | `services/polis_penawaran.go` — `Putuskan`, `Golongkan`, satu `pagari`, satu transaksi bersama jejak |
| Handlers | `POST /api/polis-life/{{id}}/keputusan` dan `…/penggolong` |

⛔ **`Confirm` di tahap penawaran TIDAK menulis apa pun.** Di flow ia hanya memindahkan kendali
ke `Decision3`; yang memindahkan kasus adalah hasil penggolong itu. Menuliskan tahap lebih awal
berarti kasus berpindah **sebelum ada yang memutuskan ke mana**.

⛔ **AC 6 dijaga uji statik.** `TestNolAturanOtomatisMenetapkanKeputusan` membaca berkas
layanannya dan menolak literal `"Confirm"`/`"Reject"`/`"Decline"` di dalamnya: keputusan
selalu **datang dari pengguna**, tidak pernah dihitung.

### Kode bersama yang disentuh — aditif, dan satu penyempitan

- `handlers/handlers.go` — **dua baris rute ditambahkan**; nol baris Claim Life disunting.
- `services/tutupkontrak_test.go` — **dipersempit**, bukan dilonggarkan. Ronde pertamanya
  menuntut `PastikanKasusTerbuka` dari **setiap** layanan pengubah; kalimat itu benar selama
  hanya ada satu modul. `polis_penawaran.go` **memang** memeriksa kasus tertutup — lewat
  `T_WORK_POLIS` dan `models.KasusPolisTertutup`, sebab layanan polis yang menanyai
  `T_WORK_CLAIM` akan selalu menjawab *"tidak ada"*. Penjaga kini menerima penjaga **yang
  disebut namanya per berkas**; yang tidak terdaftar tetap dituntut penjaga bawaan, dan **yang nol
  penjaga tetap gagal** — dibuktikan merah dengan mencabut pemeriksaannya.

### Yang BELUM ada di tiket ini

Layar **Input Offer** beserta ketiga tombolnya, dan kotak masuk `PremiumList` beserta
`Input Premium` b16472 / `Input Offer` b16964. Backend-nya siap dan berute; layarnya
bagian 2.
