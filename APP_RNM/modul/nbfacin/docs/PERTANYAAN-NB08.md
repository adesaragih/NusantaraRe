# Satu pertanyaan yang menahan enam predikat `When` — operand tanpa kutip

**Untuk:** work owner, lewat UI Pega · **Dari:** pengerjaan tiket `issues/08-registry-rules-eval.md` ·
**Tanggal:** 1 Oktober 2026

## Yang kami tanyakan

Di rule `When` berikut, nilai pembanding ditulis **tanpa tanda kutip**:

```
IsUW   (ASM-FW-GISFW-Work)  pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName = ReasFacInDirector
```

Saat Pega mengevaluasinya, apakah `ReasFacInDirector` dibaca sebagai **teks** "ReasFacInDirector",
atau sebagai **nama properti** yang nilainya diambil dari clipboard?

**Cara menjawab:** buka `IsUW` di UI Pega dan lihat kondisi A — apakah nilai sisi kanannya tampil
sebagai teks atau sebagai rujukan properti. Satu jawaban berlaku untuk keenam rule di bawah.

> 💡 **Rekomendasi agent (1 Oktober 2026)** — bukan jawaban; keputusannya tetap pada penerima.
>
> Baca operand tanpa kutip sebagai **teks**, setelah dikonfirmasi dengan satu kali lihat di UI Pega.
> Penguatnya: `ReasFacInDirector`, `ReasFacInGroupLeader`, … adalah nama antrean (workbasket) yang juga
> muncul sebagai nilai `PositionNote` di flow — nilai, bukan nama properti. Sampai dikonfirmasi, keenam
> predikat tetap ditolak keras.

> ✅ **Jawaban work owner (1 Oktober 2026)** — `KEPUTUSAN-30-09-2026.md` butir 28
> **Dibaca sebagai teks** ("operand isuw itu harusnya pakai kutip"). Diterapkan: lima predikat kini
> teks dan tidak lagi panic. **Satu pengecualian menunggu konfirmasi (A13):** `IsVisible` `= True` tetap
> ditolak — bila Pega membacanya sebagai boolean, teks `"True"` tidak sama dengan `"true"`.


## Mengapa penting

Bila dibaca sebagai properti dan properti itu tidak ada, nilainya kosong — dan operator yang **tidak
punya workbasket** akan membuat `"" = ""` bernilai benar: `IsUW` terbuka. Karena itu keenam predikat
kini **ditolak keras** di sistem baru sampai jawabannya ada (keputusan work owner 01-10-2026,
`KEPUTUSAN-30-09-2026.md` butir 24).

| Predikat | Baris tanpa kutip | Contoh |
| --- | ---: | --- |
| `IsLimitSBondKBG` | 26 | `= A1` … `= F6` |
| `IsLimitCreditCL` | 11 | idem |
| `IsUW` | 3 | `= ReasFacInDirector` · `ReasFacInGroupLeader` · `ReasFacInUnderwriting` |
| `IsVisible` | 1 | `= True` (huruf besar) |
| `IsCedingConfirmOffer` | 1 | `= Offer` |
| `IsCedingConfirmPolicy` | 1 | `= Policy` |

`[terverifikasi]` 43 baris, dihitung dari `registry_gen.go` dan dari korpus `NB FacIn\When` secara
terpisah; keduanya sepakat. Dua predikat (`IsLimitSBondKBG`, `IsLimitCreditCL`) dirujuk predikat lain,
sehingga perujuknya ikut ditolak saat dievaluasi.
