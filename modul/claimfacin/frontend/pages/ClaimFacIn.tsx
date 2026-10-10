// Halaman awal Claim Fac In - pola Claim Non Prop (disalin dari `modul/claimnonprop/frontend/pages/ClaimNonProp.tsx`,
// bukan impor): dua tab, Process dan Resolve. Tab Process bawaan = worklist pembuat (Assignment1 Input Register +
// Assignment7 Input Estimasi, `Register_Flow`); switch Teknik = Assignment3 "Choose Surveyor" (workbasket
// ReasKlaimTeknik), dapat dinyalakan hanya oleh anggota workbasket itu. Add Claim (Start -> Assignment1; harness New
// tidak diekspor) hanya saat switch Teknik mati. Tabel komite di bawah inbox (modul Komite Claim Fac In tanpa menu,
// prompt tahap 2 §2 butir 2; pola Claim Non Prop): baris dibuka DI TEMPAT lewat `onBukaModul`.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { ambilHak, buatKasus, daftarKasus, type HakPelaku, type RingkasanKasus } from '../api'
import LayarKasus from '../components/LayarKasus'
import TabelKomite from '../components/TabelKomite'
import { tampilTanggal } from '../ketikTanggal'
import { CFI } from '../labels'
import { bolehTambahKlaim, jenisDaftar, switchTeknikAktif, TAB_INBOX, tabelKomiteTampil, type TabInbox } from './inbox'

const LABEL_TAB: Record<TabInbox, string> = { proses: CFI.tabProses, selesai: CFI.tabResolve }
const JEDA_CARI_MS = 300

export default function ClaimFacIn({
  pelaku,
  onBukaModul,
  bukaKasus,
  onBeranda,
}: {
  pelaku: string
  /** `PropsRute.onBukaModul` - baris tabel komite: layar komite dibuka di tempat. */
  onBukaModul?: (modul: string, id: string) => boolean
  /** `PropsRute.bukaKasus` - satu berkas dibuka langsung; `hanyaLihat` = tampilan saja. */
  bukaKasus?: { id: string; ketuk: number; hanyaLihat?: boolean }
  /** `PropsRute.onBeranda` - Back berkas yang dibuka lewat `bukaKasus` kembali ke pemanggil. */
  onBeranda?: () => void
}) {
  const [tab, setTab] = useState<TabInbox>('proses')
  const [teknik, setTeknik] = useState(false)
  const [hak, setHak] = useState<HakPelaku | null>(null)
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  const [daftar, setDaftar] = useState<RingkasanKasus[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [buka, setBuka] = useState<string | null>(null)
  const [segar, setSegar] = useState(0)
  const [dariLuar, setDariLuar] = useState(false)
  const [hanyaLihat, setHanyaLihat] = useState(false)

  const ketukBuka = bukaKasus?.ketuk
  const idBuka = bukaKasus?.id
  const lihatBuka = bukaKasus?.hanyaLihat === true
  useEffect(() => {
    if (ketukBuka !== undefined && idBuka !== undefined) {
      setBuka(idBuka)
      setDariLuar(true)
      setHanyaLihat(lihatBuka)
    }
  }, [ketukBuka, idBuka, lihatBuka])

  useEffect(() => {
    const t = setTimeout(() => setCari(kata), JEDA_CARI_MS)
    return () => clearTimeout(t)
  }, [kata])

  useEffect(() => {
    ambilHak().then(setHak, () => setHak(null))
  }, [])

  const bolehTeknik = switchTeknikAktif(hak)
  const teknikNyala = teknik && bolehTeknik

  useEffect(() => {
    if (buka) return
    // Jawaban basi (tab / switch / kata cari lama) diabaikan - tidak menimpa daftar terbaru.
    let aktif = true
    setDaftar(null)
    daftarKasus(jenisDaftar(tab, teknikNyala), cari).then(
      (d) => {
        if (!aktif) return
        setDaftar(d)
        setGalat(null)
      },
      (g: unknown) => {
        if (aktif) setGalat(g)
      },
    )
    return () => {
      aktif = false
    }
  }, [tab, teknikNyala, cari, buka, segar])

  const baru = useCallback(() => {
    buatKasus().then(
      (l) => setBuka(l.kasus.id),
      (g: unknown) => setGalat(g),
    )
  }, [])

  if (buka) {
    return (
      <LayarKasus
        id={buka}
        pelaku={pelaku}
        hanyaLihat={hanyaLihat}
        onKembali={() => {
          setBuka(null)
          setHanyaLihat(false)
          setSegar((s) => s + 1)
          if (dariLuar) {
            setDariLuar(false)
            onBeranda?.()
          }
        }}
      />
    )
  }

  return (
    <section className="inbox claimfacin__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{CFI.judul}</h2>
      </header>
      <StripTab tab={TAB_INBOX} aktif={tab} onPilih={setTab} label={(t) => LABEL_TAB[t]} />

      <div className="toolbar">
        {tab === 'proses' && (
          <span className="claimfacin__saklar-bungkus" title={bolehTeknik ? undefined : CFI.teknikTanpaHak}>
            <span className="claimfacin__saklar-label">{CFI.inbox}</span>
            <button
              type="button"
              role="switch"
              aria-checked={teknikNyala}
              aria-labelledby="claimfacin-saklar-teknik"
              className={'claimfacin__saklar' + (teknikNyala ? ' claimfacin__saklar--nyala' : '')}
              disabled={!bolehTeknik}
              onClick={() => setTeknik(!teknikNyala)}
            >
              <span className="claimfacin__saklar-kenop" />
            </button>
            <span id="claimfacin-saklar-teknik" className="claimfacin__saklar-label">
              {CFI.wbTeknik}
            </span>
          </span>
        )}
        <input
          className="field__input claimfacin__cari"
          type="search"
          aria-label={CFI.cari}
          placeholder={CFI.cari}
          value={kata}
          onChange={(e) => setKata(e.target.value)}
        />
        <span className="toolbar__spacer" />
        {bolehTambahKlaim(tab, teknikNyala) && (
          <button type="button" className="btn btn--primary" onClick={baru}>
            {CFI.tambahKlaim}
          </button>
        )}
      </div>

      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={CFI.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={CFI.kosong} />}
      {daftar !== null && daftar.length > 0 && (
        <table className="inbox__tabel claimfacin__tabel">
          <thead>
            <tr>
              <th>{CFI.kolomID}</th>
              <th>{CFI.kolomTahap}</th>
              <th>{CFI.kolomNoClaim}</th>
              <th>{CFI.kolomPolis}</th>
              <th>{CFI.kolomTertanggung}</th>
              <th>{CFI.kolomDOL}</th>
              <th>{CFI.kolomPembuat}</th>
              <th>{CFI.kolomTanggal}</th>
            </tr>
          </thead>
          <tbody>
            {daftar.map((k) => (
              <tr key={k.id} className="inbox__baris" onClick={() => setBuka(k.id)}>
                <td>{k.id}</td>
                <td>{k.statusWork || k.label}</td>
                <td>{k.noClaim || <span className="muted">—</span>}</td>
                <td>{k.policyNo || <span className="muted">—</span>}</td>
                <td>{k.insuredName || <span className="muted">—</span>}</td>
                <td>{k.dateOfLoss ? tampilTanggal(k.dateOfLoss, 'tanggal') : <span className="muted">—</span>}</td>
                <td>{k.pembuatNama}</td>
                <td>{k.tglCreate}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {tabelKomiteTampil(hak) && <TabelKomite onBukaModul={onBukaModul} />}
    </section>
  )
}
