// Halaman awal modul - `Section/InboxProductName.xml` (PARITAS §1–§2).
//
// ⭐ Pintu masuk: `InboxProductName` tidak dibuka rule mana pun di korpus (dibuka portal); harness satu-satunya,
// `InwardProductName`, hanya dibuka tombol `Inward` b75368 yang `OTHER 1=2` (RALAT R11).
//
// Mode daftar (wadah b71246 `DATASHOW != 1`): tombol b71865 (label sel `End Period`, teks `Add`, tooltip
// `Add New Data` → `NewProductLife`), `pyGridPaginator` b72103 (10 baris/halaman, `pyRDLPageSize` b71371), grid
// RD `BrowseProduct_Life` (urut `.ID ASC` - server) dengan kolom `ID` · `Ceding` · `Treaty Number` · `Treaty Name`
// · `Create Operator` · `Last Updated Operator` dan tombol baris `View` b74798 (`SetProductName` → mode lihat).
// Mode form (`DATASHOW = 1`): `FormProduk`. Wadah grid b71574 ber-`IsFire` dengan `ALWAYS` → selalu tampil.

import { useCallback, useEffect, useRef, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambilDaftarProduk, ambilProduk, type Produk, type RingkasanProduk } from '../api'
import { UKURAN_HALAMAN_MPNL, jepitHalaman, potongHalaman, produkBaru } from '../bentuk'
import FormProduk from '../components/FormProduk'
import { GRID_MPNL, LAIN_MPNL, MENU_MPNL } from '../labels'

/** Form yang terbuka: halamannya dan `IsView` awal. */
interface FormTerbuka {
  produk: Produk
  lihat: boolean
  /** Kunci React - form baru setiap kali dibuka. */
  ke: number
}

export default function MasterProductNameLife() {
  const [daftar, setDaftar] = useState<RingkasanProduk[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [halaman, setHalaman] = useState(1)
  const [form, setForm] = useState<FormTerbuka | null>(null)
  const [galatBuka, setGalatBuka] = useState<unknown>(null)
  // Kunci pemasangan ulang form, dan nomor pembukaan terakhir (`Add`/`View`).
  const ke = useRef(0)

  const muat = useCallback(async () => {
    try {
      const d = await ambilDaftarProduk()
      setDaftar(d.daftar)
      setGalat(null)
      setHalaman((h) => jepitHalaman(h, d.daftar.length))
    } catch (e) {
      // ⛔ Galat DINYATAKAN, bukan menjadi daftar kosong.
      setGalat(e)
    }
  }, [])

  useEffect(() => {
    void muat()
  }, [muat])

  // Audit 02-10-2026: setiap pembukaan mendapat nomor baru SAAT diminta; jawaban `View` yang tiba sesudah
  // pembukaan lain (View berikutnya atau `Add`) diabaikan, jadi form tidak pernah menampilkan produk yang salah.
  function buka(produk: Produk, lihat: boolean, nomor = ++ke.current): void {
    if (nomor !== ke.current) return
    setForm({ produk, lihat, ke: nomor })
  }

  async function lihatProduk(id: string): Promise<void> {
    const nomor = ++ke.current
    setGalatBuka(null)
    try {
      // `View` b74798 → `SetProductName` (+ `SetProductNameInward`), 8 b2147 `IsView = true`.
      buka(await ambilProduk(id), true, nomor)
    } catch (e) {
      if (nomor === ke.current) setGalatBuka(e)
    }
  }

  function tutup(): void {
    // `Close` b58770 → `HideCreateLife`: `DATASHOW = 0`; refresh section membaca ulang grid.
    setForm(null)
    void muat()
  }

  const semua = daftar ?? []

  return (
    <section className="inbox mpnl">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{MENU_MPNL.kelompok}</h2>
      </header>

      {form !== null ? (
        <FormProduk key={form.ke} awal={form.produk} lihatAwal={form.lihat} onTutup={tutup} onTersimpan={tutup} />
      ) : (
        <>
          <div className="aksi-baris mpnl-aksi-grid">
            <span className="mpnl-label-sel">{GRID_MPNL.labelSelAdd}</span>
            <button type="button" className="btn btn--primary" title={GRID_MPNL.tooltipAdd} onClick={() => buka(produkBaru(), false)}>
              {GRID_MPNL.add}
            </button>
            {semua.length > 0 && (
              <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN_MPNL} total={semua.length} onPindah={setHalaman} />
            )}
          </div>
          {galatBuka !== null && <Gagal galat={galatBuka} />}
          {daftar === null && galat === null && <Memuat />}
          {galat !== null && <Gagal galat={galat} />}
          {daftar !== null && semua.length === 0 && <Kosong pesan={LAIN_MPNL.kosong} />}
          {semua.length > 0 && (
            <table className="inbox__tabel">
              <thead>
                <tr>
                  <th>{GRID_MPNL.kolomId}</th>
                  <th>{GRID_MPNL.kolomCeding}</th>
                  <th>{GRID_MPNL.kolomTreatyNumber}</th>
                  <th>{GRID_MPNL.kolomTreatyName}</th>
                  <th>{GRID_MPNL.kolomCreateOp}</th>
                  <th>{GRID_MPNL.kolomUpdateOp}</th>
                  <th className="table__actions" />
                </tr>
              </thead>
              <tbody>
                {potongHalaman(semua, halaman).map((r) => (
                  <tr key={r.id} className="inbox__baris">
                    <td>{r.id}</td>
                    <td>{r.ceding}</td>
                    <td>{r.treatyNumber}</td>
                    <td>{r.inwardName}</td>
                    <td>{r.createOp}</td>
                    <td>{r.updateOp}</td>
                    <td className="table__actions">
                      <button type="button" className="btn btn--ghost btn--sm" onClick={() => void lihatProduk(r.id)}>
                        {GRID_MPNL.view}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </>
      )}
    </section>
  )
}
