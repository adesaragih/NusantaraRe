// Halaman Template Manager (keputusan work owner 04-10-2026) - berkas templat unduhan SEMUA menu di satu tempat:
// dikelompokkan per menu, versi aktif per templat, Upload versi baru (diperiksa dulu: ekstensi, jumlah kolom,
// judul kolom yang berubah), Riwayat versi dengan Activate untuk kembali ke versi lama atau ke berkas bawaan.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../components/ui/dasar'
import { pesanGalat } from '../klien'
import { TEMPLATE_MANAGER as TM } from '../labels'
import { aktifkan, ambilDaftar, ambilRiwayat, periksa, unduhTemplat, unggah, type HasilPeriksa, type Slot, type Versi } from './api'
import { kelompokSlot, saringSlot, ukuranBerkas, versiBerikut } from './aturan'

function teksGalat(g: unknown): string | null {
  return pesanGalat(g) ?? null
}

/** Popup Upload versi baru: pilih berkas, periksa di server, lalu simpan. */
function DialogUnggah({ slot, onTutup, onTersimpan }: { slot: Slot; onTutup: () => void; onTersimpan: (pesan: string) => void }) {
  const [berkas, setBerkas] = useState<File | null>(null)
  const [hasil, setHasil] = useState<HasilPeriksa | null>(null)
  const [catatan, setCatatan] = useState('')
  const [sibuk, setSibuk] = useState<'periksa' | 'simpan' | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [riwayat, setRiwayat] = useState<Versi[]>([])

  useEffect(() => {
    ambilRiwayat(slot.kode).then(
      (r) => {
        setRiwayat(r.versi)
      },
      () => {
        setRiwayat([])
      },
    )
  }, [slot.kode])

  const pilih = (f: File) => {
    setBerkas(f)
    setHasil(null)
    setGalat(null)
    setSibuk('periksa')
    periksa(slot.kode, f).then(
      (h) => {
        setSibuk(null)
        setHasil(h)
      },
      (g: unknown) => {
        setSibuk(null)
        setGalat(g)
      },
    )
  }

  const simpan = () => {
    if (berkas === null || hasil === null || hasil.galat.length > 0 || sibuk !== null) return
    setSibuk('simpan')
    setGalat(null)
    unggah(slot.kode, berkas, catatan).then(
      (r) => {
        setSibuk(null)
        onTersimpan(TM.tersimpan(r.versi))
      },
      (g: unknown) => {
        setSibuk(null)
        setGalat(g)
      },
    )
  }

  const bolehSimpan = berkas !== null && hasil !== null && hasil.galat.length === 0 && sibuk === null
  return (
    <Modal
      judul={TM.judulUnggah}
      onTutup={onTutup}
      onKirim={simpan}
      labelBatal={TM.batal}
      lebar
      aksi={
        <button type="submit" className="btn btn--primary" disabled={!bolehSimpan}>
          {sibuk === 'simpan' ? TM.menyimpan : TM.simpan(versiBerikut(riwayat))}
        </button>
      }
    >
      <dl className="templat-manager__info">
        <dt>{TM.untuk}</dt>
        <dd>
          {slot.grup} · {slot.nama}
        </dd>
        <dt>{TM.dipakaiDi}</dt>
        <dd>{slot.dipakaiDi === '' ? '—' : slot.dipakaiDi}</dd>
        <dt>{TM.syarat}</dt>
        <dd>{TM.syaratTeks(slot.ekstensi, slot.pemisah, slot.jumlahKolom)}</dd>
      </dl>
      <label className="field">
        <span className="field__label">{TM.pilihBerkas}</span>
        <input
          className="field__input"
          type="file"
          accept={slot.ekstensi}
          onChange={(e) => {
            const f = e.target.files?.[0]
            if (f !== undefined) pilih(f)
          }}
        />
      </label>
      {sibuk === 'periksa' && <Memuat pesan={TM.memeriksa} />}
      {galat !== null && (teksGalat(galat) !== null ? <div className="alert alert--error">{teksGalat(galat)}</div> : <Gagal galat={galat} />)}
      {hasil !== null && (
        <div className="templat-manager__cek" role="status">
          <strong>{TM.hasilCek}</strong>
          <ul>
            {hasil.galat.map((g) => (
              <li key={g} className="templat-manager__cek--gagal">
                ✗ {g}
              </li>
            ))}
            {hasil.galat.length === 0 && <li className="templat-manager__cek--lolos">✓ {TM.lolos(hasil.jumlahKolom)}</li>}
            {hasil.perbedaan.length > 0 && (
              <li className="templat-manager__cek--peringatan">
                ⚠ {TM.perbedaan(hasil.perbedaan.length)}
                <ul>
                  {hasil.perbedaan.map((p) => (
                    <li key={p.kolom}>
                      {TM.kolomKe(p.kolom)}: “{p.lama}” → “{p.baru}”
                    </li>
                  ))}
                </ul>
              </li>
            )}
          </ul>
        </div>
      )}
      <label className="field">
        <span className="field__label">{TM.catatan}</span>
        <textarea
          className="field__input"
          rows={2}
          maxLength={500}
          placeholder={TM.catatanPetunjuk}
          value={catatan}
          onChange={(e) => {
            setCatatan(e.target.value)
          }}
        />
      </label>
    </Modal>
  )
}

/** Popup Riwayat versi: unduh versi mana pun, aktifkan versi lama atau berkas bawaan. */
function DialogRiwayat({ slot, onTutup, onBerubah }: { slot: Slot; onTutup: () => void; onBerubah: (pesan: string) => void }) {
  const [riwayat, setRiwayat] = useState<Versi[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)

  const muat = useCallback(() => {
    ambilRiwayat(slot.kode).then(
      (r) => {
        setRiwayat(r.versi)
      },
      (g: unknown) => {
        setGalat(g)
      },
    )
  }, [slot.kode])
  useEffect(muat, [muat])

  const jalankan = (versi: number) => {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    aktifkan(slot.kode, versi).then(
      () => {
        setSibuk(false)
        onBerubah(TM.diaktifkan(versi))
        muat()
      },
      (g: unknown) => {
        setSibuk(false)
        setGalat(g)
      },
    )
  }
  const unduh = (versi: number) => {
    unduhTemplat(slot.kode, slot.namaUnduhan.replace(/(\.[^.]+)$/, ` v${versi}$1`), versi).catch((g: unknown) => {
      setGalat(g)
    })
  }

  const bawaanAktif = riwayat !== null && !riwayat.some((v) => v.aktif)
  return (
    <Modal judul={`${TM.judulRiwayat} · ${slot.grup} · ${slot.nama}`} onTutup={onTutup} labelBatal={TM.tutup} lebar>
      {galat !== null && <Gagal galat={galat} />}
      {riwayat === null && galat === null && <Memuat pesan={TM.memuatRiwayat} />}
      {riwayat !== null && (
        <table className="inbox__tabel templat-manager__tabel">
          <thead>
            <tr>
              <th>{TM.kolomVersi}</th>
              <th>{TM.kolomBerkas}</th>
              <th>{TM.kolomUnggah}</th>
              <th>{TM.kolomCatatan}</th>
              <th className="table__actions">{TM.kolomAksi}</th>
            </tr>
          </thead>
          <tbody>
            {riwayat.map((v) => (
              <tr key={v.versi} className="inbox__baris">
                <td>v{v.versi}</td>
                <td>
                  {v.namaBerkas}
                  <span className="muted templat-manager__kecil">
                    {ukuranBerkas(v.ukuran)}
                    {v.jumlahKolom > 0 ? ` · ${v.jumlahKolom} col` : ''}
                  </span>
                </td>
                <td>
                  {v.diunggahOleh}
                  <span className="muted templat-manager__kecil">{v.tglUnggah}</span>
                </td>
                <td>{v.catatan === '' ? '—' : v.catatan}</td>
                <td className="table__actions">
                  <span className="templat-manager__aksi">
                    {v.aktif ? (
                      <span className="badge templat-manager__badge--aktif">{TM.aktif}</span>
                    ) : (
                      <button type="button" className="btn btn--ghost btn--sm" disabled={sibuk} onClick={() => jalankan(v.versi)}>
                        {TM.aktifkan}
                      </button>
                    )}
                    <button type="button" className="btn btn--ghost btn--sm" onClick={() => unduh(v.versi)}>
                      {TM.unduh}
                    </button>
                  </span>
                </td>
              </tr>
            ))}
            <tr className="inbox__baris">
              <td>{TM.bawaan}</td>
              <td>{slot.namaUnduhan}</td>
              <td className="muted">{TM.aplikasi}</td>
              <td>—</td>
              <td className="table__actions">
                <span className="templat-manager__aksi">
                  {bawaanAktif ? (
                    <span className="badge templat-manager__badge--aktif">{TM.aktif}</span>
                  ) : (
                    <button type="button" className="btn btn--ghost btn--sm" disabled={sibuk} onClick={() => jalankan(0)}>
                      {TM.aktifkan}
                    </button>
                  )}
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => unduh(0)}>
                    {TM.unduh}
                  </button>
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      )}
    </Modal>
  )
}

export default function TemplateManager() {
  const [daftar, setDaftar] = useState<Slot[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [kata, setKata] = useState('')
  const [grup, setGrup] = useState('')
  const [pesan, setPesan] = useState<string | null>(null)
  const [unggahUntuk, setUnggahUntuk] = useState<Slot | null>(null)
  const [riwayatUntuk, setRiwayatUntuk] = useState<Slot | null>(null)

  const muat = useCallback(() => {
    ambilDaftar().then(
      (d) => {
        setDaftar(d.daftar)
        setGalat(null)
      },
      (g: unknown) => {
        setGalat(g)
      },
    )
  }, [])
  useEffect(muat, [muat])

  const semuaGrup = daftar === null ? [] : kelompokSlot(daftar).map((g) => g.grup)
  const tersaring = daftar === null ? [] : saringSlot(daftar, kata, grup)

  return (
    <section className="inbox templat-manager">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{TM.judul}</h2>
      </header>
      <div className="toolbar">
        <input
          className="field__input templat-manager__cari"
          type="search"
          aria-label={TM.cari}
          placeholder={TM.cari}
          value={kata}
          onChange={(e) => {
            setKata(e.target.value)
          }}
        />
        <select
          className="field__input templat-manager__grup"
          aria-label={TM.semuaMenu}
          value={grup}
          onChange={(e) => {
            setGrup(e.target.value)
          }}
        >
          <option value="">{TM.semuaMenu}</option>
          {semuaGrup.map((g) => (
            <option key={g} value={g}>
              {g}
            </option>
          ))}
        </select>
        <span className="toolbar__spacer" />
        {daftar !== null && <span className="muted">{TM.jumlah(tersaring.length)}</span>}
      </div>

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={TM.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={TM.kosong} />}
      {daftar !== null && daftar.length > 0 && tersaring.length === 0 && <Kosong pesan={TM.tidakCocok} />}

      {kelompokSlot(tersaring).map((g) => (
        <section key={g.grup} className="templat-manager__grup-kartu">
          <h3 className="templat-manager__judul-grup">
            {g.grup} <span className="muted">({TM.jumlah(g.slot.length)})</span>
          </h3>
          <table className="inbox__tabel templat-manager__tabel">
            <thead>
              <tr>
                <th>{TM.kolomNama}</th>
                <th>{TM.kolomBerkas}</th>
                <th>{TM.kolomKolom}</th>
                <th>{TM.kolomVersi}</th>
                <th>{TM.kolomUnggah}</th>
                <th className="table__actions">{TM.kolomAksi}</th>
              </tr>
            </thead>
            <tbody>
              {g.slot.map((s) => (
                <tr key={s.kode} className="inbox__baris">
                  <td>
                    {s.nama}
                    <span className="muted templat-manager__kecil">{s.dipakaiDi}</span>
                  </td>
                  <td>{s.aktif?.namaBerkas ?? s.namaUnduhan}</td>
                  <td>{s.jumlahKolom > 0 ? s.jumlahKolom : '—'}</td>
                  <td>{s.aktif === null ? <span className="badge">{TM.bawaan}</span> : `v${s.aktif.versi}`}</td>
                  <td>
                    {s.aktif === null ? (
                      '—'
                    ) : (
                      <>
                        {s.aktif.diunggahOleh}
                        <span className="muted templat-manager__kecil">{s.aktif.tglUnggah}</span>
                      </>
                    )}
                  </td>
                  <td className="table__actions">
                    <span className="templat-manager__aksi">
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => {
                          unduhTemplat(s.kode, s.namaUnduhan).catch((e: unknown) => {
                            setGalat(e)
                          })
                        }}
                      >
                        {TM.unduh}
                      </button>
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => {
                          setPesan(null)
                          setUnggahUntuk(s)
                        }}
                      >
                        {TM.unggah}
                      </button>
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => {
                          setPesan(null)
                          setRiwayatUntuk(s)
                        }}
                      >
                        {TM.riwayat}
                      </button>
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      ))}

      {unggahUntuk !== null && (
        <DialogUnggah
          slot={unggahUntuk}
          onTutup={() => {
            setUnggahUntuk(null)
          }}
          onTersimpan={(p) => {
            setUnggahUntuk(null)
            setPesan(p)
            muat()
          }}
        />
      )}
      {riwayatUntuk !== null && (
        <DialogRiwayat
          slot={riwayatUntuk}
          onTutup={() => {
            setRiwayatUntuk(null)
            muat()
          }}
          onBerubah={(p) => {
            setPesan(p)
          }}
        />
      )}
    </section>
  )
}
