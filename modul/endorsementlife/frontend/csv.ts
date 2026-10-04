// CSV unggahan di peramban - MURNI. Dipakai popup `View Upload` (`ViewCSVResult_LifeEDM`) dan unduhan
// `Generate Data Detail` (`GenerateDataDtlLife_act`). Nilainya teks mentah, seperti daftar sementara
// Pega `TempWorkPage.ListLifePremiumDetailUpload`.
//
// ⛔ BUKAN validasi. Yang menentukan baris lolos/ditolak hanya server (`POST /unggah`, `POST /csv`) -
// pengurai ini sekadar memperlihatkan isi berkas.

/**
 * Kepala unduhan `Generate Data Detail` - `Activity/GenerateDataDtlLife_act.xml` b276
 * `CSVPropHeaders`, VERBATIM dan berurutan (termasuk `STATUS` dan `OVRR_COMM` yang tidak dipetakan
 * `SaveCSVEDMLife` 4.1 - server mengabaikan dan melaporkannya).
 */
export const KEPALA_UNDUH_CSV: readonly string[] = [
  'POLICY_NO',
  'POLICY_HOLDER',
  'CERTIFICATE_NO',
  'NAME_OF_INSURED',
  'SEX',
  'DOB',
  'AGE',
  'PLAN',
  'BEGIN_DATE',
  'EFFECTIVE_DATE',
  'EXPIRED_DATE',
  'PERIOD_MM',
  'STATUS',
  'EM_PERCENT',
  'SUM_INSURED',
  'CEDING_RETENTION',
  'SUM_REASURED',
  'SHARE_NUSANTARA_RE',
  'GROSS_PREMIUM',
  'COMM',
  'OVRR_COMM',
  'BROKERAGE_FEE',
  'NET_PREMIUM',
]

/** `FileName` b270 `DetailUpload` + `AppendTimeStampToFileName` b272 `true`. */
export function namaBerkasUnduh(kini: Date): string {
  const d = (n: number) => String(n).padStart(2, '0')
  const tgl = `${kini.getFullYear()}${d(kini.getMonth() + 1)}${d(kini.getDate())}`
  return `DetailUpload_${tgl}_${d(kini.getHours())}${d(kini.getMinutes())}${d(kini.getSeconds())}.csv`
}

/** RFC 4180: kutip ganda, koma/baris baru di dalam kutip, CRLF atau LF; BOM di depan dibuang. */
export function uraiCSV(teks: string): string[][] {
  const s = teks.charCodeAt(0) === 0xfeff ? teks.slice(1) : teks
  const hasil: string[][] = []
  let baris: string[] = []
  let medan = ''
  let kutip = false
  let ada = false
  for (let i = 0; i < s.length; i++) {
    const c = s[i]
    if (kutip) {
      if (c === '"' && s[i + 1] === '"') {
        medan += '"'
        i++
      } else if (c === '"') {
        kutip = false
      } else {
        medan += c
      }
      continue
    }
    if (c === '"') {
      kutip = true
      ada = true
    } else if (c === ',') {
      baris.push(medan)
      medan = ''
      ada = true
    } else if (c === '\n' || c === '\r') {
      if (c === '\r' && s[i + 1] === '\n') i++
      baris.push(medan)
      hasil.push(baris)
      baris = []
      medan = ''
      ada = false
    } else {
      medan += c
      ada = true
    }
  }
  if (ada || baris.length > 0) {
    baris.push(medan)
    hasil.push(baris)
  }
  return hasil
}

/** Baris data berkunci judul huruf besar (seperti server); baris yang seluruhnya kosong dilewati. */
export function barisUnggahan(teks: string): Array<Record<string, string>> {
  const [judul, ...isi] = uraiCSV(teks)
  if (judul === undefined) return []
  const kunci = judul.map((j) => j.trim().toUpperCase())
  return isi
    .filter((b) => b.some((v) => v.trim() !== ''))
    .map((b) => Object.fromEntries(kunci.map((k, i) => [k, b[i] ?? ''])))
}

function kutipCSV(v: string): string {
  return /[",\r\n]/.test(v) ? `"${v.replace(/"/g, '""')}"` : v
}

/** CSV berkepala `kepala`, CRLF, nilai yang tak ada = kosong. */
export function rakitCSV(kepala: readonly string[], baris: ReadonlyArray<Record<string, string | undefined>>): string {
  const garis = [kepala.join(','), ...baris.map((b) => kepala.map((k) => kutipCSV(b[k] ?? '')).join(','))]
  return `${garis.join('\r\n')}\r\n`
}
