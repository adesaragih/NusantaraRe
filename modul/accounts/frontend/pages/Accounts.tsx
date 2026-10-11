// Halaman Accounts (keputusan work owner 04-10-2026): daftar `T_M_ACCOUNT` dicari dan dibagi per halaman di server
// (18 ribu baris), terbaru dulu; "Add" membuka form akun baru. Akun hanya diinput sekali - nol View, nol Edit, nol
// hapus. Tampilan sama dengan Kelola User (tema disalin ke `accounts.css`).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat, Modal, UKURAN_HALAMAN } from '../../../../inti/frontend/components/ui/dasar'
import { useBolehUbah } from '../../../../inti/frontend/lib/hakMenu'
import { ambilDaftar, ambilPilihan, type Account, type GroupBusiness, type HalamanDaftar } from '../api'
import { atauStrip } from '../aturan'
import FormAccount from '../components/FormAccount'
import { ACC } from '../labels'
import { NAMA_ACC } from '../menu'

/** Jeda ketik sebelum server ditanya. */
const JEDA_MS = 300

export default function Accounts({
  mulaiBaru = false,
  onSelesaiBaru,
}: {
  mulaiBaru?: boolean
  /**
   * Diisi = mode POPUP SAJA (tombol New Group Business NB FacIn, `MenuModul.formBaruBerjendela`, work owner 10-10-2026
   * "HANYA MUNCULIN INI SAJA, MENU UTAMA ACCOUNT NYA TIDAK USAH"): hanya popup Add yang dirender, tanpa daftar;
   * Cancel / X / Save memanggil ini.
   */
  onSelesaiBaru?: () => void
} = {}) {
  const hanyaPopup = onSelesaiBaru !== undefined
  // Menu View only (M_LOGIN_GO_MENU.HAK, 04-10-2026): tanpa tombol Add; backend juga menolak tulisnya.
  const bolehUbah = useBolehUbah(NAMA_ACC)
  const [kata, setKata] = useState('')
  const [cari, setCari] = useState('')
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<HalamanDaftar | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [groupBusiness, setGroupBusiness] = useState<GroupBusiness[]>([])
  // `mulaiBaru` (PropsRute, 10-10-2026): dibuka tombol New Group Business form Opportunity NB FacIn - langsung form Add.
  const [formBuka, setFormBuka] = useState(mulaiBaru && bolehUbah)
  const [pesan, setPesan] = useState<string | null>(null)

  const muat = useCallback(() => {
    setData(null)
    ambilDaftar(cari, halaman, UKURAN_HALAMAN).then(
      (d) => {
        setData(d)
        setGalat(null)
      },
      (g: unknown) => {
        setGalat(g)
      },
    )
  }, [cari, halaman])

  useEffect(() => {
    // Mode popup saja tidak menampilkan daftar - tidak perlu dimuat.
    if (!hanyaPopup) muat()
  }, [muat, hanyaPopup])

  useEffect(() => {
    ambilPilihan().then(
      (p) => {
        setGroupBusiness(p.groupBusiness)
      },
      (g: unknown) => {
        setGalat(g)
      },
    )
  }, [])

  useEffect(() => {
    const jam = setTimeout(() => {
      setCari(kata.trim())
      setHalaman(1)
    }, JEDA_MS)
    return () => {
      clearTimeout(jam)
    }
  }, [kata])

  const tersimpan = (a: Account) => {
    setFormBuka(false)
    setPesan(ACC.tersimpan(a.idView))
    muat()
  }

  if (hanyaPopup) {
    // Akar tetap dipasang: token tema `--acc-*` dan gaya popup tinggal di `.accounts__akar`; `--popup` tanpa kotak.
    return (
      <section className="accounts__akar accounts__akar--popup">
        {bolehUbah ? (
          <FormAccount groupBusiness={groupBusiness} onTutup={onSelesaiBaru} onTersimpan={onSelesaiBaru} />
        ) : (
          <Modal judul={ACC.judulBaru} onTutup={onSelesaiBaru} labelBatal={ACC.batal}>
            <p>{ACC.tanpaHakTambah}</p>
          </Modal>
        )}
      </section>
    )
  }

  return (
    <section className="inbox accounts__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{ACC.judul}</h2>
      </header>

      <div className="toolbar">
        <input
          className="field__input accounts__cari"
          type="search"
          aria-label={ACC.cari}
          placeholder={ACC.cari}
          value={kata}
          onChange={(e) => {
            setKata(e.target.value)
          }}
        />
        <span className="toolbar__spacer" />
        {bolehUbah && (
          <button
            type="button"
            className="btn btn--primary"
            onClick={() => {
              setPesan(null)
              setFormBuka(true)
            }}
          >
            {ACC.tambah}
          </button>
        )}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {data === null && galat === null && <Memuat pesan={ACC.memuat} />}
      {galat !== null && <Gagal galat={galat} />}
      {data !== null && data.total === 0 && <Kosong pesan={cari === '' ? ACC.kosong : ACC.tidakCocok} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <table className="inbox__tabel accounts__tabel">
            <thead>
              <tr>
                <th>{ACC.kolomId}</th>
                <th>{ACC.insuredName}</th>
                <th>{ACC.groupBusiness}</th>
                <th>{ACC.owner}</th>
                <th>{ACC.createDate}</th>
              </tr>
            </thead>
            <tbody>
              {data.daftar.map((a) => (
                <tr key={a.id} className="inbox__baris">
                  <td>{atauStrip(a.idView)}</td>
                  <td>
                    {atauStrip(a.insuredName)}
                    {/* Org ID di bawah nama, bukan kolom sendiri - pola Contact ID di Kelola User. */}
                    {a.orgId !== '' && <span className="muted accounts__kecil">{a.orgId}</span>}
                  </td>
                  <td>{atauStrip(a.groupBusiness)}</td>
                  <td>{a.createOp === '' ? <span className="muted">—</span> : a.createOp}</td>
                  <td>{a.createDate === '' ? <span className="muted">—</span> : a.createDate}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <Halaman halaman={data.halaman} ukuran={data.ukuran} total={data.total} onPindah={setHalaman} />
        </>
      )}

      {formBuka && (
        <FormAccount
          groupBusiness={groupBusiness}
          onTutup={() => {
            setFormBuka(false)
          }}
          onTersimpan={tersimpan}
        />
      )}
    </section>
  )
}
