# 09: Riwayat tangga persetujuan di UI

**Status:** ready-for-agent

**Blocked by:** **00 (skema penyimpanan komite — PREFACTOR)**, 03 (penegakan wewenang + eskalasi) — riwayat harus memuat eskalasi juga

## Hasil & nilai pengguna

Sebagai **pengguna mana pun**, saya ingin melihat **riwayat lengkap tangga** pada satu kasus — siapa
memutuskan apa, kapan, dengan komentar apa — sehingga saya paham mengapa kasus berada di keadaannya
sekarang tanpa bertanya kepada orang.

Sebagai **auditor**, saya ingin eskalasi ikut terbaca di riwayat yang sama, sehingga pengetatan
wewenang tidak dilubangi diam-diam. *(User story 9, 10, 29, 30, 31 di spec)*

## Area codebase

`internal/services` (penyusunan riwayat dari entri per tingkat + catatan eskalasi),
`internal/handlers` (endpoint riwayat kasus), `frontend/` (tampilan riwayat pada layar kasus).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, 515.675 byte | `[terverifikasi]` menulis entri per tingkat |

`[terverifikasi]` Tiga `Property-Set` mengisi elemen list yang diindeks pencacah tangga, dengan
`Local.Komite = pyWorkPage.KomiteCount`:

| Properti | Nilai |
| --- | --- |
| `…AdjustmentList(IndexAdjustment).KomiteList(Local.Komite).KomiteAproval` | `pyWorkPage.AcceptStatus` |
| `…KomiteList(Local.Komite).KomiteComment` | `pyWorkPage.Comment` |
| `…KomiteList(Local.Komite).DateApprove` | `@CurrentDateTime()` |

Hal yang sama juga ditulis ke `pyWorkPage.KomiteList(Local.Komite)`.

**Jadi setiap tingkat tangga menghasilkan satu entri berisi keputusan, komentar, dan waktu** —
riwayatnya sudah ada di data; tiket ini membuatnya terbaca.

## ADR terkait

**ADR-0007** (jejak audit setiap transisi — termasuk eskalasi), **ADR-0014** (eskalasi adalah
tindakan yang direkam: siapa, kapan, dari tingkat mana ke tingkat mana), **ADR-0011** (riwayat
melekat pada **baris `AdjustmentList`**, bukan pada klaim).

## Acceptance criteria

- [ ] Riwayat menampilkan **tiap tingkat** tangga secara berurutan, dengan keputusan, komentar, dan
      waktu. *(AC 7 spec)*
- [ ] Keputusan ditampilkan sebagai **kata** (Setuju / Tolak), bukan angka dan bukan nama field.
      *(AC 29 spec)*
- [ ] **Eskalasi ikut terbaca** di riwayat yang sama: siapa memindahkan, kapan, dari tingkat mana ke
      tingkat mana. *(AC 12 spec)*
- [ ] Tingkat yang **dilewati** karena eskalasi terlihat sebagai dilewati — bukan hilang tanpa jejak.
- [ ] Riwayat melekat pada **baris `AdjustmentList`** yang diputuskan; satu klaim dengan beberapa
      baris menampilkan riwayat per baris.
- [ ] Riwayat dapat dibaca **tanpa** wewenang memutuskan — melihat bukan memutuskan.

### Sumber riwayat ⚠️ BARU 2026-09-16 — spec §9

⚠️ **Koreksi premis.** Catatan lama menyebut *"riwayatnya sudah ada di data"* — yang dimaksud adalah
**page runtime Pega**. Di sistem baru page itu **dibuang**; riwayat punya tabelnya sendiri.

- [ ] ⚠️ Riwayat tangga dibaca dari **`T_KOMITE_KOMITELIST` diurut `KOMITE_URUT`** — **bukan** dari page
      runtime maupun JSON. Test yang menemukan pembacaan dari page **gagal**. *(AC 33 spec;
      penyimpangan sadar 1)*
- [ ] Tiap baris riwayat menampilkan **anggota pemutus**, **keputusannya**, **komentarnya**, dan
      **tanggal putus** — langsung dari kolom `KOMITE_ID`, `KOMITE_APROVAL`, `KOMITE_COMMENT`,
      `DATE_APPROVE`. *(AC 31 spec)*
- [ ] Tingkat yang **belum memutus** terbaca sebagai `KOMITE_APROVAL = 0`, dan ditampilkan sebagai
      kata — bukan angka. *(AC 29, 31 spec)*

## Catatan

`[keputusan work owner]` Eskalasi menaikkan `KomiteCount` tanpa tingkat yang dilewati memberi
keputusan, sehingga entri `KomiteAproval` untuk tingkat itu **kosong**. Bentuk penampilannya —
"dilewati (eskalasi)" atau serupa — adalah keputusan implementasi; yang mengikat adalah **ia tidak
boleh tampak seperti tingkat yang belum diputuskan**.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Implementasi — 28-09-2026 (giliran 10)

### Pembacaan ulang XML

`Section/ShowTransfer.xml` grid `.KomiteList` (b29727): judul `Committee` → `.IDKomite` (= jabatan,
`CreateKMTLife_Act` b952/b1041), `Status` → `.KomiteAproval` (dropdown), `Date Approve` →
`.DateApprove`, `Comment` → `.KomiteComment`. Dipakai VERBATIM sebagai judul kolom riwayat.

### Yang dibangun

- `GET /api/komite/{id}/riwayat` — siapa pun yang teridentifikasi (melihat ≠ memutuskan): tangga dari
  `T_KOMITE_KOMITELIST` urut `KOMITE_URUT` (bukan page runtime/JSON), status **kata**
  (`Setuju`/`Tolak`/`Menunggu`/`Dilewati (eskalasi)`), anggota (`KOMITE_OPERATORID`, pengenal akun),
  tanggal putus, komentar; **eskalasi** dari jejak `T_CLAIMLF_JEJAK` (dari → ke, oleh, kapan),
  berkunci `ADJUSTMENT_ID` + kasus (ADR-0011).
- Teks jejak tingkat/eskalasi kini **satu tempat** (`awalanJejakTingkat`, `awalanJejakEskalasi`),
  dipakai penulis (tiket 02/03) dan pembaca — dikunci `TestPenulisDanPembacaJejakSatuBentuk`.
- `KataApprovalKomite` (tiket 01) kini menerjemahkan `1`/`2` dan membedakan dilewati dari menunggu.
- Layar Kasus Komite menampilkan tangga dari riwayat dengan judul VERBATIM, dan daftar eskalasi.
- `PARITAS-LAYAR-DAN-AKSI.md` milik modul ini lahir (tiket 01–09).
- ⚠️ Penjaga batas konteks Claim Life (`DateApprove` peka huruf) — medan Go dinamai `TanggalPutus`,
  kunci JSON tetap `dateApprove`; tanpa pengecualian baru.

### AC — keadaan

| AC | Keadaan |
| --- | --- |
| tiap tingkat berurutan: keputusan, komentar, waktu | ✅ |
| keputusan sebagai kata | ✅ |
| eskalasi terbaca di riwayat yang sama | ✅ |
| tingkat dilewati terlihat dilewati | ✅ `Dilewati (eskalasi)` ≠ `Menunggu` |
| melekat pada baris adjustment | ✅ satu kasus per `ADJUSTMENT_ID` (unik, 013) |
| dapat dibaca tanpa wewenang memutuskan | ✅ `TestRiwayatUntukSiapaPun` |
| dari `T_KOMITE_KOMITELIST` urut `KOMITE_URUT` | ✅ |
| anggota, keputusan, komentar, tanggal langsung dari kolomnya | ✅ |
| belum memutus = `0` sebagai kata | ✅ `Menunggu` |

### Angka

Go **596 PASS · 0 FAIL** tingkat atas; vet (+`-tags db`), gofmt bersih · vitest **358** · tsc bersih.
