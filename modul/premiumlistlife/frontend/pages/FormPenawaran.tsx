// Form penawaran — tiket 01 bagian 3.
//
// Meniru isian `Section/InputOfferLife.xml` yang DAPAT DIISI, dan kedua popup
// pilihannya `Ceding_Harness` / `PolicyHolder_Harness`. Tombol `Save Offer`
// menyimpan tanpa memindahkan tahap; keputusan Confirm/Decline tetap di
// `InputOffer.tsx`.
//
// ⛔ RALAT 01-10-2026: sel penawaran yang TERISI di layar Pega lama (Age Limit,
// Coverage Period, Sum Insured, Offering/Response/Confirmation Date, Input TBC,
// Max TBC, Status Update, Marketing Note) kini dibawa — migrasi 059. Yang kosong
// di layar itu (Insured Name, Occupation, Underwriting Policy, Re-Confirmation/
// Realization/Binding Date, Final Status) ikut dibawa (migrasi 060) - work owner
// akan menyaring yang tidak perlu.
//
// ⛔ Pilihan tertutup (System Reinsurance, Class of Business, Status) DATANG
// DARI SERVER bersama isiannya — nol daftar di berkas ini.

import { useCallback, useEffect, useState } from 'react'

import { Area, Field, Gagal, Memuat, Modal, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilPenawaranPolis,
  cariCedingPolis,
  cariPemegangPolis,
  isiDariPenawaran,
  simpanPenawaranPolis,
  type BarisRujukanPolis,
  type IsiPenawaranPolis,
  type PenawaranPolis,
  type PilihanKode,
} from '../api'
import { KOLOM_RIWAYAT_PENAWARAN, LABEL_PENAWARAN, TEKS_PILIH } from '../labels'
import AreaTeks from '../components/AreaTeks'
import IsianTanggal from '../components/IsianTanggal'
import { tanggalJamTampil } from '../tanggal'
import '../premiumlistlife.css'

/** Popup yang sedang terbuka. */
export type JenisPopupPenawaran = 'ceding' | 'pemegang'

/** Pilihan server → opsi `Pilih`. */
export function opsiDari(pilihan: PilihanKode[]): Opsi[] {
  return pilihan.map((p) => ({ value: p.kode, label: p.nama }))
}

/**
 * Menerapkan baris popup yang dipilih — `setCeding_act` / `setPolicyHolder_act`.
 *
 * ⛔ Kode DAN nama berpindah BERSAMA: server menolak separuh pasangan.
 */
export function terapkanPilihan(
  isi: IsiPenawaranPolis,
  jenis: JenisPopupPenawaran,
  baris: BarisRujukanPolis,
): IsiPenawaranPolis {
  if (jenis === 'ceding') return { ...isi, cedingCo: baris.id, cedingCoName: baris.nama }
  return { ...isi, policyHolder: baris.id, policyHolderName: baris.nama }
}

/** Tanggal riwayat untuk layar — `dd/mm/yyyy HH:MM` (02-10-2026), tanpa menebak zona. */
export function tanggalRiwayat(iso: string): string {
  return tanggalJamTampil(iso) || '—'
}

export default function FormPenawaran({ polisID }: { polisID: string }) {
  const [data, setData] = useState<PenawaranPolis | null>(null)
  const [isi, setIsi] = useState<IsiPenawaranPolis | null>(null)
  const [galatMuat, setGalatMuat] = useState<unknown>(null)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  const [popup, setPopup] = useState<JenisPopupPenawaran | null>(null)

  const terima = useCallback((p: PenawaranPolis) => {
    setData(p)
    setIsi(isiDariPenawaran(p))
  }, [])

  useEffect(() => {
    let hidup = true
    void (async () => {
      try {
        const p = await ambilPenawaranPolis(polisID)
        if (hidup) terima(p)
      } catch (e) {
        if (hidup) setGalatMuat(e)
      }
    })()
    return () => {
      hidup = false
    }
  }, [polisID, terima])

  async function simpan(): Promise<void> {
    if (isi === null || sibuk || kolomWajibKosong(isi).length > 0) return
    setSibuk(true)
    setGalatSimpan(null)
    setTersimpan(false)
    try {
      terima(await simpanPenawaranPolis(polisID, isi))
      setTersimpan(true)
    } catch (e) {
      setGalatSimpan(e)
    } finally {
      setSibuk(false)
    }
  }

  if (galatMuat !== null) return <Gagal galat={galatMuat} />
  if (data === null || isi === null) return <Memuat />

  const ubah = (medan: keyof IsiPenawaranPolis) => (v: string) => {
    setIsi({ ...isi, [medan]: v })
    setTersimpan(false)
  }
  const bisa = data.bolehDisimpan
  // ⛔ [keputusan work owner 02-10-2026] Kasus Input Premium (FlagOnGoingPolicy
  // '1') TETAP DAPAT MENGUBAH isian penawaran. Di Pega sel-sel ini ber-
  // `pyReadOnlyCondition pyWorkPage.FlagOnGoingPolicy='1'`; penyimpangan ini
  // disengaja. Yang mengunci kini hanya tahapnya (`bolehDisimpan`).
  const kunci = !bisa
  const kurang = kolomWajibKosong(isi)

  return (
    <section className="panel pl-offer">
      <h3 className="panel__title">{LABEL_PENAWARAN.judul}</h3>

      {/*
        DUA KOLOM TETAP (premiumlistlife.css), urutan sel `InputOfferLife.xml`:
        kiri identitas dan pertanggungan, kanan tanggal-tanggal dan catatan.
        Isian pendek berpasangan supaya form tidak menuntut gulir panjang.
      */}
      <div className="pl-offer__kolom-dua">
        <div>
          <h4 className="pl-offer__subjudul">Offer Data</h4>
          {/* `pyVisible NOTBLANK` — hanya tampil bila sudah ada nomornya. */}
          {data.noOffer !== '' && (
            <Field label={LABEL_PENAWARAN.noOffer} value={data.noOffer} onChange={() => {}} readOnly />
          )}
          {/* Satu grid untuk keduanya — kotak isian Ceding dan Policy Holder sama lebar. */}
          <div className={bisa ? 'pl-offer__pilih-grup' : 'pl-offer__pilih-grup pl-offer__pilih-grup--baca'}>
            <Field
              label={LABEL_PENAWARAN.cedingCoName}
              value={isi.cedingCoName}
              onChange={() => {}}
              readOnly
              required
            />
            {bisa && (
              <button type="button" className="btn btn--ghost" onClick={() => { setPopup('ceding') }}>
                {LABEL_PENAWARAN.pilihCeding}
              </button>
            )}
            <Field
              label={LABEL_PENAWARAN.policyHolderName}
              value={isi.policyHolderName}
              onChange={() => {}}
              readOnly
              required
            />
            {bisa && (
              <button type="button" className="btn btn--ghost" onClick={() => { setPopup('pemegang') }}>
                {LABEL_PENAWARAN.pilihPemegang}
              </button>
            )}
          </div>
          <div className="pl-offer__pasangan">
            <Pilih
              kosong={TEKS_PILIH}
              label={LABEL_PENAWARAN.typeCeding}
              value={isi.typeCeding}
              onChange={ubah('typeCeding')}
              opsi={opsiDari(data.pilihan.typeCeding)}
              required
            />
            {/*
              Read-only, TURUNAN System Reinsurance (`SetReinsuranceType`): XOL ->
              Non Proportional, selain itu Proportional. Diperbarui langsung saat
              System Reinsurance diganti; server menghitung ulang saat disimpan.
            */}
            <Field
              label={LABEL_PENAWARAN.jenisAsuransi}
              value={jenisAsuransiDari(isi.typeCeding)}
              onChange={() => {}}
              readOnly
            />
          </div>
          <Pilih
            kosong={TEKS_PILIH}
            label={LABEL_PENAWARAN.businessCode}
            value={isi.businessCode}
            onChange={ubah('businessCode')}
            opsi={opsiDari(data.pilihan.classOfBusiness)}
            required
          />
          {/*
            `Insured Name` (.QQName) dan `Occupation` (.JenisUsaha) DISEMBUNYIKAN —
            keputusan work owner 01-10-2026. Kolomnya (migrasi 060) tetap ada, dan
            nilainya tetap ikut terkirim apa adanya lewat `isi`, sehingga Save Offer
            tidak mengosongkan nilai yang sudah tersimpan.
          */}
          <div className="pl-offer__pasangan">
            <Field
              label={LABEL_PENAWARAN.batasUsiaPeserta}
              type="number"
              value={isi.batasUsiaPeserta}
              onChange={ubah('batasUsiaPeserta')}
              readOnly={kunci}
            />
            <Field
              label={LABEL_PENAWARAN.periodePertanggungan}
              value={isi.periodePertanggungan}
              onChange={ubah('periodePertanggungan')}
              readOnly={kunci}
            />
          </div>
          {/* Uang tetap TEKS — `type="text"`, bukan number (ADR-U-0003). */}
          <Field
            label={LABEL_PENAWARAN.sumInsured}
            value={isi.sumInsured}
            onChange={(v) => { ubah('sumInsured')(saringAngkaDesimal(v)) }}
            readOnly={kunci}
          />
          {/* Text area — terkunci pun tetap text area (03-10-2026). */}
          <AreaTeks
            label={LABEL_PENAWARAN.ketentuanUnderwriting}
            value={isi.ketentuanUnderwriting}
            onChange={ubah('ketentuanUnderwriting')}
            readOnly={kunci}
          />
        </div>
        <div>
          <h4 className="pl-offer__subjudul">Dates &amp; Status</h4>
          <div className="pl-offer__pasangan">
            {(
              [
                ['dateReceived', LABEL_PENAWARAN.dateReceived],
                ['tanggalPenawaran', LABEL_PENAWARAN.tanggalPenawaran],
                ['tanggalRespon', LABEL_PENAWARAN.tanggalRespon],
                ['tanggalKonfirmasi', LABEL_PENAWARAN.tanggalKonfirmasi],
                ['tanggalRealisasi', LABEL_PENAWARAN.tanggalRealisasi],
                ['tanggalKonfirmasiBalik', LABEL_PENAWARAN.tanggalKonfirmasiBalik],
                ['tanggalBind', LABEL_PENAWARAN.tanggalBind],
              ] as const
            ).map(([medan, label]) => (
              <IsianTanggal
                key={medan}
                label={label}
                value={isi[medan]}
                onChange={ubah(medan)}
                readOnly={kunci}
                // Email Received Date wajib — keputusan work owner 02-10-2026.
                required={medan === 'dateReceived'}
              />
            ))}
          </div>
          <div className="pl-offer__pasangan">
            {/* `Input TBC` tampil bila `.TanggalKonfirmasi != ''`. */}
            {tampilTBC(isi) && (
              <Field
                label={LABEL_PENAWARAN.tbc}
                type="number"
                value={isi.tbc}
                onChange={ubah('tbc')}
                readOnly={kunci}
              />
            )}
            {/*
              `Max TBC` tampil bila `.TBC != ''`, read-only. Dihitung LANGSUNG
              dari isian yang sedang diketik (Confirmation Date + Input TBC hari),
              sebagaimana `SetMaxTBCLife_Act` berjalan saat `.TBC` berubah. Saat
              disimpan, server menghitungnya ulang — nilai layar hanya tampilan.
            */}
            {isi.tbc !== '' && (
              <Field
                label={LABEL_PENAWARAN.tanggalTbc}
                value={hitungMaxTBC(isi.tanggalKonfirmasi, isi.tbc)}
                onChange={() => {}}
                readOnly
              />
            )}
          </div>
          <div className="pl-offer__pasangan">
            <Field
              label={LABEL_PENAWARAN.statusUpdate}
              value={isi.statusUpdate}
              onChange={ubah('statusUpdate')}
              readOnly={kunci}
            />
            <Field
              label={LABEL_PENAWARAN.statusFinal}
              value={isi.statusFinal}
              onChange={ubah('statusFinal')}
              readOnly={kunci}
            />
          </div>
          <AreaTeks
            label={LABEL_PENAWARAN.keteranganMarketing}
            value={isi.keteranganMarketing}
            onChange={ubah('keteranganMarketing')}
            readOnly={kunci}
          />
        </div>
      </div>

      <div className="pl-offer__bawah">
        <fieldset className="pl-offer__status">
          <legend className="field__label">
            {LABEL_PENAWARAN.status}
            <span className="field__req">*</span>
          </legend>
          <div className="pl-offer__status-pilihan">
            {data.pilihan.status.map((p) => (
              <label key={p.kode} className="pl-offer__pil">
                <input
                  type="radio"
                  name="status-penawaran"
                  value={p.kode}
                  checked={isi.status === p.kode}
                  disabled={!bisa}
                  onChange={() => { ubah('status')(p.kode) }}
                />
                {p.nama}
              </label>
            ))}
          </div>
        </fieldset>
        <Area
          label={LABEL_PENAWARAN.description}
          value={isi.description}
          onChange={ubah('description')}
          required
        />
      </div>

      {galatSimpan !== null && <Gagal galat={galatSimpan} />}
      {bisa && (
        <div className="pl-offer__aksi">
          {/* ⛔ Tidak dapat ditekan selama ada kolom wajib yang kosong. */}
          <button
            type="button"
            className="btn btn--primary"
            disabled={sibuk || kurang.length > 0}
            onClick={() => { void simpan() }}
          >
            {LABEL_PENAWARAN.simpan}
          </button>
          {kurang.length > 0 && (
            <span className="pl-offer__kurang" role="status">
              Required: {kurang.join(', ')}
            </span>
          )}
          {tersimpan && kurang.length === 0 && (
            <span className="pl-offer__tersimpan" role="status">
              Offer saved.
            </span>
          )}
        </div>
      )}

      <div className="pl-offer__riwayat">
        <h4 className="pl-offer__subjudul">Offer History</h4>
        <div className="pl-offer__riwayat-gulir">
          <table className="inbox__tabel">
            <thead>
              <tr>
                <th>{KOLOM_RIWAYAT_PENAWARAN.dateSuggest}</th>
                <th>{KOLOM_RIWAYAT_PENAWARAN.picSuggest}</th>
                <th>{KOLOM_RIWAYAT_PENAWARAN.isCedingConfirm}</th>
                <th>{KOLOM_RIWAYAT_PENAWARAN.initialSuggest}</th>
                <th>{KOLOM_RIWAYAT_PENAWARAN.commentSuggest}</th>
              </tr>
            </thead>
            <tbody>
              {data.riwayat.length === 0 && (
                <tr>
                  <td colSpan={5} className="pl-offer__kosong">
                    No offer history yet.
                  </td>
                </tr>
              )}
              {data.riwayat.map((r) => (
                <tr key={r.no}>
                  <td>{tanggalRiwayat(r.dateSuggest)}</td>
                  <td>{r.picSuggest || '—'}</td>
                  <td>{r.isCedingConfirm || '—'}</td>
                  <td>{r.initialSuggest || '—'}</td>
                  <td>{r.commentSuggest || '—'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {popup !== null && (
        <PopupRujukan
          jenis={popup}
          onPilih={(b) => {
            setIsi(terapkanPilihan(isi, popup, b))
            setTersimpan(false)
            setPopup(null)
          }}
          onTutup={() => { setPopup(null) }}
        />
      )}
    </section>
  )
}

/** Popup pencarian master — `Ceding_Section` / `PolicyHolder_Section`. */
function PopupRujukan({
  jenis,
  onPilih,
  onTutup,
}: {
  jenis: JenisPopupPenawaran
  onPilih: (b: BarisRujukanPolis) => void
  onTutup: () => void
}) {
  const [cari, setCari] = useState('')
  const [hasil, setHasil] = useState<BarisRujukanPolis[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)

  async function jalankan(): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    try {
      setHasil(jenis === 'ceding' ? await cariCedingPolis(cari) : await cariPemegangPolis(cari))
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  return (
    <Modal
      judul={jenis === 'ceding' ? LABEL_PENAWARAN.pilihCeding : LABEL_PENAWARAN.pilihPemegang}
      onTutup={onTutup}
      onKirim={() => { void jalankan() }}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          {LABEL_PENAWARAN.cari}
        </button>
      }
      lebar
    >
      <Field label={LABEL_PENAWARAN.cari} value={cari} onChange={setCari} autoFocus />
      {galat !== null && <Gagal galat={galat} />}
      {hasil !== null && hasil.length === 0 && <p>No matching results.</p>}
      {hasil !== null && hasil.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{LABEL_PENAWARAN.kolomId}</th>
              <th>{LABEL_PENAWARAN.kolomNama}</th>
              {jenis === 'pemegang' && <th>{LABEL_PENAWARAN.kolomBisnis}</th>}
              <th />
            </tr>
          </thead>
          <tbody>
            {hasil.map((b) => (
              <tr key={b.id}>
                <td>{b.id}</td>
                <td>{b.nama}</td>
                {jenis === 'pemegang' && <td>{b.keterangan ?? ''}</td>}
                <td>
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => { onPilih(b) }}>
                    {LABEL_PENAWARAN.pilih}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </Modal>
  )
}

/** `Input TBC` hanya tampil bila Confirmation Date terisi (`.TanggalKonfirmasi!=''`). */
export function tampilTBC(isi: IsiPenawaranPolis): boolean {
  return isi.tanggalKonfirmasi !== ''
}

/**
 * "Max TBC" — `SetMaxTBCLife_Act` langkah 2:
 * `@DateTime.addCalendar(TanggalKonfirmasi, 0,0,0, TBC, 0,0,0)`.
 *
 * Masukan: Confirmation Date `YYYY-MM-DD` dan Input TBC (teks bilangan bulat).
 * Keluaran: `DD/MM/YYYY` seperti layar Pega (29/09/2026 + 25 = 24/10/2026),
 * atau kosong bila salah satunya kosong / bukan bilangan bulat.
 *
 * ⛔ Dihitung dalam UTC: tanggal tanpa jam yang dihitung di zona setempat dapat
 * bergeser sehari di sekitar pergantian zona waktu.
 */
export function hitungMaxTBC(tanggalKonfirmasi: string, tbc: string): string {
  const cocok = /^(\d{4})-(\d{2})-(\d{2})$/.exec(tanggalKonfirmasi)
  const hari = tbc.trim()
  if (cocok === null || !/^\d+$/.test(hari)) return ''
  const t = new Date(Date.UTC(Number(cocok[1]), Number(cocok[2]) - 1, Number(cocok[3]) + Number(hari)))
  const dua = (n: number) => String(n).padStart(2, '0')
  return `${dua(t.getUTCDate())}/${dua(t.getUTCMonth() + 1)}/${String(t.getUTCFullYear())}`
}

/**
 * "Reinsurance Type" — `SetReinsuranceType`: System Reinsurance "4" (XOL) →
 * Non Proportional, selain itu (termasuk kosong) Proportional. Hanya untuk
 * tampilan langsung; nilai yang disimpan dihitung ulang server.
 */
export function jenisAsuransiDari(typeCeding: string): string {
  return typeCeding === '4' ? 'Non Proportional' : 'Proportional'
}

/**
 * Kolom wajib yang masih kosong — label layar, urut seperti di layar.
 *
 * System Reinsurance, Class of Business, Comment: `pyRequired` korpus.
 * Ceding Name, Policy Holder, Status: keputusan work owner 01-10-2026.
 *
 * ⚠️ Gerbang LAYAR saja, supaya tombol Save Offer tidak dapat ditekan. Server
 * memeriksa ulang dengan aturan yang sama (`models.KekuranganPenawaran`) dan
 * tetap menolak — layar yang dilewati tidak melonggarkan aturannya.
 */
export function kolomWajibKosong(isi: IsiPenawaranPolis): string[] {
  const periksa: [string, string][] = [
    [isi.cedingCoName, LABEL_PENAWARAN.cedingCoName],
    [isi.policyHolderName, LABEL_PENAWARAN.policyHolderName],
    [isi.typeCeding, LABEL_PENAWARAN.typeCeding],
    [isi.businessCode, LABEL_PENAWARAN.businessCode],
    [isi.dateReceived, LABEL_PENAWARAN.dateReceived],
    [isi.status, LABEL_PENAWARAN.status],
    [isi.description, LABEL_PENAWARAN.description],
  ]
  return periksa.filter(([nilai]) => nilai.trim() === '').map(([, label]) => label)
}

/**
 * Penyaring ketikan Sum Insured: hanya angka dan SATU titik desimal.
 *
 * ⛔ Bukan `type="number"`: uang tetap TEKS (ADR-U-0003) — kotak angka peramban
 * membulatkan nilai besar dan menerima `e`, `+`, `-`. Koma dibuang, bukan
 * dianggap desimal: server membaca titik sebagai pemisah desimal.
 */
export function saringAngkaDesimal(v: string): string {
  const bersih = v.replace(/[^0-9.]/g, '')
  const titik = bersih.indexOf('.')
  if (titik < 0) return bersih
  return bersih.slice(0, titik + 1) + bersih.slice(titik + 1).replace(/\./g, '')
}
