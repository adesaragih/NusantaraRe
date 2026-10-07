// Daftar kotak masuk Beranda modul NB Treaty In (keputusan work owner 06-10-2026: klik workbasket / jenis di Beranda
// menampilkan berkas yang menunggu akun, tanpa masuk menu NB Treaty In; kolom = portal + kolom yang ada datanya;
// ID membuka layar kasus langsung lewat `PropsRute.bukaKasus`). Nilai DB tampil apa adanya; tanggal dd-mm-yyyy.

import type { DaftarBeranda } from '../../../inti/frontend/modul'
import { daftarMenunggu, type RingkasanKasus } from './api'
import { KOLOM_BERANDA } from './labels'
import { sajikan } from './sajian'

/** Tanggal pertukaran `YYYY-MM-DD[ HH:MI:SS]` sebagai waktu lokal; selain itu null. */
function bacaWaktu(t: string): Date | null {
  const m = /^(\d{4})-(\d{2})-(\d{2})(?:[ T](\d{2}):(\d{2})(?::(\d{2}))?)?/.exec(t.trim())
  if (m === null) return null
  return new Date(Number(m[1]), Number(m[2]) - 1, Number(m[3]), Number(m[4] ?? 0), Number(m[5] ?? 0), Number(m[6] ?? 0))
}

/** Lama sejak `tgl` sampai `kini`, gaya kolom "Time Since Last Update": just now, 45m, 5h, 25d, 11mo, 1y 4mo ago. */
export function sejak(tgl: string, kini: Date): string {
  const w = bacaWaktu(tgl)
  if (w === null) return ''
  const detik = Math.max(0, Math.floor((kini.getTime() - w.getTime()) / 1000))
  if (detik < 60) return 'just now'
  if (detik < 3600) return `${Math.floor(detik / 60)}m ago`
  if (detik < 86400) return `${Math.floor(detik / 3600)}h ago`
  let bulan = (kini.getFullYear() - w.getFullYear()) * 12 + (kini.getMonth() - w.getMonth())
  if (kini.getDate() < w.getDate()) bulan--
  if (bulan < 1) return `${Math.floor(detik / 86400)}d ago`
  if (bulan < 12) return `${bulan}mo ago`
  const tahun = Math.floor(bulan / 12)
  const sisa = bulan % 12
  return sisa === 0 ? `${tahun}y ago` : `${tahun}y ${sisa}mo ago`
}

/** Baris daftar portal -> daftar Beranda (kolom `KOLOM_BERANDA`, urutan tetap). */
export function keDaftarBeranda(baris: readonly RingkasanKasus[], kini: Date): DaftarBeranda {
  const kolom = (Object.keys(KOLOM_BERANDA) as (keyof typeof KOLOM_BERANDA)[]).map((kunci) => ({
    kunci,
    label: KOLOM_BERANDA[kunci],
  }))
  return {
    kolom,
    baris: baris.map((b) => ({
      id: b.id,
      sel: {
        id: b.id,
        jenis: b.proportionalType,
        tertanggung: b.insuredName,
        bisnis: b.businessName,
        ceding: b.cedingCoName,
        mulai: sajikan(b.startDate, 'tanggal'),
        marketing: b.marketingName,
        status: b.nbStatus,
        pembuat: b.createOpName,
        tanggal: sajikan(b.tglCreate, 'tanggal'),
        sejak: sejak(b.tglUpdate, kini),
      },
    })),
  }
}

/** Penyedia `MenuModul.daftarBeranda` modul ini. */
export async function daftarBeranda(workbasket: string | null): Promise<DaftarBeranda> {
  return keDaftarBeranda(await daftarMenunggu(workbasket), new Date())
}
