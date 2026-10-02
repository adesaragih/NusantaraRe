// Layar keputusan penawaran — tiket 01 bagian 2.
//
// Meniru `Section/InputOfferLife.xml` + `Section/ConfirmSection.xml`, yang
// keduanya memanggil `Activity/ProtectAccept.xml` (b639/b868, b11420/b11652)
// sebelum meneruskan.
//
// ⛔ KETIGA TOMBOL TIDAK SAMA TERSEDIANYA. `Reject` muncul TEPAT SEKALI di
// seluruh flow — `Transition9` b2306 pada `Decision2`, yaitu sesudah Input
// Premium Detail. Di tahap penawaran ia tidak punya konektor, jadi layar
// TIDAK menawarkannya di sana: tombol yang pasti dijawab 409 adalah tombol
// yang mengajari orang mengabaikan galat.
//
// ⛔ `Confirm` di tahap penawaran menyerahkan kasus ke `Decision3` — decision
// table `IsFlagOnGoingPolicy` atas bendera KASUS (`"0"` Input Offer → Offer,
// tutup; `"1"` Input Premium → Premium, pindah ke Input Premium Detail).
// ⭐ GILIRAN-14 butir bq: Pega tidak menanyakannya, dan layar ini pun tidak
// lagi — backend menerapkannya di dalam `Confirm`. Kedua tombol portal kini
// berbeda perilaku seperti di sistem lama.
//
// ⚠️ Gerbang `ProtectAccept` (`models.ValidasiPenawaran`) BELUM TERSAMBUNG
// ke rute mana pun — pemanggilnya hanya uji (sensus remark 28-09-2026). Bila
// kelak tersambung, layar ini hanya menampilkan kalimatnya VERBATIM.
// (`Please choose no offer !` ter-remark di XML dan dibuang.)

import { useEffect, useState } from 'react'

import { KEPUTUSAN_POLIS } from '../labels'
import { Gagal } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilPeriodeProduksi,
  bolehRejectDiTahap,
  putuskanPenawaran,
  TAHAP_POLIS,
  type AkibatKeputusanPolis,
} from '../api'
import FormPenawaran from './FormPenawaran'
import '../premiumlistlife.css'

/** Menyusun kalimat tentang akibat sebuah keputusan. */
export function ringkasanAkibat(a: AkibatKeputusanPolis): string {
  if (a.statusWork !== '') return `Case closed — ${a.statusWork}.`
  if (a.tahapTujuan !== '') return `Case moved to ${a.tahapTujuan}.`
  return 'Decision saved.'
}

/**
 * Judul layar keputusan, menurut tahapnya.
 *
 * ⛔ SATU LAYAR, DUA TAHAP. Ketiga tombol keputusan hidup di tahap penawaran
 * DAN di tahap Input Premium Detail (`Reject` bahkan HANYA di sana), jadi
 * layar ini muncul di keduanya. Judul "Input Offer" yang tetap saat orang
 * berada di Input Premium Detail memberi tahu mereka hal yang keliru tentang
 * di mana mereka berada.
 */
export function judulKeputusan(tahap: string): string {
  // Judul tahap penawaran = nama tahapnya, VERBATIM `pyWorkStatus` (`Input Offer Life`)
  // - permintaan work owner 01-10-2026.
  return tahap === TAHAP_POLIS.detail ? TAHAP_POLIS.detail : TAHAP_POLIS.penawaran
}

export default function InputOffer({
  polisID,
  tahap,
  onSelesai,
}: {
  polisID: string
  /** Tahap berjalan — `pyWorkStatus`, bukan posisi layar. */
  tahap: string
  onSelesai: () => void
}) {
  const [akibat, setAkibat] = useState<AkibatKeputusanPolis | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  // ⛔ Periode produksi DITAMPILKAN sebelum menyimpan - AC tiket 02.
  // Kosong berarti belum terbaca; galatnya dinyatakan, bukan disembunyikan,
  // sebab 503 di sini berarti POOLDATA.TANGGAL_CLOSING kosong.
  const [periode, setPeriode] = useState('')
  const [galatPeriode, setGalatPeriode] = useState<unknown>(null)

  useEffect(() => {
    let hidup = true
    void (async () => {
      try {
        const p = await ambilPeriodeProduksi()
        if (hidup) setPeriode(p)
      } catch (e) {
        if (hidup) setGalatPeriode(e)
      }
    })()
    return () => {
      hidup = false
    }
  }, [])

  async function jalankan(kerja: () => Promise<AkibatKeputusanPolis>): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    try {
      const hasil = await kerja()
      setAkibat(hasil)
      // Kasus yang tertutup atau berpindah tidak lagi milik layar ini — dan
      // sejak butir bq setiap keputusan yang berhasil menutup atau memindahkan.
      onSelesai()
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  // ⛔ Di tahap Input Premium Detail layar ini berdiri DI BAWAH Premium List
  // Detail (rute.tsx), yang sudah punya kepala sendiri — kepala kedua dibuang,
  // periode pindah ke panel Decision.
  const diDetail = tahap === TAHAP_POLIS.detail
  const lencanaPeriode = periode !== '' && (
    <span className="pl-kepala__chip" role="status">
      <span className="pl-kepala__label">Period</span> {periodeTampil(periode)}
    </span>
  )

  return (
    <section className="polis-offer">
      {/*
        Kepala ringkas (permintaan work owner 01-10-2026): judul, lalu nomor
        kasus dan periode produksi saja. Tahap tidak diulang — judulnya
        (`judulKeputusan`) sudah menyebutnya.
      */}
      {!diDetail && (
        <header className="pl-kepala">
          <h2 className="pl-kepala__judul">{judulKeputusan(tahap)}</h2>
          <div className="pl-kepala__meta">
            <span className="pl-kepala__chip">{polisID}</span>
            {/* ⛔ Periode produksi DITAMPILKAN sebelum menyimpan — AC tiket 02. */}
            {lencanaPeriode}
          </div>
        </header>
      )}
      {galatPeriode !== null && <Gagal galat={galatPeriode} />}

      {/*
        Tiket 01 bagian 3 — isian `InputOfferLife.xml`, HANYA di tahap
        penawaran (`Assignment2`). `Confirm` di tahap ini ditolak server (409)
        selama System Reinsurance, Class of Business, dan Comment belum
        tersimpan lewat `Save Offer`.
      */}
      {tahap === TAHAP_POLIS.penawaran && <FormPenawaran polisID={polisID} />}

      <section className="panel pl-keputusan">
        <h3 className="panel__title">Decision</h3>
        {diDetail && <div className="pl-kepala__meta pl-keputusan__meta">{lencanaPeriode}</div>}
        {galat !== null && <Gagal galat={galat} />}
        {akibat !== null && (
          <p className="pl-offer__tersimpan" role="status">
            {ringkasanAkibat(akibat)}
          </p>
        )}

        <div className="pl-keputusan__aksi">
          <button
            type="button"
            className="btn btn--primary"
            disabled={sibuk}
            onClick={() => {
              void jalankan(() => putuskanPenawaran(polisID, KEPUTUSAN_POLIS.confirm))
            }}
          >
            {KEPUTUSAN_POLIS.confirm}
          </button>
          {/* ⛔ Hanya bila tahapnya punya konektornya — lihat kepala berkas. */}
          {bolehRejectDiTahap(tahap) && (
            <button
              type="button"
              className="btn"
              disabled={sibuk}
              onClick={() => {
                void jalankan(() => putuskanPenawaran(polisID, KEPUTUSAN_POLIS.reject))
              }}
            >
              {KEPUTUSAN_POLIS.reject}
            </button>
          )}
          <button
            type="button"
            className="btn btn--danger"
            disabled={sibuk}
            onClick={() => {
              void jalankan(() => putuskanPenawaran(polisID, KEPUTUSAN_POLIS.decline))
            }}
          >
            {KEPUTUSAN_POLIS.decline}
          </button>
        </div>
      </section>
    </section>
  )
}

/**
 * Periode produksi untuk layar: `YYYY-MM` dari server menjadi `MM/YYYY`.
 *
 * Bentuk yang tidak dikenal ditampilkan APA ADANYA — mengubah teks yang tidak
 * dipahami berarti menampilkan periode yang tidak pernah dikirim server.
 */
export function periodeTampil(periode: string): string {
  const m = /^(\d{4})-(\d{1,2})$/.exec(periode.trim())
  const tahun = m?.[1]
  const bulan = m?.[2]
  if (tahun === undefined || bulan === undefined) return periode
  return `${bulan.padStart(2, '0')}/${tahun}`
}
