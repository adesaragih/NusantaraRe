// Halaman Kelola User — daftar akun M_LOGIN_GO dan aksinya (keputusan work
// owner 01-10-2026): Tambah, Ubah (tab Profil dan Security), Buka kunci,
// Nonaktifkan/Aktifkan, Hapus permanen (dengan konfirmasi).
//
// Akun sendiri tidak punya tombol Nonaktifkan dan Hapus; backend tetap
// menolaknya (409), dan menolak perubahan yang menyisakan nol admin aktif.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../components/ui/dasar'
import { KELOLA_USER } from '../labels'
import {
  ambilDaftarPengguna,
  bukaKunciPengguna,
  hapusPengguna,
  setelAktifPengguna,
  type RingkasAkun,
  type RinciAkun,
} from './api'
import { saringDaftar, teksJenjang } from './aturan'
import FormPengguna from './FormPengguna'

/** Lencana status satu akun. */
function Status({ a }: { a: RingkasAkun }) {
  return (
    <span className="kelola-user__status">
      <span className={`badge ${a.aktif ? 'kelola-user__badge--aktif' : 'kelola-user__badge--nonaktif'}`}>
        {a.aktif ? KELOLA_USER.aktif : KELOLA_USER.nonaktif}
      </span>
      {a.terkunci && <span className="badge kelola-user__badge--terkunci">{KELOLA_USER.terkunci}</span>}
      {a.wajibGantiSandi && <span className="badge kelola-user__badge--wajib">{KELOLA_USER.wajibGanti}</span>}
    </span>
  )
}

export default function KelolaUser({
  akunSaya,
  onDiriBerubah,
}: {
  /** Akun yang sedang login. */
  akunSaya: string
  /** Dipanggil sesudah akun SENDIRI berubah - App membaca ulang sesi dan menu. */
  onDiriBerubah: () => void
}) {
  const [daftar, setDaftar] = useState<RingkasAkun[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [kueri, setKueri] = useState('')
  const [pesan, setPesan] = useState<string | null>(null)
  const [galatAksi, setGalatAksi] = useState<unknown>(null)
  // `undefined` = form tertutup; `null` = tambah; teks = ubah akun itu.
  const [form, setForm] = useState<string | null | undefined>(undefined)
  const [hapus, setHapus] = useState<RingkasAkun | null>(null)
  const [sibuk, setSibuk] = useState<string | null>(null)

  const muat = useCallback(() => {
    setGalat(null)
    ambilDaftarPengguna().then(setDaftar, (g: unknown) => {
      setGalat(g)
    })
  }, [])
  useEffect(() => {
    muat()
  }, [muat])

  const aksi = (akunId: string, jalan: () => Promise<unknown>, berhasil: string) => {
    if (sibuk !== null) return
    setSibuk(akunId)
    setPesan(null)
    setGalatAksi(null)
    jalan().then(
      () => {
        setSibuk(null)
        setPesan(berhasil)
        muat()
        if (akunId === akunSaya) onDiriBerubah()
      },
      (g: unknown) => {
        setSibuk(null)
        setGalatAksi(g)
      },
    )
  }

  const tersimpan = (r: RinciAkun) => {
    setForm(undefined)
    setGalatAksi(null)
    setPesan(KELOLA_USER.tersimpan(r.akunId))
    muat()
    if (r.akunId === akunSaya) onDiriBerubah()
  }

  const tampil = daftar === null ? [] : saringDaftar(daftar, kueri)

  return (
    <section className="inbox kelola-user">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{KELOLA_USER.judul}</h2>
      </header>
      <p className="muted kelola-user__sub">{KELOLA_USER.sub}</p>

      <div className="toolbar">
        <input
          className="field__input kelola-user__cari"
          type="search"
          aria-label={KELOLA_USER.cari}
          placeholder={KELOLA_USER.cari}
          value={kueri}
          onChange={(e) => {
            setKueri(e.target.value)
          }}
        />
        <span className="toolbar__spacer" />
        <button
          type="button"
          className="btn btn--primary"
          onClick={() => {
            setPesan(null)
            setForm(null)
          }}
        >
          {KELOLA_USER.tambah}
        </button>
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galatAksi !== null && <Gagal galat={galatAksi} />}
      {daftar === null && galat === null && <Memuat pesan={KELOLA_USER.memuat} />}
      {galat !== null && <Gagal galat={galat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={KELOLA_USER.kosong} />}
      {daftar !== null && daftar.length > 0 && tampil.length === 0 && <Kosong pesan={KELOLA_USER.tidakCocok} />}
      {tampil.length > 0 && (
        <table className="inbox__tabel kelola-user__tabel">
          <thead>
            <tr>
              <th>{KELOLA_USER.kolomAkun}</th>
              <th>{KELOLA_USER.kolomNama}</th>
              <th>{KELOLA_USER.kolomJabatan}</th>
              <th>{KELOLA_USER.kolomJenjang}</th>
              <th>{KELOLA_USER.kolomStatus}</th>
              <th>{KELOLA_USER.kolomLogin}</th>
              <th className="table__actions">{KELOLA_USER.kolomAksi}</th>
            </tr>
          </thead>
          <tbody>
            {tampil.map((a) => {
              const saya = a.akunId === akunSaya
              const kunciBaris = sibuk !== null
              return (
                <tr key={a.akunId} className="inbox__baris">
                  <td>
                    {a.akunId}
                    {saya && <span className="muted"> ({KELOLA_USER.anda})</span>}
                  </td>
                  <td>
                    {a.nama}
                    {/* Email di bawah nama, bukan kolom sendiri: alamat panjang mendorong kolom Aksi keluar layar laptop. */}
                    {a.email !== '' && <div className="muted">{a.email}</div>}
                  </td>
                  <td>{a.jabatan === '' ? <span className="muted">—</span> : a.jabatan}</td>
                  <td>{teksJenjang(a)}</td>
                  <td>
                    <Status a={a} />
                  </td>
                  <td>{a.loginTerakhir === '' ? <span className="muted">{KELOLA_USER.belumLogin}</span> : a.loginTerakhir}</td>
                  <td className="table__actions">
                    <span className="kelola-user__aksi">
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        disabled={kunciBaris}
                        onClick={() => {
                          setPesan(null)
                          setForm(a.akunId)
                        }}
                      >
                        {KELOLA_USER.ubah}
                      </button>
                      {a.terkunci && (
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          disabled={kunciBaris}
                          onClick={() => {
                            aksi(a.akunId, () => bukaKunciPengguna(a.akunId), KELOLA_USER.dibukaKunci(a.akunId))
                          }}
                        >
                          {KELOLA_USER.bukaKunci}
                        </button>
                      )}
                      {!(saya && a.aktif) && (
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          disabled={kunciBaris}
                          onClick={() => {
                            aksi(
                              a.akunId,
                              () => setelAktifPengguna(a.akunId, !a.aktif),
                              a.aktif ? KELOLA_USER.dinonaktifkan(a.akunId) : KELOLA_USER.diaktifkan(a.akunId),
                            )
                          }}
                        >
                          {a.aktif ? KELOLA_USER.nonaktifkan : KELOLA_USER.aktifkan}
                        </button>
                      )}
                      {!saya && (
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm kelola-user__hapus"
                          disabled={kunciBaris}
                          onClick={() => {
                            setPesan(null)
                            setGalatAksi(null)
                            setHapus(a)
                          }}
                        >
                          {KELOLA_USER.hapus}
                        </button>
                      )}
                    </span>
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      )}

      {form !== undefined && (
        <FormPengguna
          key={form ?? ''}
          akunId={form}
          akunSaya={akunSaya}
          onTutup={() => {
            setForm(undefined)
          }}
          onTersimpan={tersimpan}
        />
      )}

      {hapus !== null && (
        <Modal
          judul={KELOLA_USER.judulHapus}
          onTutup={() => {
            if (sibuk === null) setHapus(null)
          }}
          labelBatal={KELOLA_USER.batal}
          aksi={
            <button
              type="button"
              className="btn btn--danger"
              disabled={sibuk !== null}
              onClick={() => {
                const sasaran = hapus.akunId
                // Modal tetap terbuka selama penghapusan berjalan; tertutup sesudahnya.
                aksi(sasaran, () => hapusPengguna(sasaran).finally(() => setHapus(null)), KELOLA_USER.terhapus(sasaran))
              }}
            >
              {sibuk !== null ? KELOLA_USER.menghapus : KELOLA_USER.hapusPermanen}
            </button>
          }
        >
          <p>{KELOLA_USER.tanyaHapus(hapus.akunId)}</p>
        </Modal>
      )}
    </section>
  )
}
