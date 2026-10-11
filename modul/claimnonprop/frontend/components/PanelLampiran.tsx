// Disalin dari `modul/claimprop/frontend/components/PanelLampiran.tsx` (pola, bukan impor; perintah work owner 09-10-2026 "untuk attachment
// juga mengikuti dari klaim prop"): keputusan work owner yang disebut di bawah adalah keputusan layar Claim Prop.
// Tab "Lampiran" - susunan layar Pega (screenshot work owner 09-10-2026): judul (tanpa jumlah - work owner 09-10-2026),
// tombol Add attachment / Refresh / Save, grid Category - Count Attach - Upload File - View File (5 baris, "Show All"). Isinya
// `AttachCategory.pxResults` = master T_KATEGORI_DOC_KLAIM NONPROP x cacah dokumen klaim. Upload File / Add attachment =
// `GCNMSaveAttachments` (kategori baris = TempInputParam.pyCategory; Add attachment = kategori PER BERKAS `.pyCategory`,
// tabel File / Category + hapus, screenshot work owner 09-10-2026), berkas dipilih atau diseret-lepas
// (`dragDropFileUpload`), satu `InsertDocument_Act` per berkas. View File = dokumen klaim kategori itu (kolom Section `ViewAttachment`), isinya lewat
// `GetBase64Attachment` -> `GetUrlGoogleStorage_Act`. Hanya pemegang assignment yang mengunggah.

import { useCallback, useEffect, useRef, useState } from 'react'

import { Gagal, Memuat, Modal, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { ApiFailure } from '../../../../inti/frontend/klien'
import {
  ambilIsiLampiran,
  ambilLampiran,
  hapusLampiran,
  pindahKategoriLampiran,
  tautanOfficeLampiran,
  unduhLampiran,
  unggahLampiran,
  type KategoriLampiran,
  type Lampiran,
  type LampiranKasus,
} from '../api'
import { CNP } from '../labels'
import {
  BINGKAI_PENAMPIL,
  PARAM_PENAMPIL,
  PENAMPIL_OFFICE,
  bisaViewOffice,
  jenisViewOnline,
  mimeViewOnline,
} from '../penampilOffice'

/** Baris grid sebelum "Show All" (screenshot: "Showing 5 of 18"). */
const BARIS_AWAL = 5

/** Kategori bawaan berkas Add attachment (screenshot work owner: berkas baru berkategori "File"). */
const KATEGORI_BAWAAN = 'File'

/** Satu berkas di jendela unggah (`dragDropFileUpload.pxResults`): isi dan `.pyCategory`. */
interface BerkasPilihan {
  file: File
  kategori: string
}

const TAB_LAMPIRAN = ['lampiran'] as const

function IkonUnggah() {
  return (
    <svg
      width="20"
      height="20"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      aria-hidden="true"
    >
      <path d="M12 16V4M7 9l5-5 5 5M4 15v4h16v-4" />
    </svg>
  )
}

function IkonHapus() {
  return (
    <svg viewBox="0 0 24 24" width="15" height="15" aria-hidden="true" focusable="false">
      <path
        d="M4 7h16M10 11v6M14 11v6M6 7l1 13h10l1-13M9 7V4h6v3"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

function IkonLihat() {
  return (
    <svg
      width="20"
      height="20"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      aria-hidden="true"
    >
      <path d="M12 4v12M7 11l5 5 5-5M4 15v4h16v-4" />
    </svg>
  )
}

export default function PanelLampiran({
  id,
  bolehSimpan,
  onSimpan,
  hanyaLihat = false,
}: {
  id: string
  /** Tampilan saja: tanpa Upload / Delete walau server mengizinkan (View more details Komite Claim Prop). */
  hanyaLihat?: boolean
  /** Tombol Save = simpan kasus (pemegang assignment). */
  bolehSimpan: boolean
  onSimpan: () => void
}) {
  const [data, setData] = useState<LampiranKasus | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [semua, setSemua] = useState(false)
  // unggah: kategori tetap (Upload File baris) atau pilihan (Add attachment, kategori '')
  const [unggah, setUnggah] = useState<{ tetap: boolean; kategori: string } | null>(null)
  const [lihat, setLihat] = useState<KategoriLampiran | null>(null)
  const [berkas, setBerkas] = useState<BerkasPilihan[]>([])
  const [seret, setSeret] = useState(false)
  const [sibuk, setSibuk] = useState(false)
  const masukan = useRef<HTMLInputElement>(null)

  const muat = useCallback(() => {
    ambilLampiran(id).then(
      (d) => {
        setData(d)
        setGalat(null)
      },
      (g: unknown) => setGalat(g),
    )
  }, [id])

  useEffect(() => {
    muat()
  }, [muat])

  // Daftar dari server bisa `null` (backend lama / master kosong) - layar tidak boleh mati karenanya.
  const daftarKategori = data?.kategori ?? []
  const daftarLampiran = data?.lampiran ?? []
  const bolehUnggah = !hanyaLihat && data?.bolehUnggah === true
  const tampil = semua ? daftarKategori : daftarKategori.slice(0, BARIS_AWAL)
  const labelKategori = (k: string) => daftarKategori.find((x) => x.id === k)?.label ?? k

  const tutupUnggah = () => {
    setUnggah(null)
    setBerkas([])
    setSeret(false)
  }

  // Berkas dipilih / diseret-lepas DITAMBAHKAN ke daftar (multi attach); kategori bawaan "File" bila ada di master.
  const tambahBerkas = (fs: FileList | null) => {
    const bawaan = daftarKategori.some((k) => k.id === KATEGORI_BAWAAN) ? KATEGORI_BAWAAN : ''
    const baru = Array.from(fs ?? []).map((file) => ({ file, kategori: bawaan }))
    setBerkas((b) => [...b, ...baru])
    if (masukan.current) masukan.current.value = ''
  }

  const siapKirim =
    unggah !== null &&
    berkas.length > 0 &&
    (unggah.tetap ? unggah.kategori !== '' : berkas.every((b) => b.kategori !== ''))

  const kirim = () => {
    if (!unggah || !siapKirim) return
    setSibuk(true)
    unggahLampiran(
      id,
      unggah.tetap ? unggah.kategori : '',
      berkas.map((b) => b.file),
      unggah.tetap ? undefined : berkas.map((b) => b.kategori),
    ).then(
      () => {
        setSibuk(false)
        tutupUnggah()
        muat()
      },
      (g: unknown) => {
        setSibuk(false)
        setGalat(g)
        tutupUnggah()
        muat()
      },
    )
  }

  const pesanGalat = galat instanceof ApiFailure ? galat.detail.message : null
  return (
    <div>
      <StripTab tab={TAB_LAMPIRAN} aktif="lampiran" onPilih={() => undefined} label={() => CNP.lampiran} />
      <section className="panel">
        <h3 className="panel__title">{CNP.lampiran}</h3>
        {galat !== null &&
          (pesanGalat ? <div className="alert alert--error">{pesanGalat}</div> : <Gagal galat={galat} />)}
        <div className="claimnonprop__lampiran-aksi">
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={!bolehUnggah}
            onClick={() => setUnggah({ tetap: false, kategori: '' })}
          >
            {CNP.lampiranTambah}
          </button>
          <button type="button" className="btn btn--ghost btn--sm" onClick={muat}>
            {CNP.lampiranMuatUlang}
          </button>
          <button type="button" className="btn btn--ghost btn--sm" disabled={!bolehSimpan} onClick={onSimpan}>
            {CNP.save}
          </button>
        </div>
        <div className="claimnonprop__grid">
          <div className="claimnonprop__grid-alat">
            <span>
              {CNP.lampiranMenampilkan} {tampil.length} {CNP.dari} {daftarKategori.length}
            </span>
            {!semua && daftarKategori.length > BARIS_AWAL && (
              <button type="button" className="btn btn--ghost btn--sm" onClick={() => setSemua(true)}>
                {CNP.lampiranSemua}
              </button>
            )}
          </div>
          <table className="claimnonprop__tabel">
            <thead>
              <tr>
                <th>{CNP.lampiranKategori}</th>
                <th>{CNP.lampiranCacah}</th>
                <th className="claimnonprop__kolom-ikon">{CNP.lampiranUnggah}</th>
                <th className="claimnonprop__kolom-ikon">{CNP.lampiranLihat}</th>
              </tr>
            </thead>
            <tbody>
              {tampil.map((k) => (
                <tr key={k.id}>
                  <td>{k.label}</td>
                  <td>{k.countAttach}</td>
                  <td className="claimnonprop__kolom-ikon">
                    <button
                      type="button"
                      className="btn btn--sm claimnonprop__ikon-lampiran claimnonprop__ikon-lampiran--unggah"
                      aria-label={`${CNP.lampiranUnggah} ${k.label}`}
                      disabled={!bolehUnggah}
                      onClick={() => setUnggah({ tetap: true, kategori: k.id })}
                    >
                      <IkonUnggah />
                    </button>
                  </td>
                  <td className="claimnonprop__kolom-ikon">
                    <button
                      type="button"
                      className="btn btn--sm claimnonprop__ikon-lampiran claimnonprop__ikon-lampiran--lihat"
                      aria-label={`${CNP.lampiranLihat} ${k.label}`}
                      onClick={() => setLihat(k)}
                    >
                      <IkonLihat />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {unggah && (
        <Modal
          judul={unggah.tetap ? `${CNP.lampiranUnggah} - ${labelKategori(unggah.kategori)}` : CNP.lampiranTambah}
          onTutup={tutupUnggah}
          labelBatal={CNP.cancel}
          aksi={
            <button type="button" className="btn btn--primary" disabled={sibuk || !siapKirim} onClick={kirim}>
              {sibuk ? CNP.lampiranMengirim : CNP.lampiranKirim}
            </button>
          }
        >
          <div
            className={'claimnonprop__seret' + (seret ? ' claimnonprop__seret--aktif' : '')}
            onDragOver={(e) => {
              e.preventDefault()
              setSeret(true)
            }}
            onDragLeave={() => setSeret(false)}
            onDrop={(e) => {
              e.preventDefault()
              setSeret(false)
              if (!sibuk) tambahBerkas(e.dataTransfer.files)
            }}
          >
            {CNP.lampiranSeret}{' '}
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={sibuk}
              onClick={() => masukan.current?.click()}
            >
              {CNP.lampiranPilih}
            </button>
            <input
              ref={masukan}
              type="file"
              multiple
              hidden
              aria-label={CNP.lampiranPilih}
              onChange={(e) => tambahBerkas(e.target.files)}
            />
          </div>
          {berkas.length > 0 && (
            <table className="claimnonprop__tabel">
              <thead>
                <tr>
                  <th>{CNP.lampiranFile}</th>
                  <th>{CNP.lampiranKategori}</th>
                  <th className="claimnonprop__th-aksi" />
                </tr>
              </thead>
              <tbody>
                {berkas.map((b, i) => (
                  <tr key={`${b.file.name}-${String(i)}`}>
                    <td>{b.file.name}</td>
                    <td>
                      {unggah.tetap ? (
                        labelKategori(unggah.kategori)
                      ) : (
                        <select
                          className="field__input"
                          aria-label={`${CNP.lampiranKategori} ${b.file.name}`}
                          value={b.kategori}
                          disabled={sibuk}
                          onChange={(e) => {
                            const v = e.target.value
                            setBerkas((xs) => xs.map((x, j) => (j === i ? { ...x, kategori: v } : x)))
                          }}
                        >
                          <option value="">{CNP.pilih}</option>
                          {daftarKategori.map((k) => (
                            <option key={k.id} value={k.id}>
                              {k.label}
                            </option>
                          ))}
                        </select>
                      )}
                    </td>
                    <td className="claimnonprop__td-aksi">
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm claimnonprop__ikon claimnonprop__ikon--hapus"
                        aria-label={`${CNP.lampiranBuang} ${b.file.name}`}
                        disabled={sibuk}
                        onClick={() => setBerkas((xs) => xs.filter((_, j) => j !== i))}
                      >
                        <IkonHapus />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </Modal>
      )}

      {lihat && (
        <DialogLihat
          id={id}
          awal={lihat.id}
          kategori={daftarKategori}
          lampiran={daftarLampiran}
          bolehUbah={bolehUnggah}
          onUbah={muat}
          onTutup={() => setLihat(null)}
        />
      )}
    </div>
  )
}

/**
 * View File - susunan layar Pega (screenshot work owner 09-10-2026): pilih Category ("Label - (jumlah)"), tabel File /
 * Category / Upload Date / Attached By / View / centang Action / Delete, tombol Download Selected, Delete Selected,
 * Change Category. View (pdf / gambar) dan View Office Online membuka POPUP PENAMPIL layar penuh yang menggantikan
 * daftar (pola NB Treaty In). Delete / Change Category hanya pemegang assignment. Daftar dibaca dari panel (`onUbah`
 * memuat ulang sesudah perubahan).
 */
function DialogLihat({
  id,
  awal,
  kategori,
  lampiran,
  bolehUbah,
  onUbah,
  onTutup,
}: {
  id: string
  awal: string
  kategori: KategoriLampiran[]
  lampiran: Lampiran[]
  bolehUbah: boolean
  onUbah: () => void
  onTutup: () => void
}) {
  const [kat, setKat] = useState(awal)
  const [pilih, setPilih] = useState<ReadonlySet<string>>(new Set())
  const [tujuan, setTujuan] = useState<string | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [penampil, setPenampil] = useState<{ nama: string; jenis: 'office' | 'pdf' | 'gambar'; siap: boolean } | null>(
    null,
  )
  const formOffice = useRef<HTMLFormElement>(null)
  const urlOffice = useRef<HTMLInputElement>(null)
  const bingkai = useRef<HTMLIFrameElement>(null)
  const gambar = useRef<HTMLImageElement>(null)
  // Objek URL yang sedang tampil (dicabut saat penampil ditutup) dan giliran buka (jawaban penampil tertutup dibuang).
  const objekAktif = useRef<string | null>(null)
  const giliran = useRef(0)

  useEffect(
    () => () => {
      if (objekAktif.current !== null) URL.revokeObjectURL(objekAktif.current)
    },
    [],
  )

  const daftar = lampiran.filter((a) => a.kategori === kat)
  const terpilih = daftar.filter((a) => pilih.has(a.id))
  const labelKat = (k: string) => kategori.find((x) => x.id === k)?.label ?? k
  const semuaTerpilih = daftar.length > 0 && terpilih.length === daftar.length

  const gantiKategori = (k: string) => {
    setKat(k)
    setPilih(new Set())
    setTujuan(null)
  }

  const centang = (lid: string, ya: boolean) =>
    setPilih((s) => {
      const n = new Set(s)
      if (ya) n.add(lid)
      else n.delete(lid)
      return n
    })

  const jalankan = (kerja: () => Promise<void>) => {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    kerja()
      .catch((g: unknown) => setGalat(g))
      .finally(() => setSibuk(false))
  }

  // Delete / Delete Selected / Change Category: berurutan; sesudahnya panel memuat ulang (berhasil sebagian ikut tampil).
  const ubah = (kerja: () => Promise<void>) =>
    jalankan(async () => {
      try {
        await kerja()
        setPilih(new Set())
        setTujuan(null)
      } finally {
        onUbah()
      }
    })

  const tutupPenampil = () => {
    giliran.current++
    if (objekAktif.current !== null) URL.revokeObjectURL(objekAktif.current)
    objekAktif.current = null
    setPenampil(null)
  }

  /** View Office Online: URL bertanda tangan -> form GET ke bingkai popup penampil. */
  const bukaOffice = (a: Lampiran) => {
    if (sibuk) return
    const ke = ++giliran.current
    setPenampil({ nama: a.namaFile, jenis: 'office', siap: false })
    jalankan(async () => {
      try {
        const { url } = await tautanOfficeLampiran(id, a)
        if (ke !== giliran.current) return
        if (formOffice.current !== null && urlOffice.current !== null) {
          urlOffice.current.value = url
          formOffice.current.submit()
        }
        setPenampil((p) => (p === null ? p : { ...p, siap: true }))
      } catch (g: unknown) {
        if (ke === giliran.current) tutupPenampil()
        throw g
      }
    })
  }

  /** View pdf / gambar: isi dari rute unduh beridentitas -> objek URL lokal di popup penampil. */
  const bukaOnline = (a: Lampiran) => {
    const jenis = jenisViewOnline(a.mime)
    if (sibuk || jenis === null) return
    const ke = ++giliran.current
    setPenampil({ nama: a.namaFile, jenis, siap: false })
    jalankan(async () => {
      try {
        const isi = await ambilIsiLampiran(id, a)
        if (ke !== giliran.current) return
        const objekURL = URL.createObjectURL(new Blob([isi], { type: mimeViewOnline(a.mime) }))
        objekAktif.current = objekURL
        if (jenis === 'gambar' && gambar.current !== null) gambar.current.src = objekURL
        if (jenis === 'pdf' && bingkai.current !== null) bingkai.current.src = objekURL
        setPenampil((p) => (p === null ? p : { ...p, siap: true }))
      } catch (g: unknown) {
        if (ke === giliran.current) tutupPenampil()
        throw g
      }
    })
  }

  // ⛔ `key` BERBEDA untuk popup penampil dan popup daftar (bug layar beku lampiran Bordereaux 08-10-2026).
  if (penampil !== null) {
    return (
      <Modal key="penampil" judul={penampil.nama} onTutup={tutupPenampil} labelBatal={CNP.tutup} penuh>
        {!penampil.siap && <Memuat />}
        <form ref={formOffice} method="get" action={PENAMPIL_OFFICE} target={BINGKAI_PENAMPIL} hidden>
          <input ref={urlOffice} type="hidden" name={PARAM_PENAMPIL} />
        </form>
        {penampil.jenis === 'gambar' ? (
          <img ref={gambar} alt={penampil.nama} className="claimnonprop__penampil-gambar" />
        ) : (
          <iframe
            ref={bingkai}
            name={BINGKAI_PENAMPIL}
            title={penampil.nama}
            className="claimnonprop__penampil-bingkai"
          />
        )}
      </Modal>
    )
  }

  const pesanGalat = galat instanceof ApiFailure ? galat.detail.message : null
  return (
    <Modal key="daftar" judul={CNP.lampiranLihat} onTutup={onTutup} labelBatal={CNP.tutup} lebar>
      <label>
        {CNP.lampiranKategori}{' '}
        <select className="field__input" value={kat} disabled={sibuk} onChange={(e) => gantiKategori(e.target.value)}>
          {kategori.map((k) => (
            <option key={k.id} value={k.id}>
              {k.label} - ({k.countAttach})
            </option>
          ))}
        </select>
      </label>
      {galat !== null &&
        (pesanGalat ? <div className="alert alert--error">{pesanGalat}</div> : <Gagal galat={galat} />)}
      {daftar.length === 0 ? (
        <p className="muted">{CNP.lampiranKosong}</p>
      ) : (
        <table className="claimnonprop__tabel">
          <thead>
            <tr>
              <th>{CNP.lampiranNama}</th>
              <th>{CNP.lampiranKategori}</th>
              <th>{CNP.lampiranTanggal}</th>
              <th>{CNP.lampiranOleh}</th>
              <th aria-label={CNP.lampiranView} />
              <th>
                <label>
                  <input
                    type="checkbox"
                    checked={semuaTerpilih}
                    disabled={sibuk}
                    onChange={(e) => setPilih(e.target.checked ? new Set(daftar.map((a) => a.id)) : new Set())}
                  />{' '}
                  {CNP.lampiranAksi}
                </label>
              </th>
              <th>{CNP.lampiranHapus}</th>
            </tr>
          </thead>
          <tbody>
            {daftar.map((a) => (
              <tr key={a.id}>
                <td>
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    disabled={sibuk}
                    onClick={() => jalankan(() => unduhLampiran(id, a))}
                  >
                    {a.namaFile}
                  </button>
                </td>
                <td>{labelKat(a.kategori)}</td>
                <td>{a.tanggal}</td>
                <td>{a.operator}</td>
                <td>
                  <div className="claimnonprop__lampiran-aksi">
                    {a.adaObjek && jenisViewOnline(a.mime) !== null && (
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        disabled={sibuk}
                        onClick={() => bukaOnline(a)}
                      >
                        {CNP.lampiranView}
                      </button>
                    )}
                    {a.adaObjek && bisaViewOffice(a.mime) && (
                      <button
                        type="button"
                        className="btn btn--ghost btn--sm"
                        disabled={sibuk}
                        onClick={() => bukaOffice(a)}
                      >
                        {CNP.lampiranViewOffice}
                      </button>
                    )}
                  </div>
                </td>
                <td>
                  <input
                    type="checkbox"
                    aria-label={`${CNP.lampiranAksi} ${a.namaFile}`}
                    checked={pilih.has(a.id)}
                    disabled={sibuk}
                    onChange={(e) => centang(a.id, e.target.checked)}
                  />
                </td>
                <td>
                  {bolehUbah && (
                    <button
                      type="button"
                      className="btn btn--danger btn--sm"
                      disabled={sibuk}
                      onClick={() => ubah(async () => void (await hapusLampiran(id, a)))}
                    >
                      <IkonHapus /> {CNP.lampiranHapus}
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      <div className="claimnonprop__grid-alat">
        <button
          type="button"
          className="btn btn--ghost btn--sm"
          disabled={sibuk || terpilih.length === 0}
          onClick={() =>
            jalankan(async () => {
              for (const a of terpilih) await unduhLampiran(id, a)
            })
          }
        >
          {CNP.lampiranUnduhPilih}
        </button>
        {bolehUbah && (
          <>
            <button
              type="button"
              className="btn btn--danger btn--sm"
              disabled={sibuk || terpilih.length === 0}
              onClick={() =>
                ubah(async () => {
                  for (const a of terpilih) await hapusLampiran(id, a)
                })
              }
            >
              <IkonHapus /> {CNP.lampiranHapusPilih}
            </button>
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              disabled={sibuk || terpilih.length === 0}
              onClick={() => setTujuan(tujuan === null ? '' : null)}
            >
              {CNP.lampiranPindah}
            </button>
          </>
        )}
      </div>
      {bolehUbah && tujuan !== null && terpilih.length > 0 && (
        <div className="claimnonprop__grid-alat">
          <label>
            {CNP.lampiranPindahKe}{' '}
            <select
              className="field__input"
              value={tujuan}
              disabled={sibuk}
              onChange={(e) => setTujuan(e.target.value)}
            >
              <option value="">{CNP.pilih}</option>
              {kategori
                .filter((k) => k.id !== kat)
                .map((k) => (
                  <option key={k.id} value={k.id}>
                    {k.label}
                  </option>
                ))}
            </select>
          </label>
          <button
            type="button"
            className="btn btn--primary btn--sm"
            disabled={sibuk || tujuan === ''}
            onClick={() =>
              ubah(async () => {
                await pindahKategoriLampiran(
                  id,
                  tujuan,
                  terpilih.map((a) => a.id),
                )
              })
            }
          >
            {CNP.lampiranPindah}
          </button>
        </div>
      )}
    </Modal>
  )
}
