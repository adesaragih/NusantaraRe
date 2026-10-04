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

import { HASIL_NOMOR_PL, KEPUTUSAN_POLIS, KONFIRMASI_KEPUTUSAN } from '../labels'
import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilPeriodeProduksi,
  putuskanPenawaran,
  TAHAP_POLIS,
  type AkibatKeputusanPolis,
} from '../api'
import { tanggalTampil } from '../tanggal'
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
  confirmTerkunci = false,
  onSelesai,
}: {
  polisID: string
  /** Tahap berjalan — `pyWorkStatus`, bukan posisi layar. */
  tahap: string
  /**
   * Confirm DIKUNCI: form Premium List Detail punya perubahan yang belum
   * disimpan (keputusan work owner 03-10-2026) - Confirm memakai data yang
   * TERSIMPAN, jadi perubahan di layar akan diabaikan diam-diam.
   */
  confirmTerkunci?: boolean
  onSelesai: () => void
}) {
  const [akibat, setAkibat] = useState<AkibatKeputusanPolis | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  // Keputusan yang menunggu dikonfirmasi lewat popup (03-10-2026): tombol
  // Decision TIDAK langsung menjalankan apa pun — mencegah tertekan tanpa sengaja.
  const [tanya, setTanya] = useState<typeof KEPUTUSAN_POLIS.confirm | typeof KEPUTUSAN_POLIS.decline | null>(null)

  // PL Number yang baru terbit — popup hasil sebelum kembali ke kotak masuk.
  const [nomorTerbit, setNomorTerbit] = useState<string | null>(null)
  const [wpcTerbit, setWpcTerbit] = useState('')

  /** Tombol "Yes" popup: baru di sini keputusannya benar-benar dikirim. */
  function lanjutkan(): void {
    if (tanya === null) return
    // Perubahan belum disimpan muncul SESUDAH popup terbuka: Confirm tetap ditahan.
    if (tanya === KEPUTUSAN_POLIS.confirm && confirmTerkunci) {
      setTanya(null)
      return
    }
    const keputusan = tanya
    setTanya(null)
    void jalankan(() => putuskanPenawaran(polisID, keputusan))
  }
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
      // PL Number yang baru terbit DITAMPILKAN dulu; kembali ke kotak masuk
      // sesudah popupnya ditutup (keputusan work owner 03-10-2026).
      if ((hasil.plNumber ?? '').trim() !== '') {
        setWpcTerbit(hasil.wpc ?? '')
        setNomorTerbit(hasil.plNumber ?? '')
        return
      }
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
  // dan periodenya tampil di kepala Premium List Detail (02-10-2026), bukan di
  // panel Decision.
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
      {!diDetail && galatPeriode !== null && <Gagal galat={galatPeriode} />}

      {/*
        Tiket 01 bagian 3 — isian `InputOfferLife.xml`, HANYA di tahap
        penawaran (`Assignment2`). `Confirm` di tahap ini ditolak server (409)
        selama System Reinsurance, Class of Business, dan Comment belum
        tersimpan lewat `Save Offer`.
      */}
      {tahap === TAHAP_POLIS.penawaran && <FormPenawaran polisID={polisID} />}

      <section className="panel pl-keputusan">
        <h3 className="panel__title">Decision</h3>
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
            disabled={sibuk || confirmTerkunci}
            onClick={() => {
              setTanya(KEPUTUSAN_POLIS.confirm)
            }}
          >
            {KEPUTUSAN_POLIS.confirm}
          </button>
          {/*
            ⛔ Tombol `Reject` (Transition9, kembali ke Input Offer Life) SENGAJA
            tidak ditampilkan di tahap mana pun — keputusan work owner
            03-10-2026. Jalurnya di backend tetap ada (`models.TransisiPenawaran`).
          */}
          <button
            type="button"
            className="btn btn--danger"
            disabled={sibuk}
            onClick={() => {
              setTanya(KEPUTUSAN_POLIS.decline)
            }}
          >
            {KEPUTUSAN_POLIS.decline}
          </button>
        </div>
        {/* Sebab Confirm terkunci DIKATAKAN, bukan dibiarkan ditebak. */}
        {confirmTerkunci && (
          <p className="pl-offer__kurang" role="status">
            {KONFIRMASI_KEPUTUSAN.belumTersimpan}
          </p>
        )}
      </section>

      {tanya !== null && (
        <Modal
          judul={
            tanya === KEPUTUSAN_POLIS.confirm ? KONFIRMASI_KEPUTUSAN.judulConfirm : KONFIRMASI_KEPUTUSAN.judulDecline
          }
          onTutup={() => {
            setTanya(null)
          }}
          aksi={
            tanya === KEPUTUSAN_POLIS.decline ? (
              <button type="button" className="btn btn--danger" disabled={sibuk} onClick={lanjutkan}>
                {KONFIRMASI_KEPUTUSAN.ya}
              </button>
            ) : (
              <button type="button" className="btn btn--primary" disabled={sibuk} onClick={lanjutkan}>
                {KONFIRMASI_KEPUTUSAN.ya}
              </button>
            )
          }
        >
          <p>{kalimatKonfirmasi(tanya, tahap)}</p>
        </Modal>
      )}

      {/*
        Hasil Confirm yang menerbitkan PL Number — DITAMPILKAN sebelum kembali
        ke kotak masuk; setiap cara menutup (OK, X, Escape, klik luar) baru
        kemudian meninggalkan layar (keputusan work owner 03-10-2026).
      */}
      {nomorTerbit !== null && (
        <Modal
          judul={HASIL_NOMOR_PL.judul}
          labelBatal={HASIL_NOMOR_PL.ok}
          onTutup={() => {
            setNomorTerbit(null)
            onSelesai()
          }}
        >
          {/* PL Number dan WPC bersama (03-10-2026). */}
          <dl className="pl-hasil" role="status">
            <dt>{HASIL_NOMOR_PL.labelNomor}</dt>
            <dd className="pl-hasil-nomor">{nomorTerbit}</dd>
            {wpcTerbit !== '' && (
              <>
                <dt>{HASIL_NOMOR_PL.labelWpc}</dt>
                <dd className="pl-hasil-nomor">{tanggalTampil(wpcTerbit)}</dd>
              </>
            )}
          </dl>
          <p>{HASIL_NOMOR_PL.kalimat}</p>
        </Modal>
      )}
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

/** Kalimat popup konfirmasi: AKIBAT keputusannya di tahap ini (03-10-2026). */
export function kalimatKonfirmasi(keputusan: string, tahap: string): string {
  if (keputusan === KEPUTUSAN_POLIS.decline) return KONFIRMASI_KEPUTUSAN.decline
  return tahap === TAHAP_POLIS.detail ? KONFIRMASI_KEPUTUSAN.confirmDetail : KONFIRMASI_KEPUTUSAN.confirmPenawaran
}
