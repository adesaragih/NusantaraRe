// Layar rekap premium list — tiket 05a bagian 2.
//
// Meniru `Section/ShowLifePremiumSummary.xml`: satu grid rekap per `Type`
// (empat layout ber-`pyContainerVisibleWhen`), dan tombol `Submit` (b27471).
//
// ⛔ ANGKANYA DARI SERVER, TIDAK DIHITUNG DI SINI. Rumus `PREMIUM`/`BALANCE`
// bercabang per `Type` dengan tiga keanehan tanda warisan; salinan kedua di
// TypeScript adalah salinan yang tidak diuji literal 975/978/897/1789.
//
// ⛔ UANG TETAP TEKS. Tidak satu pun sel melewati `Number(...)`.
//
// ⚠️ `Submit` menyimpan nomor + rekap + salinan peserta warisan LALU menutup
// kasus Resolved-Completed (tiket 05b: `finishAssignment` b26442 →
// `Transition2` → `END52`). Sesudahnya tombolnya mati, dan layar mengatakannya.

import { useCallback, useEffect, useState } from 'react'

import { GRID_REKAP, SUMMARY_POLIS, type KolomRekap } from '../../assets/labels.premiumlist'
import { Gagal, Kosong, Memuat } from '../../components/ui/dasar'
import {
  ambilKepalaPolis,
  ambilRekapPolis,
  submitRekapPolis,
  type KepalaPolis,
  type RekapMataUangPolis,
} from '../../services/api'

/** Kolom grid untuk satu `Type`; tipe asing → tidak ada grid (bukan grid QR). */
export function kolomGridRekap(tipe: string): readonly KolomRekap[] {
  return GRID_REKAP[tipe] ?? []
}

/**
 * Isi satu sel rekap.
 *
 * ⚠️ `COB` tidak punya sumber di baris rekap (lihat `SUMMARY_POLIS`), jadi
 * ditandai kosong — bukan diisi tebakan.
 */
export function selRekap(r: RekapMataUangPolis, kolom: string, plNumber: string): string {
  let v: string | undefined
  switch (kolom) {
    case 'CURRENCY':
      v = r.currency
      break
    case 'PREMIUM':
      v = r.premium
      break
    case 'COMMISSION':
      v = r.commission
      break
    case 'BALANCE':
      v = r.balance
      break
    case 'PL_NUMBER':
      v = plNumber
      break
    case 'COB':
      v = undefined
      break
    default:
      v = r.jumlah[kolom]
  }
  return v === undefined || v.trim() === '' ? '—' : v
}

/**
 * Kalimat efek keluar — tiket 06.
 *
 * ⛔ Kegagalan efek keluar DIKATAKAN, tetapi tidak dibuat tampak seperti
 * simpan yang gagal: rekapnya sudah tersimpan, yang tertunda kirimannya.
 */
export function kalimatEfek(e: { dilewati: boolean; gagal: string[]; tidakTerantre: number }): string {
  if (e.dilewati) return SUMMARY_POLIS.efekDilewati
  if (e.gagal.length === 0) return ''
  const antre = e.tidakTerantre > 0 ? SUMMARY_POLIS.efekTidakTerantre : SUMMARY_POLIS.efekTerantre
  return `${SUMMARY_POLIS.efekGagal} ${e.gagal.join(', ')}. ${antre}`
}

export default function PremiumListSummary({ polisID }: { polisID: string }) {
  const [kepala, setKepala] = useState<KepalaPolis | null>(null)
  const [tipe, setTipe] = useState('')
  const [rekap, setRekap] = useState<RekapMataUangPolis[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(true)
  const [kirim, setKirim] = useState(false)
  const [kabar, setKabar] = useState('')
  const [selesai, setSelesai] = useState(false)

  const muat = useCallback(async () => {
    setSibuk(true)
    setGalat(null)
    try {
      const [k, r] = await Promise.all([ambilKepalaPolis(polisID), ambilRekapPolis(polisID)])
      setKepala(k)
      setTipe(r.tipe)
      setRekap(r.rekap)
    } catch (e) {
      setGalat(e)
      setRekap(null)
    } finally {
      setSibuk(false)
    }
  }, [polisID])

  useEffect(() => {
    void muat()
  }, [muat])

  async function submit(): Promise<void> {
    if (kirim || selesai) return
    setKirim(true)
    setGalat(null)
    setKabar('')
    try {
      const hasil = await submitRekapPolis(polisID)
      setRekap(hasil.rekap)
      setKepala((lama) => (lama === null ? lama : { ...lama, plNumber: hasil.nomor.nomor }))
      setSelesai(true)
      setKabar(
        `PL_NUMBER ${hasil.nomor.nomor}: ${String(hasil.rekap.length)} rekap mata uang ` +
          `tersimpan, ${String(hasil.pesertaWarisan)} peserta tersalin. ${SUMMARY_POLIS.ditutup} ` +
          kalimatEfek(hasil.efekKeluar),
      )
    } catch (e) {
      setGalat(e)
    } finally {
      setKirim(false)
    }
  }

  const kolom = kolomGridRekap(tipe)
  const plNumber = kepala?.plNumber ?? ''

  return (
    <section className="pl-summary">
      <h2 className="pl-summary__judul">{SUMMARY_POLIS.judul}</h2>
      {sibuk && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {kabar !== '' && <p role="status">{kabar}</p>}

      {rekap !== null && rekap.length === 0 && <Kosong pesan={SUMMARY_POLIS.tanpaRekap} />}

      {rekap !== null && rekap.length > 0 && kolom.length > 0 && (
        <table className="pl-summary__tabel">
          <thead>
            <tr>
              {kolom.map((k) => (
                <th key={k.kolom}>{k.judul}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rekap.map((r) => (
              <tr key={r.currency}>
                {kolom.map((k) => (
                  <td key={k.kolom}>{selRekap(r, k.kolom, plNumber)}</td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {rekap !== null && rekap.length > 0 && (
        <p className="pl-summary__catatan">{SUMMARY_POLIS.cobTanpaSumber}</p>
      )}

      <button
        type="button"
        className="pl-summary__submit"
        disabled={sibuk || kirim || selesai || rekap === null || rekap.length === 0}
        onClick={() => {
          void submit()
        }}
      >
        {SUMMARY_POLIS.submit}
      </button>
    </section>
  )
}
