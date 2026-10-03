# 06: Transactional outbox — keputusan + daftar efek dalam satu transaksi

**Status:** sebagian — AC 6 dibantah korpus (email langkah 11 diantre di setiap tingkat); menunggu work owner mencabut AC itu

**Blocked by:** 04b (rekam akseptasi — satu jalur simpan)

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin keputusan komite dan **antrean efek keluarnya** tersimpan dalam
**satu transaksi** — sehingga tidak ada efek yang hilang bila proses mati tepat setelah keputusan
disimpan. *(User story 24, 28 di spec)*

Ini fondasi jaminan "wajib berhasil"; pengirimannya sendiri ada di tiket 07.

## Area codebase

`internal/models` (entri outbox: efek apa, muatan apa, keadaan apa), `internal/repository`
(penulisan outbox dalam transaksi yang sama dengan keputusan), `internal/services` (penyusunan
daftar empat efek saat keputusan final).

Belum menyentuh `frontend/` — permukaan manusia ada di tiket 08.

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` step **6** | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` `Obj-Save` — **persist**; keempat efek berjalan **sesudahnya** |

`[terverifikasi]` Urutan efek yang harus masuk antrean, dibaca dari `<pyStepPageReference>` versi
2026-09-15:

| Urutan | Step | Efek |
| ---: | ---: | --- |
| 1 | **8** | `Call InsertJsonClaimLife_Act` |
| 2 | **10** | `Call serviceInsertArasapasClaimLife_act` |
| 3 | **11** | `Call SendEmailKlaimLife` |
| 4 | **12** | `Call HitServiceToKasirKMTLife_Act` — Kasir/pembayaran |

`[terverifikasi]` Step **10** dan **12** digerbangi `AcceptStatus = 1 && KomiteCount == KomiteLoop`
(baris 8648, 8887) — hanya pada keputusan Setuju di tingkat final.

## ADR terkait

**ADR-0015** (efek keluar wajib berhasil — transactional outbox, at-least-once; **menyimpang dari
ADR-0008**), **ADR-0013** (alamat endpoint di-lookup runtime, bukan disimpan di outbox sebagai URL),
**OQ-013** (batas transaksi dipegang Go).

## Acceptance criteria

- [x] Keputusan dan **daftar keempat efeknya** tersimpan dalam **satu transaksi database**.
      *(AC 19 spec)* — bukti: `services/komite_keputusan.go:KeputusanKomite.Putuskan` (`antreEfekKomite` di dalam `DalamTransaksi`), uji `TestUrutanKeputusanDalamSatuTransaksi`; tiga efek — InsertJson dibuang (2026-09-16)
- [x] Bila transaksi gagal, **tidak ada** keputusan tersimpan **dan tidak ada** entri outbox — tidak
      ada keadaan separuh. — bukti: uji `TestUrutanKeputusanDalamSatuTransaksi` (satu `DalamTransaksi`); `repository/efekkeluar.go:PohonKlaim.AntreEfek` menuntut `*Tx`
- [x] Bila proses mati tepat setelah commit, antrean efek **tetap ada** dan dapat diambil worker. — bukti: `repository/efekkeluar.go:PohonKlaim.AntreEfek` (baris `antre`, `JADWAL_BERIKUT` = saat) + `PungutEfek` bersaring `MODUL`; penjadwal pekerja belum ada (tiket 07)
- [x] Tiap entri outbox membawa **ID idempoten unik** sejak dibuat. *(AC 21 spec)* — bukti: `repository/efekkeluar.go:PohonKlaim.AntreEfek` (`ID` dari `SEQ_LOG_SERVICE_RNM`)
- [x] Entri outbox menyimpan **kunci kategori** endpoint, **bukan URL** — alamat di-resolve saat
      kirim (**ADR-0013**). — bukti: uji `TestMuatanOutboxKomiteTanpaURL`
- [ ] Efek hanya diantrekan pada keputusan yang benar-benar final; keputusan di tingkat bukan-
      terakhir tidak menghasilkan entri outbox. — belum: diralat — email (langkah 11) diantre pada setiap keputusan tingkat; hanya Arasapas/Kasir yang terbatas Setuju akhir (uji `TestEfekMenurutLangkah10Sampai12`)
- [x] Keputusan dilaporkan **tersimpan**, belum **tuntas** — kedua keadaan itu dibedakan.
      *(**ADR-0015**)* — bukti: `services/komite_keputusan.go:HasilKeputusanKomite` (`efekTertunda`) + kalimat layar `KasusKomite.tsx`, uji `TestKeadaanEfekKasus`

## Catatan — mengapa Komite berbeda dari Claim Life

⚠️ **ADR-0015 menyimpang dari ADR-0008.** Di Claim — Life, efek keluar boleh gagal tanpa memblokir.
Di sini tidak — pemicunya **Kasir**, integrasi **pembayaran** yang `[terverifikasi]` tidak ada di
Claim Life.

Kegagalan diam di Kasir berarti **klaim disetujui tetapi tidak pernah sampai ke pembayaran**, dan
tidak ada yang tahu.

`[keputusan work owner]` Step **14** `SetInformationData` **bukan** efek keluar — hanya temporary
penampung hasil submit. **Tidak** masuk outbox.

## Perintah verifikasi

```
go test ./internal/...
make check
```

## Implementasi — 28-09-2026 (giliran 10)

### Pembacaan ulang XML — langkah 8–12 `KomitePostAdjustment` (sesudah 7 `Obj-Save`)

| # | Activity | Precondition (`pyStepsPreCondition true`) | Nasib |
| ---: | --- | --- | --- |
| 8 | `InsertJsonClaimLife_Act` | — | ⛔ JSON — dibuang (2026-09-16) |
| 10 | `serviceInsertArasapasClaimLife_act` | `AcceptStatus = 1 && KomiteCount == KomiteLoop` (b8648), `IsPEGAPROD` | ✅ diantre pada Setuju akhir |
| 11 | `SendEmailKlaimLife` | **`IsPEGAPROD` saja** | ✅ diantre pada **setiap** keputusan |
| 12 | `HitServiceToKasirKMTLife_Act` | b8887 tingkat akhir, `Type=="TP"‖"TR"`, `IsKPR=="KPR"` | ✅ diantre pada Setuju akhir — lihat OQ |

Kunci Kasir `[terverifikasi]` `HitServiceToKasirKMTLife_Act`: `Kategori_1 "Kasir"`, `Kategori_2
"insertAllPaymentKasir"`; `Connect-REST`-nya sendiri `//` "kalau diserver dev jangan dijalanin".

### ⚠️ RALAT AC 6 tiket ini

*"Keputusan di tingkat bukan-terakhir tidak menghasilkan entri outbox"* — dibantah langkah 11: email
hanya bergerbang `IsPEGAPROD`, jadi ia berjalan pada **setiap** keputusan tingkat. Yang dibangun
mengikuti korpus; Arasapas dan Kasir tetap hanya pada Setuju akhir.

### ⚠️ `[terbuka — pemilik ekspor]` tiga `When` langkah 12

Transisi `When`-nya (lanjut/lewati) tidak terbaca dari ekspor. "Ketiganya AND" = Kasir hanya untuk
polis treaty ber-KPR — terlalu sempit untuk ditebak untuk integrasi pembayaran. Dibangun: gerbang Setuju
akhir; dua syarat lain dicatat di `komite_outbox.go`.

### Yang dibangun

- `services/komite_outbox.go` — `EfekKeputusanKomite` (murni, urutan 10 → 11 → 12), `antreEfekKomite`
  menulis ke outbox `T_LOG_SERVICE_RNM` (`MODUL = KOMITELIFE`) lewat `AntreEfek` yang **ada** (menuntut
  `*Tx` bukan nil), **di dalam transaksi `Putuskan`** — sesudah keputusan dan penyelesai akhir, sebelum
  jejak. Gagal di mana pun → tidak ada keputusan **dan** tidak ada entri.
- ID baris = `SEQ_LOG_SERVICE_RNM` — pengenal unik sejak lahir (AC 21), kunci anti-dobel tiket 07.
- Muatan: pengenal kasus/klaim/baris/akun, waktu, **kunci kategori** — nol URL, nol email
  (`TestMuatanOutboxKomiteTanpaURL` menagih himpunan medannya tepat).
- `IsPEGAPROD` tidak menggerbangi **pengantrean**, hanya pengiriman (tiket 07).
- Jawaban keputusan membawa `efekTertunda` — "tersimpan, belum tuntas" (ADR-0015); layar mengatakannya.

### AC — keadaan

| AC | Keadaan |
| --- | --- |
| keputusan + efeknya satu transaksi | ✅ |
| gagal → tidak ada keduanya | ✅ satu `DalamTransaksi` |
| proses mati sesudah commit → antrean tetap ada | ✅ baris outbox `antre` |
| ID idempoten unik | ✅ |
| kunci kategori, bukan URL | ✅ |
| efek hanya pada keputusan final | ⚠️ diralat — email setiap keputusan (langkah 11) |
| dilaporkan tersimpan, belum tuntas | ✅ `efekTertunda` |

### Angka

Go **583 PASS · 0 FAIL** tingkat atas; vet (+`-tags db`), gofmt bersih · vitest **357** · tsc bersih.

### ⛔ RALAT 28-09-2026 (temuan `/code-review` Spec) — gerbang Kasir langkah 12 TERBACA

Bagian "tiga `When` langkah 12" di atas keliru menyebutnya tak terjawab. Transisi tiap `When` terbaca: pola normal file ini `WhenTrue 2` (lanjut) / `WhenFalse 3` (lewati) — 48 lawan 32 kemunculan; baris `Type=="TP"||"TR"` langkah 12 **terbalik** (`WhenTrue 3`, `WhenFalse 2`). Kasir karena itu berjalan **hanya** bila Setuju di tingkat akhir **dan** `Type` bukan TP/TR **dan** `IsKPR == "KPR"` (dan `IsPEGAPROD`). Dibangun `KasirBerlaku`; `IS_KPR` klaim induk (`T_GENERAL_CLAIM`, 002) kini ikut dibaca kasus komite. Dikunci `TestKasirHanyaNonTreatyBerKPR` dan `TestGerbangKasirVERBATIMDariKorpus` (membaca korpus). Bacaan sebelumnya mengantre Kasir untuk polis treaty dan klaim non-KPR.

⚠️ Dicatat dari tinjauan Spec, belum dibangun (terang di PARITAS): lampiran (`AttachDocumentLife`/`SaveAttachLife`), unduh dokumen, `SetRemarkKomiteLife`/`Remarks`, syarat tampil `ShowTransfer`, `PrintAkseptasiPDF`. Suntingan aditif pada berkas Claim Life (`antrean.go`, `efekkeluar.go` ×2, penjaga uji) dilaporkan.
