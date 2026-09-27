// Panel lima total uang SATU PESERTA — A3 kelompok Detail & Tutup.
//
// Meniru `Section/ClaimLifeDetailGCNM.xml`.
//
// ⛔ RALAT LETAK 27-09-2026. Ronde pertama menaruh panel ini di tingkat
// KLAIM, berjudul "Total klaim". KELIRU: `ClaimLifeDetailGCNM.xml` berkelas
// `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL` (b84), dan kelima medan terikat
// properti berawalan TITIK pada halaman itu - `.TotalShareRNM` b20921,
// `.TotalSumInsured` b21208, `.TotalSumReasured` b21495,
// `.TotalShareRetro` b21782, `.TotalClaimAmount`. Titik di depan berarti
// "properti halaman yang sedang berjalan", dan halaman itu PESERTA.
//
// Jadi kelimanya total per PESERTA, bukan per klaim. Sebab salah bacanya
// sama dengan butir av: labelnya dibaca, IKATANNYA tidak.
//
// ⛔ KELIMA TOTAL INI TIDAK DIHITUNG DI SINI, DAN ITU KEPUTUSAN.
//
// Di Pega, kelima medan itu me-refresh lewat `CheckTotalAdjustmentClaim` —
// sepuluh pemanggilan, dua untuk tiap medan:
//
//   b20914 `Total Share Nusantara Re`  -> b20970, b21091
//   b21201 `Total Sum Insured`         -> b21262, b21377
//   b21488 `Total Sum Reasured`        -> b21546, b21664
//   b21775 `Total Share Retro`         -> b21836, b21951
//   b22063 `Total Claim Amount`        -> b22120, b22238
//
// ⚠️ RALAT 27-09-2026: kedua nomor per medan itu BUKAN dua aksi berbeda.
// Blok `postValue` (mis. b20940) tidak membawa `pyActivity` sama sekali -
// hanya `pyActivityClass`. Keduanya adalah aksi `refresh` yang SAMA,
// terserialisasi dua kali oleh Pega: sekali di `pyModes`, sekali di
// `pyActionSets`. Jadi lima medan, lima pemanggilan, sepuluh kemunculan
// teks. Ronde pertama menulis "postValue + refresh" dan itu keliru.
//
// Activity itu **tidak punya satu pun berkas rule di seluruh korpus** —
// diperiksa: nol hasil untuk `*CheckTotalAdjustment*`, dan satu-satunya
// berkas yang menyebutnya adalah section ini sendiri. Rujukan menggantung.
//
// Karena itu kita TIDAK TAHU apa yang ia hitung. "Total" terdengar seperti
// penjumlahan kolomnya, tetapi pertanyaan yang menentukan tidak terjawab:
// baris yang mana? Seluruhnya, atau hanya yang `IsCheck`? Termasuk yang
// `STS_REJECT = 2`? Jalur tolak mencabut `IsCheck`, jadi jawabannya
// berpengaruh — dan ini ANGKA UANG.
//
// Menjumlahkannya sendiri berarti menebak, dan angka uang yang ditebak jauh
// lebih berbahaya daripada angka yang dinyatakan belum ada (ADR-U-0003).
// Dilaporkan OQ-H.

import { DETAIL } from '../assets/labels'
import { BelumTersedia } from './ui/dasar'

/** Satu total: label VERBATIM dan keadaan sumbernya. */
export interface MedanTotal {
  label: string
  /** Selalu true hari ini — lihat kepala berkas. */
  belumBersumber: boolean
}

/** Nama rule yang hilang, disebut apa adanya supaya dapat dicari. */
export const RULE_TOTAL_HILANG = 'CheckTotalAdjustmentClaim'

/**
 * Menyusun kelima total.
 *
 * ⚠️ Dipisah dari komponennya supaya dapat diuji tanpa DOM, dan supaya
 * cacahnya dapat dikunci: total yang DIHILANGKAN dari layar sama merusaknya
 * dengan total yang dikarang, dan yang pertama tidak berbunyi.
 */
export function medanTotal(): MedanTotal[] {
  return [
    { label: DETAIL.totalShareNusantaraRe, belumBersumber: true },
    { label: DETAIL.totalSumInsured, belumBersumber: true },
    { label: DETAIL.totalSumReasured, belumBersumber: true },
    { label: DETAIL.totalShareRetro, belumBersumber: true },
    { label: DETAIL.totalClaimAmount, belumBersumber: true },
  ]
}

export function PanelTotalPeserta() {
  return (
    <section className="polis">
      <h3 className="polis__judul">Total peserta</h3>
      <dl className="polis__daftar">
        {medanTotal().map((m) => (
          <div key={m.label} className="polis__baris">
            <dt className="polis__label">{m.label}</dt>
            <dd className="polis__nilai">
              <BelumTersedia apa={m.label} />
            </dd>
          </div>
        ))}
      </dl>
      <p className="polis__catatan" role="note">
        Kelima total ini di sistem lama dihitung oleh rule{' '}
        <code>{RULE_TOTAL_HILANG}</code>, dan rule itu tidak ikut dalam ekspor
        yang kami terima — dirujuk sepuluh kali, nol berkasnya. Kami tidak
        menjumlahkannya sendiri: yang belum terjawab adalah baris mana yang
        ikut dihitung, dan salah menjawabnya mengubah angka uang.
      </p>
    </section>
  )
}
