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

import { KEPUTUSAN_POLIS } from '../../assets/labels.premiumlist'
import { Gagal } from '../../components/ui/dasar'
import {
  ambilPeriodeProduksi,
  bolehRejectDiTahap,
  putuskanPenawaran,
  TAHAP_POLIS,
  type AkibatKeputusanPolis,
} from '../../services/api'

/** Menyusun kalimat tentang akibat sebuah keputusan. */
export function ringkasanAkibat(a: AkibatKeputusanPolis): string {
  if (a.statusWork !== '') return `Kasus ditutup — ${a.statusWork}.`
  if (a.tahapTujuan !== '') return `Kasus berpindah ke ${a.tahapTujuan}.`
  return 'Keputusan tersimpan.'
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
  return tahap === TAHAP_POLIS.detail ? 'Input Premium Detail' : 'Input Offer'
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

  return (
    <section className="polis-offer">
      <h2 className="polis-offer__judul">{judulKeputusan(tahap)}</h2>
      <p className="polis-offer__tahap">
        {polisID} — {tahap}
      </p>

      {periode !== '' && (
        <p className="polis-offer__periode" role="status">
          Periode produksi: {periode}
        </p>
      )}
      {galatPeriode !== null && <Gagal galat={galatPeriode} />}
      {galat !== null && <Gagal galat={galat} />}
      {akibat !== null && <p role="status">{ringkasanAkibat(akibat)}</p>}

      <p className="polis-offer__aksi">
        <button
          type="button"
          disabled={sibuk}
          onClick={() => {
            void jalankan(() => putuskanPenawaran(polisID, KEPUTUSAN_POLIS.confirm))
          }}
        >
          {KEPUTUSAN_POLIS.confirm}
        </button>{' '}
        {/* ⛔ Hanya bila tahapnya punya konektornya — lihat kepala berkas. */}
        {bolehRejectDiTahap(tahap) && (
          <>
            <button
              type="button"
              disabled={sibuk}
              onClick={() => {
                void jalankan(() => putuskanPenawaran(polisID, KEPUTUSAN_POLIS.reject))
              }}
            >
              {KEPUTUSAN_POLIS.reject}
            </button>{' '}
          </>
        )}
        <button
          type="button"
          disabled={sibuk}
          onClick={() => {
            void jalankan(() => putuskanPenawaran(polisID, KEPUTUSAN_POLIS.decline))
          }}
        >
          {KEPUTUSAN_POLIS.decline}
        </button>
      </p>
    </section>
  )
}
