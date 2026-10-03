// Panel ENAM total uang SATU PESERTA — A3 kelompok Detail & Tutup.
//
// Meniru `Section/ClaimLifeDetailGCNM.xml`.
//
// ⛔ RALAT LETAK 27-09-2026. Ronde pertama menaruh panel ini di tingkat
// KLAIM, berjudul "Total klaim". KELIRU: `ClaimLifeDetailGCNM.xml` berkelas
// `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL` (b84), dan keenam medan terikat
// properti berawalan TITIK pada halaman itu - `.TotalCedingRetention` b20636,
// `.TotalShareRNM` b20921, `.TotalSumInsured` b21208, `.TotalSumReasured`
// b21495, `.TotalShareRetro` b21782, `.TotalClaimAmount` b22069. Titik di
// depan berarti "properti halaman yang sedang berjalan", dan halaman itu
// PESERTA. Sebab salah bacanya sama dengan butir av: labelnya dibaca,
// IKATANNYA tidak.
//
// ⛔ RALAT KEDUA 27-09-2026, dan ini yang lebih besar. Berkas ini sempat
// berkata: "KELIMA TOTAL INI TIDAK DIHITUNG DI SINI, DAN ITU KEPUTUSAN" —
// alasannya `CheckTotalAdjustmentClaim` nol berkas rule-nya di korpus,
// sehingga kami mengaku tidak tahu baris mana yang ikut dijumlah. DUA hal
// keliru di situ:
//
//   1. Totalnya ENAM. `Total Ceding Retention` b20629 luput karena kami
//      mencacah lewat RUJUKAN `CheckTotalAdjustmentClaim`, dan ia
//      satu-satunya total tanpa aksi refresh - jadi ia tidak ikut tercacah.
//   2. Nilainya bukan misteri. Rujukan menggantung itu nyata dan tetap
//      dilaporkan (OQ-H), tetapi yang hilang hanya pemanggil REFRESH.
//      Yang MENGHITUNG keenam angka ada di korpus:
//
//        `Activity/SavePesertaClaim.xml`        8 b4002 / 8.1 b4221 / 8.2 b4592
//        `Activity/SaveOutStandingLife_Act.xml` 23 b10638 / 23.1 b10841 / 23.2 b11067
//
//      Keduanya berprasyarat KOSONG, sehingga jawabannya: SELURUH baris
//      `.AdjustmentList` peserta itu, termasuk yang `STS_REJECT = 2`.
//
// Sebab salahnya: berhenti pada rule yang NAMANYA tertulis di section, tanpa
// menanyakan siapa lagi yang menulis medan itu. Pelajaran yang sama dengan
// "Close Claim tidak punya aksi lain" - berhenti di aksi pertama yang
// membawa activity, lalu menyimpulkan tentang seluruh tombol.
//
// Penjumlahannya sendiri ada di backend (`models.HitungTotalPeserta`), bukan
// di sini: uang dijumlah sekali, di tempat yang punya apd.Decimal, dan layar
// hanya menampilkan teks yang sudah jadi (ADR-U-0003, ADR-U-0016).

import { DETAIL } from '../labels'
import { jumlahUang, type TotalPeserta } from '../api'

/** Satu total: label VERBATIM dan nilainya sebagai TEKS. */
export interface MedanTotal {
  label: string
  /** Jumlah sebagai teks desimal apa adanya. Tidak pernah number. */
  jumlah: string
  mataUang: string
}

/**
 * Nama rule yang hilang dari ekspor.
 *
 * ⚠️ Ia TETAP disebut walau angkanya kini ada. Rujukan menggantungnya nyata -
 * sepuluh pemanggilan, nol berkas - dan pertanyaan ke pemilik ekspor tinggal
 * satu: kenapa rule refresh-nya tidak ikut diekspor. Menghapus nama ini
 * membuat OQ-H tidak dapat ditelusuri dari kode.
 */
export const RULE_TOTAL_HILANG = 'CheckTotalAdjustmentClaim'

/**
 * Menyusun keenam total dari jawaban backend.
 *
 * ⚠️ Dipisah dari komponennya supaya dapat diuji tanpa DOM, dan supaya
 * cacahnya dapat dikunci: total yang DIHILANGKAN dari layar sama merusaknya
 * dengan total yang dikarang, dan yang pertama tidak berbunyi sendiri.
 *
 * Urutannya urutan layar: b20629, b20914, b21201, b21488, b21775, b22063.
 */
export function medanTotal(total: TotalPeserta): MedanTotal[] {
  return [
    { label: DETAIL.totalCedingRetention, ...isi(total.cedingRetention) },
    { label: DETAIL.totalShareNusantaraRe, ...isi(total.shareNusantaraRe) },
    { label: DETAIL.totalSumInsured, ...isi(total.sumInsured) },
    { label: DETAIL.totalSumReasured, ...isi(total.sumReasured) },
    { label: DETAIL.totalShareRetro, ...isi(total.shareRetro) },
    { label: DETAIL.totalClaimAmount, ...isi(total.jumlahKlaim) },
  ]
}

/** Membaca satu nilai uang lewat `jumlahUang`, yang menolak angka JSON. */
function isi(uang: TotalPeserta['jumlahKlaim'] | undefined): {
  jumlah: string
  mataUang: string
} {
  return { jumlah: jumlahUang(uang), mataUang: uang?.currency ?? '' }
}

export function PanelTotalPeserta({ total }: { total: TotalPeserta }) {
  return (
    <section className="polis">
      <h3 className="polis__judul">Total peserta</h3>
      <dl className="polis__daftar">
        {medanTotal(total).map((m) => (
          <div key={m.label} className="polis__baris">
            <dt className="polis__label">{m.label}</dt>
            <dd className="polis__nilai">
              {m.jumlah === '' ? (
                <span className="polis__kosong">—</span>
              ) : (
                <>
                  {m.jumlah} {m.mataUang}
                </>
              )}
            </dd>
          </div>
        ))}
      </dl>
      <p className="polis__catatan" role="note">
        Keenam total adalah jumlah <strong>seluruh</strong> baris adjustment
        peserta ini — termasuk baris yang ditolak, persis seperti{' '}
        <code>SavePesertaClaim</code> langkah 8.1 yang berprasyarat kosong.
        Angkanya dihitung saat halaman dibaca dan tidak disimpan.
      </p>
    </section>
  )
}
