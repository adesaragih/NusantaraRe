// Rumus kepala form panel New — dua DataTransform yang dijalankan saat
// medannya BERUBAH (bukan tombol). Murni di layar; nol panggilan backend.
//
// ⛔ Modul Treaty In punya `TreatyInSetTreatyYear` sendiri; modul ini tidak
// boleh mengimpornya, jadi aturannya disalin dari ekspor korpus ini — yang
// diadu sama persis dengan milik Treaty In (7 Oktober 2026).

import { keInputTanggal } from '../../../../inti/frontend/lib/tanggalInput'

type Medan = Readonly<Record<string, string>>

/** Tanggal tersimpan / masukan → `Date` UTC tengah malam; `null` bila tak terbaca. */
function keTanggal(v: string | undefined): Date | null {
  const iso = keInputTanggal(v ?? '')
  const m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(iso)
  if (m === null) return null
  return new Date(Date.UTC(Number(m[1]), Number(m[2]) - 1, Number(m[3])))
}

/** `Date` → bentuk TERSIMPAN Pega (`YYYYMMDD`). */
function keTersimpan(d: Date): string {
  const p = (n: number) => String(n).padStart(2, '0')
  return `${String(d.getUTCFullYear())}${p(d.getUTCMonth() + 1)}${p(d.getUTCDate())}`
}

/**
 * `DataTransform/TreatyInSetTreatyYear` — saat Commencement berubah:
 *
 *   TreatyIn.TreatyYear  = @substring(TreatyIn.Commencement,0,4)
 *   TreatyIn.Termination = @addCalendar(TreatyIn.Commencement,"1",…)   +1 tahun
 *
 * ⚠️ `@addCalendar` menambah TAHUN kalender (Java `Calendar.add`): 29 Feb
 * menjadi 28 Feb tahun berikutnya, bukan 1 Maret.
 */
export function tahunTreaty(m: Medan): Record<string, string> {
  const d = keTanggal(m.Commencement)
  if (d === null) return {}
  const th = d.getUTCFullYear() + 1
  const akhirBulan = new Date(Date.UTC(th, d.getUTCMonth() + 1, 0)).getUTCDate()
  const akhir = new Date(Date.UTC(th, d.getUTCMonth(), Math.min(d.getUTCDate(), akhirBulan)))
  return { TreatyYear: String(d.getUTCFullYear()), Termination: keTersimpan(akhir) }
}

const HARI_MS = 86_400_000

/** `@DateTimeDifference(a, b, 'D')` — hari kalender dari a ke b. */
const selisihHari = (a: Date, b: Date) => Math.round((b.getTime() - a.getTime()) / HARI_MS)

/** `@divide(a, b)` — dua desimal, pembulatan setengah ke atas. */
function bagiDuaDesimal(a: number, b: number): number {
  if (b === 0) return 0
  const x = a / b
  return Math.sign(x) * Math.round(Math.abs(x) * 100 + 1e-9) / 100
}

/**
 * `DataTransform/TreatyCalculateProratePct` — saat Effective Date berubah dan
 * saat Is Pro Rate diklik:
 *
 *   ProRateDays = ProRateTotalDays = ProRatePercent = 0
 *   bila IsProRate == true:
 *     ProRateDays      = @DateTimeDifference(EDMEffective, Termination, 'D')
 *     ProRateTotalDays = @DateTimeDifference(Commencement, Termination, 'D')
 *     ProRatePercent   = @divide(ProRateDays, ProRateTotalDays) * 100
 *
 * ⭐ DIUKUR 7 Oktober 2026 atas 22 penyesuaian ber-IsProRate (seluruh
 * korpus): rumus ini cocok dengan 5 penyesuaian TERBARU (ID ≥ 1001016) —
 * hari dan persennya (`@divide` dua desimal: 134/364 → 0,37 → 37,00). Lima
 * belas yang LEBIH LAMA (ID ≤ 1000709) menyimpan hari Commencement →
 * Effective: versi rumus sebelumnya. Yang dibangun bunyi ekspor hari ini.
 */
export function proRata(m: Medan): Record<string, string> {
  const nol = { ProRateDays: '0', ProRateTotalDays: '0', ProRatePercent: '0' }
  if (m.IsProRate !== 'true') return nol
  const ef = keTanggal(m.EDMEffective)
  const mulai = keTanggal(m.Commencement)
  const akhir = keTanggal(m.Termination)
  if (ef === null || mulai === null || akhir === null) return nol
  const hari = selisihHari(ef, akhir)
  const total = selisihHari(mulai, akhir)
  return {
    ProRateDays: String(hari),
    ProRateTotalDays: String(total),
    ProRatePercent: (bagiDuaDesimal(hari, total) * 100).toFixed(2),
  }
}

/**
 * Medan yang ikut berubah bila `kunci` diubah — event `change`/`click`
 * sel kepala `Section/TreatyInNONProportional.xml`:
 *
 *   Commencement (sel 11)  → TreatyInSetTreatyYear
 *   EDMEffective (sel 30)  → TreatyCalculateProratePct
 *   IsProRate    (sel 31)  → TreatyCalculateProratePct (+ `TreatyInNonAddItem(share)`,
 *                            rumus Share — belum tersambung di layar ini)
 *
 * ⛔ Termination TIDAK punya event di ekspor: mengubahnya tidak menghitung
 * ulang Pro Rate, sama seperti Pega.
 */
export function ikutBerubah(kunci: string, m: Medan): Record<string, string> {
  switch (kunci) {
    case 'Commencement':
      return tahunTreaty(m)
    case 'EDMEffective':
    case 'IsProRate':
      return proRata(m)
    default:
      return {}
  }
}
