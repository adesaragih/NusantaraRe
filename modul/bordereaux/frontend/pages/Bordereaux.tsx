// Halaman Bordereaux (perintah work owner 04-10-2026) - padanan harness `PortalBordereaux` / `PortalBordereaux_Sec`:
// tombol Input Data (hanya menu ber-hak PENUH atau superadmin), tombol Copy Old Data (hanya superadmin),
// kotak Search / Filter, grid daftar (Edit, View, Delete menurut hak), dan form `InputBordereaux` yang menggantikan
// daftar saat dibuka.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilDaftar, ambilPilihan, hapus, type BarisDaftar, type Filter, type Halaman, type Pilihan } from '../api'
import { FILTER_KOSONG, jumlahHalaman, keIso, keKabel } from '../aturan'
import ChartBordereaux from '../components/ChartBordereaux'
import DialogCopyOld from '../components/DialogCopyOld'
import FormBordereaux from '../components/FormBordereaux'
import TabelDaftar from '../components/TabelDaftar'
import { BDX } from '../labels'

type Form = { id: string | null; lihat: boolean } | null

export default function Bordereaux() {
  const [pilihan, setPilihan] = useState<Pilihan | null>(null)
  const [saring, setSaring] = useState(false)
  const [ketik, setKetik] = useState<Filter>(FILTER_KOSONG)
  const [filter, setFilter] = useState<Filter>(FILTER_KOSONG)
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<Halaman | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [segar, setSegar] = useState(0)
  const [form, setForm] = useState<Form>(null)
  const [akanHapus, setAkanHapus] = useState<BarisDaftar | null>(null)
  const [menghapus, setMenghapus] = useState(false)
  const [galatHapus, setGalatHapus] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)
  const [copyOld, setCopyOld] = useState(false)

  useEffect(() => {
    ambilPilihan().then(setPilihan, (g: unknown) => setGalat(g))
  }, [])

  const muat = useCallback((f: Filter, h: number) => {
    ambilDaftar(f, h).then(
      (d) => {
        setData(d)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
  }, [])
  useEffect(() => {
    if (form === null) muat(filter, halaman)
  }, [muat, filter, halaman, segar, form])

  if (form !== null && pilihan !== null) {
    return (
      <section className="inbox bordereaux__akar">
        <FormBordereaux
          id={form.id}
          lihat={form.lihat}
          pilihan={pilihan}
          onTutup={(berubah) => {
            setForm(null)
            if (berubah) setSegar((n) => n + 1)
          }}
        />
      </section>
    )
  }

  const jalankanHapus = () => {
    if (akanHapus === null || menghapus) return
    setMenghapus(true)
    setGalatHapus(null)
    hapus(akanHapus.bdxId).then(
      () => {
        setMenghapus(false)
        setPesan(BDX.terhapus(akanHapus.bdxId))
        setAkanHapus(null)
        setSegar((n) => n + 1)
      },
      (g: unknown) => {
        setMenghapus(false)
        setGalatHapus(g)
      },
    )
  }

  const ubah = (s: Partial<Filter>) => setKetik((f) => ({ ...f, ...s }))
  const teks = (nama: keyof Filter, label: string) => (
    <label className="field">
      <span className="field__label">{label}</span>
      <input className="field__input" value={ketik[nama]} onChange={(e) => ubah({ [nama]: e.target.value })} />
    </label>
  )
  const total = data?.total ?? 0
  const jumlah = jumlahHalaman(total, data?.ukuran ?? 1)

  return (
    <section className="inbox bordereaux__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{BDX.judul}</h2>
        <span className="toolbar__spacer" />
        {pilihan?.copyOld === true && (
          <button
            type="button"
            className="btn btn--ghost"
            onClick={() => {
              setPesan(null)
              setCopyOld(true)
            }}
          >
            {BDX.copyOld}
          </button>
        )}
        {pilihan?.bolehBuat === true && (
          <button
            type="button"
            className="btn btn--primary"
            onClick={() => {
              setPesan(null)
              setForm({ id: null, lihat: false })
            }}
          >
            {BDX.inputData}
          </button>
        )}
      </header>

      <ChartBordereaux segar={segar} />

      <label className="bordereaux__saring-centang">
        <input type="checkbox" checked={saring} onChange={(e) => setSaring(e.target.checked)} /> {BDX.searchFilter}
      </label>
      {saring && pilihan !== null && (
        <form
          className="bordereaux__kartu bordereaux__filter"
          onSubmit={(e) => {
            e.preventDefault()
            setPesan(null)
            setHalaman(1)
            setFilter(ketik)
          }}
        >
          {teks('id', BDX.id)}
          <label className="field">
            <span className="field__label">{BDX.type}</span>
            <select className="field__input" value={ketik.type} onChange={(e) => ubah({ type: e.target.value, business: '' })}>
              <option value="">{BDX.semua}</option>
              {pilihan.type.map((t) => (
                <option key={t}>{t}</option>
              ))}
            </select>
          </label>
          <label className="field">
            <span className="field__label">{BDX.filterBusiness}</span>
            <select className="field__input" value={ketik.business} onChange={(e) => ubah({ business: e.target.value })}>
              <option value="">{BDX.semua}</option>
              {(ketik.type === '' ? [...new Set(Object.values(pilihan.business).flat())] : (pilihan.business[ketik.type] ?? [])).map((b) => (
                <option key={b}>{b}</option>
              ))}
            </select>
          </label>
          {teks('reffSoa', BDX.filterReffNoSoa)}
          {teks('reffBdx', BDX.reffNoBdx)}
          {teks('ceding', BDX.cedingCo)}
          {teks('treaty', BDX.treatyName)}
          <label className="field">
            <span className="field__label">{BDX.reportStart}</span>
            <input className="field__input" type="date" value={keIso(ketik.start)} onChange={(e) => ubah({ start: keKabel(e.target.value) })} />
          </label>
          <label className="field">
            <span className="field__label">{BDX.reportEnd}</span>
            <input className="field__input" type="date" value={keIso(ketik.end)} onChange={(e) => ubah({ end: keKabel(e.target.value) })} />
          </label>
          {teks('position', BDX.position)}
          <label className="field">
            <span className="field__label">{BDX.status}</span>
            <select className="field__input" value={ketik.status} onChange={(e) => ubah({ status: e.target.value })}>
              <option value="">{BDX.semua}</option>
              {pilihan.status.map((s) => (
                <option key={s}>{s}</option>
              ))}
            </select>
          </label>
          <div className="bordereaux__filter-aksi">
            <button
              type="button"
              className="btn btn--ghost"
              onClick={() => {
                setKetik(FILTER_KOSONG)
                setFilter(FILTER_KOSONG)
                setHalaman(1)
              }}
            >
              {BDX.reset}
            </button>
            <button type="submit" className="btn btn--primary">
              {BDX.applyFilter}
            </button>
          </div>
        </form>
      )}

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={BDX.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={filter === FILTER_KOSONG ? BDX.kosong : BDX.tidakCocok} />}
      {data !== null && data.daftar.length > 0 && (
        <TabelDaftar
          daftar={data.daftar}
          onBuka={(id, lihat) => setForm({ id, lihat })}
          onHapus={(b) => {
            setPesan(null)
            setGalatHapus(null)
            setAkanHapus(b)
          }}
        />
      )}
      {data !== null && (
        <div className="bordereaux__halaman">
          <span className="muted">{BDX.jumlah(total)}</span>
          <span className="toolbar__spacer" />
          <button type="button" className="btn btn--ghost" disabled={halaman <= 1} onClick={() => setHalaman((h) => h - 1)}>
            {BDX.sebelumnya}
          </button>
          <span>{BDX.halaman(halaman, jumlah)}</span>
          <button type="button" className="btn btn--ghost" disabled={halaman >= jumlah} onClick={() => setHalaman((h) => h + 1)}>
            {BDX.berikutnya}
          </button>
        </div>
      )}

      {copyOld && (
        <DialogCopyOld
          onTutup={(adaYangDisalin) => {
            setCopyOld(false)
            if (adaYangDisalin) setSegar((n) => n + 1)
          }}
        />
      )}

      {akanHapus !== null && (
        <Modal
          judul={BDX.judulHapus}
          onTutup={() => setAkanHapus(null)}
          onKirim={jalankanHapus}
          labelBatal={BDX.batal}
          aksi={
            <button type="submit" className="btn btn--danger" disabled={menghapus}>
              {menghapus ? BDX.menghapus : BDX.hapus}
            </button>
          }
        >
          {galatHapus !== null && <Gagal galat={galatHapus} />}
          <p>{BDX.kalimatHapus(akanHapus.bdxId)}</p>
          <p className="muted">
            {[akanHapus.type, akanHapus.typeBusiness, akanHapus.cedingName, akanHapus.treatyName].filter((s) => s !== '').join(' · ')}
          </p>
        </Modal>
      )}
    </section>
  )
}
