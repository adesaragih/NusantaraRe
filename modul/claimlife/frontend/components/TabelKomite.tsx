// Tabel komite di bawah inbox Claim Life (perintah work owner 09-10-2026: menu Komite Claim Life dihapus, "anggap menu
// itu tidak pernah ada" - komite digabung ke menu Claim Life). Isinya Inbox Komite modul Komite Claim Life yang dipindah
// ke sini (disalin, bukan diimpor): kasus yang menunggu PELAKU (worklist `KomiteRouter`, penyaringnya di server) lewat
// rute pinjaman `GET /api/komite`, plus laporan harian "perlu intervensi" untuk admin (tiket 08). Klik baris = layar
// kasus komite DI TEMPAT (`onBukaModul`); Kembali ke inbox ini. `[tidak ada di korpus]`: bagian ini tampil hanya bila
// ada kasus yang menunggu, ada laporan admin, atau pembacaannya gagal - pelaku tanpa kasus komite tidak melihatnya.
//
// ⭐ 12-10-2026 (work owner "rapikan tampilan inbox committee ini"): kulit tabel Inbox Claim Life, paginasi `Halaman` inti,
// tanggal dd/mm/yyyy HH:mm, Claim Amount berformat ribuan rata kanan, Level rata tengah, Row Status sebagai chip. Mode
// `utama` (akun tanpa antrian Claim Life): judul halaman "Committee Inbox"; tanpa kasus → `kosong` milik pemanggil.
//
// ⛔ UANG TETAP TEKS (formatNumber atas teks desimal); status baris KATA. Data dan aksi tidak berubah.

import { useEffect, useState, type ReactNode } from 'react'

import { PERAN, type KodePeran } from '../../../../inti/frontend/labels'
import { Gagal, Halaman } from '../../../../inti/frontend/components/ui/dasar'
import { DESIMAL_TAK_DIBATASI, formatNumber } from '../../../../inti/frontend/lib/format'
import {
  daftarKomite,
  laporanHarianKomite,
  MODUL_KOMITE,
  type BarisKomite,
  type HalamanKomite,
  type LaporanHarianKomite,
} from '../api'
import { TABEL_KOMITE } from '../labels'

/** Sel kosong ditandai (ADR-U-0027). */
export function selKomite(nilai: string): string {
  return nilai.trim() === '' ? '—' : nilai
}

/** Tingkat berjalan dari seluruh tingkat — `KomiteCount`/`KomiteLoop`. */
export function tingkatKomite(b: BarisKomite): string {
  return `${String(b.tingkatBerjalan)} / ${String(b.komiteLoop)}`
}

/** "YYYY-MM-DD HH:MM[:SS]" (atau RFC3339) → dd/mm/yyyy HH:mm; tak terbaca / kosong = "—". */
export function waktuKomite(nilai: string): string {
  const m = /^(\d{4})-(\d{2})-(\d{2})[T ](\d{2}):(\d{2})/.exec(nilai.trim())
  return m === null ? selKomite('') : `${m[3]}/${m[2]}/${m[1]} ${m[4]}:${m[5]}`
}

/** Claim Amount berformat ribuan (teks desimal, tanpa float). */
export function nilaiKomite(nilai: string): string {
  return nilai.trim() === '' ? '—' : formatNumber(nilai, DESIMAL_TAK_DIBATASI)
}

/**
 * Work Status kasus komite. `T_WORK_CLAIM.STATUS_WORK` kasus komite yang masih terbuka SELALU kosong (diisi hanya saat
 * ditutup), dan daftar ini hanya memuat kasus terbuka yang tingkat PELAKU-nya masih menunggu (`KOMITE_APPROVAL` menunggu,
 * repository komite_inbox) — maka kosong berarti "menunggu keputusan Anda". Nilai yang ada ditampilkan apa adanya.
 */
export function statusWorkKomite(nilai: string): string {
  return nilai.trim() === '' ? TABEL_KOMITE.menungguKeputusan : nilai
}

export default function TabelKomite({
  peran,
  onBukaModul,
  utama = false,
  kosong,
}: {
  /** Peran pelaku — laporan "perlu intervensi" hanya untuk admin (tiket 08). */
  peran: readonly KodePeran[]
  /** `PropsRute.onBukaModul` - baris tabel: layar kasus komite dibuka di tempat. */
  onBukaModul?: (modul: string, id: string) => boolean
  /** Mode halaman sendiri (akun tanpa antrian Claim Life, 12-10-2026): judul halaman "Committee Inbox". */
  utama?: boolean
  /** Ditampilkan bila tidak ada kasus komite, laporan, maupun galat (mode utama). */
  kosong?: ReactNode
}) {
  const [hal, setHal] = useState<HalamanKomite | null>(null)
  const [halaman, setHalaman] = useState(1)
  const [laporan, setLaporan] = useState<LaporanHarianKomite | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [selesai, setSelesai] = useState(false)
  const admin = peran.includes(PERAN.admin)

  useEffect(() => {
    let aktif = true
    void (async () => {
      try {
        const h = await daftarKomite(halaman)
        // ⛔ Laporan kosong DINYATAKAN — ketiadaan laporan ≠ ketiadaan masalah.
        const l = admin ? await laporanHarianKomite() : null
        if (!aktif) return
        setHal(h)
        setLaporan(l)
      } catch (e) {
        // ⛔ Galat DINYATAKAN, bukan menjadi daftar kosong.
        if (aktif) setGalat(e)
      } finally {
        if (aktif) setSelesai(true)
      }
    })()
    return () => {
      aktif = false
    }
  }, [admin, halaman])

  const buka = (id: string) => {
    setPesan(null)
    if (onBukaModul?.(MODUL_KOMITE, id) !== true) setPesan(TABEL_KOMITE.takTerpasang)
  }

  if (galat === null && (hal === null || hal.baris.length === 0) && laporan === null) {
    return selesai && kosong !== undefined ? <>{kosong}</> : null
  }

  return (
    <section className="inbox claimlife__komite" aria-labelledby="claimlife-judul-komite">
      {utama ? (
        <header className="inbox__kepala">
          <h2 id="claimlife-judul-komite" className="inbox__judul">
            {TABEL_KOMITE.judul}
          </h2>
        </header>
      ) : (
        <h3 id="claimlife-judul-komite" className="inbox__judul">
          {TABEL_KOMITE.judul}
        </h3>
      )}
      {pesan && (
        <div className="alert alert--error" role="alert">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {hal !== null && hal.baris.length > 0 && (
        <div className="inbox__gulir">
        <table className="inbox__tabel claimlife__tabel">
          <thead>
            <tr>
              <th>{TABEL_KOMITE.kasusId}</th>
              <th>{TABEL_KOMITE.tglUpdate}</th>
              <th>{TABEL_KOMITE.statusWork}</th>
              <th>{TABEL_KOMITE.nomorKlaim}</th>
              <th className="claimlife__tengah">{TABEL_KOMITE.tingkat}</th>
              <th className="claimlife__angka">{TABEL_KOMITE.nilaiKlaim}</th>
              <th>{TABEL_KOMITE.mataUang}</th>
              <th>{TABEL_KOMITE.statusBaris}</th>
            </tr>
          </thead>
          <tbody>
            {hal.baris.map((b) => (
              <tr key={b.kasusId} className="inbox__baris" onClick={() => buka(b.kasusId)}>
                <td>{b.kasusId}</td>
                <td>{waktuKomite(b.tglUpdate)}</td>
                <td>{statusWorkKomite(b.statusWork)}</td>
                <td>{selKomite(b.nomorKlaim)}</td>
                <td className="claimlife__tengah">{tingkatKomite(b)}</td>
                <td className="claimlife__angka">{nilaiKomite(b.nilaiKlaim)}</td>
                <td>{selKomite(b.mataUang)}</td>
                <td>
                  <span className="claimlife__chip">{selKomite(b.statusBaris)}</span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        </div>
      )}
      {hal !== null && hal.baris.length > 0 && (
        // Paginasi inti (Previous / Page x of N / Next), seperti Inbox Claim Life (12-10-2026).
        <Halaman halaman={hal.halaman} ukuran={hal.ukuran} total={hal.total} onPindah={setHalaman} />
      )}
      {laporan !== null && (
        <section className="inbox__intervensi">
          <h3>
            {TABEL_KOMITE.laporan} ({laporan.tanggal})
          </h3>
          {laporan.kosong ? (
            <p role="status">{TABEL_KOMITE.laporanKosong}</p>
          ) : (
            <ul>
              {laporan.baris.map((b, i) => (
                <li key={`${b.kasusId}-${b.jenis}-${String(i)}`}>
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => buka(b.kasusId)}>
                    {b.kasusId}
                  </button>{' '}
                  {b.jenis} — {b.keadaan} {TABEL_KOMITE.sejak} {b.sejak}
                </li>
              ))}
            </ul>
          )}
        </section>
      )}
    </section>
  )
}
