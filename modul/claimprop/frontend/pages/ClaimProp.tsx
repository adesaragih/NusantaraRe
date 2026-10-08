// Halaman awal Claim Prop - tampilan pola Kelola User (keputusan work owner 07-10 dan 08-10-2026): dua tab, Process dan
// Resolve. Tab Process bawaan = worklist pembuat (Assignment2 "Outstanding Claim", tanpa cek workbasket - XML apa
// adanya); switch Teknik = Assignment1 "Input Acceptation" (workbasket ReasKlaimTeknik), dapat dinyalakan hanya oleh
// anggota workbasket itu. Add Claim (Start1 -> Assignment2; harness New tidak diekspor - OQ-CP-13) hanya saat switch
// Teknik mati.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { ambilHak, buatKasus, daftarKasus, type HakPelaku, type RingkasanKasus } from '../api'
import LayarKasus from '../components/LayarKasus'
import { CP } from '../labels'
import { bolehTambahKlaim, jenisDaftar, switchTeknikAktif, TAB_INBOX, type TabInbox } from './inbox'

const LABEL_TAB: Record<TabInbox, string> = { proses: CP.tabProses, selesai: CP.tabResolve }
const JEDA_CARI_MS = 300

export default function ClaimProp({ pelaku }: { pelaku: string }) {
  const [tab, setTab] = useState<TabInbox>('proses')
  const [teknik, setTeknik] = useState(false)
  const [hak, setHak] = useState<HakPelaku | null>(null)
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  const [daftar, setDaftar] = useState<RingkasanKasus[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [buka, setBuka] = useState<string | null>(null)
  const [segar, setSegar] = useState(0)

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
    setDaftar(null)
    daftarKasus(jenisDaftar(tab, teknikNyala), cari).then(
      (d) => {
        setDaftar(d)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
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
        onKembali={() => {
          setBuka(null)
          setSegar((s) => s + 1)
        }}
      />
    )
  }

  return (
    <section className="inbox claimprop__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{CP.judul}</h2>
      </header>
      <StripTab tab={TAB_INBOX} aktif={tab} onPilih={setTab} label={(t) => LABEL_TAB[t]} />

      <div className="toolbar">
        {tab === 'proses' && (
          <span className="claimprop__saklar-bungkus" title={bolehTeknik ? undefined : CP.teknikTanpaHak}>
            <span className="claimprop__saklar-label">{CP.inbox}</span>
            <button
              type="button"
              role="switch"
              aria-checked={teknikNyala}
              aria-labelledby="claimprop-saklar-teknik"
              className={'claimprop__saklar' + (teknikNyala ? ' claimprop__saklar--nyala' : '')}
              disabled={!bolehTeknik}
              onClick={() => setTeknik(!teknikNyala)}
            >
              <span className="claimprop__saklar-kenop" />
            </button>
            <span id="claimprop-saklar-teknik" className="claimprop__saklar-label">
              {CP.wbTeknik}
            </span>
          </span>
        )}
        <input
          className="field__input claimprop__cari"
          type="search"
          aria-label={CP.cari}
          placeholder={CP.cari}
          value={kata}
          onChange={(e) => setKata(e.target.value)}
        />
        <span className="toolbar__spacer" />
        {bolehTambahKlaim(tab, teknikNyala) && (
          <button type="button" className="btn btn--primary" onClick={baru}>
            {CP.tambahKlaim}
          </button>
        )}
      </div>

      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={CP.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={CP.kosong} />}
      {daftar !== null && daftar.length > 0 && (
        <table className="inbox__tabel claimprop__tabel">
          <thead>
            <tr>
              <th>{CP.kolomID}</th>
              <th>{CP.kolomTahap}</th>
              <th>{CP.kolomNoClaim}</th>
              <th>{CP.kolomPolis}</th>
              <th>{CP.kolomTreaty}</th>
              <th>{CP.kolomPembuat}</th>
              <th>{CP.kolomTanggal}</th>
            </tr>
          </thead>
          <tbody>
            {daftar.map((k) => (
              <tr key={k.id} className="inbox__baris" onClick={() => setBuka(k.id)}>
                <td>{k.id}</td>
                <td>{k.statusWork || k.label}</td>
                <td>{k.noClaim || k.claimNoTemp || <span className="muted">—</span>}</td>
                <td>{k.policyNo || <span className="muted">—</span>}</td>
                <td>{k.treatyName || <span className="muted">—</span>}</td>
                <td>{k.pembuatNama}</td>
                <td>{k.tglCreate}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}
